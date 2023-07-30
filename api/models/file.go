package models

import "mime/multipart"

const (
	MainFolder = "main"
)

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type Path struct {
	Filename string `json:"filename"`
}
