package main

import (
	"darbelis.eu/stabas/di"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()

	router.LoadHTMLFiles(
		"templates/index.html",
		"templates/tasks.html",
		"templates/participants.html",
		"templates/settings.html",
	)

	templatesDir := os.DirFS("templates")
	menuContent, ferr := fs.ReadFile(templatesDir, "menu.html")
	if ferr != nil {
		fmt.Println("Error reading menu.html " + ferr.Error())
	}

	menuContentHtml := template.HTML(menuContent)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"menu": menuContentHtml})
	})
	router.GET("/tasks", func(c *gin.Context) {
		fmt.Println("Tasks reloaded")
		c.HTML(http.StatusOK, "tasks.html", gin.H{"menu": menuContentHtml})
	})
	router.GET("/participants", func(c *gin.Context) {
		c.HTML(http.StatusOK, "participants.html", gin.H{"menu": menuContentHtml})
	})

	router.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "settings.html", gin.H{"menu": menuContentHtml})
	})

	//api.NewTaskController()
	router.GET("/api/tasks", func(c *gin.Context) { di.TaskControllerInstance.GetAllTasks(c) })
	//router.GET("/api/tasks2", func(c *gin.Context) { testController.GetAllTasks(c) })
	router.GET("/api/tasks/:id", func(c *gin.Context) { di.TaskControllerInstance.GetTask(c) })
	router.PUT("/api/tasks", func(c *gin.Context) { di.TaskControllerInstance.AddTask(c) })
	router.POST("/api/tasks/:id/change-status", func(c *gin.Context) { di.TaskControllerInstance.ChangeStatus(c) })
	router.POST("/api/tasks/:id", func(c *gin.Context) { di.TaskControllerInstance.UpdateTask(c) })
	router.DELETE("/api/tasks/:id", func(c *gin.Context) { di.TaskControllerInstance.DeleteTask(c) })

	router.GET("/api/participants", func(c *gin.Context) { di.ParticipantsControllerInstance.GetParticipants(c) })
	router.POST("/api/participants/:id", func(c *gin.Context) { di.ParticipantsControllerInstance.UpdateParticipant(c) })
	router.PUT("/api/participants", func(c *gin.Context) { di.ParticipantsControllerInstance.AddParticipant(c) })
	router.DELETE("/api/participants/:id", func(c *gin.Context) { di.ParticipantsControllerInstance.DeleteParticipant(c) })

	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.Static("/assets/js", "./assets/js")
	router.Static("/assets/css", "./assets/css")

	err := router.Run(":8088")

	fmt.Println(err)
}
