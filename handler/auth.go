package handler

import (
	"net/http"

	"github.com/ehudthelefthand/course/db"
	"github.com/ehudthelefthand/course/model"
	"github.com/ehudthelefthand/course/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 12
)

type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(dbx *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(AuthReq)
		if err := c.BindJSON(req); err != nil {
			util.SendError(c, http.StatusBadRequest, err)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), cost)
		if err != nil {
			util.SendError(c, http.StatusInternalServerError, err)
			return
		}

		user := model.User{
			Username: req.Username,
			Password: string(hash),
			Role:     db.Student,
		}
		if err := dbx.CreateUser(&user); err != nil {
			util.SendError(c, http.StatusInternalServerError, err)
			return
		}

		token, err := generateToken(user.ID)
		if err != nil {
			util.SendError(c, http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

func Login(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(AuthReq)
		if err := c.BindJSON(req); err != nil {
			util.SendError(c, http.StatusBadRequest, err)
			return
		}
		found, err := db.GetUserByUsername(req.Username)
		if found == nil || err != nil {
			util.SendError(c, http.StatusUnauthorized, err)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(req.Password))
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, err)
			return
		}
		token, err := generateToken(found.ID)
		if err != nil {
			util.SendError(c, http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
