package api

import (
	"darbelis.eu/stabas/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskController struct {
	// temporary
	// later DB will be used
	tasks []entities.Task
}

func (controller *TaskController) GetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, controller.tasks)
}

var TaskControllerInstance = TaskController{tasks: make([]entities.Task, 0)}
