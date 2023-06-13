package v1

import (
	"net/http"
	"strconv"

	errorsapi "github.com/AsaHero/abclinic/api/errors"
	"github.com/AsaHero/abclinic/api/handlers"
	"github.com/AsaHero/abclinic/api/models"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/AsaHero/abclinic/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type dentistsHandler struct {
	config          *config.Config
	logger          *zap.Logger
	dentistsUsecase usecase.Denstists
}

func NewDentistsHandler(args handlers.HandlerArguments) http.Handler {
	handler := dentistsHandler{
		config:          args.Config,
		logger:          args.Logger,
		dentistsUsecase: args.DentistsUsecase,
	}

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Get("/{id}", handler.GetDentistsList())
	})

	return router
}

// GetDentistsList
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
func (h dentistsHandler) GetDentistsList() http.HandlerFunc {
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
			Name:      dentist.Name,
			Info:      dentist.Info,
			Img:       dentist.URL,
			Side:      dentist.Side,
			Priority:  dentist.Priority,
			Language:  dentist.Language,
		}

		render.JSON(w, r, response)
	}
}
