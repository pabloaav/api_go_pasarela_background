package reportes

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos/retenciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/commonsdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	filtros_reportes "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/reportes"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/email"
)

type ReportesService interface {
	/* 1 OBTENER CLIENTES: esto nos permite filtrar los pagos de cada clientes*/
	GetClientes(request reportedtos.RequestPagosClientes) (response administraciondtos.ResponseFacturacionPaginado, erro error)

	/* REPORTES GENERALES : reportes de todos los pagos (se generan desde el frontend) */
	GetPagosReportes(request reportedtos.RequestPagosPeriodo) (response []reportedtos.ResponsePagosPeriodo, erro error)
	ResultPagosReportes(request []reportedtos.ResponsePagosPeriodo, paginacion filtros.Paginacion) (response reportedtos.ResponseListaPagoPeriodo, erro error)

	/* REPORTES GENERADOS: estos son los reportes que se enviaran a cada cliente via email */
	GetPagosClientes(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error)               /* Todos los pagos*/
	GetRendicionClientes(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error)           /* Pagos acreditados */
	GetReversionesClientes(requestCliente administraciondtos.ResponseFacturacionPaginado, request reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error) /* Pagos revertidos */

	GetPagosClientesMensual(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error)
	GetRendicionesClientesMensual(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error)

	/* ENVIAR REPORTES : se envia por correo electronico los pagos a cada clientes*/
	SendPagosClientes(request []reportedtos.ResponseClientesReportes) (errorFile []reportedtos.ResponseCsvEmailError, erro error)
	SendLiquidacionClientes(request []reportedtos.ResultMovLiquidacion) (errorFile []reportedtos.ResponseCsvEmailError, erro error)
	//Envia el reporde de reversiones a los clientes.
	SendReporteReversiones(request []reportedtos.ResponseClientesReportes) (errorFile []reportedtos.ResponseCsvEmailError, erro error)

	SendReporteRendiciones(request []reportedtos.ResponseClientesReportes) (errorFile []reportedtos.ResponseCsvEmailError, numero_reporte_rrm uint, erro error)

	// Servicio para enviar archivos txt de retenciones
	SendRetencionestxt(request reportedtos.RequestRetencionEmail) (erro error)

	MakeControlReportes(apikeys string, fechaControlar string, token string, runEndpoint util.RunEndpoint) (response []reportedtos.ResponseControlReporte, erro error)
	SendControlReportes(request []reportedtos.ResponseControlReporte) (erro error)

	/* REPORTES DE COBRANZAS PARA CLIENTES(batch): se genera un archivo txt */
	GetPagoItems(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosItems, erro error)
	GetPagoItemsAlternativo(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosItemsAlternativo, erro error)
	BuildPagosItems(request []reportedtos.ResponsePagosItems) (response []reportedtos.ResultPagosItems)
	ValidarEsctucturaPagosItems(request []reportedtos.ResultPagosItems) error
	ValidarEsctucturaPagosBatch(request []reportedtos.ResultPagosItemsAlternativo) error
	SendPagosItems(ctx context.Context, request []reportedtos.ResultPagosItems, filtro reportedtos.RequestPagosClientes) error // tambien debe retornar la lista de pagos para insertar en la tabla pagoslotes(pagos que ya fueron enviados)

	SendPagosBatch(ctx context.Context, request []reportedtos.ResultPagosItemsAlternativo, filtro reportedtos.RequestPagosClientes) error // tambien debe retornar la lista de pagos para insertar en la tabla pagoslotes(pagos que ya fueron enviados)

	// Build pagos archivo alternativo
	BuildPagosArchivo(request []reportedtos.ResponsePagosItemsAlternativo) (response []reportedtos.ResultPagosItemsAlternativo)

	// generar comprobante de liquidacion (DPEC)
	GetRecaudacion(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosLiquidacion, erro error)
	BuildPagosLiquidacion(request []reportedtos.ResponsePagosLiquidacion) (response []reportedtos.ResultPagosLiquidacion)

	// recaudacion diaria
	GetRecaudacionDiaria(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseMovLiquidacion, erro error)
	BuildMovLiquidacion(request []reportedtos.ResponseMovLiquidacion) (response []reportedtos.ResultMovLiquidacion)

	/* REPORTES COBRANZAS Y RENDICIONES (rentas)*/
	GetCobranzas(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseCobranzas, erro error)     /* Pagos */
	GetRendiciones(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseRendiciones, erro error) /* transferencias */

	NotificacionErroresReportes(errorFile []reportedtos.ResponseCsvEmailError) (erro error)

	/* REPORTES MOVIMIENTOS-COMISIONES */
	MovimientosComisionesService(request reportedtos.RequestReporteMovimientosComisiones) (res reportedtos.ResposeReporteMovimientosComisiones, erro error)

	/* REPORTES COBRANZAS-CLIENTES */
	GetCobranzasClientesService(request reportedtos.RequestCobranzasClientes) (res reportedtos.ResponseCobranzasClientes, erro error)
	/* REPORTES RENDICIONES-CLIENTES */
	GetRendicionesClientesService(request reportedtos.RequestReporteClientes) (res reportedtos.ResponseRendicionesClientes, erro error)
	/* REPORTES RENDICIONES-CLIENTES */
	GetReversionesClientesService(request reportedtos.RequestReporteClientes) (res reportedtos.ResponseReversionesClientes, erro error)

	// Reportes Informacion general
	GetPeticiones(request reportedtos.RequestPeticiones) (response reportedtos.ResponsePeticiones, erro error)
	GetLogs(request reportedtos.RequestLogs) (response reportedtos.ResponseLogs, erro error)
	GetNotificaciones(request reportedtos.RequestNotificaciones) (response reportedtos.ResponseNotificaciones, erro error)
	GetReportesEnviadosService(request reportedtos.RequestReportesEnviados) (response reportedtos.ResponseReportesEnviados, erro error)
	GetReportesPdfService(request reportedtos.RequestReportesEnviados, cliente entities.Cliente, cuenta entities.Cuenta) (response []reportedtos.ResponseClientesReportes, erro error)

	GetPagos(clientes administraciondtos.ResponseFacturacionPaginado, request reportedtos.RequestCobranzasDiarias) (response []reportedtos.ResponseClientesReportes, erro error)
	//Reporte Mensual Tratamiento Datos
	TratamientoReporteMensualPagos(reporteCliente []reportedtos.ResponseClientesReportes, orderCobranza bool) (reporte reportedtos.ReporteMensual, erro error)
	TratamientoReporteMensualRendiciones(reporteCliente []reportedtos.ResponseClientesReportes, orderCobranza bool) (reporte reportedtos.ReporteMensual, erro error)

	// RETENCIONES
	SendReporteRetencionComprobante(RRComprobante reportedtos.RequestRRComprobante) (control bool, erro error)
	//  Ejecutar todo el proceso de liquidacion de retenciones para uno o todos los clientes
	LiquidarRetencionesService(request reportedtos.RequestReportesEnviados) (erro error)
	//  Ejecutar todo el proceso de creacion de txt retenciones
	CreateTxtRetencionesService(request reportedtos.RequestReportesEnviados) (erro error)
	// crear, guardar y enviar txt retenciones SICORE
	CreateTxtRetencionesSICOREService(request reportedtos.RequestReportesEnviados) (rutatxt string, erro error)
	// crear, guardar y enviar txt retenciones SICAR
	CreateTxtRetencionesSICARService(request reportedtos.RequestReportesEnviados) (rutatxt string, erro error)
	// crear, guardar y enviar txt comisiones F 8125
	CreateTxtForm8125Service(request reportedtos.RequestReportesEnviados) (rutatxt string, erro error)
	// Crear ruta de archivo txt SICORE en proyecto. Crear archivo txt y guardar
	CreateTxtSICORE(lines []string, nombreTXT string, fechaFin string) (rutaDetalle string, erro error)
	// escribir lineas de texto en archivo SICORE txt de retenciones
	WriteTxtSICORE(archivo *os.File, lines []string) (erro error)
	// subir archivo individual a cloud storage
	UploadTxtFile(ctx context.Context, rutaArchivo, rutaCloud string, data []byte) (erro error)
	// recibe un entities.Comprobante. Evalua si sus ComprobanteDetalle tiene el atributo retener en true. Basta que un ComprobanteDetalle sea false para no retener
	EvaluarMinimoRetencionDeComprobante(comprobante entities.Comprobante) (result bool, erro error)

	GetCuentaByApiKeyService(apikey string) (cuenta *entities.Cuenta, erro error)

	CreateExcelRetencionesService(request reportedtos.RequestReportesEnviados) (rutatxt string, erro error)
}

type reportesService struct {
	repository     ReportesRepository
	administracion administracion.Service
	util           util.UtilService
	commons        commons.Commons
	factory        ReportesFactory
	factoryEmail   ReportesSendFactory
	store          util.Store
	emailService   email.Emailservice
}

func NewService(rm ReportesRepository, adm administracion.Service, util util.UtilService, c commons.Commons, storage util.Store, email email.Emailservice) ReportesService {
	reporte := &reportesService{
		repository: rm,
		// apilinkService:   link
		administracion: adm,
		util:           util,
		commons:        c,
		factory:        &procesarReportesFactory{},
		factoryEmail:   &enviaremailFactory{},
		store:          storage,
		emailService:   email,
	}
	return reporte

}

func (s *reportesService) GetCuentaByApiKeyService(apikey string) (cuenta *entities.Cuenta, erro error) {
	return s.repository.GetCuentaByApiKeyRepository(apikey)
}

func (s *reportesService) SendReporteReversiones(request []reportedtos.ResponseClientesReportes) (errorFile []reportedtos.ResponseCsvEmailError, erro error) {

	/* en esta ruta se crearan los archivos */
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	for _, cliente := range request {

		if len(cliente.Email) == 0 {
			erro = fmt.Errorf("no esta definido el email del cliente %v", cliente.Clientes)
			errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
				Archivo: "",
				Error:   fmt.Sprintf("error al enviar archivo: no esta definido email del cliente %v", cliente.Clientes),
			})
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "EnviarMailService",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return erro
			}
		} else {
			var tipo_archivo string
			var contentType string
			var asunto string
			var nombreArchivo string
			var tipo string
			var titulo string

			cliente.TipoReporte = "reversiones"
			asunto = "Rendiciones Mensuales"
			nombreArchivo = cliente.RutaArchivo
			titulo = "reversiones"
			tipo = "mensuales presentadas"

			tipo_archivo = ".pdf"
			contentType = "application/pdf"

			var campo_adicional = []string{"reversiones"}
			var email = cliente.Email //[]string{cliente.Email}
			filtro := utildtos.RequestDatosMail{
				Email:            email,
				Asunto:           asunto,
				From:             "Wee.ar!",
				Nombre:           cliente.Clientes,
				Mensaje:          "reportes de reversiones: #0",
				CamposReemplazar: campo_adicional,
				Descripcion: utildtos.DescripcionTemplate{
					Fecha:   cliente.Fecha,
					Cliente: cliente.RazonSocial,
					Cuit:    cliente.Cuit,
				},
				Totales: utildtos.TotalesTemplate{
					Titulo:       titulo,
					TipoReporte:  tipo,
					Elemento:     "reversiones",
					Cantidad:     cliente.CantOperaciones,
					TotalCobrado: cliente.TotalCobrado,
					TotalRendido: cliente.RendicionTotal,
				},
				AdjuntarEstado: true,
				Attachment: utildtos.Attachment{
					Name:        fmt.Sprintf("%s%s", nombreArchivo, tipo_archivo),
					ContentType: contentType,
					WithFile:    true,
				},
				TipoEmail: "reporte",
			}
			/*enviar archivo csv por correo*/
			/* en el caso de no registrar error al enviar correo se guardan los datos del reporte*/
			erro = s.util.EnviarMailService(filtro)
			logs.Info(erro)
			if erro != nil {
				erro = fmt.Errorf("no se no pudo enviar rendicion al %v", cliente.Clientes)
				errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
					Archivo: filtro.Attachment.Name,
					Error:   fmt.Sprintf("servicio email: %v", erro),
				})
				logs.Error(erro.Error())
				log := entities.Log{
					Tipo:          entities.EnumLog("Error"),
					Funcionalidad: "EnviarMailService",
					Mensaje:       erro.Error(),
				}
				erro = s.util.CreateLogService(log)
				if erro != nil {
					logs.Error("error: al crear logs: " + erro.Error())
					// return erro
				}
				/* informar el error al enviar el emial pero se debe continuar enviando los siguientes archivos a otros clientes */
			}

			if len(cliente.Reporte) == 1 {

				erro = s.repository.SaveGuardarDatosReporte(cliente.Reporte[0])
				if erro != nil {
					mensaje := fmt.Errorf("no se pudieron registrar datos del reporte de pago enviado al cliente %+v", cliente.Clientes).Error()
					logs.Info(mensaje)
					log := entities.Log{
						Tipo:          entities.EnumLog("Error"),
						Funcionalidad: "SaveGuardarDatosReporte",
						Mensaje:       mensaje,
					}
					erro = s.util.CreateLogService(log)
					if erro != nil {
						logs.Error("error: al crear logs: " + erro.Error())
					}

				}

			}

			// una vez enviado el correo se elimina el archivo csv
			erro = s.commons.BorrarArchivo(ruta, fmt.Sprintf(("%s%s"), nombreArchivo, tipo_archivo))
			if erro != nil {
				logs.Error(erro.Error())
				log := entities.Log{
					Tipo:          entities.EnumLog("Error"),
					Funcionalidad: "BorrarArchivos",
					Mensaje:       erro.Error(),
				}
				erro = s.util.CreateLogService(log)
				if erro != nil {
					logs.Error("error: al crear logs: " + erro.Error())
					// return nil, erro
				}
			}
		}

	}
	return
}

func formatearStringAFechaCorrecta(fecha string) (fechaFormateada string) {
	fechaSinEspacios := strings.Split(fecha, " ")
	anoMesDia := fechaSinEspacios[0]
	fechaEnPartes := strings.Split(anoMesDia, "-")
	fechaFormateada = fechaEnPartes[2] + "/" + fechaEnPartes[1] + "/" + fechaEnPartes[0]

	return
}

func (s *reportesService) GetPagosReportes(request reportedtos.RequestPagosPeriodo) (response []reportedtos.ResponsePagosPeriodo, erro error) {
	// 1 obtener estado pendiente para luego filtrar los pagos
	filtro := filtros.PagoEstadoFiltro{
		Nombre: "pending",
	}
	estadoPendiente, err := s.administracion.GetPagoEstado(filtro)
	if err != nil {
		erro = err
		return
	}

	// 2  obtener channels para luego filtrar cada pago con su cierre lote
	// debin
	canalDebin, erro := s.util.FirstOrCreateConfiguracionService("CHANNEL_DEBIN", "Nombre del canal debin", "debin")
	if erro != nil {
		return
	}
	filtroChannelDebin := filtros.ChannelFiltro{
		Channel: canalDebin,
	}

	channelDebin, erro := s.administracion.GetChannelService(filtroChannelDebin)

	if erro != nil && channelDebin.Id < 1 {
		return
	}

	// offline
	canalOffline, erro := s.util.FirstOrCreateConfiguracionService("CHANNEL_OFFLINE", "Nombre del canal debin", "offline")
	if erro != nil {
		return
	}
	filtroChannelOffline := filtros.ChannelFiltro{
		Channel: canalOffline,
	}

	channelOffline, erro := s.administracion.GetChannelService(filtroChannelOffline)

	if erro != nil && channelOffline.Id < 1 {
		return
	}

	// 3 obtener pagos del periodo
	pagos, erro := s.repository.GetPagosReportes(request, estadoPendiente[0].ID)
	if erro != nil {
		return
	}

	var listaPagoApilink []string
	var listaPagoOffline []string
	var listaPagoPrisma []string
	for _, pago := range pagos {
		// logs.Info(pago.ID)
		var valorCupon entities.Monto
		if pago.PagoIntentos[len(pago.PagoIntentos)-1].Valorcupon == 0 {
			valorCupon = pago.PagoIntentos[len(pago.PagoIntentos)-1].Amount
		} else {
			valorCupon = pago.PagoIntentos[len(pago.PagoIntentos)-1].Valorcupon
		}

		var fechaRendicion string
		var nroreferencia string
		var comision_porcentaje float64
		var comision_porcentaje_iva float64
		var importe_comision_sobre_tap float64
		var importe_comision_sobre_tap_iva float64
		var costo_fijo_transaccion float64
		var importe_rendido float64
		if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos) > 0 {
			if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientocomisions) > 0 {
				comision_porcentaje = s.util.ToFixed((pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientocomisions[0].Porcentaje * 100), 2)
				importe_comision_sobre_tap = s.util.ToFixed((pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientocomisions[0].Monto.Float64()), 2)
			}

			if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos) > 0 {
				comision_porcentaje_iva = pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos[0].Porcentaje * 100
				importe_comision_sobre_tap_iva = s.util.ToFixed((pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos[0].Monto.Float64()), 2)
			}

			if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos) > 0 {
				for _, mov := range pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos {
					if mov.Tipo == "D" {
						fechaRendicion = pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[1].Movimientotransferencia[0].FechaOperacion.Format("02-01-2006")
						nroreferencia = pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientos[1].Movimientotransferencia[0].ReferenciaBancaria
					} else {
						importe_rendido = s.util.ToFixed((mov.Monto.Float64()), 4)
					}
				}
			}

		}
		if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID == int64(channelDebin.Id) {
			costo_fijo_transaccion = float64(pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.Channel.Channelaranceles[0].Importe)
		}
		if !pago.PagoIntentos[len(pago.PagoIntentos)-1].PaidAt.IsZero() {
			fechaPago := pago.PagoIntentos[len(pago.PagoIntentos)-1].PaidAt.Format("02-01-2006")
			response = append(response, reportedtos.ResponsePagosPeriodo{
				Cliente:                 pago.PagosTipo.Cuenta.Cliente.Cliente,
				Cuenta:                  pago.PagosTipo.Cuenta.Cuenta,
				Pagotipo:                pago.PagosTipo.Pagotipo,
				IdPago:                  pago.ID,
				ExternalReference:       pago.ExternalReference,
				Estado:                  pago.PagoEstados.Nombre,
				ChannelId:               pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID,
				ExternalId:              pago.PagoIntentos[len(pago.PagoIntentos)-1].ExternalID,
				TransactionId:           pago.PagoIntentos[len(pago.PagoIntentos)-1].TransactionID,
				Barcode:                 pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode,
				IdExterno:               pago.ExternalReference,
				MedioPago:               pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.Mediopago,
				Pagador:                 strings.ToUpper(pago.PagoIntentos[len(pago.PagoIntentos)-1].HolderName),
				DniPagador:              pago.PagoIntentos[len(pago.PagoIntentos)-1].HolderNumber,
				Cuotas:                  uint(pago.PagoIntentos[len(pago.PagoIntentos)-1].Installmentdetail.Cuota),
				FechaPago:               fechaPago,
				FechaRendicion:          fechaRendicion,
				Amount:                  s.util.ToFixed((pago.PagoIntentos[len(pago.PagoIntentos)-1].Amount.Float64()), 4),
				AmountPagado:            s.util.ToFixed((valorCupon.Float64()), 4),
				CftCoeficiente:          uint(pago.PagoIntentos[len(pago.PagoIntentos)-1].Installmentdetail.Coeficiente),
				ComisionPorcentaje:      comision_porcentaje,
				ComisionPorcentajeIva:   comision_porcentaje_iva,
				ImporteComisionSobreTap: importe_comision_sobre_tap,
				ImporteIvaComisionTap:   importe_comision_sobre_tap_iva,
				CostoFijoTransaccion:    costo_fijo_transaccion,
				ImporteRendido:          importe_rendido,
				ReferenciaBancaria:      nroreferencia,
			})
			if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID == int64(channelDebin.Id) {
				listaPagoApilink = append(listaPagoApilink, pago.PagoIntentos[len(pago.PagoIntentos)-1].ExternalID)
			} else if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID == int64(channelOffline.Id) {
				listaPagoOffline = append(listaPagoOffline, pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode)
			} else {
				listaPagoPrisma = append(listaPagoPrisma, pago.PagoIntentos[len(pago.PagoIntentos)-1].TransactionID)
			}

		}
	}

	var listasPagos []reportedtos.TipoFactory
	listasPagos = append(listasPagos, reportedtos.TipoFactory{TipoApilink: listaPagoApilink}, reportedtos.TipoFactory{TipoOffline: listaPagoOffline}, reportedtos.TipoFactory{TipoPrisma: listaPagoPrisma})

	reportes, err := s.obtenerReportes(listasPagos)
	if err != nil {
		erro = err
		return
	}
	// actualzar pagos con valores obtenidos de cada cierrelote
	for i := range response {
		for j := range reportes {
			if response[i].ExternalId == reportes[j].Pago || response[i].TransactionId == reportes[j].Pago || response[i].Barcode == reportes[j].Pago {
				response[i].Nroestablecimiento = reportes[j].NroEstablecimiento
				response[i].NroLiquidacion = reportes[j].NroLiquidacion
				response[i].FechaPresentacion = reportes[j].FechaPresentacion
				response[i].FechaAcreditacion = reportes[j].FechaAcreditacion
				response[i].ArancelPorcentaje = reportes[j].ArancelPorcentaje
				response[i].RetencionIva = reportes[j].RetencionIva
				response[i].ImporteMinimo = reportes[j].Importeminimo
				response[i].ImporteMaximo = reportes[j].Importemaximo
				response[i].ArancelPorcentajeMinimo = reportes[j].ArancelPorcentajeMinimo
				response[i].ArancelPorcentajeMaximo = reportes[j].ArancelPorcentajeMaximo
				response[i].ImporteArancel = reportes[j].ImporteArancel
				response[i].ImporteArancelIva = reportes[j].ImporteArancelIva
				response[i].ImporteCft = reportes[j].ImporteCft
				response[i].ImporteNetoCobrado = reportes[j].ImporteNetoCobrado
				response[i].Revertido = reportes[j].Revertido
				response[i].Enobservacion = reportes[j].Enobservacion
				response[i].Cantdias = reportes[j].Cantdias
			}
		}
	}

	return

}

func (s *reportesService) ResultPagosReportes(request []reportedtos.ResponsePagosPeriodo, paginacion filtros.Paginacion) (response reportedtos.ResponseListaPagoPeriodo, erro error) {
	var responseTemporal []reportedtos.ResultadoPagosPeriodo
	var contador int64
	var recorrerHasta int32
	for _, listaPago := range request {
		contador++
		resp := reportedtos.ResultadoPagosPeriodo{
			Cliente:                 listaPago.Cliente,
			Cuenta:                  listaPago.Cuenta,
			Pagotipo:                listaPago.Pagotipo,
			ExternalReference:       listaPago.ExternalReference,
			IdPago:                  listaPago.IdPago,
			Estado:                  listaPago.Estado,
			MedioPago:               listaPago.MedioPago,
			Pagador:                 listaPago.Pagador,
			Dni:                     listaPago.DniPagador,
			Cuotas:                  listaPago.Cuotas,
			Nroestablecimiento:      listaPago.Nroestablecimiento,
			NroLiquidacion:          listaPago.NroLiquidacion,
			FechaPago:               listaPago.FechaPago,
			FechaPresentacion:       listaPago.FechaPresentacion,
			FechaAcreditacion:       listaPago.FechaAcreditacion,
			FechaRendicion:          listaPago.FechaRendicion,
			Amount:                  listaPago.Amount,
			AmountPagado:            listaPago.AmountPagado,
			ArancelPorcentaje:       listaPago.ArancelPorcentaje,
			CftCoeficiente:          listaPago.CftCoeficiente,
			RetencionIva:            listaPago.RetencionIva,
			ImporteMinimo:           listaPago.ImporteMinimo,
			ImporteMaximo:           listaPago.ImporteMaximo,
			ArancelPorcentajeMinimo: listaPago.ArancelPorcentajeMinimo,
			ArancelPorcentajeMaximo: listaPago.ArancelPorcentajeMaximo,
			CostoFijoTransaccion:    listaPago.CostoFijoTransaccion,
			ImporteArancel:          listaPago.ImporteArancel,
			ImporteArancelIva:       listaPago.ImporteArancelIva,
			ImporteCft:              listaPago.ImporteCft,
			ComisionPorcentaje:      listaPago.ComisionPorcentaje,
			ComisionPorcentajeIva:   listaPago.ComisionPorcentajeIva,
			ImporteComisionSobreTap: listaPago.ImporteComisionSobreTap,
			ImporteIvaComisionTap:   listaPago.ImporteIvaComisionTap,
			ImporteRendido:          listaPago.ImporteRendido,
			ImporteNetoCobrado:      listaPago.ImporteNetoCobrado,
			ReferenciaBancaria:      listaPago.ReferenciaBancaria,
			Revertido:               listaPago.Revertido,
			Enobservacion:           listaPago.Enobservacion,
			Cantdias:                listaPago.Cantdias,
		}
		responseTemporal = append(responseTemporal, resp)
	}
	if paginacion.Number > 0 && paginacion.Size > 0 {
		response.Meta = _setPaginacion(paginacion.Number, paginacion.Size, contador)
	}
	recorrerHasta = response.Meta.Page.To
	if response.Meta.Page.CurrentPage == response.Meta.Page.LastPage {
		recorrerHasta = response.Meta.Page.Total
	}
	if recorrerHasta == 0 {
		recorrerHasta = int32(contador)
	}

	if len(responseTemporal) > 0 {
		for i := response.Meta.Page.From; i < recorrerHasta; i++ {
			response.PagosByPeriodo = append(response.PagosByPeriodo, responseTemporal[i])
			response.TotalImporteRendidio += responseTemporal[i].ImporteRendido
		}
		response.TotalImporteRendidio = s.util.ToFixed((response.TotalImporteRendidio), 2)
	}
	return
}

func (s *reportesService) obtenerReportes(listaPagos []reportedtos.TipoFactory) (response []reportedtos.ResponseFactory, erro error) {

	for _, listaPago := range listaPagos {
		var tipoReporte string
		if len(listaPago.TipoApilink) > 0 {
			tipoReporte = "debin"
		} else if len(listaPago.TipoOffline) > 0 {
			tipoReporte = "offline"
		} else if len(listaPago.TipoPrisma) > 0 {
			tipoReporte = "prisma"
		}

		if tipoReporte != "" {
			metodoProcesarReporte, err := s.factory.GetProcesarReportes(tipoReporte)
			if err != nil {
				erro = err
				return
			}

			logs.Info("Procesando reportes tipo: " + tipoReporte)

			listaReporteProcesada := metodoProcesarReporte.ResponseReportes(s, listaPago)

			response = append(response, listaReporteProcesada...)
		}
	}

	return
}

func (s *reportesService) GetClientes(request reportedtos.RequestPagosClientes) (response administraciondtos.ResponseFacturacionPaginado, erro error) {
	filtro := filtros.ClienteFiltro{
		Id:              request.Cliente,
		CargarContactos: true,
		CargarCuentas:   true,
		ClientesIds:     request.ClientesIds,
	}
	response, erro = s.administracion.GetClientesService(filtro)
	return
}

func (s *reportesService) GetPagosClientes(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error) {
	// SE DEFINEN VARIABLES TOTALES
	var fechaI time.Time
	var fechaF time.Time

	if filtro.FechaInicio.IsZero() {
		// Entro por proceso background
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}
	for _, cliente := range request.Clientes {
		var cantoperaciones uint
		var totalcomisiones, totaliva, totalretencion, totalcobrado entities.Monto

		filtroPagos := reportedtos.RequestPagosPeriodo{
			ClienteId:   uint64(cliente.Id),
			FechaInicio: fechaI,
			FechaFin:    fechaF,
		}

		// listaPagos, err := s.repository.GetPagosReportes(filtroPagos, estadoPendiente[0].ID)

		var listaPagos []reportedtos.DetallesPagosCobranza

		apilink, err := s.repository.GetCobranzasApilink(filtroPagos)
		prisma, err := s.repository.GetCobranzasPrisma(filtroPagos)
		rapipago, err := s.repository.GetCobranzasRapipago(filtroPagos)
		multipago, err := s.repository.GetCobranzasMultipago(filtroPagos)

		if err != nil {
			erro = err
			return
		}

		listaPagos = append(listaPagos, apilink...)
		listaPagos = append(listaPagos, prisma...)
		listaPagos = append(listaPagos, rapipago...)
		listaPagos = append(listaPagos, multipago...)

		var pagos []reportedtos.PagosReportes
		if len(listaPagos) > 0 {
			for _, pago := range listaPagos {
				iva := entities.Monto(pago.Iva)
				retencion := entities.Monto(pago.Retencion)
				comision := entities.Monto(pago.ComisionTotal)
				monto := entities.Monto(pago.TotalPago)

				// totales
				totalcobrado += monto
				totaliva += iva
				totalcomisiones += comision
				totalretencion += retencion
				cantoperaciones++

				fecha := pago.FechaPago.Format("02-01-2006")
				if pago.CanalPago == "DEBIN" || pago.CanalPago == "OFFLINE" {
					fecha = pago.FechaCobro.Format("02-01-2006")
				}

				pagos = append(pagos, reportedtos.PagosReportes{
					Cuenta:    pago.Cuenta,
					Id:        pago.Referencia,
					FechaPago: fecha,
					MedioPago: pago.MedioPago,
					Tipo:      pago.CanalPago,
					Estado:    pago.Pagoestado,
					Monto:     fmt.Sprintf("%v", s.util.FormatNum(monto.Float64())),
					Comision:  fmt.Sprintf("%v", s.util.FormatNum(comision.Float64())),
					Iva:       fmt.Sprintf("%v", s.util.FormatNum(iva.Float64())),
					Retencion: fmt.Sprintf("%v", s.util.FormatNum(retencion.Float64())),
				})
			}

		}
		if len(pagos) > 0 {
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:        cliente.Cliente,
				Email:           cliente.Emails, //[]string{cliente.Email},
				RazonSocial:     cliente.RazonSocial,
				Cuit:            cliente.Cuit,
				SujetoRetencion: cliente.SujetoRetencion,
				Fecha:           fechaI.Format("02-01-2006"),
				Pagos:           pagos,
				CantOperaciones: fmt.Sprintf("%v", cantoperaciones),
				TotalCobrado:    fmt.Sprintf("%v", s.util.FormatNum(totalcobrado.Float64())),
				TotalComision:   fmt.Sprintf("%v", s.util.FormatNum(totalcomisiones.Float64())),
				TotalIva:        fmt.Sprintf("%v", s.util.FormatNum(totaliva.Float64())),
				TotalRetencion:  fmt.Sprintf("%v", s.util.FormatNum(totalretencion.Float64())),
				//TipoArchivoPdf:      true,
				GuardarDatosReporte: true,
			})

			// fechaNombre := fechaI.Format("02-01-2006")
			// pagosReporte := reportedtos.ToEntityRegistroReporte(response[len(response)-1])
			// siguienteNroReporte, err := s.repository.GetLastReporteEnviadosRepository(pagosReporte, false)
			// if err != nil {
			// 	erro = err
			// 	return
			// }
			// nroReporteString := strconv.FormatUint(uint64(siguienteNroReporte), 10)
			// err = GetPagosPdf(response[len(response)-1], cliente, fechaNombre, nroReporteString)

			// if err != nil {
			// 	erro = err
			// 	return
			// }

		}
	}

	return
}

