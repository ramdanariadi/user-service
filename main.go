package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ramdanariadi/grocery-user-service/controller"
	"github.com/ramdanariadi/grocery-user-service/exception"
	"github.com/ramdanariadi/grocery-user-service/middleware"
	"github.com/ramdanariadi/grocery-user-service/model"
	"github.com/ramdanariadi/grocery-user-service/setup"
	"github.com/ramdanariadi/grocery-user-service/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	connection, err := setup.NewDbConnection()
	utils.PanicIfError(err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: connection}))
	utils.PanicIfError(err)
	err = db.AutoMigrate(&model.User{})
	utils.LogIfError(err)

	router := gin.Default()
	router.Use(gin.CustomRecovery(exception.Handler))
	userGroup := router.Group("api/v1/user")
	{
		userController := controller.NewUserController(db)
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", userController.Login)
		userGroup.POST("/token", userController.Token)
		userGroup.PUT("", middleware.Middleware, userController.Update)
		userGroup.GET("", middleware.Middleware, userController.Get)
	}

	err = router.Run(":10000")
	utils.LogIfError(err)
}
