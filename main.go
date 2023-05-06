package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ehudthelefthand/course/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/courses", listCourses(db))
	r.GET("/courses/:id", getCourse(db))

	r.Run(":8080")
}

func listCourses(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		courses, err := db.GetAllCourse()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "server error!",
			})
			return
		}
		c.IndentedJSON(http.StatusOK, courses)
	}
}

func getCourse(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": "invalid id",
			})
			return
		}
		course, err := db.GetCourse(uint(id))
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": "course not found",
			})
			return
		}
		c.IndentedJSON(http.StatusOK, course)
	}
}
