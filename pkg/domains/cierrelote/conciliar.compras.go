package cierrelote

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type conciliarCompras struct {
	utilService util.UtilService
}

func NewConciliarCompras(util util.UtilService) MetodoConciliarClMP {
	return &conciliarCompras{
		utilService: util,
	}
}

func (c *conciliarCompras) ConciliarTablas(valorCuota int64, cierreLote prismaCierreLote.ResponsePrismaCL, movimientoCabecera prismaCierreLote.ResponseMovimientoTotales, movimientoDetalle prismaCierreLote.ResponseMoviminetoDetalles) (listaCierreLoteProcesada []prismaCierreLote.ResponsePrismaCL, detalleMoviminetosIdArray []int64, cabeceraMoviminetosIdArray []int64, erro error) {
	if cierreLote.Tipooperacion == "C" {
		strNroEstablecimiento := strconv.Itoa(int(cierreLote.Nroestablecimiento))
		// logs.Info(strings.TrimSpace(movimientoDetalle.NroAutorizacionXl))
		// logs.Info(cierreLote.Codigoautorizacion)
		if cierreLote.Fechaoperacion == movimientoDetalle.FechaOrigenCompra && strings.Contains(movimientoCabecera.EstablecimientoNro, strNroEstablecimiento) && cierreLote.Nrotarjeta == movimientoDetalle.NroTarjetaXl && strings.Contains(strings.TrimSpace(movimientoDetalle.NroAutorizacionXl), cierreLote.Codigoautorizacion) && cierreLote.Nroticket == movimientoDetalle.NroCupon && valorCuota == movimientoDetalle.PlanCuota && movimientoDetalle.TipoAplicacion == "+" && movimientoCabecera.Codop == movimientoDetalle.Tipooperacion.ExternalId { //&& cierreLote.Monto.Int64() == int64(movimientoDetalle.Importe)

			if cierreLote.Monto.Int64() < int64(movimientoDetalle.Importe) || cierreLote.Monto.Int64() > int64(movimientoDetalle.Importe) {
				// NOTE corresponde a diferencias recibidas en archivos mxdetalles y prismacierrelote
				// esto puede variar por redondeos de importes de tablas
				mensaje := fmt.Sprintf("error: existe diferencia entre monto cierre lote $%v e importe movimiento detalles $%v. cl_id = %v y movimientoDetalle_id = %v ", cierreLote.Monto.Int64(), int64(movimientoDetalle.Importe), cierreLote.Id, movimientoDetalle.Id)
				logs.Error(mensaje)
				diferencia := cierreLote.Monto.Int64() - movimientoDetalle.Importe.Int64()
				valorAbsoluto := int64(math.Abs(float64(diferencia)))
				if valorAbsoluto <= 3 {
					cierreLote.Monto = movimientoDetalle.Importe
					cierreLote.MontoModificado = true
				} else {
					erro = errors.New(mensaje)
					return
				}
			}
			/////////////////////////////////////////////////////////////////////////////////////
			// if cierreLote.Monto.Int64() < int64(movimientoDetalle.Importe) || cierreLote.Monto.Int64() > int64(movimientoDetalle.Importe) {
			// 	mensaje := fmt.Sprintf("error: existe diferencia entre monto cierre lote $%v e importe movimiento detalles $%v. cl_id = %v y movimientoDetalle_id = %v ", cierreLote.Monto.Int64(), int64(movimientoDetalle.Importe), cierreLote.Id, movimientoDetalle.Id)
			// 	logs.Error(mensaje)
			// 	// diferencia := cierreLote.Monto.Int64() - movimientoDetalle.Importe.Int64()
			// 	//if diferencia > 3 || diferencia < -3 {
			// 	erro = errors.New(mensaje)
			// 	return
			// 	//}
			// 	//cierreLote.Monto = movimientoDetalle.Importe
			// }
			//////////////////////////////////////////////////////////////////////////////////
			// NOTE se controla que la comision calculada por prisma sea igual a la calculda por la pasarela
			porcentajeArancelControl := c.utilService.ToFixed(cierreLote.Channelarancel.Importe*100, 2) // comision wee
			porcentajeArancelPrisma := movimientoDetalle.PorcentDescArancel / 100                       // comision prisma
			// si los porcentajes no coinciden se marca como observacion
			if porcentajeArancelControl != porcentajeArancelPrisma {
				logs.Error(fmt.Sprintf("error: existe diferencia entre porcentaje arancel control %v y porcentaje arancel prisma %v. cl_id = %v y movimientoDetalle_id = %v", porcentajeArancelControl, porcentajeArancelPrisma, cierreLote.Id, movimientoDetalle.Id))
				cierreLote.Enobservacion = true
			}
			cierreLote.FechaPago = movimientoCabecera.FechaPago
			fecha := cierreLote.FechaCierre
			if !cierreLote.FechaCierre.Equal(movimientoCabecera.FechaPresentacion) {
				fecha = movimientoCabecera.FechaPresentacion
				cierreLote.Descripcionpresentacion = fmt.Sprintf("la fecha de cierre en CL %v se modifico por %v", cierreLote.FechaCierre, movimientoCabecera.FechaPresentacion)
				cierreLote.FechaCierre = fecha
			}
			cierreLote.Cantdias = int(cierreLote.FechaPago.Sub(fecha).Hours() / 24)
			cierreLote.PrismamovimientodetallesId = movimientoDetalle.Id
			detalleMoviminetosIdArray = append(detalleMoviminetosIdArray, movimientoDetalle.Id)
			cabeceraMoviminetosIdArray = append(cabeceraMoviminetosIdArray, int64(movimientoDetalle.PrismamovimientototalesId))

			RequestValidarCF := utildtos.RequestValidarCF{
				Cupon:  cierreLote.Monto,
				Cuotas: float64(cierreLote.Nrocuota),
				Dias:   float64(cierreLote.Cantdias),
				Tna:    cierreLote.Istallmentsinfo.Tna,
				//channelArancel.importe en prisma es el porcentaje
				ArancelMonto: cierreLote.Channelarancel.Importe,
			}
			// el valor presentado dentro del objeto representa el importe sin el costo finananciero y importe de arancel
			responseValidarCF := util.Resolve().ValidarCalculoCF(RequestValidarCF)
			valor_pres := entities.Monto(responseValidarCF.ValorPresente * 100)
			cierreLote.Valorpresentado = valor_pres
			// importe del arancel se obtiene de multiplicar el coeficiente del arancel por el monto
			cierreLote.ImportearancelCalculado = util.Resolve().ToFixed(cierreLote.Monto.Float64()*cierreLote.Channelarancel.Importe, 2) // util.Resolve().ToFixed(cierreLote.Monto.Float64()-responseValidarCF.ValorPresente, 4)
			// el vivs del arancel se calcula de multiplicar el importe arancel calculado por el 0.21
			cierreLote.Importeivaarancel = util.Resolve().ToFixed(cierreLote.ImportearancelCalculado*0.21, 2)
			// el importe del cf se obtiene de lo informado por prisma
			if movimientoDetalle.PlanCuota > 0 {
				cierreLote.ImporteCfPrisma = movimientoDetalle.ImporteCostoFinanciero.Float64()
				if strings.Contains(movimientoDetalle.IdCf, "10,5%") {
					cierreLote.ImporteIvaCfCalculado = util.Resolve().ToFixed(cierreLote.ImporteCfPrisma*0.105, 2)
				}
				if strings.Contains(movimientoDetalle.IdCf, "21%") {
					cierreLote.ImporteIvaCfCalculado = util.Resolve().ToFixed(cierreLote.ImporteCfPrisma*0.21, 2)
				}
			}

			// logs.Info("========")
			// logs.Info(cierreLote.Monto.Float64())
			// logs.Info(responseValidarCF.ValorPresente)
			// logs.Info(cierreLote.Diferenciaimporte)
			// logs.Info(cierreLote.Importeivaarancel)
			// logs.Info("========")
			cierreLote.Coeficientecalculado = responseValidarCF.ValorCoeficiente
			cierreLote.Costototalporcentaje = responseValidarCF.CostoTotalPorcentaje

			listaCierreLoteProcesada = append(listaCierreLoteProcesada, cierreLote)

		}
	}
	return
}

