package reportes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros_reportes "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/reportes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportesRepository interface {
	GetPagosReportes(request reportedtos.RequestPagosPeriodo, pendiente uint) (pagos []entities.Pago, erro error)
	GetCobranzasPrisma(request reportedtos.RequestPagosPeriodo) (cobranzas []reportedtos.DetallesPagosCobranza, erro error)
	GetCobranzasApilink(request reportedtos.RequestPagosPeriodo) (response []reportedtos.DetallesPagosCobranza, erro error)
	GetCobranzasRapipago(request reportedtos.RequestPagosPeriodo) (response []reportedtos.DetallesPagosCobranza, erro error)
	GetCobranzasMultipago(request reportedtos.RequestPagosPeriodo) (response []reportedtos.DetallesPagosCobranza, erro error)

	GetCierreLotePrisma(lista []string) (result []entities.Prismacierrelote, erro error)
	GetCierreLoteOffline(lista []string) (result []entities.Rapipagocierrelotedetalles, erro error)
	GetCierreLoteApilink(lista []string) (result []entities.Apilinkcierrelote, erro error)
	// GetCierreLoteApilinkReportes(requets filtros_reportes.ReportesFiltroApilink) (result []entities.Apilinkcierrelote, erro error)
	GetMovimiento(request reportedtos.RequestPagosPeriodo) (pagos []entities.Movimiento, erro error)
	GetRendicionReportes(request reportedtos.RequestPagosPeriodo) (pagos []entities.Movimiento, erro error)
	GetReversionesReportes(request reportedtos.RequestPagosPeriodo, filtroValidacion reportedtos.ValidacionesFiltro) (pagos []entities.Reversione, erro error)
	GetReversionesDeTransferenciaClientes(request reportedtos.RequestPagosPeriodo) (transferencias []entities.Transferencia, erro error)
	GetTransferenciasReportes(request reportedtos.RequestPagosPeriodo) (transferencias []entities.Transferencia, erro error)
	GetPeticionesReportes(request reportedtos.RequestPeticiones) (peticiones []entities.Webservicespeticione, total int64, erro error)
	GetLogs(request reportedtos.RequestLogs) (logs []entities.Log, total int64, erro error)
	GetNotificaciones(request reportedtos.RequestNotificaciones) (notificaciones []entities.Notificacione, total int64, erro error)
	SaveLotes(ctx context.Context, lotes []entities.Movimientolotes) (erro error)
	BajaMovimientoLotes(ctx context.Context, movimientos []entities.Movimientolotes, motivo_baja string) error

	//
	GetPagosBatch(request reportedtos.RequestPagosPeriodo) (pagos []entities.Pago, erro error)
	GetLotes(request reportedtos.RequestPagosPeriodo) (lote []entities.Pagolotes, erro error)  // obtener datos del lote
	GetLastLote(request reportedtos.RequestPagosPeriodo) (lote entities.Pagolotes, erro error) // Retorna ultimo lote
	GetCantidadLotes(request reportedtos.RequestPagosPeriodo) (lote int64, erro error)
	SavePagosLotes(ctx context.Context, lotes []entities.Pagolotes) (erro error)
	BajaPagosLotes(ctx context.Context, pagos []entities.Pagolotes, motivo_baja string) error

	// generar orden de liquidacion
	SaveLiquidacion(movliquidacion entities.Movimientoliquidaciones) (id uint64, erro error)
	/* REPORTES MOVIMIENTOS-COMISIONES */
	MovimientosComisionesRepository(filtro filtros_reportes.MovimientosComisionesFiltro) (response []reportedtos.ReporteMovimientosComisiones, total []reportedtos.ReporteMovimientosComisiones, erro error)

	/* REPORTES COBRANZAS-CLIENTES */
	CobranzasClientesRepository(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error)

	// Guardar datos reportes de clientes : NOTE esto solo se registra para control en el envio
	SaveGuardarDatosReporte(reporte entities.Reporte) (erro error)
	// Reportes enviados a clientes
	GetReportesEnviadosRepository(request reportedtos.RequestReportesEnviados) (listaReportes []entities.Reporte, totalFilas int64, erro error)

	GetLastReporteEnviadosRepository(request entities.Reporte, control filtros_reportes.BusquedaReporteFiltro) (siguiente uint, erro error)

	GetComprobantesRepository(request reportedtos.RequestRRComprobante) (response []entities.Comprobante, erro error)

	UpdateComprobanteRepository(entity entities.Comprobante) (entities.Comprobante, error)

	GetReportesRepository(request reportedtos.RequestGetReportes) (response []entities.Reporte, totalFilas int64, erro error)

	GetComprobantesByRrmIdRepository(request reportedtos.RequestRRComprobante) (response []entities.Comprobante, erro error)

	DeleteRetencionesAnticipadasRepository() (response []entities.MovimientoRetencion, erro error)

	GetMovimientoByIds(request reportedtos.RequestPagosPeriodo) (movimientos []entities.Movimiento, erro error)

	GetCuentaByApiKeyRepository(apikey string) (cuenta *entities.Cuenta, erro error)
	GetPagosItems(pagos_id []uint) ([]entities.Pagoitems, error)
}

type repository struct {
	SQLClient        *database.MySQLClient
	auditoriaService auditoria.AuditoriaService
	utilService      util.UtilService
}

func NewRepository(sqlClient *database.MySQLClient, a auditoria.AuditoriaService, t util.UtilService) ReportesRepository {
	return &repository{
		SQLClient:        sqlClient,
		auditoriaService: a,
		utilService:      t,
	}
}

