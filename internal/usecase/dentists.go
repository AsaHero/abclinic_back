package usecase

import (
	"context"
	"time"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
)

type Denstists interface {
	Get(ctx context.Context, id int64) (*entity.Dentists, error)
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
