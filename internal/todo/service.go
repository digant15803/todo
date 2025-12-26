package todo

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
	ErrInvalidTitle = errors.New("title cannot be empty")
)

type Service struct {
	mu    sync.Mutex
	todos map[string]Todo
}

func NewService() *Service {
	return &Service{
		todos: make(map[string]Todo),
	}
}

// Create adds a new todo item with the given title.
// Generates a unique ID for the todo.
func (s *Service) Create(title string) (Todo, error) {
	if title == "" {
		return Todo{}, ErrInvalidTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	todo := Todo{
		ID:        uuid.NewString(),
		Title:     title,
		Completed: false,
	}

	s.todos[todo.ID] = todo
	return todo, nil
}

// List returns all todo items.
func (s *Service) List() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		result = append(result, todo)
	}
	return result
}

// Delete removes the todo item with the given ID.
func (s *Service) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[id]; !exists {
		return ErrTodoNotFound
	}

	delete(s.todos, id)
	return nil
}

// Complete marks the todo item with the given ID as completed.
func (s *Service) Complete(id string) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, exists := s.todos[id]
	if !exists {
		return Todo{}, ErrTodoNotFound
	}

	if !todo.Completed {
		todo.Completed = true
		s.todos[id] = todo
	}

	return todo, nil
}

