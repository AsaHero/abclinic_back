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
	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type priceListHandler struct {
	config           *config.Config
	logger           *zap.Logger
	enforcer         *casbin.Enforcer
	priceListUsecase usecase.PriceList
}

func NewPriceListHandler(args handlers.HandlerArguments) http.Handler {
	handler := priceListHandler{
		config:           args.Config,
		logger:           args.Logger,
		enforcer:         args.Enforcer,
		priceListUsecase: args.PriceListUsecase,
	}

	policies := [][]string{
		// admin
		{"admin", "/v1/services/{group_id}", "GET"},
		{"admin", "/v1/services", "POST"},
		{"admin", "/v1/services/{id}", "(PUT)|(DELETE)"},
		{"admin", "/v1/services/groups", "GET"},
		{"admin", "/v1/services/groups/{id}", "GET"},
		{"admin", "/v1/services/groups", "POST"},
		{"admin", "/v1/services/groups/{id}", "(PUT)|(DELETE)"},

		// website
		{"website", "/v1/services/{group_id}", "GET"},
		{"website", "/v1/services/groups", "GET"},
		{"website", "/v1/services/groups/{id}", "GET"},

		// secretary
		{"secretary", "/v1/services/{group_id}", "GET"},
		{"secretary", "/v1/services", "POST"},
		{"secretary", "/v1/services/{id}", "(PUT)|(DELETE)"},
		{"secretary", "/v1/services/groups", "GET"},
		{"secretary", "/v1/services/groups/{id}", "GET"},
		{"secretary", "/v1/services/groups", "POST"},
		{"secretary", "/v1/services/groups/{id}", "(PUT)|(DELETE)"},

		// dentist
		{"dentist", "/v1/services/{group_id}", "GET"},
		{"dentist", "/v1/services/groups", "GET"},
		{"dentist", "/v1/services/groups/{id}", "GET"},
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
		r.Get("/{group_id}", handler.GetPriceListByGroup())
		r.Post("/", handler.CreateService())
		r.Put("/{id}", handler.UpdateService())
		r.Delete("/{id}", handler.DeleteService())
		r.Get("/groups", handler.GetGroupList())
		r.Get("/groups/{id}", handler.GetGroup())
		r.Post("/groups", handler.CreateServiceGroup())
		r.Put("/groups/{id}", handler.UpdateServiceGroup())
		r.Delete("/groups/{id}", handler.DeleteServiceGroups())
	})

	return router
}

// GetPriceListByGroup
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
// @Router /v1/services [POST]
// @Summary Create new service
// @Description Create new service
// @Tags Price list
// @Accept json
// @Produce json
// @Param body body models.CreateServiceRequest true "body"
// @Success 200 {object} models.GUIDResponse
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

		response := models.GUIDResponse{
			GUID: guid,
		}

		render.JSON(w, r, response)
	}
}

// UpdateService
// @Security ApiKeyAuth
// @Router /v1/services/{id} [PUT]
// @Summary Update service
// @Description Update service
// @Tags Price list
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body body models.UpdateServiceRequest true "body"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) UpdateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		request := models.UpdateServiceRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		err := h.priceListUsecase.UpdateService(ctx, &entity.Services{
			GUID:  guid,
			Name:  request.Name,
			Price: request.Price,
		})
		if err != nil {
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

// DeleteService
// @Security ApiKeyAuth
// @Router /v1/services/{id} [DELETE]
// @Summary Delete services
// @Description Delete servivies by guid
// @Tags Price list
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) DeleteService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		err := h.priceListUsecase.DeleteService(ctx, guid)
		if err != nil {
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

// GetGroupList
// @Security ApiKeyAuth
// @Router /v1/services/groups [GET]
// @Summary Get services groups
// @Description Get servivies groups
// @Tags Price list
// @Accept json
// @Produce json
// @Success 200 {object} []models.ServicesGroup
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) GetGroupList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		groups, err := h.priceListUsecase.ListServiceGroups(ctx, map[string]string{})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := []models.ServicesGroup{}

		for _, v := range groups {
			response = append(response, models.ServicesGroup{
				GUID: v.GUID,
				Name: v.Name,
			})
		}

		render.JSON(w, r, response)
	}
}

// GetGroup
// @Security ApiKeyAuth
// @Router /v1/services/groups/{id} [GET]
// @Summary Get services group
// @Description Get servivies group
// @Tags Price list
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.ServicesGroup
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) GetGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		groups, err := h.priceListUsecase.ListServiceGroups(ctx, map[string]string{"guid": guid})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := models.ServicesGroup{
			GUID: groups[0].GUID,
			Name: groups[0].Name,
		}

		render.JSON(w, r, response)
	}
}

// CreateService
// @Security ApiKeyAuth
// @Router /v1/services/groups [POST]
// @Summary Create new service group
// @Description Create new service group
// @Tags Price list
// @Accept json
// @Produce json
// @Param body body models.CreateServiceGroupRequest true "body"
// @Success 200 {object} models.GUIDResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) CreateServiceGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.CreateServiceGroupRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		guid, err := h.priceListUsecase.CreateServiceGroup(ctx, &entity.ServiceGroups{
			Name: request.Name,
		})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := models.GUIDResponse{
			GUID: guid,
		}

		render.JSON(w, r, response)
	}
}

// UpdateServiceGroup
// @Security ApiKeyAuth
// @Router /v1/services/groups/{id} [PUT]
// @Summary Update service group
// @Description Update service group
// @Tags Price list
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body body models.CreateServiceGroupRequest true "body"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) UpdateServiceGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		request := models.CreateServiceGroupRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		err := h.priceListUsecase.UpdateServiceGroup(ctx, &entity.ServiceGroups{
			GUID: guid,
			Name: request.Name,
		})
		if err != nil {
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

// DeleteServiceGroups
// @Security ApiKeyAuth
// @Router /v1/services/groups/{id} [DELETE]
// @Summary Delete services group
// @Description Delete servivies group
// @Tags Price list
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h priceListHandler) DeleteServiceGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		err := h.priceListUsecase.DeleteServiceGroup(ctx, guid)
		if err != nil {
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
