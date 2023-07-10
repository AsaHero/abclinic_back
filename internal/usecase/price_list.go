package usecase

import (
	"context"
	"time"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
)

type PriceList interface {
	CreateService(ctx context.Context, req *entity.Services) (string, error)
	ListServices(ctx context.Context, filter map[string]string) ([]*entity.Services, error)
	UpdateService(ctx context.Context, req *entity.Services) error
	DeleteService(ctx context.Context, id string) error
	CreateServiceGroup(ctx context.Context, req *entity.ServiceGroups) (string, error)
	ListServiceGroups(ctx context.Context, filter map[string]string) ([]*entity.ServiceGroups, error)
	UpdateServiceGroup(ctx context.Context, req *entity.ServiceGroups) error
	DeleteServiceGroup(ctx context.Context, id string) error
}

type priceListUsecase struct {
	BaseUsecase
	ctxTimeout    time.Duration
	serviceRepo   repository.Services
	serviceGroups repository.ServiceGroups
}

func NewPriceListUsecase(ctxTimeout time.Duration, serviceRepo repository.Services, serviceGroupdRepo repository.ServiceGroups) PriceList {
	return &priceListUsecase{
		ctxTimeout:    ctxTimeout,
		serviceRepo:   serviceRepo,
		serviceGroups: serviceGroupdRepo,
	}
}

func (u priceListUsecase) CreateService(ctx context.Context, req *entity.Services) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u.beforeCreate(&req.GUID, &req.CreatedAt, &req.UpdateAt)

	return req.GUID, u.serviceRepo.Create(ctx, req)
}
func (u priceListUsecase) ListServices(ctx context.Context, filter map[string]string) ([]*entity.Services, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.serviceRepo.List(ctx, filter)
}
func (u priceListUsecase) UpdateService(ctx context.Context, req *entity.Services) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u.beforeCreate(nil, nil, &req.UpdateAt)

	return u.serviceRepo.Update(ctx, req)
}
func (u priceListUsecase) DeleteService(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.serviceRepo.Delete(ctx, map[string]string{"guid": id})
}
func (u priceListUsecase) CreateServiceGroup(ctx context.Context, req *entity.ServiceGroups) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u.beforeCreate(&req.GUID, &req.CreatedAt, nil)

	return req.GUID, u.serviceGroups.Create(ctx, req)
}
func (u priceListUsecase) ListServiceGroups(ctx context.Context, filter map[string]string) ([]*entity.ServiceGroups, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.serviceGroups.List(ctx, filter)
}
func (u priceListUsecase) UpdateServiceGroup(ctx context.Context, req *entity.ServiceGroups) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.serviceGroups.Update(ctx, req)
}
func (u priceListUsecase) DeleteServiceGroup(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := u.serviceRepo.Delete(ctx, map[string]string{"group_id": id})
	if err != nil {
		if err.Error() != "no sql rows" {
			return err
		}
	}

	err = u.serviceGroups.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
