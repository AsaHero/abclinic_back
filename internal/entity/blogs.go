package entity

import (
	"time"
)

const (
	PublicationTypeVideo  = "video"
	PublicationTypeSwiper = "swiper"
)

type Authors struct {
	GUID      string
	Name      string
	Img       []byte
	CreatedAt time.Time
}

type Categories struct {
	GUID        string
	Title       string
	Description string
	URL         string
	CreatedAt   time.Time
}

type Publications struct {
	GUID        string
	CategoryID  string
	AuthorID    string
	Title       string
	Description string
	Type        string
	Content     []string
	CreatedAt   time.Time
}