func (s *reportesService) GetRendicionClientes(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error) {
	var (
		fechaI, fechaF time.Time
	)

	if filtro.FechaInicio.IsZero() {
		// si los filtros recibidos son ceros toman la fecha actual
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		// a las fechas se le restan un dia ya sea por backgraund o endpoint
		// fechaI = fechaI.AddDate(0, 0, int(-1))
		// fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		// fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		// fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
		fechaI = filtro.FechaInicio
		fechaF = filtro.FechaFin
	}
	logs.Info(fechaI)
	logs.Info(fechaF)

	// Para cada cliente se realiza toda la logica
	for _, cliente := range request.Clientes {

		// Filtro para transferencias, por cliente id, y rango de fechas
		filtro := reportedtos.RequestPagosPeriodo{
			ClienteId: uint64(cliente.Id),
			// ClienteId:   5,                             //Prueba con cliente 6
			FechaInicio: fechaI, // descomentar esta linea cuando se pasa a dev y produccion
			// se envian pagos del dia anterior
			FechaFin: fechaF,
			CuentaId: uint64(filtro.CuentaId),
		}

		// una lista de objetos entidad transferencia
		listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
		if err != nil {
			erro = err
			s.util.BuildLog(erro, "GetRendicionClientes")
			return
		}

		var (
			totalCobrado, total, totalReversion, totalIVA, totalComision, rendido, totalRetencion entities.Monto
			cantOperaciones                                                                       int
			transferencias                                                                        []*reportedtos.ResponseReportesRendiciones
			filtroMov                                                                             reportedtos.RequestPagosPeriodo
			totalCliente                                                                          reportedtos.ResponseTotales
			movrevertidos                                                                         []entities.Movimiento
			pagosintentos, pagosintentosrevertidos                                                []uint64
			listaMovimientoIds                                                                    []uint
		)

		// se recorren cada una de las transferencias
		if len(listaTransferencia) > 0 {
			for _, transferencia := range listaTransferencia {
				// se acumulan pagointentos de movs que no son reversion
				if !transferencia.Reversion {
					pagosintentos = append(pagosintentos, transferencia.Movimiento.PagointentosId)
				}
				// se acumulan pagointentos de movs que SI son reversion
				if transferencia.Reversion {
					pagosintentosrevertidos = append(pagosintentosrevertidos, transferencia.Movimiento.PagointentosId)
				}
			}
		}

		// en el caso de que existieran pagosintentosrevertidos, se buscan movimientos reversiones
		if len(pagosintentosrevertidos) > 0 {
			filtroRevertidos := reportedtos.RequestPagosPeriodo{
				PagoIntentos:                    pagosintentosrevertidos,
				TipoMovimiento:                  "C",
				CargarMovimientosTransferencias: true,
				CargarPagoIntentos:              true,
				CargarCuenta:                    true,
				CargarReversionReporte:          true,
				CargarComisionImpuesto:          true,
			}
			// se obtienen movimientos REVERTIDOS con un filtro que incluye una lista de pagosintentosrevertidos y el tipo C de movimientos
			movrevertidos, err = s.repository.GetMovimiento(filtroRevertidos)
			if err != nil {
				erro = err
				return
			}
		}

		// en el caso de que existieran pagosintentos, se buscan movimientos positivos
		if len(pagosintentos) > 0 {
			filtroMov = reportedtos.RequestPagosPeriodo{
				PagoIntentos:                    pagosintentos,
				TipoMovimiento:                  "C",
				CargarComisionImpuesto:          true,
				CargarMovimientosTransferencias: true,
				CargarPagoIntentos:              true,
				CargarCuenta:                    true,
				CargarMovimientosRetenciones:    true,
			}
			// se obtienen movimientos con un filtro que incluye una lista de pagointentos y el tipo C de movimientos
			mov, err := s.repository.GetMovimiento(filtroMov)
			if err != nil {
				erro = err
				return
			}
			var resulRendiciones []*reportedtos.ResponseReportesRendiciones

			// INICIO de for range movimientos tipo C
			for _, m := range mov {
				// guardar los id de los movimientos en un slice
				listaMovimientoIds = append(listaMovimientoIds, m.ID)
				cantOperaciones = cantOperaciones + 1
				cantidadBoletas := len(m.Pagointentos.Pago.Pagoitems)
				// acumulador montos de movimientos
				total += m.Monto
				// acumulador montos pagados
				totalCobrado += m.Pagointentos.Amount
				var comision entities.Monto
				var iva entities.Monto
				// si existen movimientos comisiones
				if len(m.Movimientocomisions) > 0 {
					comision = m.Movimientocomisions[len(m.Movimientocomisions)-1].Monto + m.Movimientocomisions[len(m.Movimientocomisions)-1].Montoproveedor
					iva = m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Monto + m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Montoproveedor
				} else {
					comision = 0
					iva = 0
				}
				// acumulador montos comision
				totalComision += comision
				// acumulador montos IVA
				totalIVA += iva

				// resulRendiciones slice de ResponseReportesRendiciones
				// incluye los campos que se muestran como columnas en el reporte de rendicion
				resulRendiciones = append(resulRendiciones, &reportedtos.ResponseReportesRendiciones{
					PagoIntentoId:           m.PagointentosId,
					Cuenta:                  m.Cuenta.Cuenta,
					Id:                      m.Pagointentos.Pago.ExternalReference,
					FechaCobro:              m.Pagointentos.PaidAt.Format("02-01-2006"),
					ImporteCobrado:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Pagointentos.Amount.Float64(), 2))),
					ImporteDepositado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Monto.Float64(), 2))),
					CantidadBoletasCobradas: fmt.Sprintf("%v", cantidadBoletas),
					Comision:                fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 4))),
					Iva:                     fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 4))),
					Concepto:                "Transferencia",
				})
			}
			// FIN de for range movimientos tipo C

			totalCliente = reportedtos.ResponseTotales{
				// Totales: reportedtos.Totales{
				// 	CantidadOperaciones: fmt.Sprintf("%v", cantOperaciones),
				// 	TotalCobrado:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalCobrado.Float64(), 4))),
				// 	TotalRendido:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(total.Float64(), 4))),
				// },
				Detalles: resulRendiciones,
			}
		}

		// si existen movimientos revertidos
		if len(movrevertidos) > 0 {
			for _, mr := range movrevertidos {
				// guardar los id de los movimientos en un slice
				listaMovimientoIds = append(listaMovimientoIds, mr.ID)
				totalReversion += mr.Monto
				cantOperaciones = cantOperaciones + 1
				cantidadBoletas := len(mr.Pagointentos.Pago.Pagoitems)
				var comision entities.Monto
				var iva entities.Monto
				if len(mr.Movimientocomisions) > 0 {
					comision = mr.Movimientocomisions[len(mr.Movimientocomisions)-1].Monto + mr.Movimientocomisions[len(mr.Movimientocomisions)-1].Montoproveedor
					iva = mr.Movimientoimpuestos[len(mr.Movimientoimpuestos)-1].Monto + mr.Movimientoimpuestos[len(mr.Movimientoimpuestos)-1].Montoproveedor
				} else {
					comision = 0
					iva = 0
				}
				// acumulador montos comision revertido. Se suma porque el importe ya es negativo
				totalComision += comision
				// acumulador montos IVA revertidos. Se suma porque el importe ya es negativo
				totalIVA += iva

				totalCliente.Detalles = append(totalCliente.Detalles, &reportedtos.ResponseReportesRendiciones{
					PagoIntentoId:           mr.PagointentosId,
					Cuenta:                  mr.Cuenta.Cuenta,
					Id:                      mr.Pagointentos.Pago.ExternalReference,
					FechaCobro:              mr.Pagointentos.PaidAt.Format("02-01-2006"),
					ImporteCobrado:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(mr.Pagointentos.Amount.Float64(), 2))),
					ImporteDepositado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(mr.Monto.Float64(), 2))),
					CantidadBoletasCobradas: fmt.Sprintf("%v", cantidadBoletas),
					Comision:                fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 4))),
					Iva:                     fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 4))),
					Concepto:                "Reversion",
				})
			}
		}

		// vuelvo a comparar con las transferenicas para asignar fecha de deposito
		if len(totalCliente.Detalles) > 0 {
			for _, transferencia := range listaTransferencia {
				for _, t := range totalCliente.Detalles {
					if transferencia.Movimiento.PagointentosId == t.PagoIntentoId {
						t.FechaDeposito = transferencia.FechaOperacion.Format("02-01-2006")
						t.CBUOrigen = transferencia.CbuOrigen
						t.CBUDestino = transferencia.CbuDestino
						t.ReferenciaBancaria = transferencia.ReferenciaBancaria
					}
				}
			}
		}

		// Total de las retenciones de los movimientos tipo C a partir de una lista de ids de movimientos
		filtroRetencionesMovimientos := retenciondtos.RentencionRequestDTO{
			ClienteId:          uint(filtro.ClienteId),
			ListaMovimientosId: listaMovimientoIds,
		}

		// buscar y agrupar retenciones por gravamen a partir de una lista de movimientos
		var retenciones_evaluadas []retenciondtos.RetencionAgrupada
		retenciones_evaluadas, erro = s.administracion.EvaluarRetencionesByMovimientoService(filtroRetencionesMovimientos)
		if erro != nil {
			s.util.BuildLog(erro, "GetRendicionClientes")
			return
		}

		// map para asociar nombre de gravamen e importe total por gravamen
		retenciones := make(map[string]entities.Monto)
		// gravamenes agrupados por tipo, en formato entities.Monto
		for _, ra := range retenciones_evaluadas {
			retenciones[ra.Gravamen] = ra.TotalRetencion
		}

		transferencias = totalCliente.Detalles
		// total de montos de movimientos + total montos movimientos revertidos
		rendido = total + totalReversion

		// total de las RETENCIONES
		for _, retencion_por_gravamen := range retenciones {
			totalRetencion += retencion_por_gravamen
		}

		// totales de reporte
		totalCliente = reportedtos.ResponseTotales{
			Totales: reportedtos.Totales{
				CantidadOperaciones:    fmt.Sprintf("%v", cantOperaciones),
				TotalCobrado:           fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalCobrado.Float64(), 4))),
				TotalRendido:           fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(rendido.Float64(), 4))),
				TotalComision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalComision.Float64(), 4))),
				TotalIva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalIVA.Float64(), 4))),
				TotalRetGanancias:      fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(retenciones["ganancias"].Float64(), 4))),
				TotalRetIva:            fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(retenciones["iva"].Float64(), 4))),
				TotalRetIngresosBrutos: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(retenciones["iibb"].Float64(), 4))),
				TotalRetencion:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalRetencion.Float64(), 4))),
			},
		}

		if len(transferencias) > 0 {
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:                cliente.Cliente,
				RazonSocial:             cliente.RazonSocial,
				Cuit:                    cliente.Cuit,
				Email:                   cliente.Emails,
				Fecha:                   fechaI.Format("02-01-2006"),
				Rendiciones:             transferencias,
				CantOperaciones:         totalCliente.Totales.CantidadOperaciones,
				TotalCobrado:            totalCliente.Totales.TotalCobrado,
				TotalIva:                totalCliente.Totales.TotalIva,
				TotalComision:           totalCliente.Totales.TotalComision,
				RendicionTotal:          totalCliente.Totales.TotalRendido,
				TipoArchivoPdf:          true,
				GuardarDatosReporte:     true,
				TotalRetencionGanancias: totalCliente.Totales.TotalRetGanancias,
				TotalRetencionIva:       totalCliente.Totales.TotalRetIva,
				TotalRetencionIIBB:      totalCliente.Totales.TotalRetIngresosBrutos,
				TotalRetencion:          totalCliente.Totales.TotalRetencion,
			})

			fechaNombre := fechaI.Format("02-01-2006")
			transferenciasReporte := reportedtos.ToEntityRegistroReporte(response[len(response)-1])
			var filtroReporte filtros_reportes.BusquedaReporteFiltro
			siguienteNroReporte, err := s.repository.GetLastReporteEnviadosRepository(transferenciasReporte, filtroReporte)
			if err != nil {
				erro = err
				return
			}

			nroReporteString := strconv.FormatUint(uint64(siguienteNroReporte), 10)
			// crear reporte de rendiciones en PDF
			err = GetRendicionesPdf(response[len(response)-1], cliente, fechaNombre, nroReporteString)

			if err != nil {
				erro = err
				s.util.BuildLog(erro, "GetRendicionClientes")
			}
		}
	} // fin de para cada cliente, crear un reporte de rendiciones

	return
}

func (s *reportesService) GetReversionesClientes(requestCliente administraciondtos.ResponseFacturacionPaginado, request reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error) {

	for _, cliente := range requestCliente.Clientes {
		var fechaI time.Time
		var fechaF time.Time
		if request.FechaInicio.IsZero() {
			// Entro por proceso background
			fechaInicio, fechaFin, err := s.commons.FormatFecha()
			if err != nil {
				erro = err
				return
			}
			fechaI = fechaInicio
			fechaF = fechaFin
		} else {
			fechaI = request.FechaInicio
			fechaF = request.FechaFin
		}

		// fechaI, fechaF, err := s.commons.FormatFecha()
		// if err != nil {
		// 	return
		// }
		filtro := reportedtos.RequestPagosPeriodo{
			ClienteId:       uint64(cliente.Id),
			FechaInicio:     fechaI,
			FechaFin:        fechaF,
			CargarReversion: true,
		}
		logs.Info(request.FechaInicio)

		/* listaPagos, err := s.repository.GetReversionesReportes(filtro, filtroValidacion)
		if err != nil {
			erro = err
			return
		} */

		transferencias, err := s.repository.GetReversionesDeTransferenciaClientes(filtro)
		if err != nil {
			erro = err
			return
		}
		var reverciones []reportedtos.Reversiones
		var cantOperacion int64
		var totalRevertido float64
		if len(transferencias) > 0 {
			for _, value := range transferencias {
				cantOperacion = cantOperacion + 1
				var revertido reportedtos.Reversiones
				var pagoRevertido reportedtos.PagoRevertido
				var itemsRevertido []reportedtos.ItemsRevertidos
				//var itemRevertido reportedtos.ItemsRevertidos
				var intentoPagoRevertido reportedtos.IntentoPagoRevertido

				pagoRevertido.EntityToPagoRevertido(value.Movimiento.Pagointentos.Pago)
				if len(value.Movimiento.Pagointentos.Pago.Pagoitems) > 0 {
					for _, valueItem := range value.Movimiento.Pagointentos.Pago.Pagoitems {
						var itemRevertido reportedtos.ItemsRevertidos
						itemRevertido.EntityToItemsRevertidos(valueItem)
						itemsRevertido = append(itemsRevertido, itemRevertido)
					}
				}

				revertido.Fecha = formatearStringAFechaCorrecta(value.CreatedAt.String())
				revertido.Id = value.Referencia
				revertido.Cuenta = value.Movimiento.Cuenta.Cuenta
				revertido.MedioPago = value.Movimiento.Pagointentos.Mediopagos.Mediopago

				montoitem_float := value.Movimiento.Pagointentos.Amount.Float64()
				revertido.Monto = "$ " + util.Resolve().FormatNum(montoitem_float)

				totalRevertido += value.Movimiento.Monto.Float64()

				intentoPagoRevertido.EntityToIntentoPagoRevertido(*value.Movimiento.Pagointentos)
				pagoRevertido.Items = itemsRevertido
				pagoRevertido.IntentoPago = intentoPagoRevertido

				montointento_float := value.Movimiento.Pagointentos.Amount.Float64()
				pagoRevertido.IntentoPago.ImportePagado = "$ " + util.Resolve().FormatNum(montointento_float)

				revertido.PagoRevertido = pagoRevertido

				reverciones = append(reverciones, revertido)
			}
		}

		/* if len(transferencias) > 0 {
			for _, pago := range transferencias {
				reverciones = append(reverciones, reportedtos.Reversiones{
					Cuenta:    pago.Movimiento.Cuenta.Cuenta,
					Id:        pago.Movimiento.Pagointentos.Pago.ExternalReference,
					MedioPago: pago.Movimiento.Pagointentos.Mediopagos.Mediopago,
					Monto:     fmt.Sprintf("%v", s.util.ToFixed(entities.Monto(pago.Movimiento.Pagointentos.Amount).Float64(), 4)),
				})
			}
		} */
		if len(reverciones) > 0 {
			fechaEnPartes := strings.Split(fechaI.Format("02-01-2006"), "-")
			fechaFormateada := fechaEnPartes[0] + "-" + fechaEnPartes[1] + "-" + fechaEnPartes[2]

			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:            cliente.Cliente,
				RazonSocial:         cliente.RazonSocial,
				Cuit:                cliente.Cuit,
				Email:               cliente.Emails,
				Fecha:               fechaFormateada, //fechaI.AddDate(0, 0, int(-1)).Format("02-01-2006"),
				Reversiones:         reverciones,
				TotalRevertido:      fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalRevertido, 4))),
				CantOperaciones:     fmt.Sprintf("%v", cantOperacion),
				TipoArchivoPdf:      true,
				GuardarDatosReporte: true,
			})
		}
	}

	// transformar la data adaptando al pdf

	for i, cliente := range response {
		reversionesData := transformarDatos(response[i])
		clienteData := commons.ClienteData{
			Clientes:    cliente.Clientes,
			RazonSocial: cliente.RazonSocial,
			Cuit:        cliente.Cuit,
		}
		err := commons.GetReversionesPdf(reversionesData, clienteData, cliente.Fecha)
		if err != nil {
			erro = err
			logs.Error(err.Error())
		}
	}

	return
}

func (s *reportesService) SendPagosClientes(request []reportedtos.ResponseClientesReportes) (errorFile []reportedtos.ResponseCsvEmailError, erro error) {

	/* en esta ruta se crearan los archivos */
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	for _, cliente := range request {

		if len(cliente.Email) == 0 {
			erro = fmt.Errorf("no esta definido el email del cliente %v", cliente.Clientes)
			errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
				Archivo: "",
				Error:   fmt.Sprintf("error al enviar archivo: no esta definido email del cliente %v", cliente.Clientes),
			})
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "EnviarMailService",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return erro
			}
		} else {
			var tipo_archivo string
			var contentType string
			var asunto string
			var nombreArchivo string
			var tipo string
			var titulo string
			if len(cliente.Pagos) > 0 {
				cliente.TipoReporte = "pagos"
				asunto = "Pagos realizados " + cliente.Fecha
				nombreArchivo = cliente.Clientes + "-" + cliente.Fecha
				titulo = "cobranzas"
				tipo = "cobrados"
			} else if len(cliente.Rendiciones) > 0 {
				cliente.TipoReporte = "rendiciones"
				asunto = "Recaudación WEE! " + cliente.Fecha
				nombreArchivo = cliente.Clientes + "-" + cliente.Fecha
				titulo = "rendiciones"
				tipo = "rendidos"
			} else if len(cliente.Reversiones) > 0 {
				cliente.TipoReporte = "revertidos"
				asunto = "Pagos revertidos " + cliente.Fecha
				nombreArchivo = cliente.Clientes + "-" + cliente.Fecha
				titulo = "reversiones"
				tipo = "revertidos"
			}
			if cliente.TipoArchivoPdf {
				tipo_archivo = ".pdf"
				contentType = "application/pdf"
			} else {
				tipo_archivo = ".csv"
				contentType = "application/pdf"
				metodoConvertirCvs, err := s.factoryEmail.SendEnviarEmail(cliente.TipoReporte)
				if err != nil {
					erro = err
					return
				}
				logs.Info("Procesando reportes tipo: " + cliente.TipoReporte)
				convertircvs := metodoConvertirCvs.SendReportes(ruta, nombreArchivo, cliente)
				if convertircvs != nil {
					erro = convertircvs
					return
				}
			}

			var campo_adicional = []string{"pagos"}
			var email = cliente.Email //[]string{cliente.Email}
			filtro := utildtos.RequestDatosMail{
				Email:            email,
				Asunto:           asunto,
				From:             "Wee.ar!",
				Nombre:           cliente.Clientes,
				Mensaje:          "reportes de pagos: #0",
				CamposReemplazar: campo_adicional,
				Descripcion: utildtos.DescripcionTemplate{
					Fecha:   cliente.Fecha,
					Cliente: cliente.RazonSocial,
					Cuit:    cliente.Cuit,
				},
				Totales: utildtos.TotalesTemplate{
					Titulo:         titulo,
					TipoReporte:    tipo,
					Elemento:       "pagos",
					Cantidad:       cliente.CantOperaciones,
					TotalCobrado:   cliente.TotalCobrado,
					TotalComision:  cliente.TotalComision,
					TotalIva:       cliente.TotalIva,
					TotalRendido:   cliente.RendicionTotal,
					TotalRevertido: cliente.TotalRevertido,
				},
				AdjuntarEstado: true,
				Attachment: utildtos.Attachment{
					Name:        fmt.Sprintf("%s%s", nombreArchivo, tipo_archivo),
					ContentType: contentType,
					WithFile:    true,
				},
				TipoEmail: "reporte",
			}
			if len(cliente.Rendiciones) > 0 {
				filtro.Totales.Rendicion = true
				filtro.Totales.CBUDestino = cliente.Rendiciones[0].CBUDestino
				filtro.Totales.CBUOrigen = cliente.Rendiciones[0].CBUOrigen
				filtro.Totales.ReferenciaBancaria = cliente.Rendiciones[0].ReferenciaBancaria
			}

			if cliente.SujetoRetencion {
				filtro.Totales.TotalRetencion = cliente.TotalRetencion
			}

			/*enviar archivo csv por correo*/
			/* en el caso de no registrar error al enviar correo se guardan los datos del reporte*/
			erro = s.util.EnviarMailService(filtro)
			logs.Info(erro)
			if erro != nil {
				erro = fmt.Errorf("no se no pudo enviar rendicion al %v", cliente.Clientes)
				errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
					Archivo: filtro.Attachment.Name,
					Error:   fmt.Sprintf("servicio email: %v", erro),
				})
				logs.Error(erro.Error())
				log := entities.Log{
					Tipo:          entities.EnumLog("Error"),
					Funcionalidad: "EnviarMailService",
					Mensaje:       erro.Error(),
				}
				erro = s.util.CreateLogService(log)
				if erro != nil {
					logs.Error("error: al crear logs: " + erro.Error())
					// return erro
				}
				/* informar el error al enviar el emial pero se debe continuar enviando los siguientes archivos a otros clientes */
			} else {

				if cliente.GuardarDatosReporte {
					// guardar datos del reporte
					//si el archivo se sube correctamente se registra en tabla movimientos lotes
					pagos := reportedtos.ToEntityRegistroReporte(cliente)
					var filtroReporte filtros_reportes.BusquedaReporteFiltro
					siguienteNroReporte, erro := s.repository.GetLastReporteEnviadosRepository(pagos, filtroReporte)
					if erro != nil {
						mensaje := fmt.Errorf("no se pudo obtener nro reporte para el reporte de pago enviado al cliente %+v", cliente.Clientes).Error()
						logs.Info(mensaje)
						log := entities.Log{
							Tipo:          entities.EnumLog("Error"),
							Funcionalidad: "GetLastReporteEnviadosRepository",
							Mensaje:       mensaje,
						}
						erro = s.util.CreateLogService(log)
						if erro != nil {
							logs.Error("error: al crear logs: " + erro.Error())
						}
					}
					pagos.Nro_reporte = siguienteNroReporte
					erro = s.repository.SaveGuardarDatosReporte(pagos)
					if erro != nil {
						mensaje := fmt.Errorf("no se pudieron registrar datos del reporte de pago enviado al cliente %+v", cliente.Clientes).Error()
						logs.Info(mensaje)
						log := entities.Log{
							Tipo:          entities.EnumLog("Error"),
							Funcionalidad: "SaveGuardarDatosReporte",
							Mensaje:       mensaje,
						}
						erro = s.util.CreateLogService(log)
						if erro != nil {
							logs.Error("error: al crear logs: " + erro.Error())
						}

					}

				}
			}

			// una vez enviado el correo se elimina el archivo csv
			erro = s.commons.BorrarArchivo(ruta, fmt.Sprintf(("%s%s"), nombreArchivo, tipo_archivo))
			if erro != nil {
				logs.Error(erro.Error())
				log := entities.Log{
					Tipo:          entities.EnumLog("Error"),
					Funcionalidad: "BorrarArchivos",
					Mensaje:       erro.Error(),
				}
				erro = s.util.CreateLogService(log)
				if erro != nil {
					logs.Error("error: al crear logs: " + erro.Error())
					// return nil, erro
				}
			}
		}

	}
	// erro = s.commons.BorrarDirectorio(ruta)
	// if erro != nil {
	// 	logs.Error(erro.Error())
	// 	log := entities.Log{
	// 		Tipo:          entities.EnumLog("Error"),
	// 		Funcionalidad: "BorrarDirectorio",
	// 		Mensaje:       erro.Error(),
	// 	}
	// 	erro = s.util.CreateLogService(log)
	// 	if erro != nil {
	// 		logs.Error("error: al crear logs: " + erro.Error())
	// 		// return erro
	// 	}
	// }

	return
}

func (s *reportesService) SendReporteRendiciones(request []reportedtos.ResponseClientesReportes) (errorFile []reportedtos.ResponseCsvEmailError, numero_reporte_rrm uint, erro error) {

	/* en esta ruta se crearan los archivos */
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	for _, cliente := range request {

		if len(cliente.Email) == 0 {
			erro = fmt.Errorf("no esta definido el email del cliente %v", cliente.Clientes)
			errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
				Archivo: "",
				Error:   fmt.Sprintf("error al enviar archivo: no esta definido email del cliente %v", cliente.Clientes),
			})
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "SendReporteRendiciones",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return erro
			}
		} else {

			if len(cliente.Email) > 0 {
				var tipo_archivo string
				var contentType string
				var asunto string
				var nombreArchivo string
				var tipo string
				var titulo string

				cliente.TipoReporte = "rendiciones"
				asunto = "Rendiciones Mensuales"
				nombreArchivo = cliente.RutaArchivo
				titulo = "rendiciones"
				tipo = "mensuales presentadas"

				tipo_archivo = ".pdf"
				contentType = "application/pdf"

				var campo_adicional = []string{"rendiciones"}
				// var email = cliente.Email // se comenta el envio al email del cliente
				// 1era Liquidacion: se hardcodean los emails
				var email = []string{"pablo.vicentin@telco.com.ar", "sebastian.escobar@telco.com.ar", "yasmir.yaya@telco.com.ar"}

				filtro := utildtos.RequestDatosMail{
					Email:            email,
					Asunto:           asunto,
					From:             "Wee.ar!",
					Nombre:           cliente.Clientes,
					Mensaje:          "reportes de rendiciones: #0",
					CamposReemplazar: campo_adicional,
					Descripcion: utildtos.DescripcionTemplate{
						Fecha:   cliente.Fecha,
						Cliente: cliente.RazonSocial,
						Cuit:    cliente.Cuit,
					},
					Totales: utildtos.TotalesTemplate{
						Titulo:       titulo,
						TipoReporte:  tipo,
						Elemento:     "rendiciones",
						Cantidad:     cliente.CantOperaciones,
						TotalCobrado: cliente.TotalCobrado,
						TotalRendido: cliente.RendicionTotal,
					},
					AdjuntarEstado: true,
					Attachment: utildtos.Attachment{
						Name:        fmt.Sprintf("%s%s", nombreArchivo, tipo_archivo),
						ContentType: contentType,
						WithFile:    true,
					},
					TipoEmail: "reporte",
				}

				// Enviar email con el reporte a los contacto-reportes asociados al cliente
				erro = s.util.EnviarMailService(filtro)
				logs.Info(erro)
				if erro != nil {
					erro = fmt.Errorf("no se no pudo enviar comprobante de rendiciones mensual al cliente %v", cliente.Clientes)
					errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
						Archivo: filtro.Attachment.Name,
						Error:   fmt.Sprintf("servicio email: %v", erro),
					})
					logs.Error(erro.Error())
					log := entities.Log{
						Tipo:          entities.EnumLog("Error"),
						Funcionalidad: "EnviarMailService",
						Mensaje:       erro.Error(),
					}
					erro = s.util.CreateLogService(log)
					if erro != nil {
						logs.Error("error: al crear logs: " + erro.Error())
						// return erro
					}
				}

				// en el caso de no registrar error al enviar correo se guardan los datos del reporte
				if len(cliente.Reporte) == 1 {
					// Guardar en DB los datos del reporte
					reporteToSave := cliente.Reporte[0]
					reporteToSave.RutaFile = CreateRutaFile(request[0].Clientes, request[0].RutaArchivo, "pdf")

					erro = s.repository.SaveGuardarDatosReporte(reporteToSave)
					// si se pudo guardar el reporte rrm en DB se guarda el numero de reporte
					if erro == nil {
						//guardar numero
						numero_reporte_rrm = cliente.Reporte[0].Nro_reporte
					}

					if erro != nil {
						mensaje := fmt.Errorf("no se pudieron registrar datos del comprobante de rendiciones mensual enviado al cliente %+v", cliente.Clientes).Error()
						logs.Info(mensaje)
						log := entities.Log{
							Tipo:          entities.EnumLog("Error"),
							Funcionalidad: "SaveGuardarDatosReporte",
							Mensaje:       mensaje,
						}
						erro = s.util.CreateLogService(log)
						if erro != nil {
							logs.Error("error: al crear logs: " + erro.Error())
						}
					}
				}

				// una vez enviado el correo se elimina el archivo
				erro = s.commons.BorrarArchivo(ruta, fmt.Sprintf(("%s%s"), nombreArchivo, tipo_archivo))
				if erro != nil {
					logs.Error(erro.Error())
					log := entities.Log{
						Tipo:          entities.EnumLog("Error"),
						Funcionalidad: "BorrarArchivos",
						Mensaje:       erro.Error(),
					}
					erro = s.util.CreateLogService(log)
					if erro != nil {
						logs.Error("error: al crear logs: " + erro.Error())
						// return nil, erro
					}
				}
			}
		} // fin del else
	} // fin de for _, cliente := range request
	return
}

