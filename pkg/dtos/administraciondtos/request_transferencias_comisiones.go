package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linktransferencia"
)

type RequestTransferenciasComisiones struct {
	Transferencia           linktransferencia.RequestTransferenciaCreateLink `json:"transferencia,omitempty"`
	MovimientosIdComisiones []uint64                                         `json:"movimientos_id_comisiones"`
}

type RequestComisiones struct {
	FechaInicio      time.Time `json:"fecha_inicio"`
	FechaFin         time.Time `json:"fecha_fin"`
	MovimientosId    []uint64  `json:"movimientosId"`
	RetiroAutomatico bool      `json:"retiro_automatico"`
}

type ResponseTransferenciaComisiones struct {
	Resultado string `json:"resultado"`
}

type RequestMovimientosId struct {
	MovimientosId []uint64 `json:"movimientos_id"`
	// MovimimientosIdRevertidos []uint64 `json:"movimientos_id_revertidos"`
}

type TransferenciasComisiones struct {
	MovID uint64
	Tipo  EnumMov
}

type EnumMov string

const (
	Positivo EnumMov = "positivo"
	Negativo EnumMov = "negativo"
)

type TransferenciasClienteId struct {
	Id uint
	//Ids []uint
}
