package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	fmt.Printf("Hashing password: %s\n", password)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Printf("Generated hash: %s\n", string(bytes))
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
