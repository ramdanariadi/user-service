package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/ramdanariadi/grocery-user-service/exception"
	"github.com/ramdanariadi/grocery-user-service/model"
	"log"
	"os"
	"time"
)

func VerifyToken(tokenStr string) *jwt.Token {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("INVALID_ALGORITHM")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		panic(exception.AuthenticationException{Message: "UNAUTHORIZED"})
	}
	return token
}

func GenerateToken(user *model.User, isRefreshToken bool) string {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	if isRefreshToken {
		claims["exp"] = time.Now().Add(48 * time.Hour).UnixNano()
	} else {
		claims["exp"] = time.Now().Add(10 * time.Minute).UnixNano()
	}
	//claims["authorized"] = true
	claims["userId"] = user.Id

	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("key invalid %s", secret)
	}
	return signedString
}
