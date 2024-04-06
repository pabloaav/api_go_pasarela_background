package background

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/cierrelote"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/reportes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	"github.com/robfig/cron"
)

func _buildNotificacion(service util.UtilService, erro error, tipo entities.EnumTipoNotificacion) {
	notificacion := entities.Notificacione{
		Tipo:        tipo,
		Descripcion: fmt.Sprintf("Configuración inválida. %s", erro.Error()),
	}
	service.CreateNotificacionService(notificacion)
}

func _buildPeriodicidad(service util.UtilService, nombreConfig string, valorConfig string, descripcionConfig string) (configuracion entities.Configuracione, erro error) {

	filtro := filtros.ConfiguracionFiltro{
		Nombre: nombreConfig,
	}

	// busca una configuracion con el mismo nombre recibido por parametro
	config, erro := service.GetConfiguracionService(filtro)
	// si la encuentra, setea el objeto de respuesta con tales valores recuperados de la BD
	configuracion.Nombre = config.Nombre
	configuracion.Valor = config.Valor
	configuracion.ID = config.Id

	if erro != nil {
		_buildNotificacion(service, erro, entities.NotificacionConfiguraciones)
		return
	}
	// si no encuentra configuracion, crea una con los parametro recibidos, y guarda la misma en la BD
	if configuracion.ID == 0 {

		config := administraciondtos.RequestConfiguracion{
			Nombre:      nombreConfig,
			Descripcion: descripcionConfig,
			Valor:       valorConfig,
		}

		configuracion = config.ToEntity(false)

		_, erro = service.CreateConfiguracionService(config)

		if erro != nil {

			_buildNotificacion(service, erro, entities.NotificacionConfiguraciones)

		}

	}

	return

}

