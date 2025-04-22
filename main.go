package main

import (
	"cs665/db"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server...")

	database, err := db.NewDatabase("./db/Project1.db")
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		panic(err)
	}
	fmt.Println("Database connected.")

	router := gin.Default()

	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.LoadHTMLFiles("static/index.html")

	router.GET("/", func(c *gin.Context) {
		fmt.Println("Serving index.html")
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api := router.Group("/api")
	{
		api.GET("/dishes", func(c *gin.Context) {
			fmt.Println("GET /api/dishes called")
			dishes, err := database.GetAllDishes()
			if err != nil {
				fmt.Println("Error fetching dishes:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fmt.Printf("Fetched %d dishes\n", len(dishes))
			c.JSON(http.StatusOK, dishes)
		})

		api.POST("/dishes", func(c *gin.Context) {
			fmt.Println("POST /api/dishes called")
			var d db.Dish
			if err := c.ShouldBindJSON(&d); err != nil {
				fmt.Println("Error binding JSON:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			fmt.Println("Adding dish:", d)
			if err := database.AddDish(d); err != nil {
				fmt.Println("Error adding dish:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "dish added"})
		})

		api.PUT("/dishes/:id", func(c *gin.Context) {
			fmt.Println("PUT /api/dishes/:id called")
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				fmt.Println("Invalid dish ID:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
				return
			}
			var d db.Dish
			if err := c.ShouldBindJSON(&d); err != nil {
				fmt.Println("Error binding JSON:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			d.ID = id
			fmt.Println("Updating dish:", d)
			if err := database.UpdateDish(d); err != nil {
				fmt.Println("Error updating dish:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "dish updated"})
		})

		api.DELETE("/dishes/:id", func(c *gin.Context) {
			fmt.Println("DELETE /api/dishes/:id called")
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				fmt.Println("Invalid dish ID:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
				return
			}
			fmt.Println("Deleting dish ID:", id)
			if err := database.DeleteDish(id); err != nil {
				fmt.Println("Error deleting dish:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "dish deleted"})
		})
	}

	err = router.Run(":8080")
	if err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
