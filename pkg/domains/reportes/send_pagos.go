package reportes

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/reportedtos"
)

type sendPagos struct {
	utilService util.UtilService
}

func SendPagos(util util.UtilService) Email {
	return &sendPagos{
		utilService: util,
	}
}

func (cl *sendPagos) SendReportes(ruta string, nombreArchivo string, request reportedtos.ResponseClientesReportes) (erro error) {

	if len(request.Pagos) > 0 {
		RutaFile := fmt.Sprintf("%s/%s.csv", ruta, nombreArchivo)
		/* estos datos son los que se van a escribir en el archivo */

		var slice_array = [][]string{{"TRANSACCIONES COBRADAS"},
			{"FECHA", request.Fecha, "", "", "", "", "", "", ""}}

		if request.SujetoRetencion {
			// estructura del excel cuando se trata de un cliente con retencion
			slice_array = append(slice_array, []string{"CUENTA", "REFERENCIA", "FECHA COBRO", "MEDIO DE PAGO", "TIPO", "ESTADO", "MONTO", "COMISION $", "IVA %", "RETENCION $"}) // columnas

			for _, pago := range request.Pagos {
				slice_array = append(slice_array, []string{pago.Cuenta, pago.Id, pago.FechaPago, pago.MedioPago, pago.Tipo, pago.Estado, pago.Monto, pago.Comision, pago.Iva, pago.Retencion})
			}

			/* Crear archivo csv  en la carpeta documentos/reportes*/
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "", "CANT OPERACIONES", request.CantOperaciones})
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "", "TOTAL COBRADO $", request.TotalCobrado})
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "", "TOTAL COMISION $", request.TotalComision})
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "", "TOTAL IVA $", request.TotalIva})
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "", "TOTAL RETENCION $", request.TotalRetencion})

		} else {
			// estructura del excel cuando se trata de un cliente sin retencion
			slice_array = append(slice_array, []string{"CUENTA", "REFERENCIA", "FECHA COBRO", "MEDIO DE PAGO", "TIPO", "ESTADO", "MONTO", "COMISION $", "IVA %"}) // columnas

			for _, pago := range request.Pagos {
				slice_array = append(slice_array, []string{pago.Cuenta, pago.Id, pago.FechaPago, pago.MedioPago, pago.Tipo, pago.Estado, pago.Monto, pago.Comision, pago.Iva})
			}

			/* Crear archivo csv  en la carpeta documentos/reportes*/
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "CANT OPERACIONES", request.CantOperaciones})
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "TOTAL COBRADO $", request.TotalCobrado})
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "TOTAL COMISION $", request.TotalComision})
			slice_array = append(slice_array, []string{"", "", "", "", "", "", "", "TOTAL IVA $", request.TotalIva})
		}

		erro = cl.utilService.CsvCreate(RutaFile, slice_array)
		if erro != nil {
			return
		}

	}
	return
}
