package v1

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

type authorsHandler struct {
	config       *config.Config
	logger       *zap.Logger
	enforcer     *casbin.Enforcer
	blogsUsecase usecase.Blogs
}

func NewAuthorsHandler(args handlers.HandlerArguments) http.Handler {
	handler := authorsHandler{
		config:       args.Config,
		logger:       args.Logger,
		blogsUsecase: args.BlogsUsecase,
		enforcer:     args.Enforcer,
	}

	policies := [][]string{
		// admin
		{"admin", "/v1/authors", "POST"},
		{"admin", "/v1/authors/{id}", "(PUT)|(DELETE)"},

		// dentist
		{"dentist", "/v1/authors", "POST"},
		{"dentist", "/v1/authors/{id}", "(PUT)|(DELETE)"},
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
		r.Get("/", handler.GetAuthorsList())
		
		router.Group(func(r chi.Router) {
			r.Use(middleware.Authorizer(handler.enforcer, handler.logger))
			
			r.Post("/", handler.CreateAuthor())
			r.Put("/{id}", handler.UpdateAuthor())
			r.Delete("/{id}", handler.DeleteAuthor())
		})
	})

	return router
}

// GetCategoriesList
// @Security ApiKeyAuth
// @Router /v1/authors [GET]
// @Summary Get authors
// @Description Get authors
// @Tags Author
// @Accept json
// @Produce json
// @Success 200 {object} []models.Authors
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h authorsHandler) GetAuthorsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authors, err := h.blogsUsecase.ListAuthors(ctx, map[string]string{})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := []models.Authors{}

		for _, v := range authors {

			imgInBase64 := base64.StdEncoding.EncodeToString(v.Img)

			response = append(response, models.Authors{
				GUID: v.GUID,
				Name: v.Name,
				Img:  imgInBase64,
			})
		}

		render.JSON(w, r, response)
	}
}

// CreateAuthor
// @Security ApiKeyAuth
// @Router /v1/authors [POST]
// @Summary Create new author
// @Description Create new author
// @Tags Author
// @Accept json
// @Produce json
// @Param body body models.CreateAuthorRequest true "body"
// @Success 200 {object} models.GUIDResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h authorsHandler) CreateAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.CreateAuthorRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		binaryImg, err := base64.StdEncoding.DecodeString(request.Img)
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("cannor decode content sent by 'img'"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      "cannor decode content sent by 'img''",
			})
			return
		}

		mimeType := http.DetectContentType(binaryImg)

		fmt.Println(mimeType)

		if mimeType != "image/jpeg" && mimeType != "image/png" {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("invalid content sent by 'img'"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      "invalid content sent by 'img'",
			})
			return
		}

		guid, err := h.blogsUsecase.CreateAuthors(ctx, &entity.Authors{
			Name: request.Name,
			Img:  binaryImg,
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

// UpdateAuthor
// @Security ApiKeyAuth
// @Router /v1/authors/{id} [PUT]
// @Summary Update author
// @Description Update author
// @Tags Author
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body body models.CreateAuthorRequest true "body"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h authorsHandler) UpdateAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		request := models.CreateAuthorRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		binaryImg, err := base64.StdEncoding.DecodeString(request.Img)
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("cannor decode content sent by 'img'"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      "cannor decode content sent by 'img''",
			})
			return
		}

		mimeType := http.DetectContentType(binaryImg)

		if mimeType != "image/jpeg" && mimeType != "image/png" {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("invalid content sent by 'img'"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      "invalid content sent by 'img'",
			})
			return
		}

		err = h.blogsUsecase.UpdateAuthors(ctx, &entity.Authors{
			GUID: guid,
			Name: request.Name,
			Img:  binaryImg,
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

// DeleteChapter
// @Security ApiKeyAuth
// @Router /v1/authors/{id} [DELETE]
// @Summary Delete author
// @Description Delete author
// @Tags Author
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h authorsHandler) DeleteAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		err := h.blogsUsecase.DeleteAuthors(ctx, guid)
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
