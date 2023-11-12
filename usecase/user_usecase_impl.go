package usecase

import (
	"context"
	_ "embed"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/ramdanariadi/grocery-user-service/dto"
	"github.com/ramdanariadi/grocery-user-service/exception"
	"github.com/ramdanariadi/grocery-user-service/model"
	"github.com/ramdanariadi/grocery-user-service/repository"
	"github.com/ramdanariadi/grocery-user-service/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type userUsecaseImpl struct {
	Repository repository.UserRepository
}

func NewUserService(db *gorm.DB) UserUsecase {
	second := time.Second * 5
	return &userUsecaseImpl{Repository: repository.NewUserRepositoryImpl(db, second)}
}

func (service userUsecaseImpl) Login(ctx context.Context, requestBody *dto.LoginDTO) *dto.TokenDTO {
	user, err := service.Repository.FindByEmail(ctx, requestBody.Email)
	if err != nil {
		panic(exception.ValidationException{Message: "UNAUTHORIZED"})
	}
	salt := os.Getenv("SALT")
	requestBody.Password = salt + requestBody.Password + salt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		panic(exception.ValidationException{Message: "UNAUTHORIZED"})
	}

	return &dto.TokenDTO{
		AccessToken:  utils.GenerateToken(&user, false),
		RefreshToken: utils.GenerateToken(&user, true),
		User:         &dto.ProfileDTO{UserId: user.Id, Name: user.Username, Username: user.Username, Email: user.Email, MobilePhoneNumber: user.MobilePhoneNumber},
	}
}

func (service userUsecaseImpl) Register(ctx context.Context, reqBody *dto.RegisterDTO) *dto.TokenDTO {
	salt := os.Getenv("SALT")
	reqBody.Password = salt + reqBody.Password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
	utils.LogIfError(err)

	id, err := uuid.NewUUID()
	utils.PanicIfError(err)
	user := model.User{
		Id:                id.String(),
		Email:             reqBody.Email,
		Password:          string(hashedPassword),
		Username:          reqBody.Username,
		MobilePhoneNumber: reqBody.MobilePhoneNumber,
	}
	err = service.Repository.Create(ctx, user)
	if err != nil {
		panic(exception.ValidationException{Message: "REGISTRATION_FAILED"})
	}

	return &dto.TokenDTO{
		AccessToken:  utils.GenerateToken(&user, false),
		RefreshToken: utils.GenerateToken(&user, true),
		User:         &dto.ProfileDTO{UserId: user.Id, Name: reqBody.Username, Username: reqBody.Username, Email: reqBody.Email, MobilePhoneNumber: reqBody.MobilePhoneNumber},
	}
}

func (service userUsecaseImpl) Get(ctx context.Context, userId string) *dto.ProfileDTO {
	user, err := service.Repository.FindById(ctx, userId)
	if err != nil {
		panic(exception.ValidationException{Message: "UNAUTHORIZED"})
	}

	profileDTO := dto.ProfileDTO{
		UserId:            user.Id,
		Name:              user.Name,
		Username:          user.Username,
		Email:             user.Email,
		MobilePhoneNumber: user.MobilePhoneNumber,
		ProfileImageUrl:   &user.ProfileImageUrl,
	}
	return &profileDTO
}

func (service userUsecaseImpl) Update(ctx context.Context, userId string, dto *dto.ProfileDTO) {
	user, err := service.Repository.FindById(ctx, userId)
	if err != nil {
		panic(exception.ValidationException{Message: "UNAUTHORIZED"})
	}

	log.Printf("user id %s", userId)
	log.Printf("name %s", dto.Name)
	log.Printf("mobile phone number %s", dto.MobilePhoneNumber)
	log.Printf("username %s", dto.Username)
	user.Name = dto.Name
	user.MobilePhoneNumber = dto.MobilePhoneNumber
	user.Email = dto.Email
	user.Username = dto.Username
	if dto.ProfileImageUrl != nil {
		user.ProfileImageUrl = *dto.ProfileImageUrl
	}
	err = service.Repository.Create(ctx, user)
	utils.PanicIfError(err)
}

func (service userUsecaseImpl) Token(ctx context.Context, reqBody dto.TokenDTO) *dto.TokenDTO {
	log.Printf("Refresh token %s", reqBody.RefreshToken)
	token := utils.VerifyToken(reqBody.RefreshToken)
	claims := token.Claims.(jwt.MapClaims)
	i := claims["userId"]

	user, err := service.Repository.FindById(ctx, i.(string))
	if err != nil {
		panic(exception.ValidationException{Message: "UNAUTHORIZED"})
	}
	return &dto.TokenDTO{
		AccessToken:  utils.GenerateToken(&user, false),
		RefreshToken: utils.GenerateToken(&user, true),
		User:         &dto.ProfileDTO{UserId: user.Id, Name: user.Username, Username: user.Username, Email: user.Email, MobilePhoneNumber: user.MobilePhoneNumber},
	}
}
