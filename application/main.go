package main

import (
	"darbelis.eu/stabas/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.LoadHTMLFiles(
		"templates/index.html",
		"templates/tasks.html",
		"templates/participants.html",
		"templates/settings.html",
	)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/tasks", func(c *gin.Context) {
		fmt.Println("Tasks reloaded")
		c.HTML(http.StatusOK, "tasks.html", gin.H{})
	})
	router.GET("/participants", func(c *gin.Context) {
		c.HTML(http.StatusOK, "participants.html", gin.H{})
	})

	router.GET("/api/tasks", func(c *gin.Context) { api.TaskControllerInstance.GetAllTasks(c) })
	router.PUT("/api/tasks", func(c *gin.Context) { api.TaskControllerInstance.AddTask(c) })

	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.Static("/assets/js", "./assets/js")

	err := router.Run(":8088")

	fmt.Println(err)
}
