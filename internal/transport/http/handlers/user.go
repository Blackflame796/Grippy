package handlers

import (
	entity "Grippy/internal/domain"
	"Grippy/internal/transport/http/middlewares"
	"Grippy/internal/transport/http/response/errorcode"
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
	r.Put("/update_info", h.UpdateInfo)
	r.Put("/upload-avatar", h.UploadAvatar)
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

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserKey).(*entity.Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Code:    errorcode.Unauthorized,
			Message: "Unauthorized",
		})
		return
	}

	const maxUploadSize = 10 * 1024 * 1024
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.UserAvatarTooLarge,
			Message: "File too large or invalid multipart format",
		})
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.UserAvatarFormKeyMissing,
			Message: "Missing 'avatar' key in form-data",
		})
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Code: errorcode.FailedReadFileContents, Message: "Failed to read file contents"})
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Message: "Failed to process file"})
		return
	}
	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.UserInvalidAvatarFormat,
			Message: "Only JPEG and PNG formats are allowed",
		})
		return
	}

	imageUrl, err := h.useCase.UploadAvatar(r.Context(), claims.ID, file, fileHeader.Filename, contentType)
	if err != nil {
		logger.Log.Errorf("Error uploading avatar: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    errorcode.InternalServerError,
			Message: "Failed to upload avatar",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message":    "Avatar uploaded successfully",
		"avatar_url": imageUrl,
	})
}
