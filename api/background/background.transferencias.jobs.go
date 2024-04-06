package background

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	"github.com/robfig/cron"
	uuid "github.com/satori/go.uuid"
)

func GetTransferenciasAutomaticas(cronjob *cron.Cron, periodicidad string, service administracion.Service, util util.UtilService) {

	var getTransferenciasAutomaticas = func() {
		var feriado bool

		/* NO TRANSFERIR FEIRADOS */
		filtro := filtros.ConfiguracionFiltro{
			Nombre: "FERIADOS",
		}

		// buscar la configuracion de dias feriados
		configuracion, erro := util.GetConfiguracionService(filtro)

		// si hay error se notifica pero se continua
		if erro != nil {
			_buildNotificacion(util, erro, entities.NotificacionConfiguraciones)
		}

		// si se obtiene el resultado de la configuracion para FERIADOS
		if configuracion.Id != 0 {
			feriado = _esFeriado(configuracion.Valor)
		}

		// si NO es feriado, hacer la transferencia
		if !feriado {
			ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})

			// logs.Info(clientes)

			// // esta funcion recibe los id de los clientes cargados en la configuracion
			idcliente := administraciondtos.TransferenciasClienteId{}

			requestMovimientosId, err := service.RetiroAutomaticoClientes(ctx, idcliente)

			if err != nil {
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionPagoExpirado,
					Descripcion: fmt.Sprintf("No se pudo realizar la transferencia automatica para los clientes. %s", err),
				}
				service.CreateNotificacionService(notificacion)
			}

			if err == nil && len(requestMovimientosId.MovimientosId) == 0 {
				logs.Info("las transferencias automaticas se realizaron con exito")
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionTransferenciaAutomatica,
					Descripcion: fmt.Sprintln("el proceso de transferencias automaticas se realizaron con exito. Sin Movimientos"),
				}
				service.CreateNotificacionService(notificacion)
			}

			if err == nil && len(requestMovimientosId.MovimientosId) > 0 {
				logs.Info("las transferencias automaticas se realizaron con exito")
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionTransferenciaAutomatica,
					Descripcion: fmt.Sprintln("el proceso de transferencias automaticas se realizaron con exito"),
				}
				service.CreateNotificacionService(notificacion)
			}

		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getTransferenciasAutomaticas)
}

func GetTransferenciasAutomaticasComisiones(cronjob *cron.Cron, periodicidad string, service administracion.Service, util util.UtilService) {

	var GetTransferenciasAutomaticasComisiones = func() {
		// var request administraciondtos.RequestComisiones

		request := administraciondtos.RequestComisiones{
			FechaInicio:      time.Now(),
			FechaFin:         time.Now(),
			RetiroAutomatico: true,
		}

		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
		uuid := uuid.NewV4()
		_, err := service.SendTransferenciasComisiones(ctx, uuid.String(), request)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionPagoExpirado,
				Descripcion: fmt.Sprintf("No se pudo realizar la transferencia automatica de comisiones. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		}
	}
	cronjob.AddFunc(periodicidad, GetTransferenciasAutomaticasComisiones)
}

/* FUNCIONES PROPIAS */
// retorna true o false segun la fecha actual sea feriado, comparando con un string de fechas
func _esFeriado(stringFechas string) (result bool) {
	// se toma la fecha actual en formato yyyy-mm-dd
	now := time.Now().UTC().Format("2006-01-02")
	// separador del split
	var sep string = ","
	// fechas en formato yyyy-mm-dd en tipo []string
	fechas := strings.Split(stringFechas, sep)

	result = commons.ContainStrings(fechas, now)
	return
}
