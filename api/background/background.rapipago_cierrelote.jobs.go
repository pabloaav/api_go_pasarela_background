package background

import (
	"context"
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/cierrelote"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	"github.com/robfig/cron"
)

func GetActualizarPagosCLRapipago(cronjob *cron.Cron, periodicidad string, cierrelote cierrelote.Service, service administracion.Service) {
	var getActualizarPagosCLRapipago = func() {

		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      false,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote  */
		listaPagoaRapipago, err := service.GetCierreLoteRapipagoService(filtroMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clrapipago no se puede continuar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(listaPagoaRapipago) > 0 {
			if err == nil {
				listaPagosClRapipago, err := service.BuildPagosClRapipago(listaPagoaRapipago)
				if err == nil {
					// Actualizar estados del pago y cierrelote
					logs.Info("inicio actualizacion de pagos rapipago")
					err = service.ActualizarPagosClRapipagoService(listaPagosClRapipago)
					if err != nil {
						logs.Error(err)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintln("error al actualizar estados de los pagos "),
						}
						service.CreateNotificacionService(notificacion)
					}
				}
			}

		} else {
			logs.Info("no existen pagos de rapipago para actualizar")
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("No existen pagos de rapipago para actualizar"),
			}
			service.CreateNotificacionService(notificacion)
		}

	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getActualizarPagosCLRapipago)
}

func GetNotificacionPagosCLRapipago(cronjob *cron.Cron, periodicidad string, service administracion.Service) {
	var getNotificacionPagosCLRapipago = func() {

		request := filtros.PagoEstadoFiltro{
			EstadoId: 4,
		}
		pagos, barcode, err := service.BuildNotificacionPagosCLRapipago(request)

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clrapipago no se puede continuar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(barcode) > 0 {
			if len(pagos) > 0 {
				pagosNotificar := service.NotificarPagos(pagos)
				if len(pagosNotificar) > 0 { /* actualzar estado de pagos a notificado */
					mensaje := fmt.Sprintf("los siguientes pagos se notificaron correctamente. %v", pagosNotificar)
					logs.Info(mensaje)

				}

			}

			//si la notificacion se realiza con exito se debera actualizar el campo pagonotificado en repiapagocierrolote
			err = service.ActualizarPagosClRapipagoDetallesService(barcode)
			if err != nil {
				msj := fmt.Sprintln("error al actualizar pagos en clrapipago")
				logs.Info(msj)
			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("no existen pagos por notificar"),
			}
			service.CreateNotificacionService(notificacion)
		}

	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getNotificacionPagosCLRapipago)
}

func GetConciliacionBancoRapipago(cronjob *cron.Cron, periodicidad string, service administracion.Service, movimientosBanco banco.BancoService) {
	var getConciliacionBancoRapipago = func() {

		filtroMovConciliarRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      true,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote - los que no fueron conciliados  */
		listaCierreRapipago, err := service.GetCierreLoteRapipagoService(filtroMovConciliarRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(listaCierreRapipago) > 0 {
			if err == nil {

				request := bancodtos.RequestConciliacion{
					TipoConciliacion: 1,
					ListaRapipago:    listaCierreRapipago,
				}
				// aqui hay retornar la lista de id de repipagocierre lote y los id del banco
				listaCierreRapipago, listaBancoId, err := movimientosBanco.ConciliacionPasarelaBanco(request)

				if len(listaBancoId) == 0 {
					logs.Info("no existen movimientos en banco para conciliar con pagos rapipago")
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("no existen movimientos en banco para conciliar con pagos rapipago: %s", err),
					}
					service.CreateNotificacionService(notificacion)
				} else {
					/*en el caso de error a actualizar la tabla rapipagocierrelote el proceso termina */
					err := service.UpdateCierreLoteRapipago(listaCierreRapipago.ListaRapipago)
					if err != nil {
						logs.Error(err)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote rapipago (volver a ejecutar proceso): %s", err),
						}
						service.CreateNotificacionService(notificacion)
					} else {
						// actualiza registro movimientos del banco
						// si no se actualiza los registros del banco se debera actualizar manualmente
						_, err := movimientosBanco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes movimientos del banco no se actualizaron: %v", listaBancoId))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar movimientos del banco - conciliacion rapipago(actualizar manualmente los siguientes movimientos): %s", err),
							}
							service.CreateNotificacionService(notificacion)
						}
					}

				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("No existen pagos de rapipago por conciliar"),
			}
			service.CreateNotificacionService(notificacion)
		}

	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getConciliacionBancoRapipago)
}

func GetGenerarMovimientosRapipago(cronjob *cron.Cron, periodicidad string, service administracion.Service) {
	var getGenerarMovimientosRapipago = func() {

		filtroMovMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: true,
			PagosNotificado:      true,
		}

		listaCierreMovRapipago, err := service.GetCierreLoteRapipagoService(filtroMovMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos para generar movimientos. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		// Si no se guarda ningún cierre no hace falta seguir el proce
		if len(listaCierreMovRapipago) > 0 {
			// 2 - Contruye los movimientos y hace la modificaciones necesarias para modificar los
			// pagos y demás datos necesarios en caso de error se repetira el día siguiente
			responseCierreLote, err := service.BuildRapipagoMovimiento(listaCierreMovRapipago)

			if err == nil {

				// 3 - Guarda los movimientos en la base de datos en caso de error se
				// repetira en el día siguiente
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				err = service.CreateMovimientosService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("No se pudo crear los movimientos clrapipago. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("no existen pagos para generar movimientos clrapipago"),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getGenerarMovimientosRapipago)
}

func GetCaducarPagoOffline(cronjob *cron.Cron, periodicidad string, service administracion.Service) {

	var getCaducarPagoOffline = func() {

		_, err := service.GetCaducarOfflineIntentos()

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionPagoOfflineExpirado,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de caducar pagos offline vencidos. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}

	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getCaducarPagoOffline)
}
