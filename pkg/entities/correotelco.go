package entities

import (
	"gorm.io/gorm"
)

type TelcoCorreo struct {
	gorm.Model
	Email string `json:"email"`
}
