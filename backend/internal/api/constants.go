package api

import "time"

const (
	apiRoot              = "/api"
	apiVersion           = 1
	logBatcherSize       = 100
	logBatcherFlushTimer = 5 * time.Second

	userContextKey   = "user_id"
	serverContextKey = "server_id"

	accessTokenExpiry  = time.Hour
	refreshTokenExpiry = time.Hour * 24 * 30 // 30 days
	cookieAccessToken  = "access_token"
	cookieRefreshToken = "refresh_token"
)
