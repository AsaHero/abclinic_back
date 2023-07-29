package usecase

import (
	"context"
	"time"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
)

type Blogs interface {
	CreatePublications(ctx context.Context, req *entity.Publications) (string, error)
	ListPublications(ctx context.Context, filter map[string]string) ([]*entity.Publications, error)
	UpdatePublications(ctx context.Context, req *entity.Publications) error
	DeletePublications(ctx context.Context, id string) error
	CreatePublicationsCategories(ctx context.Context, req *entity.Categories) (string, error)
	ListPublicationsCategories(ctx context.Context, filter map[string]string) ([]*entity.Categories, error)
	UpdatePublicationsCategories(ctx context.Context, req *entity.Categories) error
	DeletePublicationsCategories(ctx context.Context, id string) error
	CreateAuthors(ctx context.Context, req *entity.Authors) (string, error)
	ListAuthors(ctx context.Context, filter map[string]string) ([]*entity.Authors, error)
	UpdateAuthors(ctx context.Context, req *entity.Authors) error
	DeleteAuthors(ctx context.Context, id string) error
	GetAuthor(ctx context.Context, id string) (*entity.Authors, error)
}

type blogsUsecase struct {
	BaseUsecase
	ctxTimeout       time.Duration
	publicationsRepo repository.Publications
	categoriesRepo   repository.Categories
	authorsRepo      repository.Authors
}

func NewBlogsUsecase(ctxTimeout time.Duration, publicationsRepo repository.Publications, categoriesRepo repository.Categories, authorsRepo repository.Authors) Blogs {
	return &blogsUsecase{
		ctxTimeout:       ctxTimeout,
		publicationsRepo: publicationsRepo,
		authorsRepo:      authorsRepo,
		categoriesRepo:   categoriesRepo,
	}
}

func (u blogsUsecase) CreatePublications(ctx context.Context, req *entity.Publications) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u.beforeCreate(&req.GUID, &req.CreatedAt, nil)

	return req.GUID, u.publicationsRepo.Create(ctx, req)
}
func (u blogsUsecase) ListPublications(ctx context.Context, filter map[string]string) ([]*entity.Publications, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.publicationsRepo.List(ctx, filter)
}
func (u blogsUsecase) UpdatePublications(ctx context.Context, req *entity.Publications) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.publicationsRepo.Update(ctx, req)
}
func (u blogsUsecase) DeletePublications(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.publicationsRepo.Delete(ctx, map[string]string{"guid": id})
}
func (u blogsUsecase) CreatePublicationsCategories(ctx context.Context, req *entity.Categories) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u.beforeCreate(&req.GUID, &req.CreatedAt, nil)

	return req.GUID, u.categoriesRepo.Create(ctx, req)
}
func (u blogsUsecase) ListPublicationsCategories(ctx context.Context, filter map[string]string) ([]*entity.Categories, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.categoriesRepo.List(ctx, filter)
}
func (u blogsUsecase) UpdatePublicationsCategories(ctx context.Context, req *entity.Categories) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.categoriesRepo.Update(ctx, req)
}
func (u blogsUsecase) DeletePublicationsCategories(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := u.publicationsRepo.Delete(ctx, map[string]string{"group_id": id})
	if err != nil {
		if err.Error() != "no sql rows" {
			return err
		}
	}

	err = u.categoriesRepo.Delete(ctx, map[string]string{"guid": id})
	if err != nil {
		return err
	}

	return nil
}
func (u blogsUsecase) CreateAuthors(ctx context.Context, req *entity.Authors) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u.beforeCreate(&req.GUID, &req.CreatedAt, nil)

	return req.GUID, u.authorsRepo.Create(ctx, req)
}

func (u blogsUsecase) GetAuthor(ctx context.Context, id string) (*entity.Authors, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.authorsRepo.Get(ctx, id)
}

func (u blogsUsecase) ListAuthors(ctx context.Context, filter map[string]string) ([]*entity.Authors, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.authorsRepo.List(ctx, filter)
}
func (u blogsUsecase) UpdateAuthors(ctx context.Context, req *entity.Authors) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.authorsRepo.Update(ctx, req)
}
func (u blogsUsecase) DeleteAuthors(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := u.publicationsRepo.Delete(ctx, map[string]string{"author_id": id})
	if err != nil {
		if err.Error() != "no sql rows" {
			return err
		}
	}

	err = u.authorsRepo.Delete(ctx, map[string]string{"guid": id})
	if err != nil {
		return err
	}

	return nil
}
