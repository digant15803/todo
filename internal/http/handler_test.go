package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	todo "todo/internal/todo-api"
)

func setupTestHandler() (*http.ServeMux, *Handler) {
    service := todo.NewService()
    handler := NewHandler(service)

    mux := http.NewServeMux()
    handler.RegisterRoutes(mux)

    return mux, handler
}

// Test the CreateTodo endpoint
func TestCompleteTodoEndpoint(t *testing.T) {
	// Arrange
	mux, _ := setupTestHandler()

	// Create a todo first
	createReqBody := []byte(`{"title":"Test HTTP completion"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(createReqBody))
	createReq.Header.Set("Content-Type", "application/json")

	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createRec.Code)
	}

	var created todo.Todo
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}

	// Act: complete the todo
	completeReq := httptest.NewRequest(
		http.MethodPost,
		"/todos/"+created.ID+"/complete",
		nil,
	)
	completeRec := httptest.NewRecorder()
	mux.ServeHTTP(completeRec, completeReq)

	// Assert
	if completeRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d %s", completeRec.Code, "/todos/"+created.ID+"/complete")
	}

	var completed todo.Todo
	if err := json.NewDecoder(completeRec.Body).Decode(&completed); err != nil {
		t.Fatalf("failed to decode complete response: %v", err)
	}

	if !completed.Completed {
		t.Errorf("expected todo to be completed")
	}
}

// Test for a non-existent ID
func TestDeleteTodo_NotFound(t *testing.T) {
	// Arrange
	mux, _ := setupTestHandler()

	req := httptest.NewRequest(http.MethodDelete, "/todos/nonexistent-id", nil)
	rec := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

// Test for completing a non-existent todo
func TestCompleteTodo_NotFound(t *testing.T) {
	// Arrange
	mux, _ := setupTestHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/todos/nonexistent-id/complete",
		nil,
	)
	rec := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

// Test for invalid HTTP method
func TestCompleteTodo_InvalidMethod(t *testing.T) {
	// Arrange
	mux, _ := setupTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/todos/123/complete",
		nil,
	)
	rec := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

// Test for invalid request body
func TestCreateTodo_InvalidBody(t *testing.T) {
	// Arrange
	mux, _ := setupTestHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/todos",
		bytes.NewBuffer([]byte(`{invalid-json}`)),
	)
	rec := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

// Test for invalid JSON structure
func TestCreateTodo_InvalidJSON(t *testing.T) {
    // Arrange
    mux, _ := setupTestHandler()

    req := httptest.NewRequest(
        http.MethodPost,
        "/todos",
        bytes.NewBuffer([]byte(`{ "title": "Test Todo", }`)), // Invalid JSON with trailing comma
    )
    rec := httptest.NewRecorder()

    // Act
    mux.ServeHTTP(rec, req)

    // Assert
    if rec.Code != http.StatusBadRequest {
        t.Fatalf("expected status 400, got %d", rec.Code)
    }

    expectedError := "invalid JSON body"
    if !strings.Contains(rec.Body.String(), expectedError) {
        t.Errorf("expected error message to contain %q, got %q", expectedError, rec.Body.String())
    }
}