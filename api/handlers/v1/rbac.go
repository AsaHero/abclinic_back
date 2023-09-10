package v1

import (
	"encoding/json"
	"errors"
	"net/http"

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

type rbacHandler struct {
	handlers.BaseHandler
	rbacUsecase    usecase.Rbac
	authoerUsecase usecase.Blogs
	logger         *zap.Logger
	config         *config.Config
	enforcer       *casbin.Enforcer
}

func NewRbacHandler(options handlers.HandlerArguments) http.Handler {
	handler := rbacHandler{
		rbacUsecase:    options.RbacUsecase,
		logger:         options.Logger,
		config:         options.Config,
		enforcer:       options.Enforcer,
		authoerUsecase: options.BlogsUsecase,
	}

	policies := [][]string{
		// admin
		{"admin", "/v1/rbac/roles", "GET"},
		{"admin", "/v1/rbac/user", "GET"},
		{"admin", "/v1/rbac/users", "GET"},
		{"admin", "/v1/rbac/user", "POST"},
		{"admin", "/v1/rbac/user/{id}", "PUT"},
		{"admin", "/v1/rbac/user/{id}", "DELETE"},

		// dentist
		{"dentist", "/v1/rbac/roles", "GET"},
		{"dentist", "/v1/rbac/user", "GET"},

		// secretary
		{"secretary", "/v1/rbac/roles", "GET"},
		{"secretary", "/v1/rbac/user", "GET"},
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
		r.Use(middleware.Authorizer(handler.enforcer, handler.logger))
		// roles
		r.Get("/roles", handler.GetRoles())

		// users
		r.Get("/users", handler.GetAllUsers())
		r.Get("/user", handler.GetUserInfo())
		r.Post("/user", handler.CreateUser())
		r.Put("/user/{id}", handler.UpdateUser())
		r.Delete("/user/{id}", handler.DeleteUser())

	})

	return router
}

// GetRoles
// @Security ApiKeyAuth
// @Router /v1/rbac/roles [GET]
// @Summary Get roles
// @Description Get roles
// @Tags Rbac
// @Accept json
// @Produce json
// @Success 200 {object} []models.GetRolesResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h rbacHandler) GetRoles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := models.GetRolesResponse{
			Roles: []string{
				entity.RoleAdmin,
				entity.RoleDentist,
				entity.RoleSecretary,
				entity.RoleWebsite,
			},
		}

		render.JSON(w, r, response)
	}
}

// GetUserInfo
// @Security ApiKeyAuth
// @Router /v1/rbac/user [GET]
// @Summary Get user info
// @Description Get user info
// @Tags Rbac
// @Accept json
// @Produce json
// @Success 200 {object} []models.GetRolesResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h rbacHandler) GetUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		claims, ok := h.GetAuthData(ctx)
		if !ok {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("failed to fetch authentication data"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      "failed to fetch authentication data",
			})
		}

		userID := claims["user_id"]

		user, err := h.rbacUsecase.GetUser(ctx, map[string]string{"guid": userID})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
		}

		response := models.GetUserInfoResponse{
			GUID:      user.GUID,
			Role:      user.Role,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Username:  user.Username,
		}

		render.JSON(w, r, response)
	}
}

// CreateUser
// @Security ApiKeyAuth
// @Router /v1/rbac/user [POST]
// @Summary Create user
// @Description Create user info
// @Tags Rbac
// @Accept json
// @Produce json
// @Param body body models.CreateUserRequest true "body"
// @Success 200 {object} models.GUIDResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h rbacHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.CreateUserRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		guid, err := h.rbacUsecase.CreateUser(ctx, &entity.Users{
			Role:      request.Role,
			Firstname: request.Firstname,
			Lastname:  request.Lastname,
			Username:  request.Username,
			Password:  request.Password,
		})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
		}

		if request.Role == entity.RoleDentist {
			_, err := h.authoerUsecase.CreateAuthors(ctx, &entity.Authors{
				GUID: guid,
				Name: request.Username,
			})
			if err != nil {
				render.Render(w, r, &errorsapi.ErrResponse{
					Err:            err,
					HTTPStatusCode: http.StatusBadRequest,
					ErrorText:      err.Error(),
				})
			}
		}

		response := models.GUIDResponse{
			GUID: guid,
		}

		render.JSON(w, r, response)
	}
}

// UpdateUser
// @Security ApiKeyAuth
// @Router /v1/rbac/user [PUT]
// @Summary Update user info
// @Description Update user info
// @Tags Rbac
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body body models.UpdateUserRequest true "body"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h rbacHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		request := models.UpdateUserRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		err := h.rbacUsecase.UpdateUser(ctx, &entity.Users{
			GUID:      guid,
			Role:      request.Role,
			Firstname: request.Firstname,
			Lastname:  request.Lastname,
			Username:  request.Username,
			Password:  request.Password,
		})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
		}

		render.JSON(w, r, models.Empty{})
	}
}

// GetAllUsers
// @Security ApiKeyAuth
// @Router /v1/rbac/users [GET]
// @Summary Get all users
// @Description Get all users
// @Tags Rbac
// @Accept json
// @Produce json
// @Success 200 {object} []models.GetAllUsersResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h rbacHandler) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		users, err := h.rbacUsecase.ListUsers(ctx, map[string]string{})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
		}

		response := []models.GetAllUsersResponse{}

		for _, v := range users {
			response = append(response, models.GetAllUsersResponse{
				GUID:      v.GUID,
				Role:      v.Role,
				Firstname: v.Firstname,
				Lastname:  v.Lastname,
				Username:  v.Username,
			})
		}

		render.JSON(w, r, response)
	}
}

// UpdateUser
// @Security ApiKeyAuth
// @Router /v1/rbac/user [DELETE]
// @Summary Delete user
// @Description Delete user
// @Tags Rbac
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h rbacHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		err := h.rbacUsecase.DeleteUser(ctx, guid)
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
		}

		render.JSON(w, r, models.Empty{})
	}
}
