// Package auth для аутентификации пользователя
package auth

import (
	"errors"
	"fmt"
	"gophkeeper/internal"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

const tokenExp = time.Hour * 3
const secretKey = "someSecretSuperKey"

// BuildJWTString получить токен пользвателя, содержащий его ИД
func BuildJWTString(userID uint64) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// GetUserID получение ИД пользвателя из токена
func GetUserID(tokenString string) (uint64, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return 0, nil
		}
		internal.Logger.Infow("error in parse token", "err", err)
		return 0, err
	}

	if !token.Valid {
		return 0, nil
	}

	return claims.UserID, nil
}
