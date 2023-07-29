package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type Categories interface {
	List(ctx context.Context, filter map[string]string) ([]*entity.Categories, error)
	Create(ctx context.Context, req *entity.Categories) error
	Update(ctx context.Context, req *entity.Categories) error
	Delete(ctx context.Context, filter map[string]string) error
}
