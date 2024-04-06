package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type RequestTransferenciaAutomatica struct {
	CuentaId      uint64
	Cuenta        string
	SubcuentaId   uint64
	Subcuenta     string
	DatosClientes DatosClientes
	Request       RequestTransferenicaCliente
}
type RequestTransferenicaCliente struct {
	Transferencia         linktransferencia.RequestTransferenciaCreateLink `json:"transferencia,omitempty"`
	ListaMovimientosId    []uint64                                         `json:"lista_movimientos_id,omitempty"`
	ListaMovimientosIdNeg []uint64                                         `json:"lista_movimientos_id_neg,omitempty"`
}

type ResponseTransferenciaAutomatica struct {
	CuentaId uint64         `json:"cuentaid"`
	Cuenta   string         `json:"cuenta"`
	Origen   string         `json:"origen"`
	Destino  string         `json:"destino"`
	Importe  entities.Monto `json:"importe"`
	Error    string         `json:"error"`
}

type DatosClientes struct {
	NombreCliente string
	EmailCliente  string
}

// type RequestTransferenciaSubcuentas struct {
// 	CuentaId      uint64
// 	Cuenta        string
// 	Subcuenta     uint64
// 	DatosClientes DatosClientes
// 	Origen        linktransferencia.OrigenTransferenciaLink  `json:"origen"`
// 	Destino       linktransferencia.DestinoTransferenciaLink `json:"destino"`
// 	Importe       entities.Monto                             `json:"importe"`
// 	Moneda        linkdtos.EnumMoneda                        `json:"moneda"`
// 	Motivo        linkdtos.EnumMotivoTransferencia           `json:"motivo"`
// }
