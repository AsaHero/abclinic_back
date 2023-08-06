package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AsaHero/abclinic/api"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository/postgresql"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/pkg/logger"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
	"github.com/AsaHero/abclinic/internal/usecase"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

type App struct {
	Logger   *zap.Logger
	Config   *config.Config
	DB       *postgres.PostgresDB
	server   *http.Server
	Enforcer *casbin.Enforcer
}

func NewApp(cfg *config.Config) *App {

	enforcer, err := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
	if err != nil {
		log.Fatalf("error on casbin init: %v", err)
	}

	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		log.Fatalf("error on logger init: %v", err)
	}

	// db init
	db, err := postgres.NewPostgresDB(*cfg)
	if err != nil {
		log.Fatalf("error on db init: %v", err)
	}

	return &App{
		Logger:   logger,
		Config:   cfg,
		DB:       db,
		Enforcer: enforcer,
	}
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	// repo init
	dentistsRepo := postgresql.NewDentistsRepo(a.DB)
	serviceRepo := postgresql.NewServicesRepo(a.DB)
	serviceGroupdRepo := postgresql.NewServiceGroupsRepo(a.DB)
	artcileRepo := postgresql.NewArticlesRepo(a.DB)
	chapterRepo := postgresql.NewChaptersRepo(a.DB)
	authorsRepo := postgresql.NewAuthorsRepo(a.DB)
	categoriesRepo := postgresql.NewCategoriesRepo(a.DB)
	publicationsRepo := postgresql.NewPublicationsRepo(a.DB)
	userRepo := postgresql.NewUsersRepo(a.DB)
	refreshTokenRepo := postgresql.NewRefreshTokenRepo(a.DB)

	// usecase init
	dentistsUsecase := usecase.NewDentistsUsecase(contextTimeout, dentistsRepo)
	priceListUsecase := usecase.NewPriceListUsecase(contextTimeout, serviceRepo, serviceGroupdRepo)
	infoUsecase := usecase.NewinfoUsecase(contextTimeout, artcileRepo, chapterRepo)
	blogsUsecase := usecase.NewBlogsUsecase(contextTimeout, publicationsRepo, categoriesRepo, authorsRepo)
	rbacUsecase := usecase.NewRbacUsecase(contextTimeout, userRepo)
	refreshTokenUsecase := usecase.NewRefreshTokenService(contextTimeout, refreshTokenRepo)

	routerArgs := api.RouteArguments{
		Config:              a.Config,
		Logger:              a.Logger,
		Enforcer:            a.Enforcer,
		DentistsUsecase:     dentistsUsecase,
		PriceListUsecase:    priceListUsecase,
		InfoUsecase:         infoUsecase,
		BlogsUsecase:        blogsUsecase,
		RbacUsecase:         rbacUsecase,
		RefreshTokenUsecase: refreshTokenUsecase,
	}

	// router init
	handlers := api.NewRouter(routerArgs)

	if err = a.Enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("error during enforcer load policy: %w", err)
	}

	// server init
	a.server, err = api.NewServer(a.Config, handlers)
	if err != nil {
		return fmt.Errorf("error while server init: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close db pool
	a.DB.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http", zap.Error(err))
	}

	// logger sync
	a.Logger.Sync()
}
