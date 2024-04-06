package administracion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos/retenciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/multipagos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	BeginTx()
	RollbackTx()
	CommitTx()

	//CUENTAS
	CuentaByClientePage(cliente int64, limit, offset int) (*[]entities.Cuenta, int64, error)
	GetCuenta(filtro filtros.CuentaFiltro) (cuenta entities.Cuenta, erro error)
	CuentaByID(cuenta int64) (*entities.Cuenta, error)
	SaveCuenta(ctx context.Context, cuenta *entities.Cuenta) (bool, error)
	UpdateCuenta(ctx context.Context, cuenta entities.Cuenta) (erro error)
	DeleteCuenta(id uint64) (erro error)
	SetApiKey(ctx context.Context, cuenta entities.Cuenta) (erro error)
	GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, erro error)
	GetPaymentByExternal(filtroPago filtros.PagoFiltro) (*entities.Pago, error)

	//CLIENTES
	CreateCliente(ctx context.Context, cliente entities.Cliente) (id uint64, erro error)
	UpdateCliente(ctx context.Context, cliente entities.Cliente) (erro error)
	DeleteCliente(ctx context.Context, id uint64) (erro error)
	GetCliente(filtro filtros.ClienteFiltro) (cliente entities.Cliente, erro error)
	GetClientes(filtro filtros.ClienteFiltro) (clientes []entities.Cliente, totalFilas int64, erro error)
	GetCuentasByCliente(clienteId uint64) (cuentas []entities.Cuenta, erro error)
	GetClientesConfiguracion(filtro filtros.ClienteConfiguracionFiltro) (clientes []entities.Cliente, erro error)

	//PAGOS
	GetPagosByUUID(uuid []string) (pagos []*entities.Pago, erro error)
	GetPagos(filtro filtros.PagoFiltro) (pagos []entities.Pago, totalFilas int64, erro error)
	GetPago(filtro filtros.PagoFiltro) (pago entities.Pago, erro error)
	GetPagosIntentos(filtro filtros.PagoIntentoFiltro) (pagos []entities.Pagointento, erro error)
	GetPagosEstados(filtro filtros.PagoEstadoFiltro) (estados []entities.Pagoestado, erro error)
	GetPagosEstadosExternos(filtro filtros.PagoEstadoExternoFiltro) (estados []entities.Pagoestadoexterno, erro error)
	GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estados entities.Pagoestado, erro error)
	PagoById(pagoID int64) (*entities.Pago, error)
	SavePagotipo(tipo *entities.Pagotipo) (bool, error)
	ConsultarEstadoPagosRepository(parametrosVslido administraciondtos.ParamsValidados, filtro filtros.PagoFiltro) (entityPagos []entities.Pago, erro error)

	//ABM RUBROS
	CreateRubro(ctx context.Context, rubro entities.Rubro) (id uint64, erro error)
	UpdateRubro(ctx context.Context, rubro entities.Rubro) (erro error)
	GetRubro(filtro filtros.RubroFiltro) (rubro entities.Rubro, erro error)
	GetRubros(filtro filtros.RubroFiltro) (rubros []entities.Rubro, totalFilas int64, erro error)

	//ABM PAGOS TIPOS
	CreatePagoTipo(ctx context.Context, request entities.Pagotipo, channels []int64, cuotas []string) (id uint64, erro error)
	UpdatePagoTipo(ctx context.Context, request entities.Pagotipo, channels administraciondtos.RequestPagoTipoChannels, cuotas administraciondtos.RequestPagoTipoCuotas) (erro error)
	GetPagoTipo(filtro filtros.PagoTipoFiltro) (response entities.Pagotipo, erro error)
	GetPagosTipo(filtro filtros.PagoTipoFiltro) (response []entities.Pagotipo, totalFilas int64, erro error)
	GetPagosTipoReferences(filtro filtros.PagoTipoReferencesFilter) ([]entities.Pagotipo, error)
	DeletePagoTipo(ctx context.Context, id uint64) (erro error)

	//ABM CHANNELS
	CreateChannel(ctx context.Context, request entities.Channel) (id uint64, erro error)
	UpdateChannel(ctx context.Context, request entities.Channel) (erro error)
	GetChannel(filtro filtros.ChannelFiltro) (channel entities.Channel, erro error)
	GetChannels(filtro filtros.ChannelFiltro) (response []entities.Channel, totalFilas int64, erro error)
	DeleteChannel(ctx context.Context, id uint64) (erro error)

	//ABM CUENTAS COMISIONES
	CreateCuentaComision(ctx context.Context, request entities.Cuentacomision) (id uint64, erro error)
	UpdateCuentaComision(ctx context.Context, request entities.Cuentacomision) (erro error)
	GetCuentaComision(filtro filtros.CuentaComisionFiltro) (response entities.Cuentacomision, erro error)
	GetCuentasComisiones(filtro filtros.CuentaComisionFiltro) (response []entities.Cuentacomision, totalFilas int64, erro error)
	DeleteCuentaComision(ctx context.Context, id uint64) (erro error)

	// ABM IMPUESTOS
	GetImpuestosRepository(filtro filtros.ImpuestoFiltro) (response []entities.Impuesto, totalFilas int64, erro error)
	CreateImpuestoRepository(ctx context.Context, impuesto entities.Impuesto) (id uint64, erro error)
	UpdateImpuestoRepository(ctx context.Context, impuesto entities.Impuesto) (erro error)

	// ABM CHANNELS ARANCELES
	GetChannelArancel(filtro filtros.ChannelAranceFiltro) (response entities.Channelarancele, erro error)
	GetChannelsAranceles(filtro filtros.ChannelArancelFiltro) (response []entities.Channelarancele, totalFilas int64, erro error)
	CreateChannelsArancel(ctx context.Context, request entities.Channelarancele) (id uint64, erro error)
	UpdateChannelsArancel(ctx context.Context, request entities.Channelarancele) (erro error)
	DeleteChannelsArancel(ctx context.Context, id uint64) (erro error)
	/*
		Devuelve el saldo actual de una cuenta específica.
		Se debe informar el id de la cuenta.
	*/
	GetSaldoCuenta(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error)

	/*
		Devuelve el saldo actual de un cliente específico.
		Se debe informar una lista de cuentas del cliente.
	*/
	GetSaldoCliente(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error)

	//MOVIMIENTOS
	GetMovimientos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, totalFilas int64, erro error)
	BajaMovimiento(ctx context.Context, movimientos []*entities.Movimiento, motivoBaja string) error
	SaveCuentacomision(comision *entities.Cuentacomision) error
	GetMovimientosNegativos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, erro error)
	UpdateMovimientoMontoRepository(ctxPrueba context.Context, movimiento entities.Movimiento) (erro error)

	// GetMovimientosById(id uint64) (movimiento entities.Movimiento, erro error)

	//TRANSFERENCIAS
	CreateMovimientosTransferencia(ctx context.Context, movimiento []*entities.Movimiento) error
	CreateTransferencias(ctx context.Context, transferencias []*entities.Transferencia) (erro error)
	GetTransferencias(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferencia, totalFilas int64, erro error)
	CreateTransferenciasComisiones(ctx context.Context, transferencias []*entities.Transferenciacomisiones) (erro error)
	GetTransferenciasComisiones(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferenciacomisiones, totalFilas int64, erro error)
	GetMovimientosTransferencias(request reportedtos.RequestPagosPeriodo) (movimientos []entities.Movimiento, erro error)

	/* actualizar el estado de una transferencia con los datos conciliados del banco*/
	UpdateTransferencias(listas bancodtos.ResponseConciliacion) error

	/*
		Modifica el estado de una lista de pagos y además crea un pago estado log
	*/
	UpdateEstadoPagos(pagos []entities.Pago, pagoEstadoId uint64) (erro error)

	//ABM PLAN DE CUOTAS
	/*
		REVIEW: este codigo se debe revisar teniendo en cuenta los cambios que se hicieron en la BD para actualizar installmentdetails
		"CreatePlanCuotasByInstallmenIdRepository"
	*/
	/* Obtener el plan de cuotas */
	GetPlanCuotasByMedioPago(idMedioPago uint) (planCuotas []administraciondtos.PlanCuotasResponseDetalle, erro error)
	// obtiene todos los planes de cuotas
	GetInstallments(fechaDesde time.Time) (medioPagoInstallments []entities.Mediopagoinstallment, erro error)
	// obtengo todos los planes por id
	GetAllInstallmentsById(id uint) (installment []entities.Installment, erro error)
	// obtengo un plan de cuotas por id
	GetInstallmentById(id uint) (planCuotas entities.Installment, erro error)
	CreatePlanCuotasByInstallmenIdRepository(installmentActual, installmentNew entities.Installment, listaPlanCuotas []entities.Installmentdetail) (erro error)

	//CIERRE LOTE
	CreateCierreLoteApiLink(cierreLotes []*entities.Apilinkcierrelote) (erro error)
	CreateMovimientosCierreLote(ctx context.Context, mcl administraciondtos.MovimientoCierreLoteResponse) (erro error)
	GetPrismaCierreLotes(reversion bool) (prismaCierreLotes []entities.Prismacierrelote, erro error)
	CreateMovimientosTemporalesCierreLote(ctx context.Context, mcl administraciondtos.MovimientoTemporalesResponse) (erro error)

	// & Actualizar estados pagos y clrapipago
	ActualizarPagosClRapipagoRepository(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error)
	ActualizarPagosClRapipagoDetallesRepository(barcode []string) (erro error)

	// & Actualizar estados pagos y clmultipagos
	ActualizarPagosClMultipagosDetallesRepository(barcode []string) (erro error)

	// & Crear y actualizar pagos
	CreateCLApilinkPagosRepository(ctx context.Context, pg administraciondtos.RegistroClPagosApilink) (erro error)
	// consultar debines eliminados para ser procesados en el cierre de lote
	GetConsultarDebines(request linkdebin.RequestDebines) (cierreLotes []*entities.Apilinkcierrelote, erro error)
	//&end apilinkcierrelote

	//RI BCRA
	BuildRICuentasCliente(request ribcradtos.RICuentasClienteRequest) (ri []ribcradtos.RiCuentaCliente, erro error)
	BuildRIDatosFondo(request ribcradtos.RiDatosFondosRequest) (ri []ribcradtos.RiDatosFondos, erro error)
	BuilRIInfestaditica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error)

	//CONFIGURACIONES
	GetConfiguraciones(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, totalFilas int64, erro error)
	UpdateConfiguracion(ctx context.Context, request entities.Configuracione) (erro error)

	// UPDATE PAGOS NOTIFICACIDOS
	UpdatePagosNotificados(listaPagosNotificar []uint) (erro error)
	//NOTE -Solo se mantendra hasta que se cree el proceso automatico con rabbit
	UpdatePagosEstadoInicialNotificado(listaPagosNotificar []uint) (erro error)

	// consultar movimientos de la tabla rapipagos para luego ser procesados en el cierre de lote
	GetConsultarMovimientosRapipago(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Rapipagocierrelote, erro error)
	GetConsultarMovimientosRapipagoDetalles(filtro rapipago.RequestConsultarMovimientosRapipagoDetalles) (response []*entities.Rapipagocierrelotedetalles, erro error)
	UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error)

	// consultar movimientos de la tabla multipagos para luego ser procesados en el cierre de lote
	GetConsultarMovimientosMultipagos(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Multipagoscierrelote, erro error)
	// GetConsultarMovimientosRapipagoDetalles(filtro rapipago.RequestConsultarMovimientosRapipagoDetalles) (response []*entities.Rapipagocierrelotedetalles, erro error)
	UpdateCierreLoteMultipagos(cierreLotes []*entities.Multipagoscierrelote) (erro error)
	// & Actualizar estados pagos y clrapipago
	ActualizarPagosClMultipagosRepository(pagosclmultipagos administraciondtos.PagosClMultipagosResponse) (erro error)
	GetConsultarMovimientosMultipagosDetalles(filtro multipagos.RequestConsultarMovimientosMultipagosDetalles) (response []*entities.Multipagoscierrelotedetalles, erro error)

	//PagoTipoChannel
	GetPagosTipoChannelRepository(filtro filtros.PagoTipoChannelFiltro) (response []entities.Pagotipochannel, erro error)
	DeletePagoTipoChannel(id uint64) (erro error)
	CreatePagoTipoChannel(ctx context.Context, pagotipochannel entities.Pagotipochannel) (id uint64, erro error)

	//PETICIONES WEBSERVICES
	GetPeticionesWebServices(filtro filtros.PeticionWebServiceFiltro) (peticiones []entities.Webservicespeticione, totalFilas int64, erro error)

	// MEDIO-PAGO
	GetMedioPagoRepository(filtro filtros.FiltroMedioPago) (mediopago entities.Mediopago, erro error)

	// ARCHIVOS SIBIDOS DE "CIERRE LOTE, PRISMA PX Y PRISMA MX"
	GetCierreLoteSubidosRepository() (entityCl []entities.Prismacierrelote, erro error)
	GetPrismaPxSubidosRepository() (entityPx []entities.Prismapxcuatroregistro, erro error)
	GetPrismaMxSubidosRepository() (entityMx []entities.Prismamxtotalesmovimiento, erro error)

	// Obtener registro de cierre de lote rapipago
	ObtenerArchivoCierreLoteRapipago(nombre string) (existeArchivo bool, erro error)
	ObtenerArchivoCierreLoteMultipagos(nombre string) (existeArchivo bool, erro error)

	ObtenerCierreLoteEnDisputaRepository(estadoDisputa int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error)
	ObtenerCierreLoteContraCargoRepository(estadoReversion int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error)

	ObtenerPagosInDisputaRepository(filtro filtros.ContraCargoEnDisputa) (pagosEnDisputa []entities.Pagointento, erro error)

	// Preferencias
	PostPreferencesRepository(preferenceEntity entities.Preference) (erro error)

	// UPDATE PAGOS MOVIMIENTOS DEV // solo se utiliza para generar movimientos en ambiente sandbox y dev
	UpdatePagosDev(pagos []uint) (erro error)
	// Solicitud de Cuenta
	CreateSolicitudRepository(solicitudEntity entities.Solicitud) (erro error)

	// ? consultar repository CLlotes para herramienta wee
	// consultar movimeintos para herramienta wee
	GetConsultarClRapipagoRepository(filtro filtros.RequestClrapipago) (clrapiapgo []entities.Rapipagocierrelote, totalFilas int64, erro error)

	// consultar pagosintentos calculo de comisiones temporales
	GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error)
	GetPagosApilink(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error)
	GetPagosRapipago(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error)
	GetPagosPrisma(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error)

	// CL apilink
	UpdateCierreloteApilink(request linkdebin.RequestListaUpdateDebines) (erro error)

	GetSuccessPaymentsRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error)
	GetReportesPagoRepository(filtro filtros.PagoFiltro) (reportes entities.Reporte, erro error)

	// RETENCIONES IMPOSITIVAS
	GetRetencionesRepository(request retenciondtos.RentencionRequestDTO) (retenciones []entities.Retencion, count int64, erro error)
	GetClienteRetencionesRepository(request retenciondtos.RentencionRequestDTO) (clienteretencion []entities.ClienteRetencion, count int64, erro error)
	CreateClienteRetencionRepository(request retenciondtos.RentencionRequestDTO) (erro error)
	GetCalcularRetencionesRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)
	GetGravamensRepository(filtro retenciondtos.GravamenRequestDTO) (gravamenes []entities.Gravamen, erro error)
	GetClienteUnlinkedRetencionesRepository(request retenciondtos.RentencionRequestDTO) (retenciones []entities.Retencion, count int64, erro error)
	GetCondicionesRepository(request retenciondtos.RentencionRequestDTO) (condicions []entities.Condicion, erro error)
	CreateRetencionRepository(entity entities.Retencion) (entities.Retencion, error)
	DeleteClienteRetencionRepository(request retenciondtos.RentencionRequestDTO) (erro error)
	UpdateRetencionRepository(entity entities.Retencion) (entities.Retencion, error)
	CreateCondicionRepository(condicion entities.Condicion) (erro error)
	UpdateCondicionRepository(condicion entities.Condicion) (erro error)
	GetClienteRetencionRepository(retencion_id, cliente_id uint) (cliente_retencion entities.ClienteRetencion, erro error)
	EvaluarRetencionesByClienteRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)
	EvaluarRetencionesByMovimientosRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)
	// genera un registro de las retenciones o un credito a favor del cliente
	GenerarCertificacionRepository(comprobante []entities.Comprobante) (erro error)
	GetMovimientosRetencionesRepository(request retenciondtos.RentencionRequestDTO) (listaMovimientosId []uint, erro error)
	// sumar los montos brutos de cada pago sobre el cual se hace la retencion
	GetTotalAmountByMovimientoIdsRepository(listaMovimientosId []uint) (totalAmount uint64, erro error)
	ComprobantesRetencionesDevolverRepository(request retenciondtos.RentencionRequestDTO) (comprobantes []entities.Comprobante, erro error)
	TotalizarRetencionesMovimientosRepository(listaMovimientoIds []uint) (totalRetenciones uint64, erro error)
	GetComprobantesRepository(request retenciondtos.RentencionRequestDTO) (comprobantes []entities.Comprobante, erro error)

	//Certificados Retenciones && Certificados
	PostRetencionCertificadoRepository(certificado entities.Certificado) (erro error)
	GetRetencionClienteRepository(request retenciondtos.RentencionRequestDTO) (cliente entities.ClienteRetencion, erro error)
	GetCertificadoRepository(requestId uint) (cliente entities.Certificado, erro error)
	GetCertificadosVencimientoRepository(request retenciondtos.CertificadoVencimientoDTO) (certificados []entities.Certificado, erro error)
	GetCalcularRetencionesByTransferenciasRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)
	GetMovimientosIdsCalculoRetencionComprobante(request retenciondtos.RentencionRequestDTO) (resultado []uint, erro error)

	// NOTE pruebas crear auditoria
	CreateAuditoria(resultado entities.Auditoria) (erro error)

	CalcularRetencionesByTransferenciasSinAgruparRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)
}

