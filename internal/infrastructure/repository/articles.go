package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type Articles interface {
	Create(ctx context.Context, req *entity.Articles) error
	List(ctx context.Context, filter map[string]string) ([]*entity.Articles, error)
	Update(ctx context.Context, req *entity.Articles) error
	Delete(ctx context.Context, filter map[string]string) error
}
