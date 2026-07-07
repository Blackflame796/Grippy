package errorcode

const (
	Unauthorized          ErrorCode = "UNAUTHORIZED"
	AccessDenied          ErrorCode = "ACCESS_DENIED"
	AuthSessionNotFound   ErrorCode = "AUTH_SESSION_NOT_FOUND"
	InvalidCredentials    ErrorCode = "INVALID_CREDENTIALS"
	UsernameAlreadyExists ErrorCode = "USERNAME_ALREADY_EXISTS"
	UserAlreadyExists     ErrorCode = "USER_ALREADY_EXISTS"
	UserAlreadySignedIn   ErrorCode = "USER_ALREADY_SIGNED_IN"
	RefreshTokenMissing   ErrorCode = "REFRESH_TOKEN_MISSING"
	RefreshTokenInvalid   ErrorCode = "REFRESH_TOKEN_INVALID"
)
