package api

import "net/http"

func (cfg *ApiConfig) registerSystemRoutes(mux *http.ServeMux) {
	mux.HandleFunc(newEndpoint("GET", "/docs"), cfg.SwaggerDocsHandler)
	mux.HandleFunc(newEndpoint("GET", "/docs/openapi.json"), cfg.SwaggerSpecHandler)
	mux.HandleFunc(newEndpoint("GET", "/healthz"), cfg.HealthzHandler)
}

func (cfg *ApiConfig) registerAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc(newEndpoint("POST", "/users"), cfg.CreateUserHandler)
	mux.HandleFunc(newEndpoint("POST", "/login"), cfg.LoginHandler)
	mux.HandleFunc(newEndpoint("POST", "/refresh"), cfg.RefreshTokenHandler)
	mux.HandleFunc(newEndpoint("POST", "/revoke"), cfg.RevokeTokenHandler)

	sessionCheck := http.HandlerFunc(cfg.SessionCheckHandler)
	mux.Handle(newEndpoint("GET", "/session"), cfg.RequireAuth(sessionCheck))
}

func (cfg *ApiConfig) registerServerRoutes(mux *http.ServeMux) {
	// Public
	mux.HandleFunc(newEndpoint("GET", "/servers/public"), cfg.GetPublicServersHandler)

	// Protected
	getServers := http.HandlerFunc(cfg.GetServersHandler)
	getServerMembers := http.HandlerFunc(cfg.GetServerMembersHandler)
	createServer := http.HandlerFunc(cfg.CreateServerHandler)
	joinServer := http.HandlerFunc(cfg.JoinServerHandler)
	leaveServer := http.HandlerFunc(cfg.LeaveServerHandler)

	mux.Handle(newEndpoint("GET", "/servers"), cfg.RequireAuth(getServers))
	mux.Handle(newEndpoint("POST", "/servers"), cfg.RequireAuth(createServer))
	mux.Handle(newEndpoint("POST", "/servers/join"), cfg.RequireAuth(joinServer))
	
	// Membership required
	mux.Handle(newEndpoint("GET", "/servers/{server_id}/members"), cfg.RequireAuth(cfg.ServerMembershipMiddleware(getServerMembers)))
	mux.Handle(newEndpoint("POST", "/servers/{server_id}/leave"), cfg.RequireAuth(cfg.ServerMembershipMiddleware(leaveServer)))

	//TODO: invite code based join
	//mux.Handle(newEndpoint("POST", "/servers/join/{invite_code}"), cfg.RequireAuth(joinServerByInviteCode))
}

func (cfg *ApiConfig) registerChannelRoutes(mux *http.ServeMux) {
	getChannels := http.HandlerFunc(cfg.GetServerChannelsHandler)
	createChannel := http.HandlerFunc(cfg.CreateChannelHandler)
	deleteChannel := http.HandlerFunc(cfg.DeleteChannelHandler)

	protectedGetChannels := cfg.RequireAuth(cfg.ServerMembershipMiddleware(getChannels))
	protectedCreateChannel := cfg.RequireAuth(cfg.ServerMembershipMiddleware(createChannel))
	protectedDeleteChannel := cfg.RequireAuth(cfg.ServerMembershipMiddleware(deleteChannel))

	mux.Handle(newEndpoint("GET", "/servers/{server_id}/channels"), protectedGetChannels)
	mux.Handle(newEndpoint("POST", "/servers/{server_id}/channels"), protectedCreateChannel)
	mux.Handle(newEndpoint("DELETE", "/servers/{server_id}/channels/{channel_id}"), protectedDeleteChannel)
}
