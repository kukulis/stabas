package my_tests

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/db"
	"darbelis.eu/stabas/entities"
	"testing"
	"time"
)

func setupLiteTaskTestDB(t *testing.T) *dao.LiteTaskRepository {
	t.Helper()
	database := db.NewDatabase(":memory:")
	repo, err := dao.NewLiteTaskRepository(database)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	return repo
}

func TestNewLiteTaskRepository(t *testing.T) {
	database := db.NewDatabase(":memory:")
	repo, err := dao.NewLiteTaskRepository(database)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}

	defer func() { _ = repo.Close() }()
}

func TestLiteAddTask(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Test task",
		Result:    "",
		Sender:    1,
		Receivers: []int{2, 3},
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	id := repo.AddTask(task)

	if id == 0 {
		t.Error("Expected task ID to be set, got 0")
	}

	if task.Id != id {
		t.Errorf("Expected task.Id to be %d, got %d", id, task.Id)
	}

	if task.TaskGroup != id {
		t.Errorf("Expected TaskGroup to be set to task ID %d, got %d", id, task.TaskGroup)
	}
}

func TestLiteAddTaskWithTaskGroup(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Test task",
		Sender:    1,
		Receivers: []int{2},
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
		TaskGroup: 100,
	}

	_ = repo.AddTask(task)

	if task.TaskGroup != 100 {
		t.Errorf("Expected TaskGroup to remain 100, got %d", task.TaskGroup)
	}
}

func TestLiteFindTaskById(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Test task",
		Result:    "Test result",
		Sender:    1,
		Receivers: []int{2, 3},
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	id := repo.AddTask(task)
	found, err := repo.FindById(id)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if found.Id != id {
		t.Errorf("Expected ID %d, got %d", id, found.Id)
	}

	if found.Message != "Test task" {
		t.Errorf("Expected message 'Test task', got '%s'", found.Message)
	}

	if found.Result != "Test result" {
		t.Errorf("Expected result 'Test result', got '%s'", found.Result)
	}

	if found.Sender != 1 {
		t.Errorf("Expected sender 1, got %d", found.Sender)
	}

	if len(found.Receivers) != 2 || found.Receivers[0] != 2 || found.Receivers[1] != 3 {
		t.Errorf("Expected receivers [2, 3], got %v", found.Receivers)
	}

	if found.Status != entities.STATUS_NEW {
		t.Errorf("Expected status %d, got %d", entities.STATUS_NEW, found.Status)
	}

	if found.Version != 1 {
		t.Errorf("Expected version 1, got %d", found.Version)
	}
}

func TestLiteFindTaskByIdNotFound(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	_, err := repo.FindById(999)

	if err == nil {
		t.Error("Expected error for non-existent task, got nil")
	}
}

func TestLiteFindAllTasks(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()

	task1 := &entities.Task{Message: "Task 1", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1}
	task2 := &entities.Task{Message: "Task 2", Status: entities.STATUS_SENT, CreatedAt: &now, Version: 1}
	task3 := &entities.Task{Message: "Task 3", Status: entities.STATUS_RECEIVED, CreatedAt: &now, Version: 1}

	repo.AddTask(task1)
	repo.AddTask(task2)
	repo.AddTask(task3)

	tasks := repo.FindAll()

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

func TestLiteFindAllExcludesDeleted(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()

	task1 := &entities.Task{Message: "Task 1", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1}
	task2 := &entities.Task{Message: "Task 2", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1}
	task3 := &entities.Task{Message: "Task 3", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1}

	id1 := repo.AddTask(task1)
	repo.AddTask(task2)
	repo.AddTask(task3)

	_ = repo.DeleteTask(id1)

	tasks := repo.FindAll()

	if len(tasks) != 2 {
		t.Errorf("Expected 2 active tasks, got %d", len(tasks))
	}

	for _, task := range tasks {
		if task.Message == "Task 1" {
			t.Error("Deleted task 'Task 1' should not be in results")
		}
	}
}

func TestLiteUpdateTask(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Original message",
		Sender:    1,
		Receivers: []int{2},
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	id := repo.AddTask(task)

	task.Message = "Updated message"
	task.Result = "Updated result"
	task.Status = entities.STATUS_SENT
	task.Version = 2
	task.Receivers = []int{3, 4}

	updated, err := repo.UpdateTask(task)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.Message != "Updated message" {
		t.Errorf("Expected message 'Updated message', got '%s'", updated.Message)
	}

	if updated.Result != "Updated result" {
		t.Errorf("Expected result 'Updated result', got '%s'", updated.Result)
	}

	if updated.Status != entities.STATUS_SENT {
		t.Errorf("Expected status %d, got %d", entities.STATUS_SENT, updated.Status)
	}

	if len(updated.Receivers) != 2 || updated.Receivers[0] != 3 || updated.Receivers[1] != 4 {
		t.Errorf("Expected receivers [3, 4], got %v", updated.Receivers)
	}

	found, _ := repo.FindById(id)
	if found.Message != "Updated message" {
		t.Error("Task was not persisted correctly")
	}
}

func TestLiteUpdateTaskNotFound(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Id:        999,
		Message:   "Ghost task",
		CreatedAt: &now,
	}

	_, err := repo.UpdateTask(task)

	if err == nil {
		t.Error("Expected error when updating non-existent task, got nil")
	}
}

