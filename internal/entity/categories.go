package entity

import "time"

type Categories struct {
	GUID        string
	Title       string
	Description string
	Path        string
	URL         string
	Created_at  time.Time
}
