package repository

import (
	"context"
	"github.com/ramdanariadi/grocery-user-service/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindById(ctx context.Context, id string) (model.User, error)
	Create(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
}
