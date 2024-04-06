package entities

import (
	"gorm.io/gorm"
)

type Movimientosubcuenta struct {
	gorm.Model
	SubcuentasID       uint64     `json:"subcuentas_id"`
	MovimientosID      uint64     `json:"movimientos_id"`
	Transferido        bool       `json:"transferido"`
	Monto              Monto      `json:"monto"`
	PorcentajeAplicado float64    `json:"porcentaje_aplicado"`
	Subcuenta          Subcuenta  `json:"subcuenta" gorm:"foreignKey:SubcuentasID"`
	Movimiento         Movimiento `json:"movimiento" gorm:"foreignKey:MovimientosID"`
}

// TableName sobreescribe el nombre de la tabla
func (Movimientosubcuenta) TableName() string {
	return "movimientosubcuentas"
}
