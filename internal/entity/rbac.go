package entity

import "time"

const (
	RoleAdmin     = "admin"
	RoleWebsite   = "website"
	RoleDentist   = "dentist"
	RoleSecretary = "secretary"
)

type Users struct {
	GUID      string
	Role   string
	Firstname string
	Lastname  string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

