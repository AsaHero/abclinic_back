package models

import "mime/multipart"

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type Path struct {
	Filename string `json:"filename"`
}
