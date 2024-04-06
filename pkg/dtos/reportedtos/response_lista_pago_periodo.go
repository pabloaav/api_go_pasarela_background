package reportedtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos"

type ResponseListaPagoPeriodo struct {
	PagosByPeriodo       []ResultadoPagosPeriodo
	TotalImporteRendidio float64
	Meta                 dtos.Meta
}
