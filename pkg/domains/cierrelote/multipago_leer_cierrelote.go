package cierrelote

import (
	"bufio"
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/cierrelotemultipagosdtos"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type multipagoProcesarArchivos struct {
	utilService    util.UtilService
	administracion administracion.Service
}

func NewMPProcesarArchivo(util util.UtilService, administracion administracion.Service) MetodoProcesarArtchivos {
	return &multipagoProcesarArchivos{
		utilService:    util,
		administracion: administracion,
	}
}

// rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *multipagoProcesarArchivos) ProcesarArchivos(archivo *os.File, estadosPagoExterno []entities.Pagoestadoexterno, impuesto administraciondtos.ResponseImpuesto, clRepository Repository) (listaLogArchivo cierrelotedtos.PrismaLogArchivoResponse) {
	// logs.Info(estadosPagoExterno)

	var estado = true
	var estadoInsert = true
	var ErrorProducido string
	rutaArchivo := strings.Split(archivo.Name(), "/")
	/* recorrer y validar archivo headers, trailer , detalles*/

	// 1 verifico que el archivo no se encuentre en la base de datos
	nombreArchivo := strings.Split(rutaArchivo[len(rutaArchivo)-1], "-")

	// Control lote ya procesado pendiente
	archivoExiste, err := cl.administracion.ObtenerArchivoCierreLoteMultipagos(nombreArchivo[len(nombreArchivo)-1])
	if err != nil {
		estado = false
		estadoInsert = false
		ErrorProducido = ERROR_OBTENER_NOMBRE_ARCHIVOS + err.Error()
		logs.Error(ErrorProducido)
		logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
		logs.Info("ObtenerArchivoCierreLoteMultipagos")
	} else if archivoExiste {
		estado = false
		estadoInsert = false
		ErrorProducido = ERROR_ARCHIVO_REPETIDO
		logs.Error(ErrorProducido)
		logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
		logs.Info("archivoExiste")
	} else {
		// Cambio pagoRP A pagoMP
		pagosMP, err := RecorrerArchivoMP(archivo)
		if err != nil {
			estado = false
			estadoInsert = false
			ErrorProducido = ERROR_RECORRER_ARCHIVOS + err.Error()
			logs.Error(ErrorProducido)
			logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
			logs.Info("RecorrerArchivoMP verificar")
		} else {

			if len(pagosMP) > 0 {

				pagoMp, err := GenerarListaMP(cl, impuesto, nombreArchivo, pagosMP)
				if err != nil {
					estado = false
					estadoInsert = false
					ErrorProducido = ERROR_RECORRER_ARCHIVOS + err.Error()
					logs.Error(ErrorProducido)
					logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
					logs.Info("GenerarListaMP verificar")
				} else {
					err := clRepository.SaveTransactionPagoMP(pagoMp)
					if err != nil {
						estadoInsert = false
						ErrorProducido = ERROR_REGISTRO_EN_DB + err.Error()
						logs.Error(ErrorProducido)
						logs.Info("error en el archivo: " + nombreArchivo[len(nombreArchivo)-1])
						logs.Info("clRepository.SaveTransactionPagoMP verificar")
					}
					//  else {
					// NOTE en el caso de guardar los registros del archivo se calcula las comisiones temporales de esos pagos
					logs.Info("inicio proceso calculo de comisiones temporales de pagos multipagos")

					var barcodes []string
					var pagos []uint
					for _, mp := range pagoMp {
						for _, mpdetalle := range mp.MultipagosDetalle {
							barcodes = append(barcodes, mpdetalle.CodigoBarras)
						}
					}
					//se debe buscar pago intentos relacionados con los codigos de barra
					filtroPagoIntento := filtros.PagoIntentoFiltro{
						Barcode:        barcodes,
						CargarPago:     true,
						CargarPagoTipo: true,
					}
					pagosIntentos, err := cl.administracion.GetPagosIntentosByTransaccionIdService(filtroPagoIntento)
					if err != nil {
						logs.Info(err.Error())
						logs.Info("no se pudo obtener pagos intentos para calculos de comsiones temporales")
					} else {
						// NOTE acumular por id de pago para luego calcular comisiones temporales
						for _, pintentos := range pagosIntentos {
							pagos = append(pagos, uint(pintentos.PagosID))
						}
					}
					if len(pagos) > 0 {
						responseCierreLote, err := cl.administracion.BuildPagosCalculoTemporales(pagos)
						if err == nil {
							ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
							// crear los movimientos temoorales y actualziar campo calculado en pago intento
							// esto inidica que el pago ya fue calculado y guardado en movimientostemporales
							err = cl.administracion.CreateMovimientosTemporalesService(ctx, responseCierreLote)
							if err != nil {
								logs.Error(err)
								logs.Info("no se pudo calcular comisiones temporales de pagos multipago")
							}
						}
					}

					// }
				}
			}
		}

	}

	listaLogArchivo = cierrelotedtos.PrismaLogArchivoResponse{
		NombreArchivo:  rutaArchivo[len(rutaArchivo)-1], //archivo.Name(),
		ArchivoLeido:   estado,
		ArchivoMovido:  false,
		LoteInsert:     estadoInsert,
		ErrorProducido: ErrorProducido,
	}
	return
}

