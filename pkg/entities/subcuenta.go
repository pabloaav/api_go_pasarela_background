package entities

import (
	"gorm.io/gorm"
)

type Subcuenta struct {
	gorm.Model
	Nombre               string                `json:"nombre"`
	Email                string                `json:"email"`
	Tipo                 string                `json:"tipo"`
	CuentasID            uint                  `json:"cuentas_id"`
	Cbu                  string                `json:"cbu"`
	Porcentaje           float64               `json:"porcentaje"`
	Cuenta               Cuenta                `json:"cuenta" gorm:"foreignKey:CuentasID"`
	Movimientosubcuentas []Movimientosubcuenta `json:"movimientos_subcuentas" gorm:"foreignKey:SubcuentasID"`
}

// TableName sobreescribe el nombre de la tabla
func (Subcuenta) TableName() string {
	return "subcuentas"
}
