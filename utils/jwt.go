package utils

import (
	"errors"
	"klubRanks/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"userId":   userId,
		"exp":      time.Now().Add(config.AppConfig.JWT.Expiry).Unix(),
	})
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return 0, errors.New("token is expired or invalid")
	}

	isValid := parsedToken.Valid

	if !isValid {
		return 0, errors.New("invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// email, _ := claims["email"].(string)
	userId := int64(claims["userId"].(float64))

	return userId, nil
}