func RecorrerArchivoMP(archivo *os.File) (listaMP []cierrelotemultipagosdtos.Multipagos, erro error) {
	readScanner := bufio.NewScanner(archivo)

	var header cierrelotemultipagosdtos.Header
	var trailer cierrelotemultipagosdtos.Trailler
	var detalles []cierrelotemultipagosdtos.Detalles

	for readScanner.Scan() {

		/* identificar headers, detalles y trailler */
		c1 := len(readScanner.Text())
		/* 1 armar headers */
		if readScanner.Text()[0:8] == "00000000" && len(readScanner.Text()) == 73 {
			header = cierrelotemultipagosdtos.Header{
				IdHeader:      readScanner.Text()[0:8],
				NombreEmpresa: readScanner.Text()[8:28],
				FechaProceso:  readScanner.Text()[28:36],
				IdArchivo:     readScanner.Text()[36:56],
				FillerHeader:  readScanner.Text()[56:73],
			}
			// erro = header.ValidarHeader()
			// if erro != nil {
			// 	logs.Error("error al validar el registro: " + erro.Error())
			// 	return
			// }
			/* 2 armar trailer */
		} else if readScanner.Text()[0:8] == "99999999" && len(readScanner.Text()) == 73 {

			trailer = cierrelotemultipagosdtos.Trailler{
				IdTrailler:    readScanner.Text()[0:8],
				CantDetalles:  readScanner.Text()[8:16],
				ImporteTotal:  readScanner.Text()[16:34],
				FillerTrailer: readScanner.Text()[34:73],
			}
			// erro = trailer.ValidarTrailer()
			// if erro != nil {
			// 	logs.Error("error al validar el registro: " + erro.Error())
			// 	return
			// }
			/* detalles */
		} else if c1 == 81 {
			detallesPago := cierrelotemultipagosdtos.Detalles{
				FechaCobro:     readScanner.Text()[0:8],
				ImporteCobrado: readScanner.Text()[8:23],
				CodigoBarras:   readScanner.Text()[23:71],
				Clearing:       readScanner.Text()[73:81],
			}
			// erro = detallesPago.ValidarDetalle()
			// if erro != nil {
			// 	logs.Error("error al validar el registro: " + erro.Error())
			// 	return
			// } else {
			detalles = append(detalles, detallesPago)

			// }
		} else {
			logs.Error("error al validar el registro: " + ERROR_FORMATO_REGISTRO_RP)
			return
		}

	}

	if len(detalles) > 0 {
		headerTrailer := cierrelotemultipagosdtos.HeaderTrailer{
			Header:  header,
			Trailer: trailer,
		}

		listaMP = append(listaMP, cierrelotemultipagosdtos.Multipagos{
			MultipagosHeader:   headerTrailer,
			MultipagosDetalles: detalles,
		})
	}

	return
}

