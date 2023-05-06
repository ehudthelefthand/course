package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Course struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"desc"`
}

var courses = []Course{
	{ ID: "1", Name: "TDD", Description: "Let's make a bug free software!" },
	{ ID: "2", Name: "CI&CD", Description: "Continue deliver a customer value" },
}

func main() {
	r := gin.Default()

	r.GET("/courses", listCourses)
	r.GET("/courses/:id", getCourse)

	r.Run(":8080")
}

func listCourses(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, courses)
}

func getCourse(c *gin.Context) {
	id := c.Param("id")
	for _, course := range courses {
		if course.ID == id {
			c.IndentedJSON(http.StatusOK, course)
			return
		}	
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"message": "course not found",
	})
}