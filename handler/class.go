package handler

import (
	"net/http"
	"time"

	"github.com/ehudthelefthand/course/db"
	"github.com/ehudthelefthand/course/util"
	"github.com/gin-gonic/gin"
)

type ClassReq struct {
	CourseID  uint      `json:"course_id"`
	TrainerID uint      `json:"trainer_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Seats     int       `json:"seats"`
}

func CreateClasses(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(ClassReq)
		if err := c.BindJSON(req); err != nil {
			util.SendError(c, http.StatusBadRequest, err)
			return
		}
		course, err := db.GetCourse(req.CourseID)
		if err != nil {
			util.SendError(c, http.StatusNotFound, err)
			return
		}
		class, err := course.CreateClass(req.Start, req.End)
		if err != nil {
			util.SendError(c, http.StatusBadRequest, err)
			return
		}
		if err := class.SetSeats(req.Seats); err != nil {
			util.SendError(c, http.StatusBadRequest, err)
			return
		}

		class.Trainer.ID = req.TrainerID

		if err := db.SaveClass(class); err != nil {
			util.SendError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
