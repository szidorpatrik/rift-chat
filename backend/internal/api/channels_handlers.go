package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

var channelTypes = map[int32]string{
	0: "text",
	1: "voice",
	2: "category",
}

func (cfg *ApiConfig) GetServerChannelsHandler(w http.ResponseWriter, r *http.Request) {
	serverId := r.Context().Value(serverContextKey).(uuid.UUID)

	channels, err := cfg.DbQueries.GetChannelsForServer(r.Context(), serverId)
	if err != nil {
		cfg.LogError(r, err, EventDbError, "server_id", serverId)
		RespondWithError(w, http.StatusInternalServerError, "failed to get channels")
		return
	}

	cfg.LogInfo(r, "server channels retrieved successfully", serverId, EventChannelsRetrieved)

	res := make([]ChannelResponse, len(channels))
	for i, channel := range channels {
		res[i] = getChannelResponse(channel)
	}
	RespondWithJSON(w, http.StatusOK, res)
}

// TODO: Channels
func (cfg *ApiConfig) CreateChannelHandler(w http.ResponseWriter, r *http.Request) {
	serverId := r.Context().Value(serverContextKey).(uuid.UUID)

	var req ChannelCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.LogWarn(r, "invalid request body", serverId, EventChannelCreateFailed)
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate channel type
	if _, exists := channelTypes[req.Type]; !exists {
		cfg.LogWarn(r, "invalid channel type", serverId, EventChannelCreateFailed)
		RespondWithError(w, http.StatusBadRequest, "invalid channel type")
		return
	}

	// Validate channel name
	req.Name = sanitizeInput(req.Name, " ")
	if !channelNameRegex.MatchString(req.Name) {
		cfg.LogWarn(r, "invalid channel name", serverId, EventChannelCreateFailed)
		RespondWithError(w, http.StatusBadRequest, "invalid channel name")
		return
	}

	if req.Type == 0 {
		req.Name = whitespace.ReplaceAllString(req.Name, "-")
		req.Name = strings.ToLower(req.Name)
	}

	channelId, err := uuid.NewV7()
	if err != nil {
		cfg.LogError(r, err, EventChannelCreateFailed, "server_id", serverId)
		RespondWithError(w, http.StatusInternalServerError, "failed to create channel")
		return
	}

	params := database.CreateChannelParams{
		ID:       channelId,
		ServerID: serverId,
		Name:     req.Name,
		Type:     req.Type,
		Position: req.Position,
	}
	channel, err := cfg.DbQueries.CreateChannel(r.Context(), params)
	if err != nil {
		cfg.LogError(r, err, EventChannelCreateFailed, "server_id", serverId)
		RespondWithError(w, http.StatusInternalServerError, "failed to create channel")
		return
	}

	cfg.LogInfo(r, "channel created successfully", channelId, EventChannelCreated, "server_id", serverId)
	res := getChannelResponse(channel)
	RespondWithJSON(w, http.StatusCreated, res)
}

func (cfg *ApiConfig) DeleteChannelHandler(w http.ResponseWriter, r *http.Request) {
	serverId := r.Context().Value(serverContextKey).(uuid.UUID)
	channelIdStr := r.PathValue("channel_id")
	if channelIdStr == "" {
		cfg.LogWarn(r, "missing channel_id path parameter", serverId, EventMissingParameter)
		RespondWithError(w, http.StatusBadRequest, "channel_id path parameter is required")
		return
	}

	channelId, err := uuid.Parse(channelIdStr)
	if err != nil {
		cfg.LogWarn(r, "invalid channel ID", serverId, EventChannelDeleteFailed)
		RespondWithError(w, http.StatusBadRequest, "invalid channel ID")
		return
	}

	params := database.DeleteChannelByIdParams{
		ID:       channelId,
		ServerID: serverId,
	}

	if err := cfg.DbQueries.DeleteChannelById(r.Context(), params); err != nil {
		cfg.LogError(r, err, EventChannelDeleteFailed, "channel_id", channelId)
		RespondWithError(w, http.StatusInternalServerError, "failed to delete channel")
		return
	}

	cfg.LogInfo(r, "channel deleted successfully", channelId, EventChannelDeleted)
	RespondNoContent(w)
}

func getChannelResponse(channel database.Channel) ChannelResponse {
	return ChannelResponse{
		ID:       channel.ID,
		ServerID: channel.ServerID,
		Name:     channel.Name,
		Type:     channel.Type,
		Position: channel.Position,
	}
}
