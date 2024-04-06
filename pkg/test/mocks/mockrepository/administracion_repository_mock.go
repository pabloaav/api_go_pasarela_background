package mockrepository

import (
	"context"
	"time"

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
	"github.com/stretchr/testify/mock"
)

type MockRepositoryAdministracion struct {
	mock.Mock
}

func (mock *MockRepositoryAdministracion) BeginTx()    {}
func (mock *MockRepositoryAdministracion) CommitTx()   {}
func (mock *MockRepositoryAdministracion) RollbackTx() {}

func (mock *MockRepositoryAdministracion) PagoById(pagoID int64) (*entities.Pago, error) {
	args := mock.Called(pagoID)
	result := args.Get(0)
	return result.(*entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CuentaByClientePage(cliente int64, limit, offset int) (*[]entities.Cuenta, int64, error) {
	args := mock.Called(cliente, limit, offset)
	result := args.Get(0)
	return result.(*[]entities.Cuenta), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) CuentaByID(cuenta int64) (*entities.Cuenta, error) {
	args := mock.Called(cuenta)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}

func (mock *MockRepositoryAdministracion) SaveCuenta(ctx context.Context, cuenta *entities.Cuenta) (bool, error) {
	args := mock.Called(ctx, cuenta)
	return args.Bool(0), args.Error(1)
}

func (mock *MockRepositoryAdministracion) SavePagotipo(tipo *entities.Pagotipo) (bool, error) {
	args := mock.Called(tipo)
	return args.Bool(0), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCuentasByCliente(clienteId uint64) (cuentas []entities.Cuenta, erro error) {
	args := mock.Called(clienteId)
	result := args.Get(0)
	return result.([]entities.Cuenta), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosByUUID(uuid []string) (pagos []*entities.Pago, erro error) {
	args := mock.Called(uuid)
	result := args.Get(0)
	return result.([]*entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagos(filtro filtros.PagoFiltro) (pagos []entities.Pago, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pago), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetPagosIntentos(filtro filtros.PagoIntentoFiltro) (pagos []entities.Pagointento, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagointento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosEstados(filtro filtros.PagoEstadoFiltro) (estados []entities.Pagoestado, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagoestado), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosEstadosExternos(filtro filtros.PagoEstadoExternoFiltro) (estados []entities.Pagoestadoexterno, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagoestadoexterno), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estados entities.Pagoestado, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Pagoestado), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetSaldoCuenta(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error) {
	args := mock.Called(cuentaId)
	result := args.Get(0)
	return result.(administraciondtos.SaldoCuentaResponse), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetSaldoCliente(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error) {
	args := mock.Called(clienteId)
	result := args.Get(0)
	return result.(administraciondtos.SaldoClienteResponse), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetMovimientos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Movimiento), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) CreateMovimientosTransferencia(ctx context.Context, movimiento []*entities.Movimiento) error {
	args := mock.Called(ctx, movimiento)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateTransferencias(ctx context.Context, transferencias []*entities.Transferencia) (erro error) {
	args := mock.Called(ctx, transferencias)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateTransferenciasComisiones(ctx context.Context, transferencias []*entities.Transferenciacomisiones) (erro error) {
	args := mock.Called(ctx, transferencias)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) BajaMovimiento(ctx context.Context, movimientos []*entities.Movimiento, motivoBaja string) error {
	args := mock.Called(ctx, movimientos, motivoBaja)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateCierreLoteApiLink(cierreLotes []*entities.Apilinkcierrelote) (erro error) {
	args := mock.Called(cierreLotes)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) UpdateEstadoPagos(pagos []entities.Pago, pagoEstadoId uint64) (erro error) {
	args := mock.Called(pagos, pagoEstadoId)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetPlanCuotasByMedioPago(idMedioPago uint) (planCuotas []administraciondtos.PlanCuotasResponseDetalle, erro error) {
	args := mock.Called(idMedioPago)
	result := args.Get(0)
	return result.([]administraciondtos.PlanCuotasResponseDetalle), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPrismaCierreLotes(reversion bool) (prismaCierreLotes []entities.Prismacierrelote, erro error) {
	args := mock.Called(reversion)
	result := args.Get(0)
	return result.([]entities.Prismacierrelote), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPrismaPagoIntentos(siteTransaccionId string) (pagos entities.Pago, erro error) {
	args := mock.Called(siteTransaccionId)
	result := args.Get(0)
	return result.(entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CreateMovimientosCierreLote(ctx context.Context, movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse) (erro error) {

	args := mock.Called(ctx, movimientoCierreLote)

	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateCLApilinkPagosRepository(ctx context.Context, pg administraciondtos.RegistroClPagosApilink) (erro error) {
	args := mock.Called(ctx, pg)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) ActualizarPagosClRapipagoRepository(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error) {

	args := mock.Called(pagosclrapiapgo)

	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) SaveCuentacomision(comision *entities.Cuentacomision) error {
	args := mock.Called(comision)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetAllInstallmentsById(id uint) (installment []entities.Installment, erro error) {
	args := mock.Called(id)
	return args.Get(0).([]entities.Installment), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetInstallments(fechaDesde time.Time) (medioPagoInstallments []entities.Mediopagoinstallment, erro error) {
	args := mock.Called(fechaDesde)
	return args.Get(0).([]entities.Mediopagoinstallment), args.Error(1)
}

func (mock *MockRepositoryAdministracion) BuildRICuentasCliente(request ribcradtos.RICuentasClienteRequest) (ri []ribcradtos.RiCuentaCliente, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]ribcradtos.RiCuentaCliente), args.Error(1)
}

func (mock *MockRepositoryAdministracion) BuildRIDatosFondo(request ribcradtos.RiDatosFondosRequest) (ri []ribcradtos.RiDatosFondos, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]ribcradtos.RiDatosFondos), args.Error(1)
}
func (mock *MockRepositoryAdministracion) BuilRIInfestaditica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]ribcradtos.RiInfestadistica), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetConfiguraciones(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Configuracione), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) UpdateConfiguracion(ctx context.Context, request entities.Configuracione) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetClientes(filtro filtros.ClienteFiltro) (clientes []entities.Cliente, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Cliente), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetCliente(filtro filtros.ClienteFiltro) (cliente entities.Cliente, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Cliente), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCuenta(filtro filtros.CuentaFiltro) (cuenta entities.Cuenta, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Cuenta), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetTransferencias(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferencia, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Transferencia), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetTransferenciasComisiones(filtro filtros.TransferenciaFiltro) (transferencias []entities.Transferenciacomisiones, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Transferenciacomisiones), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) CreateCliente(ctx context.Context, cliente entities.Cliente) (id uint64, erro error) {
	args := mock.Called(ctx, cliente)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateCliente(ctx context.Context, cliente entities.Cliente) (erro error) {
	args := mock.Called(ctx, cliente)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) DeleteCliente(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateRubro(ctx context.Context, rubro entities.Rubro) (id uint64, erro error) {
	args := mock.Called(ctx, rubro)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateRubro(ctx context.Context, rubro entities.Rubro) (erro error) {
	args := mock.Called(ctx, rubro)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetRubro(filtro filtros.RubroFiltro) (rubro entities.Rubro, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Rubro), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetRubros(filtro filtros.RubroFiltro) (rubros []entities.Rubro, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Rubro), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) UpdateCuenta(ctx context.Context, cuenta entities.Cuenta) (erro error) {
	args := mock.Called(ctx, cuenta)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) DeleteCuenta(id uint64) (erro error) {
	args := mock.Called(id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreatePagoTipo(ctx context.Context, request entities.Pagotipo, channels []int64, cuotas []string) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdatePagoTipo(ctx context.Context, request entities.Pagotipo, channels administraciondtos.RequestPagoTipoChannels, cuotas administraciondtos.RequestPagoTipoCuotas) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetPagoTipo(filtro filtros.PagoTipoFiltro) (response entities.Pagotipo, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Pagotipo), args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetPagosTipo(filtro filtros.PagoTipoFiltro) (response []entities.Pagotipo, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagotipo), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) DeletePagoTipo(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateChannel(ctx context.Context, request entities.Channel) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateChannel(ctx context.Context, request entities.Channel) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetChannel(filtro filtros.ChannelFiltro) (channel entities.Channel, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Channel), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetChannels(filtro filtros.ChannelFiltro) (response []entities.Channel, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Channel), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) DeleteChannel(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateCuentaComision(ctx context.Context, request entities.Cuentacomision) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateCuentaComision(ctx context.Context, request entities.Cuentacomision) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetCuentaComision(filtro filtros.CuentaComisionFiltro) (response entities.Cuentacomision, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Cuentacomision), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCuentasComisiones(filtro filtros.CuentaComisionFiltro) (response []entities.Cuentacomision, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Cuentacomision), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) DeleteCuentaComision(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(ctx, id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) SetApiKey(ctx context.Context, cuenta entities.Cuenta) (erro error) {
	args := mock.Called(ctx, cuenta)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetPago(filtro filtros.PagoFiltro) (pago entities.Pago, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateTransferencias(listas bancodtos.ResponseConciliacion) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetImpuestosRepository(filtro filtros.ImpuestoFiltro) (response []entities.Impuesto, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Impuesto), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) UpdateImpuestoRepository(ctx context.Context, impuesto entities.Impuesto) (erro error) {
	args := mock.Called(ctx, impuesto)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetConsultarDebines(request linkdebin.RequestDebines) (cierreLotes []*entities.Apilinkcierrelote, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]*entities.Apilinkcierrelote), args.Error(1)
}
func (mock *MockRepositoryAdministracion) CreatePlanCuotasByInstallmenIdRepository(installmentActual, installmentNew entities.Installment, listaPlanCuotas []entities.Installmentdetail) (erro error) {
	args := mock.Called(installmentActual, installmentNew, listaPlanCuotas)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetInstallmentById(id uint) (planCuotas entities.Installment, erro error) {
	args := mock.Called(id)
	result := args.Get(0)
	return result.(entities.Installment), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdatePagosNotificados(listaPagosNotificar []uint) (erro error) {
	args := mock.Called()
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) UpdatePagosEstadoInicialNotificado(listaPagosNotificar []uint) (erro error) {
	args := mock.Called()
	return args.Error(0)
}
func (mock *MockRepositoryAdministracion) UpdateCierreloteApilink(request linkdebin.RequestListaUpdateDebines) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateImpuestoRepository(ctx context.Context, request entities.Impuesto) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetConsultarMovimientosRapipago(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Rapipagocierrelote, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]*entities.Rapipagocierrelote), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPeticionesWebServices(filtro filtros.PeticionWebServiceFiltro) (peticiones []entities.Webservicespeticione, totalFilas int64, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Webservicespeticione), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetPagosTipoChannelRepository(filtro filtros.PagoTipoChannelFiltro) (pagostipochannel []entities.Pagotipochannel, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Pagotipochannel), args.Error(1)
}

func (mock *MockRepositoryAdministracion) DeletePagoTipoChannel(id uint64) (erro error) {
	args := mock.Called()
	return args.Error(1)
}

func (mock *MockRepositoryAdministracion) CreatePagoTipoChannel(ctx context.Context, request entities.Pagotipochannel) (id uint64, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint64), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetChannelsAranceles(filtro filtros.ChannelArancelFiltro) (response []entities.Channelarancele, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Channelarancele), int64(args.Int(1)), args.Error(2)
}

func (mock *MockRepositoryAdministracion) CreateChannelsArancel(ctx context.Context, request entities.Channelarancele) (id uint64, erro error) {
	args := mock.Called(ctx, request)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateChannelsArancel(ctx context.Context, request entities.Channelarancele) (erro error) {
	args := mock.Called(ctx, request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) DeleteChannelsArancel(ctx context.Context, id uint64) (erro error) {
	args := mock.Called(id)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetChannelArancel(filtro filtros.ChannelAranceFiltro) (response entities.Channelarancele, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Channelarancele), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetMedioPagoRepository(filtro filtros.FiltroMedioPago) (mediopago entities.Mediopago, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Mediopago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCierreLoteSubidosRepository() (entityCl []entities.Prismacierrelote, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Prismacierrelote), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPrismaPxSubidosRepository() (entityPx []entities.Prismapxcuatroregistro, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Prismapxcuatroregistro), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPrismaMxSubidosRepository() (entityMx []entities.Prismamxtotalesmovimiento, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Prismamxtotalesmovimiento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerArchivoCierreLoteRapipago(nombre string) (existeArchivo bool, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(bool), args.Error(1)
}

func (mock *MockRepositoryAdministracion) UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetMovimientosNegativos(filtro filtros.MovimientoFiltro) (movimiento []entities.Movimiento, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Movimiento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, erro error) {
	args := mock.Called(apikey)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerPagosInDisputaRepository(filtro filtros.ContraCargoEnDisputa) (pagosEnDisputa []entities.Pagointento, erro error) {
	args := mock.Called(filtro)
	retsult := args.Get(0)
	return retsult.([]entities.Pagointento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerCierreLoteEnDisputaRepository(estadoDisputa int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error) {
	args := mock.Called(estadoDisputa, filtro)
	retsult := args.Get(0)
	return retsult.([]entities.Prismacierrelote), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerCierreLoteContraCargoRepository(estadoReversion int, filtro filtros.ContraCargoEnDisputa) (enttyClEnDsiputa []entities.Prismacierrelote, erro error) {
	args := mock.Called(estadoReversion, filtro)
	retsult := args.Get(0)
	return retsult.([]entities.Prismacierrelote), args.Error(1)
}

func (mock *MockRepositoryAdministracion) PostPreferencesRepository(preferenceEntity entities.Preference) (erro error) {
	args := mock.Called(preferenceEntity)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) CreateSolicitudRepository(solicitudEntity entities.Solicitud) (erro error) {
	args := mock.Called(solicitudEntity)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) UpdatePagosDev(pagos []uint) (erro error) {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetConsultarClRapipagoRepository(filtro filtros.RequestClrapipago) (movimientos []entities.Rapipagocierrelote, totalFilas int64, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.([]entities.Rapipagocierrelote), 1, args.Error(1)
}

func (mock *MockRepositoryAdministracion) ConsultarEstadoPagosRepository(parametrosVslido administraciondtos.ParamsValidados, filtro filtros.PagoFiltro) (entityPagos []entities.Pago, erro error) {
	args := mock.Called(parametrosVslido, filtro)
	resultado := args.Get(0)
	return resultado.([]entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pagointento), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CreateMovimientosTemporalesCierreLote(ctx context.Context, mcl administraciondtos.MovimientoTemporalesResponse) (erro error) {
	return
}

func (mock *MockRepositoryAdministracion) GetReportesPagoRepository(filtro filtros.PagoFiltro) (reportes entities.Reporte, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.(entities.Reporte), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetSuccessPaymentsRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pago), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetMovimientosTransferencias(request reportedtos.RequestPagosPeriodo) (movimientos []entities.Movimiento, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]entities.Movimiento), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetClienteRetencionesRepository(request retenciondtos.RentencionRequestDTO) (clienteretencion []entities.ClienteRetencion, count int64, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	resultTotal := args.Get(1)
	return result.([]entities.ClienteRetencion), resultTotal.(int64), args.Error(2)
}

func (mock *MockRepositoryAdministracion) CreateClienteRetencionRepository(request retenciondtos.RentencionRequestDTO) (erro error) {
	args := mock.Called(request)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetClienteUnlinkedRetencionesRepository(request retenciondtos.RentencionRequestDTO) (retenciones []entities.Retencion, count int64, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	resultTotal := args.Get(1)
	return result.([]entities.Retencion), resultTotal.(int64), args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetRetencionClienteRepository(request retenciondtos.RentencionRequestDTO) (clienteretencion entities.ClienteRetencion, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.(entities.ClienteRetencion), args.Error(1)
}

func (mock *MockRepositoryAdministracion) PostRetencionCertificadoRepository(certificado entities.Certificado) (erro error) {
	args := mock.Called(certificado)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetCertificadoRepository(requestId uint) (certificado entities.Certificado, erro error) {
	args := mock.Called(requestId)
	result := args.Get(0)
	return result.(entities.Certificado), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCondicionesRepository(request retenciondtos.RentencionRequestDTO) (condicions []entities.Condicion, erro error) {
	return
}

func (mock *MockRepositoryAdministracion) GetGravamensRepository(filtro retenciondtos.GravamenRequestDTO) (gravamenes []entities.Gravamen, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Gravamen), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CreateRetencionRepository(entity entities.Retencion) (entities.Retencion, error) {
	return entities.Retencion{}, nil
}

func (mock *MockRepositoryAdministracion) DeleteClienteRetencionRepository(request retenciondtos.RentencionRequestDTO) (erro error) {
	return nil
}

func (mock *MockRepositoryAdministracion) UpdateRetencionRepository(entity entities.Retencion) (entities.Retencion, error) {
	return entities.Retencion{}, nil
}

func (mock *MockRepositoryAdministracion) GetConsultarMovimientosRapipagoDetalles(filtro rapipago.RequestConsultarMovimientosRapipagoDetalles) (response []*entities.Rapipagocierrelotedetalles, erro error) {
	return
}

func (mock *MockRepositoryAdministracion) ActualizarPagosClRapipagoDetallesRepository(barcode []string) (erro error) {
	return
}

func (mock *MockRepositoryAdministracion) CreateCondicionRepository(condicion entities.Condicion) (erro error) {
	return
}

func (mock *MockRepositoryAdministracion) UpdateCondicionRepository(condicion entities.Condicion) (erro error) {
	return
}

func (mock *MockRepositoryAdministracion) EvaluarRetencionesByClienteRepository(request retenciondtos.RentencionRequestDTO) ([]retenciondtos.RetencionAgrupada, error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]retenciondtos.RetencionAgrupada), args.Error(1)
}

func (mock *MockRepositoryAdministracion) EvaluarRetencionesByMovimientosRepository(request retenciondtos.RentencionRequestDTO) ([]retenciondtos.RetencionAgrupada, error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]retenciondtos.RetencionAgrupada), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GenerarCertificacionRepository(comprobante []entities.Comprobante) error {
	args := mock.Called(comprobante)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetCalcularRetencionesRepository(request retenciondtos.RentencionRequestDTO) ([]retenciondtos.RetencionAgrupada, error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]retenciondtos.RetencionAgrupada), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetClienteRetencionRepository(retencion_id, cliente_id uint) (entities.ClienteRetencion, error) {
	args := mock.Called(retencion_id, cliente_id)
	result := args.Get(0)
	return result.(entities.ClienteRetencion), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetMovimientosRetencionesRepository(request retenciondtos.RentencionRequestDTO) ([]uint, error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]uint), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetRetencionesRepository(request retenciondtos.RentencionRequestDTO) ([]entities.Retencion, int64, error) {
	args := mock.Called(request)
	retenciones := args.Get(0).([]entities.Retencion)
	count := args.Get(1).(int64)
	return retenciones, count, args.Error(2)
}

func (mock *MockRepositoryAdministracion) GetTotalAmountByMovimientoIdsRepository(listaMovimientosId []uint) (uint64, error) {
	args := mock.Called(listaMovimientosId)
	return args.Get(0).(uint64), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ComprobantesRetencionesDevolverRepository(request retenciondtos.RentencionRequestDTO) ([]entities.Comprobante, error) {
	args := mock.Called(request)
	return args.Get(0).([]entities.Comprobante), args.Error(1)
}

func (mock *MockRepositoryAdministracion) TotalizarRetencionesMovimientosRepository(listaMovimientoIds []uint) (totalRetenciones uint64, erro error) {
	args := mock.Called(listaMovimientoIds)
	return args.Get(0).(uint64), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetComprobantesRepository(request retenciondtos.RentencionRequestDTO) (comprobantes []entities.Comprobante, erro error) {
	args := mock.Called(request)
	return args.Get(0).([]entities.Comprobante), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCertificadosVencimientoRepository(request retenciondtos.CertificadoVencimientoDTO) (certificados []entities.Certificado, erro error) {
	args := mock.Called(request)
	return args.Get(0).([]entities.Certificado), args.Error(1)
}
func (mock *MockRepositoryAdministracion) UpdateMovimientoMontoRepository(ctxPrueba context.Context, movimiento entities.Movimiento) (erro error) {
	args := mock.Called(ctxPrueba, movimiento)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetPagosTipoReferences(filtro filtros.PagoTipoReferencesFilter) (listaPagos []entities.Pagotipo, err error) {
	// args := mock.Called()
	// return args.Get(0).([]entities.Comprobante), args.Error(1)
	return
}

func (mock *MockRepositoryAdministracion) GetClientesConfiguracion(filtro filtros.ClienteConfiguracionFiltro) (clientes []entities.Cliente, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.([]entities.Cliente), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CreateAuditoria(resultado entities.Auditoria) (erro error) {
	return
}

func (mock *MockRepositoryAdministracion) ActualizarPagosClMultipagosRepository(pagosclmultipagos administraciondtos.PagosClMultipagosResponse) (erro error) {
	args := mock.Called(pagosclmultipagos)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) ActualizarPagosClMultipagosDetallesRepository(barcode []string) (erro error) {
	args := mock.Called(barcode)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) UpdateCierreLoteMultipagos(cierreLotes []*entities.Multipagoscierrelote) (erro error) {
	args := mock.Called(cierreLotes)
	return args.Error(0)
}

func (mock *MockRepositoryAdministracion) GetConsultarMovimientosMultipagos(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Multipagoscierrelote, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.([]*entities.Multipagoscierrelote), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetConsultarMovimientosMultipagosDetalles(filtro multipagos.RequestConsultarMovimientosMultipagosDetalles) (response []*entities.Multipagoscierrelotedetalles, erro error) {
	args := mock.Called(filtro)
	resultado := args.Get(0)
	return resultado.([]*entities.Multipagoscierrelotedetalles), args.Error(1)
}

func (mock *MockRepositoryAdministracion) ObtenerArchivoCierreLoteMultipagos(nombre string) (existeArchivo bool, erro error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(bool), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetCalcularRetencionesByTransferenciasRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]retenciondtos.RetencionAgrupada), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetMovimientosIdsCalculoRetencionComprobante(request retenciondtos.RentencionRequestDTO) (resultado []uint, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]uint), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosApilink(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]uint), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosRapipago(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]uint), args.Error(1)
}

func (mock *MockRepositoryAdministracion) GetPagosPrisma(filtro filtros.PagoIntentoFiltros) (ids []uint, erro error) {
	args := mock.Called(filtro)
	result := args.Get(0)
	return result.([]uint), args.Error(1)
}

func (mock *MockRepositoryAdministracion) CalcularRetencionesByTransferenciasSinAgruparRepository(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {
	args := mock.Called(request)
	result := args.Get(0)
	return result.([]retenciondtos.RetencionAgrupada), args.Error(1)
}
