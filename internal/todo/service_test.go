package todo

import "testing"

func TestCompleteTodo(t *testing.T) {
	service := NewService()

	// Create a todo
	todoItem, err := service.Create("Write unit tests")
	if err != nil {
		t.Fatalf("unexpected error creating todo: %v", err)
	}

	// Mark as completed
	completedTodo, err := service.Complete(todoItem.ID)
	if err != nil {
		t.Fatalf("unexpected error completing todo: %v", err)
	}

	if !completedTodo.Completed {
		t.Errorf("expected todo to be completed")
	}

	// Call Complete again (idempotency)
	completedAgain, err := service.Complete(todoItem.ID)
	if err != nil {
		t.Fatalf("unexpected error on second complete: %v", err)
	}

	if !completedAgain.Completed {
		t.Errorf("todo should remain completed")
	}
}