/*
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("============Fecha Operacion======")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.Fechaoperacion, movimientoDetalle.FechaOrigenCompra)
	fmt.Println("============Fecha Presentacion===")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.FechaCierre, movimientoCabecera.FechaPresentacion)
	fmt.Println("============Establecimiento======")
	fmt.Printf("cl: %v - movimiento: %v \n", movimientoCabecera.EstablecimientoNro, strNroEstablecimiento)

	fmt.Println("============Nro.Tarjeta==========")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.Nrotarjeta, movimientoDetalle.NroTarjetaXl)
	fmt.Println("============Nor.Atorizacion======")
	fmt.Printf("cl: %v - movimiento: %v \n", movimientoDetalle.NroAutorizacionXl, cierreLote.Codigoautorizacion)

	fmt.Println("============Lote=================")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.ExternalloteId, movimientoDetalle.Lote)

	fmt.Println("============Ticket===============")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.Nroticket, movimientoDetalle.NroCupon)

	fmt.Println("============Cuota================")
	fmt.Printf("cl: %v - movimiento: %v \n", valorCuota, movimientoDetalle.PlanCuota)

	fmt.Println("============Importe==============")
	fmt.Printf("cl: %v - movimiento: %v \n", cierreLote.Monto.Int64(), int64(movimientoDetalle.Importe))
	fmt.Println("=================================")
	fmt.Println("=================================")
*/
