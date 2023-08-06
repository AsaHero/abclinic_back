package usecase

import (
	"context"
	"time"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/token"
	"github.com/google/uuid"
)

type RefreshToken interface {
	Get(ctx context.Context, refreshToken string) (*entity.RefreshToken, error)
	Create(ctx context.Context, m *entity.RefreshToken) error
	Delete(ctx context.Context, refreshToken string) error
	GenerateToken(ctx context.Context, sub, tokenType, jwtSecret string, accessTTL, refreshTTL time.Duration, optionalFields ...map[string]interface{}) (string, string, error)
}

type refreshTokenService struct {
	ctxTimeout time.Duration
	repo       repository.RefreshTokenRepo
}

func NewRefreshTokenService(ctxTimeout time.Duration, repo repository.RefreshTokenRepo) RefreshToken {
	return &refreshTokenService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (r *refreshTokenService) beforeCreate(m *entity.RefreshToken) error {
	m.GUID = uuid.New().String()
	m.CreatedAt = time.Now().UTC()
	return nil
}

func (r *refreshTokenService) Get(ctx context.Context, refreshToken string) (*entity.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.Get(ctx, refreshToken)
}

func (r *refreshTokenService) Create(ctx context.Context, m *entity.RefreshToken) error {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	r.beforeCreate(m)
	return r.repo.Create(ctx, m)
}

func (r *refreshTokenService) Delete(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.Delete(ctx, refreshToken)
}

func (r *refreshTokenService) GenerateToken(ctx context.Context, sub, userID, jwtSecret string, accessTTL, refreshTTL time.Duration, optionalFields ...map[string]interface{}) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	accessToken, refreshToken, err := token.GenerateToken(sub, userID, jwtSecret, accessTTL, refreshTTL, optionalFields...)
	if err != nil {
		return "", "", err
	}

	m := entity.RefreshToken{
		RefreshToken: refreshToken,
		ExpiryDate:   time.Now().Add(refreshTTL),
	}

	err = r.Create(ctx, &m)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
