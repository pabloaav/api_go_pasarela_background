package background

import (
	"context"
	"errors"
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	"github.com/robfig/cron"
)

// realiza consultas al servicio de apilink. Notifica al cliente del cambio de estado de los pagos (debines)
func GetCierreLoteApiLink(cronjob *cron.Cron, periodicidad string, service administracion.Service) {

	var getCierreLoteApiLink = func() {
		listas, err := service.BuildCierreLoteApiLinkService()
		if err != nil {
			logs.Error(err)
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("error al consultar servicio apilink: %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
		if len(listas.ListaPagos) > 0 && len(listas.ListaCLApiLink) > 0 {
			ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
			err = service.CreateCLApilinkPagosService(ctx, listas)
			if err != nil {
				logs.Error(err)
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionCierreLote,
					Descripcion: fmt.Sprintf("error al crear registros en apilinkcierrelote: %s", err),
				}
				service.CreateNotificacionService(notificacion)
			} else {
				// NOTE NOTIFICAR AL USUARIO EL CAMBIO DE ESTADO DE LOS PAGOS
				filtro := linkdebin.RequestDebines{
					BancoExternalId: false,
					Pagoinformado:   true,
				}
				debines, _ := service.GetConsultarDebines(filtro)
				if len(debines) > 0 {
					// NOTE construir lote de pagos debin que se notificara al cliente
					pagos, debin, erro := service.BuildNotificacionPagosCLApilink(debines)
					if erro != nil {
						errorBuildNotificacion := errors.New("error al obtener debines para notificar al cliente")
						err = errorBuildNotificacion
						logError := entities.Log{
							Tipo:          entities.EnumLog("error"),
							Funcionalidad: "GetConsultarDebines",
							Mensaje:       errorBuildNotificacion.Error() + "-" + err.Error(),
						}
						errCrearLog := service.CreateLogService(logError)
						if errCrearLog != nil {
							logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
						}
					}
					if len(pagos) > 0 && len(debin) > 0 {
						// NOTE notificar lote de pagos a clientes
						pagosNotificar := service.NotificarPagos(pagos)
						if len(pagosNotificar) > 0 {
							//NOTE Si se envian los pagos con exito se debe actualziar el campo pagoinformado en la tabla aplilinkcierrelote
							filtro := linkdebin.RequestListaUpdateDebines{
								DebinId: debin,
							}
							// actualizar pagoinformado en tabla apilinkcierrelote
							erro := service.UpdateCierreLoteApilink(filtro)
							if erro != nil {
								logs.Error(erro)
								notificacion := entities.Notificacione{
									Tipo:        entities.NotificacionCierreLote,
									Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote apilink pagoinformado: %s", erro),
								}
								service.CreateNotificacionService(notificacion)
							}
						} else {
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionWebhook,
								Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos debines"),
							}
							service.CreateNotificacionService(notificacion)
						}

					}
				}
			}

		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getCierreLoteApiLink)
}

func GetConciliacionDebinBancoApiLink(cronjob *cron.Cron, periodicidad string, service administracion.Service, movimientosBanco banco.BancoService) {
	var getConciliacionDebinBancoApiLink = func() {

		filtro := linkdebin.RequestDebines{
			BancoExternalId:  false,
			CargarPagoEstado: true,
		}
		debines, err := service.GetDebines(filtro)
		if err != nil {
			errorBuildNotificacion := errors.New("error al obtener debines para conciliar con banco")
			err = errorBuildNotificacion
			logError := entities.Log{
				Tipo:          entities.EnumLog("error"),
				Funcionalidad: "GetConsultarDebines",
				Mensaje:       errorBuildNotificacion.Error() + "-" + err.Error(),
			}
			errCrearLog := service.CreateLogService(logError)
			if errCrearLog != nil {
				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
			}
		}
		if len(debines) > 0 {
			if err == nil {
				request := bancodtos.RequestConciliacion{
					TipoConciliacion: 2,
					ListaApilink:     debines,
				}
				listaCierreApiLinkBanco, listaBancoId, err := movimientosBanco.ConciliacionPasarelaBanco(request)
				if err != nil {
					logs.Error(err)
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("error al conciliar movimiento banco y cierre loteapilink: %s", err),
					}
					service.CreateNotificacionService(notificacion)

				} else {
					// NOTE Actualizar lista de cierreloteapilink campo banco external_id, match y fecha de acreditacion
					if len(listaCierreApiLinkBanco.ListaApilink) > 0 || len(listaCierreApiLinkBanco.ListaApilinkNoAcreditados) > 0 {
						listas := linkdebin.RequestListaUpdateDebines{
							Debines:              listaCierreApiLinkBanco.ListaApilink,
							DebinesNoAcreditados: listaCierreApiLinkBanco.ListaApilinkNoAcreditados,
						}
						erro := service.UpdateCierreLoteApilink(listas)
						if erro != nil {
							logs.Error(erro)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote apilink y conciliacion con banco: %s", erro),
							}
							service.CreateNotificacionService(notificacion)
						}
					}
					// FIXME se debe verificar si las 2 listas son iguales ?
					if len(listaBancoId) > 0 {
						_, err := movimientosBanco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes movimientos del banco no se actualizaron: %v", listaBancoId))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar registros del banco: %s", err),
							}
							service.CreateNotificacionService(notificacion)
							// en le caso de este error y si el pago no se actualizo a estados finales no afecta el cierre de apilink
							// el estado del pago se actualiza a estado final y no tendra en cuenta al consultar a apilink
							// ACCION : se debe actualizar manualmente el campo check en la tabla de movimientos de banco(no es obligatorio)
						}
					}

				}

			}

		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getConciliacionDebinBancoApiLink)
}

// Generar movimientos debines
func GetMovimientosDebinApiLink(cronjob *cron.Cron, periodicidad string, service administracion.Service) {

	var getMovimientosDebinApiLink = func() {
		filtro := linkdebin.RequestDebines{
			BancoExternalId:  true,
			CargarPagoEstado: true,
		}
		debines, err := service.GetDebines(filtro)
		if err != nil {
			errorBuildNotificacion := errors.New("error al obtener debines para generar movimientos")
			err = errorBuildNotificacion
			logError := entities.Log{
				Tipo:          entities.EnumLog("error"),
				Funcionalidad: "GetDebines",
				Mensaje:       errorBuildNotificacion.Error() + "-" + err.Error(),
			}
			errCrearLog := service.CreateLogService(logError)
			if errCrearLog != nil {
				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
			}
		}
		if len(debines) > 0 {
			responseCierreLote, err := service.BuildMovimientoApiLink(debines)
			if err == nil {
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				err = service.CreateMovimientosService(ctx, responseCierreLote)
				if err != nil {
					logs.Info("error al generar movimientos debines")
					logs.Error(err)
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("error al generar movimientos debines: %s", err),
					}
					service.CreateNotificacionService(notificacion)
				}
			}
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getMovimientosDebinApiLink)
}
