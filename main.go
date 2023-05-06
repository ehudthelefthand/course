package main

import (
	"log"

	"github.com/ehudthelefthand/course/db"
	"github.com/ehudthelefthand/course/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/courses", handler.ListCourses(db))
	r.GET("/courses/:id", handler.GetCourse(db))
	r.POST("/courses", handler.CreateCourse(db))
	r.POST("/classes", handler.CreateClasses(db))
	r.POST("/enrollments", handler.EnrollClass(db))

	r.Run(":8080")
}

func Error(c *gin.Context, status int, err error) {
	log.Println(err)
	c.JSON(status, gin.H{
		"message": err.Error(),
	})
}
