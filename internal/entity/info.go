package entity

import "time"

type Chapters struct {
	GUID      string
	Title     string
	CreatedAt time.Time
}

type Articles struct {
	GUID      string
	ChapterID string
	Info      string
	Img       string	
	Side      string
	CreatedAt time.Time
}
