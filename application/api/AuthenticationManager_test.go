package api

import (
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/my_tests"
	"testing"
)

// Tests with task.Sender = 1
func TestAuthorize_AddTask_WithSender1(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender: 1,
	}

	result := manager.Authorize(c, "HQ", "AddTask", task)

	if !result {
		t.Errorf("Expected Authorize to return true for UpdateTask with Sender=1, got false")
	}
}
func TestAuthorize_AddTask_WithSender2(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender: 2,
	}

	result := manager.Authorize(c, "HQ", "AddTask", task)

	if result {
		t.Errorf("Expected Authorize to return false for UpdateTask with Sender=1, got true")
	}
}

func TestAuthorize_UpdateTask_WithSender1(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender: 1,
	}

	result := manager.Authorize(c, "HQ", "UpdateTask", task)

	if !result {
		t.Errorf("Expected Authorize to return true for UpdateTask with Sender=1, got false")
	}
}

// Tests with task.Sender = 2
func TestAuthorize_UpdateTask_WithSender2(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender: 2,
	}

	result := manager.Authorize(c, "HQ", "UpdateTask", task)

	if result {
		t.Errorf("Expected Authorize to return false for UpdateTask with Sender=2, got true")
	}
}

// Tests with task.Receivers = [1]
func TestAuthorize_UpdateTask_WithReceivers1(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Receivers: []int{1},
	}

	result := manager.Authorize(c, "HQ", "UpdateTask", task)

	if !result {
		t.Errorf("Expected Authorize to return true for UpdateTask with Receivers=[1], got false")
	}
}

// Tests with task.Receivers = [2]
func TestAuthorize_UpdateTask_WithReceivers2(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Receivers: []int{2},
	}

	result := manager.Authorize(c, "HQ", "UpdateTask", task)

	if result {
		t.Errorf("Expected Authorize to return false for UpdateTask with Receivers=[2], got true")
	}
}
