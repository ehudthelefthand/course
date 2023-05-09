package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ehudthelefthand/course/gorm"
	"github.com/ehudthelefthand/course/handler"
	"github.com/ehudthelefthand/course/util"
	"github.com/gin-gonic/gin"
)

func RequireUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get header
		header := c.GetHeader("Authorization")
		header = strings.TrimSpace(header)
		min := len("Bearer ")
		if len(header) <= min {
			util.SendError(c, http.StatusUnauthorized, errors.New("token is require"))
			return
		}
		token := header[min:]
		claims, err := handler.VerifyToken(token)
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, err)
			return
		}
		user, err := db.GetUserByID(claims.UserID)
		if user == nil || err != nil {
			util.SendError(c, http.StatusUnauthorized, err)
			return
		}
		handler.SetUser(c, user)
	}
}
