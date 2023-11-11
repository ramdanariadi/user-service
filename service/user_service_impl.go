package user

import (
	_ "embed"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/ramdanariadi/grocery-user-service/dto"
	"github.com/ramdanariadi/grocery-user-service/exception"
	"github.com/ramdanariadi/grocery-user-service/model"
	"github.com/ramdanariadi/grocery-user-service/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type UserServiceImpl struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserServiceImpl {
	return &UserServiceImpl{DB: db}
}

func (service UserServiceImpl) Login(requestBody *dto.LoginDTO) *dto.TokenDTO {
	user := model.User{Email: requestBody.Email}
	tx := service.DB.Find(&user)

	if tx.RowsAffected < 1 {
		panic(exception.ValidationException{Message: "UNAUTHORIZED"})
	}
	salt := os.Getenv("SALT")
	requestBody.Password = salt + requestBody.Password + salt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		panic(exception.ValidationException{Message: "UNAUTHORIZED"})
	}

	return &dto.TokenDTO{
		AccessToken:  generateToken(&user, false),
		RefreshToken: generateToken(&user, true),
		User:         &dto.ProfileDTO{UserId: user.Id, Name: user.Username, Username: user.Username, Email: user.Email, MobilePhoneNumber: user.MobilePhoneNumber},
	}
}

func (service UserServiceImpl) Register(reqBody *dto.RegisterDTO) *dto.TokenDTO {
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
	tx := service.DB.Create(&user)
	if tx.Error != nil {
		panic(exception.ValidationException{Message: "REGISTRATION_FAILED"})
	}

	return &dto.TokenDTO{
		AccessToken:  generateToken(&user, false),
		RefreshToken: generateToken(&user, true),
		User:         &dto.ProfileDTO{UserId: user.Id, Name: reqBody.Username, Username: reqBody.Username, Email: reqBody.Email, MobilePhoneNumber: reqBody.MobilePhoneNumber},
	}
}

func (service UserServiceImpl) Get(userId string) *dto.ProfileDTO {
	user := model.User{Id: userId}
	tx := service.DB.Find(&user)
	if tx.RowsAffected < 1 {
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

func (service UserServiceImpl) Update(userId string, dto *dto.ProfileDTO) {
	user := model.User{Id: userId}
	tx := service.DB.Find(&user)
	if tx.RowsAffected < 1 {
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
	save := service.DB.Save(&user)
	utils.PanicIfError(save.Error)
}

func (service UserServiceImpl) Token(reqBody dto.TokenDTO) *dto.TokenDTO {
	log.Printf("Refresh token %s", reqBody.RefreshToken)
	token := VerifyToken(reqBody.RefreshToken)
	claims := token.Claims.(jwt.MapClaims)
	i := claims["userId"]
	user := model.User{Id: i.(string)}
	tx := service.DB.Find(&user)
	if tx.RowsAffected < 1 {
		panic(exception.ValidationException{Message: "UNAUTHORIZED"})
	}
	return &dto.TokenDTO{
		AccessToken:  generateToken(&user, false),
		RefreshToken: generateToken(&user, true),
		User:         &dto.ProfileDTO{UserId: user.Id, Name: user.Username, Username: user.Username, Email: user.Email, MobilePhoneNumber: user.MobilePhoneNumber},
	}
}

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

func generateToken(user *model.User, isRefreshToken bool) string {
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