func TestLiteDeleteTask(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Task to delete",
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	id := repo.AddTask(task)

	err := repo.DeleteTask(id)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tasks := repo.FindAll()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 active tasks, got %d", len(tasks))
	}

	found, err := repo.FindById(id)
	if err != nil {
		t.Fatal("Expected to find deleted task by ID")
	}

	if !found.Deleted {
		t.Error("Expected task to be marked as deleted")
	}
}

func TestLiteDeleteTaskNotFound(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	err := repo.DeleteTask(999)

	if err == nil {
		t.Error("Expected error when deleting non-existent task, got nil")
	}
}

func TestLiteUpdateTaskWithValidation(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Original",
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	id := repo.AddTask(task)
	task.Id = id
	task.Message = "Updated"
	task.Version = 2

	updated, err := repo.UpdateTaskWithValidation(task)

	if err != nil {
		t.Fatalf("Expected no error with correct version, got %v", err)
	}

	if updated.Message != "Updated" {
		t.Errorf("Expected message 'Updated', got '%s'", updated.Message)
	}
}

func TestLiteUpdateTaskWithValidationWrongVersion(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Original",
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	id := repo.AddTask(task)
	task.Id = id
	task.Message = "Updated"
	task.Version = 5

	_, err := repo.UpdateTaskWithValidation(task)

	if err == nil {
		t.Error("Expected error for wrong version, got nil")
		return
	}

	if err.Error() != "Wrong task version" {
		t.Errorf("Expected 'Wrong task version' error, got '%s'", err.Error())
	}
}

func TestLiteGetCountWithSameGroup(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()

	task1 := &entities.Task{Message: "Task 1", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1, TaskGroup: 100}
	task2 := &entities.Task{Message: "Task 2", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1, TaskGroup: 100}
	task3 := &entities.Task{Message: "Task 3", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1, TaskGroup: 200}

	repo.AddTask(task1)
	repo.AddTask(task2)
	repo.AddTask(task3)

	count := repo.GetCountWithSameGroup(100)

	if count != 2 {
		t.Errorf("Expected 2 tasks in group 100, got %d", count)
	}

	count = repo.GetCountWithSameGroup(200)

	if count != 1 {
		t.Errorf("Expected 1 task in group 200, got %d", count)
	}
}

func TestLiteGetCountWithSameGroupExcludesDeleted(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()

	task1 := &entities.Task{Message: "Task 1", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1, TaskGroup: 100}
	task2 := &entities.Task{Message: "Task 2", Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1, TaskGroup: 100}

	id1 := repo.AddTask(task1)
	repo.AddTask(task2)

	_ = repo.DeleteTask(id1)

	count := repo.GetCountWithSameGroup(100)

	if count != 1 {
		t.Errorf("Expected 1 active task in group 100 (excluding deleted), got %d", count)
	}
}

func TestLiteTaskWithEmptyReceivers(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Task with no receivers",
		Sender:    1,
		Receivers: []int{},
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	id := repo.AddTask(task)
	found, err := repo.FindById(id)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if found.Receivers == nil {
		t.Error("Expected empty array, got nil")
	}

	if len(found.Receivers) != 0 {
		t.Errorf("Expected 0 receivers, got %d", len(found.Receivers))
	}
}

func TestLiteTaskWithNilReceivers(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Task with nil receivers",
		Sender:    1,
		Receivers: nil,
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	_ = repo.AddTask(task)
	found, err := repo.FindById(task.Id)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(found.Receivers) != 0 {
		t.Errorf("Expected 0 receivers, got %d", len(found.Receivers))
	}
}