func (s *reportesService) MakeControlReportes(apikeys string, fechaControlar string, token string, runEndpoint util.RunEndpoint) (response []reportedtos.ResponseControlReporte, erro error) {
	fechaConsultar := fechaControlar

	var (
		typo = "GET"
		url  = config.URL_PASARELA_API + "/reporte/verificar-cobranzas"
	)

	mapConHeaders := map[string]string{
		"authorization": token,
	}

	mapConQuerys := map[string]string{
		"fechaConsultar": fechaConsultar,
		"apykeys":        apikeys,
	}

	var objt interface{}
	var respCobranzasCliente reportedtos.ResponseData

	responseEndpointCobranzasCliente, err := runEndpoint.RunEndpoint(typo, url, mapConHeaders, objt, mapConQuerys, false)
	if err != nil {
		erro = err
		return
	}

	jsonCobranzasCliente, err := json.Marshal(responseEndpointCobranzasCliente)
	if err != nil {
		erro = err
		return
	}

	// Deserializar la respuesta JSON en la estructura ResponseData
	err = json.Unmarshal(jsonCobranzasCliente, &respCobranzasCliente)
	if err != nil {
		fmt.Println("Error al deserializar JSON: ", err.Error())
	}

	response = respCobranzasCliente.Data

	return
}

func (s *reportesService) SendControlReportes(request []reportedtos.ResponseControlReporte) (erro error) {

	params := utildtos.RequestDatosMail{
		Email:     []string{"control.archivos.wee@telco.com.ar"},
		Asunto:    "Control Diario Cobranzas",
		Nombre:    "Prueba",
		Mensaje:   "Test",
		From:      "Wee.ar!",
		TipoEmail: "mixto",
		Template:  "reporte_control.html",
		Datos:     request,
	}

	erro = s.util.EnviarMailService(params)
	logs.Info(erro)
	if erro != nil {
		erro = fmt.Errorf("no se pudo enviar email")
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "EnviarMailService",
			Mensaje:       erro.Error(),
		}
		erro = s.util.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			// return erro
		}
	}

	return
}

func (s *reportesService) SendRetencionestxt(request reportedtos.RequestRetencionEmail) (erro error) {

	if len(request.Emails) == 0 {
		erro = errors.New("sin Emails destinos")
		return
	}

	for _, archivoTXT := range request.Archivos {

		var contentType string
		var asunto string
		var nombreArchivo string

		asunto = "Retencion TXT"
		nombreArchivo = archivoTXT

		contentType = "application/txt"

		var email = request.Emails
		filtro := utildtos.RequestDatosMail{
			Email:          email,
			Asunto:         asunto,
			From:           "Wee.ar!",
			Nombre:         "",
			Mensaje:        "En este email se adjunta uno de los reportes txt.",
			AdjuntarEstado: true,
			Attachment: utildtos.Attachment{
				Name:        fmt.Sprintf("%s", nombreArchivo),
				ContentType: contentType,
				WithFile:    true,
			},
			TipoEmail:   "mixto",
			RutaArchivo: config.DIR_BASE + "/documentos/retenciones/",
			Template:    "send_mail.html",
		}

		erro = s.util.EnviarMailService(filtro)
		logs.Info(erro)
		if erro != nil {
			erro = fmt.Errorf("no se no pudo enviar archivos txt de retenciones mensuales.")
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "EnviarMailService",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return erro
			}
		}

		// una vez enviado el correo se elimina el archivo
		erro = s.commons.BorrarArchivo(config.DIR_BASE+"/documentos/retenciones/", fmt.Sprintf(("%s"), nombreArchivo))
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "BorrarArchivos",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return nil, erro
			}
		}
	} // fin de for _, cliente := range request
	return
}

func (s *reportesService) GetPagoItems(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosItems, erro error) {
	// SE DEFINEN VARIABLES TOTALES
	var fechaI time.Time
	var fechaF time.Time
	if filtro.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}
	fecha := commons.ConvertFechaString(fechaI) // fecha de creacion del archivo

	for _, cliente := range request.Clientes {
		if cliente.ReporteBatch {
			filtro := reportedtos.RequestPagosPeriodo{
				ClienteId:   uint64(cliente.Id),
				FechaInicio: fechaI,
				FechaFin:    fechaF,
			}
			// logs.Info(filtro)
			// se debe obtener los pagos aptobados y autorizados(debin)
			// pagosItems, err := s.repository.GetPagosBatch(filtro)

			var listaPagos []reportedtos.DetallesPagosCobranza
			prisma, err := s.repository.GetCobranzasPrisma(filtro)
			apilink, err := s.repository.GetCobranzasApilink(filtro)
			rapipago, err := s.repository.GetCobranzasRapipago(filtro)
			multipago, err := s.repository.GetCobranzasMultipago(filtro)

			if err != nil {
				erro = err
				return
			}

			listaPagos = append(listaPagos, prisma...)
			listaPagos = append(listaPagos, apilink...)
			listaPagos = append(listaPagos, rapipago...)
			listaPagos = append(listaPagos, multipago...)

			// obtener lote del cliente
			lote, err := s.repository.GetLastLote(filtro)
			if err != nil {
				erro = err
				return
			}

			var idpg []uint
			// solo los pagos de tipo C y que no se informaron en algun lote
			if len(listaPagos) > 0 {

				// obteniendo los ids de los pagos
				for _, pago := range listaPagos {
					if pago.TotalPago > 0 && pago.Lote == 0 {
						idpg = append(idpg, uint(pago.Id))
					}
				}

				// obteniendo pago items
				pago_items, err := s.repository.GetPagosItems(idpg)

				if err != nil {
					erro = err
					return
				}

				// conectando pago items con pagos
				for i, pago := range listaPagos {
					for _, pago_item := range pago_items {
						if pago.Id == int(pago_item.PagosID) {
							listaPagos[i].PagoItems = append(listaPagos[i].PagoItems, pago_item)
						}
					}
				}

				// creando respuesta
				response = append(response, reportedtos.ResponsePagosItems{
					Clientes: reportedtos.ClientesResponse{
						Id:          cliente.Id,
						Cliente:     cliente.Cliente,
						RazonSocial: cliente.NombreFantasia,
						Email:       cliente.Email,
					},
					Fecha:          fecha,
					PagosCobranzas: listaPagos,
					PagLotes: reportedtos.PagLotes{
						Idpg:          idpg,
						Idcliente:     cliente.Id,
						Lote:          int(lote.Lote) + 1,
						Fechalote:     fecha,
						Cliente:       cliente.Cliente,
						NombreReporte: cliente.NombreReporte,
					},
				})
			}
		}
	}
	return
}

func (s *reportesService) GetPagoItemsAlternativo(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosItemsAlternativo, erro error) {
	// OBTENER ESTADOS DE LOS PAGOS:
	//aprobado (credito , debito y offline)
	paid, erro := s.util.FirstOrCreateConfiguracionService("PAID", "Nombre del estado aprobado", "Paid")
	if erro != nil {
		return
	}
	filtroPagosEstado := filtros.PagoEstadoFiltro{
		Nombre: paid,
	}
	estado_paid, err := s.administracion.GetPagoEstado(filtroPagosEstado)
	if err != nil {
		erro = err
		return
	}
	//si no se obtiene el estado del pago no se puede seguir
	if estado_paid[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	//autorizado (debin)
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, err := s.administracion.GetPagoEstado(filtroPagoEstado)

	if err != nil {
		erro = err
		return
	}

	//si no se obtiene el estado del pago no se puede seguir
	if pagoEstadoAcreditado[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID_AUTORIZADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID_AUTORIZADO,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	// SE DEFINEN VARIABLES TOTALES
	var pagoestados []uint
	var fechaI time.Time
	var fechaF time.Time
	if filtro.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}
	fecha := commons.ConvertFechaString(fechaI) // fecha de creacion del archivo
	pagoestados = append(pagoestados, estado_paid[0].ID, pagoEstadoAcreditado[0].ID)

	for _, cliente := range request.Clientes {
		// if cliente.ReporteBatch {
		filtro := reportedtos.RequestPagosPeriodo{
			ClienteId:   uint64(cliente.Id),
			FechaInicio: fechaI,
			FechaFin:    fechaF,
			PagoEstados: pagoestados,
		}
		// logs.Info(filtro)
		// se debe obtener los pagos aptobados y autorizados(debin)
		pagosItems, err := s.repository.GetPagosBatch(filtro)
		if err != nil {
			erro = err
			return
		}
		// obtener lote del cliente
		lote, err := s.repository.GetLastLote(filtro)
		if err != nil {
			erro = err
			return
		}

		var pagos []entities.Pago
		var idpg []uint

		// solo los pagos de tipo C y que no se informaron en algun lote
		if len(pagosItems) > 0 {
			for _, pg := range pagosItems {
				if pg.PagoIntentos[len(pg.PagoIntentos)-1].Amount > 0 {
					pagos = append(pagos, pg)
					idpg = append(idpg, pg.ID)
				}
			}
		}
		// respuesta solo si existen  pagos para ese cliente
		if len(pagos) > 0 {
			response = append(response, reportedtos.ResponsePagosItemsAlternativo{
				Clientes: reportedtos.ClientesResponse{
					Id:          cliente.Id,
					Cliente:     cliente.Cliente,
					RazonSocial: cliente.NombreFantasia,
					Email:       cliente.Email,
				},
				Fecha: fecha,
				Pagos: pagos,
				PagLotes: reportedtos.PagLotes{
					Idpg:          idpg,
					Idcliente:     cliente.Id,
					Lote:          int(lote.Lote) + 1,
					Fechalote:     fecha,
					Cliente:       cliente.Cliente,
					NombreReporte: cliente.NombreReporte,
				},
				// BatchContactos: cliente.BatchContactos,
			})
		}

		// }
	}
	return
}

func (s *reportesService) BuildPagosItems(request []reportedtos.ResponsePagosItems) (response []reportedtos.ResultPagosItems) {
	var cabeceraArchivo reportedtos.CabeceraArchivo
	var cabeceraLote reportedtos.CabeceraLote
	var colaArchivo reportedtos.ColaArchivo
	for _, pago := range request {
		// CABECERAS

		cabeceraArchivo = reportedtos.CabeceraArchivo{
			RecordCode:   "1",
			CreateDate:   pago.Fecha,
			OrigenName:   commons.EspaciosBlanco("WEE", 25, "RIGHT"),
			ClientNumber: commons.EspaciosBlanco("", 9, "RIGHT"),
			ClientName:   commons.EspaciosBlanco("", 35, "RIGHT"),
			Filler:       commons.EspaciosBlanco("", 54, "RIGHT"),
		}
		cabeceraLote = reportedtos.CabeceraLote{
			RecordCodeLote: "3",
			CreateDateLote: pago.Fecha,
			BatchNumber:    commons.AgregarCeros(6, pago.PagLotes.Lote), // esta longitud puede variar (su longitud maxima es 6)
			Description:    commons.EspaciosBlanco("", 35, "RIGHT"),
			Filler:         commons.EspaciosBlanco("", 82, "RIGHT"),
		}
		// DETALLES
		var resultItems []reportedtos.ResultItems
		var detalle1 reportedtos.DetalleTransaccion
		var detalle2 reportedtos.DetalleDescripcion
		var payment_date string
		var payment_time string
		var fileCount int64
		var totalFileAmount entities.Monto

		for _, items := range pago.PagosCobranzas {
			fechaPago := items.FechaPago

			if items.CanalPago == "DEBIN" || items.CanalPago == "OFFLINE" {
				fechaPago = items.FechaCobro
			}

			payment_date = commons.ConvertFechaString(fechaPago)
			payment_time = fmt.Sprintf("%v%v", fechaPago.Hour(), fechaPago.Minute())
			for _, pi := range items.PagoItems {
				fileCount = fileCount + 1
				totalFileAmount += pi.Amount
				detalle1 = reportedtos.DetalleTransaccion{
					RecordCodeTransaccion: "5",
					RecordSequence:        commons.AgregarCeros(5, 0),
					TransactionCode:       commons.AgregarCeros(2, 0),
					WorkDate:              commons.AgregarCeros(8, 0),
					TransferDate:          commons.AgregarCeros(8, 0),
					AccountNumber:         commons.EspaciosBlanco(pi.Description, 21, "RIGHT")[0:21], // pago items
					CurrencyCode:          commons.EspaciosBlanco("", 3, "RIGHT"),
					Amount:                commons.AgregarCeros(14, int(pi.Amount)), // pago items
					TerminalId:            commons.EspaciosBlanco("", 6, "RIGHT"),
					PaymentDate:           payment_date,                                        //Pago intento
					PaymentTime:           commons.AgregarCerosString(payment_time, 4, "LEFT"), // pago intento
					SeqNumber:             commons.AgregarCeros(4, 0),
					Filler:                commons.EspaciosBlanco("", 48, "RIGHT"),
				}
				detalle2 = reportedtos.DetalleDescripcion{
					RecordCodeLote: "6",
					BarCode:        commons.AgregarCerosString(pi.Identifier, 80, "LEFT")[0:80], // pago items
					TypeCode:       commons.EspaciosBlanco("", 1, "RIGHT"),
					Filler:         commons.EspaciosBlanco("", 50, "RIGHT"),
				}
				resultItems = append(resultItems, reportedtos.ResultItems{
					DetalleTransaccion: detalle1,
					DetalleDescripcion: detalle2,
				})
			}
		}
		// COLA DE ARCHIVO
		colaArchivo = reportedtos.ColaArchivo{
			RecordCodeCola:    "9",
			CreateDateCola:    pago.Fecha,
			TotalBatches:      commons.AgregarCeros(6, 0),
			FilePaymentCount:  commons.AgregarCeros(7, int(fileCount)),
			FilePaymentAmount: commons.AgregarCeros(12, int(totalFileAmount)), // total acumulado(detalles)
			Filler:            commons.AgregarCerosString("0", 38, "LEFT"),
			FileCount:         commons.AgregarCeros(7, 0),
			Filler2:           commons.EspaciosBlanco("", 53, "RIGHT"),
		}
		// RESPUESTA
		response = append(response, reportedtos.ResultPagosItems{
			PagLotes:        pago.PagLotes,
			CabeceraArchivo: cabeceraArchivo,
			CabeceraLote:    cabeceraLote,
			ResultItems:     resultItems,
			// DetalleTransaccion: detalle1,
			// DetalleDescripcion: detalle2,
			ColaArchivo: colaArchivo,
		})
	}
	return
}

func (s *reportesService) BuildPagosArchivo(request []reportedtos.ResponsePagosItemsAlternativo) (response []reportedtos.ResultPagosItemsAlternativo) {
	var cabeceraArchivo reportedtos.CabeceraArchivo
	var cabeceraLote reportedtos.CabeceraLote
	var colaArchivo reportedtos.ColaArchivo
	for _, pago := range request {
		// CABECERAS

		cabeceraArchivo = reportedtos.CabeceraArchivo{
			RecordCode:   "1",
			CreateDate:   pago.Fecha,
			OrigenName:   commons.EspaciosBlanco("WEE", 25, "RIGHT"),
			ClientNumber: commons.EspaciosBlanco("", 9, "RIGHT"),
			ClientName:   commons.EspaciosBlanco("", 35, "RIGHT"),
			Filler:       commons.EspaciosBlanco("", 54, "RIGHT"),
		}
		cabeceraLote = reportedtos.CabeceraLote{
			RecordCodeLote: "3",
			CreateDateLote: pago.Fecha,
			BatchNumber:    commons.AgregarCeros(6, pago.PagLotes.Lote), // esta longitud puede variar (su longitud maxima es 6)
			Description:    commons.EspaciosBlanco("", 35, "RIGHT"),
			Filler:         commons.EspaciosBlanco("", 82, "RIGHT"),
		}
		// DETALLES
		var resultItems []reportedtos.ResultItemsAlternativo
		var detalle1 reportedtos.DetalleTransaccion
		var detalle2 reportedtos.DetalleDescripcionAlternativo
		var payment_date string
		var payment_time string
		var fileCount int64
		var totalFileAmount entities.Monto
		for _, items := range pago.Pagos {
			payment_date = commons.ConvertFechaString(items.PagoIntentos[len(items.PagoIntentos)-1].PaidAt)
			payment_time = fmt.Sprintf("%v%v", items.PagoIntentos[len(items.PagoIntentos)-1].PaidAt.Hour(), items.PagoIntentos[len(items.PagoIntentos)-1].PaidAt.Minute())
			// for _, pi := range items.Pagoitems {
			fileCount = fileCount + 1
			totalFileAmount += items.PagoIntentos[len(items.PagoIntentos)-1].Amount
			detalle1 = reportedtos.DetalleTransaccion{
				RecordCodeTransaccion: "5",
				RecordSequence:        commons.AgregarCeros(5, 0),
				TransactionCode:       commons.AgregarCeros(2, 0),
				WorkDate:              commons.AgregarCeros(8, 0),
				TransferDate:          commons.AgregarCeros(8, 0),
				AccountNumber:         commons.EspaciosBlanco(items.ExternalReference, 21, "RIGHT")[0:21], // pago external reference
				CurrencyCode:          commons.EspaciosBlanco("", 3, "RIGHT"),
				Amount:                commons.AgregarCeros(14, int(items.PagoIntentos[len(items.PagoIntentos)-1].Amount)), // amount pago item
				TerminalId:            commons.EspaciosBlanco("", 6, "RIGHT"),
				PaymentDate:           payment_date,                                        //Pago intento
				PaymentTime:           commons.AgregarCerosString(payment_time, 4, "LEFT"), // pago intento
				SeqNumber:             commons.AgregarCeros(4, 0),
				Filler:                commons.EspaciosBlanco("", 48, "RIGHT"),
			}
			detalle2 = reportedtos.DetalleDescripcionAlternativo{
				RecordCodeLote: "6",
				UUID:           commons.AgregarCerosString(items.Uuid, 80, "LEFT")[0:80], // pago items
				TypeCode:       commons.EspaciosBlanco("", 1, "RIGHT"),
				Filler:         commons.EspaciosBlanco("", 50, "RIGHT"),
			}
			resultItems = append(resultItems, reportedtos.ResultItemsAlternativo{
				DetalleTransaccion: detalle1,
				DetalleDescripcion: detalle2,
			})
			// }
		}
		// COLA DE ARCHIVO
		colaArchivo = reportedtos.ColaArchivo{
			RecordCodeCola:    "9",
			CreateDateCola:    pago.Fecha,
			TotalBatches:      commons.AgregarCeros(6, 0),
			FilePaymentCount:  commons.AgregarCeros(7, int(fileCount)),
			FilePaymentAmount: commons.AgregarCeros(12, int(totalFileAmount)), // total acumulado(detalles)
			Filler:            commons.AgregarCerosString("0", 38, "LEFT"),
			FileCount:         commons.AgregarCeros(7, 0),
			Filler2:           commons.EspaciosBlanco("", 53, "RIGHT"),
		}
		// RESPUESTA
		response = append(response, reportedtos.ResultPagosItemsAlternativo{
			PagLotes:        pago.PagLotes,
			CabeceraArchivo: cabeceraArchivo,
			CabeceraLote:    cabeceraLote,
			ResultItems:     resultItems,
			// DetalleTransaccion: detalle1,
			// DetalleDescripcion: detalle2,
			ColaArchivo: colaArchivo,
			Clienteid:   pago.Clientes.Id,
			EmailsBatch: pago.BatchContactos,
		})
	}
	return
}

func (s *reportesService) ValidarEsctucturaPagosItems(request []reportedtos.ResultPagosItems) (err error) {
	var registroDescripcion reportedtos.EstructuraRegistrosBatch
	for _, items := range request {
		//validar datos de la cabecera
		err = validarRegistroCabeceraArchivo(items.CabeceraArchivo, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_CABECERA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}

		err = validarRegistroCabeceraLote(items.CabeceraLote, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_CABECERA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}

		err = validarRegistroDetalleTransaccion(items.ResultItems, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_DETALLE_TRANSACCION, err.Error())
			return errors.New(mensaje)
		}

		// err = validarRegistroDetalleDescripcion(items.DetalleDescripcion, registroDescripcion)
		// if err != nil {
		// 	mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_DETALLE_DESCRIPCION, err.Error())
		// 	return errors.New(mensaje)
		// }

		err = validarColaArchivo(items.ColaArchivo, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_COLA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}
	}
	return nil
}

func (s *reportesService) ValidarEsctucturaPagosBatch(request []reportedtos.ResultPagosItemsAlternativo) (err error) {
	var registroDescripcion reportedtos.EstructuraRegistrosBatch
	for _, items := range request {
		//validar datos de la cabecera
		err = validarRegistroCabeceraArchivo(items.CabeceraArchivo, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_CABECERA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}

		err = validarRegistroCabeceraLote(items.CabeceraLote, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_CABECERA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}

		err = validarRegistroDetalleTransaccionAlternativo(items.ResultItems, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_DETALLE_TRANSACCION, err.Error())
			return errors.New(mensaje)
		}

		// err = validarRegistroDetalleDescripcion(items.DetalleDescripcion, registroDescripcion)
		// if err != nil {
		// 	mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_DETALLE_DESCRIPCION, err.Error())
		// 	return errors.New(mensaje)
		// }

		err = validarColaArchivo(items.ColaArchivo, registroDescripcion)
		if err != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_FORMATO_REGISTRO_COLA_ARCHIVO, err.Error())
			return errors.New(mensaje)
		}
	}
	return nil
}

