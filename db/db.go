package db

import "github.com/ehudthelefthand/course/model"

type DB interface {
	GetAllCourse() ([]model.Course, error)
}
