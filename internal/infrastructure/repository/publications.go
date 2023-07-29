package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type Publications interface {
	Create(ctx context.Context, req *entity.Publications) error
	List(ctx context.Context, filter map[string]string) ([]*entity.Publications, error)
	Update(ctx context.Context, req *entity.Publications) error
	Delete(ctx context.Context, filter map[string]string) error
}
