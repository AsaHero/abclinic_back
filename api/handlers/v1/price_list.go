package v1

import (
	"encoding/json"
	"net/http"

	errorsapi "github.com/AsaHero/abclinic/api/errors"
	"github.com/AsaHero/abclinic/api/handlers"
	"github.com/AsaHero/abclinic/api/models"
	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type priceListHandler struct {
	config           *config.Config
	logger           *zap.Logger
	priceListUsecase usecase.PriceList
}

func NewPriceListHandler(args handlers.HandlerArguments) http.Handler {
	handler := priceListHandler{
		config:           args.Config,
		logger:           args.Logger,
		priceListUsecase: args.PriceListUsecase,
	}

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Get("/{group_id}", handler.GetPriceListByGroup())
		r.Post("/", nil)
		r.Put("/{id}", nil)
		r.Delete("/{id}", nil)
		r.Get("/groups", nil)
		r.Get("/groups/{id}", nil)
		r.Post("/groups", nil)
		r.Put("/groups/{id}", nil)
		r.Delete("/groups/{id}", nil)
	})

	return router
}

// GetPriceListByGroup
// @Router /v1/services/{group_id} [GET]
// @Summary Get services
// @Description Get servivies by group id
// @Tags Price list
// @Accept json
// @Produce json
// @Param group_id path string true "group_id"
// @Success 200 {object} []models.Services
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) GetPriceListByGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		groupID := chi.URLParam(r, "group_id")

		services, err := h.priceListUsecase.ListServices(ctx, map[string]string{"group_id": groupID})
		if err != nil {
			h.logger.Error("error on GetPriceListByGroup/ priceListUsecase.ListServices", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := []models.Services{}

		for _, v := range services {
			response = append(response, models.Services{
				GUID:  v.GUID,
				Name:  v.Name,
				Price: v.Price,
			})
		}

		render.JSON(w, r, response)
	}
}

// CreateService
// @Router /v1/services [POST]
// @Summary Create new service
// @Description Create new service
// @Tags Price list
// @Accept json
// @Produce json
// #Param body body models.CreateServiceRequest true "body"
// @Success 200 {object} models.GIUD
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) CreateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.CreateServiceRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Error("error on decoding request body", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		guid, err := h.priceListUsecase.CreateService(ctx, &entity.Services{
			GroupID: request.GroupID,
			Name:    request.Name,
			Price:   request.Price,
		})
		if err != nil {
			h.logger.Error("error on CreateService/ priceListUsecase.CreateService", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := models.GUID{
			GUID: guid,
		}

		render.JSON(w, r, response)
	}
}
