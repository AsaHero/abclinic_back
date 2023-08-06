package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type Users interface {
	Get(ctx context.Context, filter map[string]string) (*entity.Users, error)
	Create(ctx context.Context, req *entity.Users) error
	List(ctx context.Context, filter map[string]string) ([]*entity.Users, error)
	Update(ctx context.Context, req *entity.Users) error
	Delete(ctx context.Context, filter map[string]string) error
}
