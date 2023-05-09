package handler

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const secretKey = "SuperSecret"

type claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func generateToken(userID uint) (string, error) {
	payload := claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    "course-api",
		},
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return claim.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (*claims, error) {
	// Key Func will verify signing method
	// and return secretKey if signing method is match
	// otherwise return error
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("signing method error")
		}
		return []byte(secretKey), nil
	}
	// Parse Claims
	jwtToken, err := jwt.ParseWithClaims(token, &claims{}, keyFunc)
	if err != nil {
		return nil, err
	}
	// Check Claim Type
	claims, ok := jwtToken.Claims.(*claims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
