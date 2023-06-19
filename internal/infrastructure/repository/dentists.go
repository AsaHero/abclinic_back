package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type Denstists interface {
	Get(ctx context.Context, id int64) (*entity.Dentists, error)
	List(ctx context.Context, filter map[string]string) ([]*entity.Dentists, error)
	Update(ctx context.Context, req *entity.Dentists) error
}
