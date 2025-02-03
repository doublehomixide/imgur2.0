package jwtutils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UtilsJWT struct {
}

var jwtSecret = []byte("xyeettttttttttttttttttta")

func (u *UtilsJWT) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (u *UtilsJWT) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неизвестный метод подписи")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, err
}

func (u *UtilsJWT) ValidateAndExtractPayload(tokenString string) (jwt.MapClaims, error) {
	token, err := u.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("cant get claims")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("exp claim is missing or invalid")
	}

	if int64(exp) <= time.Now().Unix() {
		return nil, fmt.Errorf("token is expired")
	}
	return claims, nil
}