func validarRegistroCabeceraArchivo(cabeceraArchivo reportedtos.CabeceraArchivo, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {
	err := cabeceraArchivo.ValidarCabeceraArchivo(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func validarRegistroCabeceraLote(cabeceraLote reportedtos.CabeceraLote, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {
	err := cabeceraLote.ValidarCabeceraLote(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func validarRegistroDetalleTransaccion(descripcionTransaccion []reportedtos.ResultItems, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {

	for _, detalle := range descripcionTransaccion {
		if detalle.DetalleTransaccion.RecordCodeTransaccion == "5" {
			err := detalle.DetalleTransaccion.ValidarDetalleTransaccion(&registroDescripcion)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		if detalle.DetalleDescripcion.RecordCodeLote == "6" {
			err := detalle.DetalleDescripcion.ValidarDetalleDescripcion(&registroDescripcion)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

func validarRegistroDetalleTransaccionAlternativo(descripcionTransaccion []reportedtos.ResultItemsAlternativo, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {

	for _, detalle := range descripcionTransaccion {
		if detalle.DetalleTransaccion.RecordCodeTransaccion == "5" {
			err := detalle.DetalleTransaccion.ValidarDetalleTransaccion(&registroDescripcion)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		if detalle.DetalleDescripcion.RecordCodeLote == "6" {
			err := detalle.DetalleDescripcion.ValidarDetalleDescripcionAlternativo(&registroDescripcion)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

// func validarRegistroDetalleDescripcion(descripcionTransaccion []reportedtos.DetalleDescripcion, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {
// 	for _, detalle := range descripcionTransaccion {
// 		err := detalle.ValidarDetalleDescripcion(&registroDescripcion)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}
// 	}
// 	return nil
// }

func validarColaArchivo(cabeceraLote reportedtos.ColaArchivo, registroDescripcion reportedtos.EstructuraRegistrosBatch) error {
	err := cabeceraLote.ValidarColaArchivo(&registroDescripcion)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// ENVIAR ARCHIVO POR FTP ETC
func (s *reportesService) SendPagosItems(ctx context.Context, request []reportedtos.ResultPagosItems, filtro reportedtos.RequestPagosClientes) (erro error) {

	// obtener fecha de envio del archivo
	// por defecto toma el ultimo dia
	var fechaArchivo time.Time

	if filtro.FechaInicio.IsZero() {
		fechaArchivo = time.Now()
	} else {
		fechaArchivo = filtro.FechaInicio
	}

	/* en esta ruta se crearan los archivos para enviar */
	//ruta := fmt.Sprintf("..%s", config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}
	for _, pago := range request {
		// 1 CREAR EL NOMBRE DEL ARCHIVO
		// se crea con el nombre(WEE) + la fecha de hoy
		// fechaArchivo := time.Now()
		nombreArchivo := "WEE" + fechaArchivo.Format("020106")
		nombreArchivoPagosItems := commonsdtos.FileName{
			RutaBase:  ruta + "/",
			Nombre:    nombreArchivo,
			Extension: "txt",
			UsaFecha:  false,
		}
		rutaDetalle := s.commons.CreateFileName(nombreArchivoPagosItems)

		// 2 CREAR EL ARCHIVO
		_, erro = s.commons.CreateFile(rutaDetalle)
		if erro != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_CREAR_MOVIMIENTOS_LOTES, erro.Error())
			return errors.New(mensaje)
		}

		// 3 ABRIR EL ARCHIVO
		file, err := s.commons.LeerArchivo(rutaDetalle)
		if err != nil {
			err = erro
			return
		}

		// 4 ESCRIBIR EN EL ARCHIVO
		err = s.buildArchivo(file, pago)
		if err != nil {
			err = erro
			return
		}

		// 5 GUARDAR CAMBIOS y CERRAR ARCHIVO
		err = s.commons.GuardarCambios(file)
		if err != nil {
			err = erro
			return
		}

		//si el archivo se sube correctamente se registra en tabla movimientos lotes
		lotes := reportedtos.ToEntity(pago.PagLotes)
		erro = s.repository.SavePagosLotes(ctx, lotes)
		if erro != nil {
			mensaje := fmt.Errorf("no se pudieron registrar los siguiente movimientos en la tabla lotes %+v", pago.PagLotes.Idpg).Error()
			logs.Info(mensaje)
			return
		}

		// 6 ENVIAR ARCHIVO POR SFTP
		erro = s.SubirArchivo(ctx, nombreArchivoPagosItems, pago.PagLotes.NombreReporte, file)

		// 7 en el caso de que no se pueda enviar el archivo se deben dar de bajas los movimientos lotes creados
		if erro != nil {
			logs.Info("ocurrio error al enviar el archivo " + fmt.Sprintf("%v", erro))
			// 8.1 - En caso de que me tire un error se dan de bajas los movimientos lotes creados anteriormente
			err = s.repository.BajaPagosLotes(ctx, lotes, erro.Error())

			if err != nil {
				// 8.1.1 - En caso de que no se puede cancelar los movimientos aviso al usuario para que intervenga manualmente.
				mensaje_baja := fmt.Errorf("no se pudieron dar de bajas los siguientes movientos lotes %+v", pago.PagLotes.Idpg).Error()

				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionTransferencia,
					Descripcion: fmt.Sprintf("atención los siguientes movimientos de lotes no pudieron ser cancelados, movimientosId: %s", mensaje_baja),
				}
				erro = s.util.CreateNotificacionService(notificacion)
				if erro != nil {
					logs.Error(erro.Error())
				}
				erro = err
				return erro
			}
			return
		}

		// 5 UNA VEZ ENVIADO EL ARCHIVO , ELIMINAR EL ARCHIVO CREADO TEMPORALEMTE
		erro = s.commons.BorrarArchivo(ruta, fmt.Sprintf("%s.txt", nombreArchivo))
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "BorrarArchivos",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				return erro
			}
		}

	}

	// 6 BORRAR DIRECTORIO CREADO PARA EL REPORTE
	erro = s.commons.BorrarDirectorio(ruta)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.util.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			// return erro
		}
	}

	return
}

// ENVIAR ARCHIVO POR FTP BATCHS ALTERNATIVOS
func (s *reportesService) SendPagosBatch(ctx context.Context, request []reportedtos.ResultPagosItemsAlternativo, filtro reportedtos.RequestPagosClientes) (erro error) {

	// obtener fecha de envio del archivo
	// por defecto toma el ultimo dia
	var fechaArchivo time.Time

	if filtro.FechaInicio.IsZero() {
		fechaArchivo = time.Now()
	} else {
		fechaArchivo = filtro.FechaInicio
	}

	/* en esta ruta se crearan los archivos para enviar */
	//ruta := fmt.Sprintf("..%s", config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}
	for _, pago := range request {
		// 1 CREAR EL NOMBRE DEL ARCHIVO Y EMAILS
		// se crea con el nombre(WEE) + la fecha de hoy
		// fechaArchivo := time.Now()

		var email []string
		var nombreArchivo string

		clienteId := pago.Clienteid

		switch clienteId {
		case 18: // Goya
			email = []string{"mg.infor608@gmail.com", "SSINGRESOSPUBLICOS@GMAIL.COM"}
			// email = []string{"sebastian.escobar@telco.com.ar"}
			nombreArchivo = "WEE" + "GOYA" + fechaArchivo.Format("020106")
		case 19: // ACOR (Muni Corrientes)
			email = []string{"recaudacion@acor.gob.ar"}
			// email = []string{"sebastian.escobar@telco.com.ar"}
			nombreArchivo = "WEE" + "CORR" + fechaArchivo.Format("020106")
		}

		// email = pago.EmailsBatch

		nombreArchivoPagosItems := commonsdtos.FileName{
			RutaBase:  ruta + "/",
			Nombre:    nombreArchivo,
			Extension: "txt",
			UsaFecha:  false,
		}
		rutaDetalle := s.commons.CreateFileName(nombreArchivoPagosItems)

		// 2 CREAR EL ARCHIVO
		_, erro = s.commons.CreateFile(rutaDetalle)
		if erro != nil {
			mensaje := fmt.Sprintf("%v: %v", ERROR_CREAR_MOVIMIENTOS_LOTES, erro.Error())
			return errors.New(mensaje)
		}

		// 3 ABRIR EL ARCHIVO
		file, err := s.commons.LeerArchivo(rutaDetalle)
		if err != nil {
			err = erro
			return
		}

		// 4 ESCRIBIR EN EL ARCHIVO
		err = s.buildArchivoAlternativo(file, pago)
		if err != nil {
			err = erro
			return
		}

		// 5 GUARDAR CAMBIOS y CERRAR ARCHIVO
		err = s.commons.GuardarCambios(file)
		if err != nil {
			err = erro
			return
		}

		//si el archivo se sube correctamente se registra en tabla movimientos lotes
		// lotes := reportedtos.ToEntity(pago.PagLotes)
		// erro = s.repository.SavePagosLotes(ctx, lotes)
		// if erro != nil {
		// 	mensaje := fmt.Errorf("no se pudieron registrar los siguiente movimientos en la tabla lotes %+v", pago.PagLotes.Idpg).Error()
		// 	logs.Info(mensaje)
		// 	return
		// }

		// // 6 ENVIAR ARCHIVO

		asunto := "Archivo de rendicion de pagos"

		contentType := "txt"
		nombreArchivo = nombreArchivo + "." + contentType

		filtro := utildtos.RequestDatosMail{
			Email:          email,
			Asunto:         asunto,
			From:           "Wee.ar!",
			Nombre:         "",
			Mensaje:        "En este email se adjunta el archivo batch de cobranzas",
			AdjuntarEstado: true,
			Attachment: utildtos.Attachment{
				Name:        nombreArchivo,
				ContentType: contentType,
				WithFile:    true,
			},
			TipoEmail:   "mixto",
			RutaArchivo: ruta + "/",
			Template:    "send_mail.html",
		}

		erro = s.util.EnviarMailService(filtro)
		if erro != nil {
			erro = fmt.Errorf("no se pudo enviar archivo batch de rendicion de pagos")
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "EnviarMailService",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return erro
			}
		} else {
			// caso de exito
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionArchivoBatchCliente,
				Descripcion: fmt.Sprintf("Enviado archivo %v", nombreArchivo),
			}
			s.util.CreateNotificacionService(notificacion)
		}

		// erro = s.SubirArchivo(ctx, nombreArchivoPagosItems, pago.PagLotes.Cliente, file)

		// 5 UNA VEZ ENVIADO EL ARCHIVO , ELIMINAR EL ARCHIVO CREADO TEMPORALEMTE
		erro = s.commons.BorrarArchivo(ruta, nombreArchivo)
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "BorrarArchivos",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				return erro
			}
		}

	}

	// 6 BORRAR DIRECTORIO CREADO PARA EL REPORTE
	erro = s.commons.BorrarDirectorio(ruta)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.util.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			// return erro
		}
	}

	return
}

func (s *reportesService) buildArchivo(archivo *os.File, request reportedtos.ResultPagosItems) (erro error) {

	/* 	CABECERA DE ARCHIVO */
	cabArchivo := []string{request.CabeceraArchivo.RecordCode, request.CabeceraArchivo.CreateDate,
		request.CabeceraArchivo.OrigenName, request.CabeceraArchivo.ClientNumber,
		request.CabeceraArchivo.ClientName, request.CabeceraArchivo.Filler, "\n"}
	resultcabeceraArhivo := commons.JoinString(cabArchivo)

	// escribir cabecera archivo
	erro = s.commons.EscribirArchivo(resultcabeceraArhivo, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE ARCHIVO */

	/* CABECERA DE LOTE */
	cabLote := []string{request.CabeceraLote.RecordCodeLote, request.CabeceraLote.CreateDateLote,
		request.CabeceraLote.BatchNumber, request.CabeceraLote.Description, request.CabeceraLote.Filler, "\n"}
	resultcabeceraLote := commons.JoinString(cabLote)
	// escribir cabecera lote
	erro = s.commons.EscribirArchivo(resultcabeceraLote, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE LOTE */

	/* DETALLES TRANSACCION */
	for _, detalle := range request.ResultItems {
		var detalleTransaccion = []string{}
		var detalleDescripcion = []string{}
		detalleTransaccion = []string{detalle.DetalleTransaccion.RecordCodeTransaccion, detalle.DetalleTransaccion.RecordSequence, detalle.DetalleTransaccion.TransactionCode,
			detalle.DetalleTransaccion.WorkDate, detalle.DetalleTransaccion.TransferDate, detalle.DetalleTransaccion.AccountNumber, detalle.DetalleTransaccion.CurrencyCode, detalle.DetalleTransaccion.Amount,
			detalle.DetalleTransaccion.TerminalId, detalle.DetalleTransaccion.PaymentDate, detalle.DetalleTransaccion.PaymentTime, detalle.DetalleTransaccion.SeqNumber, detalle.DetalleTransaccion.Filler, "\n"}
		resultdetalleTransaccion := commons.JoinString(detalleTransaccion)
		// escribir cabecera lote
		erro = s.commons.EscribirArchivo(resultdetalleTransaccion, archivo)
		if erro != nil {
			return erro
		}

		detalleDescripcion = []string{detalle.DetalleDescripcion.RecordCodeLote, detalle.DetalleDescripcion.BarCode, detalle.DetalleDescripcion.TypeCode,
			detalle.DetalleDescripcion.Filler, "\n"}
		resultdetalleDescripcion := commons.JoinString(detalleDescripcion)
		// escribir cabecera lote
		erro = s.commons.EscribirArchivo(resultdetalleDescripcion, archivo)
		if erro != nil {
			return erro
		}
		/* 	END DETALLES TRANSACCION */
	}

	// /* DETALLES DESCRIPCION */
	// for _, detalle2 := range request.DetalleDescripcion {
	// 	var detalleDescripcion = []string{}
	// 	detalleDescripcion = []string{detalle2.RecordCodeLote, detalle2.BarCode, detalle2.TypeCode,
	// 		detalle2.Filler, "\n"}
	// 	resultdetalleDescripcion := commons.JoinString(detalleDescripcion)
	// 	// escribir cabecera lote
	// 	erro = s.commons.EscribirArchivo(resultdetalleDescripcion, archivo)
	// 	if erro != nil {
	// 		return erro
	// 	}
	// 	/* 	END DETALLES DESCRIPCION */
	// }

	/* COLA DE ARCHIVO */
	colaArchivo := []string{request.ColaArchivo.RecordCodeCola, request.ColaArchivo.CreateDateCola,
		request.ColaArchivo.TotalBatches, request.ColaArchivo.FilePaymentCount, request.ColaArchivo.FilePaymentAmount,
		request.ColaArchivo.Filler, request.ColaArchivo.FileCount, request.ColaArchivo.Filler2}
	resultcolaArchivo := commons.JoinString(colaArchivo)
	// escribir cabecera lote
	erro = s.commons.EscribirArchivo(resultcolaArchivo, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE LOTE */

	return
}

func (s *reportesService) buildArchivoAlternativo(archivo *os.File, request reportedtos.ResultPagosItemsAlternativo) (erro error) {

	/* 	CABECERA DE ARCHIVO */
	cabArchivo := []string{request.CabeceraArchivo.RecordCode, request.CabeceraArchivo.CreateDate,
		request.CabeceraArchivo.OrigenName, request.CabeceraArchivo.ClientNumber,
		request.CabeceraArchivo.ClientName, request.CabeceraArchivo.Filler, "\r", "\n"}
	resultcabeceraArhivo := commons.JoinString(cabArchivo)

	// escribir cabecera archivo
	erro = s.commons.EscribirArchivo(resultcabeceraArhivo, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE ARCHIVO */

	/* CABECERA DE LOTE */
	cabLote := []string{request.CabeceraLote.RecordCodeLote, request.CabeceraLote.CreateDateLote,
		request.CabeceraLote.BatchNumber, request.CabeceraLote.Description, request.CabeceraLote.Filler, "\r", "\n"}
	resultcabeceraLote := commons.JoinString(cabLote)
	// escribir cabecera lote
	erro = s.commons.EscribirArchivo(resultcabeceraLote, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE LOTE */

	/* DETALLES TRANSACCION */
	for _, detalle := range request.ResultItems {
		var detalleTransaccion = []string{}
		var detalleDescripcion = []string{}
		detalleTransaccion = []string{detalle.DetalleTransaccion.RecordCodeTransaccion, detalle.DetalleTransaccion.RecordSequence, detalle.DetalleTransaccion.TransactionCode,
			detalle.DetalleTransaccion.WorkDate, detalle.DetalleTransaccion.TransferDate, detalle.DetalleTransaccion.AccountNumber, detalle.DetalleTransaccion.CurrencyCode, detalle.DetalleTransaccion.Amount,
			detalle.DetalleTransaccion.TerminalId, detalle.DetalleTransaccion.PaymentDate, detalle.DetalleTransaccion.PaymentTime, detalle.DetalleTransaccion.SeqNumber, detalle.DetalleTransaccion.Filler, "\r", "\n"}
		resultdetalleTransaccion := commons.JoinString(detalleTransaccion)
		// escribir cabecera lote
		erro = s.commons.EscribirArchivo(resultdetalleTransaccion, archivo)
		if erro != nil {
			return erro
		}

		detalleDescripcion = []string{detalle.DetalleDescripcion.RecordCodeLote, detalle.DetalleDescripcion.UUID, detalle.DetalleDescripcion.TypeCode,
			detalle.DetalleDescripcion.Filler, "\r", "\n"}
		resultdetalleDescripcion := commons.JoinString(detalleDescripcion)
		// escribir cabecera lote
		erro = s.commons.EscribirArchivo(resultdetalleDescripcion, archivo)
		if erro != nil {
			return erro
		}
		/* 	END DETALLES TRANSACCION */
	}

	// /* DETALLES DESCRIPCION */
	// for _, detalle2 := range request.DetalleDescripcion {
	// 	var detalleDescripcion = []string{}
	// 	detalleDescripcion = []string{detalle2.RecordCodeLote, detalle2.BarCode, detalle2.TypeCode,
	// 		detalle2.Filler, "\n"}
	// 	resultdetalleDescripcion := commons.JoinString(detalleDescripcion)
	// 	// escribir cabecera lote
	// 	erro = s.commons.EscribirArchivo(resultdetalleDescripcion, archivo)
	// 	if erro != nil {
	// 		return erro
	// 	}
	// 	/* 	END DETALLES DESCRIPCION */
	// }

	/* COLA DE ARCHIVO */
	colaArchivo := []string{request.ColaArchivo.RecordCodeCola, request.ColaArchivo.CreateDateCola,
		request.ColaArchivo.TotalBatches, request.ColaArchivo.FilePaymentCount, request.ColaArchivo.FilePaymentAmount,
		request.ColaArchivo.Filler, request.ColaArchivo.FileCount, request.ColaArchivo.Filler2}
	resultcolaArchivo := commons.JoinString(colaArchivo)
	// escribir cabecera lote
	erro = s.commons.EscribirArchivo(resultcolaArchivo, archivo)
	if erro != nil {
		return erro
	}
	/* 	END CABECERA DE LOTE */

	return
}

func (s *reportesService) SubirArchivo(ctx context.Context, rutaArchivos commonsdtos.FileName, cliente string, archivo *os.File) (erro error) {
	// rutaDestino := config.DIR_KEY_REPORTES
	rutaDestinoReporte := strings.Replace(config.DIR_KEY_REPORTES, "*", cliente, 3)
	data, filename, filetypo, err := util.LeerArchivo(rutaDestinoReporte, rutaArchivos.RutaBase, rutaArchivos.Nombre+"."+rutaArchivos.Extension)
	if err != nil {
		erro = err
		logs.Error(err)
		return
	}

	erro = s.store.PutObject(ctx, data, filename, filetypo)
	if erro != nil {
		logs.Error("No se pudo guardar el archivo")
		return
	}

	defer archivo.Close()
	return
}

func _setPaginacion(number uint32, size uint32, total int64) (meta dtos.Meta) {
	from := (number - 1) * size
	lastPage := math.Ceil(float64(total) / float64(size))
	meta = dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}
	return
}

func (s *reportesService) GetCobranzas(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseCobranzas, erro error) {
	err := request.Validar()
	if err != nil {
		return response, err
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fecha := s.commons.ConvertirFecha(request.Date)

	// parse date
	date, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return response, errors.New(ERROR_CONVERTIR_FECHA)
	}

	// obtener estado pendiente para filtrar pagos
	filtroEstadoPago := filtros.PagoEstadoFiltro{
		Nombre: "pending",
	}
	estadoPendiente, err := s.administracion.GetPagoEstado(filtroEstadoPago)
	if err != nil {
		erro = err
		return
	}

	//buscar transferencias los pagos correspondiente al cliente(apiKey)
	filtro := reportedtos.RequestPagosPeriodo{
		ApiKey:      apikey,
		FechaInicio: date,
		FechaFin:    date,
	}

	// 3 obtener pagos del periodo
	listapagos, erro := s.repository.GetPagosReportes(filtro, estadoPendiente[0].ID)
	if erro != nil {
		return
	}

	var totalCobrado entities.Monto
	var descuentoComisionIva entities.Monto
	var totalNeto entities.Monto

	if len(listapagos) > 0 {
		var resulRendiciones []reportedtos.ResponseDetalleCobranza
		for _, m := range listapagos {
			// controlar que este pago sea un movimiento
			var comision entities.Monto
			var iva entities.Monto
			var availableAt string
			var netAmount entities.Monto
			if len(m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos) > 0 {
				if m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Tipo == "C" {
					netAmount = m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Monto
					availableAt = m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].CreatedAt.Format("2006-01-02 15:04:05")
					totalCobrado += m.PagoIntentos[len(m.PagoIntentos)-1].Amount
					totalNeto += m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Monto
					if len(m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientocomisions) > 0 && len(m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos) > 0 {
						comision = m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientocomisions[0].Monto + m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientocomisions[0].Montoproveedor
						iva = m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos[0].Monto + m.PagoIntentos[len(m.PagoIntentos)-1].Movimientos[0].Movimientoimpuestos[0].Montoproveedor
					}
					descuentoComisionIva += comision + iva
				}
			}
			// netAmount := m.Monto - (comision + iva)
			resulRendiciones = append(resulRendiciones, reportedtos.ResponseDetalleCobranza{
				InformedDate: date,
				// NOTE informacion del pago
				RequestId:         int64(m.ID),
				ExternalReference: m.ExternalReference,
				PayerName:         m.PayerName,
				Description:       m.Description,
				PaymentDate:       m.PagoIntentos[len(m.PagoIntentos)-1].PaidAt,
				Channel:           m.PagoIntentos[len(m.PagoIntentos)-1].Mediopagos.Mediopago,
				AmountPaid:        s.util.ToFixed(m.PagoIntentos[len(m.PagoIntentos)-1].Amount.Float64(), 2),
				// NOTE disponible solo cuando hay movimientos
				NetFee:      s.util.ToFixed(comision.Float64(), 2),
				IvaFee:      s.util.ToFixed(iva.Float64(), 2),
				NetAmount:   s.util.ToFixed(netAmount.Float64(), 2),
				AvailableAt: availableAt,
			})

			response = reportedtos.ResponseCobranzas{
				AccountId:      commons.AgregarCerosString(fmt.Sprintf("%v", listapagos[0].PagosTipo.CuentasID), 6, "LEFT"),
				ReportDate:     date,
				TotalCollected: s.util.ToFixed(totalCobrado.Float64(), 2),
				TotalGrossFee:  s.util.ToFixed(descuentoComisionIva.Float64(), 2),
				TotalNetAmount: s.util.ToFixed(totalNeto.Float64(), 2),
				Data:           resulRendiciones,
			}
		}
	}

	return
}

func (s *reportesService) GetRendiciones(request reportedtos.RequestCobranzas, apikey string) (response reportedtos.ResponseRendiciones, erro error) {
	err := request.Validar()
	if err != nil {
		return response, err
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fecha := s.commons.ConvertirFecha(request.Date)

	// parse date
	date, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return response, errors.New(ERROR_CONVERTIR_FECHA)
	}

	//buscar transferencias los pagos correspondiente al cliente(apiKey)
	filtro := reportedtos.RequestPagosPeriodo{
		ApiKey:      apikey,
		FechaInicio: date,
		FechaFin:    date,
	}

	// 3 obtener pagos del periodo
	// TODO se obtienen transferencias del cliente indicado en el filtro
	listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
	if err != nil {
		erro = err
		return
	}

	// FIXME esto deberia consultar la tabla movimientos
	// en la misma estan tanto los apgos acreditados como debitados  rquermiento para este reporte
	// filtroM := reportedtos.RequestPagosPeriodo{
	// 	ApiKey:                          apikey,
	// 	CargarComisionImpuesto:          true,
	// 	CargarMovimientosTransferencias: true,
	// 	CargarPagoIntentos:              true,
	// 	CargarCuenta:                    true,
	// }
	// listamovimientos, err := s.repository.GetMovimiento(filtroM)
	// if err != nil {
	// 	erro = err
	// 	return
	// }
	// logs.Info(listamovimientos)

	// totales credit(transferencias) response
	var totalCredit uint64
	var creditAmount entities.Monto

	// NOTE este es el total de sumar credit - debit
	var settlementAmount entities.Monto

	var pagosintentos []uint64
	var filtroMov reportedtos.RequestPagosPeriodo
	if len(listaTransferencia) > 0 {
		for _, transferencia := range listaTransferencia {
			pagosintentos = append(pagosintentos, transferencia.Movimiento.PagointentosId)
			totalCredit = totalCredit + 1
			creditAmount += transferencia.Movimiento.Monto * -1
			settlementAmount += transferencia.Movimiento.Monto * -1
		}
		filtroMov = reportedtos.RequestPagosPeriodo{
			PagoIntentos:                    pagosintentos,
			TipoMovimiento:                  "C",
			CargarComisionImpuesto:          true,
			CargarMovimientosTransferencias: true,
			CargarPagoIntentos:              true,
			CargarCuenta:                    true,
		}
	}

	// total debit si hubo reversiones de algun pago
	var totalDebit uint64
	var debitAmount entities.Monto
	if len(pagosintentos) > 0 {
		mov, err := s.repository.GetMovimiento(filtroMov)
		if err != nil {
			erro = err
			return
		}
		reversiones, err := s.repository.GetMovimiento(reportedtos.RequestPagosPeriodo{
			TipoMovimiento:                  "C",
			CargarComisionImpuesto:          true,
			CargarMovimientosTransferencias: true,
			CargarPagoIntentos:              true,
			CargarCuenta:                    true,
			CargarReversion:                 true,
			FechaInicio:                     date,
			FechaFin:                        date,
		})
		if err != nil {
			erro = err
			return
		}
		// en el caso de que existan reversiones ese dia se detalla en el resultado del reporte
		mov = append(mov, reversiones...)
		var resulRendiciones []reportedtos.ResponseDetalleRendiciones
		for _, m := range mov {
			var resultDebit float64
			resultCredit := s.util.ToFixed(m.Monto.Float64(), 2)
			if m.Reversion {
				totalDebit = totalDebit + 1
				debitAmount += m.Monto
				resultDebit = s.util.ToFixed(m.Monto.Float64(), 2)
				resultCredit = 0
				settlementAmount -= m.Monto
			}
			resulRendiciones = append(resulRendiciones, reportedtos.ResponseDetalleRendiciones{
				RequestId:         m.Pagointentos.PagosID,
				ExternalReference: m.Pagointentos.Pago.ExternalReference,
				Credit:            resultCredit,
				Debit:             resultDebit,
			})
		}

		response = reportedtos.ResponseRendiciones{
			AccountId:        commons.AgregarCerosString(fmt.Sprintf("%v", mov[0].CuentasId), 6, "LEFT"),
			ReportDate:       date,
			TotalCredits:     totalCredit,
			CreditAmount:     s.util.ToFixed(creditAmount.Float64(), 2),
			TotalDebits:      totalDebit,
			DebitAmount:      s.util.ToFixed(debitAmount.Float64(), 2),
			SettlementAmount: s.util.ToFixed(settlementAmount.Float64(), 2),
			Data:             resulRendiciones,
		}

	}

	return
}

func (s *reportesService) GetRecaudacion(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponsePagosLiquidacion, erro error) {
	var fechaI time.Time       // este seria la fecha de cobro
	var fechaProceso time.Time // fecha que se procesa el archivo
	var fechaF time.Time
	var lote int64
	var cantpagos int
	if filtro.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaProceso = fechaI
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaProceso = filtro.FechaInicio
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}
	for _, cliente := range request.Clientes {
		// este reporte es enviado solo a dpec
		if cliente.ReporteBatch {
			request := filtros.PagoIntentoFiltros{
				PagoEstadosIds:              []uint64{4, 7},
				CargarPago:                  true,
				CargarPagoTipo:              true,
				CargarPagoEstado:            true,
				CargarCuenta:                true,
				PagoIntentoAprobado:         true,
				CargarPagoCalculado:         true,
				FechaPagoInicio:             fechaI,
				FechaPagoFin:                fechaF,
				ClienteId:                   uint64(cliente.Id),
				CargarPagoItems:             true,
				CargarMovimientosTemporales: true,
				Channel:                     true,
			}

			//NOTE consultar los pagos intentos del dia anterior:
			// 1 Deben estar aprobados y el estado debe ser calculado
			// se deben comparar con los lotes informados en el archivo batch
			// se deben informar la misma cantidad en los 2 reportes cobranzas y liquidacion
			pagos, err := s.administracion.GetPagosIntentosCalculoComisionRepository(request)
			if err != nil {
				erro = err
				return
			}

			if len(pagos) > 0 {
				cantpagos = len(pagos)
				var pagos_id []uint64
				for _, pg := range pagos {
					pagos_id = append(pagos_id, uint64(pg.PagosID))
				}

				// obtener cantidad de lotes del cliente del dia
				//NOTE se debe verificar que las liquidaciones sean iguales a las cobranzas informadas
				filtro := reportedtos.RequestPagosPeriodo{
					ClienteId: uint64(cliente.Id),
					Pagos:     pagos_id,
				}
				lote, err = s.repository.GetCantidadLotes(filtro)
				if err != nil {
					erro = err
					return
				}
			}

			// respuesta solo si existen  pagos para ese cliente
			if int64(cantpagos) == lote {
				response = append(response, reportedtos.ResponsePagosLiquidacion{
					Clientes: reportedtos.Clientes{
						Id:          cliente.Id,
						Cliente:     cliente.Cliente,
						RazonSocial: cliente.RazonSocial,
						Email:       cliente.Emails,
						Cuit:        cliente.Cuit,
					},
					FechaCobro:   fechaI,
					FechaProceso: fechaProceso,
					Pagos:        pagos,
				})
			} else {
				erro = errors.New("error cantidad de pagos a informar es distinto a los lotes enviados en cobranzas")
				return
			}
		}
	}
	return
}

func (s *reportesService) BuildPagosLiquidacion(request []reportedtos.ResponsePagosLiquidacion) (response []reportedtos.ResultPagosLiquidacion) {
	var cabecera reportedtos.Clientes
	var fechaCobro string
	var fechaProceso string
	for _, pago := range request {

		// CABECERAS
		cabecera = pago.Clientes
		fechaCobro = pago.FechaCobro.Format("02-01-2006")
		fechaProceso = pago.FechaProceso.Format("02-01-2006")

		// DETALLES
		// var resultItemsMediopago []reportedtos.MedioPagoItems
		//  ? 1cobrado total y por medio de pago
		// credito
		var importecobrado entities.Monto
		var importecobradoCredito entities.Monto
		var importecobradoDebito entities.Monto
		var importecobradoDebin entities.Monto

		// debito
		var importeADepositar entities.Monto
		var importeADepositarCredito entities.Monto
		var importeADepositarDebito entities.Monto
		var importeADepositarDebin entities.Monto

		// cantidad boletas
		var cantidadTotalBoletas int
		var cantidadTotalBoletasCredit int
		var cantidadTotalBoletasDebito int
		var cantidadTotalBoletasDebin int

		// comision
		var comisionTotal entities.Monto
		var comisionCredito entities.Monto
		var comisionDebito entities.Monto
		var comisionDebin entities.Monto
		// iva
		var ivaTotal entities.Monto
		var ivaCredito entities.Monto
		var ivaDebito entities.Monto
		var ivaDebin entities.Monto

		// detales torales parciales
		var mcredit reportedtos.MedioPagoCredit
		var mdebito reportedtos.MedioPagoDebit
		var mdebin reportedtos.MedioPagoDebin
		for _, items := range pago.Pagos {
			importecobrado += items.Amount
			importeADepositar += items.Movimientotemporale[len(items.Movimientotemporale)-1].Monto
			cantidadTotalBoletas += len(items.Pago.Pagoitems)
			comisionTotal += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions)-1].Monto
			ivaTotal += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos)-1].Monto

			if items.Mediopagos.Channel.ID == 1 {
				importecobradoCredito += items.Amount
				importeADepositarCredito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Monto
				cantidadTotalBoletasCredit += len(items.Pago.Pagoitems)
				comisionCredito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions)-1].Monto
				ivaCredito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos)-1].Monto

				mcredit = reportedtos.MedioPagoCredit{
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoCredito.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarCredito.Float64(), 2))),
					CantidadBoletas:   fmt.Sprintf("%v", cantidadTotalBoletasCredit),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionCredito.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaCredito.Float64(), 2))),
				}
			}

			if items.Mediopagos.Channel.ID == 2 {
				importecobradoDebito += items.Amount
				importeADepositarDebito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Monto
				cantidadTotalBoletasDebito += len(items.Pago.Pagoitems)
				comisionDebito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions)-1].Monto
				ivaDebito += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos)-1].Monto

				mdebito = reportedtos.MedioPagoDebit{
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoDebito.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarDebito.Float64(), 2))),
					CantidadBoletas:   fmt.Sprintf("%v", cantidadTotalBoletasDebito),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionDebito.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaDebito.Float64(), 2))),
				}
			}

			if items.Mediopagos.Channel.ID == 4 {
				importecobradoDebin += items.Amount
				importeADepositarDebin += items.Movimientotemporale[len(items.Movimientotemporale)-1].Monto
				cantidadTotalBoletasDebin += len(items.Pago.Pagoitems)
				comisionDebin += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientocomisions)-1].Monto
				ivaDebin += items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos[len(items.Movimientotemporale[len(items.Movimientotemporale)-1].Movimientoimpuestos)-1].Monto

				mdebin = reportedtos.MedioPagoDebin{
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoDebin.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarDebin.Float64(), 2))),
					CantidadBoletas:   fmt.Sprintf("%v", cantidadTotalBoletasDebin),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionDebin.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaDebin.Float64(), 2))),
				}
			}
		}

		// RESPUESTA
		response = append(response, reportedtos.ResultPagosLiquidacion{
			Cabeceras:    cabecera,
			FechaCobro:   fechaCobro,
			FechaProceso: fechaProceso,
			MedioPagoItems: reportedtos.MedioPagoItems{
				MedioPagoCredit: mcredit,
				MedioPagoDebit:  mdebito,
				MedioPagoDebin:  mdebin,
			},
			Totales: reportedtos.TotalesALiquidar{
				ImporteCobrado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobrado.Float64(), 2))),
				ImporteADepositar:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositar.Float64(), 2))),
				CantidadTotalBoletas: fmt.Sprintf("%v", cantidadTotalBoletas),
				ComisionTotal:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionTotal.Float64(), 2))),
				IvaTotal:             fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaTotal.Float64(), 2))),
			},
		})
	}
	return
}

