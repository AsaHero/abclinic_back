package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type Authors interface {
	Create(ctx context.Context, req *entity.Authors) error
	Get(ctx context.Context, guid string) (*entity.Authors, error)
	List(ctx context.Context, filter map[string]string) ([]*entity.Authors, error)
	Update(ctx context.Context, req *entity.Authors) error
	Delete(ctx context.Context, filter map[string]string) error
}
