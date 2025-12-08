package api

import (
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/my_tests"
	"testing"
)

// DeleteTask Tests - Only checks Sender (not Receivers)
func TestAuthorize_DeleteTask_AsSender(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender: 1, // HQ is sender
	}

	result := manager.Authorize(c, "HQ", "DeleteTask", task)

	if !result {
		t.Errorf("Expected Authorize to return true for DeleteTask when user is sender, got false")
	}
}

func TestAuthorize_DeleteTask_NotSender(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender: 2, // Someone else is sender
	}

	result := manager.Authorize(c, "HQ", "DeleteTask", task)

	if result {
		t.Errorf("Expected Authorize to return false for DeleteTask when user is not sender, got true")
	}
}

func TestAuthorize_DeleteTask_AsAdmin(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender: 2, // Someone else is sender
	}

	result := manager.Authorize(c, "admin", "DeleteTask", task)

	// Admin should always be able to delete
	if !result {
		t.Errorf("Expected Authorize to return true for DeleteTask when user is admin, got false")
	}
}

// ChangeStatus Tests - Checks Sender OR Receiver

func TestAuthorize_ChangeStatus_AsSender(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender:    1, // HQ is sender
		Receivers: []int{2},
	}

	result := manager.Authorize(c, "HQ", "ChangeStatus", task)

	if !result {
		t.Errorf("Expected Authorize to return true for ChangeStatus when user is sender, got false")
	}
}

func TestAuthorize_ChangeStatus_AsReceiver(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender:    2,        // Someone else is sender
		Receivers: []int{1}, // HQ is receiver
	}

	result := manager.Authorize(c, "HQ", "ChangeStatus", task)

	if !result {
		t.Errorf("Expected Authorize to return true for ChangeStatus when user is receiver, got false")
	}
}

func TestAuthorize_ChangeStatus_NotSenderOrReceiver(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender:    2,        // Someone else is sender
		Receivers: []int{3}, // Someone else is receiver
	}

	result := manager.Authorize(c, "HQ", "ChangeStatus", task)

	if result {
		t.Errorf("Expected Authorize to return false for ChangeStatus when user is neither sender nor receiver, got true")
	}
}

func TestAuthorize_ChangeStatus_AsAdmin(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender:    2,
		Receivers: []int{3},
	}

	result := manager.Authorize(c, "admin", "ChangeStatus", task)

	// Admin should always be able to change status
	if !result {
		t.Errorf("Expected Authorize to return true for ChangeStatus when user is admin, got false")
	}
}

func TestAuthorize_ChangeStatus_NoReceivers(t *testing.T) {
	manager := NewAuthenticationManager(my_tests.NewParticipantsRepository())
	c := &TestJSONResponder{}
	task := &entities.Task{
		Sender:    1, // HQ is sender
		Receivers: []int{},
	}

	result := manager.Authorize(c, "HQ", "ChangeStatus", task)

	// Should pass because HQ is sender (even without receivers)
	if !result {
		t.Errorf("Expected Authorize to return true for ChangeStatus when user is sender with no receivers, got false")
	}
}
