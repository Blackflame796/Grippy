package handlers

import (
	"ToDoApp/pkg/logger"
	"ToDoApp/pkg/repository"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type ToDoHandler struct {
	repo *repository.ToDoRepository
}

func NewToDoHandler(repo *repository.ToDoRepository) *ToDoHandler {
	return &ToDoHandler{
		repo: repo,
	}
}

type CreateToDoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateToDoRequest struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
}

type DeleteToDoRequest struct {
	ID uuid.UUID `json:"id"`
}

func (h *ToDoHandler) RegisterRoutes(r *Router) {
	r.Post("/todos/create", h.CreateTodo)
	r.Put("/todos/update", h.UpdateToDo)
	r.Delete("/todos/delete", h.DeleteTodo)
	r.Get("/todos/get", h.GetToDoByID)
	r.Get("/todos/get_all", h.GetAllTodos)
}

func (h *ToDoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateToDoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdToDo, err := h.repo.Add(r.Context(), req.Title, req.Description)
	if err != nil {
		logger.Log.Errorf("Error creating todo: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error creating todo"}`))
		return
	}
	jsonData, err := json.Marshal(createdToDo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
	logger.Log.Infof("Todo created with id: %s", createdToDo.ID)
}

func (h *ToDoHandler) UpdateToDo(w http.ResponseWriter, r *http.Request) {
	var req UpdateToDoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedToDo, err := h.repo.Update(r.Context(), req.ID, req.Title, req.Description, req.IsCompleted)
	if err != nil {
		logger.Log.Errorf("Error updating todo: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error updating todo"}`))
		return
	}
	jsonData, err := json.Marshal(updatedToDo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	logger.Log.Infof("Todo updated with id: %s", updatedToDo.ID)
}

func (h *ToDoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var id uuid.UUID
	var err error

	// First, try to get ID from query parameters
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err = uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid ID in query parameter", http.StatusBadRequest)
			return
		}
	} else {
		var req DeleteToDoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "ID not found in query parameters or request body", http.StatusBadRequest)
			return
		}
		id = req.ID

		if id == uuid.Nil {
			http.Error(w, "ID is required in query parameters or request body", http.StatusBadRequest)
			return
		}
	}

	deletedID, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		logger.Log.Errorf("Error deleting todo: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error deleting todo"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Log.Infof("Todo deleted with id: %s", deletedID)
}

func (h *ToDoHandler) GetToDoByID(w http.ResponseWriter, r *http.Request) {
	var id uuid.UUID
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		var err error
		id, err = uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid ID in query parameter", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "ID not found in query parameter", http.StatusBadRequest)
		return
	}

	todo, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		logger.Log.Errorf("Error getting todo: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error getting todo"}`))
		return
	}

	jsonData, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	logger.Log.Infof("Todo retrieved with id: %d", id)
}

func (h *ToDoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repo.GetAll(r.Context())
	if err != nil {
		logger.Log.Errorf("Error getting todos: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error getting todos"}`))
		return
	}

	jsonData, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	logger.Log.Infof("Todos retrieved")
}