func (s *reportesService) SendLiquidacionClientes(request []reportedtos.ResultMovLiquidacion) (errorFile []reportedtos.ResponseCsvEmailError, erro error) {

	/* en esta ruta se crearan los archivos */
	ruta := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE) //dev
	// ruta := fmt.Sprintf(".%s", config.DIR_REPORTE) //prod
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	for _, cliente := range request {

		if len(cliente.Cabeceras.Email) == 0 {
			erro = fmt.Errorf("no esta definido el email del cliente %v", cliente.Cabeceras.Cliente)
			errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
				Archivo: "",
				Error:   fmt.Sprintf("error al enviar archivo: no esta definido email del cliente %v", cliente.Cabeceras.Cliente),
			})
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "EnviarMailService",
				Mensaje:       erro.Error(),
			}
			erro = s.util.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				// return erro
			}
		} else {

			logs.Info("Procesando reportes tipo: liquidacion diaria ")
			asunto := "Liquidación Wee! " + cliente.FechaProceso
			nombreArchivo := cliente.Cabeceras.Cliente + "-" + cliente.FechaProceso

			// crear en carpeta tempora

			erro = s.util.GetRecaudacionPdf(cliente, ruta, nombreArchivo)
			if erro != nil {
				return errorFile, erro
			}

			if cliente.Cabeceras.EnviarEmail {
				var campo_adicional = []string{"pagos"}
				var email = cliente.Cabeceras.Email //[]string{cliente.Email}
				filtro := utildtos.RequestDatosMail{
					Email:            email,
					Asunto:           asunto,
					From:             "Wee.ar!",
					Nombre:           "Wee.ar!",
					Mensaje:          "reportes de pagos: #0",
					CamposReemplazar: campo_adicional,
					AdjuntarEstado:   true,
					Attachment: utildtos.Attachment{
						Name:        fmt.Sprintf("%s.pdf", nombreArchivo),
						ContentType: "text/csv",
						WithFile:    true,
					},
					TipoEmail: "adjunto",
				}
				/*enviar archivo csv por correo*/
				erro = s.util.EnviarMailService(filtro)
				logs.Info(erro)
				if erro != nil {
					erro = fmt.Errorf("no se no pudo enviar rendicion al %v", cliente.Cabeceras.Cliente)
					errorFile = append(errorFile, reportedtos.ResponseCsvEmailError{
						Archivo: filtro.Attachment.Name,
						Error:   fmt.Sprintf("servicio email: %v", erro),
					})
					logs.Error(erro.Error())
					log := entities.Log{
						Tipo:          entities.EnumLog("Error"),
						Funcionalidad: "EnviarMailService",
						Mensaje:       erro.Error(),
					}
					erro = s.util.CreateLogService(log)
					if erro != nil {
						logs.Error("error: al crear logs: " + erro.Error())
						// return erro
					}
					/* informar el error al enviar el emial pero se debe continuar enviando los siguientes archivos a otros clientes */
				}
			}

			// una vez enviado el correo se elimina el archivo csv
			erro = s.commons.BorrarArchivo(ruta, fmt.Sprintf("%s.pdf", nombreArchivo))
			if erro != nil {
				logs.Error(erro.Error())
				log := entities.Log{
					Tipo:          entities.EnumLog("Error"),
					Funcionalidad: "BorrarArchivos",
					Mensaje:       erro.Error(),
				}
				erro = s.util.CreateLogService(log)
				if erro != nil {
					logs.Error("error: al crear logs: " + erro.Error())
					return nil, erro
				}
			}
		}

	}
	erro = s.commons.BorrarDirectorio(ruta)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.util.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			// return erro
		}
	}

	return
}

func (s *reportesService) GetRecaudacionDiaria(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseMovLiquidacion, erro error) {
	var fechaI time.Time // este seria la fecha de cobr
	var fechaF time.Time
	var fechaProceso time.Time // fecha que se procesa el archivo
	var fechaRendicion string
	if filtro.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaProceso = fechaI
	} else {
		fechaI = filtro.FechaInicio
		fechaF = filtro.FechaFin
		fechaProceso = filtro.FechaInicio
	}

	if filtro.FechaAdicional.IsZero() {
		fechaRendicion = fechaI.AddDate(0, 0, int(+2)).Format("02-01-2006")
	} else if filtro.CargarFechaAdicional {
		fechaRendicion = filtro.FechaAdicional.Format("02-01-2006")
	}

	for _, c := range request.Clientes {
		var orden_diaria bool
		for _, cu := range c.Cuenta {
			logs.Info("procesando orden de pago cliente" + c.Cliente)
			filtroMov := reportedtos.RequestPagosPeriodo{
				FechaInicio:                     fechaI,
				FechaFin:                        fechaF,
				TipoMovimiento:                  "C",
				CargarComisionImpuesto:          true,
				CargarMovimientosTransferencias: true,
				CargarPagoIntentos:              true,
				CargarCuenta:                    true,
				CargarMedioPago:                 true,
				CuentaId:                        uint64(cu.Id),
			}
			mov, err := s.repository.GetRendicionReportes(filtroMov)
			if err != nil {
				erro = err
				return
			}

			if !c.OrdenDiaria && len(mov) == 0 {
				orden_diaria = false
			}
			if c.OrdenDiaria {
				orden_diaria = true
			} else if !c.OrdenDiaria && len(mov) > 0 {
				orden_diaria = true
			}

			if orden_diaria {
				var detalleliquidacion []entities.Liquidaciondetalles
				if len(mov) > 0 {
					for _, m := range mov {
						detalleliquidacion = append(detalleliquidacion, entities.Liquidaciondetalles{
							PagointentosId: int64(m.PagointentosId),
							MovimientosId:  int64(m.ID),
							CuentasId:      int64(m.CuentasId),
						})
					}
				}

				liquidacion := entities.Movimientoliquidaciones{
					ClientesID:           uint64(c.Id),
					FechaEnvio:           fechaProceso.Format("2006-01-02"),
					LiquidacioneDetalles: detalleliquidacion,
				}
				// guardar y obtener numero de liquidacion(ultimo registro ingresado)
				nroliquidacion, err := s.repository.SaveLiquidacion(liquidacion)
				if err != nil {
					erro = err
					return
				}

				// buscar fecha de rendicion en transferencia (moviminetos tipo D)
				response = append(response, reportedtos.ResponseMovLiquidacion{
					Clientes: reportedtos.Cliente{
						Id:          c.Id,
						Cliente:     c.Cliente,
						RazonSocial: c.RazonSocial,
						Email:       c.Emails,
						Cuit:        c.Cuit,
						EnviarEmail: filtro.EnviarEmail,
						EnviarPdf:   false, // descomentar esta linea cuando se desea generar pdf de dpec
					},
					Cuenta:         cu.Cuenta,
					FechaRendicion: fechaRendicion,
					FechaProceso:   fechaProceso,
					Movimientos:    mov,
					NroLiquidacion: int(nroliquidacion),
				})
			}

		}
	}
	return
}

func (s *reportesService) BuildMovLiquidacion(request []reportedtos.ResponseMovLiquidacion) (response []reportedtos.ResultMovLiquidacion) {
	var cabecera reportedtos.Cliente
	var cuenta string
	var nroliquidacion int
	// var fechaCobro string
	var fechaProceso string
	var fecharendicion string
	for _, pago := range request {

		// CABECERAS
		cabecera = pago.Clientes
		cuenta = pago.Cuenta
		nroliquidacion = pago.NroLiquidacion
		fecharendicion = pago.FechaRendicion
		fechaProceso = pago.FechaProceso.Format("02-01-2006")

		// DETALLES
		// var resultItemsMediopago []reportedtos.MedioPagoItems
		//  ? 1cobrado total y por medio de pago
		// credito
		var importecobrado entities.Monto
		var importecobradoCredito entities.Monto
		var importecobradoDebito entities.Monto
		var importecobradoDebin entities.Monto

		// debito
		var importeADepositar entities.Monto
		var importeADepositarCredito entities.Monto
		var importeADepositarDebito entities.Monto
		var importeADepositarDebin entities.Monto

		// cantidad boletas
		var cantidadTotalBoletas int
		var cantidadTotalBoletasCredit int
		var cantidadTotalBoletasDebito int
		var cantidadTotalBoletasDebin int

		// comision
		var comisionTotal entities.Monto
		var comisionCredito entities.Monto
		var comisionDebito entities.Monto
		var comisionDebin entities.Monto
		// iva
		var ivaTotal entities.Monto
		var ivaCredito entities.Monto
		var ivaDebito entities.Monto
		var ivaDebin entities.Monto

		// detales torales parciales
		var mcredit reportedtos.MedioMovCredit
		var detallecredit []reportedtos.DetalleMov

		var mdebito reportedtos.MedioMovDebit
		var detalledebit []reportedtos.DetalleMov

		var mdebin reportedtos.MedioMovDebin
		var detalledebin []reportedtos.DetalleMov

		for _, items := range pago.Movimientos {
			importecobrado += items.Pagointentos.Amount
			importeADepositar += items.Monto
			cantidadTotalBoletas += len(items.Pagointentos.Pago.Pagoitems)
			comisionTotal += items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
			ivaTotal += items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor

			if items.Pagointentos.Mediopagos.ChannelsID == 1 {
				comision := items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				iva := items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor
				detallecredit = append(detallecredit, reportedtos.DetalleMov{
					Cuenta:            items.Cuenta.Cuenta,
					Referencia:        items.Pagointentos.Pago.ExternalReference,
					FechaCobro:        items.Pagointentos.PaidAt.Format("02-01-2006"),
					CantidadBoletas:   fmt.Sprintf("%v", len(items.Pagointentos.Pago.Pagoitems)),
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Pagointentos.Amount.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Monto.Float64(), 2))),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 2))),
				})
				importecobradoCredito += items.Pagointentos.Amount
				importeADepositarCredito += items.Monto
				cantidadTotalBoletasCredit += len(items.Pagointentos.Pago.Pagoitems)
				comisionCredito += items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				ivaCredito += items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor
			}

			if items.Pagointentos.Mediopagos.ChannelsID == 2 {
				comision1 := items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				iva1 := items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor
				detalledebit = append(detalledebit, reportedtos.DetalleMov{
					Cuenta:            items.Cuenta.Cuenta,
					Referencia:        items.Pagointentos.Pago.ExternalReference,
					FechaCobro:        items.Pagointentos.PaidAt.Format("02-01-2006"),
					CantidadBoletas:   fmt.Sprintf("%v", len(items.Pagointentos.Pago.Pagoitems)),
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Pagointentos.Amount.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Monto.Float64(), 2))),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision1.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva1.Float64(), 2))),
				})
				importecobradoDebito += items.Pagointentos.Amount
				importeADepositarDebito += items.Monto
				cantidadTotalBoletasDebito += len(items.Pagointentos.Pago.Pagoitems)
				comisionDebito += items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				ivaDebito += items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor

			}

			if items.Pagointentos.Mediopagos.ChannelsID == 4 {
				comision2 := items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				iva2 := items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor
				// cantidadBoletasDebin += len(items.Pagointentos.Pago.Pagoitems)
				detalledebin = append(detalledebin, reportedtos.DetalleMov{
					Cuenta:            items.Cuenta.Cuenta,
					Referencia:        items.Pagointentos.Pago.ExternalReference,
					FechaCobro:        items.Pagointentos.PaidAt.Format("02-01-2006"),
					CantidadBoletas:   fmt.Sprintf("%v", len(items.Pagointentos.Pago.Pagoitems)),
					ImporteCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Pagointentos.Amount.Float64(), 2))),
					ImporteADepositar: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(items.Monto.Float64(), 2))),
					Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision2.Float64(), 2))),
					Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva2.Float64(), 2))),
				})
				importecobradoDebin += items.Pagointentos.Amount
				importeADepositarDebin += items.Monto
				cantidadTotalBoletasDebin += len(items.Pagointentos.Pago.Pagoitems)
				comisionDebin += items.Movimientocomisions[len(items.Movimientocomisions)-1].Monto + items.Movimientocomisions[len(items.Movimientocomisions)-1].Montoproveedor
				ivaDebin += items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Monto + items.Movimientoimpuestos[len(items.Movimientoimpuestos)-1].Montoproveedor

			}
		}
		mcredit = reportedtos.MedioMovCredit{
			Detalle:              detallecredit,
			CantidaTotaldBoletas: fmt.Sprintf("%v", cantidadTotalBoletasCredit),
			TotalCobrado:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoCredito.Float64(), 2))),
			TotalaRendir:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarCredito.Float64(), 2))),
		}
		mdebito = reportedtos.MedioMovDebit{
			Detalle:              detalledebit,
			CantidaTotaldBoletas: fmt.Sprintf("%v", cantidadTotalBoletasDebito),
			TotalCobrado:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoDebito.Float64(), 2))),
			TotalaRendir:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarDebito.Float64(), 2))),
		}
		mdebin = reportedtos.MedioMovDebin{
			Detalle:              detalledebin,
			CantidaTotaldBoletas: fmt.Sprintf("%v", cantidadTotalBoletasDebin),
			TotalCobrado:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobradoDebin.Float64(), 2))),
			TotalaRendir:         fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositarDebin.Float64(), 2))),
		}

		// RESPUESTA
		response = append(response, reportedtos.ResultMovLiquidacion{
			Cabeceras:      cabecera,
			NroLiquidacion: commons.AgregarCerosString(fmt.Sprintf("%v", nroliquidacion), 6, "LEFT"),
			FechaProceso:   fechaProceso,
			Cuenta:         cuenta,
			FechaRendicion: fecharendicion,
			MedioPagoItems: reportedtos.MedioMovItems{
				MedioMovCredit: mcredit,
				MedioMovDebit:  mdebito,
				MedioMovDebin:  mdebin,
			},
			Totales: reportedtos.TotalesMovLiquidar{
				ImporteCobrado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importecobrado.Float64(), 2))),
				ImporteADepositar:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(importeADepositar.Float64(), 2))),
				CantidadTotalBoletas: fmt.Sprintf("%v", cantidadTotalBoletas),
				ComisionTotal:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisionTotal.Float64(), 2))),
				IvaTotal:             fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(ivaTotal.Float64(), 2))),
			},
		})
	}
	return
}

func (s *reportesService) NotificacionErroresReportes(errorFile []reportedtos.ResponseCsvEmailError) (erro error) {
	// var campos = []string{}
	// var slice_array = []string{}
	// for _, tr := range responseTransferencias {
	// 	slice_array = append(slice_array, tr.Cuenta, fmt.Sprintf("CuentaID%v ", tr.CuentaId), fmt.Sprintf("Origen%v ", tr.Origen), fmt.Sprintf("Destino%v ", tr.Destino), fmt.Sprintf("Importe%v ", tr.Importe), tr.Error, "\n")
	// }
	// mensaje := (strings.Join(slice_array, "-"))
	// var arrayEmail []string
	// // NOTE para  pruebas
	// // arrayEmail = append(arrayEmail, "jose.alarcon@telco.com.ar")
	// arrayEmail = append(arrayEmail, config.EMAIL_TELCO)
	// params := utildtos.RequestDatosMail{
	// 	Email:            arrayEmail,
	// 	Asunto:           "Error transferencias automaticas",
	// 	Nombre:           "Wee!!",
	// 	Mensaje:          mensaje,
	// 	CamposReemplazar: campos,
	// 	From:             "Wee.ar!",
	// 	TipoEmail:        "template",
	// }
	// erro = s.util.EnviarMailService(params)
	// if erro != nil {
	// 	logs.Info("Ocurrio un error al enviar correo notificación transferencias automaticas con error")
	// 	logs.Error(erro.Error())
	// }
	return
}

func transformarDatos(responseClienteReversion reportedtos.ResponseClientesReportes) (reversionesData []commons.ReversionData) {

	// Parsear manualmente las reversiones a la data struct para presentar el pdf
	for _, revers := range responseClienteReversion.Reversiones {
		var data commons.ReversionData
		// Datos del pago en header
		data.Pago.ReferenciaExterna = revers.PagoRevertido.ReferenciaExterna
		data.Pago.MedioPago = revers.MedioPago

		// formatear importe

		data.Pago.Monto = revers.PagoRevertido.IntentoPago.ImportePagado
		/* data.Pago.Monto = revers.PagoRevertido.IntentoPago.ImportePagado */
		data.Pago.IdPago = revers.PagoRevertido.IdPago
		data.Pago.Estado = revers.PagoRevertido.PagoEstado
		// Datos del intento de pago en herader
		data.Intento.IdIntento = revers.PagoRevertido.IntentoPago.IdIntentoPago
		data.Intento.IdTransaccion = revers.PagoRevertido.IntentoPago.IdTransaccion

		data.Intento.FechaPago = formatearStringAFechaCorrecta(revers.PagoRevertido.IntentoPago.FechaPago)
		data.Intento.Importe = revers.PagoRevertido.IntentoPago.ImportePagado

		for _, item := range revers.PagoRevertido.Items {
			var tempItem commons.ItemsReversionData
			tempItem.Cantidad = item.Cantidad
			tempItem.Descripcion = item.Descripcion
			tempItem.Identificador = item.Identificador

			montoitem_int64, _ := strconv.ParseInt(item.Monto, 10, 64)
			montoitem_float := entities.Monto(montoitem_int64).Float64()
			tempItem.Monto = util.Resolve().FormatNum(montoitem_float)
			tempItem.Monto = "$ " + tempItem.Monto
			data.Items = append(data.Items, tempItem)
		}
		data.Fecha = revers.Fecha
		reversionesData = append(reversionesData, data)
	}
	////
	return
}

// for _, revers := range responseClienteReversion.Reversiones {
// 	var data commons.ReversionData
// 	// Datos del pago en header
// 	data.Pago.ReferenciaExterna = revers.PagoRevertido.ReferenciaExterna
// 	data.Pago.MedioPago = revers.MedioPago
// 	data.Pago.Monto = revers.Monto
// 	data.Pago.IdPago = revers.PagoRevertido.IdPago
// 	data.Pago.Estado = revers.PagoRevertido.PagoEstado
// 	// Datos del intento de pago en herader
// 	data.Intento.IdIntento = revers.PagoRevertido.IntentoPago.IdIntentoPago
// 	data.Intento.IdTransaccion = revers.PagoRevertido.IntentoPago.IdTransaccion
// 	fecha := strings.Split(revers.PagoRevertido.IntentoPago.FechaPago, " ")
// 	data.Intento.FechaPago = fecha[0]
// 	data.Intento.Importe = revers.PagoRevertido.IntentoPago.ImportePagado

// 	for _, item := range revers.PagoRevertido.Items {
// 		var tempItem commons.ItemsReversionData
// 		tempItem.Cantidad = item.Cantidad
// 		tempItem.Descripcion = item.Descripcion
// 		tempItem.Identificador = item.Identificador
// 		tempItem.Monto = item.Monto
// 		data.Items = append(data.Items, tempItem)
// 	}

//		reversionesData = append(reversionesData, data)
//	}
func (s *reportesService) MovimientosComisionesService(request reportedtos.RequestReporteMovimientosComisiones) (res reportedtos.ResposeReporteMovimientosComisiones, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	fechaInicioString := s.commons.GetDateFirstMoment(fechaInicioTime)
	fechaFinString := s.commons.GetDateLastMoment(fechaFinTime)

	filtro := filtros_reportes.MovimientosComisionesFiltro{
		FechaInicio: fechaInicioString,
		FechaFin:    fechaFinString,
		ClienteId:   request.ClienteId,
		Number:      request.Number,
		Size:        request.Size,
	}

	resultado, total, err := s.repository.MovimientosComisionesRepository(filtro)
	if err != nil {
		erro = errors.New(ERROR_MOVIMIENTOS_COMISIONES)
		return
	}

	res.Reportes = total
	res.SetTotales()

	res.Reportes = resultado
	for i := 0; i < len(res.Reportes); i++ {
		var reporte = res.Reportes[i]
		res.Reportes[i].PorcentajeComision = s.util.ToFixed((res.Reportes[i].PorcentajeComision), 4)
		res.Reportes[i].Subtotal = (reporte.MontoComision + reporte.MontoImpuesto)
	}

	if request.Size != 0 {
		res.LastPage = int(math.Ceil(float64(len(total)) / float64(request.Size)))
	}

	return
}

func (s *reportesService) GetCobranzasClientesService(request reportedtos.RequestCobranzasClientes) (res reportedtos.ResponseCobranzasClientes, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	fechaInicioString := s.commons.GetDateFirstMoment(fechaInicioTime)
	fechaFinString := s.commons.GetDateLastMoment(fechaFinTime)

	filtro := filtros_reportes.CobranzasClienteFiltro{
		FechaInicio: fechaInicioString,
		FechaFin:    fechaFinString,
		ClienteId:   request.ClienteId,
	}

	resultado, err := s.repository.CobranzasClientesRepository(filtro)
	if err != nil {
		erro = errors.New(ERROR_COBRANZAS_CLIENTES)
		return
	}

	var fechas []string
	for _, pago := range resultado {

		// fecha, _ := time.Parse("2006-01-02T00:00:00Z", pago.FechaPago)
		fecha_pago := pago.FechaPago.Format("2006-01-02")
		if !contains(fechas, fecha_pago) {
			fechas = append(fechas, fecha_pago)
			var cobranza reportedtos.DetallesCobranza

			cobranza.Fecha = fecha_pago
			cobranza.Nombre = (pago.Cliente + "-" + pago.FechaPago.Format("02-01-2006"))
			for _, OtroPago := range resultado {
				// fechaOtro, _ := time.Parse("2006-01-02T00:00:00Z", OtroPago.FechaPago)
				fecha_pagoOtro := OtroPago.FechaPago.Format("2006-01-02")
				if fecha_pago == fecha_pagoOtro {
					cobranza.Pagos = append(cobranza.Pagos, OtroPago)
					cobranza.Registros += 1
					cobranza.Subtotal += OtroPago.TotalPago
				}
			}

			res.Cobranzas = append(res.Cobranzas, cobranza)

		}

	}

	res.CantidadCobranzas = len(res.Cobranzas)

	sort.Slice(res.Cobranzas, func(i, j int) bool {
		if res.Cobranzas[i].Fecha != "" && res.Cobranzas[j].Fecha != "" {
			return s.commons.ConvertirFecha(res.Cobranzas[i].Fecha) > s.commons.ConvertirFecha(res.Cobranzas[j].Fecha)
		}
		return false

	})

	for _, cob := range res.Cobranzas {
		res.Total += cob.Subtotal
	}

	return
}

func (s *reportesService) GetRendicionesClientesService(request reportedtos.RequestReporteClientes) (res reportedtos.ResponseRendicionesClientes, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	filtro := reportedtos.RequestPagosPeriodo{
		ClienteId:     uint64(request.ClienteId),
		CuentaId:      uint64(request.CuentaId),
		FechaInicio:   fechaInicioTime,
		FechaFin:      fechaFinTime,
		OrdenadoFecha: true,
	}

	// TODO se obtienen transferencias del cliente indicado en el filtro
	listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
	if err != nil {
		erro = err
		return
	}

	var pagosintentos []uint64

	type fechaDeposito struct {
		pagoIntentoId uint64
		fechaDeposito string
	}
	var depositos []fechaDeposito
	var depositosRev []fechaDeposito
	var filtroMov reportedtos.RequestPagosPeriodo
	var pagosintentosrevertidos []uint64
	var movrevertidos []entities.Movimiento
	var controlIds []string

	if len(listaTransferencia) > 0 {
		for _, transferencia := range listaTransferencia {
			transferenciaDepo := fechaDeposito{
				pagoIntentoId: transferencia.Movimiento.PagointentosId,
				fechaDeposito: transferencia.FechaOperacion.Format("02-01-2006"),
			}
			if !transferencia.Reversion {
				pagosintentos = append(pagosintentos, transferencia.Movimiento.PagointentosId)
				depositos = append(depositos, transferenciaDepo)
			}
			if transferencia.Reversion {
				pagosintentosrevertidos = append(pagosintentosrevertidos, transferencia.Movimiento.PagointentosId)
				depositosRev = append(depositosRev, transferenciaDepo)
			}
		}
		filtroMov = reportedtos.RequestPagosPeriodo{
			PagoIntentos:                    pagosintentos,
			TipoMovimiento:                  "C",
			CargarComisionImpuesto:          true,
			CargarMovimientosTransferencias: true,
			CargarPagoIntentos:              true,
			CargarCuenta:                    true,
			CargarCliente:                   true,
			OrdenadoFecha:                   true,
			CargarRetenciones:               true,
			CargarMovimientosRetenciones:    true,
		}
	}

	// se obtienen movimientos revertidos, en el caso de que existieran reversiones
	if len(pagosintentosrevertidos) > 0 {
		filtroRevertidos := reportedtos.RequestPagosPeriodo{
			PagoIntentos:                    pagosintentosrevertidos,
			TipoMovimiento:                  "C",
			CargarMovimientosTransferencias: true,
			CargarPagoIntentos:              true,
			CargarCuenta:                    true,
			CargarReversionReporte:          true,
			CargarCliente:                   true,
			OrdenadoFecha:                   true,
			CargarComisionImpuesto:          true,
			CargarRetenciones:               true,
			CargarMovimientosRetenciones:    true,
		}
		movrevertidos, err = s.repository.GetMovimiento(filtroRevertidos)
		if err != nil {
			erro = err
			return
		}
	}

	var fechas []string

	// se obtienen movimientos a partir de pagointentos positivos, y se recorre cada movimiento
	if len(pagosintentos) > 0 {
		// se obtienen MOVIMIENTOS
		mov, err := s.repository.GetMovimiento(filtroMov)
		if err != nil {
			erro = err
			return
		}

		// En este punto se tienen movimientos C positivos y movimientos revertidos en las var mov y movrevertidos

		for _, m_fecha := range mov {
			// vars para acumular retenciones por gravamen
			var (
				movRetencionGanancias, movRetencionIVA, movRetencionIIBB entities.Monto
			)

			var fecha_deposito string
			// se obtiene el valor de la var fecha_deposito que determina la fecha de cada reporte de rendicion
			for _, deposito := range depositos {
				if m_fecha.Pagointentos.ID == uint(deposito.pagoIntentoId) {
					fecha_deposito = deposito.fechaDeposito
				}
			}
			// para no repetir reportes, se controla por la fecha y se va guardando en slice para comparar
			if !contains(fechas, fecha_deposito) {
				// se guarda la feha en slice string fechas
				fechas = append(fechas, fecha_deposito)

				// Var que representa un Reporte de rendicion
				var rendiciones reportedtos.DetallesRendicion

				// cada rendicion tiene una fecha especifica distinta
				rendiciones.Fecha = fecha_deposito
				rendiciones.Nombre = (m_fecha.Cuenta.Cliente.Cliente + "-" + fecha_deposito)

				entityControl := entities.Reporte{
					Tiporeporte:    "rendiciones",
					Fecharendicion: fecha_deposito,
				}

				opcionesBusqueda := filtros_reportes.BusquedaReporteFiltro{
					SigNumero: true,
				}
				// teniendo en cuenta la fecha unica se obtiene el numero de reporte correspondiente
				nroReporteUint, err := s.repository.GetLastReporteEnviadosRepository(entityControl, opcionesBusqueda)
				if err != nil {
					erro = err
					return
				}
				if nroReporteUint != 0 {
					nroReporteString := strconv.FormatUint(uint64(nroReporteUint), 10)
					rendiciones.NroReporte = nroReporteString
				}

				// Cada iteracion de este for equivale a cada operacion que conforman un reporte de cliente
				for _, m := range mov {
					// var para acumular retenciones por movimiento
					var totalRetencionPorMovimiento entities.Monto

					var fecha_deposito_otro string
					for _, deposito := range depositos {
						if m.Pagointentos.ID == uint(deposito.pagoIntentoId) {
							fecha_deposito_otro = deposito.fechaDeposito
						}
					}
					if fecha_deposito == fecha_deposito_otro {
						// Se calcula la comision e impuesto de cada movimiento
						var comision entities.Monto
						var iva entities.Monto
						if len(m.Movimientocomisions) > 0 {
							comision = m.Movimientocomisions[len(m.Movimientocomisions)-1].Monto + m.Movimientocomisions[len(m.Movimientocomisions)-1].Montoproveedor
							iva = m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Monto + m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Montoproveedor
						} else {
							comision = 0
							iva = 0
						}

						// Para cada movimiento (m) se suman sus retenciones si las tuviere
						if len(m.Movimientoretencions) > 0 {
							for _, movimiento_retencion := range m.Movimientoretencions {
								totalRetencionPorMovimiento += movimiento_retencion.ImporteRetenido
							}
							movRetencionGanancias += importeRetencionByName("ganancias", m.Movimientoretencions)
							movRetencionIVA += importeRetencionByName("iva", m.Movimientoretencions)
							movRetencionIIBB += importeRetencionByName("iibb", m.Movimientoretencions)
						}

						cantidadBoletas := len(m.Pagointentos.Pago.Pagoitems)
						if !contains(controlIds, fmt.Sprint(m.PagointentosId)) {
							controlIds = append(controlIds, fmt.Sprint(m.PagointentosId))

							rendicion := reportedtos.ResponseReportesRendiciones{
								PagoIntentoId:           m.PagointentosId,
								Cuenta:                  m.Cuenta.Cuenta,
								Id:                      m.Pagointentos.Pago.ExternalReference,
								FechaCobro:              m.Pagointentos.PaidAt.Format("02-01-2006"),
								ImporteCobrado:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Pagointentos.Amount.Float64(), 2))),
								ImporteDepositado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Monto.Float64(), 2))),
								CantidadBoletasCobradas: fmt.Sprintf("%v", cantidadBoletas),
								Comision:                fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 4))),
								Iva:                     fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 4))),
								Concepto:                "Transferencia",
								FechaDeposito:           fecha_deposito_otro,
								Retenciones:             fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalRetencionPorMovimiento.Float64(), 2))),
							}

							rendiciones.Rendiciones = append(rendiciones.Rendiciones, rendicion)
							auxCobrado := m.Pagointentos.Amount.Float64()
							auxDepositado := m.Monto.Float64()

							rendiciones.TotalCobrado += auxCobrado
							rendiciones.TotalRendido += auxDepositado
							rendiciones.TotalComision += comision.Float64()
							rendiciones.TotalIva += iva.Float64()
							rendiciones.CantidadOperaciones += 1

						}

					}

				} // Fin de for _, m := range mov. Fin de cada operacion de un reporte
				// Totalizar las retenciones por cada reporte, por tipo de gravamen
				rendiciones.TotalRetGanancias = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(movRetencionGanancias.Float64(), 2)))
				rendiciones.TotalRetIVA = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(movRetencionIVA.Float64(), 2)))
				rendiciones.TotalRetIIBB = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(movRetencionIIBB.Float64(), 2)))
				res.DetallesRendiciones = append(res.DetallesRendiciones, rendiciones)

			}
		} // Fin de for _, m_fecha := range mov. Fin de un reporte

	} // Fin de if len(pagosintentos) > 0

	if len(movrevertidos) > 0 {
		var fechasReversion []string

		for _, mr := range movrevertidos {

			// vars para acumular retenciones por gravamen
			var (
				movRetencionGanancias, movRetencionIVA, movRetencionIIBB entities.Monto
			)

			var fecha_deposito string
			for _, deposito := range depositosRev {
				if mr.Pagointentos.ID == uint(deposito.pagoIntentoId) {
					fecha_deposito = deposito.fechaDeposito
				}
			}

			if !contains(fechasReversion, fecha_deposito) {
				fechasReversion = append(fechasReversion, fecha_deposito)

				fechas = append(fechas, fecha_deposito)
				var rendiciones reportedtos.DetallesRendicion
				for i := 0; i < len(res.DetallesRendiciones); i++ {
					if res.DetallesRendiciones[i].Fecha == fecha_deposito {
						rendiciones = res.DetallesRendiciones[i]
					}
				}

				for _, mr2 := range movrevertidos {
					var fecha_deposito_otro string
					for _, deposito := range depositosRev {
						if mr2.Pagointentos.ID == uint(deposito.pagoIntentoId) {
							fecha_deposito_otro = deposito.fechaDeposito
						}
					}
					if fecha_deposito == fecha_deposito_otro {

						// var para acumular retenciones por movimiento
						var totalRetencionPorMovimiento entities.Monto

						var comision entities.Monto
						var iva entities.Monto
						if len(mr2.Movimientocomisions) > 0 {
							comision = mr2.Movimientocomisions[len(mr2.Movimientocomisions)-1].Monto + mr2.Movimientocomisions[len(mr2.Movimientocomisions)-1].Montoproveedor
							iva = mr2.Movimientoimpuestos[len(mr2.Movimientoimpuestos)-1].Monto + mr2.Movimientoimpuestos[len(mr2.Movimientoimpuestos)-1].Montoproveedor
						} else {
							comision = 0
							iva = 0
						}

						// Para cada movimiento (m) se suman sus retenciones si las tuviere
						if len(mr2.Movimientoretencions) > 0 {
							for _, movimiento_retencion := range mr2.Movimientoretencions {
								totalRetencionPorMovimiento += movimiento_retencion.ImporteRetenido
							}
							movRetencionGanancias += importeRetencionByName("ganancias", mr2.Movimientoretencions)
							movRetencionIVA += importeRetencionByName("iva", mr2.Movimientoretencions)
							movRetencionIIBB += importeRetencionByName("iibb", mr2.Movimientoretencions)
						}

						rendicion := reportedtos.ResponseReportesRendiciones{
							PagoIntentoId:     mr2.PagointentosId,
							Cuenta:            mr2.Cuenta.Cuenta,
							Id:                mr2.Pagointentos.Pago.ExternalReference,
							ImporteDepositado: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(mr2.Monto.Float64(), 2))),
							Concepto:          "Reversion",
							FechaDeposito:     fecha_deposito_otro,
							Comision:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 2))),
							Iva:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 2))),
							Retenciones:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalRetencionPorMovimiento.Float64(), 2))),
						}

						rendiciones.Rendiciones = append(rendiciones.Rendiciones, rendicion)
						auxDepositado := mr2.Monto.Float64()

						rendiciones.TotalRendido += auxDepositado
						rendiciones.CantidadOperaciones += 1
						rendiciones.TotalReversion += auxDepositado
					}

				}

				for i := 0; i < len(res.DetallesRendiciones); i++ {
					if res.DetallesRendiciones[i].Fecha == fecha_deposito {
						res.DetallesRendiciones[i] = rendiciones
					}
				}

			}

		} // Fin de for _, mr := range movrevertidos

		// fmt.Print(depositosRev)

	}

	sort.Slice(res.DetallesRendiciones, func(i, j int) bool {
		if res.DetallesRendiciones[i].Fecha != "" && res.DetallesRendiciones[j].Fecha != "" {
			return s.commons.ConvertirFecha(res.DetallesRendiciones[i].Fecha) > s.commons.ConvertirFecha(res.DetallesRendiciones[j].Fecha)
		}
		return false
	})

	for _, rendicionesFinales := range res.DetallesRendiciones {
		res.CantidadRegistros += 1
		res.Total += rendicionesFinales.TotalRendido
	}

	return
}

