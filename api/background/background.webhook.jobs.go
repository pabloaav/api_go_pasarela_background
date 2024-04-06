package background

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	webhooks_dtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	"github.com/robfig/cron"
)

func GetNotificacionPagosWebhook(cronjob *cron.Cron, periodicidad string, service administracion.Service) {
	var getNotificacionPagosWebhook = func() {

		filtroWebhook := webhooks_dtos.RequestWebhook{
			DiasPago:         25,
			PagosNotificado:  false,
			EstadoFinalPagos: true,
		}
		pagos, err := service.BuildNotificacionPagosService(filtroWebhook)
		if err == nil {
			pagosNotificar, err := service.CreateNotificacionPagosService(pagos)
			if err == nil {
				if len(pagosNotificar) > 0 {
					pagosupdate := service.NotificarPagos(pagosNotificar)
					if len(pagosupdate) > 0 { /* actualzar estado de pagos a notificado */
						err = service.UpdatePagosNoticados(pagosupdate)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes pagos que se notificaron al cliente no se actualizaron: %v", pagosupdate))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionWebhook,
								Descripcion: fmt.Sprintf("webhook: Error al actualizar estado de pagos a notificado .: %s", err),
							}
							service.CreateNotificacionService(notificacion)

						}
					} else {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionWebhook,
							Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos"),
						}
						service.CreateNotificacionService(notificacion)
					}

				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionWebhook,
						Descripcion: fmt.Sprintln("webhook: No existen pagos por notificar"),
					}
					service.CreateNotificacionService(notificacion)
				}
			}
		}
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de notificacion de pagos. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getNotificacionPagosWebhook)
}
func NotificarPagosWebhook(cronjob *cron.Cron, periodicidad string, service administracion.Service) {
	var notificarPagosWebHook = func() {
		err:=service.NotificarPagosWebhookSinNotificarService()
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de notificacion de pagos. %s", err),
			}
			erro:=service.CreateNotificacionService(notificacion)
			if erro != nil {
				logs.Error(erro)
			}
		}
	}
	// add job to cron
	cronjob.AddFunc(periodicidad, notificarPagosWebHook)
}