package usecase

import (
	"context"
	// "errors"
	"time"

	"github.com/AsaHero/abclinic/internal/entity"
	errorspkg "github.com/AsaHero/abclinic/internal/errors"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/validation"
)

type Rbac interface {
	UserExists(ctx context.Context, id string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, *entity.Users, error)
	GetUser(ctx context.Context, filter map[string]string) (*entity.Users, error)
	CreateUser(ctx context.Context, req *entity.Users) (string, error)
	ListUsers(ctx context.Context, filter map[string]string) ([]*entity.Users, error)
	UpdateUser(ctx context.Context, req *entity.Users) error
	DeleteUser(ctx context.Context, id string) error
}

type rbacUsecase struct {
	BaseUsecase
	usersRepo  repository.Users
	ctxTimeout time.Duration
}

func NewRbacUsecase(ctxTimeout time.Duration, usersRepo repository.Users) Rbac {
	return &rbacUsecase{
		usersRepo:  usersRepo,
		ctxTimeout: ctxTimeout,
	}
}

func (u rbacUsecase) UserExists(ctx context.Context, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	var isExists bool = true
	_, err := u.usersRepo.Get(ctx, map[string]string{"guid": id})
	if err != nil {
		if err == errorspkg.ErrorNotFound {
			isExists = false
		} else {
			return false, err
		}
	}

	return isExists, nil

}
func (u rbacUsecase) UsernameExists(ctx context.Context, username string) (bool, *entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	var isExists bool = true
	user, err := u.usersRepo.Get(ctx, map[string]string{"username": username})
	if err != nil {
		if err == errorspkg.ErrorNotFound {
			isExists = false
		} else {
			return false, nil, err
		}
	}

	return isExists, user, nil

}
func (u rbacUsecase) GetUser(ctx context.Context, filter map[string]string) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.usersRepo.Get(ctx, filter)
}
func (u rbacUsecase) CreateUser(ctx context.Context, req *entity.Users) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	passwordHash, err := validation.HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	req.Password = passwordHash

	u.BaseUsecase.beforeCreate(&req.GUID, &req.CreatedAt, &req.UpdatedAt)

	return req.GUID, u.usersRepo.Create(ctx, req)
}
func (u rbacUsecase) ListUsers(ctx context.Context, filter map[string]string) ([]*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.usersRepo.List(ctx, filter)
}
func (u rbacUsecase) UpdateUser(ctx context.Context, req *entity.Users) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	passwordHash, err := validation.HashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = passwordHash

	u.BaseUsecase.beforeCreate(nil, nil, &req.UpdatedAt)

	return u.usersRepo.Update(ctx, req)
}
func (u rbacUsecase) DeleteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.usersRepo.Delete(ctx, map[string]string{"guid": id})
}