type repository struct {
	SQLClient        *database.MySQLClient
	auditoriaService auditoria.AuditoriaService
	utilService      util.UtilService
}

func NewRepository(sqlClient *database.MySQLClient, a auditoria.AuditoriaService, t util.UtilService) Repository {
	return &repository{
		SQLClient:        sqlClient,
		auditoriaService: a,
		utilService:      t,
	}
}

func (r *repository) BeginTx() {
	r.SQLClient.TX = r.SQLClient.DB
	r.SQLClient.DB = r.SQLClient.Begin()
}
func (r *repository) CommitTx() {
	r.SQLClient.Commit()
	r.SQLClient.DB = r.SQLClient.TX
}
func (r *repository) RollbackTx() {
	r.SQLClient.Rollback()
	r.SQLClient.DB = r.SQLClient.TX
}

// func (r *repository) GetMovimientosById(id uint64) (movimiento entities.Movimiento, erro error) {
// 	resp := r.SQLClient.Model(entities.Movimiento{})

// 	if id > 0 {
// 		resp.Where("id = ?", id)
// 	}
// 	resp.Preload("Cuenta")

// 	resp.First(&movimiento)

// 	if resp.Error != nil {

// 		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
// 			/*REVIEW en el caso de que no devuelva un elemento */
// 			erro = nil //fmt.Errorf(RESULTADO_NO_ENCONTRADO)
// 			return
// 		}

// 		erro = fmt.Errorf("error al obtener movimieto por id")

//		}
//		return
//	}
func (r *repository) GetCuenta(filtro filtros.CuentaFiltro) (cuenta entities.Cuenta, erro error) {

	resp := r.SQLClient.Model(entities.Cuenta{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.DistintoId > 0 {
		resp.Where("id != ?", filtro.DistintoId)
	}

	if len(filtro.Cbu) > 0 {
		resp.Where("cbu = ?", filtro.Cbu)
	}

	if len(filtro.Cvu) > 0 {
		resp.Where("cvu = ?", filtro.Cvu)
	}

	resp.Preload("Cliente")

	resp.First(&cuenta)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			/*REVIEW en el caso de que no devuelva un elemento */
			erro = nil //fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCuenta: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) PagoById(pagoID int64) (*entities.Pago, error) {
	var pago entities.Pago

	res := r.SQLClient.Model(entities.Pago{}).Preload("PagosResultados.Pagoresultadodetalle").Find(&pago, pagoID)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontraron pagos con id %d", pagoID)
	}

	return &pago, nil
}

func (r *repository) CuentaByClientePage(cliente int64, limit, offset int) (*[]entities.Cuenta, int64, error) {
	var cuentas []entities.Cuenta
	var count int64

	res := r.SQLClient.Model(entities.Cuenta{}).Preload("Pagotipos")
	res.Where("clientes_id = ?", cliente)
	res.Count(&count)
	res.Limit(limit)
	res.Offset(offset)
	res.Find(&cuentas)

	if res.RowsAffected <= 0 {
		return nil, 0, fmt.Errorf("no se encontraron cuentas para el cliente con id %d", cliente)
	}

	return &cuentas, count, nil
}

func (r *repository) CuentaByID(cuentaID int64) (*entities.Cuenta, error) {
	var cuenta entities.Cuenta

	res := r.SQLClient.Model(entities.Cuenta{}).Preload("Pagotipos").Preload("Cuentacomisiones").Find(&cuenta, cuentaID)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cuenta con id %d", cuentaID)
	}

	return &cuenta, nil
}

func (r *repository) SaveCuenta(ctx context.Context, cuenta *entities.Cuenta) (bool, error) {
	res := r.SQLClient.WithContext(ctx)
	if cuenta.ID == 0 {
		res = res.Create(&cuenta)
	} else {
		res = res.Model(&cuenta).Updates(cuenta)
	}

	if res.RowsAffected <= 0 {
		return false, fmt.Errorf("error al guardar cuenta: %s", res.Error.Error())
	}

	err := r.auditarAdministracion(res.Statement.Context, true)
	if err != nil {
		return false, fmt.Errorf("no es posible la auditoría: %v", err)
	}

	return true, nil
}

func (r *repository) SetApiKey(ctx context.Context, cuenta entities.Cuenta) (erro error) {

	entidad := entities.Cuenta{
		Model: gorm.Model{ID: cuenta.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Select("apikey").Updates(cuenta)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_APIKEY)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "SetApiKey",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, erro error) {
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

func (r *repository) UpdateCuenta(ctx context.Context, cuenta entities.Cuenta) (erro error) {

	entidad := entities.Cuenta{
		Model: gorm.Model{ID: cuenta.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,apikey,created_at,deleted_at").Select("*").Updates(cuenta)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) DeleteCuenta(id uint64) (erro error) {

	entidad := entities.Cuenta{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) SavePagotipo(tipo *entities.Pagotipo) (bool, error) {
	if tipo.ID == 0 {
		res := r.SQLClient.Create(tipo)
		if res.RowsAffected <= 0 {
			return false, fmt.Errorf("error al crear tipo de pago: %s", res.Error.Error())
		}
		return true, nil
	}
	res := r.SQLClient.Model(entities.Pagotipo{}).Where("id = ?", tipo.ID).Updates(tipo)
	if res.RowsAffected <= 0 {
		return false, fmt.Errorf("error al actualizar tipo de pago: %s", res.Error.Error())
	}

	return true, nil
}

func (r *repository) ConsultarEstadoPagosRepository(parametrosVslido administraciondtos.ParamsValidados, filtro filtros.PagoFiltro) (entityPagos []entities.Pago, erro error) {
	/*
		los fitros que recibe son:
		- 1 uuid
		- arrays de uuid
		- rango de fecha
		- external reference
	*/
	resp := r.SQLClient.Model(entities.Pago{})

	if filtro.PagoEstadosId != 0 {
		resp.Where("pagoestados_id <> ?", filtro.PagoEstadosId)
	}
	if len(filtro.PagosTipoIds) > 0 {
		resp.Where("pagostipo_id in (?)", filtro.PagosTipoIds)
	}

	if parametrosVslido.Uuuid || parametrosVslido.Uuids {
		resp.Where("uuid in ? ", filtro.Uuids)
	}

	if parametrosVslido.ExternalReference {
		resp.Where("external_reference = ?", filtro.ExternalReference)
	}

	if parametrosVslido.RangoFecha {
		resp.Where("created_at BETWEEN ? AND ?", filtro.Fecha[0], filtro.Fecha[1])
	}
	if filtro.CargarPagoTipos {
		resp.Preload("PagosTipo")
	}
	if filtro.CargarPagoEstado {
		resp.Preload("PagoEstados")
	}
	if filtro.CargaPagoIntentos {
		stateComments := []string{"approved", "INICIADO"}
		resp.Preload("PagoIntentos", "state_comment in ? ", stateComments)
		resp.Preload("PagoIntentos.Mediopagos")
		resp.Preload("PagoIntentos.Mediopagos.Channel")
		resp.Preload("PagoIntentos.Movimientos", "tipo = ?", "C")
		resp.Preload("PagoIntentos.Movimientos.Movimientocomisions")
		resp.Preload("PagoIntentos.Movimientos.Movimientoimpuestos")

	}
	resp.Find(&entityPagos)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_PAGO)
		return
	}
	return
}

func (r *repository) SaveCuentacomision(comision *entities.Cuentacomision) error {
	if comision.ID == 0 {
		res := r.SQLClient.Create(comision)
		if res.RowsAffected <= 0 {
			return fmt.Errorf("error al crear comision: %s", res.Error.Error())
		}
		return nil
	}
	res := r.SQLClient.Model(entities.Cuentacomision{}).Where("id = ?", comision.ID).Updates(comision)
	if res.RowsAffected <= 0 {
		return fmt.Errorf("error al actualizar comision: %s", res.Error.Error())
	}
	return nil
}

func (r *repository) GetPagosByUUID(uuids []string) (pagos []*entities.Pago, erro error) {

	if len(uuids) > 0 {
		resp := r.SQLClient.Preload("PagoIntentos", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC")
		}).Preload("PagoTipos").Where("uuid IN ?", uuids).Find(&pagos)
		if resp.Error != nil {
			erro = resp.Error
		}
	}

	return
}

func (r *repository) GetPagosEstados(filtro filtros.PagoEstadoFiltro) (estados []entities.Pagoestado, erro error) {

	resp := r.SQLClient.Model(entities.Pagoestado{})

	if filtro.BuscarPorFinal {
		resp.Where("final = ?", filtro.Final)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("estado", filtro.Nombre)
	}
	if filtro.EstadoId != 0 {
		resp.Where("id = ?", filtro.EstadoId)
	}
	resp.Find(&estados)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO_ESTADO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosEstados",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return

}

func (r *repository) GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estado entities.Pagoestado, erro error) {

	resp := r.SQLClient.Model(entities.Pagoestado{})

	if filtro.BuscarPorFinal {
		if filtro.Final {

			resp.Where("final = ?", filtro.Final)
		}
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("estado", filtro.Nombre)
	}

	if filtro.EstadoId > 0 {
		resp.Where("id = ?", filtro.EstadoId)
	}

	resp.First(&estado)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_PAGO_ESTADO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagoEstado",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return

}

func (r *repository) GetSaldoCuenta(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error) {

	resp := r.SQLClient.Table("movimientos as m").Select("m.cuentas_id, sum(m.monto) as total").
		Where("m.deleted_at IS NULL").Group("m.cuentas_id").Having("m.cuentas_id", cuentaId).Scan(&saldo)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_SALDO_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetSaldoCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetSaldoCliente(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error) {

	resp := r.SQLClient.Table("movimientos as m").
		Select("m.cuentas_id, sum(m.monto) as total").
		Joins("inner join cuentas as c on c.id = m.cuentas_id").
		Joins("left join clientes as cl on cl.id = c.clientes_id").
		Where("m.deleted_at IS NULL").Where("cl.id", clienteId).Group("m.cuentas_id").Scan(&saldo)

	if resp.Error != nil {
		erro = resp.Error
	}

	return
}

func (r *repository) GetCuentasByCliente(clienteId uint64) (cuentas []entities.Cuenta, erro error) {

	resp := r.SQLClient.Where("clientes_id", clienteId).Find(&cuentas)

	if resp.Error != nil {
		erro = resp.Error
	}

	return
}

func (r *repository) UpdateEstadoPagos(pagos []entities.Pago, pagoEstadoId uint64) (erro error) {

	var estadosLogs []entities.Pagoestadologs

	for i := range pagos {
		pagoEstado := entities.Pagoestadologs{
			PagosID:       int64(pagos[i].ID),
			PagoestadosID: pagos[i].PagoestadosID,
		}
		estadosLogs = append(estadosLogs, pagoEstado)
	}

	erro = r.SQLClient.Transaction(func(tx *gorm.DB) error {

		// Creo los logs de estados
		if err := tx.Create(&estadosLogs).Error; err != nil {

			erro := fmt.Errorf(ERROR_CREAR_ESTADO_LOGS)

			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       err.Error(),
				Funcionalidad: "UpdateEstadoPagos",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				logs.Error(err.Error())
			}

			return erro
		}
		// Modifico los estados de los pagos
		if err := tx.Model(&pagos).Omit(clause.Associations).UpdateColumns(entities.Pago{PagoestadosID: int64(pagoEstadoId), Model: gorm.Model{UpdatedAt: time.Now()}}).Error; err != nil {

			erro := fmt.Errorf(ERROR_UPDATE_PAGO)

			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       err.Error(),
				Funcionalidad: "UpdateEstadoPagos",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				logs.Error(err.Error())
			}

			return erro
		}

		return nil
	})

	return
}

