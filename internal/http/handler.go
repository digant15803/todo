package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"errors"
	// "fmt"

	"todo/internal/todo"
)

type Handler struct {
	service *todo.Service
}

type CreateTodoRequest struct {
    Title string `json:"title"` // The title of the todo item
}

func NewHandler(service *todo.Service) *Handler {
	return &Handler{service: service}
}

// Helper function to handle HTTP errors
func handleError(w http.ResponseWriter, err error, statusCode int) {
    http.Error(w, err.Error(), statusCode)
}

// Registers the HTTP routes for the todo handler
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/todos", h.ServeHTTP)
	mux.HandleFunc("/todos/", h.ServeHTTP)
}

// Routes the incoming HTTP requests to the appropriate handler methods
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Normalize the path to handle "/todos" and "/todos/" equivalently
    path := strings.TrimSuffix(r.URL.Path, "/")

    // Route: "/todos" -> List todos
    if path == "/todos" {
        switch r.Method {
        case http.MethodGet:
            h.listTodos(w, r) // Handle GET /todos
        case http.MethodPost:
            h.createTodo(w, r) // Handle POST /todos
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
        return
    }

    // Route: "/todos/{id}" or "/todos/{id}/complete"
    if strings.HasPrefix(path, "/todos/") {
		// fmt.Println(path)
        h.handleTodoByID(w, r, strings.TrimPrefix(path, "/todos/"))
        return
    }

    // Default: Not Found
    http.NotFound(w, r)
}

// Handles requests for specific todo items by ID
func (h *Handler) handleTodoByID(w http.ResponseWriter, r *http.Request, path string) {
	parts := strings.Split(path, "/")

	// Ensure we have at least the ID part
	if len(parts) < 1 || parts[0] == "" {
		handleError(w, errors.New("missing todo id"), http.StatusBadRequest)
		return
	}

	id := parts[0]

	// fmt.Println("Handling todo with ID:", id, "Path parts:", parts)

	// POST /todos/{id}/complete
	if len(parts) == 2 && parts[1] == "complete" {
		if r.Method != http.MethodPost {
			handleError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
			return
		}
		h.completeTodo(w, id)
		return
	}

	// DELETE /todos/{id}
	if len(parts) == 1 && r.Method == http.MethodDelete {
        h.deleteTodo(w, id)
        return
    }

    // Default: Method Not Allowed
    handleError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
}

// Creates a new todo item
func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var request CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, errors.New("invalid JSON body"), http.StatusBadRequest)
		return
	}

	todoItem, err := h.service.Create(request.Title)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todoItem)
}

// Deletes a todo item by ID
func (h *Handler) deleteTodo(w http.ResponseWriter, id string) {
    if err := h.service.Delete(id); err != nil {
        if err == todo.ErrTodoNotFound {
            handleError(w, err, http.StatusNotFound)
            return
        }
        handleError(w, err, http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// Lists all todo items
func (h *Handler) listTodos(w http.ResponseWriter, r *http.Request) {
	todos := h.service.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Marks a todo item as completed by ID
func (h *Handler) completeTodo(w http.ResponseWriter, id string) {
	updatedTodo, err := h.service.Complete(id)
	if err != nil {
		if err == todo.ErrTodoNotFound {
			handleError(w, err, http.StatusNotFound)
			return
		}
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

