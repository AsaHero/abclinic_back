package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	errorsapi "github.com/AsaHero/abclinic/api/errors"
	"github.com/AsaHero/abclinic/api/handlers"
	"github.com/AsaHero/abclinic/api/middleware"
	"github.com/AsaHero/abclinic/api/models"
	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/usecase"
	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type dentistsHandler struct {
	config          *config.Config
	logger          *zap.Logger
	enforcer        *casbin.Enforcer
	dentistsUsecase usecase.Denstists
}

func NewDentistsHandler(args handlers.HandlerArguments) http.Handler {
	handler := dentistsHandler{
		config:          args.Config,
		logger:          args.Logger,
		enforcer:        args.Enforcer,
		dentistsUsecase: args.DentistsUsecase,
	}

	policies := [][]string{
		// admin
		{"admin", "/v1/dentists/{id}", "PUT"},
	}

	for _, v := range policies {
		_, err := handler.enforcer.AddPolicy(v)
		if err != nil {
			handler.logger.Error("error while adding policies to the casbin", zap.Error(err))
			return nil
		}
	}

	handler.enforcer.SavePolicy()

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Get("/", handler.GetDentistsList())
		r.Get("/{id}", handler.GetDentist())

	})
	
	router.Group(func(r chi.Router) {
		r.Use(middleware.Authorizer(handler.enforcer, handler.logger))
		r.Put("/{id}", handler.UpdateDentist())
	})
	
	return router
}

// GetDentists
// @Router /v1/dentists/{id} [GET]
// @Summary Get one dentist
// @Description Get one dentist by ID
// @Tags dentists
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.GetDentistsListResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h dentistsHandler) GetDentist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		param := chi.URLParam(r, "id")

		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			h.logger.Error("error on parsing param to int", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		dentist, err := h.dentistsUsecase.Get(ctx, id)
		if err != nil {
			h.logger.Error("error on GetDentistsList/dentistsUsecase.Get", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := models.GetDentistsListResponse{
			ID:        dentist.ID,
			CloneName: dentist.CloneName,
			Img:       dentist.URL,
			Priority:  dentist.Priority,
			Side:      dentist.Side,
			Name:      dentist.Name,
			Info:      dentist.Info,
		}

		render.JSON(w, r, response)
	}
}

// GetDentistsList
// @Router /v1/dentists [GET]
// @Summary Get dentist list
// @Description List of dentists
// @Tags dentists
// @Accept json
// @Produce json
// @Success 200 {object} []models.GetDentistsListResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h dentistsHandler) GetDentistsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dentists, err := h.dentistsUsecase.List(ctx, map[string]string{})
		if err != nil {
			h.logger.Error("error on GetDentistsList/dentistsUsecase.Get", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := []models.GetDentistsListResponse{}
		for _, v := range dentists {
			response = append(response, models.GetDentistsListResponse{
				ID:        v.ID,
				CloneName: v.CloneName,
				Name:      v.Name,
				Info:      v.Info,
				Img:       v.URL,
				Side:      v.Side,
				Priority:  v.Priority,
			})
		}

		render.JSON(w, r, response)
	}
}

// GetDentistsList
// @Security ApiKeyAuth
// @Router /v1/dentists/{id} [PUT]
// @Summary Update dentist data
// @Description Update dentists data by ID
// @Tags dentists
// @Accept json
// @Produce json
// @Param body body models.UpdateDentistRequest true "body"
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h dentistsHandler) UpdateDentist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.UpdateDentistRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Error("error on decoding request body", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		param := chi.URLParam(r, "id")

		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			h.logger.Error("error on parsing param to int", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		err = h.dentistsUsecase.Update(ctx, &entity.Dentists{
			ID:   id,
			Name: request.Name,
			Info: request.Info,
			URL:  request.Img,
		})
		if err != nil {
			h.logger.Error("error on UpdateDentist/dentistsUsecase.Update", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		render.JSON(w, r, models.Empty{})
	}
}