func (r *repository) GetPago(filtro filtros.PagoFiltro) (pago entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	if !filtro.FechaPagoFin.IsZero() {
		resp.Where("cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaPagoInicio, filtro.FechaPagoFin)
	}

	if len(filtro.Ids) > 0 {
		resp.Where("id IN ?", filtro.Ids)
	}

	if filtro.PagoEstadosId > 0 {
		resp.Where("pagoestados_id", filtro.PagoEstadosId)
	}

	if len(filtro.TiempoExpiracion) > 0 {
		resp.Where("timestampdiff(day, created_at, now() ) >= ?", filtro.TiempoExpiracion)
	}

	if filtro.CargaPagoIntentos {
		resp.Preload("PagoIntentos", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC")
		})
	}
	if len(filtro.Uuids) > 0 {
		resp.Where("uuid IN ?", filtro.Uuids)
	}

	if len(filtro.ExternalReference) > 0 {
		resp.Where("external_reference = ?", filtro.ExternalReference)
	}

	if filtro.CargaMedioPagos {
		if filtro.CargarChannel {
			resp.Preload("PagoIntentos.Mediopagos.Channel")
		} else {
			resp.Preload("PagoIntentos.Mediopagos")
		}
	}

	if filtro.CargarPagoTipos {
		resp.Preload("PagosTipo")
	}

	if filtro.CargarCuenta {
		if filtro.CuentaId > 0 {
			resp.Preload("PagosTipo.Cuenta", "id = ?", filtro.CuentaId)
		} else {
			resp.Preload("PagosTipo.Cuenta")
		}
	}

	if filtro.CargarPagoEstado {
		resp.Preload("PagoEstados")
	}

	if len(filtro.Fecha) > 0 {
		resp.Where("updated_at BETWEEN ? AND ?", filtro.Fecha[0], filtro.Fecha[1])
	}

	resp.First(&pago)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPago",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetPago: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return

}

func (r *repository) GetPagos(filtro filtros.PagoFiltro) (pagos []entities.Pago, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	if !filtro.FechaPagoFin.IsZero() {
		if filtro.FiltroFechaPaid {
			resp.Preload("PagoIntentos").Joins("JOIN pasarela.pagointentos as p1 ON (pagos.id = p1.pagos_id) LEFT OUTER JOIN pasarela.pagointentos p2 ON (pagos.id = p2.pagos_id AND (p1.created_at < p2.created_at OR (p1.created_at = p2.created_at AND p1.id < p2.id)))").Where("p2.id IS NULL").Where("cast(p1.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaPagoInicio, filtro.FechaPagoFin).Order("p1.paid_at ASC")
			if filtro.Ordenar {
				if filtro.Descendente {
					resp.Order("pagos.PagoIntentos.paid_at DESC")
				}
				if !filtro.Descendente {
					resp.Order("pagos.PagoIntentos.paid_at ASC")
				}
			}

		} else {
			resp.Where("pagos.created_at BETWEEN cast(? as datetime) AND cast(? as datetime)", filtro.FechaPagoInicio, filtro.FechaPagoFin)
			if filtro.Ordenar {
				if filtro.Descendente {
					resp.Order("pagos.created_at DESC")
				}
				if !filtro.Descendente {
					resp.Order("pagos.created_at ASC")
				}
			}
		}
	}

	if filtro.BuscarNotificado {
		if !filtro.Notificado {
			resp.Where("notificado = ?", filtro.Notificado)
		}
	}

	if len(filtro.Fecha) > 0 {
		resp.Where("updated_at BETWEEN ? AND ?", filtro.Fecha[0], filtro.Fecha[1])
	}

	if len(filtro.Ids) > 0 {

		resp.Where("id IN ?", filtro.Ids)
	}

	if len(filtro.Uuids) > 0 {
		resp.Where("uuid IN ?", filtro.Uuids)
	}

	if len(filtro.PagoEstadosIds) > 0 {
		resp.Where("pagoestados_id IN ?", filtro.PagoEstadosIds)
	}

	if filtro.PagoEstadosId > 0 {
		resp.Where("pagoestados_id", filtro.PagoEstadosId)
	}

	if filtro.PagosTipoId > 0 {
		resp.Where("pagostipo_id = ?", filtro.PagosTipoId)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("payer_name LIKE ?", "%"+filtro.Nombre+"%")
	}

	if len(filtro.ExternalReference) > 0 {
		resp.Where("external_reference LIKE ?", "%"+filtro.ExternalReference+"%")

	}

	if len(filtro.PagosTipoIds) > 0 {
		resp.Where("pagostipo_id IN ?", filtro.PagosTipoIds)
	}

	if !filtro.VisualizarPendientes {
		filtro := filtros.PagoEstadoFiltro{
			Nombre: "pending",
		}

		estadoPendiente, err := r.GetPagoEstado(filtro)

		if err != nil {
			erro = err
			return
		}

		resp.Where("pagoestados_id != ?", estadoPendiente.ID)
	}

	if len(filtro.TiempoExpiracion) > 0 {
		resp.Where("timestampdiff(day, created_at, now() ) >= ?", filtro.TiempoExpiracion)
	}
	if len(filtro.TiempoExpiracionSecondDueDate) > 0 {
		resp.Where("timestampdiff(day, second_due_date, now() ) >= ?", filtro.TiempoExpiracionSecondDueDate)
	}
	if filtro.CargarCuenta {
		if filtro.CuentaId > 0 {
			resp.Preload("PagosTipo.Cuenta", "cuentas.id = ?", filtro.CuentaId).Joins("INNER JOIN pagotipos as pt on pt.id = pagos.pagostipo_id INNER JOIN cuentas as c on c.id = pt.cuentas_id").Where("c.id = ?", filtro.CuentaId)
		} else {
			resp.Preload("PagosTipo.Cuenta")
		}
	}

	if filtro.CargaPagoIntentos {
		if filtro.CargaMedioPagos {
			if filtro.MedioPagoId > 0 {
				if filtro.CargarChannel {
					//NOTE: SE MODIFICO LA CONSULTA m.id por m.channels_id
					resp.Preload("PagoIntentos.Mediopagos.Channel").Joins("LEFT JOIN pagointentos as pi on pagos.id = pi.pagos_id  INNER JOIN mediopagos as m on m.id = pi.mediopagos_id INNER JOIN channels as ch on ch.id = m.channels_id").Where("m.channels_id = ?", filtro.MedioPagoId).Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO")
				} else {
					resp.Preload("PagoIntentos.Mediopagos").Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id inner join INNER JOIN mediopagos as m on m.id = pi.mediopagos_id").Where("m.channels_id = ?", filtro.MedioPagoId).Where("pi.state_comment = ? OR pi.state_comment = ?", "approved", "INICIADO")
				}
			} else {
				if filtro.CargarChannel {
					resp.Preload("PagoIntentos", "state_comment = ? OR state_comment = ?", "approved", "INICIADO")
					resp.Preload("PagoIntentos.Mediopagos.Channel")
				} else {
					resp.Preload("PagoIntentos.Mediopagos").Joins("INNER JOIN pagointentos as pi on pagos.id = pi.pagos_id").Where("pi.external_id != ''").Order("pi.paid_at ASC")
				}
			}
		} else {
			resp.Preload("PagoIntentos", func(db *gorm.DB) *gorm.DB {
				return db.Order("id DESC")
			})
		}

	}

	if filtro.CargarPagoTipos {
		resp.Preload("PagosTipo")
	}

	if filtro.CargarPagoEstado {
		resp.Preload("PagoEstados")
	}

	if filtro.CargarPagosItems {
		resp.Preload("Pagoitems")
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	if filtro.OrdenarPorPaidAtPagointento {
		resp.Preload("PagoIntentos", func(db *gorm.DB) *gorm.DB {
			return db.Order("paid_at DESC")
		})
	}

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetPaymentByExternal(filtroPago filtros.PagoFiltro) (*entities.Pago, error) {

	var pago entities.Pago

	res := r.SQLClient.Model(entities.Pago{}).Preload("Pagoitems").Preload("PagoIntentos")

	if filtroPago.CargaMedioPagos {
		res = res.Preload("PagoIntentos.Mediopagos")
	}

	if filtroPago.CargarPagoEstado {
		res = res.Preload("PagoEstados")
	}

	res = res.Where("external_reference = ?", filtroPago.ExternalReference).Find(&pago)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no existe pago con external %s", filtroPago.ExternalReference)
	}

	return &pago, nil
}

func (r *repository) GetPlanCuotasByMedioPago(idMedioPago uint) (planCuotas []administraciondtos.PlanCuotasResponseDetalle, erro error) {
	var details []entities.Installmentdetail
	response := r.SQLClient.Model(entities.Installmentdetail{}).Joins("Installment").Joins("LEFT JOIN mediopagos ON (Installment.id = mediopagos.installments_id) AND mediopagos.id = ?", idMedioPago).Find(&details)
	if response.Error != nil {
		erro = response.Error
	}
	for _, v := range details {
		planCuotas = append(planCuotas, administraciondtos.PlanCuotasResponseDetalle{
			InstallmentsID: v.InstallmentsID,
			Cuota:          uint(v.Cuota),
			Tna:            v.Tna,
			Tem:            v.Tem,
			Coeficiente:    v.Coeficiente,
		})
	}
	return
}

func (r *repository) GetInstallments(fechaDesde time.Time) (medioPagoInstallments []entities.Mediopagoinstallment, erro error) {
	res := r.SQLClient.Table("mediopagoinstallments as mpi")
	res.Preload("Installments")
	res.Preload("Installments.Installmentdetail")
	res.Find(&medioPagoInstallments)
	if res.Error != nil {
		logs.Info(res.Error)
		erro = errors.New(ERROR_CREAR_INSTALLMENT_DETAILS)
		return
	}
	return
}

func (r *repository) GetAllInstallmentsById(id uint) (installment []entities.Installment, erro error) {
	res := r.SQLClient.Model(entities.Installment{}).Where("mediopagoinstallments_id = ?", id).Find(&installment)
	if res.Error != nil {
		erro = errors.New(ERROR_CONSULTA_INSTALLMENT)
		return
	}
	return
}

func (r *repository) GetInstallmentById(id uint) (installment entities.Installment, erro error) {
	res := r.SQLClient.Model(entities.Installment{}).Where("mediopagoinstallments_id = ?", id).Order("created_at desc").First(&installment)
	if res.Error != nil {
		erro = errors.New(ERROR_CONSULTA_INSTALLMENT)
		return
	}
	return
}

func (r *repository) CreatePlanCuotasByInstallmenIdRepository(installmentActual, installmentNew entities.Installment, listaPlanCuotas []entities.Installmentdetail) (erro error) {
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(entities.Installment{}).Omit(clause.Associations).Where("id = ?", installmentActual.ID).Update("vigencia_hasta", installmentActual.VigenciaHasta)
		if res.Error != nil {
			logs.Info(res.Error)
			return errors.New(ERROR_ACTUALIZAR_INSTALLMENT)
		}
		installmentNew.Installmentdetail = listaPlanCuotas

		res = tx.Create(&installmentNew)
		if res.Error != nil {
			logs.Info(res.Error)
			return errors.New(ERROR_CREAR_INSTALLMENT_DETAILS)
		}
		return nil
	})
}

func (r *repository) GetPagosIntentos(filtro filtros.PagoIntentoFiltro) (pagosIntentos []entities.Pagointento, erro error) {

	resp := r.SQLClient.Model(entities.Pagointento{})

	if len(filtro.ExternalIds) > 0 {

		resp.Where("external_id IN ?", filtro.ExternalIds)
	}

	if filtro.PagoIntentoAprobado {
		resp.Where("paid_at <> ?", "0000-00-00 00:00:00")
	}

	if filtro.ExternalId {
		resp.Where("external_id <>  ? OR external_id <>  ?", "", "0")
		// resp.Where("external_id <> 0")
	}

	if len(filtro.TransaccionesId) > 0 {
		resp.Where("transaction_id IN (?)", filtro.TransaccionesId)
	}
	if len(filtro.TicketNumber) > 0 {
		resp.Where("ticket_number IN (?)", filtro.TicketNumber)
	}

	if len(filtro.CodigoAutorizacion) > 0 {
		resp.Where("authorization_code IN (?)", filtro.CodigoAutorizacion)
	}

	if len(filtro.Barcode) > 0 {
		resp.Where("barcode IN (?)", filtro.Barcode)
	}

	if len(filtro.PagosId) > 0 {
		resp.Where("pagos_id IN (?)", filtro.PagosId)
	}

	if filtro.ChannelIdFiltro != 0 {
		resp.Joins("INNER JOIN mediopagos as mp ON mp.id = pagointentos.mediopagos_id and mp.channels_id = ?", filtro.ChannelIdFiltro)
	}

	if filtro.Channel {

		resp.Preload("Mediopagos")
		resp.Preload("Mediopagos.Channel")
	}

	if filtro.PagoEstadoIdFiltro != 0 {
		resp.Joins("INNER JOIN pagos as p ON p.id = pagointentos.pagos_id and p.pagoestados_id = ?", filtro.PagoEstadoIdFiltro)
	}

	if filtro.CargarPago {

		resp.Preload("Pago")

	}

	if filtro.CargarPagoTipo {
		resp.Preload("Pago.PagosTipo")
		if filtro.CargarCuenta {
			resp.Preload("Pago.PagosTipo.Cuenta.Subcuentas")
			if filtro.CargarCliente {
				resp.Preload("Pago.PagosTipo.Cuenta.Cliente")
				if filtro.CargarImpuestos {
					resp.Preload("Pago.PagosTipo.Cuenta.Cliente.Iva")
					resp.Preload("Pago.PagosTipo.Cuenta.Cliente.Iibb")
				}
			}
			if filtro.CargarCuentaComision {
				resp.Preload("Pago.PagosTipo.Cuenta.Cuentacomisions")
			}
		}
	}

	if filtro.CargarPagoEstado {
		resp.Preload("Pago.PagoEstados")
	}

	if filtro.CargarMovimientos {
		resp.Preload("Movimientos")
	}

	if filtro.CargarInstallmentdetail {
		resp.Preload("Installmentdetail")
	}

	resp.Find(&pagosIntentos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO_INTENTO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosIntentos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetConfiguraciones(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Configuracione{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("nombre like ?", fmt.Sprintf("%%%s%%", filtro.Nombre))
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))

	}

	resp.Find(&configuraciones)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONFIGURACIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConfiguraciones",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) UpdateConfiguracion(ctx context.Context, request entities.Configuracione) (erro error) {

	entidad := entities.Configuracione{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at,nombre").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_CONFIGURACIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateConfiguracion",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, request)

	return
}

