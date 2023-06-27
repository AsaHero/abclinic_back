package usecase

import (
	"time"

	"github.com/google/uuid"
)

type BaseUsecase struct {
}

func (u BaseUsecase) beforeCreate(guid *string, created_at *time.Time, update_at *time.Time) {
	if guid != nil {
		*guid = uuid.New().String()
	}

	if created_at != nil {
		*created_at = time.Now().Local()
	}

	if update_at != nil {
		*update_at = time.Now().Local()
	}

}