func GenerarListaMP(cl *multipagoProcesarArchivos, impuesto administraciondtos.ResponseImpuesto, nombreArchivo []string, listaMP []cierrelotemultipagosdtos.Multipagos) (multipagosRegistros []entities.Multipagoscierrelote, erro error) {

	/*
		Autor: Jose Alarcon
		Descripción: Calcular comisión de pago MULTIPAGOS : se calcula 1.00 por cada tranasacción con un minimo de 40
		Estos valores se calculan con la tabla channelsarancel y IVA(de la tabla configuiracion)
		Se debe analizar si la fecha de vigencia es menor a la fecha de proceso del archivo
	*/
	// 1 obtener medio de pago MULTIPAGOS
	filtroChannel := filtros.ChannelFiltro{
		Channel: "MULTIPAGOS",
	}

	channel, erro := cl.administracion.GetChannelService(filtroChannel)
	if erro != nil && int64(channel.Id) < 0 {
		logs.Info(erro)
		return
	}

	// 2 obtener channels aracel
	filtroChannelArancel := filtros.ChannelArancelFiltro{
		CargarAllMedioPago: true,
		CargarChannel:      true,
		ChannelId:          channel.Id,
	}
	channelArancel, erro := cl.administracion.GetChannelsArancelService(filtroChannelArancel)
	logs.Info(channelArancel)
	if erro != nil {
		return
	}

	// esto permitira obtener el channels arancel del proveedor vigente
	var arancel administraciondtos.ResponseChannelsAranceles
	fecha_proceso := ConvertirFormatoFechaRapipago(listaMP[0].MultipagosHeader.Header.FechaProceso)
	fecha_proceso_archivo, _ := time.Parse("2006-01-02", fecha_proceso)
	for _, ch := range channelArancel.ChannelArancel {
		fecha_arancel, _ := time.Parse("2006-01-02T00:00:00Z", ch.Fechadesde)
		if fecha_proceso_archivo.Equal(fecha_arancel) || fecha_proceso_archivo.After(fecha_arancel) {
			arancel = administraciondtos.ResponseChannelsAranceles{
				Importe:       ch.Importe,
				Importeminimo: ch.Importeminimo,
				Importemaximo: ch.Importemaximo,
			}
		}

	}

	var importeTotalHeader float64
	var cabecera entities.Multipagoscierrelote
	// var detalle []entities.Rapipagocierrelotedetalles
	for _, lista := range listaMP {
		// result := ValidarTrailerDetalles(lista)
		// if !result {
		// 	erro = errors.New(ERROR_FORMATO_REGISTRO_RP_TRAILER_DETALLES)
		// 	return
		// } else {
		validar := commons.NewAlgoritmoVerificacion()
		// fecha_proceso := ConvertirFormatoFechaRapipago(lista.RapipagoHeader.Header.FechaProceso)
		fecha_clearing := ConvertirFormatoFechaRapipagoAcreditacion(lista.MultipagosDetalles[0].Clearing)
		// fecha_proc, _ := time.Parse("2006-01-02T00:00:00Z", fecha_proceso)
		fecha_acreditacion, err := time.Parse("2006-01-02", fecha_clearing)
		if err != nil {
			erro = err
			return nil, errors.New(erro.Error())
		}
		cant_dias, err := validar.CalcularDiasEntreFechas(fecha_proceso, fecha_clearing)
		if err != nil {
			erro = err
			return nil, errors.New(erro.Error())
		}

		cant_registros, err := strconv.ParseInt(lista.MultipagosHeader.Trailer.CantDetalles, 10, 64)
		if err != nil {
			erro = err
			return nil, errors.New(erro.Error())
		}
		importe_total, err := strconv.ParseInt(lista.MultipagosHeader.Trailer.ImporteTotal, 10, 64)
		if err != nil {
			erro = err
			return nil, errors.New(erro.Error())
		}
		cabecera = entities.Multipagoscierrelote{
			NombreArchivo:     nombreArchivo[len(nombreArchivo)-1],
			IdHeader:          strings.TrimSpace(lista.MultipagosHeader.Header.IdHeader),
			FechaProceso:      fecha_proceso,
			NombreEmpresa:     strings.TrimSpace(lista.MultipagosHeader.Header.NombreEmpresa),
			IdArchivo:         strings.TrimSpace(lista.MultipagosHeader.Header.IdArchivo),
			FillerHeader:      strings.TrimSpace(lista.MultipagosHeader.Header.FillerHeader),
			IdTrailer:         strings.TrimSpace(lista.MultipagosHeader.Trailer.IdTrailler),
			CantDetalles:      cant_registros,
			Fechaacreditacion: fecha_acreditacion,
			Cantdias:          cant_dias,
			ImporteTotal:      importe_total,
			ImporteMinimo:     arancel.Importeminimo,
			ImporteMaximo:     arancel.Importemaximo,
			Coeficiente:       arancel.Importe,
		}

		for _, detalleMP := range lista.MultipagosDetalles {
			fecha_cobro := ConvertirFormatoFechaRapipago(detalleMP.FechaCobro)
			importe_cobrado, err := strconv.ParseInt(detalleMP.ImporteCobrado, 10, 64)
			if err != nil {
				erro = err
				return nil, errors.New(erro.Error())
			}

			// format fecha de clearing obtenida de los detalles del archivo
			fecha_clearing_detalle := ConvertirFormatoFechaRapipagoAcreditacion(detalleMP.Clearing)

			/*
				Descripción:  Se calcula comision cobrada por rapipago
				se suman los importes de los detalles y se acumula en el header
			*/
			importeTotalCalculado := calcularComisionMP(cl, arancel, impuesto, importe_cobrado)

			mpdetalle := entities.Multipagoscierrelotedetalles{
				FechaCobro:       fecha_cobro,
				ImporteCobrado:   importe_cobrado,
				CodigoBarras:     strings.TrimSpace(detalleMP.CodigoBarras),
				ImporteCalculado: importeTotalCalculado,
				Clearing:         fecha_clearing_detalle,
			}
			importeTotalHeader += mpdetalle.ImporteCalculado
			cabecera.ImporteTotalCalculado = importeTotalHeader
			cabecera.MultipagosDetalle = append(cabecera.MultipagosDetalle, &mpdetalle)
		}
		/* armar lista de registros */
		// }
	}
	// cabecera.ImporteTotalCalculado = importeTotalCalculado

	multipagosRegistros = append(multipagosRegistros, cabecera)

	return
}

func calcularComisionMP(cl *multipagoProcesarArchivos, arancel administraciondtos.ResponseChannelsAranceles, impuesto administraciondtos.ResponseImpuesto, valor int64) (calculado float64) {

	valorCalculo := cl.utilService.ToFixed((float64(valor) / 100), 2)

	// calcular comision y verificar el minimo
	comision := valorCalculo * arancel.Importe

	// verificar si el valor es menor al minimo(20)
	if comision < arancel.Importeminimo {
		iva := arancel.Importeminimo * impuesto.Porcentaje
		descuento := arancel.Importeminimo + iva
		valorfinalcalculado := valorCalculo - descuento
		calculado = cl.utilService.ToFixed(valorfinalcalculado*100, 2)
	} else if comision > arancel.Importeminimo {
		iva := comision * impuesto.Porcentaje
		descuento := comision + iva
		valorfinalcalculado := valorCalculo - descuento
		calculado = cl.utilService.ToFixed(valorfinalcalculado*100, 2)
	}
	return
}
