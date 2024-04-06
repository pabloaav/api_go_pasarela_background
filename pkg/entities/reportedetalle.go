package entities

import (
	"gorm.io/gorm"
)

type Reportedetalle struct {
	gorm.Model
	ReportesId int64 `json:"reportes_id"`
	PagosId    string
	Monto      string
	Mediopago  string
	Estadopago string
	Reporte    Reporte `json:"reportes" gorm:"foreignKey:ReportesId"`
	Pago       Pago    `json:"pago" gorm:"foreignKey:PagosId;references:ExternalReference"`
	// Pagointentos []Pagointento `json:"pagointento" gorm:"foreignKey:PagosId;references:ExternalID"`
}
