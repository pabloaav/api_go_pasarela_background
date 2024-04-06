package entities

import (
	"gorm.io/gorm"
)

type Multipagoscierrelotedetalles struct {
	gorm.Model
	MultipagoscierrelotesId int64 `json:"multipagoscierrelotes_id"`
	FechaCobro              string
	ImporteCobrado          int64
	CodigoBarras            string
	Clearing                string
	ImporteCalculado        float64
	Match                   bool
	Enobservacion           bool
	Pagoinformado           bool
	MultipagosCabecera      Multipagoscierrelote `json:"multipagoscierrelotes" gorm:"foreignKey:MultipagoscierrelotesId"`
}
