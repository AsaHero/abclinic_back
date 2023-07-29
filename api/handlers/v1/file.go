package v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	errorsapi "github.com/AsaHero/abclinic/api/errors"
	"github.com/AsaHero/abclinic/api/handlers"
	"github.com/AsaHero/abclinic/api/models"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type filesHandler struct {
	logger *zap.Logger
	config *config.Config
}

func NewFilesHandler(option *handlers.HandlerArguments) http.Handler {
	handler := filesHandler{
		logger: option.Logger,
		config: option.Config,
	}

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {

		// file
		r.Post("/", handler.UploadFile())
	})
	return router
}

// UploadFile
// @Router /v1/files [POST]
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} models.Path
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (handler *filesHandler) UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, header, err := r.FormFile("file")
		if err != nil {
			handler.logger.Error("error r.FormFile", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}

		file := models.File{
			File: header,
		}

		fName := uuid.New()
		file.File.Filename = fmt.Sprintf("%s_%s", fName.String(), file.File.Filename)
		dst, _ := os.Getwd()

		src, err := file.File.Open()
		if err != nil {
			handler.logger.Error("cannot open file from form-data", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}
		defer src.Close()

		out, err := os.Create(dst + "/" + file.File.Filename)
		if err != nil {
			handler.logger.Error("cannot crate new file", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}
		defer out.Close()

		_, err = io.Copy(out, src)
		if err != nil {
			handler.logger.Error("cannot copy file", zap.Error(err))
			render.Render(w, r, &errorsapi.ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorText:      err.Error(),
			})
			return
		}

		cdnConfig := &aws.Config{
			Credentials: credentials.NewStaticCredentials(
				handler.config.CDN.AwsAccessKeyID,
				handler.config.CDN.AwsSecretAccessKey,
				"",
			),
			Endpoint: aws.String(handler.config.CDN.AwsEndpoint),
			Region:   aws.String("us-east-1"),
		}

		newSession := session.New(cdnConfig)
		cdnClient := s3.New(newSession)

		// Get file size and read the file content into a buffer
		fileInfo, _ := out.Stat()
		var size int64 = fileInfo.Size()
		buffer := make([]byte, size)
		out.Read(buffer)

		object := s3.PutObjectInput{
			Bucket:             aws.String(handler.config.CDN.BucketName),
			Key:                aws.String(file.File.Filename),
			Body:               bytes.NewReader(buffer),
			ContentLength:      aws.Int64(size),
			ContentType:        aws.String(http.DetectContentType(buffer)),
			ContentDisposition: aws.String("attachment"),
			ACL:                aws.String("public-read"),
		}

		fmt.Printf("%v\n", object)
		_, err = cdnClient.PutObject(&object)
		if err != nil {
			fmt.Println(err.Error())
		}

		path := models.Path{
			Filename: fmt.Sprintf("%s/%s/%s", cdnClient.Endpoint, handler.config.CDN.BucketName, file.File.Filename),
		}

		os.Remove(dst + "/" + file.File.Filename)

		render.JSON(w, r, path)
	}
}
