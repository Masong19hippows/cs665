package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Serve the favicon (from ./static/favicon.ico)
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("static/favicon.ico")
	})

	// Load the HTML file from the static directory
	router.LoadHTMLFiles("static/index.html")

	// Serve the index page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Start the server
	router.Run(":8080")
}
