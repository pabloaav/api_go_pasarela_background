package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"

type SaldoCuentaResponse struct {
	CuentasId uint64         `json:"cuentas_id"`
	Total     entities.Monto `json:"total"`
}