func TestLiteTaskTimeFields(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	createdAt := time.Now().Add(-5 * time.Hour)
	sentAt := time.Now().Add(-4 * time.Hour)
	receivedAt := time.Now().Add(-3 * time.Hour)
	executingAt := time.Now().Add(-2 * time.Hour)
	finishedAt := time.Now().Add(-1 * time.Hour)

	task := &entities.Task{
		Message:     "Task with all time fields",
		Status:      entities.STATUS_FINISHED,
		CreatedAt:   &createdAt,
		SentAt:      &sentAt,
		ReceivedAt:  &receivedAt,
		ExecutingAt: &executingAt,
		FinishedAt:  &finishedAt,
		Version:     1,
	}

	id := repo.AddTask(task)
	found, err := repo.FindById(id)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if found.CreatedAt == nil || !found.CreatedAt.Equal(createdAt) {
		t.Error("CreatedAt not preserved correctly")
	}

	if found.SentAt == nil || !found.SentAt.Equal(sentAt) {
		t.Error("SentAt not preserved correctly")
	}

	if found.ReceivedAt == nil || !found.ReceivedAt.Equal(receivedAt) {
		t.Error("ReceivedAt not preserved correctly")
	}

	if found.ExecutingAt == nil || !found.ExecutingAt.Equal(executingAt) {
		t.Error("ExecutingAt not preserved correctly")
	}

	if found.FinishedAt == nil || !found.FinishedAt.Equal(finishedAt) {
		t.Error("FinishedAt not preserved correctly")
	}

	if found.ClosedAt != nil {
		t.Error("ClosedAt should be nil")
	}
}

func TestLiteTaskWithNilTimeFields(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:     "Task with minimal time fields",
		Status:      entities.STATUS_NEW,
		CreatedAt:   &now,
		SentAt:      nil,
		ReceivedAt:  nil,
		ExecutingAt: nil,
		FinishedAt:  nil,
		ClosedAt:    nil,
		Version:     1,
	}

	id := repo.AddTask(task)
	found, err := repo.FindById(id)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if found.CreatedAt == nil {
		t.Error("CreatedAt should not be nil")
	}

	if found.SentAt != nil {
		t.Error("SentAt should be nil")
	}

	if found.ReceivedAt != nil {
		t.Error("ReceivedAt should be nil")
	}

	if found.ExecutingAt != nil {
		t.Error("ExecutingAt should be nil")
	}

	if found.FinishedAt != nil {
		t.Error("FinishedAt should be nil")
	}

	if found.ClosedAt != nil {
		t.Error("ClosedAt should be nil")
	}
}

func TestLiteUpdateTaskPreservesTaskGroup(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()
	task := &entities.Task{
		Message:   "Original",
		Status:    entities.STATUS_NEW,
		CreatedAt: &now,
		Version:   1,
	}

	_ = repo.AddTask(task)

	originalGroup := task.TaskGroup

	task.Message = "Updated"
	task.TaskGroup = 0

	updated, err := repo.UpdateTask(task)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.TaskGroup != originalGroup {
		t.Errorf("Expected TaskGroup to be preserved as %d, got %d", originalGroup, updated.TaskGroup)
	}
}

func TestLiteMultipleTaskOperations(t *testing.T) {
	repo := setupLiteTaskTestDB(t)
	defer func() { _ = repo.Close() }()

	now := time.Now()

	task1 := &entities.Task{Message: "Task 1", Sender: 1, Receivers: []int{2}, Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1, TaskGroup: 100}
	task2 := &entities.Task{Message: "Task 2", Sender: 2, Receivers: []int{3}, Status: entities.STATUS_SENT, CreatedAt: &now, Version: 1, TaskGroup: 100}
	task3 := &entities.Task{Message: "Task 3", Sender: 1, Receivers: []int{3}, Status: entities.STATUS_NEW, CreatedAt: &now, Version: 1, TaskGroup: 200}

	id1 := repo.AddTask(task1)
	_ = repo.AddTask(task2)
	id3 := repo.AddTask(task3)

	task1.Message = "Task 1 Updated"
	task1.Version = 2
	_, _ = repo.UpdateTask(task1)

	_ = repo.DeleteTask(id3)

	all := repo.FindAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 active tasks, got %d", len(all))
	}

	found1, _ := repo.FindById(id1)
	if found1.Message != "Task 1 Updated" {
		t.Error("Task 1 update failed")
	}

	count := repo.GetCountWithSameGroup(100)
	if count != 2 {
		t.Errorf("Expected 2 tasks in group 100, got %d", count)
	}

	count = repo.GetCountWithSameGroup(200)
	if count != 0 {
		t.Errorf("Expected 0 active tasks in group 200 (deleted), got %d", count)
	}

	deleted, _ := repo.FindById(id3)
	if !deleted.Deleted {
		t.Error("Task 3 should be marked as deleted")
	}
}
