package utilfake

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"

type TableDriverBuildComisiones struct {
	TituloPrueba          string
	WantTable             error
	RequestMovimiento     *entities.Movimiento
	RequestCuentaComision *[]entities.Cuentacomision
	RequestIva            *entities.Impuesto
	ImporteSolicitado     entities.Monto
}

// const ERROR_CALCULO_COMISION = "error de validación: no se pudo obtener calculo de comisiones"
type TableDriverBuildMovimientosSubcuentas struct {
	TituloPrueba      string
	WantTable         error
	RequestMovimiento *entities.Movimiento
	RequestCuenta     entities.Cuenta
}
