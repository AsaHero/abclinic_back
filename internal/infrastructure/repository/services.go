package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type Services interface {
	Create(ctx context.Context, req *entity.Services) error
	List(ctx context.Context, filter map[string]string) ([]*entity.Services, error)
	Update(ctx context.Context, req *entity.Services) error
	Delete(ctx context.Context, filter map[string]string) error
}
