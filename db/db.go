package db

import (
	"log"
	"time"

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

func (db *DB) Reset() error {
	err := db.db.Migrator().DropTable(
		&User{},
		&Course{},
		&Class{},
		&ClassStudent{},
	)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (db *DB) AutoMigrate() error {
	err := db.db.Migrator().AutoMigrate(
		&User{},
		&Course{},
		&Class{},
		&ClassStudent{},
	)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

type Course struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
}

type Class struct {
	ID        uint `gorm:"primaryKey"`
	CourseID  uint
	Course    Course
	TrainerID uint
	Trainer   User
	Start     time.Time
	End       time.Time
	Seats     int
	Students  []ClassStudent
}

type ClassStudent struct {
	ID        uint `gorm:"primaryKey"`
	ClassID   uint
	StudentID uint
	Student   User
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string
}

func (db *DB) CreateUser(u *model.User) error {
	user := User{
		Username: u.Username,
		Password: u.Password,
	}
	if err := db.db.Create(&user).Error; err != nil {
		return err
	}
	u.ID = user.ID
	return nil
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

func (db *DB) SaveClass(cls *model.Class) error {
	class := Class{
		CourseID:  cls.Course.ID,
		TrainerID: cls.Trainer.ID,
		Start:     cls.Start,
		End:       cls.End,
		Seats:     cls.Seats,
	}
	if err := db.db.Save(&class).Error; err != nil {
		return err
	}
	cls.ID = class.ID
	return nil
}

func (db *DB) GetClass(id uint) (*model.Class, error) {
	var class Class
	err := db.db.Preload("Course").
		Preload("Trainer").
		Preload("Students.Student").
		First(&class, id).Error
	if err != nil {
		return nil, err
	}

	students := []model.Student{}
	for _, stu := range class.Students {
		students = append(students, model.Student{
			ID:   stu.StudentID,
			Name: stu.Student.Username,
		})
	}

	return &model.Class{
		ID: class.ID,
		Course: model.Course{
			ID:          class.Course.ID,
			Name:        class.Course.Name,
			Description: class.Course.Description,
		},
		Trainer: model.Trainer{
			ID:   class.Trainer.ID,
			Name: class.Trainer.Username,
		},
		Start:    class.Start,
		End:      class.End,
		Seats:    class.Seats,
		Students: students,
	}, nil
}

func (db *DB) GetStudent(id uint) (*model.Student, error) {
	var student User
	if err := db.db.First(&student, id).Error; err != nil {
		return nil, err
	}
	return &model.Student{
		ID:   student.ID,
		Name: student.Username,
	}, nil
}

func (db *DB) CreateClassStudent(studentID uint, classID uint) error {
	classStudent := ClassStudent{
		StudentID: studentID,
		ClassID:   classID,
	}
	return db.db.Create(&classStudent).Error
}

func (db *DB) GetUserByUsername(username string) (*model.User, error) {
	var user User
	err := db.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (db *DB) GetUserByID(id uint) (*model.User, error) {
	var user User
	err := db.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
