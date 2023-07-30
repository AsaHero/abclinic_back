package models

import "mime/multipart"

const (
	MainFolder = "main"
)

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type Path struct {
	URL string `json:"url"`
}
