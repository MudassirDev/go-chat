package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	ISSUER string = "go-chat"
)

func HashPassword(password string) (string, error) {
	rawPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(rawPassword), nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateJWT(userID int64, jwtSecret string, expiryTime time.Duration) (string, error) {
	jwtSecretKey := []byte(jwtSecret)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    ISSUER,
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return jwtToken.SignedString(jwtSecretKey)
}

func VerifyJWT(jwtSecret, token string) (int64, error) {
	claims := jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return 0, err
	}

	if issuer != ISSUER {
		return 0, fmt.Errorf("issuer don't match")
	}

	rawId, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(rawId)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}
