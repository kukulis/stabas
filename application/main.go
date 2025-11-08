package main

import (
	"darbelis.eu/stabas/api"
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
		"templates/test.html",
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
	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.html", gin.H{"menu": menuContentHtml})
	})
	router.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "settings.html", gin.H{"menu": menuContentHtml})
	})

	router.GET("/api/tasks", func(c *gin.Context) { api.TaskControllerInstance.GetAllTasks(c) })
	router.GET("/api/tasks/:id", func(c *gin.Context) { api.TaskControllerInstance.GetTask(c) })
	router.PUT("/api/tasks", func(c *gin.Context) { api.TaskControllerInstance.AddTask(c) })
	router.POST("/api/tasks/:id/change-status", func(c *gin.Context) { api.TaskControllerInstance.ChangeStatus(c) })
	router.POST("/api/tasks/:id", func(c *gin.Context) { api.TaskControllerInstance.UpdateTask(c) })
	router.DELETE("/api/tasks/:id", func(c *gin.Context) { api.TaskControllerInstance.DeleteTask(c) })

	router.GET("/api/participants", func(c *gin.Context) { api.ParticipantControllerInstance.GetParticipants(c) })
	router.POST("/api/participants/:id", func(c *gin.Context) { api.ParticipantControllerInstance.UpdateParticipant(c) })

	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.Static("/assets/js", "./assets/js")

	err := router.Run(":8088")

	fmt.Println(err)
}
