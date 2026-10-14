package auth

import (
	"realtime-server/internal/common/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.GetEnv("JWT_SECRET", ""))

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 60).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
