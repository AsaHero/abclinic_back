package usecase

import (
	"context"
	"time"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
)

type Denstists interface {
	Get(ctx context.Context, id int64) (*entity.Dentists, error)
	List(ctx context.Context, filter map[string]string) ([]*entity.Dentists, error)
	Update(ctx context.Context, req *entity.Dentists) error
}

type dentistsUsecase struct {
	ctxTimeout   time.Duration
	dentistsRepo repository.Denstists
}

func NewDentistsUsecase(ctxTimeout time.Duration, dentistsRepo repository.Denstists) Denstists {
	return &dentistsUsecase{
		ctxTimeout:   ctxTimeout,
		dentistsRepo: dentistsRepo,
	}
}

func (u *dentistsUsecase) Get(ctx context.Context, id int64) (*entity.Dentists, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.dentistsRepo.Get(ctx, id)
}

func (u *dentistsUsecase) List(ctx context.Context, filter map[string]string) ([]*entity.Dentists, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.dentistsRepo.List(ctx, filter)
}

func (u *dentistsUsecase) Update(ctx context.Context, req *entity.Dentists) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.dentistsRepo.Update(ctx, req)
}
