package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	errorsapi "github.com/AsaHero/abclinic/api/errors"
	"github.com/AsaHero/abclinic/api/handlers"
	"github.com/AsaHero/abclinic/api/models"
	errorspkg "github.com/AsaHero/abclinic/internal/errors"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/pkg/token"
	"github.com/AsaHero/abclinic/internal/pkg/validation"
	"github.com/AsaHero/abclinic/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type authHandler struct {
	rbacUsecase          usecase.Rbac
	reshreshTokenUsecase usecase.RefreshToken
	logger               *zap.Logger
	config               *config.Config
}

func NewAuthHandler(option handlers.HandlerArguments) http.Handler {
	handler := authHandler{
		rbacUsecase: option.RbacUsecase,
		logger:      option.Logger,
		config:      option.Config,
	}

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		// public apis
		r.Post("/login", handler.Login())
		r.Post("/reshresh", handler.RefreshToken())

	})
	return router
}

// Login
// @Router /v1/login [POST]
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "body"
// @Success 200 {object} models.LoginResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h authHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.LoginRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("cannot parse request body"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      "cannot parse request body",
			})
			return
		}

		isExists, user, err := h.rbacUsecase.UsernameExists(ctx, request.Username)
		if err != nil {
			h.logger.Error("error on Login/ rbacUsecase.UsernameExists", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		if !isExists {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("invalid username"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      "invalid username",
			})
			return
		}

		ok := validation.CheckPasswordHash(request.Password, user.Password)
		if !ok {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("invalid password"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      "invalid password",
			})
			return
		}

		access, refresh, err := h.reshreshTokenUsecase.GenerateToken(ctx, user.Role, user.GUID, h.config.Token.Secret, h.config.Token.AccessTTL, h.config.Token.RefreshTTL)
		if err != nil {
			h.logger.Error("error on Login/ token.GenerateToken", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := models.LoginResponse{
			AccessToken:  access,
			RefreshToken: refresh,
		}

		render.JSON(w, r, response)
	}
}

func (h authHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody := models.RefreshTokenRequest{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			h.logger.Error("investorHandler/RefreshToken/Decode", zap.Error(err))
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		tokenEntity, err := h.reshreshTokenUsecase.Get(ctx, requestBody.RefreshToken)
		if err != nil && !errors.Is(err, errorspkg.ErrorNotFound) {
			h.logger.Error("investorHandler/RefreshToken/Get", zap.Error(err))
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}
		if errors.Is(err, errorspkg.ErrorNotFound) {
			h.logger.Error("investorHandler/RefreshToken/GenerateToken", zap.Error(err))
			http.Error(w, "no such token", http.StatusBadRequest)
			return
		}

		claims, err := token.ParseJwtToken(tokenEntity.RefreshToken, h.config.Token.Secret)
		if err != nil {
			h.logger.Error("investorHandler/RefreshToken/ParseJwtToken", zap.Error(err))
			http.Error(w, "invalid authorization token", http.StatusBadRequest)
			return
		}

		sub, ok := claims["sub"]
		if !ok {
			h.logger.Error("investorHandler/RefreshToken/sub", zap.Error(err))
			http.Error(w, "not authorized", http.StatusBadRequest)
			return
		}

		role, ok := sub.(string)
		if !ok {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("failed to fetch authentication data(role)"),
				HTTPStatusCode: http.StatusUnauthorized,
				ErrorText:      "failed to fetch authentication data (role)",
			})
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			h.logger.Error("investorHandler/RefreshToken/sub", zap.Error(err))
			http.Error(w, "not authorized", http.StatusBadRequest)
			return
		}

		userIDstr, ok := userID.(string)
		if !ok {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("failed to fetch authentication data (id)"),
				HTTPStatusCode: http.StatusUnauthorized,
				ErrorText:      "failed to fetch authentication data (id)",
			})
			return
		}

		accessToken, refreshToken, err := h.reshreshTokenUsecase.GenerateToken(
			ctx,
			role,
			userIDstr,
			h.config.Token.Secret,
			h.config.Token.AccessTTL,
			h.config.Token.RefreshTTL,
		)
		if err != nil {
			h.logger.Error("investorHandler/RefreshToken/GenerateToken", zap.Error(err))
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		err = h.reshreshTokenUsecase.Delete(ctx, tokenEntity.RefreshToken)
		if err != nil {
			h.logger.Error("investorHandler/RefreshToken/Delete", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var response = models.RefreshTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		render.JSON(w, r, response)
	}
}
