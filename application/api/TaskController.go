package api

import (
	"darbelis.eu/stabas/entities"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
	"net/http"
)

type TaskController struct {
	// temporary
	// later DB will be used
	tasks []*entities.Task
	maxId int
}

func (controller *TaskController) GetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, controller.tasks)
}

func (controller *TaskController) AddTask(c *gin.Context) {
	buf, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, "error reading buffer "+err.Error())
		return
	}

	task := entities.NewTask()

	err = json.API.Unmarshal(buf, &task)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error parsing json"+err.Error())
		return
	}

	controller.maxId++
	task.Id = controller.maxId

	controller.tasks = append(controller.tasks, task)

	c.JSON(http.StatusOK, task.Id)
}

func (controller *TaskController) ChangeStatus(c *gin.Context) {
	id := c.Param("id")
	status := c.Query("status")

	// TODO

	c.JSON(http.StatusOK, "TODO change status of task "+id+" to "+status)
}

var TaskControllerInstance = TaskController{tasks: make([]*entities.Task, 0), maxId: 0}
