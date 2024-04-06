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

func GetActualizarPagosCLMultipagos(cronjob *cron.Cron, periodicidad string, cierrelote cierrelote.Service, service administracion.Service) {
	var getActualizarPagosCLMultipago = func() {

		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      false,
		}
		/* obtener lista pagos multipago encontrados en el tabla rapipagoscierrelote  */
		listaPagoMultipago, err := service.GetCierreLoteMultipagosService(filtroMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clmultipagos no se puede continuar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(listaPagoMultipago) > 0 {
			if err == nil {
				listaPagosClMultipago, err := service.BuildPagosClMultipagos(listaPagoMultipago)
				if err == nil {
					// Actualizar estados del pago y cierrelote
					logs.Info("inicio actualizacion de pagos multipago")
					err = service.ActualizarPagosClMultipagosService(listaPagosClMultipago)
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
			logs.Info("no existen pagos de multipagos para actualizar")
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("No existen pagos de multipagos para actualizar"),
			}
			service.CreateNotificacionService(notificacion)
		}

	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getActualizarPagosCLMultipago)
}

func GetNotificacionPagosCLMultipagos(cronjob *cron.Cron, periodicidad string, service administracion.Service) {
	var getNotificacionPagosCLMultipago = func() {

		request := filtros.PagoEstadoFiltro{
			EstadoId: 4,
		}
		pagos, barcode, err := service.BuildNotificacionPagosCLMultipagos(request)

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clmultipagos no se puede continuar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(barcode) > 0 {
			if len(pagos) > 0 {
				pagosNotificar := service.NotificarPagos(pagos)
				if len(pagosNotificar) == 0 {
					mensaje := fmt.Sprintf("los siguientes pagos se notificaron correctamente. %v", pagosNotificar)
					logs.Info(mensaje)
				}

			}

			err := service.ActualizarPagosClMultipagosDetallesService(barcode)
			if err != nil {
				msj := fmt.Sprintln("los pagos se notificaron correctamente pero no se pudo actualziar estados en la tabla multipagos")
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
	cronjob.AddFunc(periodicidad, getNotificacionPagosCLMultipago)
}

func GetConciliacionBancoMultipagos(cronjob *cron.Cron, periodicidad string, service administracion.Service, movimientosBanco banco.BancoService) {
	var getConciliacionBancoMultipago = func() {

		filtroMovConciliarMultipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      true,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote - los que no fueron conciliados  */
		listaCierreMultipago, err := service.GetCierreLoteMultipagosService(filtroMovConciliarMultipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(listaCierreMultipago) > 0 {
			if err == nil {

				request := bancodtos.RequestConciliacion{
					TipoConciliacion: 4,
					ListaMultipagos:  listaCierreMultipago,
				}
				// aqui hay retornar la lista de id de repipagocierre lote y los id del banco
				listaCierreMultipago, listaBancoId, err := movimientosBanco.ConciliacionPasarelaBanco(request)

				if len(listaBancoId) == 0 {
					logs.Info("no existen movimientos en banco para conciliar con pagos multipagos")
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("no existen movimientos en banco para conciliar con pagos multipagos: %s", err),
					}
					service.CreateNotificacionService(notificacion)
				} else {
					/*en el caso de error a actualizar la tabla rapipagocierrelote el proceso termina */
					err := service.UpdateCierreLoteMultipagos(listaCierreMultipago.ListaMultipagos)
					if err != nil {
						logs.Error(err)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote multipagos (volver a ejecutar proceso): %s", err),
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
								Descripcion: fmt.Sprintf("error al actualizar movimientos del banco - conciliacion multipagos (actualizar manualmente los siguientes movimientos): %s", err),
							}
							service.CreateNotificacionService(notificacion)
						}
					}

				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("No existen pagos de multipagos por conciliar"),
			}
			service.CreateNotificacionService(notificacion)
		}

	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getConciliacionBancoMultipago)
}

func GetGenerarMovimientosMultipagos(cronjob *cron.Cron, periodicidad string, service administracion.Service) {
	var GetGenerarMovimientosMultipagos = func() {

		filtroMovMovMultipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: true,
			PagosNotificado:      true,
		}

		listaCierreMovMultipago, err := service.GetCierreLoteMultipagosService(filtroMovMovMultipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos para generar movimientos. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		// Si no se guarda ningún cierre no hace falta seguir el proce
		if len(listaCierreMovMultipago) > 0 {
			// 2 - Contruye los movimientos y hace la modificaciones necesarias para modificar los
			// pagos y demás datos necesarios en caso de error se repetira el día siguiente
			responseCierreLote, err := service.BuildMultipagosMovimiento(listaCierreMovMultipago)

			if err == nil {

				// 3 - Guarda los movimientos en la base de datos en caso de error se
				// repetira en el día siguiente
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				err = service.CreateMovimientosService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("No se pudo crear los movimientos clmultipagos. %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("no existen pagos para generar movimientos clmultipagos"),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, GetGenerarMovimientosMultipagos)
}