// ABM CLIENTES
func (r *repository) GetClientes(filtro filtros.ClienteFiltro) (clientes []entities.Cliente, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Cliente{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.ClientesIds) > 0 {
		resp.Where("id in ?", filtro.ClientesIds)
	}

	if filtro.CargarImpuestos {
		resp.Preload("Iva")
		resp.Preload("Iibb")
	}

	if filtro.CargarCuentas {
		resp.Preload("Cuentas")
	}

	if filtro.CargarRubros {
		resp.Preload("Cuentas.Rubro")
	}

	if filtro.RetiroAutomatico {
		resp.Where("retiro_automatico = ?", filtro.RetiroAutomatico)
	}

	if filtro.SplitCuentas {
		resp.Where("split_cuentas = ?", filtro.SplitCuentas)
	}

	if filtro.CargarContactos {
		resp.Preload("Contactosreportes")
	}

	if filtro.CargarSubcuentas {
		resp.Preload("Cuentas.Subcuentas")
	}

	if filtro.SujetoRetencion {
		resp.Where("clientes.sujeto_retencion", filtro.SujetoRetencion)
		resp.Preload("Retenciones.Condicion.Gravamen").Preload("Retenciones.Channel")
	}

	if filtro.Formulario8125 {
		resp.Where("clientes.formulario_8125", filtro.Formulario8125)
	}

	// if filtro.CargarCuentaComision {
	// 	resp.Preload("Cuentas.Cuentacomisions")
	// }

	// if filtro.CargarTiposPago {
	// 	resp.Preload("Cuentas.Pagotipos")
	// }

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))

	}

	resp.Find(&clientes)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CARGAR_CLIENTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetClientes",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetClientes: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetCliente(filtro filtros.ClienteFiltro) (cliente entities.Cliente, erro error) {

	resp := r.SQLClient.Model(entities.Cliente{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.Nombre != "" {
		resp.Where("cliente = ?", filtro.Nombre)
	}

	if filtro.DistintoId > 0 {
		resp.Where("id = ?", filtro.DistintoId)
	}

	if filtro.UserId > 0 {
		resp.Joins("JOIN clienteusers ON clientes.id = clienteusers.clientes_id AND clienteusers.user_id = ?", filtro.UserId)

	}

	if len(filtro.Cuit) > 0 {
		resp.Where("cuit = ?", filtro.Cuit)
	}

	if filtro.RetiroAutomatico {
		resp.Where("retiro_automatico = ?", filtro.RetiroAutomatico)
	}

	if filtro.CargarImpuestos {
		resp.Preload("Iva")
		resp.Preload("Iibb")
	}

	if filtro.CargarCuentas {
		resp.Preload("Cuentas")
	}

	if filtro.CargarRubros {
		resp.Preload("Cuentas.Rubro")
	}
	if filtro.CargarCuentaComision {
		resp.Preload("Cuentas.Cuentacomisions.Channel")
		resp.Preload("Cuentas.Cuentacomisions.ChannelArancel")
	}

	// if filtro.CargarCuentaComision {
	// 	resp.Preload("Cuentas.Cuentacomisions")
	// }

	if filtro.CargarTiposPago {
		resp.Preload("Cuentas.Pagotipos")
	}

	if filtro.CargarContactos {
		resp.Preload("Contactosreportes")
	}

	resp.First(&cliente)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			/*REVIEW en el caso de que no devuelva un elemento */
			erro = nil //fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_CLIENTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCliente",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCliente: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return

}

func (r *repository) CreateCliente(ctx context.Context, cliente entities.Cliente) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&cliente)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CLIENTE)
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

	id = uint64(cliente.ID)

	erro = r.auditarAdministracion(result.Statement.Context, id)
	if erro != nil {
		return id, erro
	}

	return
}

func (r *repository) UpdateCliente(ctx context.Context, cliente entities.Cliente) (erro error) {

	entidad := entities.Cliente{
		Model: gorm.Model{ID: cliente.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(cliente)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_MODIFICAR_CLIENTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateCliente",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, cliente)

	return
}

func (r *repository) DeleteCliente(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Cliente{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CLIENTE)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteCliente",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

// ABM RUBROS
func (r *repository) GetRubros(filtro filtros.RubroFiltro) (rubros []entities.Rubro, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Rubro{})

	if len(filtro.Rubro) > 0 {
		resp.Where("rubro like ?", fmt.Sprintf("%%%s%%", filtro.Rubro))
	}

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&rubros)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CARGAR_RUBROS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRubros",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetRubros: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetRubro(filtro filtros.RubroFiltro) (rubro entities.Rubro, erro error) {

	resp := r.SQLClient.Model(entities.Rubro{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.Rubro) > 0 {
		resp.Where("rubro", filtro.Rubro)
	}

	resp.First(&rubro)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_RUBROS)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRubro",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetRubro: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return

}

func (r *repository) CreateRubro(ctx context.Context, rubro entities.Rubro) (id uint64, erro error) {
	if rubro.ID > 0 {
		rubro.ID = 0
	}
	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&rubro)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_RUBRO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateRubro",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(rubro.ID)

	erro = r.auditarAdministracion(result.Statement.Context, id)

	return
}

func (r *repository) UpdateRubro(ctx context.Context, rubro entities.Rubro) (erro error) {

	entidad := entities.Rubro{
		Model: gorm.Model{ID: rubro.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(rubro)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_RUBRO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateRubro",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, rubro)

	return
}

// ABM PAGO TIPOS
func (r *repository) GetPagosTipo(filtro filtros.PagoTipoFiltro) (response []entities.Pagotipo, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Pagotipo{})

	if len(filtro.PagoTipo) > 0 {
		resp.Where("pagotipo = ?", filtro.PagoTipo)
	}

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.CargarCuenta {
		resp.Preload("Cuenta")
	}

	if filtro.IdCuenta > 0 {
		resp.Where("cuentas_id = ?", filtro.IdCuenta)
	}

	if filtro.CargarTipoPagoChannels {
		resp.Preload("Pagotipochannel.Channel")
		resp.Preload("Pagotipoinstallment")
	}
	//Se obtienen los pagos que no poseen su estado inicial notificado, para notificar al cliente en el proceso background,
	if filtro.CargarPagosEstadoInicialNotificado {
		resp.Preload("Pagos", "pagoestados_id IN ? AND estado_inicial_notificado = ?  AND cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.PagoEstadosIds, 0, filtro.FechaPagoInicio, filtro.FechaPagoFin)
		if filtro.FiltroMediopagosID {
			resp.Preload("Pagos.PagoIntentos", "mediopagos_id IN (1,2,3,4,7,8,9,12,13,14,15,16,21,22,23,24,26,28,29,30,31,32,33,34,35,36)")
		}
		resp.Preload("Pagos.PagoEstados")
		resp.Preload("Pagos.PagoIntentos.Mediopagos.Channel")
	}

	//  filtro cargar los pagos y para filtrar por estado , pagos de los ultimos 3 dias para notificar al usuario
	if filtro.CargarPagos {
		resp.Preload("Pagos", "pagoestados_id IN ? AND notificado = ? AND cast(created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.PagoEstadosIds, filtro.CargarPagosNotificado, filtro.FechaPagoInicio, filtro.FechaPagoFin)
		resp.Preload("Pagos.PagoEstados")
		resp.Preload("Pagos.PagoIntentos.Mediopagos.Channel")
	}

	// if filtro.CargarPagosIntentos {
	// 	if len(filtro.ExternalId) > 0 {
	// 		resp.Joins("INNER JOIN pagos as pg on pagotipos.id = pg.pagostipo_id INNER JOIN pagointentos as pi on pg.id = pi.pagos_id").
	// 			Where("pi.external_id IN ?", filtro.ExternalId)
	// 	}
	// 	resp.Preload("Pagos.PagoIntentos", func(db *gorm.DB) *gorm.DB {
	// 		return db.Where("pagointentos.external_id IN ?", filtro.ExternalId)
	// 	})

	// }

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CARGAR_PAGO_TIPO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosTipo",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetPagosTipo: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetPagosTipoReferences(filtro filtros.PagoTipoReferencesFilter) ([]entities.Pagotipo, error) {
	var pagosTipos []entities.Pagotipo

	resp := r.SQLClient.Model(entities.Pagotipo{})

	if filtro.IdCuenta > 0 {
		resp.Where("cuentas_id = ?", filtro.IdCuenta)
	}

	//  filtro cargar los pagos
	if filtro.CargarPagos {

		// buscar por external references
		if len(filtro.ExternalReferences) > 0 {
			resp.Preload("Pagos", "external_reference IN (?)", filtro.ExternalReferences)
		}

		// buscar por id de pago
		if len(filtro.PagosId) > 0 {
			resp.Preload("Pagos", "id IN (?)", filtro.PagosId)
		}

		resp.Preload("Pagos.PagoEstados")
		resp.Preload("Pagos.PagoIntentos.Mediopagos.Channel")

	}

	resp.Find(&pagosTipos)

	if resp.Error != nil {

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosTipo",
		}

		if err := r.utilService.CreateLogService(log); err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetPagosTipo: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

		return nil, fmt.Errorf(ERROR_CARGAR_PAGO_TIPO)
	}

	return pagosTipos, nil
}

func (r *repository) GetPagoTipo(filtro filtros.PagoTipoFiltro) (response entities.Pagotipo, erro error) {

	resp := r.SQLClient.Model(entities.Pagotipo{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.CargarCuenta {
		resp.Preload("Cuenta")
	}

	if len(filtro.PagoTipo) > 0 {
		resp.Where("pagotipo", filtro.PagoTipo)
	}

	if filtro.CargarTipoPagoChannels {
		resp.Preload("Pagotipochannel.Channel")
		resp.Preload("Pagotipoinstallment")
	}

	resp.First(&response)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CARGAR_PAGO_TIPO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagoTipo",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetPagoTipo: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return

}

func (r *repository) CreatePagoTipo(ctx context.Context, request entities.Pagotipo, channel []int64, cuotas []string) (id uint64, erro error) {

	r.BeginTx()
	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_PAGO_TIPO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreatePagoTipo",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		r.RollbackTx()
		return
	}

	id = uint64(request.ID)

	// agregar pagotipochannels
	for _, ch := range channel {
		entidadChannel := entities.Pagotipochannel{
			PagotiposId: uint(id),
			ChannelsId:  uint(ch),
		}

		resultpagotipochannels := r.SQLClient.WithContext(ctx).Omit("id").Create(&entidadChannel)

		if resultpagotipochannels.Error != nil {
			erro = fmt.Errorf(ERROR_CREAR_RUBRO)
			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       result.Error.Error(),
				Funcionalidad: "Createpagotipochannels",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
				logs.Error(mensaje)
			}
			r.RollbackTx()
			return
		}

	}

	// agregar cuotas
	for _, c := range cuotas {
		entidadCuota := entities.Pagotipointallment{
			PagotiposId: uint(id),
			Cuota:       c,
		}

		resultcuotas := r.SQLClient.WithContext(ctx).Omit("id").Create(&entidadCuota)

		if resultcuotas.Error != nil {
			erro = fmt.Errorf(ERROR_CREAR_RUBRO)
			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       result.Error.Error(),
				Funcionalidad: "Createpagotipoinstallment",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
				logs.Error(mensaje)
			}
			r.RollbackTx()
			return
		}

	}

	erro = r.auditarAdministracion(result.Statement.Context, request)
	r.CommitTx()

	return
}

func (r *repository) UpdatePagoTipo(ctx context.Context, request entities.Pagotipo, channels administraciondtos.RequestPagoTipoChannels, cuotas administraciondtos.RequestPagoTipoCuotas) (erro error) {

	r.BeginTx()
	entidad := entities.Pagotipo{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_PAGO_TIPO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdatePagoTipo",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		r.RollbackTx()
		return
	}

	if len(channels.Add) > 0 {

		for _, ch := range channels.Add {
			entidadChannel := entities.Pagotipochannel{
				PagotiposId: entidad.ID,
				ChannelsId:  uint(ch),
			}

			resultpagotipochannels := r.SQLClient.WithContext(ctx).Omit("id").Create(&entidadChannel)

			if resultpagotipochannels.Error != nil {
				erro = fmt.Errorf(ERROR_CREAR_RUBRO)
				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       result.Error.Error(),
					Funcionalidad: "Createpagotipochannels",
				}

				err := r.utilService.CreateLogService(log)

				if err != nil {
					mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
					logs.Error(mensaje)
				}
				r.RollbackTx()
				return
			}
		}
	}
	if len(channels.Delete) > 0 {
		for _, cdelete := range channels.Delete {
			result := r.SQLClient.WithContext(ctx).Where("pagotipos_id = ? AND channels_id = ?", entidad.ID, cdelete).Delete(&entities.Pagotipochannel{})

			if result.Error != nil {

				erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       result.Error.Error(),
					Funcionalidad: "DeletePagoTipo",
				}

				err := r.utilService.CreateLogService(log)

				if err != nil {
					mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
					logs.Error(mensaje)
				}
				r.RollbackTx()
				return
			}
		}
	}
	if len(cuotas.Add) > 0 {

		for _, c := range cuotas.Add {
			entidadChannel := entities.Pagotipointallment{
				PagotiposId: entidad.ID,
				Cuota:       c,
			}

			resultpagotipochannels := r.SQLClient.WithContext(ctx).Omit("id").Create(&entidadChannel)

			if resultpagotipochannels.Error != nil {
				erro = fmt.Errorf(ERROR_CREAR_RUBRO)
				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       result.Error.Error(),
					Funcionalidad: "CreatePagotipointallment",
				}

				err := r.utilService.CreateLogService(log)

				if err != nil {
					mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
					logs.Error(mensaje)
				}
				r.RollbackTx()
				return
			}
		}
	}
	if len(cuotas.Delete) > 0 {
		for _, cudelete := range cuotas.Delete {
			result := r.SQLClient.WithContext(ctx).Where("pagotipos_id = ? AND cuota = ?", entidad.ID, cudelete).Delete(&entities.Pagotipointallment{})

			if result.Error != nil {

				erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

				log := entities.Log{
					Tipo:          entities.Error,
					Mensaje:       result.Error.Error(),
					Funcionalidad: "DeletePagoTipo",
				}

				err := r.utilService.CreateLogService(log)

				if err != nil {
					mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
					logs.Error(mensaje)
				}
				r.RollbackTx()
				return
			}
		}
	}
	erro = r.auditarAdministracion(result.Statement.Context, request)
	r.CommitTx()
	return
}

