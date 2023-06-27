package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "github.com/AsaHero/abclinic/api/docs"
	"github.com/AsaHero/abclinic/api/handlers"
	v1 "github.com/AsaHero/abclinic/api/handlers/v1"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/usecase"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type RouteArguments struct {
	Config           *config.Config
	Logger           *zap.Logger
	DentistsUsecase  usecase.Denstists
	PriceListUsecase usecase.PriceList
}

// NewRoute
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(args RouteArguments) http.Handler {
	handlersArgs := handlers.HandlerArguments{
		Config:           args.Config,
		Logger:           args.Logger,
		DentistsUsecase:  args.DentistsUsecase,
		PriceListUsecase: args.PriceListUsecase,
	}

	router := chi.NewRouter()
	router.Use(chimiddleware.RealIP, chimiddleware.Logger, chimiddleware.Recoverer)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/dentists", v1.NewDentistsHandler(handlersArgs))
		r.Mount("/services", v1.NewPriceListHandler(handlersArgs))
		r.Mount("/categorises", http.NotFoundHandler())
	})

	// declare swagger api route
	router.Get("/swagger/*", httpSwagger.Handler())
	return router
}
