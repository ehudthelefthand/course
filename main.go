package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ehudthelefthand/course/db"
	"github.com/ehudthelefthand/course/model"
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
	r.POST("/courses", createCourse(db))
	r.POST("/classes", createClasses(db))
	r.POST("/enrollments", enrollClass(db))

	r.Run(":8080")
}

func listCourses(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		courses, err := db.GetAllCourse()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
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

func createCourse(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(model.Course)
		if err := c.BindJSON(req); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": "can not parse course",
			})
			return
		}
		if err := db.CreateCourse(req); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.IndentedJSON(http.StatusOK, req)
	}
}

type ClassReq struct {
	CourseID  uint      `json:"course_id"`
	TrainerID uint      `json:"trainer_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Seats     int       `json:"seats"`
}

func createClasses(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(ClassReq)
		if err := c.BindJSON(req); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		course, err := db.GetCourse(req.CourseID)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": "coures not found",
			})
			return
		}
		class, err := course.CreateClass(req.Start, req.End)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		if err := class.SetSeats(req.Seats); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		class.Trainer.ID = req.TrainerID

		if err := db.SaveClass(class); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	}
}

type EnrollmentReq struct {
	StudentID uint `json:"student_id"`
	ClassID   uint `json:"class_id"`
}

func enrollClass(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(EnrollmentReq)
		if err := c.BindJSON(req); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		class, err := db.GetClass(req.ClassID)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		student, err := db.GetStudent(req.StudentID)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		if err := class.AddStudent(*student); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		if err := db.CreateClassStudent(student.ID, class.ID); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	}
}
