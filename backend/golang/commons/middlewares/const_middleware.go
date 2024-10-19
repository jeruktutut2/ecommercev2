package middlewares

type StringCustomType string

const (
	RequestIdKey     StringCustomType = "requestId"
	IdKey            StringCustomType = "id"
	PermissionKey    StringCustomType = "permission"
	UsernameKey      StringCustomType = "username"
	XRefreshTokenKey StringCustomType = "xRefreshToken"
	TokenIdKey       StringCustomType = "tokenId"
	SessionIdKey     StringCustomType = "sessionId"
)
