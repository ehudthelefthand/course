package db

import (
	"github.com/ehudthelefthand/course/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	db *gorm.DB
}

func NewDB() (*DB, error) {
	url := "host=localhost user=peagolang password=supersecret dbname=peagolang port=54329 sslmode=disable"
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

type Course struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
}

func (db *DB) CreateCourse(c *model.Course) error {
	course := Course{
		Name:        c.Name,
		Description: c.Description,
	}
	if err := db.db.Create(&course).Error; err != nil {
		return err
	}
	c.ID = course.ID
	return nil
}

func (db *DB) GetCourse(id uint) (*model.Course, error) {
	var course Course
	if err := db.db.First(&course, id).Error; err != nil {
		return nil, err
	}
	return &model.Course{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
	}, nil
}

func (db *DB) GetAllCourse() ([]model.Course, error) {
	var courses []Course
	if err := db.db.Find(&courses).Error; err != nil {
		return nil, err
	}

	result := []model.Course{}
	for _, course := range courses {
		result = append(result, model.Course{
			ID:          course.ID,
			Name:        course.Name,
			Description: course.Description,
		})
	}

	return result, nil
}
