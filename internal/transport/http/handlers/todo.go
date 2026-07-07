package handlers

import (
	entity "Grippy/internal/domain"
	"Grippy/internal/repository"
	"Grippy/internal/transport/http/middlewares"
	"Grippy/internal/transport/http/response/errorcode"
	"Grippy/internal/transport/http/router"
	todo_usecase "Grippy/internal/usecase/todo"
	"Grippy/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type ToDoHandler struct {
	useCase *todo_usecase.ToDoUseCase
}

func NewToDoHandler(repo *repository.ToDoRepository) *ToDoHandler {
	useCase := todo_usecase.NewToDoUseCase(repo)
	return &ToDoHandler{
		useCase: useCase,
	}
}

func (h *ToDoHandler) RegisterRoutes(r *router.Router) {
	r.Post("/create", h.CreateTodo)
	r.Put("/update", h.UpdateTodo)
	r.Delete("/delete", h.DeleteTodo)
	r.Get("/get", h.GetToDoByID)
	r.Get("/get_all", h.GetAllTodos)
}

func (h *ToDoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserKey).(*entity.Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Code:    errorcode.Unauthorized,
			Message: "Unauthorized",
		})
		return
	}

	var req todo_usecase.CreateToDoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.InvalidRequestBody,
			Message: "Invalid request body",
		})
		return
	}

	responseToDo, err := h.useCase.Create(r.Context(), claims.ID, req)
	if err != nil {
		logger.Log.Errorf("Error creating todo: %v", err)

		switch {
		case errors.Is(err, todo_usecase.ErrInvalidInput):
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    errorcode.TodoInvalidInput,
				Message: err.Error(),
			})
		case errors.Is(err, todo_usecase.ErrDuplicateTitle):
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Code:    errorcode.TodoDuplicateTitle,
				Message: err.Error(),
			})
		default:
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Code:    errorcode.InternalServerError,
				Message: "Failed to create todo",
			})
		}
		return
	}

	writeJSON(w, http.StatusCreated, responseToDo)
}

func (h *ToDoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var req todo_usecase.UpdateToDoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.InvalidRequestBody,
			Message: "Invalid request body",
		})
		return
	}

	claims, ok := r.Context().Value(middlewares.UserKey).(*entity.Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Code:    errorcode.Unauthorized,
			Message: "Unauthorized",
		})
		return
	}

	responseToDo, err := h.useCase.Update(r.Context(), req, claims)
	if err != nil {
		logger.Log.Errorf("Error updating todo: %v", err)

		switch {
		case errors.Is(err, todo_usecase.ErrTodoNotFound):
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    errorcode.TodoNotFound,
				Message: "Todo not found",
			})
		case errors.Is(err, todo_usecase.ErrForbidden):
			writeJSON(w, http.StatusForbidden, ErrorResponse{
				Code:    errorcode.TodoForbiddenAction,
				Message: "You cannot update this todo",
			})
		case errors.Is(err, todo_usecase.ErrInvalidInput):
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    errorcode.TodoInvalidInput,
				Message: err.Error(),
			})
		default:
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Code:    errorcode.InternalServerError,
				Message: "Failed to update todo",
			})
		}
		return
	}

	writeJSON(w, http.StatusOK, responseToDo)
}

func (h *ToDoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var id uuid.UUID

	if idStr != "" {
		var err error
		id, err = uuid.Parse(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    errorcode.TodoInvalidInput,
				Message: "Invalid ID format in query parameter",
			})
			return
		}
	} else {
		var req todo_usecase.DeleteToDoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    errorcode.InvalidRequestBody,
				Message: "ID required in query or body",
			})
			return
		}
		id = req.ID
	}

	if id == uuid.Nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.TodoInvalidInput,
			Message: "ID is required",
		})
		return
	}

	claims, ok := r.Context().Value(middlewares.UserKey).(*entity.Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Code:    errorcode.Unauthorized,
			Message: "Unauthorized",
		})
		return
	}

	deletedID, err := h.useCase.Delete(r.Context(), id, claims)
	if err != nil {
		logger.Log.Errorf("Error deleting todo: %v", err)

		switch {
		case errors.Is(err, todo_usecase.ErrTodoNotFound):
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    errorcode.TodoNotFound,
				Message: "Todo not found",
			})
		case errors.Is(err, todo_usecase.ErrForbidden):
			writeJSON(w, http.StatusForbidden, ErrorResponse{
				Code:    errorcode.TodoForbiddenAction,
				Message: "You cannot delete this todo",
			})
		default:
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Code:    errorcode.InternalServerError,
				Message: "Failed to delete todo",
			})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Todo deleted with id: %s", deletedID),
	})
}

func (h *ToDoHandler) GetToDoByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.TodoInvalidInput,
			Message: "ID parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    errorcode.TodoInvalidInput,
			Message: "Invalid UUID format",
		})
		return
	}

	claims, ok := r.Context().Value(middlewares.UserKey).(*entity.Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Code:    errorcode.Unauthorized,
			Message: "Unauthorized",
		})
		return
	}

	todo, err := h.useCase.GetToDoByID(r.Context(), id, claims)
	if err != nil {
		logger.Log.Errorf("Error getting todo: %v", err)

		if errors.Is(err, todo_usecase.ErrTodoNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    errorcode.TodoNotFound,
				Message: "Todo not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    errorcode.InternalServerError,
			Message: "Failed to get todo",
		})
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *ToDoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		} else {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    errorcode.TodoInvalidInput,
				Message: "Invalid limit value",
			})
			return
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		} else {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    errorcode.TodoInvalidInput,
				Message: "Invalid offset value",
			})
			return
		}
	}

	claims, ok := r.Context().Value(middlewares.UserKey).(*entity.Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Code:    errorcode.Unauthorized,
			Message: "Unauthorized",
		})
		return
	}

	todos, err := h.useCase.GetAllToDo(r.Context(), claims, limit, offset)
	if err != nil {
		logger.Log.Errorf("Error getting todos: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    errorcode.InternalServerError,
			Message: "Failed to fetch todos",
		})
		return
	}

	writeJSON(w, http.StatusOK, todos)
}
