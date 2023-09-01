package usecase

import (
	"context"
	"time"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
)

type InfoUsecase interface {
	CreateArticle(ctx context.Context, req *entity.Articles) (string, error)
	ListArticles(ctx context.Context, filter map[string]string) ([]*entity.Articles, error)
	UpdateArticles(ctx context.Context, req *entity.Articles) error
	DeleteArticles(ctx context.Context, id string) error
	CreateArticlesChapter(ctx context.Context, req *entity.Chapters) (string, error)
	ListArticlesChapters(ctx context.Context, filter map[string]string) ([]*entity.Chapters, error)
	UpdateArticlesChapter(ctx context.Context, req *entity.Chapters) error
	DeleteArticlesChapter(ctx context.Context, id string) error
}

type infoUsecase struct {
	BaseUsecase
	ctxTimeout   time.Duration
	articlesRepo repository.Articles
	chaptersRepo repository.Chapters
}

func NewinfoUsecase(ctxTimeout time.Duration, articlesRepo repository.Articles, chaptersRepo repository.Chapters) InfoUsecase {
	return &infoUsecase{
		ctxTimeout:   ctxTimeout,
		articlesRepo: articlesRepo,
		chaptersRepo: chaptersRepo,
	}
}

func (u infoUsecase) CreateArticle(ctx context.Context, req *entity.Articles) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u.beforeCreate(&req.GUID, &req.CreatedAt, nil)

	return req.GUID, u.articlesRepo.Create(ctx, req)
}
func (u infoUsecase) ListArticles(ctx context.Context, filter map[string]string) ([]*entity.Articles, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.articlesRepo.List(ctx, filter)
}
func (u infoUsecase) UpdateArticles(ctx context.Context, req *entity.Articles) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.articlesRepo.Update(ctx, req)
}
func (u infoUsecase) DeleteArticles(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.articlesRepo.Delete(ctx, map[string]string{"guid": id})
}
func (u infoUsecase) CreateArticlesChapter(ctx context.Context, req *entity.Chapters) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u.beforeCreate(&req.GUID, &req.CreatedAt, nil)

	return req.GUID, u.chaptersRepo.Create(ctx, req)
}
func (u infoUsecase) ListArticlesChapters(ctx context.Context, filter map[string]string) ([]*entity.Chapters, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.chaptersRepo.List(ctx, filter)
}
func (u infoUsecase) UpdateArticlesChapter(ctx context.Context, req *entity.Chapters) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return u.chaptersRepo.Update(ctx, req)
}
func (u infoUsecase) DeleteArticlesChapter(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := u.articlesRepo.Delete(ctx, map[string]string{"chapter_id": id})
	if err != nil {
		if err.Error() != "no sql rows" {
			return err
		}
	}

	err = u.chaptersRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