func (s *reportesService) GetReversionesClientesService(request reportedtos.RequestReporteClientes) (res reportedtos.ResponseReversionesClientes, erro error) {
	erro = request.Validar()
	if erro != nil {
		return
	}
	// convert fecha 01-01-2022 a 2020-01-01
	fechaInicio := s.commons.ConvertirFecha(request.FechaInicio)
	fechaFin := s.commons.ConvertirFecha(request.FechaFin)

	// parse date
	fechaInicioTime, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	// parse date
	fechaFinTime, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		erro = errors.New(ERROR_CONVERTIR_FECHA)
		return
	}

	filtro := reportedtos.RequestPagosPeriodo{
		ClienteId:     uint64(request.ClienteId),
		FechaInicio:   fechaInicioTime,
		FechaFin:      fechaFinTime,
		OrdenadoFecha: true,
	}
	filtro_validacion := reportedtos.ValidacionesFiltro{
		Inicio:  true,
		Fin:     true,
		Cliente: true,
	}

	listaPagos, err := s.repository.GetReversionesReportes(filtro, filtro_validacion)
	if err != nil {
		erro = err
		return
	}

	var pagosRevertidos []uint64
	for _, pagosRevertido := range listaPagos {
		pagosRevertidos = append(pagosRevertidos, uint64(pagosRevertido.PagointentosID))
	}

	type fechaDeposito struct {
		pagoIntentoId uint64
		fechaDeposito string
	}

	var depositosRev []fechaDeposito
	var pagosintentosrevertidos []uint64

	// TODO se obtienen transferencias del cliente indicado en el filtro
	listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
	if err != nil {
		erro = err
		return
	}
	if len(listaPagos) > 0 {
		if len(listaTransferencia) > 0 {
			for _, transferencia := range listaTransferencia {
				if transferencia.Reversion {
					transferenciaDepo := fechaDeposito{
						pagoIntentoId: transferencia.Movimiento.PagointentosId,
						fechaDeposito: transferencia.FechaOperacion.Format("02-01-2006"),
					}
					pagosintentosrevertidos = append(pagosintentosrevertidos, transferencia.Movimiento.PagointentosId)
					depositosRev = append(depositosRev, transferenciaDepo)
				}
			}
		}

		var fechas []string

		for _, m_fecha := range listaPagos {
			var fecha_deposito string
			for _, deposito := range depositosRev {
				if m_fecha.PagointentosID == uint(deposito.pagoIntentoId) {
					fecha_deposito = deposito.fechaDeposito
				}
			}
			if !contains(fechas, fecha_deposito) {
				fechas = append(fechas, fecha_deposito)
				var reversiones reportedtos.DetallesReversiones
				reversiones.Fecha = fecha_deposito
				reversiones.Nombre = (m_fecha.PagoIntento.Pago.PagosTipo.Cuenta.Cliente.Cliente + "-" + fecha_deposito)
				for _, value := range listaPagos {

					var fecha_deposito_otro string
					for _, deposito := range depositosRev {
						if value.PagointentosID == uint(deposito.pagoIntentoId) {
							fecha_deposito_otro = deposito.fechaDeposito
						}
					}
					if fecha_deposito == fecha_deposito_otro {
						var revertido reportedtos.Reversiones
						var pagoRevertido reportedtos.PagoRevertido
						var itemsRevertido []reportedtos.ItemsRevertidos
						//var itemRevertido reportedtos.ItemsRevertidos
						var intentoPagoRevertido reportedtos.IntentoPagoRevertido
						revertido.EntityToReversiones(value)
						pagoRevertido.EntityToPagoRevertido(value.PagoIntento.Pago)
						if len(value.PagoIntento.Pago.Pagoitems) > 0 {
							for _, valueItem := range value.PagoIntento.Pago.Pagoitems {
								var itemRevertido reportedtos.ItemsRevertidos
								itemRevertido.EntityToItemsRevertidos(valueItem)
								itemsRevertido = append(itemsRevertido, itemRevertido)
							}
						}
						intentoPagoRevertido.EntityToIntentoPagoRevertido(value.PagoIntento)
						pagoRevertido.Items = itemsRevertido
						pagoRevertido.IntentoPago = intentoPagoRevertido
						revertido.PagoRevertido = pagoRevertido
						reversiones.Reversiones = append(reversiones.Reversiones, revertido)

						auxMonto, _ := strconv.ParseFloat(revertido.Monto, 64)
						reversiones.TotalMonto += auxMonto
						reversiones.CantidadOperaciones += 1
					}
				}

				res.DetallesReversiones = append(res.DetallesReversiones, reversiones)
				res.CantidadRegistros += 1
				res.Total += reversiones.TotalMonto
			}
		}
	}
	sort.Slice(res.DetallesReversiones, func(i, j int) bool {
		if res.DetallesReversiones[i].Fecha != "" && res.DetallesReversiones[j].Fecha != "" {
			return s.commons.ConvertirFecha(res.DetallesReversiones[i].Fecha) > s.commons.ConvertirFecha(res.DetallesReversiones[j].Fecha)
		}
		return false
	})
	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s *reportesService) GetPeticiones(request reportedtos.RequestPeticiones) (response reportedtos.ResponsePeticiones, erro error) {

	var fechaI time.Time
	var fechaF time.Time
	if request.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	}
	if request.Vendor == "" {
		request.Vendor = "ApiLink"
	}

	listaPeticiones, total, err := s.repository.GetPeticionesReportes(request)
	if err != nil {
		erro = err
		return
	}

	if len(listaPeticiones) > 0 {
		var resultPeticiones []reportedtos.ResponseDetallePeticion
		for _, peticion := range listaPeticiones {

			resultPeticiones = append(resultPeticiones, reportedtos.ResponseDetallePeticion{
				Operacion: peticion.Operacion,
				Fecha:     peticion.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		response = reportedtos.ResponsePeticiones{
			FechaComienzo:   request.FechaInicio.Format("2006-01-02"),
			FechaFin:        request.FechaFin.Format("2006-01-02"),
			TotalPeticiones: int(total),
			Data:            resultPeticiones,
			LastPage:        int(math.Ceil(float64(total) / float64(request.Size))),
		}

	}

	return

}

func (s *reportesService) GetLogs(request reportedtos.RequestLogs) (response reportedtos.ResponseLogs, erro error) {

	var fechaI time.Time
	var fechaF time.Time
	if request.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	}

	listaLogs, total, err := s.repository.GetLogs(request)
	if err != nil {
		erro = err
		return
	}

	if len(listaLogs) > 0 {
		var resultLogs []reportedtos.ResponseDetalleLog
		for _, log := range listaLogs {

			resultLogs = append(resultLogs, reportedtos.ResponseDetalleLog{
				Mensaje:       log.Mensaje,
				Fecha:         log.CreatedAt.Format("2006-01-02 15:04:05"),
				Funcionalidad: log.Funcionalidad,
				Tipo:          string(log.Tipo),
			})
		}

		response = reportedtos.ResponseLogs{
			FechaComienzo: request.FechaInicio.Format("2006-01-02"),
			FechaFin:      request.FechaFin.Format("2006-01-02"),
			TotalLogs:     int(total),
			Data:          resultLogs,
			LastPage:      int(math.Ceil(float64(total) / float64(request.Size))),
		}

	}

	return

}

func (s *reportesService) GetNotificaciones(request reportedtos.RequestNotificaciones) (response reportedtos.ResponseNotificaciones, erro error) {

	var fechaI time.Time
	var fechaF time.Time
	if request.FechaInicio.IsZero() {
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	}

	listaNotif, total, err := s.repository.GetNotificaciones(request)
	if err != nil {
		erro = err
		return
	}

	if len(listaNotif) > 0 {
		var resultNotifs []reportedtos.ResponseDetalleNotificaciones
		for _, notif := range listaNotif {

			resultNotifs = append(resultNotifs, reportedtos.ResponseDetalleNotificaciones{
				Descripcion: notif.Descripcion,
				Fecha:       notif.CreatedAt.Format("2006-01-02 15:04:05"),
				Tipo:        string(notif.Tipo),
			})
		}

		response = reportedtos.ResponseNotificaciones{
			FechaComienzo:       request.FechaInicio.Format("2006-01-02"),
			FechaFin:            request.FechaFin.Format("2006-01-02"),
			TotalNotificaciones: int(total),
			Data:                resultNotifs,
			LastPage:            int(math.Ceil(float64(total) / float64(request.Size))),
		}

	}

	return

}

func (s *reportesService) GetReportesEnviadosService(request reportedtos.RequestReportesEnviados) (response reportedtos.ResponseReportesEnviados, erro error) {

	// se recibe respuesta del repositorio con datos de base de datos. o un error
	listaReportes, totalFilas, erro := s.repository.GetReportesEnviadosRepository(request)

	if erro != nil {
		erro = errors.New(erro.Error())
		return
	}

	// pasar las entidades a DTO response correspondiente

	var resTemporal reportedtos.ResponseReporteEnviado

	for _, reporte := range listaReportes {

		resTemporal.EntityToDto(reporte)

		response.Reportes = append(response.Reportes, resTemporal)
	}

	// paginacion
	if request.Number > 0 && request.Size > 0 {
		response.Meta = setPaginacion(uint32(request.Number), uint32(request.Size), totalFilas)
	}

	return
}

func (s *reportesService) GetReportesPdfService(request reportedtos.RequestReportesEnviados, cliente entities.Cliente, cuenta entities.Cuenta) (response []reportedtos.ResponseClientesReportes, erro error) {

	request_rendiciones := reportedtos.RequestReporteClientes{
		FechaInicio: s.commons.ConvertirFechaToDDMMYYYY(request.FechaInicio[:10]),
		FechaFin:    s.commons.ConvertirFechaToDDMMYYYY(request.FechaFin[:10]),
		ClienteId:   int(cliente.ID),
		CuentaId:    int(cuenta.ID),
	}

	responseRendicionesClientes, _ := s.GetRendicionesClientesService(request_rendiciones)

	listaReportesClientes := responseRendicionesClientes.DetallesRendiciones
	if len(listaReportesClientes) == 0 {
		mensaje := fmt.Sprintf("la lista de reportes de rendiciones de clientes esta vacía para el cliente %s", cliente.Cliente)
		s.util.BuildLog(errors.New(mensaje), "GetReportesPdfService")
		return
	}

	registroReporte := entities.Reporte{
		Tiporeporte:    "rrm", // Reporte Rendiciones Mensual
		Fecharendicion: s.commons.ConvertirFechaToDDMMYYYY(request.FechaFin[:10]),
	}
	filtroReporte := filtros_reportes.BusquedaReporteFiltro{
		SigNumero: true,
		LastRrm:   true,
	}

	siguienteNroReporte, err := s.repository.GetLastReporteEnviadosRepository(registroReporte, filtroReporte)
	if err != nil {
		erro = err
		return
	}

	registroReporte.Nro_reporte = siguienteNroReporte
	registroReporte.Cliente = cliente.Cliente

	// Guardar fehcas del periodo: incio y fin
	registroReporte.PeriodoInicio, _ = s.commons.DateStringToTimeFirstMoment(request.FechaInicio)
	registroReporte.PeriodoFin, _ = s.commons.DateStringToTimeLastMoment(request.FechaFin)

	// se crea el reporte mensual de rendiciones rrm
	ruta, totales, err := s.ReportesRendicionesPDF(listaReportesClientes, cliente, cuenta, siguienteNroReporte, request.FechaInicio, request.FechaFin)

	if err != nil {
		erro = err
		logs.Error(err.Error())
	}

	registroReporte.Totalcobrado = totales.TotalCobranza
	registroReporte.Totalrendido = totales.TotalRendicion
	registroReporte.TotalRetencionGanancias = totales.TotalRetencionGanancias
	registroReporte.TotalRetencionIva = totales.TotalRetencionIVA
	registroReporte.TotalRetencionIibb = totales.TotalRetencionIngresosBrutos
	registroReporte.TotalRetenido = totales.TotalRetenido

	reportes := []entities.Reporte{registroReporte}

	reporteDatos := reportedtos.ResponseClientesReportes{
		Clientes:        cliente.Cliente,
		Email:           []string{cliente.Email},
		Cuit:            cliente.Cuit,
		Fecha:           time.Now().Format("02-01-2006"),
		RutaArchivo:     ruta,
		Reporte:         reportes, // slice de reportes del cliente
		CantOperaciones: fmt.Sprint(totales.TotalOperaciones),
	}

	//Subir archivo generado al cloud
	rutaDirectorio := config.DIR_BASE + config.DIR_REPORTE + "/" + ruta + ".pdf"
	// Leer el contenido del archivo
	data, err := ioutil.ReadFile(rutaDirectorio)
	if err != nil {
		erro = err
		return
	}

	// SUBIR CLOUD STORAGE

	nombreCarpeta := cliente.Cliente
	rutaCloud := "retenciones/comprobantes/" + nombreCarpeta
	// rutaDetalle: origen del archivo
	// rutaCloud: destino del archivo
	ctx := context.Background()

	err = s.UploadTxtFile(ctx, rutaDirectorio, rutaCloud, data)
	if err != nil {
		erro = err
		s.util.BuildLog(erro, "Fallo al subir RRM")
		return
	}

	response = append(response, reporteDatos)

	return
}

func setPaginacion(number uint32, size uint32, total int64) (meta dtos.Meta) {
	from := (number - 1) * size
	lastPage := math.Ceil(float64(total) / float64(size))

	meta = dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}

	return
}

func (s *reportesService) TratamientoReporteMensualPagos(reporteCliente []reportedtos.ResponseClientesReportes, ordenCobranza bool) (reporte reportedtos.ReporteMensual, erro error) {

	var totalReporteMensual float64
	var totalComisionMensual float64

	for _, reporteCliente := range reporteCliente {

		FormatTotalComision := intercambiarCommasDots(reporteCliente.TotalComision)

		valorTotalComisionCliente, err := strconv.ParseFloat(FormatTotalComision, 64)
		if err != nil {
			erro = err
			return
		}

		totalComisionMensual += valorTotalComisionCliente

		FormatTotal := intercambiarCommasDots(reporteCliente.TotalCobrado)

		valorTotalCobradoCliente, err := strconv.ParseFloat(FormatTotal, 64)
		if err != nil {
			erro = err
			return
		}
		totalReporteMensual += valorTotalCobradoCliente

		reporteClienteMensual := reportedtos.ReporteMensualCliente{
			Cliente:                 reporteCliente.Clientes,
			TotalMensual:            reporteCliente.TotalCobrado,
			TotalOperacionesCliente: 0,
			TotalComisionMensual:    reporteCliente.TotalComision,
		}

		for _, pago := range reporteCliente.Pagos {

			reporte.TotalOperaciones += 1
			reporteClienteMensual.TotalOperacionesCliente += 1

			// Control Si se registro algun pago de ese día
			var controlDia bool
			var indiceDia int
			for indice, cobranzasDia := range reporte.CobranzasDías {
				if cobranzasDia.FechaCobranzas == pago.FechaPago {
					controlDia = true
					indiceDia = indice
				}
			}
			FormatMonto := intercambiarCommasDots(pago.Monto)
			totalConversion, err := strconv.ParseFloat(FormatMonto, 10)
			if err != nil {
				erro = err
				return
			}

			FormatMontoComision := intercambiarCommasDots(pago.Comision)
			totalConversionComision := float64(0)
			if FormatMontoComision != "" {
				totalConversionComision, err = strconv.ParseFloat(FormatMontoComision, 10)
				if err != nil {
					erro = err
					return
				}
			}

			if controlDia {

				pagoCliente := reportedtos.PagoCliente{
					Cliente: reporteCliente.Clientes,
					Pago:    pago,
				}

				reporte.CobranzasDías[indiceDia].PagosDia = append(reporte.CobranzasDías[indiceDia].PagosDia, pagoCliente)
				reporte.CobranzasDías[indiceDia].CobranzaTotalDia += totalConversion
				reporte.CobranzasDías[indiceDia].OperacionesTotalDia += 1
				reporte.CobranzasDías[indiceDia].ComisionTotalDia += totalConversionComision
			} else {
				// Sino registra con datos iniciales
				pagoInicial := reportedtos.PagoCliente{
					Cliente: reporteCliente.Clientes,
					Pago:    pago,
				}
				arrayInicial := []reportedtos.PagoCliente{pagoInicial}

				cobranzaDia := reportedtos.ReporteMensualData{
					FechaCobranzas:      pago.FechaPago,
					CobranzaTotalDia:    totalConversion,
					PagosDia:            arrayInicial,
					OperacionesTotalDia: 1,
				}
				reporte.CobranzasDías = append(reporte.CobranzasDías, cobranzaDia)
			}

		}

		reporte.TotalCliente = append(reporte.TotalCliente, reporteClienteMensual)
	}

	reporte.TotalCobranza = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalReporteMensual, 2)))
	reporte.TotalComision = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalComisionMensual, 2)))
	sort.Slice(reporte.CobranzasDías, func(i, j int) bool {
		if reporte.CobranzasDías[i].FechaCobranzas != "" && reporte.CobranzasDías[j].FechaCobranzas != "" {
			return s.commons.ConvertirFecha(reporte.CobranzasDías[i].FechaCobranzas) < s.commons.ConvertirFecha(reporte.CobranzasDías[j].FechaCobranzas)
		}
		return false

	})

	sort.Slice(reporte.TotalCliente, func(i, j int) bool {
		if ordenCobranza {
			FormatTotal := intercambiarCommasDots(reporte.TotalCliente[i].TotalMensual)
			FormatTotal2 := intercambiarCommasDots(reporte.TotalCliente[j].TotalMensual)
			valorTotalCliente, _ := strconv.ParseFloat(FormatTotal, 64)
			valorTotalCliente2, _ := strconv.ParseFloat(FormatTotal2, 64)
			return valorTotalCliente > valorTotalCliente2
		} else {
			return reporte.TotalCliente[i].Cliente < reporte.TotalCliente[j].Cliente
		}

	})

	return
}

func intercambiarCommasDots(cadena string) (respuesta string) {

	intermedio := strings.Replace(cadena, ".", "", (-1))
	respuesta = strings.Replace(intermedio, ",", ".", (-1))
	return
}

// Trae mas informacion de la DB y luego separa por cliente, en vez de realizar X llamadas por clientes a la DB.
func (s *reportesService) GetPagosClientesMensual(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error) {

	var fechaI time.Time
	var fechaF time.Time

	if filtro.FechaInicio.IsZero() {
		// Entro por proceso background
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}

	var idsClientes []uint64
	for _, cliente := range request.Clientes {
		idsClientes = append(idsClientes, uint64(cliente.Id))
	}

	filtroPagos := reportedtos.RequestPagosPeriodo{
		ClientesIds: idsClientes,
		FechaInicio: fechaI,
		FechaFin:    fechaF,
	}

	var listaPagos []reportedtos.DetallesPagosCobranza
	prisma, err := s.repository.GetCobranzasPrisma(filtroPagos)
	apilink, err := s.repository.GetCobranzasApilink(filtroPagos)
	rapipago, err := s.repository.GetCobranzasRapipago(filtroPagos)

	if err != nil {
		erro = err
		return
	}

	listaPagos = append(listaPagos, prisma...)
	listaPagos = append(listaPagos, apilink...)
	listaPagos = append(listaPagos, rapipago...)

	for _, cliente := range request.Clientes {
		var cantoperaciones int64
		var totalcobrado float64
		// var totalIVA entities.Monto
		var totalComision entities.Monto

		var pagos []reportedtos.PagosReportes
		if len(listaPagos) > 0 {
			for _, pago := range listaPagos {
				if pago.Cliente == cliente.Cliente {

					monto := entities.Monto(pago.TotalPago).Float64()
					comision := entities.Monto(pago.Comision)

					totalcobrado += monto
					totalComision += comision
					cantoperaciones = cantoperaciones + 1

					fecha := pago.FechaPago
					// si el pago es debin o offline se debe tomar la fecha de cobro
					if pago.CanalPago == "DEBIN" || pago.CanalPago == "OFFLINE" {
						fecha = pago.FechaCobro
					}

					pagos = append(pagos, reportedtos.PagosReportes{
						Cuenta:    pago.Cuenta,
						Id:        pago.Referencia,
						FechaPago: fecha.Format("02-01-2006"),
						MedioPago: pago.MedioPago,
						Tipo:      pago.CanalPago,
						Estado:    pago.Pagoestado,
						Monto:     s.util.FormatNum(monto),
						Comision:  s.util.FormatNum(comision.Float64()),
					})
				}
			}
		}
		if len(pagos) > 0 {
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:        cliente.Cliente,
				Email:           cliente.Emails, //[]string{cliente.Email},
				RazonSocial:     cliente.RazonSocial,
				Cuit:            cliente.Cuit,
				Fecha:           fechaI.Format("02-01-2006"),
				Pagos:           pagos,
				CantOperaciones: fmt.Sprintf("%v", cantoperaciones),
				TotalCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalcobrado, 2))),
				TipoArchivoPdf:  true,
				TotalComision:   fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalComision.Float64(), 2))),
			})

		}

	}

	return
}

func SepararPagos(request []entities.Pago, clienteId int64) (response []entities.Pago, requestFiltrado []entities.Pago, erro error) {

	requestFiltrado = request
	var indices []int
	for indice, pago := range request {
		if pago.PagosTipo.Cuenta.ClientesID == clienteId {
			response = append(response, pago)
			indices = append(indices, indice)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(indices)))

	for _, indX := range indices {
		requestFiltrado = RemoverPago(requestFiltrado, indX)
	}

	return
}

func RemoverPago(request []entities.Pago, indX int) (requestFiltrado []entities.Pago) {
	if len(request) <= (indX) {
		requestFiltrado = append(request[:(len(request))])
	} else {
		requestFiltrado = append(request[:indX], request[indX+1:]...)
	}

	return
}

func (s *reportesService) GetRendicionesClientesMensual(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestPagosClientes) (response []reportedtos.ResponseClientesReportes, erro error) {

	var fechaI time.Time
	var fechaF time.Time
	if filtro.FechaInicio.IsZero() {
		// si los filtros recibidos son ceros toman la fecha actual
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		// a las fechas se le restan un dia ya sea por backgraund o endpoint
		// fechaI = fechaI.AddDate(0, 0, int(-1))
		// fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		// fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		// fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
		fechaI = filtro.FechaInicio
		fechaF = filtro.FechaFin
	}
	logs.Info(fechaI)
	logs.Info(fechaF)
	for _, cliente := range request.Clientes {
		var totalCobrado entities.Monto
		var total entities.Monto
		var totalReversion entities.Monto
		var cantOperaciones int
		var totalIVA entities.Monto
		var totalComision entities.Monto

		// totales generales
		var rendido entities.Monto

		filtro := reportedtos.RequestPagosPeriodo{
			ClienteId: uint64(cliente.Id),
			// ClienteId:   5,                             //Prueba con cliente 6
			FechaInicio: fechaI, // descomentar esta linea cuando se pasa a dev y produccion
			// se envian pagos del dia anterior
			FechaFin: fechaF,
		}

		// TODO se obtienen transferencias del cliente indicado en el filtro
		listaTransferencia, err := s.repository.GetTransferenciasReportes(filtro)
		if err != nil {
			erro = err
			return
		}
		var pagos []*reportedtos.ResponseReportesRendiciones
		var movrevertidos []entities.Movimiento
		var pagosintentos []uint64
		var pagosintentosrevertidos []uint64
		var filtroMov reportedtos.RequestPagosPeriodo
		var totalCliente reportedtos.ResponseTotales
		if len(listaTransferencia) > 0 {
			for _, transferencia := range listaTransferencia {
				if !transferencia.Reversion {
					pagosintentos = append(pagosintentos, transferencia.Movimiento.PagointentosId)
				}
				if transferencia.Reversion {
					pagosintentosrevertidos = append(pagosintentosrevertidos, transferencia.Movimiento.PagointentosId)
				}
			}
			filtroMov = reportedtos.RequestPagosPeriodo{
				PagoIntentos:                    pagosintentos,
				TipoMovimiento:                  "C",
				CargarComisionImpuesto:          true,
				CargarMovimientosTransferencias: true,
				CargarPagoIntentos:              true,
				CargarCuenta:                    true,
			}
		}

		// en el caso de que existieran reversiones
		if len(pagosintentosrevertidos) > 0 {
			filtroRevertidos := reportedtos.RequestPagosPeriodo{
				PagoIntentos:                    pagosintentosrevertidos,
				TipoMovimiento:                  "C",
				CargarMovimientosTransferencias: true,
				CargarPagoIntentos:              true,
				CargarCuenta:                    true,
				CargarReversionReporte:          true,
			}
			movrevertidos, err = s.repository.GetMovimiento(filtroRevertidos)
			if err != nil {
				erro = err
				return
			}
		}

		if len(pagosintentos) > 0 {
			mov, err := s.repository.GetMovimiento(filtroMov)
			if err != nil {
				erro = err
				return
			}
			var resulRendiciones []*reportedtos.ResponseReportesRendiciones
			for _, m := range mov {
				cantOperaciones = cantOperaciones + 1
				cantidadBoletas := len(m.Pagointentos.Pago.Pagoitems)
				total += m.Monto
				totalCobrado += m.Pagointentos.Amount
				var comision entities.Monto
				var iva entities.Monto
				if len(m.Movimientocomisions) > 0 {
					comision = m.Movimientocomisions[len(m.Movimientocomisions)-1].Monto
					iva = m.Movimientoimpuestos[len(m.Movimientoimpuestos)-1].Monto
				} else {
					comision = 0
					iva = 0
				}
				totalComision += comision
				totalIVA += iva

				resulRendiciones = append(resulRendiciones, &reportedtos.ResponseReportesRendiciones{
					PagoIntentoId:           m.PagointentosId,
					Cuenta:                  m.Cuenta.Cuenta,
					Id:                      m.Pagointentos.Pago.ExternalReference,
					FechaCobro:              m.Pagointentos.PaidAt.Format("02-01-2006"),
					ImporteCobrado:          fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Pagointentos.Amount.Float64(), 2))),
					ImporteDepositado:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(m.Monto.Float64(), 2))),
					CantidadBoletasCobradas: fmt.Sprintf("%v", cantidadBoletas),
					Comision:                fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision.Float64(), 4))),
					Iva:                     fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(iva.Float64(), 4))),
					Concepto:                "Transferencia",
				})
			}

			totalCliente = reportedtos.ResponseTotales{
				// Totales: reportedtos.Totales{
				// 	CantidadOperaciones: fmt.Sprintf("%v", cantOperaciones),
				// 	TotalCobrado:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalCobrado.Float64(), 4))),
				// 	TotalRendido:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(total.Float64(), 4))),
				// },
				Detalles: resulRendiciones,
			}
		}
		// vuelvo a comparar con las transferenicas para asignar fecha de deposito
		if len(totalCliente.Detalles) > 0 {
			for _, transferencia := range listaTransferencia {
				for _, t := range totalCliente.Detalles {
					if transferencia.Movimiento.PagointentosId == t.PagoIntentoId {
						t.FechaDeposito = transferencia.FechaOperacion.Format("02-01-2006")
					}
				}
			}
		}

		if len(movrevertidos) > 0 {
			for _, mr := range movrevertidos {
				totalReversion += mr.Monto
				cantOperaciones = cantOperaciones + 1
				totalCliente.Detalles = append(totalCliente.Detalles, &reportedtos.ResponseReportesRendiciones{
					PagoIntentoId:     mr.PagointentosId,
					Cuenta:            mr.Cuenta.Cuenta,
					FechaCobro:        mr.Pagointentos.PaidAt.Format("02-01-2006"),
					Id:                mr.Pagointentos.Pago.ExternalReference,
					ImporteDepositado: fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(mr.Monto.Float64(), 2))),
					Concepto:          "Reversion",
				})

			}

		}

		if len(totalCliente.Detalles) > 0 {
			for _, transferencia := range listaTransferencia {
				for _, t := range totalCliente.Detalles {
					if transferencia.Movimiento.PagointentosId == t.PagoIntentoId {
						t.FechaDeposito = transferencia.FechaOperacion.Format("02-01-2006")
					}
				}
			}
		}
		pagos = totalCliente.Detalles

		rendido = total + totalReversion
		totalCliente = reportedtos.ResponseTotales{
			Totales: reportedtos.Totales{
				CantidadOperaciones: fmt.Sprintf("%v", cantOperaciones),
				TotalCobrado:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalCobrado.Float64(), 4))),
				TotalRendido:        fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(rendido.Float64(), 4))),
				TotalComision:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalComision.Float64(), 4))),
				TotalIva:            fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalIVA.Float64(), 4))),
			},
		}
		//&ultimo paso
		if len(pagos) > 0 {
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:            cliente.Cliente,
				RazonSocial:         cliente.RazonSocial,
				Cuit:                cliente.Cuit,
				Email:               cliente.Emails,
				Fecha:               fechaI.Format("02-01-2006"),
				Rendiciones:         pagos,
				CantOperaciones:     totalCliente.Totales.CantidadOperaciones,
				TotalCobrado:        totalCliente.Totales.TotalCobrado,
				TotalIva:            totalCliente.Totales.TotalIva,
				TotalComision:       totalCliente.Totales.TotalComision,
				RendicionTotal:      totalCliente.Totales.TotalRendido,
				TipoArchivoPdf:      false,
				GuardarDatosReporte: false,
			})

		}
	}
	return

}