func BackgroudServices(service administracion.Service, cierrelote cierrelote.Service, util util.UtilService, movimientosBanco banco.BancoService, reportes reportes.ReportesService, runEndpoint util.RunEndpoint) {

	c := cron.New()

	/* Begin Cierre Lote Prisma */

	confCierreLotePrismaRapipagoArchivo, err := _buildPeriodicidad(util, "PERIODICIDAD_PRISMA_RAPIPAGO_CIERRE_LOTE", "0 35 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad procesar archivos de rapipago y prisma")

	if err != nil {
		panic(err)
	}

	/*  ejecutar proceso archivos prisma  y rapipago */
	GetArchivosTxtCierreLote(c, confCierreLotePrismaRapipagoArchivo.Valor, cierrelote, service)

	/* INICIO PROCESO PARA OBTENER LOS ARCHIVOS TXT DE S3 Y GUARDAR EN LA DB LA INFORMACION OBTENIDA DEL TXT */
	confPrismaCierreLote, err := _buildPeriodicidad(util, "PERIODICIDAD_PRISMA_CIERRE_LOTE", "0 00 13 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de prisma")

	if err != nil {
		panic(err)
	}

	/*  ejecutar proceso */
	GetArchivosTxtCierreLote(c, confPrismaCierreLote.Valor, cierrelote, service)

	/* INICIO PROCESAR TABLA MOVIMIENTOS MX */
	confProcesarTablaMx, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESAR_TABLA_MX", "0 00 14 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de procesar tabla mx")

	if err != nil {
		panic(err)
	}

	/*  ejecutar proceso */
	GetTablaMovimientosMXCierreLote(c, confProcesarTablaMx.Valor, cierrelote, service)

	/* INICIO PROCESAR TABLA PAGOS PX*/
	confProcesarTablaPx, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESAR_TABLA_PX", "0 00 14 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de procesar tabla px")

	if err != nil {
		panic(err)
	}

	/*  ejecutar proceso */
	GetPagosPXCierreLote(c, confProcesarTablaPx.Valor, cierrelote, service)

	/* INICIO PROCESO CONCILIAR CIERRE LOTE Y MOVIMIENTOS PRISMA*/
	confProcesoConciliarClMx, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_MX", "0 30 14 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de proceso conciliacion cl-Mx")

	if err != nil {
		panic(err)
	}

	/*  ejecutar proceso */
	GetConciliarCLMX(c, confProcesoConciliarClMx.Valor, cierrelote, service)

	/* TODO -> INICIO PROCESO CONCILIAR CIERRE LOTE Y PAGOS PRISMA*/
	confProcesoConciliarClPx, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_PX", "0 00 15 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de proceso conciliacion cl-px")

	if err != nil {
		panic(err)
	}
	GetConciliarCLPX(c, confProcesoConciliarClPx.Valor, cierrelote, service)

	// /* TODO -> INICIO PROCESO CONCILIAR CON EL BANCO*/
	confProcesoConciliarClBanco, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_BANCO", "0 00 05 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de proceso conciliacion cl-banco")

	if err != nil {
		panic(err)
	}
	/*  ejecutar proceso */
	GetConciliacionBancoPrisma(c, confProcesoConciliarClBanco.Valor, false, cierrelote, service)

	confProcesoConciliarClBancoReversion, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_BANCO_REVERSION", "0 30 05 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de proceso conciliacion cl-banco reversion")

	if err != nil {
		panic(err)
	}
	/*  ejecutar proceso */
	GetConciliacionBancoPrisma(c, confProcesoConciliarClBancoReversion.Valor, true, cierrelote, service)

	/* TODO -> INICIO PROCESO GENERAR MOVIMIENTOS */
	confProcesoBuildMovimiento, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_PAGOS_PSP", "0 00 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad que concilia cl con pagos pasarela")

	if err != nil {
		panic(err)
	}
	// movimiento_reversion = false
	/*  ejecutar proceso */
	GetGenerarMovimientosPrisma(c, confProcesoBuildMovimiento.Valor, false, cierrelote, service)

	confProcesoBuildMovimientoRevertido, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_CONCILIACION_CL_REVERTIDOS_PAGOS_PSP", "0 30 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad que concilia cl revertidos con pagos pasarela")

	if err != nil {
		panic(err)
	}
	// movimiento_reversion = true
	/*  ejecutar proceso */
	GetGenerarMovimientosPrisma(c, confProcesoBuildMovimientoRevertido.Valor, true, cierrelote, service)

	/* Begin Cierre Lote Rapipago */

	/*
		PROCESO 2:
			se realiza conciliacion de los archivos de cierre de lotes con los movimientos de banco
		PROCESO 3:
			se crean los diferentes objetos para registrar los movimineto recibidos en el cierre de lote
	*/

	// TODO INICIO PROCESO RAPIPAGO
	// & 1 Se procesan los archivos se guardan en la tabla rapipago
	// & 2 Se actualizan estados de los pagos con los encontrado en cierrelote(archivo recibido) -> EL estado APROBADO indica que el pagador fue a un rapipago
	// & 3 Se notifica el cambia de estado al cliente(se ejecuta webhook)
	// & 4 Conciliar con los movimientos ingresados en banco
	// & 5 Generar movimientos

	// & 2
	confRapipagoCierreLote, err := _buildPeriodicidad(util, "PERIODICIDAD_RAPIPAGO_CIERRE_LOTE_PARTE2", "0 50 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
	if err != nil {
		panic(err)
	}
	GetActualizarPagosCLRapipago(c, confRapipagoCierreLote.Valor, cierrelote, service)

	//& 3
	confRapipagoCierreLoteNotificarPagos, err := _buildPeriodicidad(util, "PERIODICIDAD_RAPIPAGO_CIERRE_LOTE_PARTE3", "0 00 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
	if err != nil {
		panic(err)
	}
	GetNotificacionPagosCLRapipago(c, confRapipagoCierreLoteNotificarPagos.Valor, service)

	//& 4
	confRapipagoCierreLoteConciliarBanco, err := _buildPeriodicidad(util, "PERIODICIDAD_RAPIPAGO_CIERRE_LOTE_PARTE4", "0 00 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
	if err != nil {
		panic(err)
	}
	GetConciliacionBancoRapipago(c, confRapipagoCierreLoteConciliarBanco.Valor, service, movimientosBanco)

	// & 5
	confRapipagoCierreLoteGenerarMovimientos, err := _buildPeriodicidad(util, "PERIODICIDAD_RAPIPAGO_CIERRE_LOTE_PARTE5", "0 30 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
	if err != nil {
		panic(err)
	}
	GetGenerarMovimientosRapipago(c, confRapipagoCierreLoteGenerarMovimientos.Valor, service)

	// TODO INICIO PROCESO MULTIPAGOS
	// & 1 Se procesan los archivos se guardan en la tabla multipagoscierrelotes
	// & 2 Se actualizan estados de los pagos con los encontrado en cierrelote(archivo recibido) -> EL estado APROBADO indica que el pagador fue a un multipagos
	// & 3 Se notifica el cambia de estado al cliente(se ejecuta webhook) de no tener notificado online el pagointento
	// & 4 Conciliar con los movimientos ingresados en banco (Movimientos con referencia "multipagos")
	// & 5 Generar movimientos

	// & 2
	confMultipagoCierreLote, err := _buildPeriodicidad(util, "PERIODICIDAD_MULTIPAGO_CIERRE_LOTE_PARTE2", "0 52 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
	if err != nil {
		panic(err)
	}
	GetActualizarPagosCLMultipagos(c, confMultipagoCierreLote.Valor, cierrelote, service)

	//& 3
	confMultipagoCierreLoteNotificarPagos, err := _buildPeriodicidad(util, "PERIODICIDAD_MULTIPAGO_CIERRE_LOTE_PARTE3", "0 02 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
	if err != nil {
		panic(err)
	}
	GetNotificacionPagosCLMultipagos(c, confMultipagoCierreLoteNotificarPagos.Valor, service)

	//& 4
	confMultipagoCierreLoteConciliarBanco, err := _buildPeriodicidad(util, "PERIODICIDAD_MULTIPAGO_CIERRE_LOTE_PARTE4", "0 02 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
	if err != nil {
		panic(err)
	}
	GetConciliacionBancoMultipagos(c, confMultipagoCierreLoteConciliarBanco.Valor, service, movimientosBanco)

	// & 5
	confMultipagooCierreLoteGenerarMovimientos, err := _buildPeriodicidad(util, "PERIODICIDAD_MULTIPAGO_CIERRE_LOTE_PARTE5", "0 32 06 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de rapipago")
	if err != nil {
		panic(err)
	}
	GetGenerarMovimientosMultipagos(c, confMultipagooCierreLoteGenerarMovimientos.Valor, service)

	// // 0 */5 7-16 * * * : de 7 a 16 hs se ejecuta cada 5 minutos
	// // 0 0 */1 17-22 * * : de 17 a 22 hs se ajecuta cada 1 hora
	//  ^  PROCESO APILINK
	// ^ PASO 1 Y 2: CONSULTA SERVICIO APILINK Y NOTIFICACION AL CLIENTE
	confApilinkCierreLoteMorning, err := _buildPeriodicidad(util, "PERIODICIDAD_APILINK_CIERRE_LOTE_MORNING", "0 0 */3 7-12 * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
	if err != nil {
		panic(err)
	}
	GetCierreLoteApiLink(c, confApilinkCierreLoteMorning.Valor, service)

	confApilinkCierreLoteAfternoon, err := _buildPeriodicidad(util, "PERIODICIDAD_APILINK_CIERRE_LOTE_AFTERNOON", "0 0 */5 13-22 * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
	if err != nil {
		panic(err)
	}
	GetCierreLoteApiLink(c, confApilinkCierreLoteAfternoon.Valor, service)

	confApilinkCierreNight, err := _buildPeriodicidad(util, "PERIODICIDAD_APILINK_CIERRE_LOTE_NIGHT", "0 00 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
	if err != nil {
		panic(err)
	}
	GetCierreLoteApiLink(c, confApilinkCierreNight.Valor, service)

	// ^ PASO 3: Conciliacion pagos debin con banco
	confConciliacionBancoDebin, err := _buildPeriodicidad(util, "PERIODICIDAD_CONCILIACION_PAGOSDEBIN_BANCO", "0 15 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
	if err != nil {
		panic(err)
	}
	GetConciliacionDebinBancoApiLink(c, confConciliacionBancoDebin.Valor, service, movimientosBanco)

	// ^ PASO 4: Generar movimientos debines
	confGenerarMovimientosDebines, err := _buildPeriodicidad(util, "PERIODICIDAD_GENERAR_MOV_DEBINES", "0 30 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de cierre de lote de apilink")
	if err != nil {
		panic(err)
	}
	GetMovimientosDebinApiLink(c, confGenerarMovimientosDebines.Valor, service)

	// * WEBHOOK. NOTIFICACION DE PAGOS A CLIENTES */
	confNotificacionPagos, err := _buildPeriodicidad(util, "PERIODICIDAD_NOTIFICACION_PAGOS", "0 0 */12 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad de notificacion de pagos a los clientes")
	if err != nil {
		panic(err)
	}
	GetNotificacionPagosWebhook(c, confNotificacionPagos.Valor, service)

	// NOTE TRANSFERENCIAS RETIROS AUTOMATICOS CLIENTES
	confRetiroAutomatico, err := _buildPeriodicidad(util, "PERIODICIDAD_RETIRO_AUTOMATICO", "0 00 09 * * 1-5", "Periodicidad (en formato cron) en que el sistema realiza automaticamente las transferencias a los clientes.")
	if err != nil {
		panic(err)
	}
	GetTransferenciasAutomaticas(c, confRetiroAutomatico.Valor, service, util)

	// NOTE TRANSFERENCIAS RETIROS AUTOMATICOS COMISIONES
	confRetiroAutomaticoComisiones, err := _buildPeriodicidad(util, "PERIODICIDAD_RETIRO_AUTOMATICO_COMISIONES", "0 30 13 * * 1-5", "Periodicidad (en formato cron) en que el sistema realiza automaticamente las transferencias a los clientes.")
	if err != nil {
		panic(err)
	}
	GetTransferenciasAutomaticasComisiones(c, confRetiroAutomaticoComisiones.Valor, service, util)

	//  ENVIO DE REPORTES: PAGOS , RENDICION
	confSendPagos, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_PAGOS", "0 10 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar archivos de pagos a los clientes")
	if err != nil {
		panic(err)
	}
	SendPagosClientes(c, confSendPagos.Valor, service, reportes)

	confSendArchivoRendicion, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_ARCHIVO_RENDICION", "0 00 10 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar archivos de rendicion a los clientes")
	if err != nil {
		panic(err)
	}
	SendRendicionesClientes(c, confSendArchivoRendicion.Valor, service, reportes)

	// NOTE envio de reversiones
	confSendArchivoRendicion, err = _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_ARCHIVO_REVERSION", "0 15 10 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar archivos de rendicion a los clientes")
	if err != nil {
		panic(err)
	}
	SendReversionesClientes(c, confSendArchivoRendicion.Valor, service, reportes)

	// * ENVIO DE ARCHIVOS BATCH: SOLO DPEC
	confSendArchivoBatch, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_ARCHIVO_BATCH", "0 45 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar archivos batch a los clientes")
	if err != nil {
		panic(err)
	}
	SendBatch(c, confSendArchivoBatch.Valor, service, reportes)

	// * ENVIO DE ARCHIVOS BATCH PAGOS (GOYA)
	confSendArchivoBatchPagos, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_ARCHIVO_BATCH_PAGOS", "0 55 07 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar archivos batch pagos a los clientes")
	if err != nil {
		panic(err)
	}
	SendBatchPagos(c, confSendArchivoBatchPagos.Valor, service, reportes)

	// // Caducar pagos con metodos offline que expiran
	confCaducarPagosOfflineExpirados, _ := _buildPeriodicidad(util, "PERIODICIDAD_CADUCAR_PAGOSOFFLINE_EXPIRADOS", "0 00 01 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad caducar pagos con pagointentos offline procesando vencidos")
	GetCaducarPagoOffline(c, confCaducarPagosOfflineExpirados.Valor, service)

	// ENVIO DE REPORTE DIARIO PAGOS
	confReporteDiarioTelco, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_REPORTE_PAGOS_DIARIO", "0 38 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar reporte diario de pagos a correos Telco en la DB")
	if err != nil {
		panic(err)
	}
	SendPagosDiariosTelco(c, confReporteDiarioTelco.Valor, service, reportes)

	// GENERACION DE MOVIMIENTOS TEMPORALES DE PAGOS
	confMovimientosTemporales, err := _buildPeriodicidad(util, "PERIODICIDAD_GENERAR_MOVIMIENTOS_TEMPORALES", "0 50 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad generar movimientos temporales con calculos de comisiones")
	if err != nil {
		panic(err)
	}
	GenerarMovimientosTemporalesPagos(c, confMovimientosTemporales.Valor, service)
	// logs.Info(confCaducarPagosOfflineExpirados)
	// if err != nil {
	// 	panic(err)
	// }

	// movimiento_reversion = false
	/*  ejecutar proceso */

	confProcesoBuildReporteControlCobranzas, err := _buildPeriodicidad(util, "PERIODICIDAD_PROCESO_REPORTE_CONTROL", "0 15 08 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad que envia reporte control de reportes cobranzas")

	if err != nil {
		panic(err)
	}

	SendReporteControlCobranzas(c, confProcesoBuildReporteControlCobranzas.Valor, reportes, runEndpoint)

	// /* End Reportes y Envios */

	// /* Begin Otros Procesos */

	// /*  ejecutar proceso */
	// GetCaducarPagoOffline(c, "0 40 12 * * *", service)

	/* End Otros Procesos */
	//****************************************************************
	// ENVIO DE NOTIFICACIONES A CLIENTES
	SendNotificacionesDiario, err := _buildPeriodicidad(util, "PERIODICIDAD_ENVIAR_NOTIFICACIONES_CLIENTES_DIARIO", "0 00 19 * * *", "Periodicidad (en formato cron) en que se ejecuta la funcionalidad enviar notificaciones de los pagos a los clientes correspondiente, en caso que falle el envio normal")
	if err != nil {
		panic(err)
	}
	NotificarPagosWebhook(c, SendNotificacionesDiario.Valor, service)
	//****************************************************************

	// Iniciar el proceso cron
	c.Start()
}
