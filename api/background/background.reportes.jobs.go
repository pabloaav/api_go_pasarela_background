package background

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/reportes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	"github.com/robfig/cron"
)

func GenerarMovimientosTemporalesPagos(cronjob *cron.Cron, periodicidad string, service administracion.Service) {

	var GenerarCalculoMovimientosTemporales = func() {

		fechaActual := time.Now()
		// 1 Se debe consultar los pagos en estado aprobado y procesando
		request := filtros.PagoIntentoFiltros{
			PagoEstadosIds:      []uint64{4, 7},
			CargarPago:          true,
			CargarPagoTipo:      true,
			CargarPagoEstado:    true,
			CargarCuenta:        true,
			PagoIntentoAprobado: true,
			FechaPagoInicio:     fechaActual,
			FechaPagoFin:        fechaActual,
		}
		pagos, err := service.GetPagosCalculoMovTemporalesService(request)
		if len(pagos) > 0 {
			// se deben aplicar las comisiones a los pagos: aplicar al ultimo pago intento
			responseCierreLote, err := service.BuildPagosCalculoTemporales(pagos)
			if err == nil {
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				// crear los movimientos temoorales y actualziar campo calculado en pago intento
				// esto inidica que el pago ya fue calculado y guardado en movimientostemporales
				err = service.CreateMovimientosTemporalesService(ctx, responseCierreLote)
				if err != nil {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintf("No se pudo enviar crear los movimientos temporales. %s", err.Error()),
					}
					service.CreateNotificacionService(notificacion)
				}

			}

		}

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionSendEmailCsv,
				Descripcion: fmt.Sprintf("No se pudo recuperar los pagos del repositorio. %s", err.Error()),
			}
			service.CreateNotificacionService(notificacion)
		} else if len(pagos) < 1 {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionSendEmailCsv,
				Descripcion: fmt.Sprintf("No hay pagos para generar movimientos temporales"),
			}
			service.CreateNotificacionService(notificacion)
		}

		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionSendEmailCsv,
			Descripcion: fmt.Sprintf("Se generaron los movimientos temporales y sus comisiones"),
		}
		service.CreateNotificacionService(notificacion)
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, GenerarCalculoMovimientosTemporales)

}

func SendPagosClientes(cronjob *cron.Cron, periodicidad string, service administracion.Service, reportes reportes.ReportesService) {

	var sendPagosClientes = func() {
		// 1 obtener lista de cliente
		request := reportedtos.RequestPagosClientes{}
		clientes, err := reportes.GetClientes(request)
		if len(clientes.Clientes) > 0 {
			if err == nil {
				// obtener los pagos por cliente
				listaPagosClientes, err := reportes.GetPagosClientes(clientes, request)
				if err != nil {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintf("No se pudo obtener pagos para crear archivos csv de clientes. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}

				// enviar los pagos a clientes
				if len(listaPagosClientes) > 0 {
					listaErro, err := reportes.SendPagosClientes(listaPagosClientes)
					if err != nil {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionSendEmailCsv,
							Descripcion: fmt.Sprintf("No se pudo enviar archivos csv de clientes. %s", err),
						}
						service.CreateNotificacionService(notificacion)
					} else if len(listaErro) > 0 {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionSendEmailCsv,
							Descripcion: fmt.Sprintf("No se pudo enviar archivos csv de clientes. %s", listaErro),
						}
						service.CreateNotificacionService(notificacion)
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintln("No existen pagos de clientes para enviar a email"),
					}
					service.CreateNotificacionService(notificacion)
				}
			}

		}
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de enviar archivo de pagos. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, sendPagosClientes)
}

func SendRendicionesClientes(cronjob *cron.Cron, periodicidad string, service administracion.Service, reportes reportes.ReportesService) {

	var sendRendicionesClientes = func() {
		// 1 obtener lista de cliente
		request := reportedtos.RequestPagosClientes{}
		clientes, err := reportes.GetClientes(request)
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente los transferidos
				listaRendicionClientes, err := reportes.GetRendicionClientes(clientes, request)
				if err != nil {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintf("No se pudo obtener pagos para crear archivos de rendicion csv de clientes. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}

				if len(listaRendicionClientes) > 0 {
					listaErro, err := reportes.SendPagosClientes(listaRendicionClientes)
					if err != nil {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionSendEmailCsv,
							Descripcion: fmt.Sprintf("No se pudo enviar archivos rendicion csv de clientes. %s", err),
						}
						service.CreateNotificacionService(notificacion)
					} else if len(listaErro) > 0 {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionSendEmailCsv,
							Descripcion: fmt.Sprintf("No se pudo enviar archivos csv de clientes. %s", listaErro),
						}
						service.CreateNotificacionService(notificacion)
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintln("No existen rendiciones de clientes para enviar a email"),
					}
					service.CreateNotificacionService(notificacion)
				}
			}

		}
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de enviar archivo de pagos. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, sendRendicionesClientes)
}

