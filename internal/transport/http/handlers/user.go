package handlers

import (
	entity "Grippy/internal/domain"
	"Grippy/internal/transport/http/middlewares"
	"Grippy/internal/transport/http/router"
	user_usecase "Grippy/internal/usecase/user"
	"Grippy/pkg/logger"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	useCase *user_usecase.UserUseCase
}

func NewUserHandler(uc *user_usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		useCase: uc,
	}
}

func (h *UserHandler) RegisterRoutes(r *router.Router) {
	r.Post("/update_info", h.UpdateInfo)
}

func (h *UserHandler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserKey).(*entity.Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	var req user_usecase.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	user, err := h.useCase.Update(r.Context(), req, claims)
	if err != nil {
		logger.Log.Errorf("Update error: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Update failed",
		})
		return
	}

	writeJSON(w, http.StatusOK, user)
	logger.Log.Infof("User with id = %s updated successfully", user.ID.String())
}
