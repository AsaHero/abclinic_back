package v1

import (
	"encoding/json"
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

type infoHandler struct {
	config      *config.Config
	logger      *zap.Logger
	enforcer    *casbin.Enforcer
	infoUsecase usecase.InfoUsecase
}

func NewInfoHandler(args handlers.HandlerArguments) http.Handler {
	handler := infoHandler{
		config:      args.Config,
		logger:      args.Logger,
		enforcer:    args.Enforcer,
		infoUsecase: args.InfoUsecase,
	}

	policies := [][]string{
		// admin
		{"admin", "/v1/articles", "POST"},
		{"admin", "/v1/articles/{id}", "PUT"},
		{"admin", "/v1/articles/{id}", "DELETE"},
		{"admin", "/v1/articles/chapter", "POST"},
		{"admin", "/v1/articles/chapter/{id}", "PUT"},
		{"admin", "/v1/articles/chapter/{id}", "DELETE"},

		// secretary
		{"secretary", "/v1/articles", "POST"},
		{"secretary", "/v1/articles/{id}", "PUT"},
		{"secretary", "/v1/articles/{id}", "DELETE"},
		{"secretary", "/v1/articles/chapter", "POST"},
		{"secretary", "/v1/articles/chapter/{id}", "PUT"},
		{"secretary", "/v1/articles/chapter/{id}", "DELETE"},
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
		r.Get("/{id}", handler.GetArticlesByChapter())
		r.Get("/chapter", handler.GetChapterList())
		r.Get("/chapter/{id}", handler.GetChpater())

	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.Authorizer(handler.enforcer, handler.logger))

		r.Post("/", handler.CreateArticle())
		r.Put("/{id}", handler.UpdateArticle())
		r.Delete("/{id}", handler.DeleteArticle())
		r.Post("/chapter", handler.CreateChapter())
		r.Put("/chapter/{id}", handler.UpdateChapter())
		r.Delete("/chapter/{id}", handler.DeleteChapter())
	})
	return router
}

// GetArticlesByChapter
// @Security ApiKeyAuth
// @Router /v1/articles/{id} [GET]
// @Summary Get articles
// @Description Get articles by chapter id
// @Tags Info
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} []models.Article
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) GetArticlesByChapter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		chapterID := chi.URLParam(r, "id")

		articles, err := h.infoUsecase.ListArticles(ctx, map[string]string{"chapter_id": chapterID})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := []models.Article{}

		for _, v := range articles {
			response = append(response, models.Article{
				GUID: v.GUID,
				Text: v.Info,
				Img:  v.Img,
				Side: v.Side,
			})
		}

		render.JSON(w, r, response)
	}
}

// CreateArticle
// @Security ApiKeyAuth
// @Router /v1/articles [POST]
// @Summary Create new article
// @Description Create new article
// @Tags Info
// @Accept json
// @Produce json
// @Param body body models.CreateArticleRequest true "body"
// @Success 200 {object} models.GUIDResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) CreateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.CreateArticleRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Error("error on decoding request body", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		guid, err := h.infoUsecase.CreateArticle(ctx, &entity.Articles{
			ChapterID: request.ChapterID,
			Info:      request.Text,
			Img:       request.Img,
			Side:      request.Side,
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

// UpdateArticle
// @Security ApiKeyAuth
// @Router /v1/articles/{id} [PUT]
// @Summary Update article
// @Description Update article
// @Tags Info
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body body models.UpdateArticleRequest true "body"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) UpdateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		request := models.UpdateArticleRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		err := h.infoUsecase.UpdateArticles(ctx, &entity.Articles{
			GUID: guid,
			Info: request.Text,
			Img:  request.Img,
			Side: request.Side,
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

// DeleteArticle
// @Security ApiKeyAuth
// @Router /v1/articles/{id} [DELETE]
// @Summary Delete article
// @Description Delete article by guid
// @Tags Info
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) DeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		err := h.infoUsecase.DeleteArticles(ctx, guid)
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

// GetChapterList
// @Security ApiKeyAuth
// @Router /v1/articles/chapter [GET]
// @Summary Get article chapters
// @Description Get article chapters
// @Tags Info
// @Accept json
// @Produce json
// @Success 200 {object} []models.Chapter
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) GetChapterList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		chapters, err := h.infoUsecase.ListArticlesChapters(ctx, map[string]string{})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := []models.Chapter{}

		for _, v := range chapters {
			response = append(response, models.Chapter{
				GUID: v.GUID,
				Name: v.Title,
			})
		}

		render.JSON(w, r, response)
	}
}

// GetChpater
// @Security ApiKeyAuth
// @Router /v1/articles/chapter/{id} [GET]
// @Summary Get article chapter
// @Description Get article chapter
// @Tags Info
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.ServicesGroup
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) GetChpater() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		chapters, err := h.infoUsecase.ListArticlesChapters(ctx, map[string]string{"guid": guid})
		if err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		response := models.Chapter{
			GUID: chapters[0].GUID,
			Name: chapters[0].Title,
		}

		render.JSON(w, r, response)
	}
}

// CreateChapter
// @Security ApiKeyAuth
// @Router /v1/articles/chapter [POST]
// @Summary Create new article chapter
// @Description Create new article chapter
// @Tags Info
// @Accept json
// @Produce json
// @Param body body models.CreateChapterRequest true "body"
// @Success 200 {object} models.GUIDResponse
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) CreateChapter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := models.CreateChapterRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		guid, err := h.infoUsecase.CreateArticlesChapter(ctx, &entity.Chapters{
			Title: request.Name,
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

// UpdateChapter
// @Security ApiKeyAuth
// @Router /v1/articles/chapter/{id} [PUT]
// @Summary Update article chapter
// @Description Update articles chapter
// @Tags Info
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body body models.CreateChapterRequest true "body"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) UpdateChapter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		request := models.CreateChapterRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		err := h.infoUsecase.UpdateArticlesChapter(ctx, &entity.Chapters{
			GUID:  guid,
			Title: request.Name,
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
// @Router /v1/articles/chapter/{id} [DELETE]
// @Summary Delete article cahpter
// @Description Delete article chapter
// @Tags Info
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Empty
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h infoHandler) DeleteChapter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		guid := chi.URLParam(r, "id")

		err := h.infoUsecase.DeleteArticlesChapter(ctx, guid)
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