func (r *repository) DeletePagoTipo(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Pagotipo{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeletePagotipointallment",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}
	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

// ABM CHANNELS
func (r *repository) GetChannels(filtro filtros.ChannelFiltro) (response []entities.Channel, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Channel{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.Channel) > 0 {
		resp.Where("channel = ?", filtro.Channel)
	} else if len(filtro.Channels) > 0 {
		resp.Where("channel IN ?", filtro.Channels)
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetChannels",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetChannels: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetChannel(filtro filtros.ChannelFiltro) (channel entities.Channel, erro error) {

	resp := r.SQLClient.Model(entities.Channel{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.CargarMedioPago {
		resp.Preload("Mediopagos")
	}

	if len(filtro.Channel) > 0 {
		resp.Where("channel = ?", strings.ToUpper(filtro.Channel))
	} else if len(filtro.Channels) > 0 {
		resp.Where("channel IN ?", filtro.Channels)
	}

	resp.First(&channel)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) CreateChannel(ctx context.Context, request entities.Channel) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CHANNEL)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	id = uint64(request.ID)

	erro = r.auditarAdministracion(result.Statement.Context, request)

	return
}

func (r *repository) UpdateChannel(ctx context.Context, request entities.Channel) (erro error) {

	entidad := entities.Channel{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, request)
	return
}

func (r *repository) DeleteChannel(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Channel{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

// ABM CUENTA COMISION
func (r *repository) GetCuentasComisiones(filtro filtros.CuentaComisionFiltro) (response []entities.Cuentacomision, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Cuentacomision{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.CuentaId > 0 {
		resp.Where("cuentas_id = ?", filtro.CuentaId)
	}

	if filtro.ChannelId > 0 {
		resp.Where("channels_id = ?", filtro.ChannelId)
	}

	if filtro.CargarCuenta {
		resp.Preload("Cuenta")
	}

	if filtro.CargarChannel {
		resp.Preload("Channel")
	}

	if filtro.Channelarancel {
		resp.Preload("ChannelArancel")
	}

	if len(filtro.CuentaComision) > 0 {
		resp.Where("cuentacomision = ?", filtro.CuentaComision)
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCuentasComisiones",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCuentasComisiones: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) GetCuentaComision(filtro filtros.CuentaComisionFiltro) (response entities.Cuentacomision, erro error) {

	resp := r.SQLClient.Model(entities.Cuentacomision{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.CuentaComision) > 0 {
		resp.Where("cuentacomision = ?", filtro.CuentaComision)
	}

	if filtro.CuentaId > 0 {
		resp.Where("cuentas_id = ?", filtro.CuentaId)
	}
	// if filtro.ChannelId > 0 {
	// 	resp.Where("channels_id = ?", filtro.ChannelId).Order("vigencia_desde desc")
	// }

	if filtro.ChannelId > 0 {
		resp.Where("channels_id = ?", filtro.ChannelId)
	}

	if !filtro.FechaPagoVigencia.IsZero() {
		resp.Where("vigencia_desde <= ?", filtro.FechaPagoVigencia).Order("vigencia_desde desc")
	}

	if filtro.CargarCuenta {
		resp.Preload("Cuenta")
	}
	if filtro.CargarChannel {
		resp.Preload("Channel")
	}
	if filtro.Channelarancel {
		resp.Preload("ChannelArancel")
	}
	// if filtro.Mediopagoid > 0 {
	// 	resp.Where("mediopagoid = ?", filtro.Mediopagoid)
	// }

	if filtro.ExaminarPagoCuota {
		resp.Where("mediopagoid = ? AND pagocuota = ?", filtro.Mediopagoid, filtro.PagoCuota)
	}

	resp.First(&response)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCuentaComision",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) CreateCuentaComision(ctx context.Context, request entities.Cuentacomision) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CUENTA_COMISION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateCuentaComision",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(request.ID)

	erro = r.auditarAdministracion(result.Statement.Context, request)
	return
}

func (r *repository) UpdateCuentaComision(ctx context.Context, request entities.Cuentacomision) (erro error) {

	entidad := entities.Cuentacomision{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateCuentaComision",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, request)

	return
}

func (r *repository) GetImpuestosRepository(filtro filtros.ImpuestoFiltro) (response []entities.Impuesto, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Impuesto{})
	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}
	if len(filtro.Tipo) > 0 {
		resp.Where("tipo = ? and activo = ?", strings.ToUpper(filtro.Tipo), 1)
	}
	if filtro.OrdenarPorFecha {
		resp.Order("fechadesde asc")
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_IMPUESTO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetImpuestosRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetImpuestosRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) CreateImpuestoRepository(ctx context.Context, request entities.Impuesto) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_IMPUESTO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateImpuestoRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(request.ID)

	return
}

func (r *repository) UpdateImpuestoRepository(ctx context.Context, request entities.Impuesto) (erro error) {

	entidad := entities.Impuesto{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CHANNEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateImpuestoRepository",
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

func (r *repository) DeleteCuentaComision(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Cuentacomision{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteCuentaComision",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

func (r *repository) auditarAdministracion(ctx context.Context, resultado interface{}) error {
	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)

	if audit.Query == "" {
		audit.Operacion = "delete"
	} else {
		audit.Operacion = strings.ToLower(audit.Query[:6])
	}

	audit.Origen = "pasarela.administracion"

	res, _ := json.Marshal(resultado)
	audit.Resultado = string(res)

	err := r.auditoriaService.Create(&audit)

	if err != nil {
		return fmt.Errorf("auditoria: %w", err)
	}

	return nil
}

/* update pagosm notificados*/
func (r *repository) UpdatePagosNotificados(listaPagosNotificar []uint) (erro error) {

	result := r.SQLClient.Table("pagos").Where("id IN ?", listaPagosNotificar).Updates(map[string]interface{}{"notificado": 1})
	if result.Error != nil {
		erro := fmt.Errorf("no se puedo actualizar los pagos notificados")
		return erro
	}
	if result.RowsAffected <= 0 {
		logs.Info("caso de no actualizacion de pagos notificados 0")
		return nil
	}
	logs.Info("cantidad de pagos actualizados con exito " + fmt.Sprintf("%v", result.RowsAffected))

	return nil

}

// NOTE -Solo se mantendra hasta que se cree el proceso automatico con rabbit
func (r *repository) UpdatePagosEstadoInicialNotificado(listaPagosNotificar []uint) (erro error) {

	result := r.SQLClient.Table("pagos").Where("id IN ?", listaPagosNotificar).Updates(map[string]interface{}{"estado_inicial_notificado": 1})
	if result.Error != nil {
		erro := fmt.Errorf("no se puedo actualizar los pagos estado_inicial_notificado")
		return erro
	}
	if result.RowsAffected <= 0 {
		logs.Info("caso de no actualizacion de pagos notificados 0")
		return nil
	}
	logs.Info("cantidad de pagos actualizados con exito " + fmt.Sprintf("%v", result.RowsAffected))

	return nil

}
func (r *repository) GetConsultarMovimientosRapipago(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Rapipagocierrelote, erro error) {

	resp := r.SQLClient.Model(entities.Rapipagocierrelote{})

	if filtro.CargarMovConciliados {
		resp.Where("banco_external_id != ?", 0)
	} else {
		resp.Where("banco_external_id = ?", 0)
	}

	if filtro.PagosNotificado {
		resp.Where("pago_actualizado != ?", 0)
	} else {
		resp.Where("pago_actualizado = ?", 0)
	}

	resp.Preload("RapipagoDetalle")

	resp.Find(&response)
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_RAPIPAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConsultarMovimientosRapipago",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetConsultarMovimientosRapipago: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return response, erro
}

func (r *repository) GetConsultarMovimientosRapipagoDetalles(filtro rapipago.RequestConsultarMovimientosRapipagoDetalles) (response []*entities.Rapipagocierrelotedetalles, erro error) {
	resp := r.SQLClient.Model(entities.Rapipagocierrelotedetalles{})

	if filtro.PagosInformados {
		resp.Where("pagoinformado != ?", 0)
	} else {
		resp.Where("pagoinformado = ?", 0)
	}

	resp.Preload("RapipagoCabecera")

	resp.Find(&response)
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_RAPIPAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConsultarMovimientosRapipago",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetConsultarMovimientosRapipago: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return response, erro
}

func (r *repository) UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error) {

	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		for _, valueCL := range cierreLotes {
			resp := tx.Model(entities.Rapipagocierrelote{}).Where("id = ?", valueCL.ID).UpdateColumns(map[string]interface{}{"banco_external_id": valueCL.BancoExternalId, "enobservacion": valueCL.Enobservacion, "difbancocl": valueCL.Difbancocl})
			if resp.Error != nil {
				logs.Info(resp.Error)
				erro = errors.New("error: al actualizar tabla de cierre de lote rapipago")
				return erro
			}
		}

		for _, valueCL := range cierreLotes {
			for _, detalle := range valueCL.RapipagoDetalle {
				resp := tx.Model(entities.Rapipagocierrelotedetalles{}).Where("id = ?", detalle.ID).UpdateColumns(map[string]interface{}{"match": detalle.Match, "enobservacion": detalle.Enobservacion})
				if resp.Error != nil {
					logs.Info(resp.Error)
					erro = errors.New("error: al actualizar tabla de cierre de lote rapipago detalle")
					return erro
				}
			}
		}
		return nil
	})

	return
}

// func (r *repository) UpdateCierreloteAndMoviminetosRepository(entityCierreLote []entities.Prismacierrelote, listaIdsCabecera []int64, listaIdsDetalle []int64) (erro error) {
// 	r.SQLClient.Transaction(func(tx *gorm.DB) error {
// 		for _, valueCL := range entityCierreLote {
// 			resp := tx.Model(entities.Prismacierrelote{}).Where("id = ?", valueCL.ID).UpdateColumns(map[string]interface{}{"prismamovimientodetalles_id": valueCL.PrismamovimientodetallesId, "fecha_pago": valueCL.FechaPago})
// 			if resp.Error != nil {
// 				logs.Info(resp.Error)
// 				erro = errors.New("error: al actualizar tabla de cierre de lote")
// 				return erro
// 			}
// 		}

// 		if err := tx.Model(&entities.Prismamovimientototale{}).Where("id in (?)", listaIdsCabecera).UpdateColumns(map[string]interface{}{"match": 1}).Error; err != nil {
// 			logs.Info(err)
// 			erro = errors.New("error: al actualizar tabla Prisma Movimientos cabecera")
// 			return erro
// 		}
// 		if err := tx.Model(&entities.Prismamovimientodetalle{}).Where("id in (?)", listaIdsDetalle).UpdateColumns(map[string]interface{}{"match": 1}).Error; err != nil {
// 			logs.Info(err)
// 			erro = errors.New("error: al actualizar tabla Prisma Movimientos detalle")
// 			return erro
// 		}
// 		return nil
// 	})

// 	return
// }

func (r *repository) GetPeticionesWebServices(filtro filtros.PeticionWebServiceFiltro) (peticiones []entities.Webservicespeticione, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Webservicespeticione{})

	if filtro.Id != 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.OrdenarPorFechaInv {
		resp.Order("updated_at asc")
	}

	if filtro.Operacion != "" {
		resp.Where("operacion", filtro.Operacion)
	}
	if filtro.Vendor != "" {
		resp.Where("vendor = ?", filtro.Vendor)
	}

	if len(filtro.Fecha) > 0 {
		resp.Where("updated_at BETWEEN ? AND ?", filtro.Fecha[0], filtro.Fecha[1])
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))

	}

	resp.Find(&peticiones)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO_ESTADO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPeticionesWebServices",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return

}

func (r *repository) GetPagosTipoChannelRepository(filtro filtros.PagoTipoChannelFiltro) (pagostipochannel []entities.Pagotipochannel, erro error) {

	resp := r.SQLClient.Model(entities.Pagotipochannel{})

	if filtro.PagoTipoId > 0 {
		resp.Where("pagotipos_id", filtro.PagoTipoId)
	}

	if filtro.ChannelId > 0 {
		resp.Where("channels_id= ?", filtro.ChannelId)
	}

	resp.Find(&pagostipochannel)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO_ESTADO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagosTipoChannelRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return

}

func (r *repository) DeletePagoTipoChannel(id uint64) (erro error) {

	entidad := entities.Pagotipochannel{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteCuenta",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) CreatePagoTipoChannel(ctx context.Context, request entities.Pagotipochannel) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_IMPUESTO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreatePagoTipoChannel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(request.ID)

	return
}

// ABM CHANNELS ARANCELES
func (r *repository) GetChannelsAranceles(filtro filtros.ChannelArancelFiltro) (response []entities.Channelarancele, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Channelarancele{})

	if filtro.RubrosId > 0 {
		resp.Where("rubros_id = ?", filtro.RubrosId)
	}
	if filtro.CargarRubro {
		resp.Preload("Rubro")
	}
	if filtro.PagoCuota {
		resp.Where("pagocuota = ? ", filtro.PagoCuota)
	} else {
		resp.Where("pagocuota = ? ", filtro.PagoCuota)
	}
	if !filtro.CargarAllMedioPago {
		if filtro.MedioPagoId > 0 {
			resp.Where("mediopagoid = ?", filtro.MedioPagoId)
		}
		if filtro.MedioPagoId == 0 {
			resp.Where("mediopagoid = 0")
		}
	}

	if filtro.ChannelId > 0 {
		resp.Where("channels_id = ?", filtro.ChannelId)
	}
	if filtro.CargarChannel {
		resp.Preload("Channel")
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CARGAR_TOTAL_FILAS)
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Find(&response)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetChannelsAranceles",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetChannelsAranceles: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return
}

func (r *repository) CreateChannelsArancel(ctx context.Context, request entities.Channelarancele) (id uint64, erro error) {

	result := r.SQLClient.WithContext(ctx).Omit("id").Create(&request)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CUENTA_COMISION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateChannelsArancel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = uint64(request.ID)

	erro = r.auditarAdministracion(result.Statement.Context, request)
	return
}

func (r *repository) UpdateChannelsArancel(ctx context.Context, request entities.Channelarancele) (erro error) {

	entidad := entities.Channelarancele{
		Model: gorm.Model{ID: request.ID},
	}

	if entidad.ID == 0 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	result := r.SQLClient.WithContext(ctx).Model(&entidad).Omit("id,created_at,deleted_at").Select("*").Updates(request)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_MODIFICAR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "UpdateChannelsArancel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, request)

	return
}

func (r *repository) DeleteChannelsArancel(ctx context.Context, id uint64) (erro error) {

	entidad := entities.Channelarancele{
		Model: gorm.Model{ID: uint(id)},
	}

	result := r.SQLClient.WithContext(ctx).Delete(&entidad)

	if result.Error != nil {

		erro = fmt.Errorf(ERROR_BAJAR_CUENTA_COMISION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "DeleteChannelsArancel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	erro = r.auditarAdministracion(result.Statement.Context, id)
	return
}

func (r *repository) GetChannelArancel(filtro filtros.ChannelAranceFiltro) (response entities.Channelarancele, erro error) {

	resp := r.SQLClient.Model(entities.Channelarancele{})

	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if filtro.RubrosId > 0 {
		resp.Where("rubros_id = ?", filtro.RubrosId)
	}

	if filtro.ChannelId > 0 {
		resp.Where("channels_id = ?", filtro.ChannelId)
	}

	if filtro.CargarRubro {
		resp.Preload("Rubro")
	}
	if filtro.CargarChannel {
		resp.Preload("Channel")
	}

	if filtro.OrdernarChannel {
		resp.Order("fechadesde desc")
	}

	resp.First(&response)

	if resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}

		erro = fmt.Errorf(ERROR_CHANNEL_ARANCEL)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetChannelArancel",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetMedioPagoRepository(filtro filtros.FiltroMedioPago) (mediopago entities.Mediopago, erro error) {
	resp := r.SQLClient.Model(entities.Mediopago{})
	if filtro.IdMedioPago > 0 {
		resp.Where("id = ?", filtro.IdMedioPago)
	}
	resp.Find(&mediopago)
	if resp.Error != nil {
		erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMedioPagoRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	if resp.RowsAffected <= 0 {
		erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMedioPagoRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetCierreLoteSubidosRepository() (entityCl []entities.Prismacierrelote, erro error) {
	resp := r.SQLClient.Table("prismacierrelotes as cl")
	resp.Select("cl.nombrearchivolote, cl.created_at, cl.deleted_at")
	resp.Unscoped()
	resp.Group("cl.nombrearchivolote, cl.created_at")
	resp.Order("cl.created_at desc")
	resp.Find(&entityCl)
	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}
		erro = fmt.Errorf(ERROR_CIERRE_LOTE)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCierreLoteSubidosRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetPrismaPxSubidosRepository() (entityPx []entities.Prismapxcuatroregistro, erro error) {
	resp := r.SQLClient.Table("prismapxcuatroregistros as px")
	resp.Select("px.nombrearchivo, px.created_at, px.deleted_at")
	resp.Unscoped()
	resp.Group("px.nombrearchivo, px.created_at")
	resp.Order("px.created_at desc")
	resp.Find(&entityPx)
	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}
		erro = fmt.Errorf(ERROR_PRISMA_PX)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPrismaPxSubidosRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetPrismaMxSubidosRepository() (entityMx []entities.Prismamxtotalesmovimiento, erro error) {
	resp := r.SQLClient.Table("prismamxtotalesmovimientos as mx")
	resp.Select("mx.nombrearchivo, mx.created_at, mx.deleted_at")
	resp.Unscoped()
	resp.Group("mx.nombrearchivo, mx.created_at")
	resp.Order("mx.created_at desc")
	resp.Find(&entityMx)
	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			erro = fmt.Errorf(RESULTADO_NO_ENCONTRADO)
			return
		}
		erro = fmt.Errorf(ERROR_PRISMA_MX)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPrismaPxSubidosRepository",
		}
		err := r.utilService.CreateLogService(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) ObtenerArchivoCierreLoteRapipago(nombre string) (existeArchivo bool, erro error) {
	//consultar movimiento por nombre de archivo- verificar si existe
	//si existe retornar el nombre del archivo
	var result entities.Rapipagocierrelote
	res := r.SQLClient.Table("rapipagocierrelotes").Where("nombre_archivo = ?", nombre).Find(&result)

	if res.Error != nil {
		erro = res.Error
		return false, erro
	}
	if res.RowsAffected > 0 {
		existeArchivo = true
		return existeArchivo, nil
	} else {
		existeArchivo = false
		return existeArchivo, nil
	}
}

func (r *repository) ObtenerArchivoCierreLoteMultipagos(nombre string) (existeArchivo bool, erro error) {
	//consultar movimiento por nombre de archivo- verificar si existe
	//si existe retornar el nombre del archivo
	var result entities.Multipagoscierrelote
	res := r.SQLClient.Table("multipagoscierrelotes").Where("nombre_archivo = ?", nombre).Find(&result)

	if res.Error != nil {
		erro = res.Error
		return false, erro
	}
	if res.RowsAffected > 0 {
		existeArchivo = true
		return existeArchivo, nil
	} else {
		existeArchivo = false
		return existeArchivo, nil
	}
}

func (r *repository) ObtenerPagosInDisputaRepository(filtro filtros.ContraCargoEnDisputa) (pagosEnDisputa []entities.Pagointento, erro error) {
	resp := r.SQLClient.Table("pagointentos as pi")

	if len(filtro.TransactionId) > 0 {
		resp.Where("pi.transaction_id in (?)", filtro.TransactionId)

	}
	if filtro.CargarPagos {
		resp.Joins("inner join pagos as p on p.id = pi.pagos_id")
		resp.Preload("Pago")
	}

	if filtro.CargarTiposPago {
		resp.Joins("inner join pagotipos as ptip on ptip.id = p.pagostipo_id")
		resp.Preload("Pago.PagosTipo")
	}

	if filtro.CargarCuentas {
		resp.Joins("inner join cuentas as c on c.id = ptip.cuentas_id and c.id = ? and c.clientes_id = ? ", filtro.IdCuenta, filtro.IdCliente)
		resp.Preload("Pago.PagosTipo.Cuenta")
	}

	resp.Find(&pagosEnDisputa)
	if resp.Error != nil {
		erro = errors.New(ERROR_OBTENER_PAGOS_DISPUTA)
	}

	return
}

func (r *repository) PostPreferencesRepository(preferenceEntity entities.Preference) (erro error) {
	err := r.SQLClient.Create(&preferenceEntity).Error
	if err != nil {
		erro = errors.New("no se pudo guardar la preferencia en la base de datos")
		return
	}

	return
}

/* update pagosm notificados*/
func (r *repository) UpdatePagosDev(pagos []uint) (erro error) {

	result := r.SQLClient.Table("pagos").Where("id IN ?", pagos).Updates(map[string]interface{}{"pagoestados_id": 7})
	if result.Error != nil {
		erro := fmt.Errorf("no se puedo actualizar los pagos")
		return erro
	}
	if result.RowsAffected <= 0 {
		logs.Info("caso de no actualizacion de pagos notificados 0")
		return nil
	}
	logs.Info("cantidad de pagos actualizados con exito " + fmt.Sprintf("%v", result.RowsAffected))

	return nil
}

func (r *repository) CreateSolicitudRepository(solicitud entities.Solicitud) (erro error) {
	// Guardar la solicitud en la base de datos en su tabla correspondiente
	result := r.SQLClient.Model(entities.Solicitud{}).Create(&solicitud)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_SOLICITUD_CUENTA)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateSolicitud",
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

// ? consultar clrapipago para herramietna wee
func (r *repository) GetConsultarClRapipagoRepository(filtro filtros.RequestClrapipago) (clrapiapgo []entities.Rapipagocierrelote, totalFilas int64, erro error) {

	resp := r.SQLClient.Unscoped().Model(entities.Rapipagocierrelote{})

	if filtro.CodigoBarra != "" {
		resp.Joins("INNER JOIN rapipagocierrelotedetalles as rpcl_detalles on rpcl_detalles.rapipagocierrelotes_id = rapipagocierrelotes.id")
		resp.Where("rpcl_detalles.codigo_barras = ? ", filtro.CodigoBarra)
	} else {

		if !filtro.FechaInicio.IsZero() && !filtro.FechaFin.IsZero() {
			resp.Where("cast(rapipagocierrelotes.created_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaInicio, filtro.FechaFin)
		}
	}
	resp.Preload("RapipagoDetalle", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	})

	// if filtro.Id > 0 {
	// 	resp.Where("id = ?", filtro.Id)
	// }

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf("error al cargar total filas de la columna")
		}

	}

	if filtro.Number > 0 && filtro.Size > 0 {
		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))
	}

	resp.Order("created_at desc").Find(&clrapiapgo)

	return
}

func (r *repository) GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error) {

	resp := r.SQLClient.Model(entities.Pagointento{})

	if !filtro.FechaPagoFin.IsZero() {
		resp.Where("cast(pagointentos.paid_at as date) BETWEEN cast(? as date) AND cast(? as date)", filtro.FechaPagoInicio, filtro.FechaPagoFin)

	}

	if filtro.CargarPagoCalculado {
		resp.Where("calculado = ?", 1)
	} else {
		resp.Where("calculado = ?", 0)
	}

	if filtro.ClienteId > 0 {
		resp.Preload("Pago.PagosTipo.Cuenta", "cuentas.clientes_id = ?", filtro.ClienteId).Joins("INNER JOIN pagos as p on p.id = pagointentos.pagos_id INNER JOIN pagotipos as pt on pt.id = p.pagostipo_id INNER JOIN cuentas as cu on cu.id = pt.cuentas_id").Where("cu.clientes_id = ?", filtro.ClienteId)
	}

	if len(filtro.PagoEstadosIds) > 0 {
		resp.Preload("Pago", "pagos.pagoestados_id IN ?", filtro.PagoEstadosIds).Joins("INNER JOIN pagos as pg on pg.id = pagointentos.pagos_id").Where("pg.pagoestados_id IN ?", filtro.PagoEstadosIds)
		// resp.Where("pagoestados_id IN ?", filtro.PagoEstadosIds)
	}
	if filtro.PagoIntentoAprobado {

		resp.Where("paid_at <> ?", "0000-00-00 00:00:00")
	}

	if filtro.CargarMovimientosTemporales {
		resp.Preload("Movimientotemporale.Movimientocomisions")
		resp.Preload("Movimientotemporale.Movimientoimpuestos")
	}

	if filtro.CargarPagoItems {
		resp.Preload("Pago.Pagoitems")
	}

	if filtro.Channel {
		resp.Preload("Mediopagos.Channel")
	}
	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetPagosApilink(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error) {
	resp := r.SQLClient.Model(entities.Pago{})
	resp.Joins("INNER JOIN pagointentos PI ON pagos.id = PI.pagos_id")
	resp.Joins("INNER JOIN apilinkcierrelotes ACL ON ACL.debin_id = PI.external_id")

	resp.Select("pagos.id")
	resp.Where("date(ACL.fecha_cobro) BETWEEN date(?) AND date(?)", filtro.FechaPagoInicio, filtro.FechaPagoFin)
	resp.Where("PI.calculado != 1")
	resp.Find(&ids)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
			return
		}
	}

	return
}

func (r *repository) GetPagosRapipago(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error) {
	resp := r.SQLClient.Model(entities.Pago{})
	resp.Joins("INNER JOIN pagointentos PI ON pagos.id = PI.pagos_id")
	resp.Joins("INNER JOIN rapipagocierrelotedetalles RP ON RP.codigo_barras = PI.barcode")

	resp.Select("pagos.id")
	resp.Where("date(RP.fecha_cobro) BETWEEN date(?) AND date(?)", filtro.FechaPagoInicio, filtro.FechaPagoFin)
	resp.Where("PI.calculado != 1")
	resp.Find(&ids)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
			return
		}
	}

	return
}

