package handlers

import (
	"Grippy/internal/transport/http/response/errorcode"
	"Grippy/internal/transport/http/router"
	auth_usecase "Grippy/internal/usecase/auth"
	"Grippy/pkg/logger"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type AuthHandler struct {
	useCase *auth_usecase.AuthUseCase
}

type ErrorResponse struct {
	Code    errorcode.ErrorCode `json:"code"`
	Message string              `json:"message"`
}

func NewAuthHandler(uc *auth_usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: uc,
	}
}

func (h *AuthHandler) RegisterRoutes(r *router.Router) {
	r.Post("/sign_up", h.SignUp)
	r.Post("/sign_in", h.SignIn)
	r.Post("/refresh_token", h.Refresh)
	r.Post("/logout", h.Logout)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil && cookie != nil && cookie.Value != "" {
		err := h.useCase.ValidateRefreshToken(r.Context(), cookie.Value)

		if err == nil {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Code:    errorcode.UserAlreadySignedIn,
				Message: "Already signed in. Please logout first to create a new account.",
			})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/auth",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	var req auth_usecase.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.InvalidRequestBody,
			Message: "Invalid request body",
		})
		return
	}

	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.InvalidRequestBody,
			Message: "Email and password are required",
		})
		return
	}

	user, err := h.useCase.Register(r.Context(), req)
	if err != nil {
		logger.Log.Errorf("Registration error: %v", err)

		if errors.Is(err, auth_usecase.ErrUserAlreadyExists) {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Code:    errorcode.UserAlreadyExists,
				Message: "User with this email already exists",
			})
			return
		}

		if errors.Is(err, auth_usecase.ErrUsernameAlreadyExists) {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Code:    errorcode.UsernameAlreadyExists,
				Message: "This username is already taken",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    errorcode.InternalServerError,
			Message: "Registration failed",
		})
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil && cookie != nil && cookie.Value != "" {
		err := h.useCase.ValidateRefreshToken(r.Context(), cookie.Value)

		if err == nil {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Code:    errorcode.UserAlreadySignedIn,
				Message: "Already signed in. Please logout first to create a new account.",
			})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/auth",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	var req auth_usecase.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.InvalidRequestBody,
			Message: "Invalid request body",
		})
		return
	}

	tokenPair, err := h.useCase.Login(r.Context(), req)
	if err != nil {
		logger.Log.Errorf("Login error: %v", err)

		if errors.Is(err, auth_usecase.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{
				Code:    errorcode.InvalidCredentials,
				Message: "Invalid email or password",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    errorcode.InternalServerError,
			Message: "Authentication failed",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth",
	})

	writeJSON(w, http.StatusOK, map[string]string{"access_token": tokenPair.AccessToken})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Code:    errorcode.RefreshTokenMissing,
			Message: "Missing refresh token cookie",
		})
		return
	}

	tokenPair, err := h.useCase.Refresh(r.Context(), cookie.Value, "user")
	if err != nil {
		logger.Log.Errorf("Refresh error: %v", err)

		if errors.Is(err, auth_usecase.ErrSessionNotFound) {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{
				Code:    errorcode.AuthSessionNotFound,
				Message: "Session expired or not found",
			})
			return
		}

		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Code:    errorcode.RefreshTokenInvalid,
			Message: "Refresh failed",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth",
	})

	writeJSON(w, http.StatusOK, map[string]string{"access_token": tokenPair.AccessToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Already logged out",
		})
		return
	}

	err = h.useCase.Logout(r.Context(), cookie.Value)
	if err != nil {
		logger.Log.Errorf("Logout use case error: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