func SendBatch(cronjob *cron.Cron, periodicidad string, service administracion.Service, reportes reportes.ReportesService) {

	var sendBatch = func() {

		// ctx := getContextAuditable(c)
		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})

		// 1 obtener lista de cliente
		request := reportedtos.RequestPagosClientes{}
		clientes, err := reportes.GetClientes(request)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionSendEmailCsv,
				Descripcion: fmt.Sprintf("No se pudo obtener clientes para construir el archivo. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {
				// obtener los pagos/pagoitems por cliente
				// NOTE solo se obtiene los que son movimientos -> pagos autorizados
				listaPagosItems, err := reportes.GetPagoItems(clientes, request)
				if err != nil {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionArchivoBatchCliente,
						Descripcion: fmt.Sprintf("No se pudo obtener lista de pagos de clientes para procesar archivo batch. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}
				// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo
				if len(listaPagosItems) > 0 {
					// se crea la estructura correspondiente para el archivo
					resultpagositems := reportes.BuildPagosItems(listaPagosItems)
					if len(resultpagositems) > 0 {
						err := reportes.ValidarEsctucturaPagosItems(resultpagositems) // validar estructura antes de crear el archivo
						if err != nil {
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionArchivoBatchCliente,
								Descripcion: fmt.Sprintf("la estructura del archivo batch creado es incorrecta. %s", err),
							}
							service.CreateNotificacionService(notificacion)
						} else {
							// enviar archivo por sftp
							err := reportes.SendPagosItems(ctx, resultpagositems, request)
							if err != nil {
								notificacion := entities.Notificacione{
									Tipo:        entities.NotificacionArchivoBatchCliente,
									Descripcion: fmt.Sprintf("no se puedo enviar el archivo batch a clientes. %s", err),
								}
								service.CreateNotificacionService(notificacion)
							} else {
								logs.Info("Archivo batch enviado con exito")
								notificacion := entities.Notificacione{
									Tipo:        entities.NotificacionArchivoBatchCliente,
									Descripcion: fmt.Sprintln("el archivo batch se envio con exito"),
								}
								service.CreateNotificacionService(notificacion)
							}

						}
					} else {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionArchivoBatchCliente,
							Descripcion: fmt.Sprintf("existe un error al costruir el archivo batch del cliente. %s", err),
						}
						service.CreateNotificacionService(notificacion)
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionArchivoBatchCliente,
						Descripcion: fmt.Sprintf("no existen pagos batch para informar al cliente. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}
			}

		}

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionArchivoBatchCliente,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de enviar archivo batch a clientes. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, sendBatch)
}

func SendBatchPagos(cronjob *cron.Cron, periodicidad string, service administracion.Service, reportes reportes.ReportesService) {

	var sendBatchPagos = func() {

		fechaHoy := time.Now()
		year, month, day := fechaHoy.Date()
		fechaInicio := time.Date(year, month, day, 00, 00, 00, 000, fechaHoy.Location())
		fechaFin := time.Date(year, month, day, 23, 59, 59, 000, fechaInicio.Location())

		// fechaInicio, _ := time.Parse("2006-01-02 15:04:05", "2023-12-21 00:00:00")
		// fechaFin, _ := time.Parse("2006-01-02 15:04:05", "2023-12-21 23:59:00")

		ctx := context.Background()

		request := reportedtos.RequestPagosClientes{
			ClientesIds: []uint{18, 19}, //Goya y Corrientes Acor
			FechaInicio: fechaInicio,
			FechaFin:    fechaFin,
		}

		// 1 obtener lista de cliente
		clientes, err := reportes.GetClientes(request)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionArchivoBatchCliente,
				Descripcion: fmt.Sprintf("No se pudo obtener clientes para construir el archivo. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {
				// obtener los pagos/pagoitems por cliente
				listaPagosItems, err := reportes.GetPagoItemsAlternativo(clientes, request)
				if err != nil {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionArchivoBatchCliente,
						Descripcion: fmt.Sprintf("No se pudo obtener lista de pagos de clientes para procesar archivo batch. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}
				// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo
				if len(listaPagosItems) > 0 {
					resultpagositems := reportes.BuildPagosArchivo(listaPagosItems)
					if len(resultpagositems) > 0 {
						err := reportes.ValidarEsctucturaPagosBatch(resultpagositems) // validar estructura antes de crear el archivo
						if err != nil {
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionArchivoBatchCliente,
								Descripcion: fmt.Sprintf("no se puedo enviar el archivo batch a clientes. %s", err),
							}
							service.CreateNotificacionService(notificacion)
						} else {
							err := reportes.SendPagosBatch(ctx, resultpagositems, request)
							if err != nil {
								notificacion := entities.Notificacione{
									Tipo:        entities.NotificacionArchivoBatchCliente,
									Descripcion: fmt.Sprintf("no se puedo enviar el archivo batch a clientes. %s", err),
								}
								service.CreateNotificacionService(notificacion)
							} else {
								// caso de exito
								notificacion := entities.Notificacione{
									Tipo:        entities.NotificacionArchivoBatchCliente,
									Descripcion: "Proceso de envio de email con batch exitoso.",
								}
								service.CreateNotificacionService(notificacion)
							}

						}
					}
				} else {
					// caso sin resultados
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionArchivoBatchCliente,
						Descripcion: fmt.Sprintf("no se puedo enviar el archivo batch a clientes. %s", "No existen pagos para informar"),
					}
					service.CreateNotificacionService(notificacion)

				}
			}

		} else {
			// caso error en la busqueda de clientes
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionArchivoBatchCliente,
				Descripcion: fmt.Sprintf("no se puedo enviar el archivo batch a clientes. %s", "No se encontro clientes a informar"),
			}
			service.CreateNotificacionService(notificacion)

		}

	}
	// add job to cron
	cronjob.AddFunc(periodicidad, sendBatchPagos)

}