func (r *repository) GetPagosPrisma(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error) {
	resp := r.SQLClient.Model(entities.Pago{})
	resp.Joins("INNER JOIN pagointentos PI ON pagos.id = PI.pagos_id")

	resp.Select("pagos.id")
	resp.Where("PI.card_last_four_digits != '' AND PI.paid_at IS NOT NULL AND PI.calculado != 1")
	resp.Where("date(PI.paid_at) BETWEEN date(?) AND date(?)", filtro.FechaPagoInicio, filtro.FechaPagoFin)
	resp.Find(&ids)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetPagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
			return
		}
	}

	return
}

func (r *repository) GetSuccessPaymentsRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error) {

	resp := r.SQLClient.Model(entities.Pago{})

	// esto hace que cargue los pagointentos del pago, pero solo aquellos que tienen algun valor setado en el campo paid_at, y no los que tienen todo en cero ese campo
	resp.Preload("PagoIntentos", "paid_at")
	resp.Joins("INNER JOIN pagoestados AS PE ON pagos.pagoestados_id = PE.id INNER JOIN pagointentos AS PI ON pagos.id = PI.pagos_id INNER JOIN pagotipos AS PT ON pagos.pagostipo_id = PT.id INNER JOIN cuentas AS C ON PT.cuentas_id = C.id INNER JOIN pagoitems AS PIT ON pagos.id = PIT.pagos_id")
	// es necesario el DISTINCT porque al hacer join con pagoitems se duplican los registros del resultado de la consulta
	resp.Select("DISTINCT pagos.*")
	// Los estados de pagos exitosos son 4 y 7
	resp.Where("pagos.pagoestados_id IN ?", []uint{4, 7}).Where("PI.paid_at LIKE ?", filtro.Fecha[0]+"%").Where("C.id = ?", filtro.CuentaId)

	resp.Find(&pagos)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetSuccessPaymentsRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetReportesPagoRepository(filtro filtros.PagoFiltro) (reportes entities.Reporte, erro error) {

	resp := r.SQLClient.Model(entities.Reporte{})

	filtroCuenta := filtros.CuentaFiltro{
		Id: uint(filtro.CuentaId),
	}

	// Averiguar el cliente segun la cuenta
	cuenta, err := r.GetCuenta(filtroCuenta)

	if err != nil {
		mensaje := fmt.Sprintf("se produjo el siguiente error: %s.", err.Error())
		logs.Error(mensaje)
		return
	}

	// si no existe la cuenta, el objeto se encuentra vacio
	if cuenta.ID == 0 {
		erro = fmt.Errorf("no se encuentra la cuenta con el id requerido")
		logs.Error(erro.Error())
		return
	}

	// El nombre del cliente para pregunatr en la tabla reportes por ese cliente
	cliente := cuenta.Cliente.Cliente

	resp.Preload("Reportedetalle")

	resp.Where("tiporeporte = ? AND cliente = ? AND fechacobranza = ?", "pagos", cliente, filtro.Fecha)

	resp.Find(&reportes)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_PAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetReportesPagoRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetRetencionesRepository(request retenciondtos.RentencionRequestDTO) (retenciones []entities.Retencion, count int64, erro error) {

	resp := r.SQLClient.Model(entities.Retencion{})

	resp.Preload("Channel")
	resp.Preload("Condicion.Gravamen")

	if request.RetencionId != 0 {
		resp.Where("id = ?", request.RetencionId)
	}

	resp.Find(&retenciones).Count(&count)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CARGAR_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetRetencionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetRetencionesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetRetencionClienteRepository(request retenciondtos.RentencionRequestDTO) (clienteretencion entities.ClienteRetencion, erro error) {

	resp := r.SQLClient.Model(entities.ClienteRetencion{})
	resp.Preload("Cliente")
	resp.Where("cliente_id = ?", request.ClienteId)
	resp.Where("retencion_id = ?", request.RetencionId)

	resp.First(&clienteretencion)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf("error al obtener las retenciones-clientes")
	}

	return
}

func (r *repository) GetCertificadoRepository(requestId uint) (certificado entities.Certificado, erro error) {

	resp := r.SQLClient.Model(entities.Certificado{})
	resp.Where("id = ?", requestId)
	resp.First(&certificado)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf("error al obtener certificado")

	}

	return
}

func (r *repository) GetCertificadosVencimientoRepository(request retenciondtos.CertificadoVencimientoDTO) (certificados []entities.Certificado, erro error) {

	resp := r.SQLClient.Table("certificados")

	resp.Select("*,cliente_retencions_id as retencion_id, (SELECT max(fecha_caducidad) FROM certificados WHERE cliente_retencions_id = retencion_id ORDER BY fecha_caducidad DESC) as fecha_caducidad")
	resp.Group("cliente_retencions_id")
	resp.Preload("ClienteRetencion.Cliente")
	resp.Preload("ClienteRetencion.Retencion.Condicion.Gravamen")

	// fechaComparacion := time.Now()
	// if !request.Fecha_comparacion.IsZero() {
	// 	fechaComparacion = request.Fecha_comparacion
	// }
	// resp.Where("fecha_caducidad < ?", fechaComparacion)

	if len(request.Ids) > 0 {
		resp.Where("id in ?", request.Ids)
	}

	resp.Find(&certificados)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf("error al obtener certificado")

	}

	return
}