func (s *reportesService) TratamientoReporteMensualRendiciones(reporteCliente []reportedtos.ResponseClientesReportes, ordenCobranza bool) (reporte reportedtos.ReporteMensual, erro error) {

	var totalReporteMensual float64
	var totalComisionMensual float64

	for _, reporteCliente := range reporteCliente {

		FormatTotal := intercambiarCommasDots(reporteCliente.RendicionTotal)

		valorTotalRendidoCliente, err := strconv.ParseFloat(FormatTotal, 64)
		if err != nil {
			erro = err
			return
		}
		totalReporteMensual += valorTotalRendidoCliente

		FormatTotalComision := intercambiarCommasDots(reporteCliente.TotalComision)

		valorTotalComisionCliente, err := strconv.ParseFloat(FormatTotalComision, 64)
		if err != nil {
			erro = err
			return
		}

		totalComisionMensual += valorTotalComisionCliente

		reporteClienteMensual := reportedtos.ReporteMensualCliente{
			Cliente:                 reporteCliente.Clientes,
			TotalMensual:            reporteCliente.RendicionTotal,
			TotalOperacionesCliente: 0,
			TotalComisionMensual:    reporteCliente.TotalComision,
		}

		for _, rendicion := range reporteCliente.Rendiciones {

			reporte.TotalOperaciones += 1
			reporteClienteMensual.TotalOperacionesCliente += 1

			// Control Si se registro algun pago de ese día
			var controlDia bool
			var indiceDia int
			for indice, rendicionesDia := range reporte.RendicionesDías {
				if rendicionesDia.FechaRendicion == rendicion.FechaDeposito {
					controlDia = true
					indiceDia = indice
				}
			}
			FormatMonto := intercambiarCommasDots(rendicion.ImporteDepositado)
			totalConversion, err := strconv.ParseFloat(FormatMonto, 10)
			if err != nil {
				erro = err
				return
			}

			FormatMontoComision := intercambiarCommasDots(rendicion.Comision)
			totalConversionComision := float64(0)
			if FormatMontoComision != "" {
				totalConversionComision, err = strconv.ParseFloat(FormatMontoComision, 10)
				if err != nil {
					erro = err
					return
				}
			}

			if controlDia {

				rendicionCliente := reportedtos.RendicionCliente{
					Cliente:   reporteCliente.Clientes,
					Rendicion: *rendicion,
				}
				reporte.RendicionesDías[indiceDia].RendicionesDia = append(reporte.RendicionesDías[indiceDia].RendicionesDia, rendicionCliente)
				reporte.RendicionesDías[indiceDia].RendicionTotalDia += totalConversion
				reporte.RendicionesDías[indiceDia].OperacionesTotalDia += 1
				reporte.RendicionesDías[indiceDia].ComisionTotalDia += totalConversionComision
			} else {
				// Sino registra con datos iniciales
				rendicionInicial := reportedtos.RendicionCliente{
					Cliente:   reporteCliente.Clientes,
					Rendicion: *rendicion,
				}
				arrayInicial := []reportedtos.RendicionCliente{rendicionInicial}

				rendicionDia := reportedtos.ReporteMensualDataRendiciones{
					FechaRendicion:      rendicion.FechaDeposito,
					RendicionTotalDia:   totalConversion,
					RendicionesDia:      arrayInicial,
					OperacionesTotalDia: 1,
					ComisionTotalDia:    totalConversionComision,
				}
				reporte.RendicionesDías = append(reporte.RendicionesDías, rendicionDia)
			}

		}

		reporte.TotalCliente = append(reporte.TotalCliente, reporteClienteMensual)
	}

	reporte.TotalRendicion = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalReporteMensual, 2)))
	reporte.TotalComision = fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalComisionMensual, 2)))
	sort.Slice(reporte.RendicionesDías, func(i, j int) bool {
		if reporte.RendicionesDías[i].FechaRendicion != "" && reporte.RendicionesDías[j].FechaRendicion != "" {
			return s.commons.ConvertirFecha(reporte.RendicionesDías[i].FechaRendicion) < s.commons.ConvertirFecha(reporte.RendicionesDías[j].FechaRendicion)
		}
		return false

	})

	sort.Slice(reporte.TotalCliente, func(i, j int) bool {
		if ordenCobranza {
			FormatTotal := intercambiarCommasDots(reporte.TotalCliente[i].TotalMensual)
			FormatTotal2 := intercambiarCommasDots(reporte.TotalCliente[j].TotalMensual)
			valorTotalCliente, _ := strconv.ParseFloat(FormatTotal, 64)
			valorTotalCliente2, _ := strconv.ParseFloat(FormatTotal2, 64)
			return valorTotalCliente > valorTotalCliente2
		} else {
			return reporte.TotalCliente[i].Cliente < reporte.TotalCliente[j].Cliente
		}

	})

	return
}

func (s *reportesService) GetPagos(request administraciondtos.ResponseFacturacionPaginado, filtro reportedtos.RequestCobranzasDiarias) (response []reportedtos.ResponseClientesReportes, erro error) {

	// OBTENER ESTADOS DE LOS PAGOS:

	filtroEstadoPendiente := filtros.PagoEstadoFiltro{
		Nombre: "pending",
	}
	estadoPendiente, err := s.administracion.GetPagoEstado(filtroEstadoPendiente)
	if err != nil {
		erro = err
		return
	}

	//aprobado (credito , debito y offline)
	paid, erro := s.util.FirstOrCreateConfiguracionService("PAID", "Nombre del estado aprobado", "Paid")
	if erro != nil {
		return
	}
	filtroPagosEstado := filtros.PagoEstadoFiltro{
		Nombre: paid,
	}
	estado_paid, err := s.administracion.GetPagoEstado(filtroPagosEstado)
	if err != nil {
		erro = err
		return
	}
	//si no se obtiene el estado del pago no se puede seguir
	if estado_paid[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	//autorizado (debin)
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, err := s.administracion.GetPagoEstado(filtroPagoEstado)

	if err != nil {
		erro = err
		return
	}

	//si no se obtiene el estado del pago no se puede seguir
	if pagoEstadoAcreditado[0].ID < 1 {
		erro = fmt.Errorf(ERROR_PAGO_ESTADO_ID_AUTORIZADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "GetPagosClientes",
			Mensaje:       ERROR_PAGO_ESTADO_ID_AUTORIZADO,
		}
		err := s.util.CreateLogService(log)
		if err != nil {
			erro = err
			logs.Info("GetPagosClientes reportes clientes." + erro.Error())
		}
		return
	}

	// SE DEFINEN VARIABLES TOTALES
	var pagoestados []uint
	var fechaI time.Time
	var fechaF time.Time
	pagoestados = append(pagoestados, estado_paid[0].ID, pagoEstadoAcreditado[0].ID)
	if filtro.FechaInicio.IsZero() {
		// Entro por proceso background
		fechaI, fechaF, erro = s.commons.FormatFecha()
		if erro != nil {
			return
		}
		fechaI = fechaI.AddDate(0, 0, int(-1))
		fechaF = fechaF.AddDate(0, 0, int(-1))
	} else {
		fechaI = filtro.FechaInicio.AddDate(0, 0, int(-1))
		fechaF = filtro.FechaFin.AddDate(0, 0, int(-1))
	}
	var total float64
	var comisiontotal float64
	var cantoperacionestotal int64
	var totales reportedtos.TotalCobranzasDiarias
	for _, cliente := range request.Clientes {
		var cantoperaciones int64
		var totalcobrado float64
		var totalcomisiones float64
		filtroPagos := reportedtos.RequestPagosPeriodo{
			ClienteId:   uint64(cliente.Id),
			FechaInicio: fechaI,
			FechaFin:    fechaF,
			PagoEstados: pagoestados,
		}

		listaPagos, err := s.repository.GetPagosReportes(filtroPagos, estadoPendiente[0].ID)
		if err != nil {
			erro = err
			return
		}
		var pagos []reportedtos.PagosReportes
		if len(listaPagos) > 0 {
			for _, pago := range listaPagos {
				// monto := s.util.ToFixed((pago.Amount.Float64()), 4)
				cantoperaciones = cantoperaciones + 1
				totalcobrado += pago.PagoIntentos[len(pago.PagoIntentos)-1].Amount.Float64()
				medio_pago, _ := s.commons.RemoveAccents(pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.Mediopago)
				var comision float64
				var mov_temporal entities.Movimientotemporale
				if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientotemporale) > 0 {
					mov_temporal = pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientotemporale[len(pago.PagoIntentos[len(pago.PagoIntentos)-1].Movimientotemporale)-1]
				} else {
					logs.Info(fmt.Sprintf("No se encontro movimientotemporale para el pago %v", pago.ID))
				}
				if len(mov_temporal.Movimientocomisions) > 0 {
					comision = mov_temporal.Movimientocomisions[len(mov_temporal.Movimientocomisions)-1].Monto.Float64()
				}
				totalcomisiones += comision

				pagos = append(pagos, reportedtos.PagosReportes{
					Cuenta:    pago.PagosTipo.Cuenta.Cuenta,
					Id:        pago.ExternalReference,
					FechaPago: pago.PagoIntentos[len(pago.PagoIntentos)-1].PaidAt.Format("02-01-2006"),
					MedioPago: medio_pago,
					Tipo:      pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.Channel.Nombre,
					Estado:    string(pago.PagoEstados.Nombre),
					Monto:     fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(pago.PagoIntentos[len(pago.PagoIntentos)-1].Amount.Float64(), 2))),
					Comision:  fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comision, 2))),
				})
			}
		}
		if len(pagos) > 0 {
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:        cliente.Cliente,
				Email:           cliente.Emails, //[]string{cliente.Email},
				RazonSocial:     cliente.RazonSocial,
				Cuit:            cliente.Cuit,
				Fecha:           fechaI.Format("02-01-2006"),
				Pagos:           pagos,
				CantOperaciones: fmt.Sprintf("%v", cantoperaciones),
				TotalCobrado:    fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalcobrado, 2))),
				TotalComision:   fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(totalcomisiones, 2))),
				TipoArchivoPdf:  true,
			})
		}
		total += totalcobrado
		comisiontotal += totalcomisiones
		cantoperacionestotal += cantoperaciones
	}
	//Verifico que existan pagos que mostrar, antes de crear el pdf y enviar los correos.
	if len(response) > 0 {
		totales = reportedtos.TotalCobranzasDiarias{
			Total:               fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(total, 2))),
			CantidadOperaciones: fmt.Sprintf("%v", cantoperacionestotal),
			ComisionTotal:       fmt.Sprintf("%v", s.util.FormatNum(s.util.ToFixed(comisiontotal, 2))),
		}

		fechaNombre := fechaI.Format("02-01-2006")
		cliente := administraciondtos.ResponseFacturacion{
			Cliente:     "WEE",
			RazonSocial: "TelCo",
			Cuit:        "30716550849",
		}
		err = GetCobranzasPdf(response, cliente, fechaNombre, totales)
		if err != nil {
			erro = err
			logs.Error(err.Error())
		} else {
			// si se crea el pdf agregar informacion para enviar al correo
			var pg []reportedtos.PagosReportes
			pg = append(pg, reportedtos.PagosReportes{
				Cuenta: response[len(response)-1].Pagos[len(response[len(response)-1].Pagos)-1].Cuenta,
			})

			emailsTelco := filtro.Email
			if len(filtro.Email) == 0 {
				emailsTelco, err = s.util.GetCorreosTelco()
				if err != nil {
					erro = err
					logs.Info("GetPagosClientes reportes clientes." + erro.Error())

					return
				}
			}

			response = []reportedtos.ResponseClientesReportes{}
			response = append(response, reportedtos.ResponseClientesReportes{
				Clientes:        cliente.Cliente,
				Fecha:           fechaI.Format("02-01-2006"),
				RazonSocial:     cliente.RazonSocial,
				Cuit:            cliente.Cuit,
				Email:           emailsTelco,
				Pagos:           pg,
				CantOperaciones: totales.CantidadOperaciones,
				TotalCobrado:    totales.Total,
				TipoArchivoPdf:  true,
			})
		}
		return
	}
	return
}

func (s *reportesService) SendReporteRetencionComprobante(RRComprobante reportedtos.RequestRRComprobante) (control bool, erro error) {

	// obtener una lista de comprobantes desde segun filtro request RRComprobante
	listaComprobantes, erro := s.repository.GetComprobantesRepository(RRComprobante)

	if erro != nil {
		s.util.BuildLog(erro, "SendReporteRetencionComprobante")
		return
	}

	if len(listaComprobantes) == 0 {
		request_cliente_id := strconv.FormatInt(RRComprobante.Cliente_id, 10)
		request_comprobante_id := strconv.FormatInt(RRComprobante.ComprobanteId, 10)
		message := fmt.Sprintf("no existen comprobantes generados para el reporte. Id Cliente: %s, Id Comprobante: %s", request_cliente_id, request_comprobante_id)
		erro = errors.New(message)
		s.util.BuildLog(erro, "SendReporteRetencionComprobante")
		return
	}

	// para cada registro de comprobantes, crear un pdf reporte y enviar un email
	for _, comprobante := range listaComprobantes {
		// declaracion de variables
		var (
			tipo_archivo, contentType, asunto, nombreArchivo string
		)

		filtro_gravamen := retenciondtos.GravamenRequestDTO{
			Gravamen: comprobante.Gravamen,
		}
		// obtener el gravamen del comprobante para  saber el codigo de impuesto
		gravamenes, err := s.administracion.GetGravamenesService(filtro_gravamen)
		if err != nil {
			erro = err
			s.util.BuildLog(erro, "SendReporteRetencionComprobante")
			return
		}

		if len(gravamenes) < 1 {
			erro = errors.New("no se obtuvieron gravamenes en SendReporteRetencionComprobante")
			s.util.BuildLog(erro, "SendReporteRetencionComprobante")
			return
		}

		// el numero de reporte de tipo rrm correspondiente formateado en string
		nro_reporte_rrm := s.util.GenerarNumeroComprobante2(comprobante.ReporteId)
		// crear un report en pdf con la informacion del comprobante.
		// Segundo parametro: el codigo de impuesto
		// Tercer parametro: el numero de reporte de tipo rrm correspondiente formateado en string
		// El cuarto parametro es la fecha del fin del periodo de liquidacion
		fechaFinPeriodo := s.commons.ConvertirFechaToDDMMYYYY(RRComprobante.FechaFin[:10])
		nombreArchivo, err = RRComprobantePDF(comprobante, gravamenes[0].CodigoGravamen, nro_reporte_rrm, fechaFinPeriodo)
		if err != nil {
			erro = err
			s.util.BuildLog(erro, "SendReporteRetencionComprobante")
			return
		}

		// ENVIAR COMPROBANTE POR EMAIL. Esta seccion de codigo solo se realiza si el total agrupado por impuesto de un comprobante es mayor al minimo imponible configurado para ese gravamen
		corresponde_retener, err := s.EvaluarMinimoRetencionDeComprobante(comprobante)
		if err != nil {
			erro = err
			s.util.BuildLog(erro, "SendReporteRetencionComprobante")
			return
		}
		var emails []string
		tipo_archivo = ".pdf"
		contentType = "application/pdf"
		cliente := comprobante.Cliente
		var mensaje string

		if !corresponde_retener {
			emailsControl := []string{"pablo.vicentin@telco.com.ar", "yasmir.yaya@telco.com.ar", "sebastian.escobar@telco.com.ar"}

			asunto = "DEVOLUCION - Comprobante de Retenciones"
			mensaje = "Corresponde devolver los montos de retenciones detallados en el comprobante adjunto"
			for _, email := range emailsControl {
				emails = append(emails, email)
			}
		}

		// Envio de email sujeto a si corresponde retener en este comprobante
		if corresponde_retener {
			asunto = "Comprobante de Retenciones"
			mensaje = "Comprobante de retencion en PDF"
			emailsControl := []string{"pablo.vicentin@telco.com.ar", "yasmir.yaya@telco.com.ar", "sebastian.escobar@telco.com.ar"}

			// for _, email := range *cliente.Contactosreportes {
			// 	emails = append(emails, email.Email)
			// }
			for _, email := range emailsControl {
				emails = append(emails, email)
			}
		}

		filtro := utildtos.RequestDatosMail{
			Email:   emails,
			Asunto:  asunto,
			From:    "Wee.ar!",
			Nombre:  cliente.Cliente,
			Mensaje: mensaje,
			Descripcion: utildtos.DescripcionTemplate{
				Fecha:   time.Now().Format("2006-01-02"),
				Cliente: cliente.Razonsocial,
				Cuit:    cliente.Cuit,
			},
			Attachment: utildtos.Attachment{
				Name:        fmt.Sprintf("%s%s", nombreArchivo, tipo_archivo),
				ContentType: contentType,
				WithFile:    true,
			},
			TipoEmail:   "adjunto",
			RutaArchivo: (config.DIR_BASE + config.DIR_COMP_RETENCIONES),
		}

		// enviar archivo por correo
		erro = s.util.EnviarMailService(filtro)
		if erro != nil {
			s.util.BuildLog(erro, "SendReporteRetencionComprobante")
			return
		}

		nombreFile := administraciondtos.ArchivoResponse{
			NombreArchivo: (nombreArchivo + tipo_archivo),
		}
		nombreFiles := []administraciondtos.ArchivoResponse{nombreFile}

		_, erro = s.administracion.SubirArchivosCloud(context.Background(), (config.DIR_BASE + config.DIR_COMP_RETENCIONES), nombreFiles, (config.DIR_BASE + "/retenciones/comprobantes/" + cliente.Cliente))
		if erro != nil {
			s.util.BuildLog(erro, "SendReporteRetencionComprobante")
			return
		}

		// actualizar el campo emitido_el para marcar el comprobante como reportado y enviado
		if comprobante.EmitidoEl.IsZero() {
			// guardar la ruta donde se guarda en el cloud storage
			comprobante.RutaFile = CreateRutaFile(comprobante.Cliente.Cliente, nombreArchivo, tipo_archivo)
			comprobante.EmitidoEl, erro = time.Parse(time.RFC3339, RRComprobante.FechaFin)
			_, erro = s.repository.UpdateComprobanteRepository(comprobante)
			if erro != nil {
				s.util.BuildLog(erro, "SendReporteRetencionComprobante")
				return
			}
		}
	}

	control = true
	return
}

func (s *reportesService) LiquidarRetencionesService(request reportedtos.RequestReportesEnviados) (erro error) {

	filtroCliente := filtros.ClienteFiltro{
		Id: request.ClienteId,
		// SujetoRetencion: true,
		CargarCuentas: true,
	}

	clientes, _, erro := s.administracion.ObtenerClientesSinDTOService(filtroCliente)

	if erro != nil {
		s.util.BuildLog(erro, "LiquidarRetencionesService")
		return
	}

	for _, cliente := range clientes {
		if !cliente.SujetoRetencion {
			mensaje := fmt.Sprintf("el cliente %s no esta sujeto a retencion", cliente.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "LiquidarRetencionesService")
			return
		}

		fechaInicioFirtsMoment, _ := s.commons.DateYMDtoDateFirstMoment(request.FechaInicio)
		fechaFinLastMoment, _ := s.commons.DateYMDtoDateLastMoment(request.FechaFin)

		// Control para el caso de retenciones borradas y no coincidencia con reporte retenciones
		// request_control := retenciondtos.RentencionRequestDTO{
		// 	FechaInicio: fechaInicioFirtsMoment,
		// 	FechaFin:    fechaFinLastMoment,
		// 	ClienteId:   cliente.Id,
		// 	// NumeroReporteRrm: numero_reporte_rrm,
		// }
		// si el sujeto retencion no tiene movimientos-retenciones en el periodo, se continua con el proximo cliente
		// movimientos_ids, err := s.administracion.GetMovimientosRetencionesService(request_control)
		// if err != nil {
		// 	erro = err
		// 	s.util.BuildLog(erro, "LiquidarRetencionesService")
		// 	erro = nil // limpiar error para que no vuelva a la ruta con error
		// 	continue
		// }
		// if len(movimientos_ids) < 1 {
		// 	mensaje := fmt.Sprintf("error en LiquidarRetencionesService. el cliente %s no tiene retenciones en el periodo consultado", cliente.Cliente)
		// 	erro = errors.New(mensaje)
		// 	s.util.BuildLog(erro, "LiquidarRetencionesService")
		// 	erro = nil // limpiar error para que no vuelva a la ruta con error
		// 	continue
		// }

		request.Cliente = cliente.Cliente
		request.TipoReporte = "rendiciones"
		request.FechaInicio = fechaInicioFirtsMoment
		request.FechaFin = fechaFinLastMoment

		for _, cuenta := range *cliente.Cuentas {

			// PASO 1: se crea el reporte rendiciones mensual y se guarda en una carpeta temporal del proyecto. Se sube al cloud
			listaReportesCliente, err := s.GetReportesPdfService(request, cliente, cuenta)
			if err != nil {
				erro = err
				s.util.BuildLog(erro, "LiquidarRetencionesService")
				erro = nil // limpiar error para que no vuelva a la ruta con error
				continue
			}
			if len(listaReportesCliente) < 1 {
				mensaje := fmt.Sprintf("la lista de reportes de rendiciones de clientes esta vacía para el cliente %s", cliente.Cliente)
				s.util.BuildLog(errors.New(mensaje), "LiquidarRetencionesService")
				erro = nil // limpiar error para que no vuelva a la ruta con error
				continue
			}

			// PASO 2: enviar el reporte RRM de retenciones y guardar en la BD
			listaErro, numero_reporte_rrm, err := s.SendReporteRendiciones(listaReportesCliente)

			if err != nil {
				erro = err
				s.util.BuildLog(erro, "LiquidarRetencionesService")
				return
			}
			if len(listaErro) > 0 {
				s.util.BuildLog(erro, "LiquidarRetencionesService")
				return
			}

			// Si el cliente no es sujeto de retencion no se hace comprobante de retencion
			// No se puede generar comprobante de retencion sin movimientos
			if !cliente.SujetoRetencion {
				continue
			}

			// PASO 3: Generar Certificacion
			request_certificacion := retenciondtos.RentencionRequestDTO{
				FechaInicio:      fechaInicioFirtsMoment,
				FechaFin:         fechaFinLastMoment,
				ClienteId:        cliente.ID,
				NumeroReporteRrm: numero_reporte_rrm,
			}

			// la certificacion esta vinculada a un comprobante o reporte de rendiciones mensual (rrm)
			err = s.administracion.GenerarCertificacionService(request_certificacion)

			if err != nil {
				erro = err
				s.util.BuildLog(erro, "LiquidarRetencionesService")
				return
			}

			// PASO 4: enviar coprobante de retencion y subir cloud storage
			// filtro necesario para llamar a SendReporteRetencionComprobante
			request_comprobante_retencion := reportedtos.RequestRRComprobante{
				Cliente_id:  int64(cliente.ID),
				FechaInicio: fechaInicioFirtsMoment,
				FechaFin:    fechaFinLastMoment,
			}
			var successful bool
			// Generar comprobante y reporte pdf de retencion por gravamen. Envio de email y adjunto pdf
			successful, err = s.SendReporteRetencionComprobante(request_comprobante_retencion)

			if err != nil {
				erro = err
				s.util.BuildLog(erro, "LiquidarRetencionesService")
				return
			}

			if !successful {
				mensaje := fmt.Sprintf("error en SendReporteRetencionComprobante al generar comprobante de retencion para el clientes %s", cliente.Cliente)
				erro = errors.New(mensaje)
				s.util.BuildLog(erro, "LiquidarRetencionesService")
				return
			}

		} // Fin de for _,cuenta := range *cliente.Cuentas

	} // Fin for _, cliente := range clientes

	return
}

func (s *reportesService) CreateTxtRetencionesService(request reportedtos.RequestReportesEnviados) (erro error) {
	rutaSicoreTxt, erro := s.CreateTxtRetencionesSICOREService(request)
	if erro != nil {
		return
	}

	rutaSicarTxt, erro := s.CreateTxtRetencionesSICARService(request)
	if erro != nil {
		return
	}

	sicoreName := strings.Split(rutaSicoreTxt, "/")
	sicarName := strings.Split(rutaSicarTxt, "/")

	archivos := []string{sicoreName[len(sicoreName)-1], sicarName[len(sicarName)-1]}
	emails := []string{"sebastian.escobar@telco.com.ar", "pablo.vicentin@telco.com.ar", "yasmir.yaya@telco.com.ar"}

	if len(request.Emails) > 0 {
		emails = request.Emails
	}

	requestEmail := reportedtos.RequestRetencionEmail{
		Archivos: archivos,
		Emails:   emails,
	}

	fmt.Print(requestEmail)
	erro = s.SendRetencionestxt(requestEmail)
	if erro != nil {
		return
	}

	return
}

// Crear ruta de archivo txt SICORE en proyecto. Crear archivo txt y guardar
func (s *reportesService) CreateTxtSICORE(lines []string, nombreTXT string, fechaFin string) (rutaDetalle string, erro error) {
	// ruta del archivo en el proyecto "/documentos/retenciones"
	ruta := fmt.Sprintf(config.DIR_BASE + "/documentos/retenciones")
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	// diferenciar nombres de archivo
	if strings.Contains(strings.ToUpper(nombreTXT), "NOVEDADES") {
		rutaDetalle = ruta + "/" + nombreTXT + ".txt"
	} else {
		nombreArchivoTxtSICORE := commonsdtos.FileName{
			RutaBase:        ruta + "/",
			Nombre:          nombreTXT,
			Extension:       "txt",
			UsaFecha:        false,
			FechaEspecifica: fechaFin,
		}
		rutaDetalle = s.commons.CreateFileName(nombreArchivoTxtSICORE)
	}

	// CREAR ARCHIVO
	file, err := os.Create(rutaDetalle)
	if err != nil {
		erro = err
		return
	}

	// ESCRIBIR EN EL ARCHIVO
	err = s.WriteTxtSICORE(file, lines)
	if err != nil {
		erro = err
		return
	}

	// DIFERIR CIERRE DE ARCHIVO al finalizar func
	defer file.Close()

	// Leer el contenido del archivo
	data, err := ioutil.ReadFile(rutaDetalle)
	if err != nil {
		erro = err
		return
	}

	// SUBIR CLOUD STORAGE
	rutaCloud := "retenciones/comprobantes/" + nombreTXT
	// rutaDetalle: origen del archivo
	// rutaCloud: destino del archivo
	ctx := context.Background()

	err = s.UploadTxtFile(ctx, rutaDetalle, rutaCloud, data)
	if err != nil {
		erro = err
		s.util.BuildLog(erro, "CreateTxtSICORE")
		return
	}

	return
}

// escribir lineas de texto en archivo SICORE txt de retenciones
func (s *reportesService) WriteTxtSICORE(archivo *os.File, lines []string) (erro error) {

	// incluir un salto de linea al final de cada registro
	for i := 0; i < len(lines); i++ {
		lines[i] += "\n"
	}

	datos_escribir := commons.JoinString(lines)
	_, erro = archivo.Write([]byte(datos_escribir))
	if erro != nil {
		return erro
	}
	return
}

