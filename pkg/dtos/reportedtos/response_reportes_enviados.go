package reportedtos

import (
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type ResponseReportesEnviados struct {
	Reportes []ResponseReporteEnviado `json:"reportes"`
	Meta     dtos.Meta                `json:"meta"`
}

type ResponseReporteEnviado struct {
	Id             uint
	Cliente        string
	Tiporeporte    string
	Totalcobrado   string
	Totalrendido   string
	Fechacobranza  string
	Fecharendicion string
	Fechacreacion  string
	Nro_reporte    string
	Detalles       []ResponseReporteDetalle
}

type ResponseReporteDetalle struct {
	Id         uint
	PagosId    string
	Monto      string
	Mediopago  string
	Estadopago string
	Comision   float64
	Iva        float64
	Retencion  float64
}

func (rre *ResponseReporteEnviado) EntityToDto(entity entities.Reporte) {
	DDMMYYYY := "02-01-2006"
	rre.Id = entity.ID
	rre.Cliente = entity.Cliente
	rre.Tiporeporte = entity.Tiporeporte
	rre.Totalcobrado = entity.Totalcobrado
	rre.Totalrendido = entity.Totalrendido
	rre.Fechacobranza = entity.Fechacobranza
	rre.Fecharendicion = entity.Fecharendicion
	rre.Fechacreacion = entity.CreatedAt.Format(DDMMYYYY)
	rre.Nro_reporte = strconv.FormatUint(uint64(entity.Nro_reporte), 10)
	rre.Detalles = []ResponseReporteDetalle{}

	if len(entity.Reportedetalle) > 0 {
		respDetalleTemp := ResponseReporteDetalle{}
		for _, detalle := range entity.Reportedetalle {
			respDetalleTemp.ReporteDetalleToDto(detalle)
			rre.Detalles = append(rre.Detalles, respDetalleTemp)
		}
	}

}

func (rrd *ResponseReporteDetalle) ReporteDetalleToDto(entity entities.Reportedetalle) {
	rrd.Id = entity.ID
	rrd.PagosId = entity.PagosId
	rrd.Monto = entity.Monto
	rrd.Mediopago = entity.Mediopago
	rrd.Estadopago = entity.Estadopago

	if len(entity.Pago.PagoIntentos) > 0 {
		pagoIntento := entity.Pago.PagoIntentos[len(entity.Pago.PagoIntentos)-1]

		if len(pagoIntento.Movimientotemporale) > 0 {
			movTemp := pagoIntento.Movimientotemporale[0]

			if len(movTemp.Movimientocomisions) > 0 {
				movComision := movTemp.Movimientocomisions[0]
				rrd.Comision = (movComision.Monto + movComision.Montoproveedor).Float64()
			}

			if len(movTemp.Movimientoimpuestos) > 0 {
				movIva := movTemp.Movimientoimpuestos[0]
				rrd.Iva = (movIva.Monto + movIva.Montoproveedor).Float64()
			}

			if len(movTemp.Movimientoretenciontemporales) > 0 {
				movRetencion := movTemp.Movimientoretenciontemporales[0]
				rrd.Retencion = movRetencion.ImporteRetenido.Float64()
			}
		}
	}
}