func (r *repository) GetClienteRetencionesRepository(request retenciondtos.RentencionRequestDTO) (clienteretencion []entities.ClienteRetencion, count int64, erro error) {

	resp := r.SQLClient.Model(&entities.ClienteRetencion{}).Where("cliente_id = ?", request.ClienteId).Preload("Certificados").Preload("Retencion.Condicion.Gravamen").Preload("Retencion.Channel").Find(&clienteretencion)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CARGAR_CLIENTE_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetClienteRetencionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetClienteRetencionesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	count = int64(len(clienteretencion))

	return
}

func (r *repository) GetClienteUnlinkedRetencionesRepository(request retenciondtos.RentencionRequestDTO) (retenciones []entities.Retencion, count int64, erro error) {

	subQuery := r.SQLClient.Table("cliente_retencions").Where("cliente_id = (?) and deleted_at IS NULL", request.ClienteId).Select("retencion_id")

	resp := r.SQLClient.Where("id NOT IN (?)", subQuery).Preload("Condicion.Gravamen").Preload("Channel").Find(&retenciones)

	if len(retenciones) != 0 {
		count = int64(len(retenciones))
	}

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CARGAR_CLIENTE_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetClienteRetencionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetClienteRetencionesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) CreateClienteRetencionRepository(request retenciondtos.RentencionRequestDTO) (erro error) {

	existe_cliente, erro := r.checkIfModelsExistsByValue("clientes", "id", request.ClienteId)
	if erro != nil {
		erro = errors.New("el registro de cliente no se encuentra")
		return
	}
	existe_retencion, erro := r.checkIfModelsExistsByValue("retencions", "id", request.RetencionId)
	if erro != nil {
		erro = errors.New("el registro de retencion no se encuentra")
		return
	}
	if !existe_cliente || !existe_retencion {
		erro = errors.New("alguno de los registros que intenta vincular no existe")
		return
	}
	// crear un objeto de la tabla intermedia
	cliente_retencion := entities.ClienteRetencion{
		ClienteId:   request.ClienteId,
		RetencionId: request.RetencionId,
	}
	// crear
	resp := r.SQLClient.Create(&cliente_retencion)

	// manejar error en caso de que no sea por clave duplicada
	if resp.Error != nil && !strings.Contains(resp.Error.Error(), "Duplicate entry") {
		erro = fmt.Errorf(ERROR_CREAR_CLIENTE_RETENCION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "CreateClienteRetencionRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. CreateClienteRetencionRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	// caso en que ya existe el registro pero esta borrado logicamente y se desea volver a crearlo
	if resp.Error != nil && strings.Contains(resp.Error.Error(), "Duplicate entry") {
		resp := r.SQLClient.Unscoped().Model(&cliente_retencion).Where("retencion_id = ? and cliente_id = ?", request.RetencionId, request.ClienteId).Update("deleted_at", nil)

		if resp.Error != nil {
			erro = fmt.Errorf(ERROR_CREAR_CLIENTE_RETENCION)
			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       resp.Error.Error(),
				Funcionalidad: "CreateClienteRetencionRepository",
			}

			err := r.utilService.CreateLogService(log)

			if err != nil {
				mensaje := fmt.Sprintf("Crear Log: %s. CreateClienteRetencionRepository: %s", err.Error(), resp.Error.Error())
				logs.Error(mensaje)
			}
			return
		}
	}

	return
}

func (r *repository) PostRetencionCertificadoRepository(certificado entities.Certificado) (erro error) {
	err := r.SQLClient.Create(&certificado).Error
	if err != nil {
		erro = errors.New("no se pudo guardar el certificado en la base de datos")
		return
	}

	return
}

func (r *repository) checkIfModelsExistsByValue(tableName string, column string, value interface{}) (exists bool, erro error) {

	erro = r.SQLClient.Table(tableName).Select("count(*) > 0").Where(column+" = ?", value).Find(&exists).Error

	if erro != nil {
		erro = errors.New(ERROR_DB_ACCESS)
	}
	return
}

func (r *repository) GetCalcularRetencionesRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {

	resp := r.SQLClient.Table("movimiento_retencions AS MR").
		Select("MR.cliente_id, SUM(MR.importe_retenido) AS total_retencion, G.gravamen, MR.efectuada as efectuada, R.codigo_regimen as codigo_regimen, R.monto_minimo as minimo").
		Joins("JOIN retencions AS R ON MR.retencion_id = R.id").
		Joins("JOIN condicions AS C ON R.condicions_id = C.id").
		Joins("JOIN gravamens AS G ON C.gravamens_id = G.id").
		Where("MR.created_at BETWEEN ? AND ?", request.FechaInicio, request.FechaFin).
		Where("MR.deleted_at IS NULL").
		Group("MR.cliente_id, G.gravamen, efectuada, codigo_regimen, minimo")

	if request.ClienteId > 0 {
		resp.Where("MR.cliente_id", request.ClienteId)
	}

	if request.GravamensId > 0 {
		resp.Where("G.id", request.GravamensId)
	}

	resp.Scan(&resultado)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CALCULAR_CLIENTE_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCalcularRetencionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCalcularRetencionesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetCalcularRetencionesByTransferenciasRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {

	subquery2 := r.SQLClient.Table("transferencias").Select("pagointentos.id").
		Joins("JOIN movimientos ON transferencias.movimientos_id = movimientos.id").
		Joins("JOIN pagointentos ON movimientos.pagointentos_id = pagointentos.id").
		Joins("JOIN cuentas ON cuentas.id = movimientos.cuentas_id").
		Joins("JOIN clientes ON clientes.id = cuentas.clientes_id").
		Where("transferencias.fecha_operacion BETWEEN ? AND ?", request.FechaInicio, request.FechaFin).
		Where("cuentas.clientes_id", request.ClienteId)

	subquery1 := r.SQLClient.Table("movimientos").Select("movimientos.id").
		Where("pagointentos_id IN (?)", subquery2).
		Where("movimientos.tipo = 'C'")

	resp := r.SQLClient.Table("movimiento_retencions AS MR").
		Select("MR.cliente_id, SUM(MR.importe_retenido) AS total_retencion, SUM(MR.monto) AS total_monto, G.gravamen, MR.efectuada as efectuada, R.codigo_regimen as codigo_regimen, R.monto_minimo as minimo").
		Joins("JOIN retencions AS R ON MR.retencion_id = R.id").
		Joins("JOIN condicions AS C ON R.condicions_id = C.id").
		Joins("JOIN gravamens AS G ON C.gravamens_id = G.id").
		Where("MR.movimiento_id IN (?)", subquery1).
		Where("MR.deleted_at IS NULL").
		Group("MR.cliente_id, G.gravamen, efectuada, codigo_regimen, minimo")

	resp.Scan(&resultado)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CALCULAR_CLIENTE_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCalcularRetencionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCalcularRetencionesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetMovimientosIdsCalculoRetencionComprobante(request retenciondtos.RentencionRequestDTO) (resultado []uint, erro error) {

	subquery2 := r.SQLClient.Table("transferencias").Select("pagointentos.id").
		Joins("JOIN movimientos ON transferencias.movimientos_id = movimientos.id").
		Joins("JOIN pagointentos ON movimientos.pagointentos_id = pagointentos.id").
		Joins("JOIN cuentas ON cuentas.id = movimientos.cuentas_id").
		Joins("JOIN clientes ON clientes.id = cuentas.clientes_id").
		Where("transferencias.fecha_operacion BETWEEN ? AND ?", request.FechaInicio, request.FechaFin).
		Where("cuentas.clientes_id = ?", request.ClienteId)

	subquery1 := r.SQLClient.Table("movimientos").Select("movimientos.id").
		Where("pagointentos_id IN (?)", subquery2).
		Where("movimientos.tipo = 'C'")

	// obtener los movimientos ids implicados en el calculo del comprobante de retencion
	resp := r.SQLClient.Table("movimiento_retencions AS MR").
		Select("MR.movimiento_id").
		Joins("JOIN retencions AS R ON MR.retencion_id = R.id").
		Joins("JOIN condicions AS C ON R.condicions_id = C.id").
		Joins("JOIN gravamens AS G ON C.gravamens_id = G.id").
		Where("MR.movimiento_id IN (?)", subquery1).
		Where("MR.deleted_at IS NULL")

	resp.Scan(&resultado)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CALCULAR_CLIENTE_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMovimientosIdsCalculoRetencionComprobante",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetMovimientosIdsCalculoRetencionComprobante: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GetGravamensRepository(filtro retenciondtos.GravamenRequestDTO) (gravamenes []entities.Gravamen, erro error) {

	resp := r.SQLClient.Model(&entities.Gravamen{})
	if filtro.Gravamen != "" {
		resp.Where("gravamens.gravamen", filtro.Gravamen)
	}

	resp.Find(&gravamenes)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CARGAR_GRAVAMENES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetGravamensRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetGravamensRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetCondicionesRepository(request retenciondtos.RentencionRequestDTO) (condicions []entities.Condicion, erro error) {

	resp := r.SQLClient.Model(&entities.Condicion{})

	if request.GravamensId > 0 {
		resp.Where("gravamens_id = ?", request.GravamensId)
	}

	if request.CargarGravamenes {
		resp.Preload("Gravamen")
	}

	resp.Find(&condicions)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CARGAR_CONDICIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCondicionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCondicionesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) CreateRetencionRepository(entity entities.Retencion) (entities.Retencion, error) {

	existe_condicion, erro := r.checkIfModelsExistsByValue("condicions", "id", entity.CondicionsId)
	if erro != nil {
		return entities.Retencion{}, errors.New("el registro de condicion impositiva no se encuentra")

	}
	existe_channel, erro := r.checkIfModelsExistsByValue("channels", "id", entity.ChannelsId)
	if erro != nil {
		return entities.Retencion{}, errors.New("el registro del canal de pago no se encuentra")

	}
	if !existe_condicion || !existe_channel {
		return entities.Retencion{}, errors.New("debe enviar una condicion impositiva y un canal de pago válidos")

	}

	resp := r.SQLClient.Create(&entity)

	if resp.RowsAffected <= 0 {
		erro = fmt.Errorf(ERROR_CREAR_RETENCION)
		return entities.Retencion{}, erro
	}

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_RETENCION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "CreateRetencionRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. CreateRetencionRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
		return entities.Retencion{}, erro
	}

	return entity, nil
}

func (r *repository) DeleteClienteRetencionRepository(request retenciondtos.RentencionRequestDTO) (erro error) {

	existe_cliente, erro := r.checkIfModelsExistsByValue("clientes", "id", request.ClienteId)
	if erro != nil {
		erro = errors.New("el registro de cliente no se encuentra")
		return
	}
	existe_retencion, erro := r.checkIfModelsExistsByValue("retencions", "id", request.RetencionId)
	if erro != nil {
		erro = errors.New("el registro de retencion no se encuentra")
		return
	}
	if !existe_cliente || !existe_retencion {
		erro = errors.New("alguno de los registros que intenta eliminar no existe")
		return
	}
	// crear un objeto de la tabla intermedia
	cliente_retencion := entities.ClienteRetencion{
		ClienteId:   request.ClienteId,
		RetencionId: request.RetencionId,
	}
	resp := r.SQLClient.Where("cliente_id = ? and retencion_id = ?", cliente_retencion.ClienteId, cliente_retencion.RetencionId).Delete(&cliente_retencion)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_DELETE_RETENCION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "DeleteClienteRetencionService",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. DeleteClienteRetencionService: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	return
}

func (r *repository) UpdateRetencionRepository(entity entities.Retencion) (entities.Retencion, error) {
	var erro error

	existe_retencion, erro := r.checkIfModelsExistsByValue("retencions", "id", entity.ID)

	if erro != nil {
		return entities.Retencion{}, erro
	}

	if !existe_retencion {
		erro = errors.New("el registro que intenta modificar no existe")
		return entities.Retencion{}, erro
	}

	resp := r.SQLClient.Omit("created_at").Save(&entity)

	if resp.RowsAffected <= 0 {
		erro = fmt.Errorf(ERROR_UPDATE_RETENCION)
		return entities.Retencion{}, erro
	}

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_UPDATE_RETENCION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "UpdateRetencionRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. UpdateRetencionRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
		return entities.Retencion{}, erro
	}

	return entity, nil
}

func (r *repository) CreateCondicionRepository(condicion entities.Condicion) (erro error) {

	resp := r.SQLClient.Create(&condicion)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CONDICION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "CreateCondicionRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. CreateCondicionRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	return
}

func (r *repository) UpdateCondicionRepository(condicion entities.Condicion) (erro error) {

	existe_condicion, erro := r.checkIfModelsExistsByValue("condicions", "id", condicion.ID)

	if erro != nil {
		return
	}

	if !existe_condicion {
		erro = errors.New("el registro que intenta modificar no existe")
		return
	}

	resp := r.SQLClient.Omit("created_at").Save(&condicion)

	if resp.RowsAffected <= 0 {
		erro = fmt.Errorf(ERROR_UPDATE_CONDICION)
		return
	}

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_UPDATE_RETENCION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "UpdateCondicionRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. UpdateCondicionRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
		return
	}

	return
}

func (r *repository) GetClienteRetencionRepository(retencion_id, cliente_id uint) (cliente_retencion entities.ClienteRetencion, erro error) {

	resp := r.SQLClient.Model(&entities.ClienteRetencion{})

	if retencion_id > 0 {
		resp.Where("retencion_id = ?", retencion_id)
	}

	if cliente_id > 0 {
		resp.Where("cliente_id = ?", cliente_id)
	}

	resp.Find(&cliente_retencion)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_GET_CLIENTE_RETENCION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetClienteRetencionRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetClienteRetencionRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) EvaluarRetencionesByClienteRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {

	query := fmt.Sprintf("cliente_id, SUM(MR.importe_retenido) AS total_retencion, R.monto_minimo AS minimo, G.gravamen as gravamen, '%s' as fecha_inicio, '%s' as fecha_fin", request.FechaInicio, request.FechaFin)

	resp := r.SQLClient.Table("movimiento_retencions AS MR").
		Joins("JOIN retencions AS R ON MR.retencion_id = R.id").
		Joins("JOIN condicions AS C ON R.condicions_id = C.id").
		Joins("JOIN gravamens AS G ON C.gravamens_id = G.id").
		Select(query).
		Where("MR.created_at BETWEEN ? AND ?", request.FechaInicio, request.FechaFin).
		// Where("R.monto_minimo > 0").
		Where("MR.cliente_id = ?", request.ClienteId).
		Where("C.exento = 0"). // where no esta exento del gravamen
		Group("minimo, gravamen")

	resp.Scan(&resultado)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CALCULAR_CLIENTE_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "EvaluarRetencionesByClienteRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. EvaluarRetencionesByClienteRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) EvaluarRetencionesByMovimientosRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {

	query := "cliente_id, SUM(MR.importe_retenido) AS total_retencion, G.gravamen as gravamen, R.monto_minimo AS minimo, MR.efectuada as efectuada"

	resp := r.SQLClient.Table("movimiento_retencions AS MR").
		Joins("JOIN retencions AS R ON MR.retencion_id = R.id").
		Joins("JOIN condicions AS C ON R.condicions_id = C.id").
		Joins("JOIN gravamens AS G ON C.gravamens_id = G.id").
		Select(query).
		Where("MR.cliente_id = ?", request.ClienteId).
		Where("MR.movimiento_id IN ?", request.ListaMovimientosId).
		Where("C.exento = 0"). // where no esta exento del gravamen
		Group("gravamen, minimo, efectuada")

	resp.Scan(&resultado)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CALCULAR_CLIENTE_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "EvaluarRetencionesByClienteRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. EvaluarRetencionesByClienteRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) GenerarCertificacionRepository(comprobantes []entities.Comprobante) (erro error) {
	for _, c := range comprobantes {
		erro = r.SQLClient.Create(&c).Error
		if erro != nil {
			break
		}
		// el numero del comprobante se genera a partir del id AI del comprobante
		numero := r.utilService.GenerarNumeroComprobante1(c.ID)
		erro = r.SQLClient.Model(&c).UpdateColumn("numero", numero).Error
		if erro != nil {
			break
		}
	}

	if erro != nil {
		erro = fmt.Errorf(ERROR_CREAR_CERTIFICACION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       erro.Error(),
			Funcionalidad: "GenerarCertificacionRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GenerarCertificacionRepository: %s", err.Error(), erro.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetMovimientosRetencionesRepository(request retenciondtos.RentencionRequestDTO) (listaMovimientosId []uint, erro error) {
	resp := r.SQLClient.Model(&entities.MovimientoRetencion{})

	if request.ClienteId > 0 {
		resp.Where("cliente_id = ?", request.ClienteId)
	}

	if len(request.FechaInicio) > 0 && len(request.FechaFin) > 0 {
		resp.Where("created_at BETWEEN ? AND ?", request.FechaInicio, request.FechaFin)
	}

	resp.Select("movimiento_id").Distinct().Find(&listaMovimientosId)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_GET_MOVIMIENTOS_RETENCION)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetMovimientosRetencionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetMovimientosRetencionesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetTotalAmountByMovimientoIdsRepository(listaMovimientosId []uint) (totalAmount uint64, erro error) {

	resp := r.SQLClient.Table("movimiento_retencions AS MR").
		Select("SUM(MR.monto)").
		Where("MR.movimiento_id IN ?", listaMovimientosId)

	resp.Scan(&totalAmount)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_GET_TOTAL_AMOUNT)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetTotalAmountByMovimientoIdsRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetTotalAmountByMovimientoIdsRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) ComprobantesRetencionesDevolverRepository(request retenciondtos.RentencionRequestDTO) (comprobantes []entities.Comprobante, erro error) {

	resp := r.SQLClient.Preload("ComprobanteDetalles", "created_at BETWEEN ? AND ?", request.FechaInicio, request.FechaFin)

	if request.ClienteId > 0 {
		resp.Where("cliente_id = ?", request.ClienteId)
	}

	resp.Find(&comprobantes)

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_RETENCIONES_DEVOLVER)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "RetencionesDevolverRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. RetencionesDevolverRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	var comprobantesDevolver []entities.Comprobante
	for _, c := range comprobantes {
		for _, detalle := range c.ComprobanteDetalles {
			if !detalle.Retener {
				comprobantesDevolver = append(comprobantesDevolver, c)
				break
			}
		}
	}
	comprobantes = comprobantesDevolver
	return
}

func (r *repository) TotalizarRetencionesMovimientosRepository(listaMovimientoIds []uint) (totalRetenciones uint64, erro error) {

	resp := r.SQLClient.Model(&entities.MovimientoRetencion{}).Select("sum(importe_retenido) as totalRetenciones")

	if len(listaMovimientoIds) > 0 {
		resp.Where("movimiento_id IN ?", listaMovimientoIds)
	}

	resp.Find(&totalRetenciones)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_GET_TOTAL_AMOUNT)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "TotalizarRetencionesMovimientosRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. TotalizarRetencionesMovimientosRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *repository) GetComprobantesRepository(request retenciondtos.RentencionRequestDTO) (comprobantes []entities.Comprobante, erro error) {
	// se usa para el caso de que el campo emitido_el este en tiempo cero
	// zeroTime := time.Time{}
	resp := r.SQLClient.Model(entities.Comprobante{})
	resp2 := r.SQLClient.Model(entities.Reporte{})
	if request.ClienteId > 0 {
		resp.Where("cliente_id = ?", request.ClienteId)
	}

	if request.ComprobanteId > 0 {
		resp.Where("id = ?", request.ComprobanteId)
	}

	resp.Preload("ComprobanteDetalles")

	resp.Find(&comprobantes)

	var reporte entities.Reporte
	var nro_reporte_rrm uint
	for i := range comprobantes {
		if nro_reporte_rrm == comprobantes[i].ReporteId {
			comprobantes[i].Reporte = reporte
		} else {
			nro_reporte_rrm = comprobantes[i].ReporteId
			resp2.Where("nro_reporte = ? AND tiporeporte = ?", comprobantes[i].ReporteId, "rrm")
			resp2.Find(&reporte)
			comprobantes[i].Reporte = reporte
		}
	}

	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_GET_COMPROBANTES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetComprobantesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetComprobantesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}

func (r *repository) UpdateMovimientoMontoRepository(ctx context.Context, movimiento entities.Movimiento) (erro error) {
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {

		resp := tx.WithContext(ctx).Model(&movimiento).Update("monto", movimiento.Monto)
		if resp.Error != nil {
			return resp.Error
		}
		return nil
	})
}

