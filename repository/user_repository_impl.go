package repository

import (
	"context"
	"errors"
	"github.com/ramdanariadi/grocery-user-service/model"
	"gorm.io/gorm"
	"time"
)

type userRepositoryImpl struct {
	DB             *gorm.DB
	ContextTimeOut time.Duration
}

func NewUserRepositoryImpl(DB *gorm.DB, contextTimeout time.Duration) UserRepository {
	return &userRepositoryImpl{DB: DB, ContextTimeOut: contextTimeout}
}

func (repository userRepositoryImpl) FindByEmail(ctx context.Context, email string) (user model.User, e error) {
	ctx, cancelFunc := context.WithTimeout(ctx, repository.ContextTimeOut)
	defer cancelFunc()
	user.Email = email
	if repository.DB.WithContext(ctx).Find(&user).RowsAffected < 1 {
		e = errors.New("USER_NOT_FOUND")
	}

	return user, e
}

func (repository userRepositoryImpl) FindById(ctx context.Context, id string) (user model.User, e error) {
	ctx, cancelFunc := context.WithTimeout(ctx, repository.ContextTimeOut)
	defer cancelFunc()
	user.Id = id
	if repository.DB.WithContext(ctx).Find(&user).RowsAffected < 1 {
		e = errors.New("USER_NOT_FOUND")
	}

	return user, e
}

func (repository userRepositoryImpl) Create(ctx context.Context, user model.User) error {
	ctx, cancelFunc := context.WithTimeout(ctx, repository.ContextTimeOut)
	defer cancelFunc()

	return repository.DB.Save(&user).Error
}

func (repository userRepositoryImpl) Update(ctx context.Context, user model.User) error {
	ctx, cancelFunc := context.WithTimeout(ctx, repository.ContextTimeOut)
	defer cancelFunc()

	return repository.DB.Save(&user).Error
}
