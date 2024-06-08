package types

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return User{}, fmt.Errorf("could not hash the password, try again later")
	}

	return User{
		Username:     registerUser.Username,
		PasswordHash: string(hashedPass),
	}, nil
}

func ValidatePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return false
	}

	return true
}

func CreateToken(user User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)

	secret := "mySecret"

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Errorf("signingString error %w", err)
		return tokenString
	}

	return tokenString
}
