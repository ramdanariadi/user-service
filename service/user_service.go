package user

import "github.com/ramdanariadi/grocery-user-service/dto"

type Service interface {
	Login(dto *dto.LoginDTO) *dto.TokenDTO
	Register(dto *dto.RegisterDTO) *dto.TokenDTO
	Token(dto dto.TokenDTO) *dto.TokenDTO
	Update(userId string, dto *dto.ProfileDTO)
	Get(userId string) *dto.ProfileDTO
}
