package api

type EventType int

var eventName = map[EventType]string{
	EventLogin:                       "user_login",
	EventLoginFailed:                 "user_login_failed",
	EventRegister:                    "user_register",
	EventRegisterFailed:              "user_register_failed",
	EventAccessTokenCreated:          "access_token_created",
	EventAccessTokenRefreshed:        "access_token_refreshed",
	EventAccessTokenFailedSign:       "access_token_failed_sign",
	EventAccessTokenInvalid:          "access_token_invalid",
	EventBearerTokenInvalid:          "bearer_token_invalid",
	EventRefreshTokenCreated:         "refresh_token_created",
	EventRefreshTokenCreatedFailed:   "refresh_token_created_failed",
	EventRefreshTokenRevoked:         "refresh_token_revoked",
	EventRefreshTokenMissing:         "refresh_token_missing",
	EventRefreshTokenInvalid:         "refresh_token_invalid",
	EventRefreshTokenExpired:         "refresh_token_expired",
	EventRefreshTokenRevokedDetected: "refresh_token_revoked_detected",
	EventRefreshTokenNotFound:        "refresh_token_not_found",
	EventHTTPRequest:                 "http_request",
	EventDbError:                     "db_error",
	EventSessionCheckSuccess:         "session_check_success",
	EventSessionCheckFailed:          "session_check_failed",
	EventServerCreateFailed:          "server_create_failed",
	EventServerCreated:               "server_created",
	EventServerJoinSuccess:           "server_join_success",
	EventServerJoinFailed:            "server_join_failed",
}

const (
	EventLogin EventType = iota
	EventLoginFailed
	EventRegister
	EventRegisterFailed
	EventAccessTokenCreated
	EventAccessTokenRefreshed
	EventAccessTokenFailedSign
	EventAccessTokenInvalid
	EventBearerTokenInvalid
	EventRefreshTokenCreated
	EventRefreshTokenCreatedFailed
	EventRefreshTokenRevoked
	EventRefreshTokenMissing
	EventRefreshTokenInvalid
	EventRefreshTokenExpired
	EventRefreshTokenRevokedDetected
	EventRefreshTokenNotFound
	EventHTTPRequest
	EventDbError
	EventSessionCheckSuccess
	EventSessionCheckFailed
	EventServerCreateFailed
	EventServerCreated
	EventServerJoinSuccess
	EventServerJoinFailed
)

func (e EventType) String() string {
	return eventName[e]
}
