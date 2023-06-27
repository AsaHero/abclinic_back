package repository

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
)

type ServiceGroups interface {
	List(ctx context.Context, filter map[string]string) ([]*entity.ServiceGroups, error)
	Create(ctx context.Context, req *entity.ServiceGroups) error
	Update(ctx context.Context, req *entity.ServiceGroups) error
	Delete(ctx context.Context, id string) error
}
