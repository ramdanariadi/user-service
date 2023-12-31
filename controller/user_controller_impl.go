package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ramdanariadi/grocery-user-service/dto"
	"github.com/ramdanariadi/grocery-user-service/exception"
	service "github.com/ramdanariadi/grocery-user-service/usecase"
	"github.com/ramdanariadi/grocery-user-service/utils"
	"gorm.io/gorm"
)

type UserControllerImpl struct {
	UserService service.UserUsecase
}

func NewUserController(db *gorm.DB) UserController {
	return &UserControllerImpl{UserService: service.NewUserService(db)}
}

func (controller *UserControllerImpl) Register(ctx *gin.Context) {
	registerDTO := dto.RegisterDTO{}
	err := ctx.ShouldBind(&registerDTO)
	utils.PanicIfError(err)
	tokenDTO := controller.UserService.Register(ctx, &registerDTO)
	ctx.JSON(200, gin.H{"data": tokenDTO})
}

func (controller *UserControllerImpl) Get(ctx *gin.Context) {
	value, exists := ctx.Get("userId")
	if !exists {
		panic(exception.AuthenticationException{Message: "FORBIDDEN"})
	}
	profileDTO := controller.UserService.Get(ctx, value.(string))
	ctx.JSON(200, gin.H{"data": profileDTO})
}

func (controller *UserControllerImpl) Update(ctx *gin.Context) {
	updateProfileDTO := dto.ProfileDTO{}
	err := ctx.ShouldBind(&updateProfileDTO)
	if err != nil {
		panic(exception.ValidationException{Message: "BAD_REQUEST"})
	}

	value, exists := ctx.Get("userId")
	if !exists {
		panic(exception.AuthenticationException{Message: "FORBIDDEN"})
	}
	controller.UserService.Update(ctx, value.(string), &updateProfileDTO)
	ctx.JSON(200, gin.H{})
}

func (controller *UserControllerImpl) Login(ctx *gin.Context) {
	loginDTO := dto.LoginDTO{}
	err := ctx.ShouldBind(&loginDTO)
	utils.PanicIfError(err)
	tokenDTO := controller.UserService.Login(ctx, &loginDTO)
	ctx.JSON(200, gin.H{"data": tokenDTO})
}

func (controller *UserControllerImpl) Token(ctx *gin.Context) {
	tokenDTO := dto.TokenDTO{}
	err := ctx.ShouldBind(&tokenDTO)
	utils.PanicIfError(err)
	token := controller.UserService.Token(ctx, tokenDTO)
	ctx.JSON(200, gin.H{"data": token})
}
