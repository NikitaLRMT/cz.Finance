package utils

import (
	"errors"
	"time"

	"cz.Finance/backend/models"

	"github.com/dgrijalva/jwt-go"
)

// JWTClaims содержит данные для JWT токена
type JWTClaims struct {
	ID    int64  `json:"id"`
	Email string `json:"email,omitempty"`
	jwt.StandardClaims
}

// GenerateJWT создает новый JWT токен
func GenerateJWT(userID int64, email string, secret string, expiresIn time.Duration) (string, time.Time, error) {
	expirationTime := time.Now().Add(expiresIn)

	claims := &JWTClaims{
		ID:    userID,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ValidateJWT проверяет JWT токен и возвращает данные из него
func ValidateJWT(tokenString string, secret string) (*models.TokenClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неподдерживаемый метод подписи токена")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("неверный токен")
	}

	return &models.TokenClaims{
		UserID: claims.ID,
		Email:  claims.Email,
	}, nil
}
