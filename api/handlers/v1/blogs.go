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

type blogsHandler struct {
	config       *config.Config
	logger       *zap.Logger
	enforcer     *casbin.Enforcer
	blogsUsecase usecase.Blogs
}

func NewBlogsHandler(args handlers.HandlerArguments) http.Handler {
	handler := blogsHandler{
		config:       args.Config,
		logger:       args.Logger,
		enforcer:     args.Enforcer,
		blogsUsecase: args.BlogsUsecase,
	}

	policies := [][]string{
		// admin
		{"admin", "/v1/blogs", "POST"},
		{"admin", "/v1/blogs/{id}", "(PUT)|(DELETE)"},
		{"admin", "/v1/blogs/{id}/publication", "POST"},
		{"admin", "/v1/blogs/publication/{id}", "(PUT)|(DELETE)"},

		// dentist
		{"dentist", "/v1/blogs", "POST"},
		{"dentist", "/v1/blogs/{id}", "(PUT)|(DELETE)"},
		{"dentist", "/v1/blogs/{id}/publication", "POST"},
		{"dentist", "/v1/blogs/publication/{id}", "(PUT)|(DELETE)"},
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
		r.Get("/", handler.GetCategoriesList())
		r.Get("/{id}/publication", handler.GetPublicationsList())

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authorizer(handler.enforcer, handler.logger))
			// category

			r.Post("/", handler.CreateCategory())
			r.Put("/{id}", handler.UpdateCategory())
			r.Delete("/{id}", handler.DeleteCategory())

			// publication
			r.Post("/{id}/publication", handler.CreatePublication())
			r.Put("/publication/{id}", handler.UpdatePublication())
			r.Delete("/publication/{id}", handler.DeletePublication())
		})

	})

	return router
}

// GetCategoriesList
// @Security ApiKeyAuth
// @Router /v1/blogs/{id}/publication [GET]
// @Summary Get publications
// @Description Get publications by category id
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} []models.Publications
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h blogsHandler) GetPublicationsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		categoryID := chi.URLParam(r, "id")

		publications, err := h.blogsUsecase.ListPublications(ctx, map[string]string{"category_id": categoryID})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := []models.Publications{}

		for _, v := range publications {
			author, err := h.blogsUsecase.GetAuthor(ctx, v.AuthorID)
			if err != nil {
				h.logger.Error("error on GetPublicationsList/ blogsUsecase.ListPublicationsAuthors", zap.Error(err))
			}

			publication := models.Publications{
				GUID:       v.GUID,
				CategoryID: v.CategoryID,
				Author: models.Authors{
					GUID: author.GUID,
					Name: author.Name,
					// Photo: author.Img,
				},
				Title: v.Title,
				Text:  v.Description,
			}

			if v.Type == entity.PublicationTypeVideo {
				publication.Video = v.Content[0]
			}

			if v.Type == entity.PublicationTypeSwiper {
				for _, imgs := range v.Content {
					publication.Img = append(publication.Img, models.Contents{
						URL: imgs,
					})
				}
			}

			response = append(response, publication)
		}

		render.JSON(w, r, response)
	}
}

// CreatePublication
// @Security ApiKeyAuth
// @Router /v1/blogs/{id}/publication [POST]
// @Summary Create new publication
// @Description Create new publication
// @Tags Blogs
// @Accept json
// @Produce json
// @Param body body models.CreatePublicationRequest true "body"
// @Success 200 {object} models.GUIDResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h blogsHandler) CreatePublication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		categoryId := chi.URLParam(r, "id")

		request := models.CreatePublicationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Error("error on decoding request body", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		publication := &entity.Publications{
			CategoryID:  categoryId,
			AuthorID:    request.AuthorID,
			Title:       request.Title,
			Description: request.Text,
			Type:        request.Type,
		}

		if publication.Type == entity.PublicationTypeSwiper {
			for _, v := range request.Img {
				publication.Content = append(publication.Content, v.URL)
			}
		} else if publication.Type == entity.PublicationTypeVideo {
			publication.Content = append(publication.Content, request.Video)
		} else {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("error: invalid type of content"),
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      "invalid type of content",
			})
			return
		}

		guid, err := h.blogsUsecase.CreatePublications(ctx, publication)
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

// UpdatePublication
// @Security ApiKeyAuth
// @Router /v1/blogs/publication/{id} [PUT]
// @Summary Update publication
// @Description Update publication
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body body models.UpdatePublicationRequest true "body"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h blogsHandler) UpdatePublication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		request := models.UpdatePublicationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		publication := &entity.Publications{
			GUID:        guid,
			Title:       request.Title,
			Description: request.Text,
			Type:        request.Type,
		}

		if publication.Type == entity.PublicationTypeSwiper {
			for _, v := range request.Img {
				publication.Content = append(publication.Content, v.URL)
			}
		} else if publication.Type == entity.PublicationTypeVideo {
			publication.Content = append(publication.Content, request.Video)
		} else {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            errors.New("error: invalid type of content"),
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      "invalid type of content",
			})
			return
		}

		err := h.blogsUsecase.UpdatePublications(ctx, publication)
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

// DeletePublication
// @Security ApiKeyAuth
// @Router /v1/blogs/publication/{id} [DELETE]
// @Summary Delete publication
// @Description Delete publication by guid
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h blogsHandler) DeletePublication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		err := h.blogsUsecase.DeletePublications(ctx, guid)
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

// GetCategoriesList
// @Security ApiKeyAuth
// @Router /v1/blogs [GET]
// @Summary Get publication categories
// @Description Get publication categories
// @Tags Blogs
// @Accept json
// @Produce json
// @Success 200 {object} []models.Categories
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h blogsHandler) GetCategoriesList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		categories, err := h.blogsUsecase.ListPublicationsCategories(ctx, map[string]string{})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := []models.Categories{}

		for _, v := range categories {
			response = append(response, models.Categories{
				GUID:        v.GUID,
				Title:       v.Title,
				Description: v.Description,
				Img:         v.URL,
			})
		}

		render.JSON(w, r, response)
	}
}

// CreateCategory
// @Security ApiKeyAuth
// @Router /v1/blogs [POST]
// @Summary Create new publication category
// @Description Create new publication category
// @Tags Blogs
// @Accept json
// @Produce json
// @Param body body models.CreateCategoryRequest true "body"
// @Success 200 {object} models.GUIDResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h blogsHandler) CreateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.CreateCategoryRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		guid, err := h.blogsUsecase.CreatePublicationsCategories(ctx, &entity.Categories{
			Title:       request.Title,
			Description: request.Description,
			URL:         request.Img,
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

// UpdateCategory
// @Security ApiKeyAuth
// @Router /v1/blogs/{id} [PUT]
// @Summary Update publication category
// @Description Update publication category
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body body models.CreateCategoryRequest true "body"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h blogsHandler) UpdateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		request := models.CreateCategoryRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		err := h.blogsUsecase.UpdatePublicationsCategories(ctx, &entity.Categories{
			GUID:        guid,
			Title:       request.Title,
			Description: request.Description,
			URL:         request.Img,
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
// @Router /v1/blogs/{id} [DELETE]
// @Summary Delete publication category
// @Description Delete publication category
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h blogsHandler) DeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		err := h.blogsUsecase.DeletePublicationsCategories(ctx, guid)
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
