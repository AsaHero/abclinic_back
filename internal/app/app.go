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
	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
	Config *config.Config
	DB     *postgres.PostgresDB
	server *http.Server
}

func NewApp(cfg *config.Config) *App {
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
		Logger: logger,
		Config: cfg,
		DB:     db,
	}
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	// repo init
	dentistsRepo := postgresql.NewDentistsRepo(a.DB)

	// usecase init
	dentistsUsecase := usecase.NewDentistsUsecase(contextTimeout, dentistsRepo)

	routerArgs := api.RouteArguments{
		Config:          a.Config,
		Logger:          a.Logger,
		DentistsUsecase: dentistsUsecase,
	}

	// router init
	handlers := api.NewRouter(routerArgs)

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
		a.Logger.Error("shutdown server http: %v", zap.Error(err))
	}

	// logger sync
	a.Logger.Sync()
}
