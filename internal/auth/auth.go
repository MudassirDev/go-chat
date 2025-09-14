package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func CreateJWT(userID uuid.UUID, jwtSecret string, expiryTime time.Duration) (string, error) {
	jwtSecretKey := []byte(jwtSecret)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    ISSUER,
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return jwtToken.SignedString(jwtSecretKey)
}

func VerifyJWT(jwtSecret, token string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != ISSUER {
		return uuid.Nil, fmt.Errorf("issuer don't match")
	}

	rawId, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(rawId)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
