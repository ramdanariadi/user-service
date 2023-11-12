package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ramdanariadi/grocery-user-service/exception"
	"github.com/ramdanariadi/grocery-user-service/utils"
	"strings"
)

func Middleware(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	split := strings.Split(header, " ")
	if len(split) < 2 {
		panic(exception.AuthenticationException{Message: "UNAUTHORIZED"})
	}

	if strings.Compare(split[0], "Bearer") != 0 {
		panic(exception.AuthenticationException{Message: "UNAUTHORIZED"})
	}

	token := utils.VerifyToken(split[1])

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"]
	ctx.Set("userId", userId.(string))
	ctx.Next()
}