func (r *repository) GetCuentaByApiKeyRepository(apikey string) (cuenta *entities.Cuenta, erro error) {
	resp := r.SQLClient.Model(entities.Cuenta{}).Where("apikey = ?", apikey)
	resp.Preload("Pagotipos")
	resp.Find(&cuenta)
	if resp.Error != nil {
		logs.Error("error al consultar cuenta: " + resp.Error.Error())
		erro = errors.New(ERROR_CONSULTAR_CUENTA)
		return
	}
	if resp.RowsAffected <= 0 {
		logs.Error("no existe cuenta")
		erro = errors.New(ERROR_CONSULTAR_CUENTA)
		return
	}
	return
}

func (r *repository) GetReportesRepository(request reportedtos.RequestGetReportes) (response []entities.Reporte, totalFilas int64, erro error) {
	queryGorm := r.SQLClient.Model(entities.Reporte{})

	if request.FechaInicio != "" && request.FechaFin != "" {
		queryGorm.Unscoped().Where("reportes.created_at BETWEEN ? AND ?", request.FechaInicio, request.FechaFin)
	}

	if request.NroReporte != 0 {
		queryGorm.Where("reportes.nro_reporte = ?", request.NroReporte)
	}

	if request.TipoReporte != "todos" {
		queryGorm.Where("reportes.tiporeporte = ?", request.TipoReporte)
	}

	if len(request.Cliente) != 0 {
		queryGorm.Where("reportes.cliente LIKE ?", "%"+request.Cliente+"%")
	}

	if !request.PeriodoInicio.IsZero() && !request.PeriodoFin.IsZero() {
		queryGorm.Where("reportes.periodo_inicio = ? AND reportes.periodo_fin = ?", request.PeriodoInicio, request.PeriodoFin)
	}

	if request.LastRrm {
		queryGorm.Order("fecharendicion desc").Limit(1)
	}

	// Paginacion
	if request.Number > 0 && request.Size > 0 {

		// Ejecutar y contar las filas devueltas
		queryGorm.Count(&totalFilas)

		if queryGorm.Error != nil {
			erro = fmt.Errorf("no se pudo cargar el total de filas de la consulta")
			return
		}

		offset := (request.Number - 1) * request.Size
		queryGorm.Limit(int(request.Size))
		queryGorm.Offset(int(offset))
	}

	queryGorm.Find(&response)

	// capturar error query DB
	if queryGorm.Error != nil {

		erro = fmt.Errorf("repositorio: no se puedieron obtener los registros de reportes")

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       queryGorm.Error.Error(),
			Funcionalidad: "GetReportesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), queryGorm.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetPagosReportes(request reportedtos.RequestPagosPeriodo, pendiente uint) (pagos []entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	if !request.FechaFin.IsZero() {
		resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
			Where("cast(pint.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	// if request.BuscarFechaPagointento {
	// 	if !request.FechaFin.IsZero() {
	// 		resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
	// 			Where("cast(pint.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	// 	}
	// }

	resp.Where("pagoestados_id != ?", pendiente)
	resp.Preload("PagoEstados")

	if len(request.PagoEstados) > 0 {
		resp.Where("pagoestados_id IN ?", request.PagoEstados)
	}

	if len(request.ClientesIds) > 0 {
		resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id in ?", request.ClientesIds).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id in ?", request.ClientesIds)
	} else {
		if request.ClienteId > 0 {
			resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
		} else if len(request.ApiKey) > 0 {
			resp.Preload("PagosTipo.Cuenta", "cuentas.apikey = ?", request.ApiKey).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id").Where("c.apikey = ?", request.ApiKey)
		} else {
			resp.Preload("PagosTipo.Cuenta.Cliente").Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Order("cli.id DESC")
		}
	}
	resp.Preload("PagoIntentos.Mediopagos.Channel.Channelaranceles").Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id").
		Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO").
		Order("pi.created_at DESC")
	resp.Preload("PagoIntentos.Installmentdetail")
	resp.Preload("PagoIntentos.Movimientos.Movimientocomisions")
	resp.Preload("PagoIntentos.Movimientos.Movimientoimpuestos")
	resp.Preload("PagoIntentos.Movimientos.Movimientotransferencia")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientocomisions")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientoimpuestos")
	resp.Preload("PagoIntentos.Movimientotemporale.Movimientoretenciontemporales")

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCobranzasPrisma(request reportedtos.RequestPagosPeriodo) (cobranzas []reportedtos.DetallesPagosCobranza, erro error) {
	resp := r.constructQuery(request, "", "PI.paid_at")
	resp.Where("PI.card_last_four_digits != '' AND PI.paid_at IS NOT NULL")
	resp.Find(&cobranzas)

	return
}

func (r *repository) GetCobranzasRapipago(request reportedtos.RequestPagosPeriodo) (response []reportedtos.DetallesPagosCobranza, erro error) {
	resp := r.constructQuery(request, "INNER JOIN rapipagocierrelotedetalles RP ON RP.codigo_barras = PI.barcode", "RP.fecha_cobro")
	resp.Find(&response)

	return
}

func (r *repository) GetCobranzasMultipago(request reportedtos.RequestPagosPeriodo) (response []reportedtos.DetallesPagosCobranza, erro error) {
	resp := r.constructQuery(request, "INNER JOIN multipagoscierrelotedetalles MTD ON MTD.codigo_barras = PI.barcode", "MTD.fecha_cobro")
	resp.Find(&response)

	return
}

func (r *repository) GetCobranzasApilink(request reportedtos.RequestPagosPeriodo) (response []reportedtos.DetallesPagosCobranza, erro error) {
	resp := r.constructQuery(request, "INNER JOIN apilinkcierrelotes ACL ON ACL.debin_id = PI.external_id", "ACL.fecha_cobro")

	resp.Find(&response)

	return
}

func (r *repository) constructQuery(filtro reportedtos.RequestPagosPeriodo, joinTable string, fechaCobroColumn string) *gorm.DB {
	resp := r.SQLClient.Model(entities.Pago{})

	resp.Joins("INNER JOIN pagotipos PT ON pagos.pagostipo_id = PT.id")
	resp.Joins("INNER JOIN cuentas C ON C.id = PT.cuentas_id")
	resp.Joins("INNER JOIN clientes CTS ON C.clientes_id = CTS.id")
	resp.Joins("INNER JOIN pagoestados PE ON pagos.pagoestados_id = PE.id")
	resp.Joins("INNER JOIN pagointentos PI ON pagos.id = PI.pagos_id")
	resp.Joins("INNER JOIN mediopagos MP ON MP.id = PI.mediopagos_id")
	resp.Joins("INNER JOIN channels CH ON CH.id = MP.channels_id")

	resp.Joins("LEFT JOIN movimientotemporales MT ON MT.pagointentos_id = PI.id")
	resp.Joins("LEFT JOIN movimientocomisionetemporales MCT ON MCT.movimientotemporales_id = MT.id")
	resp.Joins("LEFT JOIN movimientoimpuestotemporales MIT ON MIT.movimientotemporales_id = MT.id")
	resp.Joins("LEFT JOIN movimiento_retenciontemporales MRT ON MRT.movimientotemporales_id = MT.id")
	resp.Joins("LEFT JOIN pagolotes PL ON PL.pagos_id = pagos.id")

	if len(joinTable) > 0 {
		resp.Joins(joinTable)
	}

	querySelect := `
		pagos.id, 
		pagos.payer_name, 
		pagos.payer_email, 
		pagos.first_total as total_pago, 
		pagos.external_reference as referencia , 
		PE.nombre as Pagoestado ,
		pagos.description AS descripcion , 
		C.cuenta, 
		CTS.cliente, 
		cast(PI.paid_at as date)  as fecha_pago, 
		MP.mediopago as medio_pago, 
		CH.nombre as canal_pago, 
		(MIT.monto + MIT.montoproveedor) AS iva,
		(MCT.monto + MCT.montoproveedor) AS comision_total, 
		MCT.monto AS comision,
		SUM(MRT.importe_retenido) AS retencion,
		PL.lote, 
		` + fechaCobroColumn

	resp.Select(querySelect)
	resp.Where("DATE("+fechaCobroColumn+") BETWEEN DATE(?) AND DATE(?)", filtro.FechaInicio, filtro.FechaFin)
	resp.Group("PI.id")

	if filtro.ClienteId != 0 && filtro.CuentaId != 0 {
		resp.Where("CTS.id = ?", filtro.ClienteId)
		resp.Where("C.id = ?", filtro.CuentaId)
	} else {
		if filtro.ClienteId != 0 {
			resp.Where("CTS.id = ?", filtro.ClienteId)
		}
		if filtro.CuentaId != 0 {
			resp.Joins("INNER JOIN cuentas Cu ON Cu.id = PT.cuentas_id")
			resp.Where("Cu.id = ?", filtro.CuentaId)
		}
	}

	return resp
}

func (r *repository) GetPagosItems(pagos_id []uint) ([]entities.Pagoitems, error) {
	var pago_items []entities.Pagoitems
	resp := r.SQLClient.Table("pagoitems")

	if len(pagos_id) > 0 {
		resp.Where("pagos_id IN ?", pagos_id)
	}

	resp.Find(&pago_items)

	if resp.Error != nil {
		return pago_items, resp.Error
	}

	return pago_items, nil
}

func (r *repository) GetCierreLotePrisma(lista []string) (result []entities.Prismacierrelote, erro error) {

	resp := r.SQLClient.Model(entities.Prismacierrelote{})

	resp.Unscoped()

	resp.Where("externalcliente_id IN ?", lista)

	resp.Preload("Prismamovimientodetalle.MovimientoCabecera")

	resp.Preload("Prismatrdospagos")
	resp.Order("fecha_pago")
	resp.Find(&result)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_PRISMA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLotePrisma",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCierreLoteOffline(lista []string) (result []entities.Rapipagocierrelotedetalles, erro error) {

	// resp := r.SQLClient.Model(entities.Rapipagocierrelotedetalles{}).Where("codigo_barras IN ?", lista)

	resp := r.SQLClient.Unscoped().Where("codigo_barras IN ?", lista)

	resp.Preload("RapipagoCabecera", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	})

	// resp.Preload("RapipagoCabecera")

	resp.Find(&result)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_OFFLINE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLoteOffline",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

// func (r *repository) GetCierreLoteApilinkReportes(requets filtros_reportes.ReportesFiltroApilink) (result []entities.Apilinkcierrelote, erro error) {

// 	resp := r.SQLClient.Unscoped()

// 	if requets.Conciliado {
// 		resp.Where("banco_external_id != ?", 0)
// 	}

// 	if !requets.Informado {
// 		resp.Where("pagoinformado = ?", requets.Informado)
// 	}

// 	resp.Find(&result)

// 	if resp.Error != nil {

// 		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_APILINK)

// 		log := entities.Log{
// 			Tipo:          entities.Error,
// 			Mensaje:       resp.Error.Error(),
// 			Funcionalidad: "GetCierreLoteApilinkReportes",
// 		}

// 		err := r.utilService.CreateLogService(log)

// 		if err != nil {
// 			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
// 			logs.Error(mensaje)
// 		}
// 	}

// 	return
// }

func (r *repository) GetCierreLoteApilink(lista []string) (result []entities.Apilinkcierrelote, erro error) {

	resp := r.SQLClient.Unscoped().Where("debin_id IN ?", lista)

	resp.Find(&result)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_CIERRELOTE_APILINK)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLoteApilink",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetRendicionReportes(request reportedtos.RequestPagosPeriodo) (pagos []entities.Movimiento, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	if !request.FechaFin.IsZero() {
		resp.Where("cast(movimientos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if request.CuentaId > 0 {
		resp.Preload("Cuenta", "cuentas.id = ?", request.CuentaId).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id").Where("c.id = ?", request.CuentaId)
	}

	if request.ClienteId > 0 {
		resp.Preload("Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
	}

	if request.PagoIntento > 0 {
		resp.Where("pagointentos_id = ?", request.PagoIntento)
	}

	if len(request.TipoMovimiento) > 0 {
		resp.Where("tipo = ?", request.TipoMovimiento)
	}

	if request.CargarMedioPago {
		resp.Preload("Pagointentos.Mediopagos")
	}

	resp.Preload("Pagointentos.Pago.Pagoitems")
	resp.Preload("Movimientotransferencia")
	resp.Preload("Movimientocomisions")
	resp.Preload("Movimientoimpuestos")
	// resp.Preload("Movimientolotes")

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRendicionReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}
func (r *repository) GetReversionesDeTransferenciaClientes(request reportedtos.RequestPagosPeriodo) (transferencias []entities.Transferencia, erro error) {

	resp := r.SQLClient.Model(entities.Transferencia{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(transferencias.fecha_operacion as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	if request.ClienteId > 0 {
		resp.Preload("Movimiento.Cuenta", "cuentas.clientes_id = ?", request.ClienteId).
			Joins("INNER JOIN movimientos as mov on mov.id = transferencias.movimientos_id INNER JOIN pagointentos as pi on pi.id = mov.pagointentos_id INNER JOIN pagos as p on p.id = pi.pagos_id INNER JOIN cuentas as cu on cu.id = mov.cuentas_id").
			Where("cu.clientes_id = ?", request.ClienteId)
	}
	if request.CargarReversion {
		resp.Where("transferencias.reversion = ?", 1)
	}

	resp.Preload("Movimiento.Pagointentos")
	resp.Preload("Movimiento.Cuenta.Cliente")
	resp.Preload("Movimiento.Pagointentos.Pago")
	resp.Preload("Movimiento.Pagointentos.Pago.Pagoitems")
	resp.Preload("Movimiento.Pagointentos.Pago.PagoEstados")
	resp.Preload("Movimiento.Pagointentos.Mediopagos")

	if request.OrdenadoFecha {
		resp.Order("id DESC")
	}

	resp.Find(&transferencias)
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_TRANSFERENCIAS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetReversionesDeTransferenciaClientes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetReversionesReportes(request reportedtos.RequestPagosPeriodo, filtroValidacion reportedtos.ValidacionesFiltro) (pagos []entities.Reversione, erro error) {

	resp := r.SQLClient.Model(entities.Reversione{})

	if filtroValidacion.Fin && filtroValidacion.Inicio {
		resp.Where("cast(reversiones.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	if request.ClienteId > 0 {
		resp.Preload("PagoIntento.Pago.PagosTipo.Cuenta", "cuentas.clientes_id = ?", request.ClienteId).
			Joins("INNER JOIN pagointentos as pt on pt.id = reversiones.pagointentos_id INNER JOIN pagos as p on p.id = pt.pagos_id INNER JOIN pagotipos as pa on pa.id = p.pagostipo_id INNER JOIN cuentas as cu on cu.id = pa.cuentas_id").
			Where("cu.clientes_id = ?", request.ClienteId)
	}

	resp.Preload("PagoIntento.Mediopagos")
	resp.Preload("PagoIntento.Pago.PagosTipo.Cuenta.Cliente")
	resp.Preload("PagoIntento.Pago.Pagoitems")
	resp.Preload("PagoIntento.Pago.PagoEstados")

	if request.OrdenadoFecha {
		resp.Order("id DESC")
	}

	resp.Find(&pagos)
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetReversionesReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetTransferenciasReportes(request reportedtos.RequestPagosPeriodo) (transferencias []entities.Transferencia, erro error) {

	resp := r.SQLClient.Model(entities.Transferencia{})

	if !request.FechaFin.IsZero() {
		resp.Where("cast(transferencias.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	if request.ClienteId > 0 {
		resp.Preload("Movimiento.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN movimientos as mv on mv.id = transferencias.movimientos_id INNER JOIN cuentas as c on c.id = mv.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
	}

	if request.CuentaId > 0 {
		resp.Where("c.id = ?", request.CuentaId)
	}

	if len(request.ApiKey) > 0 {
		resp.Preload("Movimiento.Cuenta", "cuentas.apikey = ?", request.ApiKey).Joins("INNER JOIN movimientos as mv on mv.id = transferencias.movimientos_id INNER JOIN cuentas as c on c.id = mv.cuentas_id").Where("c.apikey = ?", request.ApiKey)
	}
	// resp.Preload("Movimiento.Pagointentos.Pago.Pagoitems")
	// resp.Preload("Movimiento.Movimientocomisions")
	// resp.Preload("Movimiento.Movimientoimpuestos")
	// resp.Preload("Movimientolotes")

	if request.OrdenadoFecha {
		resp.Order("fecha_operacion DESC")
	}

	resp.Find(&transferencias)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetTransferenciasReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetPeticionesReportes(request reportedtos.RequestPeticiones) (peticiones []entities.Webservicespeticione, total int64, erro error) {

	resp := r.SQLClient.Model(entities.Webservicespeticione{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	resp.Where("vendor = ?", request.Vendor)
	resp.Where("operacion != ?", "Autenticacion(genera token)")

	resp.Count(&total)
	if request.Number > 0 && request.Size > 0 {

		offset := (request.Number - 1) * request.Size
		resp.Limit(int(request.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&peticiones)

	return
}

func (r *repository) GetLogs(request reportedtos.RequestLogs) (logs []entities.Log, total int64, erro error) {

	resp := r.SQLClient.Model(entities.Log{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	resp.Order("created_at DESC")
	resp.Count(&total)
	if request.Number > 0 && request.Size > 0 {

		offset := (request.Number - 1) * request.Size
		resp.Limit(int(request.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&logs)

	return
}

func (r *repository) GetNotificaciones(request reportedtos.RequestNotificaciones) (notificaciones []entities.Notificacione, total int64, erro error) {

	resp := r.SQLClient.Model(entities.Notificacione{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}
	resp.Order("created_at DESC")
	resp.Count(&total)
	if request.Number > 0 && request.Size > 0 {

		offset := (request.Number - 1) * request.Size
		resp.Limit(int(request.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&notificaciones)

	return
}

func (r *repository) GetLastLote(request reportedtos.RequestPagosPeriodo) (lote entities.Pagolotes, erro error) {

	resp := r.SQLClient.Model(entities.Pagolotes{})

	if request.ClienteId > 0 {
		resp.Where("clientes_id = ?", request.ClienteId)
	}
	resp.Last(&lote)

	if resp.RowsAffected <= 0 {
		lote = entities.Pagolotes{}
	}
	return
}

func (r *repository) SaveLotes(ctx context.Context, lotes []entities.Movimientolotes) (erro error) {
	//Si no se realiza toda la operación entonces vuelve todo a como estaba antes de empezar.
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		// 1 - creo los movimientos lotes
		if len(lotes) > 0 {
			res := tx.WithContext(ctx).Create(&lotes)
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_GUARDAR_LOTES)
			}
		}
		return nil
	})
}

func (r *repository) BajaMovimientoLotes(ctx context.Context, movimientos []entities.Movimientolotes, motivo_baja string) error {

	resp := r.SQLClient.WithContext(ctx).Model(&movimientos).Omit(clause.Associations).UpdateColumns(map[string]interface{}{"updated_at": time.Now(), "deleted_at": time.Now(), "motivo_baja": motivo_baja})

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro := fmt.Errorf(ERROR_BAJAR_MOVIMIENTOS_LOTES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "BajaMovimientoLotes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
		return erro
	}
	return nil
}

func (r *repository) GetMovimiento(request reportedtos.RequestPagosPeriodo) (movimientos []entities.Movimiento, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	if !request.FechaInicio.IsZero() && !request.FechaFin.IsZero() {
		resp.Where("cast(movimientos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if request.PagoIntento > 0 {
		resp.Where("pagointentos_id = ?", request.PagoIntento)
	}

	if len(request.PagoIntentos) > 0 {
		resp.Where("pagointentos_id IN ?", request.PagoIntentos)
	}

	if len(request.ApiKey) > 0 {
		resp.Preload("Cuenta", "cuentas.apikey = ?", request.ApiKey).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id").Where("c.apikey = ?", request.ApiKey)
	}

	if request.CargarMedioPago {
		resp.Preload("Pagointentos.Mediopagos")
	}

	if len(request.TipoMovimiento) > 0 {
		resp.Where("tipo = ?", request.TipoMovimiento)
	}

	if request.CargarReversionReporte {
		resp.Where("reversion = ?", true)
		resp.Where("monto < ?", 0)
	}

	if request.CargarReversion {
		resp.Where("reversion = ?", true)
		resp.Where("cast(movimientos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if request.CargarCuenta {
		resp.Preload("Pagointentos.Pago.Pagoitems")
	}

	if request.CargarCuenta {
		resp.Preload("Cuenta")
	}

	if request.CargarComisionImpuesto {
		resp.Preload("Movimientocomisions")
		resp.Preload("Movimientoimpuestos")
	}

	if request.CargarMovimientosTransferencias {
		resp.Preload("Movimientotransferencia")
	}

	if request.CargarCliente {
		resp.Preload("Cuenta.Cliente")
	}

	if request.ClienteId != 0 {
		resp.Joins("JOIN cuentas ON movimientos.cuentas_id = cuentas.id")
		resp.Joins("JOIN clientes ON clientes.id = cuentas.clientes_id")
		resp.Where("clientes.id = ? ", request.ClienteId)
	}

	if request.OrdenadoFecha {
		resp.Order("pagointentos_id Desc")
	}

	// cargar MovimientoRetencion
	if request.CargarMovimientosRetenciones {
		resp.Preload("Movimientoretencions.Retencion.Condicion.Gravamen")
	}

	// cargar solo movimientos que tengan MovimientoRetencion
	if request.CaragarSoloMovimientoRetencion {
		resp.Preload("Movimientoretencions.Retencion.Condicion.Gravamen").Joins("JOIN movimiento_retencions ON movimientos.id = movimiento_retencions.movimiento_id").Distinct()
	}

	if request.CargarPagoIntentos {
		resp.Preload("Pagointentos")
	}

	resp.Find(&movimientos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_MOVIMIENTOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMovimiento",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetMovimientoByIds(request reportedtos.RequestPagosPeriodo) (movimientos []entities.Movimiento, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	if len(request.IdsMovimientos) > 0 {
		resp.Where("movimientos.id IN ?", request.IdsMovimientos)
	}

	if request.CargarCliente {
		resp.Preload("Cuenta.Cliente")
	}

	// cargar solo movimientos que tengan MovimientoRetencion
	if request.CaragarSoloMovimientoRetencion {
		resp.Preload("Movimientoretencions.Retencion.Condicion.Gravamen").Joins("JOIN movimiento_retencions ON movimientos.id = movimiento_retencions.movimiento_id").Distinct()
	}

	if request.CargarPagoIntentos {
		resp.Preload("Pagointentos")
	}

	resp.Find(&movimientos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_MOVIMIENTOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMovimiento",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

// pagos items batch
func (r *repository) GetPagosBatch(request reportedtos.RequestPagosPeriodo) (pagos []entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	// se cambio la busqueda al campo paid_at de pago intento
	// if !request.FechaFin.IsZero() {
	// 	resp.Where("cast(pagos.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	// }

	if !request.FechaFin.IsZero() {
		resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pint on pagos.id = pint.pagos_id").
			Where("cast(pint.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if len(request.PagoEstados) > 0 {
		resp.Where("pagoestados_id IN ?", request.PagoEstados)
	}

	if request.ClienteId > 0 {
		// resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
		resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as cu on cu.id = pt.cuentas_id").Where("cu.clientes_id = ?", request.ClienteId)
	}
	// resp.Preload("PagoIntentos")
	// buscar solo sobre el pago intento aprobado
	resp.Preload("PagoIntentos").Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id").
		Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO")

	resp.Preload("Pagoitems")
	resp.Preload("Pagolotes")

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRendicionReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCantidadLotes(request reportedtos.RequestPagosPeriodo) (lote int64, erro error) {

	var pagoslotes []entities.Pagolotes
	resp := r.SQLClient.Model(entities.Pagolotes{})

	if !request.FechaFin.IsZero() {
		resp.Where("cast(pagolotes.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", request.FechaInicio, request.FechaFin)
	}

	if len(request.Pagos) > 0 {
		resp.Where("pagos_id IN ?", request.Pagos)
	}

	// if request.ClienteId > 0 {
	// 	// resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN cuentas as c on c.id = movimientos.cuentas_id INNER JOIN clientes as cli on cli.id = c.clientes_id").Where("cli.id = ?", request.ClienteId)
	// 	resp.Preload("PagosTipo.Cuenta.Cliente", "clientes.id = ?", request.ClienteId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as cu on cu.id = pt.cuentas_id").Where("cu.clientes_id = ?", request.ClienteId)
	// }

	resp.Find(&pagoslotes)

	lote = resp.RowsAffected

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRendicionReportes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

// func (r *repository) GetPagosLotes(request reportedtos.RequestPagosPeriodo) (lote entities.Pagolotes, erro error) {

// 	resp := r.SQLClient.Model(entities.Pagolotes{})

// 	if request.ClienteId > 0 {
// 		resp.Where("clientes_id = ?", request.ClienteId)
// 	}
// 	resp.Last(&lote)

// 	if resp.RowsAffected <= 0 {
// 		lote = entities.Pagolotes{}
// 	}
// 	return
// }

func (r *repository) GetLotes(request reportedtos.RequestPagosPeriodo) (lote []entities.Pagolotes, erro error) {

	resp := r.SQLClient.Model(entities.Pagolotes{})

	if len(request.Pagos) > 0 {
		resp.Where("pagos_id = ?", request.Pagos)
	}

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONSULTAR_PAGOS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetLotes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	resp.Find(&lote)
	return
}

func (r *repository) SavePagosLotes(ctx context.Context, lotes []entities.Pagolotes) (erro error) {
	//Si no se realiza toda la operación entonces vuelve todo a como estaba antes de empezar.
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		// 1 - creo los movimientos lotes
		if len(lotes) > 0 {
			res := tx.WithContext(ctx).Create(&lotes)
			if res.Error != nil {
				logs.Info(res.Error)
				return errors.New(ERROR_GUARDAR_LOTES)
			}
		}
		return nil
	})
}

func (r *repository) BajaPagosLotes(ctx context.Context, pagos []entities.Pagolotes, motivo_baja string) error {

	resp := r.SQLClient.WithContext(ctx).Model(&pagos).Omit(clause.Associations).UpdateColumns(map[string]interface{}{"updated_at": time.Now(), "deleted_at": time.Now(), "motivo_baja": motivo_baja})

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro := fmt.Errorf(ERROR_BAJAR_MOVIMIENTOS_LOTES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "BajaPagosLotes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			logs.Error(err)
		}
		return erro
	}
	return nil
}

func (r *repository) SaveLiquidacion(movliquidacion entities.Movimientoliquidaciones) (id uint64, erro error) {

	result := r.SQLClient.Create(&movliquidacion)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_MOVIMIENTOS_LIQUIDACIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateCliente",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}
	id = uint64(movliquidacion.ID)

	return
}

func (r *repository) MovimientosComisionesRepository(filtro filtros_reportes.MovimientosComisionesFiltro) (response []reportedtos.ReporteMovimientosComisiones, total []reportedtos.ReporteMovimientosComisiones, erro error) {

	resp := r.SQLClient.Model(entities.Movimiento{})

	resp.Joins("INNER JOIN pagointentos PI ON movimientos.pagointentos_id = PI.id INNER JOIN pagos P ON PI.pagos_id = P.id INNER JOIN movimientoimpuestos MI ON movimientos.id = MI.movimientos_id INNER JOIN movimientocomisiones AS MC ON movimientos.id = MC.movimientos_id INNER JOIN cuentacomisions CC ON MC.cuentacomisions_id = CC.id INNER JOIN cuentas C ON C.id = movimientos.cuentas_id INNER JOIN clientes CTS ON C.clientes_id = CTS.id")
	resp.Select("movimientos.id AS id,case movimientos.tipo when 'C' then 'CREDITO' else 'DEBITO' end AS tipo, PI.amount AS monto_pago, movimientos.monto AS monto_movimiento, (MC.monto + MC.montoproveedor) AS monto_comision, (MC.porcentaje + MC.porcentajeproveedor) AS porcentaje_comision , (MI.monto + MI.montoproveedor) AS monto_impuesto, MI.porcentaje AS porcentaje_impuesto, MC.montoproveedor AS monto_comisionproveedor,MC.porcentajeproveedor AS porcentaje_comisionproveedor , MI.montoproveedor AS monto_impuestoproveedor, CTS.cliente AS nombre_cliente, movimientos.created_at")

	resp.Where("movimientos.created_at BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)

	if filtro.ClienteId != 0 {
		resp.Where("CTS.id = ?", filtro.ClienteId)
		if filtro.CuentaId != 0 {
			resp.Where("C.id = ?", filtro.CuentaId)
		}
	}

	resp.Order("movimientos.id DESC")

	resp.Find(&total)
	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	// manejo y log del error en la consulta a base de datos
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_MOVIMIENTOS_COMISIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "MovimientosComisionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) CobranzasClientesRepository(filtro filtros_reportes.CobranzasClienteFiltro) (response []reportedtos.DetallesPagosCobranza, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	resp.Joins("INNER JOIN pagotipos PT ON pagos.pagostipo_id = PT.id INNER JOIN cuentas C ON C.id = PT.cuentas_id INNER JOIN clientes CTS ON C.clientes_id = CTS.id INNER JOIN pagoestados PE ON pagos.pagoestados_id = PE.id")
	resp.Joins("INNER JOIN pagointentos PI ON pagos.id = PI.pagos_id INNER JOIN mediopagos MP ON MP.id = PI.mediopagos_id INNER JOIN channels CH ON CH.id = MP.channels_id")
	resp.Select("pagos.*, pagos.first_total as total_pago, pagos.external_reference as referencia , PE.nombre as Pagoestado ,pagos.description AS descripcion , C.cuenta, CTS.cliente, cast(PI.paid_at as date)  as fecha_pago, MP.mediopago as medio_pago, CH.nombre as canal_pago")

	resp.Where("pagos.created_at BETWEEN ? AND ?", filtro.FechaInicio, filtro.FechaFin)
	resp.Where("PE.nombre = ? OR PE.nombre = ?", "AUTORIZADO", "APROBADO")

	resp.Where("PI.id in ( SELECT  MAX(id) FROM pagointentos as PI GROUP BY PI.pagos_id	)")

	if filtro.ClienteId != 0 {
		resp.Where("CTS.id = ?", filtro.ClienteId)
	}

	resp.Order("pagos.id DESC")

	resp.Find(&response)

	// manejo y log del error en la consulta a base de datos
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_COBRANZAS_CLIENTES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "CobranzasClientesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) SaveGuardarDatosReporte(reporte entities.Reporte) (erro error) {

	result := r.SQLClient.Create(&reporte)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_REGISTRO_REPORTE)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "SaveGuardarDatosReporte",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	return
}

func (r *repository) GetReportesEnviadosRepository(request reportedtos.RequestReportesEnviados) (listaReportes []entities.Reporte, totalFilas int64, erro error) {

	queryGorm := r.SQLClient.Model(entities.Reporte{})

	queryGorm.Unscoped().Where("created_at BETWEEN ? AND ?", request.FechaInicio, request.FechaFin)

	// se checkea el filtro TipoReporte
	if request.TipoReporte != "todos" {
		queryGorm.Where("tiporeporte = ?", request.TipoReporte)
	}

	// filtro por cliente
	if len(request.Cliente) != 0 {
		queryGorm.Where("cliente LIKE ?", "%"+request.Cliente+"%")
	}

	queryGorm.Preload("Reportedetalle.Pago.PagoIntentos.Movimientotemporale.Movimientocomisions")
	queryGorm.Preload("Reportedetalle.Pago.PagoIntentos.Movimientotemporale.Movimientoimpuestos")
	queryGorm.Preload("Reportedetalle.Pago.PagoIntentos.Movimientotemporale.Movimientoretenciontemporales")
	// queryGorm.Preload("Reportedetalle.Pago")

	// Paginacion
	if request.Number > 0 && request.Size > 0 {

		// Ejecutar y contar las filas devueltas
		queryGorm.Count(&totalFilas)

		if queryGorm.Error != nil {
			erro = fmt.Errorf("no se pudo cargar el total de filas de la consulta")
			return
		}

		offset := (request.Number - 1) * request.Size
		queryGorm.Limit(int(request.Size))
		queryGorm.Offset(int(offset))
	}

	queryGorm.Order("created_at desc").Find(&listaReportes)

	// capturar error query DB
	if queryGorm.Error != nil {

		erro = fmt.Errorf("repositorio: no se puedieron obtener los registros de reportes enviados")

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       queryGorm.Error.Error(),
			Funcionalidad: "GetReportesEnviadosRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), queryGorm.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetLastReporteEnviadosRepository(request entities.Reporte, control filtros_reportes.BusquedaReporteFiltro) (siguiente uint, erro error) {

	if !control.SigNumero {

		var controlEntity []entities.Reporte
		queryGorm := r.SQLClient.Model(entities.Reporte{})

		queryGorm.Where("tiporeporte = ?", request.Tiporeporte)

		if len(request.Cliente) != 0 {
			queryGorm.Where("cliente = ?", request.Cliente)
		}

		if request.Fechacobranza != "" {
			queryGorm.Where("fechacobranza = ?", request.Fechacobranza)
		}

		if request.Fecharendicion != "" {
			queryGorm.Where("fecharendicion = ?", request.Fecharendicion)
		}

		// Busca el último si no encuentra coincidencia

		queryGorm.Order("created_at asc")

		queryGorm.Find(&controlEntity)

		if len(controlEntity) > 0 {
			siguiente = controlEntity[0].Nro_reporte
			return
		}

		if control.Control {
			return
		}
	}
	var lastReporte entities.Reporte
	queryLast := r.SQLClient.Model(entities.Reporte{})

	queryLast.Where("tiporeporte = ?", request.Tiporeporte)

	queryLast.Order("created_at desc")

	if control.LastRrm {
		queryLast.Order("reportes.nro_reporte desc")
	}

	queryLast.First(&lastReporte)

	siguiente = 1

	if queryLast.RowsAffected > 0 {
		siguiente = lastReporte.Nro_reporte + 1
	}

	return
}

func (r *repository) GetComprobantesRepository(request reportedtos.RequestRRComprobante) (response []entities.Comprobante, erro error) {
	zeroTime := time.Time{}
	resp := r.SQLClient.Model(entities.Comprobante{})

	resp.Where("comprobantes.emitido_el IS NULL OR comprobantes.emitido_el = ?", zeroTime)

	resp.Preload("Cliente.Contactosreportes")
	resp.Preload("ComprobanteDetalles")

	if request.Cliente_id != 0 {
		resp.Where("comprobantes.cliente_id = ? ", request.Cliente_id)
	}

	if request.ComprobanteId != 0 {
		resp.Where("comprobantes.id = ? ", request.ComprobanteId)
	}

	resp.Find(&response)

	// manejo y log del error en la consulta a base de datos
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_GET_COMPROBANTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetComprobantesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) UpdateComprobanteRepository(entity entities.Comprobante) (entities.Comprobante, error) {
	var erro error

	resp := r.SQLClient.Omit("created_at").Omit(clause.Associations).Updates(&entity)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_UPDATE_COMPROBANTE)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "UpdateComprobanteRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. UpdateComprobanteRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
		return entities.Comprobante{}, erro
	}

	if resp.RowsAffected <= 0 {
		erro = fmt.Errorf(ERROR_UPDATE_COMPROBANTE)
		return entities.Comprobante{}, erro
	}

	return entity, nil
}

func (r *repository) GetComprobantesByRrmIdRepository(request reportedtos.RequestRRComprobante) (response []entities.Comprobante, erro error) {
	// zeroTime := time.Time{}
	resp := r.SQLClient.Model(entities.Comprobante{})

	resp.Where("comprobantes.emitido_el IS NOT NULL")

	resp.Preload("ComprobanteDetalles")

	if request.Cliente_id != 0 {
		resp.Where("comprobantes.cliente_id = ? ", request.Cliente_id)
	}

	if request.ComprobanteId != 0 {
		resp.Where("comprobantes.id = ? ", request.ComprobanteId)
	}

	if request.ReporteId != 0 {
		resp.Where("comprobantes.reporte_id = ? ", request.ReporteId)
	}

	// filtro para diferenciar por tipo de gravamen en el comprobante
	if len(request.GravamenesIn) > 0 {
		resp.Where("comprobantes.gravamen IN ?", request.GravamenesIn)
	}

	resp.Find(&response)

	// manejo y log del error en la consulta a base de datos
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_GET_COMPROBANTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetComprobantesByRrmIdRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) DeleteRetencionesAnticipadasRepository() (response []entities.MovimientoRetencion, erro error) {

	subQuery := r.SQLClient.Table("movimiento_retencions").
		Select("movimiento_retencions.id").
		Joins("JOIN movimientos ON movimiento_retencions.movimiento_id = movimientos.id").
		Joins("JOIN pagointentos ON movimientos.pagointentos_id = pagointentos.id").
		Where("pagointentos.paid_at <> '0000-00-00 00:00:00' AND pagointentos.paid_at < '2023-11-01 00:00:00'")
	resp := r.SQLClient.Model(&response).Where("movimiento_retencions.id IN (?)", subQuery).Update("movimiento_retencions.deleted_at", time.Now())

	if resp.Error != nil {

		erro = fmt.Errorf("error al borrar retenciones anticipadas")

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "DeleteRetencionesAnticipadasRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}
