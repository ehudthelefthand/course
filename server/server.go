package server

import (
	"github.com/ehudthelefthand/course/db"
	"github.com/ehudthelefthand/course/gorm"
	"github.com/ehudthelefthand/course/handler"
	"github.com/ehudthelefthand/course/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(db db.DB, gorm *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/courses", handler.ListCourses(db))
	r.GET("/courses/:id", handler.GetCourse(gorm))
	r.POST("/courses", handler.CreateCourse(gorm))
	r.POST("/classes", handler.CreateClasses(gorm))
	r.POST("/enrollments", middleware.RequireUser(gorm), handler.EnrollClass(gorm))
	r.POST("/register", handler.Register(gorm))
	r.POST("/login", handler.Login(gorm))

	return r
}
