package main

import (
	"darbelis.eu/stabas/di"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using default values")
	}

	// Set CheckAuthorization from environment
	checkAuthStr := strings.ToLower(os.Getenv("CHECK_AUTHORIZATION"))
	if checkAuthStr == "true" {
		di.AuthenticationManager.CheckAuthorization = true
		fmt.Println("Authorization checking is ENABLED")
	} else {
		di.AuthenticationManager.CheckAuthorization = false
		fmt.Println("Authorization checking is DISABLED")
	}

	// TODO check for certificates
	// build if missing one

	router := gin.Default()

	router.LoadHTMLFiles(
		"templates/index.html",
		"templates/tasks.html",
		"templates/participants.html",
		"templates/settings.html",
		"templates/login.html",
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

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.POST("/api/login", func(c *gin.Context) { di.AuthenticationControllerInstance.Login(c) })
	router.GET("/api/user", func(c *gin.Context) { di.AuthenticationControllerInstance.User(c) })

	//api.NewTaskController()
	router.GET("/api/tasks", func(c *gin.Context) { di.TaskControllerInstance.GetAllTasks(c) })
	router.GET("/api/groups", func(c *gin.Context) { di.TaskControllerInstance.GetTasksGroups(c) })
	router.GET("/api/tasks/:id", func(c *gin.Context) { di.TaskControllerInstance.GetTask(c) })
	router.POST("/api/tasks", func(c *gin.Context) { di.TaskControllerInstance.AddTask(c) })
	router.POST("/api/tasks/:id/change-status", func(c *gin.Context) { di.TaskControllerInstance.ChangeStatus(c) })
	router.PUT("/api/tasks/:id", func(c *gin.Context) { di.TaskControllerInstance.UpdateTask(c) })
	router.DELETE("/api/tasks/:id", func(c *gin.Context) { di.TaskControllerInstance.DeleteTask(c) })

	router.GET("/api/participants", func(c *gin.Context) { di.ParticipantsControllerInstance.GetParticipants(c) })
	router.GET("/api/participants/:id", func(c *gin.Context) { di.ParticipantsControllerInstance.GetParticipant(c) })
	router.POST("/api/participants/:id/regenerate-password", func(c *gin.Context) { di.ParticipantsControllerInstance.RegeneratePassword(c) })
	router.PUT("/api/participants/:id", func(c *gin.Context) { di.ParticipantsControllerInstance.UpdateParticipant(c) })
	router.POST("/api/participants", func(c *gin.Context) { di.ParticipantsControllerInstance.AddParticipant(c) })
	router.DELETE("/api/participants/:id", func(c *gin.Context) { di.ParticipantsControllerInstance.DeleteParticipant(c) })

	router.GET("/api/settings", func(c *gin.Context) { di.SettingsControllerInstance.GetSettings(c) })
	router.POST("/api/settings", func(c *gin.Context) { di.SettingsControllerInstance.UpdateSettings(c) })

	router.Static("/assets/js", "./assets/js")
	router.Static("/assets/css", "./assets/css")
	router.Static("/assets/img", "./assets/img")

	err = router.RunTLS(":8443", "./tls/server.crt", "./tls/server.key")

	fmt.Println(err)
}
