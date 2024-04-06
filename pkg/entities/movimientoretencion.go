package entities

import "gorm.io/gorm"

type MovimientoRetencion struct {
	gorm.Model
	MovimientoId    uint64    `json:"movimiento_id"`
	RetencionId     uint      `json:"retencion_id"`
	Retencion       Retencion `gorm:"foreignKey:RetencionId"`
	ClienteId       uint64    `json:"cliente_id"`
	Monto           Monto     `json:"monto"`
	ImporteRetenido Monto     `json:"importe_retenido"`
	Efectuada       bool      `json:"efectuada"`
}

type MovimientosRetenciones []MovimientoRetencion 

// obtener de un []MovimientoRetencion, un mov retencion por el id del gravamen 
func (mrs MovimientosRetenciones) GetByGravamenId(id uint) (mov_ret MovimientoRetencion, result bool){
	if len(mrs) == 0 {
		return
	}
	for _, m := range mrs {
		if m.Retencion.Condicion.Gravamen.ID == id {
			mov_ret = m
			result = true
			break
		}
	}
	return
}