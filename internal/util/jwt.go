package util

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-strong-secret-key")

func GenerateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(authHeader string) (int64, error) {
	log.Printf("Попытка разбора токена: %s", authHeader)

	// 1. Разделяем строку по пробелу, чтобы отделить "Bearer" от токена
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, fmt.Errorf("неверный формат заголовка авторизации")
	}
	tokenString := parts[1]

	// 2. Парсим сам токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("Ошибка разбора токена: %v", err)
		return 0, err
	}

	// 3. Проверяем, что токен валиден и извлекаем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Claims["user_id"] приходит как float64, нужно конвертировать
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			return int64(userIDFloat), nil
		}
	}

	return 0, fmt.Errorf("невалидный токен")
}
