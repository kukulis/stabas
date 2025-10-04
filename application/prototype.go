package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.LoadHTMLFiles("templates/prototype.html")
	router.GET("/prototype", func(c *gin.Context) {
		c.HTML(http.StatusOK, "prototype.html", gin.H{})
	})

	err := router.Run(":8088")

	fmt.Println(err)
}
