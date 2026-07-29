package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

func (cfg *ApiConfig) GetServersHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userContextKey).(uuid.UUID)

	servers, err := cfg.DbQueries.GetServersByUserId(r.Context(), userId)
	if err != nil {
		cfg.LogError(r, err, EventDbError)
		RespondWithError(w, http.StatusInternalServerError, "could not retrieve servers")
		return
	}

	res := make([]ServerResponse, len(servers))
	for i, server := range servers {
		res[i] = getServerResponse(server)
	}

	RespondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) GetPublicServersHandler(w http.ResponseWriter, r *http.Request) {
	servers, err := cfg.DbQueries.GetPublicServers(r.Context())
	if err != nil {
		cfg.LogError(r, err, EventDbError)
		RespondWithError(w, http.StatusInternalServerError, "could not retrieve public servers")
		return
	}

	res := make([]ServerResponse, len(servers))
	for i, server := range servers {
		res[i] = getServerResponse(server)
	}

	RespondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) CreateServerHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userContextKey).(uuid.UUID)

	req, err := decodeJSONBody[ServerRequest](r)
	if err != nil {
		cfg.LogWarn(r, "invalid request body", userId, EventServerCreateFailed)
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Name = sanitizeInput(req.Name, " ")
	req.Description = sanitizeInput(req.Description, " ")

	if !serverNameRegex.MatchString(req.Name) {
		cfg.LogWarn(r, "invalid server name", userId, EventServerCreateFailed)
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("server name must be %d-%d characters", serverNameMinLength, serverNameMaxLength))
		return
	}

	if !serverDescRegex.MatchString(req.Description) {
		cfg.LogWarn(r, "invalid server description", userId, EventServerCreateFailed)
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("server description must be at most %d characters", serverDescMaxLength))
		return
	}

	serverId, err := uuid.NewV7()
	if err != nil {
		cfg.LogError(r, err, EventServerCreateFailed)
		RespondWithError(w, http.StatusInternalServerError, "failed to generate server ID")
		return
	}

	createServerParams := buildCreateServerParams(req, userId, serverId)

	tx, err := cfg.DB.BeginTx(r.Context(), nil)
	if err != nil {
		cfg.LogError(r, err, EventDbError)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	defer tx.Rollback()
	qtx := cfg.DbQueries.WithTx(tx)

	createdServer, err := qtx.CreateServer(r.Context(), createServerParams)
	if err != nil {
		cfg.LogError(r, err, EventServerCreateFailed)
		RespondWithError(w, http.StatusInternalServerError, "failed to create server in database")
		return
	}

	joinServerParams := buildJoinServerParams(serverId, userId)
	if err := qtx.JoinServer(r.Context(), joinServerParams); err != nil {
		cfg.LogError(r, err, EventServerJoinFailed)
		RespondWithError(w, http.StatusInternalServerError, "failed to join server")
		return
	}

	if err := tx.Commit(); err != nil {
		cfg.LogError(r, err, EventDbError)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	cfg.LogInfo(r, "server created successfully", createdServer.ID, EventServerCreated)
	cfg.LogInfo(r, "user joined server successfully", userId, EventServerJoinSuccess, "server_id", createdServer.ID)

	res := getServerResponse(createdServer)
	RespondWithJSON(w, http.StatusCreated, res)
}

func (cfg *ApiConfig) JoinServerHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userContextKey).(uuid.UUID)
	req, err := decodeJSONBody[ServerJoinRequest](r)
	if err != nil {
		cfg.LogWarn(r, "invalid request body", userId, EventServerJoinFailed)
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	serverId := req.ServerId
	if serverId == uuid.Nil {
		cfg.LogWarn(r, "invalid server ID", userId, EventServerJoinFailed)
		RespondWithError(w, http.StatusBadRequest, "invalid server ID")
		return
	}

	server, err := cfg.DbQueries.GetServerById(r.Context(), serverId)
	if err != nil {
		cfg.LogWarn(r, "server not found", userId, EventServerJoinFailed, "server_id", serverId)
		RespondWithError(w, http.StatusNotFound, "server not found")
		return
	}

	if !server.IsPublic {
		cfg.LogWarn(r, "attempt to join private server", userId, EventServerJoinFailed, "server_id", serverId)
		RespondWithError(w, http.StatusForbidden, "you cannot join a private server without an invite code")
		return
	}

	params := buildJoinServerParams(serverId, userId)

	if err := cfg.DbQueries.JoinServer(r.Context(), params); err != nil {

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			cfg.LogWarn(r, "user already a member of the server", userId, EventServerJoinFailed, "server_id", serverId)
			RespondWithError(w, http.StatusConflict, "you are already a member of this server")
			return
		}
		cfg.LogError(r, err, EventServerJoinFailed)
		RespondWithError(w, http.StatusInternalServerError, "failed to join server")
		return
	}

	cfg.LogInfo(r, "server join successful", userId, EventServerJoinSuccess, "server_id", serverId)
	res := getServerResponse(server)
	RespondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) GetServerMembersHandler(w http.ResponseWriter, r *http.Request) {
	serverId := r.Context().Value(serverContextKey).(uuid.UUID)
	userId := r.Context().Value(userContextKey).(uuid.UUID)

	members, err := cfg.DbQueries.GetServerMembers(r.Context(), serverId)
	if err != nil {
		cfg.LogError(r, err, EventDbError)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	cfg.LogInfo(r, "server members listed", userId, ServerMembersListed)

	res := make([]MembersResponse, len(members))
	for i, member := range members {
		res[i] = getMembersResponse(member)
	}
	RespondWithJSON(w, http.StatusOK, res)
}

//TODO: invite codes, leaving servers, deleting servers, kicking members, etc.

func getServerResponse(server database.Server) ServerResponse {
	return ServerResponse{
		ID:            server.ID.String(),
		Name:          server.Name,
		OwnerId:       server.OwnerID.String(),
		Description:   server.Description.String,
		BackgroundUrl: server.BackgroundUrl.String,
		IconUrl:       server.IconUrl.String,
		IsPublic:      server.IsPublic,
		CreatedAt:     server.CreatedAt.Time,
		UpdatedAt:     server.UpdatedAt.Time,
	}
}

func getMembersResponse(member database.GetServerMembersRow) MembersResponse {
	return MembersResponse{
		ID:          member.UserID,
		Username:    member.Username,
		Status:      member.Status.String,
		AvatarUrl:   member.AvatarUrl.String,
		Description: member.Description.String,
	}
}

func buildCreateServerParams(req ServerRequest, userId, serverId uuid.UUID) database.CreateServerParams {
	var serverDesc sql.NullString
	if req.Description != "" {
		serverDesc = sql.NullString{String: sanitizeInput(req.Description, " "), Valid: true}
	}

	return database.CreateServerParams{
		ID:          serverId,
		OwnerID:     userId,
		Name:        sanitizeInput(req.Name, " "),
		Description: serverDesc,
		IsPublic:    req.IsPublic,
	}
}

func buildJoinServerParams(serverId, userId uuid.UUID) database.JoinServerParams {
	return database.JoinServerParams{
		ServerID: serverId,
		UserID:   userId,
	}
}