func SendPagosDiariosTelco(cronjob *cron.Cron, periodicidad string, service administracion.Service, reportes reportes.ReportesService) {

	var sendPagosClientesDiario = func() {
		// 1 obtener lista de cliente
		request_cliente := reportedtos.RequestPagosClientes{}

		clientes, err := reportes.GetClientes(request_cliente)

		request := reportedtos.RequestCobranzasDiarias{}
		if len(clientes.Clientes) > 0 {
			if err == nil {
				// obtener los pagos por cliente
				listaPagosClientes, err := reportes.GetPagos(clientes, request)
				if err != nil {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintf("No se pudo obtener pagos para crear reporte diario de clientes. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}

				// enviar los pagos a clientes
				if len(listaPagosClientes) > 0 {
					listaErro, err := reportes.SendPagosClientes(listaPagosClientes)
					if err != nil {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionSendEmailCsv,
							Descripcion: fmt.Sprintf("No se pudo enviar reporte diario de clientes. %s", err),
						}
						service.CreateNotificacionService(notificacion)
					} else if len(listaErro) > 0 {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionSendEmailCsv,
							Descripcion: fmt.Sprintf("No se pudo enviar  reporte diario de clientes. %s", listaErro),
						}
						service.CreateNotificacionService(notificacion)
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintln("No existen pagos de clientes en reporte diario para enviar a email"),
					}
					service.CreateNotificacionService(notificacion)
				}
			}

		}
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de enviar archivo de pagos. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, sendPagosClientesDiario)
}

func SendReporteControlCobranzas(cronjob *cron.Cron, periodicidad string, reportes reportes.ReportesService, runEndpoint util.RunEndpoint) {

	var EnvioReporteControlCobranzas = func() {

		fechaHoy := time.Now().AddDate(0, 0, -1).Format("02-01-2006")

		apikeysParaControl := []string{
			"4e9af107-5420-414c-9e1f-811e18d8c895",
			"cf4963ba-bc89-48e0-89db-0962398788f9",
			"e32da2af-158a-45e7-bcbf-bd42f4413049",
			"70f299cc-918d-4705-b306-d22774d83ee8",
			"27d02a31-ca96-427f-b60d-bafa8036ef02",
		}

		apikeys := strings.Join(apikeysParaControl, ",")

		response, err := reportes.MakeControlReportes(apikeys, fechaHoy, "token", runEndpoint)

		if err != nil {
			fmt.Println("Error al enviar email", err.Error())
		}

		err = reportes.SendControlReportes(response)
		if err != nil {
			fmt.Println("Error al enviar email", err.Error())
		}

	}

	// add job to cron
	cronjob.AddFunc(periodicidad, EnvioReporteControlCobranzas)

}

func SendReversionesClientes(cronjob *cron.Cron, periodicidad string, service administracion.Service, reportes reportes.ReportesService) {

	var sendReversionesClientes = func() {
		// 1 obtener lista de cliente
		request := reportedtos.RequestPagosClientes{}
		clientes, err := reportes.GetClientes(request)
		if len(clientes.Clientes) > 0 {
			if err == nil {
				listaReversionesClientes, err := reportes.GetReversionesClientes(clientes, request)
				if err != nil {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintf("No se pudo obtener pagos para crear archivos de reversiones csv de clientes. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}
				if len(listaReversionesClientes) <= 0 {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintf("no existen reversiones para informar a los clientes. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}

				// enviar los reverciones a clientes
				listaErro, err := reportes.SendPagosClientes(listaReversionesClientes)
				logs.Info(listaErro)
				if err != nil {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintf("No se pudo enviar archivos batch pagos a clientes. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}
				if len(listaErro) > 0 {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionSendEmailCsv,
						Descripcion: fmt.Sprintf("No se pudo enviar archivos batch pagos a clientes. %s", listaErro),
					}
					service.CreateNotificacionService(notificacion)
				}
			}

		}
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de enviar archivo de pagos. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, sendReversionesClientes)
}
