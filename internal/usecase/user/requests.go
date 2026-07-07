package user_usecase

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
