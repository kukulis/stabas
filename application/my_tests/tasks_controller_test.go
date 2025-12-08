package my_tests

import (
	"bytes"
	"darbelis.eu/stabas/api"
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var testableTime, _ = time.Parse(time.RFC3339, "2025-12-05T15:04:05Z")

func setupTestController() (*api.TaskController, *dao.TasksRepository) {
	tasksRepository := NewTasksRepository()
	participantsRepository := NewParticipantsRepository()

	timeProvider := util.FixedTimeProvider{Time: testableTime}
	authManager := api.NewAuthenticationManager(participantsRepository)
	authManager.CheckAuthorization = false
	controller := api.NewTaskController(tasksRepository, participantsRepository, timeProvider, authManager)

	return controller, tasksRepository
}

func TestGetAllTasks(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.GetAllTasks(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var tasks []*entities.Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(tasks) == 0 {
		t.Error("Expected non-empty tasks list")
	}

	// TODO assert that received tasks list is equal to the exact list
}

func TestGetTasksGroups(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.GetTasksGroups(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var groupedTasks []*entities.Task
	err := json.Unmarshal(w.Body.Bytes(), &groupedTasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(groupedTasks) == 0 {
		t.Error("Expected non-empty grouped tasks list")
	}
	// TODO assert more exact tasks groups ( regarding data in the repository_factory.go )
}

func TestGetTask_Success(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	controller.GetTask(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var task entities.Task
	err := json.Unmarshal(w.Body.Bytes(), &task)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if task.Id != 1 {
		t.Errorf("Expected task ID 1, got %d", task.Id)
	}
}

func TestGetTask_InvalidId(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	controller.GetTask(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "9999"}}

	controller.GetTask(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAddTask(t *testing.T) {
	controller, repo := setupTestController()

	newTask := map[string]interface{}{
		"message":   "New test task",
		"result":    "",
		"sender":    1,
		"receivers": []int{2},
		"status":    entities.STATUS_NEW,
	}

	body, _ := json.Marshal(newTask)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("", "/", bytes.NewBuffer(body))

	initialCount := len(repo.FindAll())

	controller.AddTask(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var createdTask entities.Task
	err := json.Unmarshal(w.Body.Bytes(), &createdTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if createdTask.Id == 0 {
		t.Error("Expected task to have an ID assigned")
	}

	if createdTask.Message != "New test task" {
		t.Errorf("Expected message 'New test task', got '%s'", createdTask.Message)
	}

	if *createdTask.CreatedAt != testableTime {
		t.Errorf("Expected CreatedAt to be %v, got %v", testableTime, *createdTask.CreatedAt)
	}

	finalCount := len(repo.FindAll())
	if finalCount != initialCount+1 {
		t.Errorf("Expected %d tasks, got %d", initialCount+1, finalCount)
	}
}

func TestUpdateTask_Simple(t *testing.T) {
	controller, repo := setupTestController()

	existingTask, _ := repo.FindById(1)
	updatedTask := map[string]interface{}{
		"message":   "Updated message",
		"result":    "Updated result",
		"sender":    existingTask.Sender,
		"receivers": existingTask.Receivers,
		"status":    existingTask.Status,
		"version":   existingTask.Version + 1,
	}

	body, _ := json.Marshal(updatedTask)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))

	controller.UpdateTask(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result entities.Task
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.Message != "Updated message" {
		t.Errorf("Expected message 'Updated message', got '%s'", result.Message)
	}
}

func TestUpdateTask_InvalidId(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}
	c.Request = httptest.NewRequest("PUT", "/tasks/invalid", bytes.NewBuffer([]byte("{}")))

	controller.UpdateTask(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTask_SplitWithMultipleReceivers(t *testing.T) {
	controller, repo := setupTestController()

	initialTasksCount := len(repo.FindAll())

	updatedTask := map[string]interface{}{
		"message":   "Task for multiple receivers",
		"result":    "",
		"sender":    1,
		"receivers": []int{2, 3, 4},
		"status":    entities.STATUS_NEW,
		"version":   2,
	}

	body, _ := json.Marshal(updatedTask)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "8"}}
	c.Request = httptest.NewRequest("", "/", bytes.NewBuffer(body))

	controller.UpdateTask(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)

		errorResult := map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &errorResult)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		t.Errorf("Received error is [%s]", errorResult["error"])
		return
	}

	var result entities.Task
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(result.Receivers) != 1 {
		t.Errorf("Expected 1 receiver on parent task, got %d", len(result.Receivers))
	}

	if result.TaskGroup != 8 {
		t.Errorf("Expected task group to be 1, got %d", result.TaskGroup)
	}

	if len(result.Children) != 2 {
		t.Errorf("Expected 2 child tasks, got %d", len(result.Children))
	}

	finalTasksCount := len(repo.FindAll())
	expectedCount := initialTasksCount + 2
	if finalTasksCount != expectedCount {
		t.Errorf("Expected %d tasks total, got %d", expectedCount, finalTasksCount)
	}

	for i, child := range result.Children {
		if child.TaskGroup != result.TaskGroup {
			t.Errorf("Child %d has wrong task group: expected %d, got %d", i, result.TaskGroup, child.TaskGroup)
		}
		if len(child.Receivers) != 1 {
			t.Errorf("Child %d should have 1 receiver, got %d", i, len(child.Receivers))
		}
		if child.Message != result.Message {
			t.Errorf("Child %d has different message than parent", i)
		}
	}
}

func TestUpdateTask_VersionConflict(t *testing.T) {
	controller, repo := setupTestController()

	existingTask, _ := repo.FindById(1)
	updatedTask := map[string]interface{}{
		"message":   "Updated message",
		"result":    "",
		"sender":    existingTask.Sender,
		"receivers": existingTask.Receivers,
		"status":    existingTask.Status,
		"version":   existingTask.Version,
	}

	body, _ := json.Marshal(updatedTask)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	// it seems method parameter is not important
	// also url is not important
	c.Request = httptest.NewRequest("", "/", bytes.NewBuffer(body))

	controller.UpdateTask(c)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d (conflict), got %d", http.StatusConflict, w.Code)
	}
}

func TestUpdateTask_ModifyReceiverNonNewStatus(t *testing.T) {
	controller, repo := setupTestController()

	existingTask, _ := repo.FindById(7)

	updatedTask := map[string]interface{}{
		"message":   "Task with multiple receivers",
		"result":    "",
		"sender":    existingTask.Sender,
		"receivers": []int{1},
		"status":    entities.STATUS_SENT,
		"version":   existingTask.Version + 1,
	}

	body, _ := json.Marshal(updatedTask)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "7"}}
	c.Request = httptest.NewRequest("", "/", bytes.NewBuffer(body))

	controller.UpdateTask(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d (bad request) when updating non-NEW task receivers, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTask_MultipleReceiversWithParent(t *testing.T) {
	controller, repo := setupTestController()

	existingTask, _ := repo.FindById(2)

	updatedTask := map[string]interface{}{
		"message":    "Task with multiple receivers and parent",
		"result":     "",
		"sender":     existingTask.Sender,
		"receivers":  []int{2, 3, 4},
		"status":     entities.STATUS_NEW,
		"version":    existingTask.Version + 1,
		"task_group": 1,
	}

	body, _ := json.Marshal(updatedTask)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "2"}}
	c.Request = httptest.NewRequest("PUT", "/tasks/2", bytes.NewBuffer(body))

	controller.UpdateTask(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d (bad request) when updating task with parent and multiple receivers, got %d", http.StatusBadRequest, w.Code)
	}
}

// TODO test case with a non existant sender or receiver

func TestUpdateTask_MultipleReceiversWithChildren(t *testing.T) {
	controller, repo := setupTestController()

	existingTask, _ := repo.FindById(1)

	child := entities.NewTask()
	child.Message = "Child task"
	child.Sender = 1
	child.Receivers = []int{5}
	child.TaskGroup = 1
	repo.AddTask(child)

	existingTask.Children = []*entities.Task{child}
	_, err := repo.UpdateTask(existingTask)
	if err != nil {
		t.Fatalf("Expected no error when updating task with children, got: %v", err)
	}

	updatedTask := map[string]interface{}{
		"message":   "Task with multiple receivers and children",
		"result":    "",
		"sender":    existingTask.Sender,
		"receivers": []int{2, 3, 4},
		"status":    entities.STATUS_NEW,
		"version":   existingTask.Version + 1,
	}

	body, _ := json.Marshal(updatedTask)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))

	controller.UpdateTask(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d (bad request) when updating task with children and multiple receivers, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	controller, repo := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	initialCount := len(repo.FindAll())
	controller.DeleteTask(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	finalCount := len(repo.FindAll())
	if finalCount != initialCount-1 {
		t.Errorf("Expected %d tasks after deletion, got %d", initialCount-1, finalCount)
	}

	task, err := repo.FindById(1)
	if err != nil {
		t.Fatalf("Task should still exist in repo: %v", err)
	}
	if !task.Deleted {
		t.Error("Task should be marked as deleted")
	}
}

func TestDeleteTask_InvalidId(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	controller.DeleteTask(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChangeStatus_Success(t *testing.T) {
	controller, repo := setupTestController()

	task, _ := repo.FindById(3)
	currentStatus := task.Status
	newStatus := currentStatus + 1

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "3"}}
	c.Request = httptest.NewRequest("POST", "/tasks/3/status?status="+string(rune(newStatus+'0')), nil)

	controller.ChangeStatus(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var updatedTask entities.Task
	err := json.Unmarshal(w.Body.Bytes(), &updatedTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedTask.Status != newStatus {
		t.Errorf("Expected status %d, got %d", newStatus, updatedTask.Status)
	}

	statusTime := updatedTask.GetStatusTime()
	if statusTime == nil {
		t.Error("Expected status time to be set")
	} else {
		if time.Since(*statusTime) > time.Second {
			t.Error("Status time should be recent")
		}
	}
}

func TestChangeStatus_InvalidId(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}
	c.Request = httptest.NewRequest("POST", "/tasks/invalid/status?status=2", nil)

	controller.ChangeStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChangeStatus_InvalidStatus(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("POST", "/tasks/1/status?status=invalid", nil)

	controller.ChangeStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChangeStatus_InvalidTransition(t *testing.T) {
	controller, repo := setupTestController()

	task, _ := repo.FindById(1)
	currentStatus := task.Status
	invalidNewStatus := currentStatus + 2

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("POST", "/tasks/1/status?status="+string(rune(invalidNewStatus+'0')), nil)

	controller.ChangeStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChangeStatus_TaskNotFound(t *testing.T) {
	controller, _ := setupTestController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "9999"}}
	c.Request = httptest.NewRequest("POST", "/tasks/9999/status?status=2", nil)

	controller.ChangeStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChangeStatus_NewTaskWithNoReceivers(t *testing.T) {
	controller, repo := setupTestController()

	// Create a NEW task with no receivers
	newTask := entities.NewTask()
	newTask.Message = "Task without receivers"
	newTask.Sender = 1
	newTask.Receivers = []int{}
	newTask.Status = entities.STATUS_NEW
	repo.AddTask(newTask)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(newTask.Id)}}
	c.Request = httptest.NewRequest("POST", "/tasks/"+strconv.Itoa(newTask.Id)+"/status?status="+strconv.Itoa(entities.STATUS_SENT), nil)

	controller.ChangeStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d when changing NEW task with no receivers to SENT, got %d", http.StatusBadRequest, w.Code)
		return
	}

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse["error"] == "" {
		t.Error("Expected error message in response")
	}
}