func (s *reportesService) UploadTxtFile(ctx context.Context, rutaArchivo, rutaCloud string, data []byte) (erro error) {

	fileType := strings.Replace(filepath.Ext(rutaArchivo), ".", "", -1)
	// Obtener el nombre del archivo sin la ruta
	nombreArchivo := filepath.Base(rutaArchivo)
	// Eliminar la extensión
	nombreSinExtension := strings.TrimSuffix(nombreArchivo, "."+fileType)
	filenameWithoutExt := rutaCloud + "/" + nombreSinExtension
	fmt.Print(filenameWithoutExt)
	// Example:
	// filenameWithoutExt: "retenciones/comprobantes/SICORE/SICORE_04102023200503"
	// filetype: “txt”
	err := s.store.PutObject(ctx, data, filenameWithoutExt, fileType)
	if erro != nil {
		erro = err
		return
	}

	return
}

// recibe un entities.Comprobante. Evalua si sus ComprobanteDetalle tiene el atributo retener en true
// Basta que un ComprobanteDetalle sea false para no retener
func (s *reportesService) EvaluarMinimoRetencionDeComprobante(comprobante entities.Comprobante) (result bool, erro error) {
	if len(comprobante.ComprobanteDetalles) == 0 {
		mensaje := fmt.Sprintf("el comprobante de id %d del cliente %s no posee detalles", comprobante.ID, comprobante.RazonSocial)
		erro = errors.New(mensaje)
		return
	}
	result = true // se supone que se debe retener a menos que se encuentre uno falso a continuacion
	for _, cd := range comprobante.ComprobanteDetalles {
		if !cd.Retener {
			result = cd.Retener
			break
		}
	}
	return
}

func (s *reportesService) CreateTxtRetencionesSICARService(request reportedtos.RequestReportesEnviados) (rutatxt string, erro error) {

	// Paso 1: Buscar los movimientos con retención asociados entre las fechas
	fechaInicio, _ := time.Parse("2006-01-02", request.FechaInicio)
	fechaFin, _ := time.Parse("2006-01-02", request.FechaFin)
	fechaFin = commons.GetDateLastMomentTime(fechaFin) // que sea el ultimo momento de la fecha

	filtroCliente := filtros.ClienteFiltro{
		Id:              request.ClienteId,
		SujetoRetencion: true,
	}

	clientes, _, erro := s.administracion.ObtenerClientesSinDTOService(filtroCliente)
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICARService")
		return
	}

	var idsClientesSicar []uint
	// filtrar clientes que tienen al menos uno de los gravamenes de este reporte
	for _, cte := range clientes {
		res := cte.HasRetencionByGravamenName([]string{"iibb"})
		if res {
			idsClientesSicar = append(idsClientesSicar, cte.ID)
		}
	}
	var movimientos_ids_por_cliente []uint
	for _, id := range idsClientesSicar {
		request_movs_retencion := retenciondtos.RentencionRequestDTO{
			FechaInicio: fechaInicio.Format("2006-01-02 15:04:05"),
			FechaFin:    fechaFin.Format("2006-01-02 15:04:05"),
			ClienteId:   id,
		}
		movimientos_ids_por_cliente_temp, _ := s.administracion.GetMovimientosIdsCalculoRetencionComprobanteService(request_movs_retencion)
		movimientos_ids_por_cliente = append(movimientos_ids_por_cliente, movimientos_ids_por_cliente_temp...)
	}

	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICARService")
		return
	}
	filtroMovimientos := reportedtos.RequestPagosPeriodo{
		IdsMovimientos:                 movimientos_ids_por_cliente,
		CargarCliente:                  true,
		CaragarSoloMovimientoRetencion: true,
		CargarPagoIntentos:             true,
	}

	// TRAER SOLO MOVIMIENTOS QUE SEA HAYAN INFORMADO EN COMPROBANTE DE RETENCION
	movimientos, erro := s.repository.GetMovimientoByIds(filtroMovimientos)

	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICARService")
		return
	}
	if len(movimientos) < 1 {
		mensaje := fmt.Sprintf("error en CreateTxtRetencionesSICARService. No existen movimientos en el periodo requerido %s a %s", request.FechaInicio, request.FechaFin)
		erro = errors.New(mensaje)
		s.util.BuildLog(erro, "CreateTxtRetencionesSICARService")
		return
	}

	// Paso 2: Formar las líneas de objetos LineData para cada movimiento con retencion SICAR
	// Datos fijos de LineData para SICAR
	lineData := LineDataIIBB{
		OrigenComprobante: "1",
		TipoComprobante:   "1",
		NroComprobante:    "000000000000", // 12 caracteres
		TipoRegimen:       "02",
		Jurisdiccion:      "905",
	}

	// obtener los gravamenes para saber su codigo
	gravamenes, erro := s.administracion.GetGravamenesService(retenciondtos.GravamenRequestDTO{})
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICARService")
		return
	}
	// separar en un map por nombre de gravamen
	var map_gravamenes = make(map[string]retenciondtos.GravamenResponseDTO)
	for _, g := range gravamenes {
		map_gravamenes[g.Gravamen] = g
	}

	// var lines guarda cada una de las lineas de retencion
	var lines []string
	renglon := 0 // contador de lineas
	idGravamenIIBB := map_gravamenes["iibb"].Id

	// crear una linea del archivo txt
	// para cada movimiento
	for _, movimiento := range movimientos {
		// que el movimiento tenga MovimientoRetencion
		if len(movimiento.Movimientoretencions) > 0 {
			// que el movimiento tenga MovimientoRetencion cuya Retencion sea de tipo IIBB
			mov_ret_iibb, tiene_ret_iibb := entities.MovimientosRetenciones(movimiento.Movimientoretencions).GetByGravamenId(idGravamenIIBB)

			if tiene_ret_iibb {
				renglon++
				renglonStr := strconv.Itoa(renglon)
				totalPago := strconv.FormatFloat(mov_ret_iibb.Monto.Float64(), 'f', 2, 64)
				// fechaPago := movimiento.Pagointentos.PaidAt.Format("02/01/2006")
				fechaPago := fechaFin.Format("02/01/2006") // se pone la fecha de fin del periodo
				CuitContribuyente := movimiento.Cuenta.Cliente.Cuit
				alicuota := strconv.FormatFloat(mov_ret_iibb.Retencion.Alicuota, 'f', 2, 64)
				montoRetenido := strconv.FormatFloat(mov_ret_iibb.ImporteRetenido.Float64(), 'f', 2, 64)

				line := NewLineBuilder().
					SetValueString(renglonStr, 5).SetComma(1).
					SetValueString(lineData.OrigenComprobante, 1).SetComma(1).
					SetValueString(lineData.TipoComprobante, 1).SetComma(1).
					SetValueString(lineData.NroComprobante, 12).SetComma(1).
					SetString(CuitContribuyente).SetComma(1).
					SetString(fechaPago).SetComma(1).
					SetString(totalPago).SetComma(1).
					SetString(alicuota).SetComma(1).
					SetString(montoRetenido).SetComma(1).
					SetValueString(lineData.TipoRegimen, 3).SetComma(1).
					SetString(lineData.Jurisdiccion).
					Build()

				lines = append(lines, line)
			}
		}
	} // Fin for _, movimiento := range movimientos

	// Paso 5: Crear archivo txt, escribir datos y subir a cloud storage
	fecha, erro := s.commons.DateTimeToYYYYMM(fechaInicio, fechaFin)
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICARService")
		return
	}
	nombre_txt_rentas := "CTES_RG 202_Novedades" + "_" + fecha
	rutaDetalle, erro := s.CreateTxtSICORE(lines, nombre_txt_rentas, "")
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICARService")
		return
	}

	rutatxt = rutaDetalle

	return
}

func (s *reportesService) CreateTxtRetencionesSICOREService(request reportedtos.RequestReportesEnviados) (rutatxt string, erro error) {
	// Var para contener lineas de archivo retencion txt
	var lines []string

	// Paso 1: buscar clientes sujeto de retencion
	filtroCliente := filtros.ClienteFiltro{
		Id:              request.ClienteId,
		SujetoRetencion: true,
	}

	clientes, _, erro := s.administracion.ObtenerClientesSinDTOService(filtroCliente)
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
		return
	}

	var clientesSicore []entities.Cliente
	// filtrar clientes que tienen al menos uno de los gravamenes de este reporte
	for _, cte := range clientes {
		res := cte.HasRetencionByGravamenName([]string{"iva", "ganancias"})
		if res {
			clientesSicore = append(clientesSicore, cte)
		}
	}

	// obtener los gravamenes para saber su codigo
	gravamenes, err := s.administracion.GetGravamenesService(retenciondtos.GravamenRequestDTO{})
	if err != nil {
		erro = err
		s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
		return
	}
	// separar en un map por nombre de gravamen
	var map_gravamenes = make(map[string]retenciondtos.GravamenResponseDTO)
	for _, g := range gravamenes {
		map_gravamenes[g.Gravamen] = g
	}

	// PARA CADA CLIENTE SUJETO DE RETENCION
	for _, cliente := range clientesSicore {
		if !cliente.SujetoRetencion {
			mensaje := fmt.Sprintf("el cliente %s no esta sujeto a retencion", cliente.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "CreateTxtRetencionesService")
			return
		}

		// Paso 2: Buscar reporte rrm para el sujeto de retención
		var sujeto_retencion = cliente
		fechaInicio, _ := s.commons.DateStringToTimeFirstMoment(request.FechaInicio)
		fechaFin, _ := s.commons.DateStringToTimeLastMoment(request.FechaFin)
		requestGetReportes := reportedtos.RequestGetReportes{
			TipoReporte:   "rrm",
			Cliente:       sujeto_retencion.Cliente,
			PeriodoInicio: fechaInicio,
			PeriodoFin:    fechaFin,
			LastRrm:       true,
		}

		reportes_rrm, _, err := s.repository.GetReportesRepository(requestGetReportes)

		if len(reportes_rrm) == 0 {
			mensaje := fmt.Sprintf("no se encontro reporte mensual de rendiciones para el cliente %s ", sujeto_retencion.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionRetenciones,
				Descripcion: mensaje,
			}
			s.administracion.CreateNotificacionService(notificacion)
			continue
		}

		if err != nil {
			erro = err
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			return
		}
		var reporte_rrm = reportes_rrm[0]

		// Paso 3: Buscar los comprobantes de retención asociados al reporte rrm
		requestRRComprobante := reportedtos.RequestRRComprobante{
			ReporteId:    reporte_rrm.Nro_reporte,
			Cliente_id:   int64(sujeto_retencion.ID),
			GravamenesIn: []string{"ganancias", "iva"}, // retenciones AFIP
		}

		// para SICORE se buscan solo los comprobantes de retenciones AFIP
		comprobantes, err := s.repository.GetComprobantesByRrmIdRepository(requestRRComprobante)
		if err != nil {
			erro = err
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			return
		}

		if len(comprobantes) == 0 {
			mensaje := fmt.Sprintf("no se encontraron comprobantes de rendiciones para el cliente %s ", sujeto_retencion.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			continue
		}

		// CONTROL DE MINIMO. Si el cliente no supera el minimo no debe crearse lineas de datos en el reporte SICORE
		comprobantesSlice := entities.Comprobantes(comprobantes)
		comprobantesARetener := comprobantesSlice.GetComprobantesARetener()

		// si en algun caso se dio no retener, no se retiene
		if len(comprobantesARetener) < 1 {
			mensaje := fmt.Sprintf("no corresponde txt SICORE para el cliente %s ", sujeto_retencion.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			continue
		}

		// Paso 4: Formar las líneas de objetos LineData para un sujeto de retención
		for _, comprobante := range comprobantesARetener {
			for _, cd := range comprobante.ComprobanteDetalles {
				lineData := LineData{
					Codigo:                     "06", // fijo
					FechaRrm:                   strings.Replace(reporte_rrm.Fecharendicion, "-", "/", -1),
					NumeroRrm:                  s.util.GenerarNumeroComprobante2(reporte_rrm.Nro_reporte),
					ImporteRrm:                 s.util.QuitarPuntoNumeroString(reporte_rrm.Totalcobrado),
					CodigoGravamen:             map_gravamenes[cd.Gravamen].CodigoGravamen,
					CodigoRegimen:              cd.CodigoRegimen,
					CodigoEsRetencion:          "1", // fijo
					ImporteComprobanteCabecera: entities.Monto(comprobante.Importe).Float64(),
					FechaComprobanteCabecera:   s.util.DateTimeToStringFormatDMY(comprobante.EmitidoEl),
					CodigoCondicion:            "01", // fijo
					TotalRetenido:              cd.TotalRetencion.Float64(),
					PorcentajeExclusion:        "000,00", // fijo
					TipoDocumento:              "80",     // fijo
					CuitCliente:                sujeto_retencion.Cuit,
					NumeroCertificadoOriginal:  "00000000000000",
				}

				// crear linea con Builder
				line := CreateLine(lineData)
				// Guardar en slice de string
				lines = append(lines, line) // var lines definida en scope global de la funcion CreateTxtRetencionesService
			}
		}

	} // Fin de for _, cliente := range clientes

	// Paso 5: Crear archivo txt, escribir datos y subir a cloud storage
	rutaDetalle, erro := s.CreateTxtSICORE(lines, "SICORE", request.FechaFin)
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
		return
	}

	rutatxt = rutaDetalle

	return
}

func (s *reportesService) CreateTxtForm8125Service(request reportedtos.RequestReportesEnviados) (rutatxt string, erro error) {

	//

	fechaInicio, _ := time.Parse("2006-01-02", request.FechaInicio)
	fechaFin, _ := time.Parse("2006-01-02", request.FechaFin)
	fechaFin = commons.GetDateLastMomentTime(fechaFin) // que sea el ultimo momento de la fecha

	// var lines guarda cada una de las lineas del txt
	var lines []string
	var linesBody []string

	renglon := 0 // contador de lineas

	// Paso 1 Buscar los clientes a los cuales se necesita informar a la AFIP con el formulario 8125

	filtroCliente := filtros.ClienteFiltro{
		Formulario8125: true,
	}

	respuestaClientes, erro := s.administracion.GetClientesService(filtroCliente)
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtForm8125Service")
		return
	}

	for _, cliente := range respuestaClientes.Clientes {

		MontoCliente := int64(0)
		ComisionCliente := int64(0)

		var linesDetalles []string

		// Paso 2: Buscar los movimientos con comisiones entre las fechas
		filtroMovimientos := reportedtos.RequestPagosPeriodo{
			FechaInicio:            fechaInicio,
			FechaFin:               fechaFin,
			ClienteId:              uint64(cliente.Id),
			CargarCliente:          true,
			CargarComisionImpuesto: true,
		}

		movimientos, err := s.repository.GetMovimiento(filtroMovimientos)
		if err != nil {
			erro = err
			s.util.BuildLog(erro, "CreateTxtForm8125Service")
			return
		}

		// crear una linea detalle en archivo txt
		// para cada movimiento
		for _, movimiento := range movimientos {
			// que el movimiento tenga MovimientoComisions
			if len(movimiento.Movimientocomisions) > 0 {

				renglon++
				montoTruncado := math.Trunc(movimiento.Monto.Float64())

				comisionMov := entities.Monto(movimiento.Movimientocomisions[len(movimiento.Movimientocomisions)-1].Monto) + entities.Monto(movimiento.Movimientocomisions[len(movimiento.Movimientocomisions)-1].Montoproveedor)
				comisionTruncada := math.Trunc(comisionMov.Float64())

				signoMovimiento := 0 // Positivo
				if montoTruncado < 0 {
					signoMovimiento = 1 // Negativo
				}
				montoMovimientoFormateado := int64(math.Abs(montoTruncado))

				MontoCliente += int64(montoTruncado)
				ComisionCliente += int64(comisionTruncada)

				detalleData := DetalleData{
					TipoRegistro:            "03",
					MetodologiaAcreditacion: "01",
					TipoCuenta:              "13",
					NumeroIdentificacion:    movimiento.Cuenta.Cbu,
					SignoMonto:              fmt.Sprint(signoMovimiento),
					Monto:                   fmt.Sprint(montoMovimientoFormateado),
				}

				line := NewLineBuilder().
					SetValueString(detalleData.TipoRegistro, 2).
					SetValueString(detalleData.MetodologiaAcreditacion, 2).
					SetValueString(detalleData.TipoCuenta, 2).
					SetValueString(detalleData.NumeroIdentificacion, 22).
					SetValueString(detalleData.SignoMonto, 1).
					SetValueString(detalleData.Monto, 12).
					Build()

				linesDetalles = append(linesDetalles, line)

			}
		} // Fin for _, movimiento := range movimientos

		signo := 0 // Positivo
		if MontoCliente < 0 {
			signo = 1 // Negativo
		}

		renglon++
		vendedorData := VendedorData{
			TipoRegistro:           "02",
			TipoIdentificacion:     "80",
			IdentificacionVendedor: cliente.Cuit,
			CodigoRubro:            "07",
			SignoTotal:             fmt.Sprint(signo),
			MontoTotal:             fmt.Sprint(MontoCliente),
			ImporteComision:        fmt.Sprint(ComisionCliente),
		}

		line := NewLineBuilder().
			SetValueString(vendedorData.TipoRegistro, 2).
			SetValueString(vendedorData.TipoIdentificacion, 2).
			SetValueString(vendedorData.IdentificacionVendedor, 11).
			SetValueString(vendedorData.CodigoRubro, 2).
			SetValueString(vendedorData.SignoTotal, 1).
			SetValueString(vendedorData.MontoTotal, 12).
			SetValueString(vendedorData.ImporteComision, 12).
			Build()

		linesBody = append(linesBody, line)
		linesBody = append(linesBody, linesDetalles...)
	} // Fin for _, cliente := range respuestaClientes.Clientes

	fechaFormateada := fechaFin.Format("200601")
	reporteBusqueda := entities.Reporte{
		Tiporeporte: "formulario8125",
	}
	opcionesBusqueda := filtros_reportes.BusquedaReporteFiltro{
		SigNumero: true,
	}

	nroReporte, erro := s.repository.GetLastReporteEnviadosRepository(reporteBusqueda, opcionesBusqueda)
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtForm8125Service")
		return
	}

	renglon++
	cabeceraData := CabeceraData{
		TipoRegistro:      "01",
		CuitInformante:    "30716550849",
		PeriodoInformado:  fechaFormateada,
		Secuencia:         "00",
		Denominacion:      "CORRIENTES TELECOMUNICACIONES SAPEM",
		Hora:              "000000",
		CodigoImpuesto:    "0103",
		CodigoConcepto:    "830",
		NumeroVerificador: fmt.Sprint(nroReporte),
		NumeroFormulario:  "8125",
		NumeroVersion:     "00100",
		Establecimiento:   "00",
		CantidadRegistros: fmt.Sprint(renglon),
	}

	line := NewLineBuilder().
		SetValueString(cabeceraData.TipoRegistro, 2).
		SetValueString(cabeceraData.CuitInformante, 11).
		SetStringSpaced(cabeceraData.PeriodoInformado, 6).
		SetValueString(cabeceraData.Secuencia, 2).
		SetStringSpaced(cabeceraData.Denominacion, 200).
		SetValueString(cabeceraData.Hora, 6).
		SetValueString(cabeceraData.CodigoImpuesto, 4).
		SetValueString(cabeceraData.CodigoConcepto, 3).
		SetValueString(cabeceraData.NumeroVerificador, 6).
		SetValueString(cabeceraData.NumeroFormulario, 4).
		SetValueString(cabeceraData.NumeroVersion, 5).
		SetValueString(cabeceraData.Establecimiento, 2).
		SetValueString(cabeceraData.CantidadRegistros, 10).
		Build()

	lines = append(lines, line)

	lines = append(lines, linesBody...)

	// Paso 4: Crear archivo txt, escribir datos y subir a cloud storage
	rutaDetalle, erro := s.CreateTxtSICORE(lines, "Form8125", "")
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtForm8125Service")
		return
	}

	rutatxt = rutaDetalle

	erro = os.Remove(rutaDetalle)
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICARService")
		return
	}
	return
}

// para un slice de entities.MovimientoRetencion devuelve el importe de la retencion por el nombre del gravamen correspondiente
func importeRetencionByName(gravamen_name string, mr []entities.MovimientoRetencion) (importe entities.Monto) {
	for _, item := range mr {
		if gravamen_name == item.Retencion.Condicion.Gravamen.Gravamen {
			importe = item.ImporteRetenido
			break
		}
	}
	return
}

// genera la ruta que se guarda en la tabla de reportes y comprobantes de retenciones
// la extension se puede recibir con o sin punto
func CreateRutaFile(carpeta, filename, extension string) (rutaFile string) {
	// Eliminar el punto si la extensión lo tiene
	extension = strings.TrimPrefix(extension, ".")
	rutaFile = fmt.Sprintf("/%s/%s.%s", carpeta, filename, extension)
	return
}

func (s *reportesService) CreateExcelRetencionesService(request reportedtos.RequestReportesEnviados) (rutatxt string, erro error) {
	// Var para contener lineas de archivo retencion txt
	var lines []string
	var linesExcel []LineDataExcel

	// Paso 1: buscar clientes sujeto de retencion
	filtroCliente := filtros.ClienteFiltro{
		Id:              request.ClienteId,
		SujetoRetencion: true,
	}

	clientes, _, erro := s.administracion.ObtenerClientesSinDTOService(filtroCliente)
	if erro != nil {
		s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
		return
	}

	var clientesSicore []entities.Cliente
	// filtrar clientes que tienen al menos uno de los gravamenes de este reporte
	for _, cte := range clientes {
		res := cte.HasRetencionByGravamenName([]string{"iva", "ganancias", "iibb"})
		if res {
			clientesSicore = append(clientesSicore, cte)
		}
	}

	// obtener los gravamenes para saber su codigo
	gravamenes, err := s.administracion.GetGravamenesService(retenciondtos.GravamenRequestDTO{})
	if err != nil {
		erro = err
		s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
		return
	}
	// separar en un map por nombre de gravamen
	var map_gravamenes = make(map[string]retenciondtos.GravamenResponseDTO)
	for _, g := range gravamenes {
		map_gravamenes[g.Gravamen] = g
	}

	// PARA CADA CLIENTE SUJETO DE RETENCION
	for _, cliente := range clientesSicore {
		if !cliente.SujetoRetencion {
			mensaje := fmt.Sprintf("el cliente %s no esta sujeto a retencion", cliente.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "CreateTxtRetencionesService")
			continue
		}

		// Paso 2: Buscar reporte rrm para el sujeto de retención
		var sujeto_retencion = cliente
		fechaInicio, _ := s.commons.DateStringToTimeFirstMoment(request.FechaInicio)
		fechaFin, _ := s.commons.DateStringToTimeLastMoment(request.FechaFin)
		requestGetReportes := reportedtos.RequestGetReportes{
			TipoReporte:   "rrm",
			Cliente:       sujeto_retencion.Cliente,
			PeriodoInicio: fechaInicio,
			PeriodoFin:    fechaFin,
		}

		reportes_rrm, _, err := s.repository.GetReportesRepository(requestGetReportes)

		if len(reportes_rrm) == 0 {
			mensaje := fmt.Sprintf("no se encontro reporte mensual de rendiciones para el cliente %s ", sujeto_retencion.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionRetenciones,
				Descripcion: mensaje,
			}
			s.administracion.CreateNotificacionService(notificacion)
			continue
		}

		if err != nil {
			erro = err
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			continue
		}
		var reporte_rrm entities.Reporte
		var comprobantes []entities.Comprobante

		reporte_rrm = reportes_rrm[0]

		// Paso 3: Buscar los comprobantes de retención asociados al reporte rrm
		requestRRComprobante := reportedtos.RequestRRComprobante{
			ReporteId:    reporte_rrm.Nro_reporte,
			Cliente_id:   int64(sujeto_retencion.ID),
			GravamenesIn: []string{"ganancias", "iva", "iibb"}, // retenciones AFIP
		}

		// para SICORE se buscan solo los comprobantes de retenciones AFIP
		comprobantes, err = s.repository.GetComprobantesByRrmIdRepository(requestRRComprobante)
		if err != nil {
			erro = err
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			continue
		}

		if len(comprobantes) == 0 {
			mensaje := fmt.Sprintf("no se encontraron comprobantes de rendiciones para el cliente %s ", sujeto_retencion.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			continue
		}

		// CONTROL DE MINIMO. Si el cliente no supera el minimo no debe cerarse lineas de datos en el reporte SICORE
		// TODO: caso ambiguo en el que un impuesto supera el minimo y otro no
		var noRetener bool
		for _, comprobante := range comprobantes {
			for _, cd := range comprobante.ComprobanteDetalles {
				if !cd.Retener {
					noRetener = true
				}
			}
		}

		// si en algun caso se dio no retener, no se retiene
		if noRetener {
			mensaje := fmt.Sprintf("no corresponde txt SICORE para el cliente %s ", sujeto_retencion.Cliente)
			erro = errors.New(mensaje)
			s.util.BuildLog(erro, "CreateTxtRetencionesSICOREService")
			continue
		}

		// Paso 4: Formar las líneas de objetos LineData para un sujeto de retención
		for _, comprobante := range comprobantes {
			for _, cd := range comprobante.ComprobanteDetalles {
				lineData := LineData{
					Codigo:                     "06", // fijo
					FechaRrm:                   strings.Replace(reporte_rrm.Fecharendicion, "-", "/", -1),
					NumeroRrm:                  s.util.GenerarNumeroComprobante2(reporte_rrm.Nro_reporte),
					ImporteRrm:                 s.util.QuitarPuntoNumeroString(reporte_rrm.Totalcobrado),
					CodigoGravamen:             map_gravamenes[cd.Gravamen].CodigoGravamen,
					CodigoRegimen:              cd.CodigoRegimen,
					CodigoEsRetencion:          "1", // fijo
					ImporteComprobanteCabecera: entities.Monto(comprobante.Importe).Float64(),
					FechaComprobanteCabecera:   s.util.DateTimeToStringFormatDMY(comprobante.EmitidoEl),
					CodigoCondicion:            "01", // fijo
					TotalRetenido:              cd.TotalRetencion.Float64(),
					PorcentajeExclusion:        "000,00", // fijo
					TipoDocumento:              "80",     // fijo
					CuitCliente:                sujeto_retencion.Cuit,
					NumeroCertificadoOriginal:  comprobante.Numero,
				}

				// crear linea con Builder
				line := CreateLine(lineData)
				// Guardar en slice de string
				lines = append(lines, line) // var lines definida en scope global de la funcion CreateTxtRetencionesService
			}
		}

		request_control := retenciondtos.RentencionRequestDTO{
			FechaInicio: request.FechaInicio,
			FechaFin:    request.FechaFin,
			ClienteId:   cliente.ID,
		}
		request_control.ValidarFechas()

		movimientos_ids, err := s.administracion.CalcularRetencionesByTransferenciasSinAgruparService(request_control)
		if err != nil {
			// s.utilService.BuildLog(erro, "GenerarCertificacionService. Error al obtener ids de movimientos con retenciones")
			return
		}

		for _, mov := range movimientos_ids {
			lineExcel := LineDataExcel{
				NOM_PROV:   sujeto_retencion.Cliente,
				IDENTIFTRI: sujeto_retencion.Cuit,
				FEC_RET:    strings.Replace(comprobantes[len(comprobantes)-1].EmitidoEl.Format("2006-01-02"), "-", "/", -1),
				COD_RET:    "01",
				N_COMP:     s.util.GenerarNumeroComprobante2(reporte_rrm.Nro_reporte),
				N_CERTIFIC: _getNumeroComprobanteByImpuesto(comprobantes, mov.Gravamen),
				IMP_PAGO:   strconv.FormatFloat(mov.TotalMonto.Float64(), 'f', 2, 64),
				IMP_RETEN:  strconv.FormatFloat(mov.TotalRetencion.Float64(), 'f', 2, 64),
			}
			linesExcel = append(linesExcel, lineExcel)

		}

	} // Fin de for _, cliente := range clientes

	ruta := fmt.Sprintf(config.DIR_BASE + "/documentos/retenciones")
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		err = os.MkdirAll(ruta, 0755)
		if err != nil {
			erro = err
			return
		}
	}
	err = ReporteExcel(ruta, "ExcelRetenciones", linesExcel)
	if err != nil {
		s.util.BuildLog(err, "CreateExcelRetenciones")
		erro = err
		return
	}

	erro = nil
	return
}

func ReporteExcel(ruta string, nombreArchivo string, lines []LineDataExcel) (erro error) {

	if len(lines) > 0 {
		RutaFile := fmt.Sprintf("%s/%s.csv", ruta, nombreArchivo)
		/* estos datos son los que se van a escribir en el archivo */
		var slice_array = [][]string{
			{"COD_PROV", "NOM_PROV", "IDENTIFTRI", "INGR_PROV", "FEC_RET", "COD_RET", "T_COMP", "N_COMP", "N_CERTIFIC", "IMP_PAGO", "IMP_RETEN", "COD_PROVE", "T_RETEN"}, // columnas
		}
		for _, line := range lines {
			slice_array = append(slice_array, []string{line.COD_PROV, line.NOM_PROV, line.IDENTIFTRI, line.INGR_PROV, line.FEC_RET, line.COD_RET, line.T_COMP, line.N_COMP, line.N_CERTIFIC, line.IMP_PAGO, line.IMP_RETEN, line.COD_PROVE, line.T_RETEN})
		}
		erro = CsvCreate(RutaFile, slice_array)
		if erro != nil {
			return
		}

	}
	return
}

// convertir datos  a excel // utilizado para reportes excel
func CsvCreate(name string, data [][]string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = ';'
	defer w.Flush()

	for _, d := range data {
		err := w.Write(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func _getNumeroComprobanteByImpuesto(comprobantes []entities.Comprobante, impuesto string) (numero string) {
	for _, c := range comprobantes {
		if c.Gravamen == impuesto {
			numero = c.Numero
			return
		}
	}
	return
}
