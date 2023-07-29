package handlers

import (
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/usecase"
	"go.uber.org/zap"
)

type HandlerArguments struct {
	Config           *config.Config
	Logger           *zap.Logger
	DentistsUsecase  usecase.Denstists
	PriceListUsecase usecase.PriceList
	InfoUsecase      usecase.InfoUsecase
	BlogsUsecase     usecase.Blogs
}
