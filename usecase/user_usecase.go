package usecase

import (
	"context"
	"github.com/ramdanariadi/grocery-user-service/dto"
)

type UserUsecase interface {
	Login(ctx context.Context, dto *dto.LoginDTO) *dto.TokenDTO
	Register(ctx context.Context, dto *dto.RegisterDTO) *dto.TokenDTO
	Token(ctx context.Context, dto dto.TokenDTO) *dto.TokenDTO
	Update(ctx context.Context, userId string, dto *dto.ProfileDTO)
	Get(ctx context.Context, userId string) *dto.ProfileDTO
}
