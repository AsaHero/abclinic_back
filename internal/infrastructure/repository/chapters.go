package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type Chapters interface {
	List(ctx context.Context, filter map[string]string) ([]*entity.Chapters, error)
	Create(ctx context.Context, req *entity.Chapters) error
	Update(ctx context.Context, req *entity.Chapters) error
	Delete(ctx context.Context, id string) error
}
