package api

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "github.com/AsaHero/abclinic/api/docs"
	"github.com/AsaHero/abclinic/api/handlers"
	v1 "github.com/AsaHero/abclinic/api/handlers/v1"
	"github.com/AsaHero/abclinic/api/middleware"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/usecase"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type RouteArguments struct {
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

// NewRoute
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(args RouteArguments) http.Handler {
	handlersArgs := handlers.HandlerArguments{
		Config:              args.Config,
		Logger:              args.Logger,
		Enforcer:            args.Enforcer,
		DentistsUsecase:     args.DentistsUsecase,
		PriceListUsecase:    args.PriceListUsecase,
		InfoUsecase:         args.InfoUsecase,
		BlogsUsecase:        args.BlogsUsecase,
		RbacUsecase:         args.RbacUsecase,
		RefreshTokenUsecase: args.RefreshTokenUsecase,
	}

	router := chi.NewRouter()
	router.Use(chimiddleware.RealIP, chimiddleware.Logger, chimiddleware.Recoverer)
	// router.Use(chimiddleware.Timeout(args.ContextTimeout))
	router.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Route("/v1", func(r chi.Router) {
		r.Use(middleware.AuthContext(args.Config.Token.Secret))
		r.Mount("/", v1.NewAuthHandler(handlersArgs))
		r.Mount("/dentists", v1.NewDentistsHandler(handlersArgs))
		r.Mount("/services", v1.NewPriceListHandler(handlersArgs))
		r.Mount("/articles", v1.NewInfoHandler(handlersArgs))
		r.Mount("/blogs", v1.NewBlogsHandler(handlersArgs))
		r.Mount("/authors", v1.NewAuthorsHandler(handlersArgs))
		r.Mount("/file", v1.NewFilesHandler(handlersArgs))
		r.Mount("/rbac", v1.NewRbacHandler(handlersArgs))
	})

	// declare swagger api route
	router.Get("/swagger/*", httpSwagger.Handler())
	return router
}
