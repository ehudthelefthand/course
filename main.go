package main

import (
	"errors"
	"log"

	"github.com/ehudthelefthand/course/gorm"
	"github.com/ehudthelefthand/course/model"
	"github.com/ehudthelefthand/course/server"
	"github.com/gin-gonic/gin"
)

type mockDB struct{}

func (m *mockDB) GetAllCourse() ([]model.Course, error) {
	// return []model.Course{
	// 	{ID: 1, Name: "This a mock", Description: "Mock!"},
	// 	{ID: 2, Name: "This a mock 2", Description: "Mock 2!"},
	// }, nil

	return nil, errors.New("this is an error")
}

func main() {
	gorm, err := gorm.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	// if err := db.Reset(); err != nil {
	// 	log.Fatal(err)
	// }
	if err := gorm.AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	mdb := &mockDB{}

	r := server.Init(mdb, gorm)

	r.Run(":8080")
}

func Error(c *gin.Context, status int, err error) {
	log.Println(err)
	c.JSON(status, gin.H{
		"message": err.Error(),
	})
}
