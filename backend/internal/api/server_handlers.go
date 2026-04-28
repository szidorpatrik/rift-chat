package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

const (
	serverNameMinLength = 3
	serverNameMaxLength = 100
	serverDescMaxLength = 255
)

// Allows any printable character, including Unicode/emojis.
// Length: 3 to 100 characters.
var serverNameRegex = regexp.MustCompile(fmt.Sprintf(`^[^\x00-\x1F\x7F]{%d,%d}$`, serverNameMinLength, serverNameMaxLength))

// Allows any printable character, including Unicode/emojis.
// Length: 0 to 255 characters.
var serverDescRegex = regexp.MustCompile(fmt.Sprintf(`^[^\x00-\x1F\x7F]{0,%d}$`, serverDescMaxLength))

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
		res[i] = ServerResponse{
			ID:          server.ID.String(),
			Name:        server.Name,
			Description: server.Description.String,
			IsPublic:    server.IsPublic,
			CreatedAt:   server.CreatedAt.Time,
			UpdatedAt:   server.UpdatedAt.Time,
		}
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
		res[i] = ServerResponse{
			ID:          server.ID.String(),
			Name:        server.Name,
			Description: server.Description.String,
			IsPublic:    server.IsPublic,
			CreatedAt:   server.CreatedAt.Time,
		}
	}

	RespondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) CreateServerHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userContextKey).(uuid.UUID)

	var req ServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.LogWarn(r, "invalid request body", userId, EventServerCreateFailed)
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

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

	var serverDesc sql.NullString

	if req.Description != "" {
		serverDesc = sql.NullString{
			String: req.Description,
			Valid:  true,
		}
	} else {
		serverDesc = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	createServerParams := database.CreateServerParams{
		ID:          serverId,
		OwnerID:     userId,
		Name:        req.Name,
		Description: serverDesc,
		IsPublic:    req.IsPublic,
	}

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

	joinServerParams := database.JoinServerParams{
		ServerID: serverId,
		UserID:   userId,
	}
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

	res := ServerResponse{
		ID:          createdServer.ID.String(),
		Name:        createdServer.Name,
		Description: createdServer.Description.String,
		IsPublic:    createdServer.IsPublic,
		CreatedAt:   createdServer.CreatedAt.Time,
		UpdatedAt:   createdServer.UpdatedAt.Time,
	}
	RespondWithJSON(w, http.StatusCreated, res)
}
