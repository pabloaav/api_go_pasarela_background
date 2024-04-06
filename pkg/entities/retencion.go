package entities

import (
	"time"

	"gorm.io/gorm"
)

type Retencion struct {
	gorm.Model
	CondicionsId      uint               `json:"condicions_id"`
	Condicion         Condicion          `gorm:"foreignKey:CondicionsId"`
	ChannelsId        uint               `json:"channels_id"`
	Channel           Channel            `gorm:"foreignKey:ChannelsId"`
	Alicuota          float64            `json:"alicuota"`
	AlicuotaOpcional  float64            `json:"alicuota_opcional"`
	Rg2854            bool               `json:"rg2854"`
	Minorista         bool               `json:"minorista"`
	MontoMinimo       float64            `json:"monto_minimo"`
	Descripcion       string             `json:"descripcion"`
	CodigoRegimen     string             `json:"codigo"`
	Clientes          []Cliente          `gorm:"many2many:cliente_retencions;"`
	ClienteRetencions []ClienteRetencion `json:"cliente_retencions" gorm:"foreignKey:retencion_id"`
	FechaValidezDesde time.Time             `json:"fecha_validez_desde"`
	FechaValidezHasta time.Time             `json:"fecha_validez_hasta"`
}

// Función para verificar si una fecha está en un rango de fechas.
func (r Retencion) IsCurrent(fechaPaidAt time.Time) bool {
	return fechaPaidAt.After(r.FechaValidezDesde) && (r.FechaValidezHasta.IsZero() || fechaPaidAt.Before(r.FechaValidezHasta))
}
