package handlers

import (
	"context"

	"github.com/AsaHero/abclinic/api/middleware"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/usecase"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

type HandlerArguments struct {
	Config              *config.Config
	Logger              *zap.Logger
	Enforcer            *casbin.Enforcer
	DentistsUsecase     usecase.Denstists
	PriceListUsecase    usecase.PriceList
	InfoUsecase         usecase.InfoUsecase
	BlogsUsecase        usecase.Blogs
	RbacUsecase         usecase.Rbac
	RefreshTokenUsecase usecase.RefreshToken
}

type BaseHandler struct{}

func (h *BaseHandler) GetAuthData(ctx context.Context) (map[string]string, bool) {
	data, ok := ctx.Value(middleware.CtxKeyAuthData).(map[string]string)
	return data, ok
}
