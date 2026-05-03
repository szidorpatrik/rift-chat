package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/quic-go/webtransport-go"
)

func (cfg *ApiConfig) WebTransportHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userContextKey).(uuid.UUID)

	session, err := cfg.WTServer.Upgrade(w, r)
	if err != nil {
		cfg.LogError(r, err, EventWebTransportRequest, "user_id", userId)
		RespondWithError(w, http.StatusInternalServerError, "failed to upgrade connection")
		return
	}

	cfg.LogInfo(r, "user connected via WebTransport", userId, EventWebTransportRequest)
	go cfg.handleWTSession(session, userId)
}

func (cfg *ApiConfig) handleWTSession(session *webtransport.Session, userId uuid.UUID) {
	// Ensure the connection is cleaned up when this function exits
	defer session.CloseWithError(0, "client disconnected")

	ctx := session.Context()

	// 1. Open a one-way stream from the Server to the Client (Downlink)
	downlink, err := session.OpenUniStream()
	if err != nil {
		cfg.LogError(nil, err, EventWebTransportSession, "user_id", userId, "error", "failed to open downlink stream")
		return
	}
	defer downlink.Close()

	// 2. Subscribe to this specific user's event stream in Valkey
	// Every user gets their own topic, e.g., "user_events:019de585-..."
	userTopic := fmt.Sprintf("user_events:%s", userId.String())
	pubsub := cfg.Valkey.Subscribe(ctx, userTopic)
	defer pubsub.Close()

	cfg.LogInfo(nil, "WebTransport session established and subscribed", userId, EventWebTransportSession)

	// 3. Listen for events from Valkey and push them to the client
	valkeyChannel := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			// The client closed the browser or lost connection
			return

		case msg := <-valkeyChannel:
			// We received a message from Valkey!
			// We append a newline character so the Vue frontend can easily parse incoming JSON lines.
			payload := []byte(msg.Payload + "\n")

			if _, err := downlink.Write(payload); err != nil {
				cfg.LogError(nil, err, EventWebTransportSession, "user_id", userId, "error", "failed to write to WebTransport stream")
				return
			}
		}
	}
}
