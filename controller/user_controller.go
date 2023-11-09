package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Token(ctx *gin.Context)
	Update(ctx *gin.Context)
	Get(ctx *gin.Context)
}