func (r *repository) GetClientesConfiguracion(filtro filtros.ClienteConfiguracionFiltro) (clientes []entities.Cliente, erro error) {
	resp := r.SQLClient.Model(entities.Cliente{})
	resp.Where("configuracion_retiro_automatico = ?", filtro.Configuracion).Find(&clientes)

	if resp.RowsAffected <= 0 {
		erro = errors.New("No se encontraron registros")
		return
	}

	return
}

func (r *repository) CreateAuditoria(resultado entities.Auditoria) (erro error) {

	err := r.auditoriaService.Create(&resultado)

	if err != nil {
		return fmt.Errorf("auditoria: %w", err)
	}

	return nil
}

func (r *repository) GetConsultarMovimientosMultipagos(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Multipagoscierrelote, erro error) {

	resp := r.SQLClient.Model(entities.Multipagoscierrelote{})

	if filtro.CargarMovConciliados {
		resp.Where("banco_external_id != ?", 0)
	} else {
		resp.Where("banco_external_id = ?", 0)
	}

	if filtro.PagosNotificado {
		resp.Where("pago_actualizado != ?", 0)
	} else {
		resp.Where("pago_actualizado = ?", 0)
	}

	resp.Preload("MultipagosDetalle")

	resp.Find(&response)
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_RAPIPAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConsultarMovimientosMultipagos",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetConsultarMovimientosMultipagos: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return response, erro
}

func (r *repository) GetConsultarMovimientosMultipagosDetalles(filtro multipagos.RequestConsultarMovimientosMultipagosDetalles) (response []*entities.Multipagoscierrelotedetalles, erro error) {
	resp := r.SQLClient.Model(entities.Multipagoscierrelotedetalles{})

	if filtro.PagosInformados {
		resp.Where("pagoinformado != ?", 0)
	} else {
		resp.Where("pagoinformado = ?", 0)
	}

	resp.Preload("MultipagosCabecera")

	resp.Find(&response)
	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_RAPIPAGO)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConsultarMovimientosMultipagosDetalles",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetConsultarMovimientosMultipagosDetalles: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}

	}

	return response, erro
}

func (r *repository) UpdateCierreLoteMultipagos(cierreLotes []*entities.Multipagoscierrelote) (erro error) {

	r.SQLClient.Transaction(func(tx *gorm.DB) error {
		for _, valueCL := range cierreLotes {
			resp := tx.Model(entities.Multipagoscierrelote{}).Where("id = ?", valueCL.ID).UpdateColumns(map[string]interface{}{"banco_external_id": valueCL.BancoExternalId, "enobservacion": valueCL.Enobservacion, "difbancocl": valueCL.Difbancocl})
			if resp.Error != nil {
				logs.Info(resp.Error)
				erro = errors.New("error: al actualizar tabla de cierre de lote rapipago")
				return erro
			}
		}

		for _, valueCL := range cierreLotes {
			for _, detalle := range valueCL.MultipagosDetalle {
				resp := tx.Model(entities.Multipagoscierrelotedetalles{}).Where("id = ?", detalle.ID).UpdateColumns(map[string]interface{}{"match": detalle.Match, "enobservacion": detalle.Enobservacion})
				if resp.Error != nil {
					logs.Info(resp.Error)
					erro = errors.New("error: al actualizar tabla de cierre de lote rapipago detalle")
					return erro
				}
			}
		}
		return nil
	})

	return
}

func (r *repository) CalcularRetencionesByTransferenciasSinAgruparRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {

	subquery2 := r.SQLClient.Table("transferencias").Select("pagointentos.id").
		Joins("JOIN movimientos ON transferencias.movimientos_id = movimientos.id").
		Joins("JOIN pagointentos ON movimientos.pagointentos_id = pagointentos.id").
		Joins("JOIN cuentas ON cuentas.id = movimientos.cuentas_id").
		Joins("JOIN clientes ON clientes.id = cuentas.clientes_id").
		Where("transferencias.fecha_operacion BETWEEN ? AND ?", request.FechaInicio, request.FechaFin).
		Where("cuentas.clientes_id", request.ClienteId)

	subquery1 := r.SQLClient.Table("movimientos").Select("movimientos.id").
		Where("pagointentos_id IN (?)", subquery2).
		Where("movimientos.tipo = 'C'")

	resp := r.SQLClient.Table("movimiento_retencions AS MR").
		Select("MR.cliente_id, MR.importe_retenido AS total_retencion, MR.monto AS total_monto, G.gravamen, MR.efectuada as efectuada, R.codigo_regimen as codigo_regimen, R.monto_minimo as minimo").
		Joins("JOIN retencions AS R ON MR.retencion_id = R.id").
		Joins("JOIN condicions AS C ON R.condicions_id = C.id").
		Joins("JOIN gravamens AS G ON C.gravamens_id = G.id").
		Where("MR.movimiento_id IN (?)", subquery1).
		Where("MR.deleted_at IS NULL")

	resp.Scan(&resultado)

	// manejo del error en la consulta
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CALCULAR_CLIENTE_RETENCIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetCalcularRetencionesRepository",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. GetCalcularRetencionesRepository: %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}
