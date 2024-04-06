package administracion

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos/retenciondtos"
	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos/ribcra"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/utildtos"
	webhooks "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	//PAGOS
	GetPagoByID(pagoID int64) (*entities.Pago, error)
	PostPagotipo(ctx context.Context, pagotipo *entities.Pagotipo) (bool, error)
	GetPagosService(filtro filtros.PagoFiltro) (response administraciondtos.ResponsePagos, erro error)
	GetPagosConsulta(string, administraciondtos.RequestPagosConsulta) (*[]administraciondtos.ResponsePagosConsulta, error)
	ConsultarEstadoPagosService(requestValid administraciondtos.ParamsValidados, apiKey string, request administraciondtos.RequestPagosConsulta) (responsePagoEstado []administraciondtos.ResponseSolicitudPago, registrosAfectados bool, erro error)

	///
	//GetEstadoPagoRepository()()
	///

	GetPaymentByExternalService(filtroPago filtros.PagoFiltro) (pago administraciondtos.ResponsePago, err error)

	// PAGOS INTENTOS
	GetPagosIntentosByTransaccionIdService(filtroPagoIntento filtros.PagoIntentoFiltro) (pagosIntentos []entities.Pagointento, erro error)

	//TRANSFERENCIAS
	GetTransferencias(filtro filtros.TransferenciaFiltro) (response administraciondtos.TransferenciaRespons, erro error)
	//BuildTransferenciaCliente(ctx context.Context, requerimientoId string, request administraciondtos.RequestTransferenicaCliente, cuentaId uint64) (response linktransferencia.ResponseTransferenciaCreateLink, erro error)
	/*conciliacion banco actualizar campos match con servicio banco */
	UpdateTransferencias(listas bancodtos.ResponseConciliacion) error
	SendTransferenciasComisiones(ctx context.Context, requerimientoId string, req administraciondtos.RequestComisiones) (res administraciondtos.ResponseTransferenciaComisiones, erro error)

	//CLIENTES
	CreateClienteService(ctx context.Context, request administraciondtos.ClienteRequest) (id uint64, erro error)
	UpdateClienteService(ctx context.Context, cliente administraciondtos.ClienteRequest) (erro error)
	DeleteClienteService(ctx context.Context, id uint64) (erro error)
	GetClienteService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacion, erro error)
	GetClientesService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacionPaginado, erro error)
	GetClientesConfiguracionService(filtro filtros.ClienteConfiguracionFiltro) (response administraciondtos.ResponseClientesConfiguracion, erro error)
	ObtenerClientesSinDTOService(filtro filtros.ClienteFiltro) (response []entities.Cliente, total int64, erro error)

	//RUBROS
	CreateRubroService(ctx context.Context, request administraciondtos.RubroRequest) (id uint64, erro error)
	UpdateRubroService(ctx context.Context, request administraciondtos.RubroRequest) (erro error)
	GetRubroService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubro, erro error)
	GetRubrosService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubros, erro error)

	//ABM PAGOS TIPOS
	CreatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (id uint64, erro error)
	UpdatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (erro error)
	GetPagoTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagoTipo, erro error)
	GetPagosTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagosTipo, erro error)
	DeletePagoTipoService(ctx context.Context, id uint64) (erro error)

	//ABM CHANNELS
	CreateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (id uint64, erro error)
	UpdateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (erro error)
	GetChannelService(filtro filtros.ChannelFiltro) (channel administraciondtos.ResponseChannel, erro error)
	GetChannelsService(filtro filtros.ChannelFiltro) (response administraciondtos.ResponseChannels, erro error)
	DeleteChannelService(ctx context.Context, id uint64) (erro error)

	//ABM CUENTA COMISSIONES
	CreateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (id uint64, erro error)
	UpdateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (erro error)
	GetCuentaComisionService(filtro filtros.CuentaComisionFiltro) (channel administraciondtos.ResponseCuentaComision, erro error)
	GetCuentasComisionService(filtro filtros.CuentaComisionFiltro) (response administraciondtos.ResponseCuentasComision, erro error)
	DeleteCuentaComisionService(ctx context.Context, id uint64) (erro error)

	//CUENTAS
	PostCuentaComision(ctx context.Context, comision *entities.Cuentacomision) error
	PostCuenta(ctx context.Context, cuenta administraciondtos.CuentaRequest) (bool, error)
	GetCuenta(filtro filtros.CuentaFiltro) (response administraciondtos.ResponseCuenta, erro error)
	GetCuentasByCliente(cliente int64, number, size int) (*dtos.Meta, *dtos.Links, *[]entities.Cuenta, error)
	UpdateCuentaService(ctx context.Context, request administraciondtos.CuentaRequest) (erro error)
	SetApiKeyService(ctx context.Context, request *administraciondtos.CuentaRequest) (erro error)
	DeleteCuentaService(ctx context.Context, id uint64) (erro error)
	GetCuentaByApiKeyService(apikey string) (reult bool, erro error)

	/* impuestos */
	PostImpuestoService(ctx context.Context, filtro administraciondtos.ImpuestoRequest) (id uint64, erro error)
	GetImpuestosService(filtro filtros.ImpuestoFiltro) (response administraciondtos.ResponseImpuestos, erro error)
	UpdateImpuestoService(ctx context.Context, filtro administraciondtos.ImpuestoRequest) (erro error)

	// CHANNELS ARANCELES
	GetChannelsArancelService(filtro filtros.ChannelArancelFiltro) (response administraciondtos.ResponseChannelsArancel, erro error)
	CreateChannelsArancelService(ctx context.Context, request administraciondtos.RequestChannelsAranncel) (id uint64, erro error)
	UpdateChannelsArancelService(ctx context.Context, request administraciondtos.RequestChannelsAranncel) (erro error)
	DeleteChannelsArancelService(ctx context.Context, id uint64) (erro error)
	GetChannelArancelService(filtro filtros.ChannelAranceFiltro) (response administraciondtos.ResponseChannelsAranceles, erro error)

	/*
		Devuelve el saldo de una cuenta específica
	*/
	GetSaldoCuentaService(cuentaId uint64) (saldo administraciondtos.SaldoCuentaResponse, erro error)

	/*
		Devuelve el saldo de un cliente específico
	*/
	GetSaldoClienteService(clienteId uint64) (saldo administraciondtos.SaldoClienteResponse, erro error)

	//MOVIMIENTOS
	GetMovimientosAcumulados(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoAcumuladoResponsePaginado, erro error)
	GetMovimientos(filtro filtros.MovimientoFiltro) (movimientoResponse administraciondtos.MovimientoPorCuentaResponsePaginado, erro error)
	// Se encarga de instanciar movimientos en base a los debines que se pasa por parametro, estos movimientos solo se generan si se cumple con ciertos parametros
	BuildMovimientoApiLink(listaCierre []*entities.Apilinkcierrelote) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	CreateMovimientosService(ctx context.Context, mcl administraciondtos.MovimientoCierreLoteResponse) (erro error)
	BuildPrismaMovimiento(reversion bool) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	CreateMovimientosTemporalesService(ctx context.Context, mcl administraciondtos.MovimientoTemporalesResponse) (erro error)
	// GetMovimientosById(id uint64) (movimiento entities.Movimiento, erro error)

	/*
		Crea una notificación para los usuarios del sistema
	*/
	CreateNotificacionService(notificacion entities.Notificacione) error

	CreateLogService(log entities.Log) error

	/*
		Busca una lista de pagos estados. Si busca por final es true verifica si el estado el final.
	*/
	GetPagosEstadosService(buscarPorFinal, final bool) (estados []entities.Pagoestado, erro error)
	// servicio para consultar estadopago
	GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estadoPago []entities.Pagoestado, erro error)
	/*
		se obtiene una lista de pagos estado externos
	*/
	GetPagosEstadosExternoService(filtro filtros.PagoEstadoExternoFiltro) (estadosExternos []entities.Pagoestadoexterno, erro error)
	/*
		Construye y guarda una lista de cierre de lotes para api link
		Este proceso crea el cierre de lote a partir de las informaciones consultadas en apilink
	*/
	BuildCierreLoteApiLinkService() (response administraciondtos.RegistroClPagosApilink, erro error)

	/*
		Obtiene los pagointentos que no tienen su "apilinkcierrelote" correspondiente, crea estos cierrelotes en base a datos traidos de APILINK y los devuelve para guardar y actualizar datos en la DB
	*/
	BuildDebinNotRegisteredApiLinkService(request cierrelotedtos.RequestDebinNotRegisteredApilink) (res administraciondtos.RegistroClPagosApilink, erro error)

	/*
		permite obtener los planes de cuotas vigentes para un medio de pago
	*/
	GetPlanCuotas(idMedioPago uint) (response []administraciondtos.PlanCuotasResponseDetalle, erro error)

	/*
		obtiene los intereses de todos los planes existentes para informarlos
	*/
	GetInteresesPlanes(fecha string) (planes []administraciondtos.PlanCuotasResponse, erro error)
	/*
		obtiene todos los planes decuotas por id de installment
	*/
	GetAllInstallmentsById(installments_id int64) (planesCuotas []administraciondtos.InstallmentsResponse, erro error)
	//------------------------------------------------RI BCRA-----------------------------------------
	/*
		Crea archivo txt para enviar al BCRA
		https://telcodev.atlassian.net/secure/RapidBoard.jspa?rapidView=17&projectKey=PP&modal=detail&selectedIssue=PP-122
	*/
	RIInfestadistica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error)
	/*
		Construye la información para supervisión prevista en la sección 69.1 (presentación de informaciones al banco central).
		https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/blob/development/document/administracion/bcraregimeninformativo/bcra_regimen_informativo.md
	*/
	GetInformacionSupervision(request ribcradtos.GetInformacionSupervisionRequest) (ri ribcradtos.RiInformacionSupervisionReponse, erro error)
	/*
		Guarda en un archivo zip la información de supervisión prevista en la sección 69.1 (presentación de informaciones al banco central).
		https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/blob/development/document/administracion/bcraregimeninformativo/bcra_regimen_informativo.md
	*/
	BuildInformacionSupervision(request ribcradtos.BuildInformacionSupervisionRequest) (ruta string, erro error)
	/*
		Construye la información para estadística prevista en la sección 69.2 (presentación de informaciones al banco central).
		https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/blob/development/document/administracion/bcraregimeninformativo/bcra_regimen_informativo.md
	*/
	GetInformacionEstadistica(request ribcradtos.GetInformacionEstadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error)
	/*
		Guarda en un archivo zip la información de estadistica prevista en la sección 69.2 (presentación de informaciones al banco central).
		https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/blob/development/document/administracion/bcraregimeninformativo/bcra_regimen_informativo.md
	*/
	BuildInformacionEstadistica(request ribcradtos.BuildInformacionEstadisticaRequest) (ruta string, erro error)

	RIGuardarArchivos(request ribcradtos.RIGuardarArchivosRequest) (erro error)
	//------------------------------------------------RI BCRA-----------------------------------------

	/*
		Modifica el estado de pagos expirados
	*/
	ModificarEstadoPagosExpirados() (erro error)

	/*
		Realiza automaticamente las transferencias de acuerdo con el período informado en configuraciones
	*/
	RetiroAutomaticoClientes(ctx context.Context, request administraciondtos.TransferenciasClienteId) (response administraciondtos.RequestMovimientosId, erro error)
	RetiroAutomaticoClientesSubcuentas(ctx context.Context) (response administraciondtos.RequestMovimientosId, erro error)

	//CONFIGURACIONES
	GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (response administraciondtos.ResponseConfiguraciones, erro error)
	UpdateConfiguracionService(ctx context.Context, config administraciondtos.RequestConfiguracion) (erro error)
	UpdateConfiguracionSendEmailService(ctx context.Context, request administraciondtos.RequestConfiguracion) (erro error)

	//Send Mails
	SendSolicitudCuenta(request administraciondtos.SolicitudCuentaRequest) (erro error)

	// plan de cuotas INSTALLMENTS
	CreatePlanCuotasService(request administraciondtos.RequestPlanCuotas) (erro error)
	// notificaciones de pagos
	BuildNotificacionPagosService(request webhooks.RequestWebhook) (listaPagos []entities.Pagotipo, erro error)
	BuildNotificacionPagosWithReferences(request webhooks.RequestWebhookReferences) ([]entities.Pagotipo, error)
	BuildNotificacionPagosCLRapipago(filtro filtros.PagoEstadoFiltro) (response []webhooks.WebhookResponse, barcode []string, erro error)
	CreateNotificacionPagosService(listaPagos []entities.Pagotipo) (response []webhooks.WebhookResponse, erro error)
	NotificarPagos(listaPagosNotificar []webhooks.WebhookResponse) (pagoupdate []uint)
	UpdatePagosNoticados(listaPagosNotificar []uint) (erro error)

	// & apilink cierrelote
	CreateCLApilinkPagosService(ctx context.Context, mcl administraciondtos.RegistroClPagosApilink) (erro error)
	CreateCierreLoteApiLink(cierreLotes []*entities.Apilinkcierrelote) (erro error)
	GetDebines(request linkdebin.RequestDebines) (response []*entities.Apilinkcierrelote, erro error)
	GetConsultarDebines(request linkdebin.RequestDebines) (response []linkdebin.ResponseDebinesEliminados, erro error)
	BuildNotificacionPagosCLApilink(request []linkdebin.ResponseDebinesEliminados) (response []webhooks.WebhookResponse, debinID []uint64, erro error)
	UpdateCierreLoteApilink(request linkdebin.RequestListaUpdateDebines) (erro error)

	// & rapipagocierrelote
	GetCierreLoteRapipagoService(filtro rapipago.RequestConsultarMovimientosRapipago) (listaCierreRapipago []*entities.Rapipagocierrelote, erro error)
	UpdateCierreLoteRapipago(cierreLotes []*entities.Rapipagocierrelote) (erro error)
	BuildPagosClRapipago(listaPagosClRapipago []*entities.Rapipagocierrelote) (pagosclrapiapgo administraciondtos.PagosClRapipagoResponse, erro error)
	BuildRapipagoMovimiento(listaCierre []*entities.Rapipagocierrelote) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	ActualizarPagosClRapipagoService(pagosclrapiapgo administraciondtos.PagosClRapipagoResponse) (erro error)
	ActualizarPagosClRapipagoDetallesService(barcode []string) (erro error)

	// & multipagocierrelote
	GetCierreLoteMultipagosService(filtro rapipago.RequestConsultarMovimientosRapipago) (response []*entities.Multipagoscierrelote, erro error)
	UpdateCierreLoteMultipagos(cierreLotes []*entities.Multipagoscierrelote) (erro error)
	BuildPagosClMultipagos(listaPagosClMultipagos []*entities.Multipagoscierrelote) (pagosclmultipagos administraciondtos.PagosClMultipagosResponse, erro error)
	BuildMultipagosMovimiento(listaCierre []*entities.Multipagoscierrelote) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	ActualizarPagosClMultipagosService(pagosclmultipagos administraciondtos.PagosClMultipagosResponse) (erro error)
	ActualizarPagosClMultipagosDetallesService(barcode []string) (erro error)

	BuildNotificacionPagosCLMultipagos(filtro filtros.PagoEstadoFiltro) (response []webhooks.WebhookResponse, barcode []string, erro error)

	//PagoTipoChannel
	GetPagosTipoChannelService(filtro filtros.PagoTipoChannelFiltro) (response []entities.Pagotipochannel, erro error)
	DeletePagoTipoChannelService(ctx context.Context, id uint64) (erro error)
	CreatePagoTipoChannel(ctx context.Context, request administraciondtos.RequestPagoTipoChannel) (id uint64, erro error)

	//Busca lista de peticiones web services
	GetPeticionesService(filtro filtros.PeticionWebServiceFiltro) (peticiones administraciondtos.ResponsePeticionesWebServices, erro error)

	//SubirArchivos
	SubirArchivos(ctx context.Context, rutaArchivos string, listaArchivo []administraciondtos.ArchivoResponse) (countArchivo int, erro error)
	SubirArchivosCloud(ctx context.Context, rutaArchivos string, listaArchivo []administraciondtos.ArchivoResponse, directorio string) (countArchivo int, erro error)

	// Archivos Subidos
	ObtenerArchivosSubidos(filtro filtros.Paginacion) (lisArchivosSubidos administraciondtos.ResponseArchivoSubido, erro error)

	// obtener archivo de cierreloterapipago
	ObtenerArchivoCierreLoteRapipago(nombre string) (result bool, err error)

	// obtener archivo de cierrelotemultipagos
	ObtenerArchivoCierreLoteMultipagos(nombre string) (result bool, err error)

	// reversiones o contracargo //
	// obtener cierre lote con reversiones en disputa
	GetCierreLoteEnDisputaServices(estadoDisputa int, request filtros.ContraCargoEnDisputa) (cierreLoteDisputa []cierrelotedtos.ResponsePrismaCL, erro error)
	// obtener informacion de los pagos relacionados con los cierre de lotes en disputa
	GetPagosByTransactionIdsServices(filtro filtros.ContraCargoEnDisputa, cierreLoteDisputa []cierrelotedtos.ResponsePrismaCL) (listaRevertidos administraciondtos.ResponseOperacionesContracargo, erro error)

	/* Preferences*/
	PostPreferencesService(request administraciondtos.RequestPreferences) (erro error)

	/*Obtener pagos para pruebas -> utilizada para generar movimientos en dev*/
	GetPagosDevService(filtro filtros.PagoFiltro) (response []entities.Pago, erro error)
	UpdatePagosDevService(pagos []entities.Pago) (pg []uint, erro error)

	BuildPagosMovDev(pagos []uint) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error)
	//? Servicios que permiten consultar datos de cllote pora herramieta wee
	GetConsultarClRapipagoService(filtro filtros.RequestClrapipago) (response administraciondtos.ResponseCLRapipago, erro error)

	GetCaducarOfflineIntentos() (intentosCaducados int, erro error)

	// servicio para consultar contruir y generar movimientos temporales
	GetPagosCalculoMovTemporalesService(filtro filtros.PagoIntentoFiltros) (pagosid []uint, erro error)
	BuildPagosCalculoTemporales(pagos []uint) (movimientoCierreLote administraciondtos.MovimientoTemporalesResponse, erro error)
	GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error)

	ConciliacionPagosReportesService(filtro filtros.PagoFiltro) (valoresNoEncontrados []string, erro error)

	// RETENCIONES IMPOSITIVAS
	// retorna todas las retenciones
	GetRetencionesService(request retenciondtos.RentencionRequestDTO, getDTO bool) (response retenciondtos.RentencionesResponseDTO, erro error)
	// retorna las retenciones asociadas a un cliente segun id del cliente.
	GetClienteRetencionesService(request retenciondtos.RentencionRequestDTO) (response retenciondtos.RentencionesResponseDTO, erro error)
	// retorn las retenciones no vinculadas a un cliente
	GetClienteUnlinkedRetencionesService(request retenciondtos.RentencionRequestDTO) (response retenciondtos.RentencionesResponseDTO, erro error)
	// crea una asociacion cliente y retencion. Retorna un cliente con una nueva retencion asociada.
	CreateClienteRetencionService(request retenciondtos.RentencionRequestDTO) (erro error)
	// calcula las retenciones por cliente en un periodo
	GetCalcularRetencionesService(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)
	// crea las retenciones por cada movimiento generado segun su monto y tipo de retencion
	BuildRetenciones(movimiento *entities.Movimiento, importe entities.Monto, pagointento entities.Pagointento, clientes []entities.Cliente) (erro error)
	BuildRetencionesTemporales(movimiento *entities.Movimientotemporale, importe entities.Monto, pagointento entities.Pagointento, clientes []entities.Cliente) (erro error)
	// Obtener las condiones frente a los gravamenes
	GetCondicionesService(request retenciondtos.RentencionRequestDTO) (response []retenciondtos.CondicionResponseDTO, erro error)
	// crear una nueva retencion
	CreateRetencionService(request retenciondtos.PostRentencionRequestDTO, isUpdate bool) (response []entities.Retencion, erro error)
	// Modificar una retencion
	UpdateRetencionService(request retenciondtos.PostRentencionRequestDTO, isUpdate bool) (response entities.Retencion, erro error)
	// eliminar la relacion entre un cliente y una retencion
	DeleteClienteRetencionService(request retenciondtos.RentencionRequestDTO) (erro error)
	// crear una condicion impositiva
	UpSertCondicionService(request retenciondtos.CondicionRequestDTO, isUpdate bool) (erro error)
	// obtener todos los gravamenes
	GetGravamenesService(filtro retenciondtos.GravamenRequestDTO) (response []retenciondtos.GravamenResponseDTO, erro error)
	// Comprobar si las Retenciones exeden un monto minimo en caso de que la retencion evaluada lo tenga
	ComprobarMinimoRetencion(retencion entities.Retencion, clienteId uint) (result bool, monto entities.Monto, erro error)
	// Comprueba si un cliente tiene asignada una retencion
	GetClienteRetencionService(retencion_id uint, cliente_id uint) (cliente_retencion entities.ClienteRetencion, erro error)
	// dado un cliente determina si hay que realizar retenciones en un periodo especifico
	EvaluarRetencionesByClienteService(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)
	// retorna total de retenciones agrupado por gravamen, segun una lista de movimientos recibida, para un cliente determinado
	EvaluarRetencionesByMovimientoService(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)
	// registros de retenciones calculadas para un cliente en un periodo determinado
	GenerarCertificacionService(request retenciondtos.RentencionRequestDTO) (erro error)
	// obtener una lista de movimientos id de la tabla movimiento_retencions segun filtro recibido por parametro
	GetMovimientosRetencionesService(request retenciondtos.RentencionRequestDTO) (listaMovimientosId []uint, erro error)
	// retenciones practicadas no aplicables
	ComprobantesRetencionesDevolverService(request retenciondtos.RentencionRequestDTO) (comprobantesdto retenciondtos.DevolverRetencionesDTO, erro error)
	// sumar las retenciones de movimientos. @params: listaMovimientoIds []uint
	TotalizarRetencionesMovimientosService(listaMovimientoIds []uint) (totalRetenciones entities.Monto, erro error)
	// Obtener informacion de comprobantes de retenciones de un cliente con sus detalles
	GetComprobantesService(request retenciondtos.RentencionRequestDTO) (comprobantesdto []retenciondtos.ComprobanteResponseDTO, erro error)

	/* 	Retenciones Certificados */
	PostRetencionesCertificadosService(request retenciondtos.RetencionCertificadoRequestDTO) (erro error)
	ValidarRetencion(request retenciondtos.RetencionCertificadoRequestDTO) (clienteId uint, cliente_name string, erro error)
	GetCertificadoService(certificadoId uint) (certificado entities.Certificado, erro error)
	GetCertificadoCloudService(ctx context.Context, nombreFile string) (err error)
	LeerContenidoDirectorio(datos entities.Certificado) (file retenciondtos.CertificadoFileDTO, erro error)

	/* Vencimiento Certificados */
	NotificarVencimientoCertificadosService(request retenciondtos.CertificadoVencimientoDTO) (erro error)

	//Generar comisiones manual
	CreateComisionManualService(request administraciondtos.RequestComisionManual) (err error)
	//CreateAuditoriaService(request administraciondtos.RequestAuditoria) (err error)

	//Notificar pagos webhook sin notificar
	//NOTE Se eliminara cuando se aplique rabbit
	NotificarPagosWebhookSinNotificarService() (err error)

	GetMovimientosIdsCalculoRetencionComprobanteService(request retenciondtos.RentencionRequestDTO) (resultado []uint, erro error)

	// obtener contenido de un comprobante de retencion recuperado desde el cloud storage
	LeerContenidoComprobanteRetencion(datos retenciondtos.ComprobanteResponseDTO) (file retenciondtos.ComprobanteFileDTO, erro error)
	// obtener contenido de un archivo de reporte de rendicion mensual recuperado desde el cloud storage
	LeerContenidoReporteRendicionMensual(request retenciondtos.RentencionRequestDTO) (file retenciondtos.ComprobanteFileDTO, erro error)

	CalcularRetencionesByTransferenciasSinAgruparService(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error)

	GetDebinService(request linkdebin.RequestGetDebinLink) (response linkdebin.ResponseGetDebinLink, err error)
}

// variable que va a manejar la instancia del servicio
var admService *service

type service struct {
	repository     Repository
	apilinkService apilink.AplinkService
	commonsService commons.Commons
	utilService    util.UtilService
	webhook        webhook.RemoteRepository
	store          util.Store
}

func NewService(r Repository, s apilink.AplinkService, c commons.Commons, u util.UtilService, webhook webhook.RemoteRepository, storage util.Store) Service {
	admService = &service{
		repository:     r,
		apilinkService: s,
		commonsService: c,
		utilService:    u,
		webhook:        webhook,
		store:          storage,
	}
	return admService
}

// Resolve devuelve la instancia antes creada
func Resolve() *service {
	return admService
}

func _setPaginacion(number uint32, size uint32, total int64) (meta dtos.Meta) {
	from := (number - 1) * size
	lastPage := math.Ceil(float64(total) / float64(size))

	meta = dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}

	return

}

// func (s *service) GetMovimientosById(id uint64) (movimiento entities.Movimiento, erro error) {
// 	movimiento, erro = s.repository.GetMovimientosById(id)
// 	if erro != nil {
// 		return
// 	}

//		return
//	}

func (s *service) GetPaymentByExternalService(filtroPago filtros.PagoFiltro) (pago administraciondtos.ResponsePago, err error) {

	pagoEntity, err := s.repository.GetPaymentByExternal(filtroPago)
	if err != nil {
		return
	}

	var ultimoPagoIntento entities.Pagointento
	for _, pagointento := range pagoEntity.PagoIntentos {
		if len(pagointento.ExternalID) > 10 {
			ultimoPagoIntento = pagointento
		}
	}

	pago.Identificador = pagoEntity.ID
	pago.ExternalReference = pagoEntity.ExternalReference
	pago.PayerName = pagoEntity.PayerName
	pago.Estado = pagoEntity.PagoEstados.Nombre
	pago.PagoIntento.ExternalID = ultimoPagoIntento.ExternalID
	pago.PagoIntento.ID = uint64(ultimoPagoIntento.ID)
	pago.Amount = ultimoPagoIntento.Amount
	pago.FechaPago = ultimoPagoIntento.PaidAt
	pago.UltimoPagoIntentoId = uint64(ultimoPagoIntento.ID)

	return
}
func (s *service) GetDebinService(request linkdebin.RequestGetDebinLink) (response linkdebin.ResponseGetDebinLink, err error) {
	uuid := s.commonsService.NewUUID()
	response, erro := s.apilinkService.GetDebinApiLinkService(uuid, request)
	if erro != nil {
		return
	}
	return
}
func (s *service) GetCuenta(filtro filtros.CuentaFiltro) (response administraciondtos.ResponseCuenta, erro error) {

	resp, erro := s.repository.GetCuenta(filtro)
	if erro != nil {
		return
	}

	response.FromCuenta(resp)

	return
}

func (s *service) GetPagoByID(pagoID int64) (*entities.Pago, error) {
	return s.repository.PagoById(pagoID)
}

func (s *service) GetPagosConsulta(apikey string, req administraciondtos.RequestPagosConsulta) (*[]administraciondtos.ResponsePagosConsulta, error) {
	err := req.IsValid()
	if err != nil {
		return nil, fmt.Errorf("validación %w", err)
	}
	var pagotiposIds []uint64
	var rangoFechas []string
	//obtengo los pagostipos relacionada con
	cuenta, err := s.repository.GetCuentaByApiKey(apikey)
	if err != nil {
		return nil, errors.New("error: " + err.Error())
	}
	for _, values := range *cuenta.Pagotipos {
		pagotiposIds = append(pagotiposIds, uint64(values.ID))
	}

	if len(req.FechaDesde) > 0 {
		fechaDesde, err := time.Parse("02-01-2006", req.FechaDesde)
		if err != nil {
			return nil, fmt.Errorf("formato de fecha desde incorrecto: %w", err)
		}
		fechaHasta, err := time.Parse("02-01-2006", req.FechaHasta)
		if err != nil {
			return nil, fmt.Errorf("formato de fecha hasta incorrecto: %w", err)
		}
		if fechaHasta.Sub(fechaDesde) < 0 {
			return nil, fmt.Errorf("periodo de consulta incorrecto")
		}
		if fechaHasta.Sub(fechaDesde).Hours()/24 > 7 {
			return nil, fmt.Errorf("período de consulta mayor a 7 días")
		}
		rangoFechas = append(rangoFechas, fechaDesde.Format("2006-01-02")+" 00:00:00", fechaHasta.Format("2006-01-02")+" 23:59:59")
		// cuenta, err := s.repository.GetCuentaByApiKey(apikey)
		// if err != nil {
		// 	return nil, errors.New("error: " + err.Error())
		// }
		// for _, values := range *cuenta.Pagotipos {
		// 	pagotiposIds = append(pagotiposIds, uint64(values.ID))
		// }

	}

	var uuidList []string

	if len(req.Uuid) > 0 {
		uuidList = append(uuidList, req.Uuid)
	}
	if len(req.Uuids) > 0 {
		uuidList = append(uuidList, req.Uuids...)
	}

	filtro := filtros.PagoFiltro{
		Uuids:                uuidList,
		ExternalReference:    req.ExternalReference,
		CargarPagoEstado:     true,
		Fecha:                rangoFechas,
		Notificado:           true,
		PagosTipoIds:         pagotiposIds,
		VisualizarPendientes: true,
	}

	pago, _, err := s.repository.GetPagos(filtro)
	if err != nil {
		return nil, fmt.Errorf("consultando a la base de datos: %w", err)
	}

	res := make([]administraciondtos.ResponsePagosConsulta, len(pago))

	for i, p := range pago {
		res[i].SetPago(p)
	}

	return &res, nil
}

func (s *service) ConsultarEstadoPagosService(requestValid administraciondtos.ParamsValidados, apiKey string, request administraciondtos.RequestPagosConsulta) (responsePagoEstado []administraciondtos.ResponseSolicitudPago, registrosAfectados bool, erro error) {
	err := request.IsValid()
	if err != nil {
		erro = fmt.Errorf("validación %w", err)
		return
	}
	var pagotiposIds []uint64
	var rangoFechas []string
	var uuidList []string
	var referenciaExterna string
	//obtengo los pagostipos relacionada con
	cuenta, err := s.repository.GetCuentaByApiKey(apiKey)
	if err != nil {
		erro = errors.New("error: " + err.Error())
		return
	}
	for _, values := range *cuenta.Pagotipos {
		pagotiposIds = append(pagotiposIds, uint64(values.ID))
	}

	if requestValid.Uuuid {
		uuidList = append(uuidList, request.Uuid)
	}
	if requestValid.ExternalReference {
		referenciaExterna = request.ExternalReference
	}
	if requestValid.RangoFecha {
		if len(request.FechaDesde) > 0 && len(request.FechaHasta) > 0 {
			fechaDesde, err := time.Parse("02-01-2006", request.FechaDesde)
			if err != nil {
				erro = fmt.Errorf("formato de fecha desde incorrecto: %w", err)
				return
			}
			fechaHasta, err := time.Parse("02-01-2006", request.FechaHasta)
			if err != nil {
				erro = fmt.Errorf("formato de fecha hasta incorrecto: %w", err)
				return
			}
			if fechaHasta.Sub(fechaDesde) < 0 {
				erro = fmt.Errorf("periodo de consulta incorrecto")
				return
			}
			if fechaHasta.Sub(fechaDesde).Hours()/24 > 7 {
				erro = fmt.Errorf("período de consulta mayor a 7 días")
				return
			}
			rangoFechas = append(rangoFechas, fechaDesde.Format("2006-01-02")+" 00:00:00", fechaHasta.Format("2006-01-02")+" 23:59:59")
		}
	}
	if requestValid.Uuids {
		uuidList = append(uuidList, request.Uuids...)
	}

	filtroEstado := filtros.PagoEstadoFiltro{
		Nombre: "PENDING",
	}
	entityPagoEstado, err := s.repository.GetPagoEstado(filtroEstado)
	if err != nil {
		erro = errors.New("no se pudo obtener estado de pago" + err.Error())
		return
	}

	filtro := filtros.PagoFiltro{
		Uuids:             uuidList,
		ExternalReference: referenciaExterna,
		Fecha:             rangoFechas,
		PagosTipoIds:      pagotiposIds,
		PagoEstadosId:     uint64(entityPagoEstado.ID),
		CargarPagoTipos:   true,
		CargarPagoEstado:  true,
		CargaPagoIntentos: true,
		//Notificado:           true,
		//VisualizarPendientes: true,

	}
	entityPagos, err := s.repository.ConsultarEstadoPagosRepository(requestValid, filtro)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	if len(entityPagos) == 0 {
		registrosAfectados = false
		return
	}
	for _, value := range entityPagos {
		var temporalPagoEstado administraciondtos.ResponseSolicitudPago
		temporalPagoEstado.SolicitudEntityToDtos(value)
		// temporalPagoEstado.PagoIntento[0].GrossFee = s.utilService.ToFixed(temporalPagoEstado.PagoIntento[0].GrossFee, 2)
		// temporalPagoEstado.PagoIntento[0].NetFee = s.utilService.ToFixed(temporalPagoEstado.PagoIntento[0].NetFee, 2)
		// temporalPagoEstado.PagoIntento[0].FeeIva = s.utilService.ToFixed(temporalPagoEstado.PagoIntento[0].FeeIva, 2)
		responsePagoEstado = append(responsePagoEstado, temporalPagoEstado)
	}
	registrosAfectados = true
	return
}

func (s *service) GetPagosIntentosByTransaccionIdService(filtroPagoIntento filtros.PagoIntentoFiltro) (pagosIntentos []entities.Pagointento, erro error) {
	pagosIntentos, err := s.repository.GetPagosIntentos(filtroPagoIntento)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}

	return
}

func (s *service) GetCuentasByCliente(cliente int64, number, size int) (*dtos.Meta, *dtos.Links, *[]entities.Cuenta, error) {
	from := (number - 1) * size
	data, total, e := s.repository.CuentaByClientePage(cliente, size, from)
	lastPage := math.Ceil(float64(total) / float64(size))

	meta := dtos.Meta{
		Page: dtos.Page{
			CurrentPage: int32(number),
			From:        int32(from),
			LastPage:    int32(lastPage),
			PerPage:     int32(size),
			To:          int32(number * size),
			Total:       int32(total),
		},
	}

	links := dtos.Links{
		First: fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cliente, 1, size),
		Last:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cliente, meta.Page.LastPage, size),
		Next:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cliente, (number + 1), size),
		Prev:  fmt.Sprintf(config.APP_HOST+"/administracion/cuentas?cliente=%d&number=%d&size=%d", cliente, (number - 1), size),
	}

	return &meta, &links, data, e
}

func (s *service) PostCuenta(ctx context.Context, request administraciondtos.CuentaRequest) (ok bool, err error) {

	err = request.IsVAlid(false)

	if err != nil {
		return
	}
	/* 	Se comenta la validacion de CVU/CBU repetido. 28-12-2022 */
	//Valido si ya existe una cuenta con el cbu/cvu registrado
	// var filtro filtros.CuentaFiltro
	// if len(request.Cbu) > 0 {
	// 	filtro.Cbu = request.Cbu
	// 	response, err := s.repository.GetCuenta(filtro)

	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	if response.ID > 0 {
	// 		return false, fmt.Errorf(ERROR_CBU_REGISTRADO)
	// 	}

	// } else {
	// 	filtro.Cvu = request.Cvu
	// 	response, err := s.repository.GetCuenta(filtro)

	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	if response.ID > 0 {
	// 		return false, fmt.Errorf(ERROR_CVU_REGISTRADO)
	// 	}
	// }
	//Creo un uuid automaticamente
	request.Apikey = s.commonsService.NewUUID()

	cuenta := request.ToCuenta()

	ok, err = s.repository.SaveCuenta(ctx, &cuenta)
	if err != nil {
		return false, err
	}

	return ok, err
}

func (s *service) SetApiKeyService(ctx context.Context, request *administraciondtos.CuentaRequest) (erro error) {

	if request.Id < 1 {
		return fmt.Errorf(ERROR_ID)
	}

	request.Apikey = s.commonsService.NewUUID()

	cuenta := request.ToCuenta()

	erro = s.repository.SetApiKey(ctx, cuenta)

	if erro != nil {
		return erro
	}

	return
}

func (s *service) UpdateCuentaService(ctx context.Context, request administraciondtos.CuentaRequest) (erro error) {

	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}
	filtro := filtros.CuentaFiltro{
		DistintoId: request.Id,
	}

	if len(request.Cbu) > 0 {
		filtro.Cbu = request.Cbu
		response, err := s.repository.GetCuenta(filtro)

		if err != nil {
			return err
		}
		if response.ID > 0 {
			return fmt.Errorf(ERROR_CBU_REGISTRADO)
		}

	} else {
		filtro.Cvu = request.Cvu
		response, err := s.repository.GetCuenta(filtro)

		if err != nil {
			return err
		}
		if response.ID > 0 {
			return fmt.Errorf(ERROR_CVU_REGISTRADO)
		}
	}

	cuenta := request.ToCuenta()

	erro = s.repository.UpdateCuenta(ctx, cuenta)

	if erro != nil {
		return erro
	}

	return
}

func (s *service) DeleteCuentaService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	erro = s.repository.DeleteCuenta(id)
	if erro != nil {
		logs.Error(erro)
		return
	}

	if erro != nil {
		return erro
	}

	return
}

func (s *service) GetCuentaByApiKeyService(apikey string) (result bool, erro error) {
	cuenta, err := s.repository.GetCuentaByApiKey(apikey)
	if err != nil {
		log := entities.Log{
			Tipo:          "info",
			Funcionalidad: "GetCuentaByApiKey",
			Mensaje:       err.Error(),
		}
		err = s.utilService.CreateLogService(log)
		if err != nil {
			logs.Error("error al intentar registrar logs de erro en GetCuentaByApiKey")
		}
		erro = errors.New("api-key invalido")
		return
	}
	result = false
	if len(cuenta.Apikey) > 0 {
		result = true
	}
	return
}

func (s *service) GetImpuestosService(filtro filtros.ImpuestoFiltro) (response administraciondtos.ResponseImpuestos, erro error) {
	impuestosEntity, totalFilas, err := s.repository.GetImpuestosRepository(filtro)
	if err != nil {
		return
	}
	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, totalFilas)
	}
	for _, valueImpuesto := range impuestosEntity {
		impuesto := administraciondtos.ResponseImpuesto{}
		impuesto.FromImpuesto(valueImpuesto)
		response.Impuestos = append(response.Impuestos, impuesto)
	}
	return
}

func (s *service) PostImpuestoService(ctx context.Context, filtro administraciondtos.ImpuestoRequest) (id uint64, erro error) {

	erro = filtro.Validar()
	if erro != nil {
		return 0, erro
	}

	impuesto := filtro.ToImpuesto(false)

	id, erro = s.repository.CreateImpuestoRepository(ctx, impuesto)
	if erro != nil {
		return 0, erro
	}

	return
}

func (s *service) UpdateImpuestoService(ctx context.Context, request administraciondtos.ImpuestoRequest) (erro error) {

	erro = request.Validar()
	if erro != nil {
		return erro
	}

	impuesto := request.ToImpuesto(true)

	return s.repository.UpdateImpuestoRepository(ctx, impuesto)

}

func (s *service) PostPagotipo(ctx context.Context, pagotipo *entities.Pagotipo) (bool, error) {
	var err error

	res, err := s.repository.SavePagotipo(pagotipo)
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}
	return res, err
}

func (s *service) PostCuentaComision(ctx context.Context, comision *entities.Cuentacomision) error {

	err := s.repository.SaveCuentacomision(comision)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return err
}

func (s *service) CreateNotificacionService(notificacion entities.Notificacione) error {

	err := s.utilService.CreateNotificacionService(notificacion)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateLogService(log entities.Log) error {
	err := s.utilService.CreateLogService(log)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetPagosEstadosService(buscarPorFinal, final bool) (estados []entities.Pagoestado, erro error) {
	filtro := filtros.PagoEstadoFiltro{BuscarPorFinal: buscarPorFinal, Final: final}
	estados, erro = s.repository.GetPagosEstados(filtro)

	return
}

func (s *service) GetPagoEstado(filtro filtros.PagoEstadoFiltro) (estados []entities.Pagoestado, erro error) {
	estados, erro = s.repository.GetPagosEstados(filtro)

	return
}

func (s *service) GetPagosEstadosExternoService(filtro filtros.PagoEstadoExternoFiltro) (estadosExternos []entities.Pagoestadoexterno, erro error) {
	estadosExternos, erro = s.repository.GetPagosEstadosExternos(filtro)
	if erro != nil {
		return
	}
	return
}

func (s *service) GetPagosService(filtro filtros.PagoFiltro) (response administraciondtos.ResponsePagos, erro error) {

	pagos, total, erro := s.repository.GetPagos(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	var listaPagoIntentos []uint64

	// recorrer cada uno de los pagos obtenidos
	for _, p := range pagos {
		r := administraciondtos.ResponsePago{
			Identificador:     p.ID,
			Fecha:             p.CreatedAt,
			ExternalReference: p.ExternalReference,
			PayerName:         p.PayerName,
		}

		if p.PagoEstados.ID > 0 {
			r.Estado = string(p.PagoEstados.Estado)
			r.NombreEstado = p.PagoEstados.Nombre

		}

		if p.PagosTipo.ID > 0 {
			r.Pagotipo = p.PagosTipo.Pagotipo
			if p.PagosTipo.Cuenta.ID > 0 {
				r.Cuenta = p.PagosTipo.Cuenta.Cuenta
			}
		}

		if len(p.PagoIntentos) > 0 {
			r.Amount = p.PagoIntentos[len(p.PagoIntentos)-1].Amount
			r.FechaPago = p.PagoIntentos[len(p.PagoIntentos)-1].PaidAt
			if p.PagoIntentos[len(p.PagoIntentos)-1].Mediopagos.ID > 0 && p.PagoIntentos[len(p.PagoIntentos)-1].Mediopagos.Channel.ID > 0 {
				r.Channel = p.PagoIntentos[len(p.PagoIntentos)-1].Mediopagos.Channel.Channel
				r.NombreChannel = p.PagoIntentos[len(p.PagoIntentos)-1].Mediopagos.Channel.Nombre
			}
			r.UltimoPagoIntentoId = uint64(p.PagoIntentos[len(p.PagoIntentos)-1].ID)
			listaPagoIntentos = append(listaPagoIntentos, uint64(p.PagoIntentos[len(p.PagoIntentos)-1].ID))
		}

		if filtro.CargarPagosItems {
			var listaItems []administraciondtos.PagoItems
			for _, items := range p.Pagoitems {
				listaItems = append(listaItems, administraciondtos.PagoItems{
					Descripcion:   items.Description,
					Identificador: items.Identifier,
					Cantidad:      int64(items.Quantity),
					Monto:         float64(items.Amount),
				})
			}
			r.PagoItems = listaItems
		}

		response.Pagos = append(response.Pagos, r)
	}

	FiltroMovimientos := filtros.MovimientoFiltro{
		PagoIntentosIds: listaPagoIntentos,
		CuentaId:        filtro.CuentaId,
	}

	movimientos, _, erro := s.repository.GetMovimientos(FiltroMovimientos)

	if erro != nil {
		return
	}

	var listaMovimientos []uint64
	for i := range movimientos {
		listaMovimientos = append(listaMovimientos, uint64(movimientos[i].ID))
	}

	filtroTransferencias := filtros.TransferenciaFiltro{
		MovimientosIds: listaMovimientos,
	}

	transferencias, _, erro := s.repository.GetTransferencias(filtroTransferencias)

	if erro != nil {
		return
	}

	for i := range response.Pagos {

		for j := range transferencias {
			if response.Pagos[i].UltimoPagoIntentoId == transferencias[j].Movimiento.PagointentosId {
				response.Pagos[i].ReferenciaBancaria = transferencias[j].ReferenciaBancaria
				response.Pagos[i].TransferenciaId = uint64(transferencias[j].ID)
				response.Pagos[i].FechaTransferencia = transferencias[j].FechaOperacion.Format("02-01-2006")
			}
		}
		// si el estado es PAID se acumula en un atributo de la struct ResponsePagos
		if response.Pagos[i].Estado == "PAID" {
			response.SaldoPendiente += response.Pagos[i].Amount
		}

		// si el estado es ACCREDITED y no tiene fecha de transferencia se acumula en un atributo de la struct ResponsePagos
		if response.Pagos[i].Estado == "ACCREDITED" && len(response.Pagos[i].FechaTransferencia) == 0 {
			response.SaldoDisponible += response.Pagos[i].Amount
		}
	}

	return

}

func (s *service) GetPlanCuotas(idMedioPago uint) (response []administraciondtos.PlanCuotasResponseDetalle, erro error) {
	response, erro = s.repository.GetPlanCuotasByMedioPago(idMedioPago)
	if erro != nil {
		erro = errors.New("problema al obtener plan de cuotas - " + erro.Error())
	}
	return
}

func (s *service) GetInteresesPlanes(fecha string) (planes []administraciondtos.PlanCuotasResponse, erro error) {
	// fmt.Printf("Fecha actual %v", fecha)
	fechaActual, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	response, erro := s.repository.GetInstallments(fechaActual)
	if erro != nil {
		return
	}
	for _, valueMedioPagoInstallment := range response {
		var details []administraciondtos.PlanCuotasResponseDetalle
		var installmentTemp administraciondtos.PlanCuotasResponse
		for _, valueInstallment := range valueMedioPagoInstallment.Installments {
			// logs.Info(valueInstallment.VigenciaHasta)
			if valueInstallment.VigenciaHasta == nil {
				installmentTemp = administraciondtos.PlanCuotasResponse{
					Id:                      valueInstallment.ID,
					Descripcion:             valueInstallment.Descripcion,
					MediopagoinstallmentsID: valueInstallment.MediopagoinstallmentsID,
				}
				for _, valueInstalmentDetail := range valueInstallment.Installmentdetail {
					details = append(details, administraciondtos.PlanCuotasResponseDetalle{
						InstallmentsID: valueInstallment.ID,
						Cuota:          uint(valueInstalmentDetail.Cuota),
						Tna:            valueInstalmentDetail.Tna,
						Tem:            valueInstalmentDetail.Tem,
						Coeficiente:    valueInstalmentDetail.Coeficiente,
					})
				}
				break
			}
			// logs.Info("============================\n")
			// fmt.Printf("after %v - before %v \n", valueInstallment.VigenciaDesde.After(fechaActual), valueInstallment.VigenciaDesde.Before(fechaActual))
			// logs.Info("============================\n")

			// fmt.Printf("%v--%v--%v \n", valueInstallment.VigenciaDesde, fechaActual, valueInstallment.VigenciaHasta)
			// logs.Info("============================\n")
			// logs.Info("============================\n")

			if (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.Equal(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Equal(*valueInstallment.VigenciaHasta)) {
				installmentTemp = administraciondtos.PlanCuotasResponse{
					Id:                      valueInstallment.ID,
					Descripcion:             valueInstallment.Descripcion,
					MediopagoinstallmentsID: valueInstallment.MediopagoinstallmentsID,
				}
				for _, valueInstalmentDetail := range valueInstallment.Installmentdetail {
					details = append(details, administraciondtos.PlanCuotasResponseDetalle{
						InstallmentsID: valueInstallment.ID,
						Cuota:          uint(valueInstalmentDetail.Cuota),
						Tna:            valueInstalmentDetail.Tna,
						Tem:            valueInstalmentDetail.Tem,
						Coeficiente:    valueInstalmentDetail.Coeficiente,
					})
				}
				break
			}

		}
		planes = append(planes, administraciondtos.PlanCuotasResponse{
			Id:                      installmentTemp.Id,
			Descripcion:             installmentTemp.Descripcion,
			MediopagoinstallmentsID: installmentTemp.MediopagoinstallmentsID,
			Installmentdetail:       details,
		})
	}
	return
}

func (s *service) GetAllInstallmentsById(installments_id int64) (planesCuotas []administraciondtos.InstallmentsResponse, erro error) {
	plancuotas, err := s.repository.GetAllInstallmentsById(uint(installments_id))
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	for _, value := range plancuotas {
		var temporalPlanCuotas administraciondtos.InstallmentsResponse
		temporalPlanCuotas.EntityToDtos(value)
		planesCuotas = append(planesCuotas, temporalPlanCuotas)
	}
	return
}

func (s *service) ModificarEstadoPagosExpirados() (erro error) {

	// Busco el tiempo de expiración de los pagos si no existe lo creo con valor de 30 dias

	filtroConf := filtros.ConfiguracionFiltro{
		Nombre: "TIEMPO_EXPIRACION_PAGOS",
	}

	configuracion, erro := s.utilService.GetConfiguracionService(filtroConf)

	if erro != nil {
		return
	}

	if configuracion.Id == 0 {

		config := administraciondtos.RequestConfiguracion{
			Nombre:      "TIEMPO_EXPIRACION_PAGOS",
			Descripcion: "Tiempo en días para que expire un pago que está en estado pending ",
			Valor:       "30",
		}

		_, erro = s.utilService.CreateConfiguracionService(config)

		if erro != nil {
			return
		}

	}

	// Busco el pagoEstado con nobre de pending

	filtroPending := filtros.PagoEstadoFiltro{
		Nombre: "Pending",
	}

	pagoEstadoPending, erro := s.repository.GetPagoEstado(filtroPending)

	if erro != nil {
		erro = fmt.Errorf("no se pudo obtener el  id de estado de pago pendiente")
		return
	}

	// Busco los pagos que están en el estado pending y que están expirados

	filtroPagos := filtros.PagoFiltro{
		PagoEstadosId:                 uint64(pagoEstadoPending.ID),
		VisualizarPendientes:          true,
		TiempoExpiracionSecondDueDate: "15",
	}

	pagos, _, erro := s.repository.GetPagos(filtroPagos)

	if erro != nil {
		return
	}

	if len(pagos) == 0 {
		erro = fmt.Errorf("no se encontraron pagos para actualizar")
		return
	}

	// Busco el pago estado expirado
	filtroExpired := filtros.PagoEstadoFiltro{
		Nombre: "Expired",
	}

	pagoEstadoExpired, erro := s.repository.GetPagoEstado(filtroExpired)

	if erro != nil {
		erro = fmt.Errorf("no se pudo obtener el id de estado de pago expirado")
		return
	}

	erro = s.repository.UpdateEstadoPagos(pagos, uint64(pagoEstadoExpired.ID))

	if erro != nil {
		erro = fmt.Errorf("no se pudo actualizar el estado de los pagos")
		return
	}
	return

}

func (s *service) GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (response administraciondtos.ResponseConfiguraciones, erro error) {

	configuraciones, total, erro := s.repository.GetConfiguraciones(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, c := range configuraciones {

		r := administraciondtos.ResponseConfiguracion{}
		r.FromEntity(c)

		response.Data = append(response.Data, r)
	}

	return
}

func (s *service) UpdateConfiguracionService(ctx context.Context, request administraciondtos.RequestConfiguracion) (erro error) {

	erro = request.IsValid(true)

	if erro != nil {
		return
	}

	config := request.ToEntity(true)

	erro = s.repository.UpdateConfiguracion(ctx, config)

	if erro != nil {
		return
	}
	return

}

// ABM CLIENTES
func (s *service) GetClienteService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacion, erro error) {
	cliente, erro := s.repository.GetCliente(filtro)

	if erro != nil {
		return
	}
	var cli administraciondtos.ResponseFacturacion
	cli.FromEntity(cliente)
	response = cli

	return
}

func (s *service) GetClientesService(filtro filtros.ClienteFiltro) (response administraciondtos.ResponseFacturacionPaginado, erro error) {

	clientes, total, erro := s.repository.GetClientes(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, cliente := range clientes {

		var cli administraciondtos.ResponseFacturacion
		cli.FromEntity(cliente)

		response.Clientes = append(response.Clientes, cli)
	}

	return
}

func (s *service) ObtenerClientesSinDTOService(filtro filtros.ClienteFiltro) (response []entities.Cliente, total int64, erro error) {

	response, total, erro = s.repository.GetClientes(filtro)

	if erro != nil {
		return
	}

	return
}

func (s *service) CreateClienteService(ctx context.Context, request administraciondtos.ClienteRequest) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	filtro := filtros.ClienteFiltro{
		Cuit: request.Cuit,
	}

	response, erro := s.repository.GetCliente(filtro)

	if erro != nil {
		return
	}

	if response.ID > 0 {
		erro = fmt.Errorf(ERROR_CLIENTE_REGISTRADO)
		return
	}

	cliente := request.ToCliente(false)

	return s.repository.CreateCliente(ctx, cliente)

}

func (s *service) UpdateClienteService(ctx context.Context, cliente administraciondtos.ClienteRequest) (erro error) {

	erro = cliente.IsVAlid(true)

	if erro != nil {
		return
	}

	filtro := filtros.ClienteFiltro{
		DistintoId: cliente.Id,
		Cuit:       cliente.Cuit,
	}

	cliente_existente_sin_modificar, erro := s.repository.GetCliente(filtro)

	if erro != nil {
		return
	}

	if cliente_existente_sin_modificar.ID == 0 {
		erro = fmt.Errorf(ERROR_CARGAR_CLIENTE)
		return
	}

	clienteModificado := cliente.ToCliente(true)

	erro = s.repository.UpdateCliente(ctx, clienteModificado)
	if erro != nil {
		logs.Error(erro)
		return
	}

	return
}

func (s *service) DeleteClienteService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id del cliente es invalido")
		return
	}

	erro = s.repository.DeleteCliente(ctx, id)
	if erro != nil {
		logs.Error(erro)
		return
	}

	return
}

//ABMRUBROS

func (s *service) GetRubroService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubro, erro error) {
	rubro, erro := s.repository.GetRubro(filtro)

	if erro != nil {
		return
	}
	response = administraciondtos.ResponseRubro{
		Id:    rubro.ID,
		Rubro: rubro.Rubro,
	}

	return
}

func (s *service) GetRubrosService(filtro filtros.RubroFiltro) (response administraciondtos.ResponseRubros, erro error) {

	rubros, total, erro := s.repository.GetRubros(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, rubro := range rubros {

		r := administraciondtos.ResponseRubro{
			Id:    rubro.ID,
			Rubro: rubro.Rubro,
		}

		response.Rubros = append(response.Rubros, r)
	}

	return
}

func (s *service) CreateRubroService(ctx context.Context, request administraciondtos.RubroRequest) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	rubro := request.ToRubro(false)

	return s.repository.CreateRubro(ctx, rubro)

}

func (s *service) UpdateRubroService(ctx context.Context, rubro administraciondtos.RubroRequest) (erro error) {

	erro = rubro.IsVAlid(true)

	if erro != nil {
		return
	}

	rubroModificado := rubro.ToRubro(true)

	return s.repository.UpdateRubro(ctx, rubroModificado)

}

//ABM PAGO TIPOS

func (s *service) GetPagoTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagoTipo, erro error) {

	pagoTipo, erro := s.repository.GetPagoTipo(filtro)

	if erro != nil {
		return
	}

	response.FromPagoTipo(pagoTipo)

	return
}

func (s *service) GetPagosTipoService(filtro filtros.PagoTipoFiltro) (response administraciondtos.ResponsePagosTipo, erro error) {

	pagosTipo, total, erro := s.repository.GetPagosTipo(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, pagoTipo := range pagosTipo {
		var ch []administraciondtos.CanalesPago
		var cuotas []administraciondtos.CuotasPago
		for _, channel := range pagoTipo.Pagotipochannel {
			c := administraciondtos.CanalesPago{
				ChannelsId: channel.Channel.ID,
				Channel:    channel.Channel.Channel,
				Nombre:     channel.Channel.Nombre,
			}
			ch = append(ch, c)
		}

		for _, cuota := range pagoTipo.Pagotipoinstallment {
			cuo := administraciondtos.CuotasPago{
				Nro: cuota.Cuota,
			}
			cuotas = append(cuotas, cuo)
		}

		r := administraciondtos.ResponsePagoTipo{}
		r.IncludedChannels = ch
		r.IncludedInstallments = cuotas
		r.FromPagoTipo(pagoTipo)

		response.PagosTipo = append(response.PagosTipo, r)
	}

	return
}

func (s *service) CreatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	pagoTipo := request.ToPagoTipo(false)

	for _, channel := range request.IncludedChannels {
		filtro := filtros.ChannelFiltro{
			Id: uint(channel),
		}
		ch, err := s.repository.GetChannel(filtro)
		if err != nil && ch.ID == 0 {
			erro = fmt.Errorf("el id del channels es invalido")
			return 0, erro
		}
	}

	return s.repository.CreatePagoTipo(ctx, pagoTipo, request.IncludedChannels, request.IncludedInstallments)

}

func (s *service) UpdatePagoTipoService(ctx context.Context, request administraciondtos.RequestPagoTipo) (erro error) {

	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}

	pagoTipoModificado := request.ToPagoTipo(true)

	filtro := filtros.PagoTipoFiltro{
		Id:                     pagoTipoModificado.ID,
		CargarTipoPagoChannels: true,
	}
	pagotipo, err := s.repository.GetPagoTipo(filtro)
	if err != nil {
		erro = err
		return
	}

	var channels []int64
	var cuotas []string
	for _, p := range pagotipo.Pagotipochannel {
		channels = append(channels, int64(p.Channel.ID))
	}
	for _, p := range pagotipo.Pagotipoinstallment {
		cuotas = append(cuotas, p.Cuota)
	}

	channelAdd, channelDelete := commons.DifferenceInt(request.IncludedChannels, channels)
	updateChannels := administraciondtos.RequestPagoTipoChannels{
		Add:    channelAdd,
		Delete: channelDelete,
	}

	cuotasAdd, cuotasDelete := commons.DifferenceString(request.IncludedInstallments, cuotas)
	updateCuotas := administraciondtos.RequestPagoTipoCuotas{
		Add:    cuotasAdd,
		Delete: cuotasDelete,
	}

	return s.repository.UpdatePagoTipo(ctx, pagoTipoModificado, updateChannels, updateCuotas)

}

func (s *service) DeletePagoTipoService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id del pago tipo es invalido")
		return
	}

	return s.repository.DeletePagoTipo(ctx, id)

}

//ABM CHANNEL

func (s *service) GetChannelService(filtro filtros.ChannelFiltro) (response administraciondtos.ResponseChannel, erro error) {

	channel, erro := s.repository.GetChannel(filtro)

	if erro != nil {
		return
	}

	response.FromChannel(channel)

	return
}

func (s *service) GetChannelsService(filtro filtros.ChannelFiltro) (response administraciondtos.ResponseChannels, erro error) {

	channels, total, erro := s.repository.GetChannels(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, channel := range channels {

		r := administraciondtos.ResponseChannel{}
		r.FromChannel(channel)

		response.Channels = append(response.Channels, r)
	}

	return
}

func (s *service) CreateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	channel := request.ToChannel(false)

	return s.repository.CreateChannel(ctx, channel)

}

func (s *service) UpdateChannelService(ctx context.Context, request administraciondtos.RequestChannel) (erro error) {

	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}

	channelModificado := request.ToChannel(true)

	return s.repository.UpdateChannel(ctx, channelModificado)

}

func (s *service) DeleteChannelService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id del channel es invalido")
		return
	}

	return s.repository.DeleteChannel(ctx, id)

}

//ABM CUENTAS COMISION

func (s *service) GetCuentaComisionService(filtro filtros.CuentaComisionFiltro) (response administraciondtos.ResponseCuentaComision, erro error) {

	cuentaComision, erro := s.repository.GetCuentaComision(filtro)

	if erro != nil {
		return
	}

	response.FromCuentaComision(cuentaComision)

	return
}

func (s *service) GetCuentasComisionService(filtro filtros.CuentaComisionFiltro) (response administraciondtos.ResponseCuentasComision, erro error) {

	cuentasComsion, total, erro := s.repository.GetCuentasComisiones(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, cuentaComision := range cuentasComsion {

		r := administraciondtos.ResponseCuentaComision{}
		r.FromCuentaComision(cuentaComision)

		response.CuentasComision = append(response.CuentasComision, r)
	}

	return
}

func (s *service) CreateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	filtro := filtros.ChannelFiltro{
		Id: uint(request.ChannelsId),
	}
	ch, err := s.repository.GetChannel(filtro)
	if err != nil && ch.ID == 0 {
		erro = fmt.Errorf("el id del channels es invalido")
		return 0, erro
	}

	cuentaComision := request.ToCuentaComision(false)

	return s.repository.CreateCuentaComision(ctx, cuentaComision)

}

func (s *service) UpdateCuentaComisionService(ctx context.Context, request administraciondtos.RequestCuentaComision) (erro error) {
	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}

	filtro := filtros.ChannelFiltro{
		Id: uint(request.ChannelsId),
	}
	ch, err := s.repository.GetChannel(filtro)
	if err != nil && ch.ID == 0 {
		erro = fmt.Errorf("el id del channels es invalido")
		return erro
	}

	cuentaComisionModificada := request.ToCuentaComision(true)

	return s.repository.UpdateCuentaComision(ctx, cuentaComisionModificada)

}

func (s *service) DeleteCuentaComisionService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id de cuenta comision es invalido")
		return
	}

	return s.repository.DeleteCuentaComision(ctx, id)

}

// Solicitud de cuenta
func (s *service) SendSolicitudCuenta(solicitudRequest administraciondtos.SolicitudCuentaRequest) (erro error) {

	// Valido los datos de entrada
	erro = solicitudRequest.IsValid()

	if erro != nil {

		return

	}

	entidadSolicitud := solicitudRequest.ToSolicitudEntity()
	// guadar la data de solicitud por medio del repository
	erro = s.repository.CreateSolicitudRepository(entidadSolicitud)
	if erro != nil {
		return
	}
	// Busco en la base de datos el email al cual será enviada la solicitud
	filtro := filtros.ConfiguracionFiltro{
		Nombre: "EMAIL_SOLICITUD_CUENTA",
	}

	configuracion, erro := s.utilService.GetConfiguracionService(filtro)

	if erro != nil {
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionConfiguraciones,
			Descripcion: fmt.Sprintf("Configuración inválida. %s", erro.Error()),
		}
		s.CreateNotificacionService(notificacion)
		return
	}

	// En caso de que no exista lo crea con el correo de prueba
	if configuracion.Id == 0 {

		config := administraciondtos.RequestConfiguracion{
			Nombre:      "EMAIL_SOLICITUD_CUENTA",
			Descripcion: "Email al cual se va enviar la solicitud de cuenta",
			Valor:       "developmenttelco@gmail.com",
		}

		_, erro = s.utilService.CreateConfiguracionService(config)

		if erro != nil {
			s._buildNotificacion(erro, entities.NotificacionConfiguraciones, fmt.Sprintf("Configuración inválida. %s", erro.Error()))
			return fmt.Errorf(ERROR_SOLICITUD_CUENTA)
		}

	}

	// Crear mensaje
	to := []string{
		configuracion.Valor,
	}

	from := config.EMAIL_FROM

	t, erro := template.ParseFiles("../api/views/solicitud_cuenta.html")
	if erro != nil {
		s._buildNotificacion(erro, entities.NotificacionSolicitudCuenta, fmt.Sprintf("no se pudo crear el template. %s", erro.Error()))
		return fmt.Errorf(ERROR_SOLICITUD_CUENTA)
	}
	buf := new(bytes.Buffer)
	erro = t.Execute(buf, solicitudRequest)
	if erro != nil {
		s._buildNotificacion(erro, entities.NotificacionSolicitudCuenta, fmt.Sprintf("no se pudo ejecutar el template. %s", erro.Error()))
		return fmt.Errorf(ERROR_SOLICITUD_CUENTA)
	}

	message := s.commonsService.CreateMessage(to, from, buf.String(), "Solicitud de Cuenta")

	password := config.EMAIL_PASS
	smtpHost := config.SMTPHOST
	smtpPort := config.SMTPPORT
	address := smtpHost + ":" + smtpPort

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 4- Sending email.
	erro = smtp.SendMail(address, auth, from, to, []byte(message))

	if erro != nil {
		s._buildNotificacion(erro, entities.NotificacionSolicitudCuenta, fmt.Sprintf("%s. No se pudo enviar el email de solicitud de cuenta. Solicitante: %s ,  Email: %s", erro.Error(), solicitudRequest.Razonsocial, solicitudRequest.Email))
		return fmt.Errorf(ERROR_SOLICITUD_CUENTA)
	}

	return
}

func (s *service) UpdateConfiguracionSendEmailService(ctx context.Context, request administraciondtos.RequestConfiguracion) (erro error) {

	//Busco todos los clientes activos del sistema para recuperar sus correos
	filtroCliente := filtros.ClienteFiltro{}

	response, erro := s.GetClientesService(filtroCliente)

	if len(response.Clientes) < 1 {
		return
	}
	//Cargo la lista de correos para enviar
	var emails []string
	for _, c := range response.Clientes {
		if !tools.EsStringVacio(c.Email) {
			emails = append(emails, c.Email)
		}
	}

	// Crear mensaje
	to := emails

	from := config.EMAIL_FROM

	t, erro := template.ParseFiles("../api/views/terminos_condiciones.html")

	ruta := administraciondtos.TerminosCondiciones{
		Ruta: config.RUTA_BASE_HOME_PAGE + "terminos-politicas",
	}

	if erro != nil {
		s.utilService.BuildLog(erro, "SendTerminosCondiciones")
		return fmt.Errorf(ERROR_ENVIAR_EMAIL_TERMINOS_CONDICIONES)
	}
	buf := new(bytes.Buffer)
	erro = t.Execute(buf, ruta)
	if erro != nil {
		s.utilService.BuildLog(erro, "SendTerminosCondiciones")
		return fmt.Errorf(ERROR_ENVIAR_EMAIL_TERMINOS_CONDICIONES)
	}

	message := s.commonsService.CreateMessage(to, from, buf.String(), "Actualización Terminos y Condiciones")

	password := config.EMAIL_PASS
	smtpHost := config.SMTPHOST
	smtpPort := config.SMTPPORT
	address := smtpHost + ":" + smtpPort

	//Modifico los terminos y condiciones para luego enviar el mensaje
	s.repository.BeginTx()

	defer func() {
		if erro != nil {
			s.repository.RollbackTx()
		}
		s.repository.CommitTx()
	}()

	erro = s.UpdateConfiguracionService(ctx, request)

	if erro != nil {
		// s.repository.RollbackTx()
		return
	}

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 4- Sending email.
	erro = smtp.SendMail(address, auth, from, to, []byte(message))

	if erro != nil {
		s.repository.RollbackTx()
		s._buildNotificacion(erro, entities.NotivicacionEnvioEmail, fmt.Sprintf("%s. No se pudo enviar el email de actualización de terminos y condiciones a los clientes.", erro.Error()))
		return fmt.Errorf(ERROR_ENVIAR_EMAIL_TERMINOS_CONDICIONES)
	}

	// s.repository.CommitTx()

	return
}

func (s *service) _buildNotificacion(erro error, tipo entities.EnumTipoNotificacion, descripcion string) {
	notificacion := entities.Notificacione{
		Tipo:        tipo,
		Descripcion: descripcion,
	}
	s.CreateNotificacionService(notificacion)
}

func (s *service) RetiroAutomaticoClientes(ctx context.Context, request administraciondtos.TransferenciasClienteId) (movimientosidcomision administraciondtos.RequestMovimientosId, erro error) {

	//Buscar todos los clientes que tengan habilitado la opción de retiro automático
	// ? se habilita transferencias para automaticas para varios clientes
	filtroCliente := filtros.ClienteFiltro{
		Id:               request.Id,
		RetiroAutomatico: true,
		CargarCuentas:    true,
	}

	clientes, _, erro := s.repository.GetClientes(filtroCliente)

	if erro != nil {
		return
	}

	// variable creada en el caso de que existan errore al enviar transferenicas a cuentas
	// var responseTransferencias []administraciondtos.ResponseTransferenciaAutomatica
	//Si no hay ningun cliente habilitado no debe hacer nada
	if len(clientes) > 0 {
		//TODO Aquí se podría usar una go routina para buscar el motivo y la cuenta de telco.
		//Busco el motivo por defecto para las transferencias
		filtroMotivo := filtros.ConfiguracionFiltro{
			Nombre: "MOTIVO_TRANSFERENCIA_CLIENTE",
		}

		motivo, erro := s.utilService.GetConfiguracionService(filtroMotivo)

		if erro != nil {
			return administraciondtos.RequestMovimientosId{}, erro
		}

		// Busco la cuenta de telco para las transferencias

		filtroCuenta := filtros.ConfiguracionFiltro{
			Nombre: "CBU_CUENTA_TELCO",
		}

		cbu, erro := s.utilService.GetConfiguracionService(filtroCuenta)

		if erro != nil {
			return administraciondtos.RequestMovimientosId{}, erro
		}

		if len(cbu.Valor) < 1 {
			erro = fmt.Errorf("no se pudo encontrar el cbu de la cuenta de origen")
			return administraciondtos.RequestMovimientosId{}, erro
		}

		var listaTransferencias []administraciondtos.RequestTransferenciaAutomatica

		for _, c := range clientes {
			for _, cu := range *c.Cuentas {
				if cu.TransferenciaAutomatica {
					filtroMovimiento := filtros.MovimientoFiltro{
						AcumularPorPagoIntentos: true,
						CuentaId:                uint64(cu.ID),
						CargarPago:              true,
						CargarPagoEstados:       true,
						CargarPagoIntentos:      true,
						CargarMedioPago:         true,
						// FechaInicio:             fechaI.AddDate(0, 0, int(-cu.DiasRetiroAutomatico)),
						// FechaFin:                fechaF,
						CargarMovimientosNegativos: true,
					}
					// filtroMovimiento.CuentaId = uint64(cu.ID)
					movimientos, err := s.GetMovimientosAcumulados(filtroMovimiento)
					if err != nil {
						erro = err
						return administraciondtos.RequestMovimientosId{}, erro
					}
					if len(movimientos.Acumulados) > 0 {
						var listaIdsMovimiento []uint64
						var listaMovimientosIdNeg []uint64
						var acumulado entities.Monto
						var acumuladoNeg entities.Monto
						for _, ma := range movimientos.Acumulados {
							acumulado += ma.Acumulado
							for _, m := range ma.Movimientos {
								listaIdsMovimiento = append(listaIdsMovimiento, uint64(m.Id))
							}
						}
						for _, mn := range movimientos.MovimientosNegativos {
							acumuladoNeg += mn.Monto
							listaMovimientosIdNeg = append(listaMovimientosIdNeg, uint64(mn.Id))
						}
						if acumuladoNeg < 0 {
							acumulado += acumuladoNeg
						}
						request := administraciondtos.RequestTransferenciaAutomatica{
							CuentaId: uint64(cu.ID),
							Cuenta:   cu.Cuenta,
							DatosClientes: administraciondtos.DatosClientes{
								NombreCliente: c.Cliente,
								EmailCliente:  c.Email,
							},
							Request: administraciondtos.RequestTransferenicaCliente{
								Transferencia: linktransferencia.RequestTransferenciaCreateLink{
									Origen: linktransferencia.OrigenTransferenciaLink{
										Cbu: cbu.Valor,
									},
									Destino: linktransferencia.DestinoTransferenciaLink{
										Cbu:            cu.Cbu,
										EsMismoTitular: false,
									},
									Importe: acumulado,
									Moneda:  linkdtos.Pesos,
									Motivo:  linkdtos.EnumMotivoTransferencia(motivo.Valor),
								},
								ListaMovimientosId:    listaIdsMovimiento,
								ListaMovimientosIdNeg: listaMovimientosIdNeg,
							},
						}

						listaTransferencias = append(listaTransferencias, request)
					}
				}
			} // fin for range cuentas
		} // fin for range clientes
		uuid := uuid.NewV4()
		_, err := s.commonsService.IsValidUUID(uuid.String())

		if err != nil {
			erro = err
			return administraciondtos.RequestMovimientosId{}, erro
		}
		// NOTE obtener token apilink
		scopes := []linkdtos.EnumScopeLink{linkdtos.TransferenciasBancariasInmediatas}
		token, err := s.apilinkService.GetTokenApiLinkService(uuid.String(), scopes)
		if err != nil {
			erro = err
			return administraciondtos.RequestMovimientosId{}, erro
		}
		var idmovcomisiones administraciondtos.RequestMovimientosId
		for _, t := range listaTransferencias {
			if t.Request.Transferencia.Importe > 0 {
				time.Sleep(5 * time.Second)
				logs.Info(fmt.Sprintf("ejecutando transferencia automatica cliente: %s", t.DatosClientes.NombreCliente))
				fmt.Println("Han pasado 5 segundos. Continuando con la acción.")
				// uuid := uuid.NewV4()
				response, erro := s.BuildTransferencia(ctx, uuid.String(), t.Request, t.CuentaId, t.DatosClientes, token.AccessToken)
				if erro != nil {
					aviso := entities.Notificacione{
						Tipo:        entities.NotificacionTransferenciaAutomatica,
						Descripcion: fmt.Sprintf("atención no se pudo realizar transferencia automatica de la cuenta: %s. Error: %s", t.Cuenta, erro.Error()),
						UserId:      0,
					}
					erro := s.utilService.CreateNotificacionService(aviso)
					if erro != nil {
						logs.Error(erro.Error() + "no se pudo crear notificación en BuildTransferencia")
					}
				} else {
					// acumular id de movimientos transferidos
					idmovcomisiones.MovimientosId = append(idmovcomisiones.MovimientosId, response.MovimientosIdTransferidos...)
					//idmovcomisiones.MovimimientosIdRevertidos = append(idmovcomisiones.MovimimientosIdRevertidos, response.MovimientosIdReversiones...)
				}
			}
		}
		movimientosidcomision = idmovcomisiones
	}
	return
}

func (s *service) RetiroAutomaticoClientesSubcuentas(ctx context.Context) (movimientosidcomision administraciondtos.RequestMovimientosId, erro error) {

	/*NOTE - Proceso que ejecutara transferencias automaticas para clientes con subcuentas */
	filtroCliente := filtros.ClienteFiltro{
		CargarCuentas:    true,
		SplitCuentas:     true,
		CargarSubcuentas: true,
	}

	clientes, _, erro := s.repository.GetClientes(filtroCliente)
	if erro != nil {
		return
	}

	// solo debe cargar losc clientes que tengan configurado split_cuentas
	if len(clientes) > 0 {
		//TODO Aquí se podría usar una go routina para buscar el motivo y la cuenta de telco.
		//Busco el motivo por defecto para las transferencias
		filtroMotivo := filtros.ConfiguracionFiltro{
			Nombre: "MOTIVO_TRANSFERENCIA_CLIENTE",
		}

		motivo, erro := s.utilService.GetConfiguracionService(filtroMotivo)

		if erro != nil {
			return administraciondtos.RequestMovimientosId{}, erro
		}

		// Busco la cuenta de telco para las transferencias

		filtroCuenta := filtros.ConfiguracionFiltro{
			Nombre: "CBU_CUENTA_TELCO",
		}

		cbu, erro := s.utilService.GetConfiguracionService(filtroCuenta)

		if erro != nil {
			return administraciondtos.RequestMovimientosId{}, erro
		}

		if len(cbu.Valor) < 1 {
			erro = fmt.Errorf("no se pudo encontrar el cbu de la cuenta de origen")
			return administraciondtos.RequestMovimientosId{}, erro
		}

		var listaTransferencias []administraciondtos.RequestTransferenciaAutomatica

		for _, c := range clientes {
			if c.RetiroAutomatico {
				for _, cu := range *c.Cuentas {
					if len(cu.Subcuentas) > 0 {
						filtroMovimiento := filtros.MovimientoFiltro{
							AcumularPorPagoIntentos:     true,
							CuentaId:                    uint64(cu.ID),
							CargarPago:                  true,
							CargarPagoEstados:           true,
							CargarPagoIntentos:          true,
							CargarMedioPago:             true,
							CargarMovimientosNegativos:  true,
							CargarMovimientosSubcuentas: true,
						}
						movimientos, err := s.GetMovimientosSubcuentas(filtroMovimiento)
						if err != nil {
							erro = err
							return administraciondtos.RequestMovimientosId{}, erro
						}
						if len(movimientos.Acumulados) > 0 {
							for _, sub := range cu.Subcuentas { //recorrer subcuentas para cada cuenta
								var acumulado entities.Monto // acumular para cada subcuenta
								var acumuladoNeg entities.Monto
								var listaIdsMovimiento []uint64
								var listaMovimientosIdNeg []uint64

								// NOTE 1 - acumular monto por cada subcuenta
								for _, ma := range movimientos.Acumulados {
									listaIdsMovimiento = append(listaIdsMovimiento, uint64(ma.Id))
									for _, ms := range ma.MovimientoSubcuentas { // recorrer los movimientossubcuentas
										if uint64(sub.ID) == ms.SubcuentasID {
											acumulado += ms.Monto
										}
									}
								}
								// NOTE 2 - acumular monto negativos o reversiones por cada subcuenta
								// FIXME Ver casos de movnegativos
								for _, mn := range movimientos.MovimientosNegativos {
									listaMovimientosIdNeg = append(listaMovimientosIdNeg, uint64(mn.Id))
									for _, m := range mn.MovimientoSubcuentas {
										if uint64(sub.ID) == m.SubcuentasID {
											acumuladoNeg += m.Monto
										}
									}
								}
								if acumuladoNeg < 0 { // Reversiones o mov negativos
									acumulado += acumuladoNeg
								}
								request := administraciondtos.RequestTransferenciaAutomatica{
									CuentaId:    uint64(cu.ID),
									SubcuentaId: uint64(sub.ID),
									Subcuenta:   sub.Nombre,
									DatosClientes: administraciondtos.DatosClientes{
										NombreCliente: sub.Nombre,
										EmailCliente:  sub.Email,
									},
									Request: administraciondtos.RequestTransferenicaCliente{
										Transferencia: linktransferencia.RequestTransferenciaCreateLink{
											Origen: linktransferencia.OrigenTransferenciaLink{
												Cbu: cbu.Valor,
											},
											Destino: linktransferencia.DestinoTransferenciaLink{
												Cbu:            cu.Cbu,
												EsMismoTitular: false,
											},
											Importe: acumulado,
											Moneda:  linkdtos.Pesos,
											Motivo:  linkdtos.EnumMotivoTransferencia(motivo.Valor),
										},
										ListaMovimientosId:    listaIdsMovimiento,
										ListaMovimientosIdNeg: listaMovimientosIdNeg,
									},
								}
								listaTransferencias = append(listaTransferencias, request)
							}
						}
					}
				} // fin for range cuentas
			}
		} // fin for range clientes

		logs.Info(listaTransferencias)

		var idmovcomisiones administraciondtos.RequestMovimientosId
		for _, t := range listaTransferencias {
			if t.Request.Transferencia.Importe > 0 {
				uuid := uuid.NewV4()
				response, erro := s.BuildTransferenciaSubcuentas(ctx, uuid.String(), t.Request, t.CuentaId, t.DatosClientes)
				if erro != nil {
					aviso := entities.Notificacione{
						Tipo:        entities.NotificacionTransferenciaAutomatica,
						Descripcion: fmt.Sprintf("atención no se pudo realizar transferencia automatica de la cuenta: %s. Error: %s", t.Subcuenta, erro.Error()),
						UserId:      0,
					}
					erro := s.utilService.CreateNotificacionService(aviso)
					if erro != nil {
						logs.Error(erro.Error() + "no se pudo crear notificación en BuildTransferencia")
					}
				} else {
					// acumular id de movimientos transferidos
					idmovcomisiones.MovimientosId = append(idmovcomisiones.MovimientosId, response.MovimientosIdTransferidos...)
					//idmovcomisiones.MovimimientosIdRevertidos = append(idmovcomisiones.MovimimientosIdRevertidos, response.MovimientosIdReversiones...)
				}
			}
		}
		movimientosidcomision = idmovcomisiones
	}
	return
}

func (s *service) CreatePlanCuotasService(request administraciondtos.RequestPlanCuotas) (erro error) {

	/* cnovertir datos para ser procesados */
	installmentId, fechaVigencia, err := procesarRequest(request.InstalmentsId, request.VigenciaDesde)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}

	installmentActual, err := s.repository.GetInstallmentById(uint(installmentId))
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	// fmt.Printf("%v", installmentActual.VigenciaDesde.After(fechaVigencia))
	// if installmentActual.VigenciaDesde.After(fechaVigencia) || installmentActual.VigenciaDesde.Equal(fechaVigencia) {}
	if installmentActual.VigenciaHasta != nil {
		erro = errors.New(ERROR_CONSULTA_INSTALLMENT)
		return
	}
	fechaHasta := fechaVigencia.Add(time.Hour * 24 * -1)
	installmentActual.VigenciaHasta = &fechaHasta
	fmt.Printf("%v - %v \n", fechaVigencia, fechaHasta)

	installmentNew := entities.Installment{
		MediopagoinstallmentsID: installmentActual.MediopagoinstallmentsID,
		Descripcion:             installmentActual.Descripcion,
		Issuer:                  installmentActual.Issuer,
		VigenciaDesde:           fechaVigencia,
	}

	openFile, err := os.Open(request.RutaFile)
	if err != nil {
		erro = errors.New(ERROR_LEER_ARCHIVO)
		return
	}
	readFile := csv.NewReader(openFile)
	readFile.Comma = ';'
	readFile.FieldsPerRecord = -1
	flag := false
	var listPlanescuotas []entities.Installmentdetail
	for {
		registro, err := readFile.Read()
		if err != nil && err != io.EOF {
			erro = errors.New(ERROR_LEER_ARCHIVO)
			return
		}
		if err == io.EOF {
			break
		}
		if !flag {
			flag = true
			listPlanescuotas = append(listPlanescuotas, entities.Installmentdetail{
				InstallmentsID: 0,
				Activo:         false,
				Cuota:          1,
				Tna:            0,
				Tem:            0,
				Coeficiente:    1,
				Fechadesde:     fechaVigencia,
			})
			continue
		}
		cuota, tna, tem, coeficiente, err := procesarRegistro(registro)
		if err != nil {
			erro = err
			return
		}
		listPlanescuotas = append(listPlanescuotas, entities.Installmentdetail{
			InstallmentsID: 0,
			Activo:         false,
			Cuota:          cuota,
			Tna:            tna,
			Tem:            tem,
			Coeficiente:    coeficiente,
			Fechadesde:     fechaVigencia,
		})

	}
	openFile.Close()
	err = s.repository.CreatePlanCuotasByInstallmenIdRepository(installmentActual, installmentNew, listPlanescuotas)
	if err != nil {
		erro = errors.New(err.Error())
		logs.Error(erro)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       err.Error(),
			Funcionalidad: "CreatePlanCuotasByInstallmenIdRepository - repository",
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Error(err)
		}
		return
	}

	err = s.commonsService.BorrarDirectorio(request.RutaFile)
	if err != nil {
		erro = errors.New(err.Error())
		logs.Error(erro)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       err.Error(),
			Funcionalidad: "BorrarDirectorio - commonService",
		}
		err := s.utilService.CreateLogService(log)
		if err != nil {
			logs.Error(err)
		}
		return
	}
	return
}

func procesarRegistro(planCuota []string) (cuota int64, tna, tem, coeficiente float64, erro error) {

	cuota, err := strconv.ParseInt(planCuota[0], 10, 64)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	tna, err = strconv.ParseFloat(planCuota[1][0:len(planCuota[1])-1], 64)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	tem, err = strconv.ParseFloat(planCuota[2][0:len(planCuota[2])-1], 64)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	coeficiente, err = strconv.ParseFloat(planCuota[3], 64)

	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	return
}

func procesarRequest(installmentid, fecha string) (idInstallments int, fechaVigencia time.Time, erro error) {
	idInstallments, err := strconv.Atoi(installmentid)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	fechaVigencia, err = time.Parse("2006-01-02T00:00:00Z", fecha)
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}
	return
}

// func formatFecha() (fechaI time.Time, fechaF time.Time, erro error) {
// startTime := time.Now()
// fechaConvert := startTime.Format("2006-01-02") //YYYY.MM.DD
// fec := strings.Split(fechaConvert, "-")
//
// dia, err := strconv.Atoi(fec[len(fec)-1])
// if err != nil {
// erro = errors.New(ERROR_CONVERSION_DATO)
// return
// }
//
// mes, err := strconv.Atoi(fec[1])
// if err != nil {
// erro = errors.New(ERROR_CONVERSION_DATO)
// return
// }
//
// anio, err := strconv.Atoi(fec[0])
// if err != nil {
// erro = errors.New(ERROR_CONVERSION_DATO)
// return
// }
//
// fechaI = time.Date(anio, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)
// fechaF = time.Date(anio, time.Month(mes), dia, 23, 59, 59, 0, time.UTC)
//
// return
// }

func (s *service) GetPeticionesService(filtro filtros.PeticionWebServiceFiltro) (peticiones administraciondtos.ResponsePeticionesWebServices, erro error) {

	peticionesRes, total, erro := s.repository.GetPeticionesWebServices(filtro)
	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		peticiones.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, p := range peticionesRes {
		peticiones.Peticiones = append(peticiones.Peticiones, administraciondtos.ResponsePeticionWebServices{
			Operacion: p.Operacion,
			Vendor:    string(p.Vendor),
		})
	}

	return peticiones, nil
}

func (s *service) GetPagosTipoChannelService(filtro filtros.PagoTipoChannelFiltro) (response []entities.Pagotipochannel, erro error) {
	return s.repository.GetPagosTipoChannelRepository(filtro)
}

func (s *service) DeletePagoTipoChannelService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf(ERROR_ID)
		return
	}

	erro = s.repository.DeletePagoTipoChannel(id)
	if erro != nil {
		logs.Error(erro)
		return
	}

	if erro != nil {
		return erro
	}

	return
}

func (s *service) CreatePagoTipoChannel(ctx context.Context, request administraciondtos.RequestPagoTipoChannel) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	filtro := filtros.PagoTipoFiltro{
		Id: request.PagoTipoId,
	}

	_, erro = s.repository.GetPagoTipo(filtro)

	if erro != nil {
		return
	}

	filtro2 := filtros.ChannelFiltro{
		Id: request.ChannelId,
	}

	_, erro = s.repository.GetChannel(filtro2)

	if erro != nil {
		return
	}

	filtro3 := filtros.PagoTipoChannelFiltro{
		PagoTipoId: request.PagoTipoId,
		ChannelId:  request.ChannelId,
	}

	res, erro := s.repository.GetPagosTipoChannelRepository(filtro3)

	if erro != nil {
		return
	}
	if len(res) > 0 {
		return uint64(res[0].ID), nil
	}

	pagotipochannel := request.ToPagoTipoChannel(false)

	return s.repository.CreatePagoTipoChannel(ctx, pagotipochannel)

}

func (s *service) SubirArchivos(ctx context.Context, rutaArchivos string, listaArchivo []administraciondtos.ArchivoResponse) (countArchivo int, erro error) {
	/*
		por ultimo se mueven los archvios del directorio temporal
		a un directorio en minio dondo se almacenan los archivos
		de cierre de lote registrado en la DB
	*/
	var rutaDestino string
	for _, archivo := range listaArchivo {
		/*
			se lee el contenido del archivo y se obtiene su contenido se le pasa:
			- ruta destino
			- ruta origen del archivo
			- nombre del archivo
		*/
		data, filename, filetypo, err := util.LeerDatosArchivo(rutaDestino, rutaArchivos, archivo.NombreArchivo)
		filename = config.DIR_KEY + filename
		if err != nil {
			logs.Error(err)
		}
		/*	necesito la data, nombre del archivo y el tipo */
		filenameWithoutExt := filename[:len(filename)-len(filepath.Ext(filename))]
		erro := s.store.PutObject(ctx, data, filenameWithoutExt, filetypo)
		if erro != nil {
			logs.Error("No se pudo guardar el archivo")
		}

	}
	nombreDirectorio := config.DIR_KEY

	for _, archivoValue := range listaArchivo {
		/* antes de borrar el archivo se verifica si:
		el archivo fue leido, si fue movido y si la informacion que contiene el archivo se inserto en la db */
		erro = s.commonsService.BorrarArchivo(rutaArchivos, archivoValue.NombreArchivo)
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "BorrarArchivos",
				Mensaje:       erro.Error(),
			}
			erro = s.utilService.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				return 0, erro
			}
		}

		key := nombreDirectorio + "/" + archivoValue.NombreArchivo //config.DIR_KEY + "/" + archivoValue.NombreArchivo
		erro = s.store.DeleteObject(ctx, key)
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "DeleteObject",
				Mensaje:       erro.Error(),
			}
			erro = s.utilService.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				return 0, erro
			}
		}
	}
	erro = s.commonsService.BorrarDirectorio(rutaArchivos)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return 0, erro
		}
	}

	return 1, erro
}

func (s *service) SubirArchivosCloud(ctx context.Context, rutaArchivos string, listaArchivo []administraciondtos.ArchivoResponse, directorio string) (countArchivo int, erro error) {

	var rutaDestino string
	for _, archivo := range listaArchivo {
		/*
			se lee el contenido del archivo y se obtiene su contenido se le pasa:
			- ruta destino
			- ruta origen del archivo
			- nombre del archivo
		*/
		data, filename, filetypo, err := util.LeerDatosArchivo(rutaDestino, rutaArchivos, archivo.NombreArchivo)

		filename = directorio + "/" + archivo.NombreArchivo
		split := strings.Split(filename, "/")
		filenameServidor := strings.Join(split[1:(len(split))], "/")

		if err != nil {
			logs.Error(err)
		}
		// /*	necesito la data, nombre del archivo y el tipo */
		filenameWithoutExt := filenameServidor[:len(filenameServidor)-len(filepath.Ext(filenameServidor))]
		fmt.Print(filenameWithoutExt, data, filetypo)
		erro = s.store.PutObject(ctx, data, filenameWithoutExt, filetypo)
		if erro != nil {
			logs.Error("No se pudo guardar el archivo")
			return
		}

	}
	// nombreDirectorio := directorio

	for _, archivoValue := range listaArchivo {
		/* antes de borrar el archivo se verifica si:
		el archivo fue leido, si fue movido y si la informacion que contiene el archivo se inserto en la db */
		erro = s.commonsService.BorrarArchivo(rutaArchivos, archivoValue.NombreArchivo)
		if erro != nil {
			logs.Error(erro.Error())
			log := entities.Log{
				Tipo:          entities.EnumLog("Error"),
				Funcionalidad: "BorrarArchivos",
				Mensaje:       erro.Error(),
			}
			erro = s.utilService.CreateLogService(log)
			if erro != nil {
				logs.Error("error: al crear logs: " + erro.Error())
				return 0, erro
			}
		}
	}
	erro = s.commonsService.BorrarDirectorio(rutaArchivos)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return 0, erro
		}
	}

	return 1, erro
}

func (s *service) GetChannelsArancelService(filtro filtros.ChannelArancelFiltro) (response administraciondtos.ResponseChannelsArancel, erro error) {

	channelaranc, total, erro := s.repository.GetChannelsAranceles(filtro)

	logs.Info(channelaranc)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, ca := range channelaranc {

		r := administraciondtos.ResponseChannelsAranceles{}
		r.FromChannelArancel(ca)

		response.ChannelArancel = append(response.ChannelArancel, r)
	}

	return
}

func (s *service) CreateChannelsArancelService(ctx context.Context, request administraciondtos.RequestChannelsAranncel) (id uint64, erro error) {

	erro = request.IsVAlid(false)

	if erro != nil {
		return
	}

	filtro := filtros.ChannelFiltro{
		Id: uint(request.ChannelsId),
	}
	ch, err := s.repository.GetChannel(filtro)
	if err != nil && ch.ID == 0 {
		erro = fmt.Errorf("el id del channels es invalido")
		return 0, erro
	}

	channelsArancel := request.ToChannelsArancel(false)

	return s.repository.CreateChannelsArancel(ctx, channelsArancel)

}

func (s *service) UpdateChannelsArancelService(ctx context.Context, request administraciondtos.RequestChannelsAranncel) (erro error) {
	erro = request.IsVAlid(true)

	if erro != nil {
		return
	}

	filtro := filtros.ChannelFiltro{
		Id: uint(request.ChannelsId),
	}
	ch, err := s.repository.GetChannel(filtro)
	if err != nil && ch.ID == 0 {
		erro = fmt.Errorf("el id del channels es invalido")
		return erro
	}

	channelArancelModificada := request.ToChannelsArancel(true)

	return s.repository.UpdateChannelsArancel(ctx, channelArancelModificada)

}

func (s *service) DeleteChannelsArancelService(ctx context.Context, id uint64) (erro error) {

	if id < 1 {
		erro = fmt.Errorf("el id de channel arancel es invalido")
		return
	}

	return s.repository.DeleteChannelsArancel(ctx, id)

}

func (s *service) GetChannelArancelService(filtro filtros.ChannelAranceFiltro) (response administraciondtos.ResponseChannelsAranceles, erro error) {

	channel_arance, erro := s.repository.GetChannelArancel(filtro)

	if erro != nil {
		return
	}

	response.FromChArancel(channel_arance)

	return
}

func (s *service) ObtenerArchivosSubidos(filtro filtros.Paginacion) (lisArchivosSubidos administraciondtos.ResponseArchivoSubido, erro error) {
	var contador int64
	var recorrerHasta int32
	var listaTemporalArchivo []administraciondtos.ArchivoSubido
	entityCl, err := s.repository.GetCierreLoteSubidosRepository()
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New(err.Error())
		return
	}
	if len(entityCl) > 0 {
		for _, valueCL := range entityCl {
			contador++
			var listaClTemporal administraciondtos.ArchivoSubido
			listaClTemporal.EntityClToDtos(&valueCL)
			listaTemporalArchivo = append(listaTemporalArchivo, listaClTemporal)
		}
	}

	entityPx, err := s.repository.GetPrismaPxSubidosRepository()
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New(err.Error())
		return
	}
	if len(entityPx) > 0 {
		for _, valuePx := range entityPx {
			contador++
			var listaPxTemporal administraciondtos.ArchivoSubido
			listaPxTemporal.EntityPxToDtos(&valuePx)
			listaTemporalArchivo = append(listaTemporalArchivo, listaPxTemporal)
		}
	}

	entityMx, err := s.repository.GetPrismaMxSubidosRepository()
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New(err.Error())
		return
	}
	if len(entityMx) > 0 {
		for _, valueMx := range entityMx {
			contador++
			var listaMxTemporal administraciondtos.ArchivoSubido
			listaMxTemporal.EntityMxToDtos(&valueMx)
			listaTemporalArchivo = append(listaTemporalArchivo, listaMxTemporal)
		}
	}

	sort.Slice(listaTemporalArchivo, func(i, j int) bool {
		return listaTemporalArchivo[i].FechaSubida.Before(listaTemporalArchivo[j].FechaSubida)
	})

	if filtro.Number > 0 && filtro.Size > 0 {
		lisArchivosSubidos.Meta = _setPaginacion(filtro.Number, filtro.Size, contador)
	}
	recorrerHasta = lisArchivosSubidos.Meta.Page.To
	if lisArchivosSubidos.Meta.Page.CurrentPage == lisArchivosSubidos.Meta.Page.LastPage {
		recorrerHasta = lisArchivosSubidos.Meta.Page.Total
	}
	if len(listaTemporalArchivo) > 0 {
		for i := lisArchivosSubidos.Meta.Page.From; i < recorrerHasta; i++ {
			lisArchivosSubidos.ArchivosSubidos = append(lisArchivosSubidos.ArchivosSubidos, listaTemporalArchivo[i])
		}
	}
	return
}

func (s *service) ObtenerArchivoCierreLoteRapipago(nombre string) (archivo bool, err error) {

	archivo, err = s.repository.ObtenerArchivoCierreLoteRapipago(nombre)

	if err != nil {
		return
	}

	return
}

func (s *service) ObtenerArchivoCierreLoteMultipagos(nombre string) (archivo bool, err error) {

	archivo, err = s.repository.ObtenerArchivoCierreLoteMultipagos(nombre)

	if err != nil {
		return
	}

	return
}

func (s *service) GetCierreLoteEnDisputaServices(estadoDisputa int, request filtros.ContraCargoEnDisputa) (cierreLoteDisputa []cierrelotedtos.ResponsePrismaCL, erro error) {

	var clTemporal cierrelotedtos.ResponsePrismaCL
	operacionesEnDisputa, err := s.repository.ObtenerCierreLoteEnDisputaRepository(estadoDisputa, request)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	operacionesContracargo, err := s.repository.ObtenerCierreLoteContraCargoRepository(estadoDisputa, request)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}

	for _, valueCl := range operacionesEnDisputa {
		clTemporal.EntityToDtos(valueCl)
		cierreLoteDisputa = append(cierreLoteDisputa, clTemporal)
	}

	for _, valueCl := range operacionesContracargo {
		clTemporal.EntityToDtos(valueCl)
		cierreLoteDisputa = append(cierreLoteDisputa, clTemporal)
	}

	return
}

func (s *service) GetPagosByTransactionIdsServices(filtro filtros.ContraCargoEnDisputa, cierreLoteDisputa []cierrelotedtos.ResponsePrismaCL) (listaRevertidos administraciondtos.ResponseOperacionesContracargo, erro error) {
	for _, value := range cierreLoteDisputa {
		filtro.TransactionId = append(filtro.TransactionId, value.ExternalclienteID)
	}
	// obtener cuenta y pago tipos

	// obtengo pagos por pago tipo id

	listaPagosRevertidos, err := s.repository.ObtenerPagosInDisputaRepository(filtro)
	if err != nil {
		logs.Error(err.Error())
		erro = errors.New(err.Error())
		return
	}
	if len(listaPagosRevertidos) > 0 {
		var pagostemporal []administraciondtos.ResponsePagoCC
		var pagosTipoTemporal []administraciondtos.ResponsePagotipoCC

		listaRevertidos.Cuenta.Id = listaPagosRevertidos[0].Pago.PagosTipo.Cuenta.ID
		listaRevertidos.Cuenta.Cuenta = listaPagosRevertidos[0].Pago.PagosTipo.Cuenta.Cuenta
		var pagoTipoId uint
		for _, value := range listaPagosRevertidos {
			if value.Pago.PagosTipo.ID != pagoTipoId {
				pagoTipoId = value.Pago.PagosTipo.ID
				pagosTipoTemporal = append(pagosTipoTemporal, administraciondtos.ResponsePagotipoCC{
					Id:       value.Pago.PagosTipo.ID,
					Pagotipo: value.Pago.PagosTipo.Pagotipo,
				})

			}
		}
		for _, valuePago := range listaPagosRevertidos {
			pagostemporal = append(pagostemporal, administraciondtos.ResponsePagoCC{
				Id:                  valuePago.Pago.ID,
				PagostipoID:         valuePago.Pago.PagostipoID,
				Fecha:               valuePago.Pago.FirstDueDate,
				ExternalReference:   valuePago.Pago.ExternalReference,
				PayerName:           strings.ToUpper(valuePago.Pago.PayerName),
				Estado:              "",
				NombreEstado:        "",
				Amount:              valuePago.Pago.FirstTotal,
				FechaPago:           valuePago.Pago.CreatedAt,
				Channel:             "",
				NombreChannel:       "",
				UltimoPagoIntentoId: uint64(valuePago.ID),
				TransferenciaId:     0,
				ReferenciaBancaria:  "",
				PagoIntento: administraciondtos.ResponsePagoIntentoCC{
					Id:                   valuePago.ID,
					MediopagosId:         uint(valuePago.MediopagosID),
					InstallmentdetailsId: uint(valuePago.InstallmentdetailsID),
					ExternalId:           "",
					PaidAt:               valuePago.PaidAt,
					ReortAt:              valuePago.ReportAt,
					IsAvailable:          valuePago.IsAvailable,
					Amount:               valuePago.Amount,
					Valorcupon:           valuePago.Valorcupon,
					StateComent:          valuePago.StateComment,
					Barcode:              "",
					BarcodeUrl:           "",
					AvailableAt:          valuePago.AvailableAt,
					RevertedAt:           valuePago.RevertedAt,
					HolderName:           strings.ToUpper(valuePago.HolderName),
					HolderEmail:          valuePago.HolderEmail,
					HolderType:           valuePago.HolderType,
					HolderNumber:         valuePago.HolderNumber,
					HolderCbu:            "",
					TicketNumber:         valuePago.TicketNumber,
					AuthorizationCode:    valuePago.AuthorizationCode,
					CardLastFourDigits:   valuePago.CardLastFourDigits,
					TransactionId:        valuePago.TransactionID,
					SiteId:               "",
				},
			})
		}

		for Key, valuePagoTipoTemp := range pagosTipoTemporal {
			for _, valuePagoTemp := range pagostemporal {
				if valuePagoTipoTemp.Id == uint(valuePagoTemp.PagostipoID) {
					pagosTipoTemporal[Key].Pagos = append(pagosTipoTemporal[Key].Pagos, valuePagoTemp)
				}
			}
		}
		listaRevertidos.Cuenta.PagoTipo = append(listaRevertidos.Cuenta.PagoTipo, pagosTipoTemporal...)
		for _, valueCL := range cierreLoteDisputa {
			for key1, valuePI := range listaRevertidos.Cuenta.PagoTipo {
				for key, valuePago := range valuePI.Pagos {
					if valueCL.ExternalclienteID == valuePago.PagoIntento.TransactionId {
						// listaRevertidos.Cuenta.PagoTipo[key1].Pagos[key].PagoIntento.CierreLote = administraciondtos.ResponseCLCC(valueCL)
						listaRevertidos.Cuenta.PagoTipo[key1].Pagos[key].PagoIntento.CierreLote = administraciondtos.ResponseCLCC{
							Id:                         valueCL.Id,
							PagoestadoexternosId:       valueCL.PagoestadoexternosId,
							ChannelarancelesId:         valueCL.ChannelarancelesId,
							ImpuestosId:                valueCL.ImpuestosId,
							PrismamovimientodetallesId: valueCL.PrismamovimientodetallesId,
							PrismamovimientodetalleId:  0,
							PrismatrdospagosId:         valueCL.PrismatrdospagosId,
							BancoExternalId:            valueCL.BancoExternalId,
							Tiporegistro:               valueCL.Tiporegistro,
							PagosUuid:                  valueCL.PagosUuid,
							ExternalmediopagoId:        valueCL.ExternalmediopagoId,
							Nrotarjeta:                 valueCL.Nrotarjeta,
							Tipooperacion:              valueCL.Tipooperacion,
							Fechaoperacion:             valueCL.Fechaoperacion,
							Monto:                      valueCL.Monto,
							Montofinal:                 valueCL.Montofinal,
							Codigoautorizacion:         valueCL.Codigoautorizacion,
							Nroticket:                  valueCL.Nroticket,
							SiteID:                     valueCL.SiteID,
							ExternalloteId:             valueCL.ExternalloteId,
							Nrocuota:                   valueCL.Nrocuota,
							FechaCierre:                valueCL.FechaCierre,
							Nroestablecimiento:         valueCL.Nroestablecimiento,
							ExternalclienteID:          valueCL.ExternalclienteID,
							Nombrearchivolote:          valueCL.Nombrearchivolote,
							Match:                      valueCL.Match,
							FechaPago:                  valueCL.FechaPago,
							Disputa:                    valueCL.Disputa,
							Reversion:                  valueCL.Reversion,
							DetallemovimientoId:        valueCL.DetallemovimientoId,
							DetallepagoId:              valueCL.DetallepagoId,
							Descripcioncontracargo:     valueCL.Descripcioncontracargo,
							ExtbancoreversionId:        valueCL.ExtbancoreversionId,
							Conciliado:                 valueCL.Conciliado,
							Estadomovimiento:           valueCL.Estadomovimiento,
							Descripcionbanco:           valueCL.Descripcionbanco,
						}
					}
				}
			}
		}
	}

	return
}

func (s *service) PostPreferencesService(request administraciondtos.RequestPreferences) (erro error) {
	clienteId, err := strconv.Atoi(request.ClientId)
	if err != nil {
		erro = errors.New("cliente id no válido")
		return
	}

	preferenceEntity := entities.Preference{
		ClientesId:     uint(clienteId),
		Maincolor:      request.MainColor,
		Secondarycolor: request.SecondaryColor,
		Logo:           request.RutaLogo,
	}

	err = s.repository.PostPreferencesRepository(preferenceEntity)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	return
}

func (s *service) GetPagosDevService(filtro filtros.PagoFiltro) (response []entities.Pago, erro error) {

	response, _, erro = s.repository.GetPagos(filtro)

	if erro != nil {
		return
	}

	return

}

func (s *service) UpdatePagosDevService(response []entities.Pago) (pg []uint, erro error) {
	// return s.repository.UpdatePagosNotificados(listaPagosNotificar)

	for _, pago := range response {
		pg = append(pg, pago.ID)
	}

	erro = s.repository.UpdatePagosDev(pg)
	if erro != nil {
		return nil, erro
	}

	return pg, nil

}

func (s *service) BuildPagosMovDev(pagos []uint) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {

	// deben existir pagos
	if len(pagos) < 1 {
		erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE)
		return
	}

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		Channel:                 true,
		CargarPago:              true,
		CargarPagoTipo:          true,
		CargarCuenta:            true,
		CargarCliente:           true,
		CargarCuentaComision:    true,
		CargarImpuestos:         true,
		ExternalId:              true,
		CargarInstallmentdetail: true,
		PagosId:                 pagos,
	}
	// * 4 - Busco los pagos intentos que corresponden a los pagos
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)
	if erro != nil {
		return
	}

	// * 5 Obtener los pagos intentos
	for i := range pagosIntentos {
		movimientoCierreLote.ListaPagoIntentos = append(movimientoCierreLote.ListaPagoIntentos, pagosIntentos[i])
	}

	// * 6 - Busco el estado acreditado
	filtroPagoEstado := filtros.PagoEstadoFiltro{
		Nombre: config.MOVIMIENTO_ACCREDITED,
	}

	pagoEstadoAcreditado, erro := s.repository.GetPagoEstado(filtroPagoEstado)
	logs.Info(pagoEstadoAcreditado)

	if erro != nil {
		return
	}

	// var monto_pagado entities.Monto
	// * 8 - Modifico los pagos, creo los logs de los estados de pagos y creo los movimientos
	for i := range movimientoCierreLote.ListaPagoIntentos {
		/* * para el calculo de la comision fitrar por el id del channel y el id de la cuentar*/

		var pagoCuotas bool
		var examinarPagoCuota bool
		if movimientoCierreLote.ListaPagoIntentos[i].Installmentdetail.Cuota > 1 {
			pagoCuotas = true
			examinarPagoCuota = true
		}
		var idMedioPago uint
		if movimientoCierreLote.ListaPagoIntentos[i].MediopagosID == 30 {
			idMedioPago = uint(movimientoCierreLote.ListaPagoIntentos[i].MediopagosID)
			pagoCuotas = true
			examinarPagoCuota = true
		}

		filtroComisionChannel := filtros.CuentaComisionFiltro{
			CargarCuenta:      true,
			ChannelId:         uint(movimientoCierreLote.ListaPagoIntentos[i].Mediopagos.ChannelsID),
			CuentaId:          movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.ID,
			Mediopagoid:       idMedioPago,
			ExaminarPagoCuota: examinarPagoCuota,
			PagoCuota:         pagoCuotas,
			Channelarancel:    true,
			FechaPagoVigencia: movimientoCierreLote.ListaPagoIntentos[i].PaidAt,
		}

		logs.Info(filtroComisionChannel)

		cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
		listaCuentaComision := append([]entities.Cuentacomision{}, cuentaComision)

		// modificar la cuentacomision segun le channel id
		// listaCuentaComision := movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuentacomisions
		if len(listaCuentaComision) < 1 {
			erro = fmt.Errorf("no se pudo encontrar una comision para la cuenta %s del cliente %s", movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuenta, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Cliente)
			s.utilService.CreateNotificacionService(
				entities.Notificacione{
					Tipo:        entities.NotificacionCierreLote,
					Descripcion: erro.Error(),
				},
			)
			return
		}

		movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, movimientoCierreLote.ListaPagoIntentos[i].Pago)

		// * crear el log de estado de pago acreditado
		pagoEstadoLog := entities.Pagoestadologs{
			PagosID:       movimientoCierreLote.ListaPagoIntentos[i].PagosID,
			PagoestadosID: int64(pagoEstadoAcreditado.ID),
		}
		movimientoCierreLote.ListaPagosEstadoLogs = append(movimientoCierreLote.ListaPagosEstadoLogs, pagoEstadoLog)

		if movimientoCierreLote.ListaPagoIntentos[i].Pago.PagoestadosID == int64(pagoEstadoAcreditado.ID) {

			var importe entities.Monto
			importe = movimientoCierreLote.ListaPagoIntentos[i].Amount
			if movimientoCierreLote.ListaPagoIntentos[i].Valorcupon != 0 {
				importe = movimientoCierreLote.ListaPagoIntentos[i].Valorcupon
			}

			movimiento := entities.Movimiento{}
			// monto_pagado = movimientoCierreLote.ListaPagoIntentos[i].Amount
			// if movimientoCierreLote.ListaPagoIntentos[i].Valorcupon > 0 {
			// 	monto_pagado = movimientoCierreLote.ListaPagoIntentos[i].Valorcupon
			// }
			movimiento.AddCredito(uint64(movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.CuentasID), uint64(movimientoCierreLote.ListaPagoIntentos[i].ID), importe)

			s.utilService.BuildComisiones(&movimiento, &listaCuentaComision, movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Iva, importe)

			movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)
			movimientoCierreLote.ListaPagoIntentos[i].AvailableAt = movimientoCierreLote.ListaPagoIntentos[i].CreatedAt
		}

	}
	return
}

// ? Implementacion de servicio conusltar cierrelote para herramienta wee

func (s *service) GetConsultarClRapipagoService(filtro filtros.RequestClrapipago) (response administraciondtos.ResponseCLRapipago, erro error) {

	clrapiapgo, total, erro := s.repository.GetConsultarClRapipagoRepository(filtro)

	if erro != nil {
		return
	}

	if filtro.Number > 0 && filtro.Size > 0 {
		response.Meta = _setPaginacion(filtro.Number, filtro.Size, total)
	}

	for _, cl := range clrapiapgo {

		var detalles []administraciondtos.ClRapipagoDetalle
		for _, detalle := range cl.RapipagoDetalle {
			detalles = append(detalles, administraciondtos.ClRapipagoDetalle{
				FechaCobro:       detalle.Clearing,
				ImporteCobrado:   uint64(detalle.ImporteCobrado),
				ImporteCalculado: uint64(detalle.ImporteCalculado),
				CodigoBarras:     detalle.CodigoBarras,
				Conciliado:       detalle.Match,
			})
		}

		fechaProceso := s.commonsService.ConvertirFormatoFecha(cl.FechaProceso)
		r := administraciondtos.CLRapipago{
			IdArchivo:                cl.NombreArchivo,
			FechaProceso:             fechaProceso,
			Detalles:                 uint64(cl.CantDetalles),
			ImporteTotal:             uint64(cl.ImporteTotal),
			ImporteTotalCalculado:    uint64(cl.ImporteTotalCalculado),
			IdBanco:                  uint64(cl.BancoExternalId),
			FechaAcreditacion:        cl.Fechaacreditacion.Format("2006-01-02"),
			CantidadDiasAcreditacion: uint64(cl.Cantdias),
			ImporteMinimo:            uint64(cl.ImporteMinimo),
			Coeficiente:              cl.Coeficiente,
			EnObservacion:            cl.Enobservacion,
			DiferenciaBanco:          s.utilService.ToFixed(cl.Difbancocl, 2) / 100,
			FechaCreacion:            cl.CreatedAt.Format("2006-01-02"),
			ClRapipagoDetalle:        detalles,
		}

		response.ClRapipago = append(response.ClRapipago, r)
	}

	return
}

func (s *service) GetCaducarOfflineIntentos() (intentosCaducados int, erro error) {

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		CargarPago: true,
		// CargarPagoEstado:   true,
		// Channel:            true,
		PagoEstadoIdFiltro: 2, // Estado Processing
		ChannelIdFiltro:    3, // Channel Offline
	}
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)
	if erro != nil {
		return
	}

	var pagosActualizar []entities.Pago
	const idEstadoExpired = 6 // Estado Expired

	for _, pagoIntento := range pagosIntentos {
		lastFechaVencimiento := pagoIntento.Pago.FirstDueDate
		if pagoIntento.Pago.SecondDueDate.After(pagoIntento.Pago.FirstDueDate) {
			lastFechaVencimiento = pagoIntento.Pago.SecondDueDate
		}

		// Fecha tomada en cuenta para caducar (Fecha actual -5 días)
		fechaControl := time.Now().AddDate(0, 0, -5)
		if fechaControl.After(lastFechaVencimiento) {
			pagosActualizar = append(pagosActualizar, pagoIntento.Pago)
		}

	}

	if len(pagosActualizar) > 0 {
		erro = s.repository.UpdateEstadoPagos(pagosActualizar, uint64(idEstadoExpired))
		if erro != nil {
			return
		}
	}

	return len(pagosActualizar), nil
}

func (s *service) GetPagosCalculoMovTemporalesService(filtro filtros.PagoIntentoFiltros) (pagosid []uint, erro error) {

	if filtro.FechaPagoInicio.IsZero() {
		var fechaI time.Time
		var fechaF time.Time
		// si los filtros recibidos son ceros toman la fecha actual
		fechaI, fechaF, erro = s.commonsService.FormatFecha()
		if erro != nil {
			return
		}
		// a las fechas se le restan un dia ya sea por backgraund o endpoint
		filtro.FechaPagoInicio = fechaI.AddDate(0, 0, int(-1))
		filtro.FechaPagoFin = fechaF.AddDate(0, 0, int(-1))

	} else {
		filtro.FechaPagoInicio = filtro.FechaPagoInicio.AddDate(0, 0, int(-1))
		filtro.FechaPagoFin = filtro.FechaPagoFin.AddDate(0, 0, int(-1))
	}
	logs.Info(filtro.FechaPagoInicio)
	logs.Info(filtro.FechaPagoFin)

	pagos, erro := s.repository.GetPagosIntentosCalculoComisionRepository(filtro)

	// apilink, erro := s.repository.GetPagosApilink(filtro)
	// pagosid = append(pagosid, apilink...)

	// rapipago, erro := s.repository.GetPagosRapipago(filtro)
	// pagosid = append(pagosid, rapipago...)

	// prisma, erro := s.repository.GetPagosPrisma(filtro)
	// pagosid = append(pagosid, prisma...)

	if erro != nil {
		return
	}

	for _, pg := range pagos {
		pagosid = append(pagosid, uint(pg.PagosID))
	}

	return
}

func (s *service) GetPagosIntentosCalculoComisionRepository(filtro filtros.PagoIntentoFiltros) (pagos []entities.Pagointento, erro error) {
	pagos, erro = s.repository.GetPagosIntentosCalculoComisionRepository(filtro)
	if erro != nil {
		return
	}
	return
}

func (s *service) BuildPagosCalculoTemporales(pagos []uint) (movimientoCierreLote administraciondtos.MovimientoTemporalesResponse, erro error) {

	// deben existir pagos
	if len(pagos) < 1 {
		erro = fmt.Errorf(ERROR_LISTA_CIERRE_LOTE)
		return
	}

	filtroPagoIntento := filtros.PagoIntentoFiltro{
		Channel:                 true,
		CargarPago:              true,
		CargarPagoTipo:          true,
		CargarCuenta:            true,
		CargarCliente:           true,
		CargarCuentaComision:    true,
		CargarImpuestos:         true,
		CargarInstallmentdetail: true,
		PagoIntentoAprobado:     true,
		PagosId:                 pagos,
	}
	// * 1 - Busco los pagos intentos que corresponden a los pagos y cargar toda informacion necesaria para el calculo de comisiones
	pagosIntentos, erro := s.repository.GetPagosIntentos(filtroPagoIntento)
	if erro != nil {
		return
	}

	if len(pagos) != len(pagosIntentos) {
		erro = fmt.Errorf(ERROR_LISTA_PAGOS_INTENTOS)
		return
	}

	// buscar los clientes que sean sujeto de retencion y cargar sus retenciones
	filtroCLiente := filtros.ClienteFiltro{
		SujetoRetencion: true,
		CargarCuentas:   true,
	}

	clientes, _, err := s.repository.GetClientes(filtroCLiente)
	if err != nil {
		return
	}

	// var monto_pagado entities.Monto
	// * 8 - Modifico los pagos, creo los logs de los estados de pagos y creo los movimientos
	for i := range pagosIntentos {
		/* * para el calculo de la comision fitrar por el id del channel y el id de la cuentar*/

		var pagoCuotas bool
		var examinarPagoCuota bool
		if pagosIntentos[i].Installmentdetail.Cuota > 1 {
			pagoCuotas = true
			examinarPagoCuota = true
		}
		var idMedioPago uint
		if pagosIntentos[i].MediopagosID == 30 {
			idMedioPago = uint(pagosIntentos[i].MediopagosID)
			pagoCuotas = true
			examinarPagoCuota = true
		}

		filtroComisionChannel := filtros.CuentaComisionFiltro{
			CargarCuenta:      true,
			ChannelId:         uint(pagosIntentos[i].Mediopagos.ChannelsID),
			CuentaId:          pagosIntentos[i].Pago.PagosTipo.Cuenta.ID,
			Mediopagoid:       idMedioPago,
			ExaminarPagoCuota: examinarPagoCuota,
			PagoCuota:         pagoCuotas,
			Channelarancel:    true,
			FechaPagoVigencia: pagosIntentos[i].PaidAt,
		}

		logs.Info(filtroComisionChannel)

		cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
		listaCuentaComision := append([]entities.Cuentacomision{}, cuentaComision)

		// modificar la cuentacomision segun le channel id
		// listaCuentaComision := movimientoCierreLote.ListaPagoIntentos[i].Pago.PagosTipo.Cuenta.Cuentacomisions
		if len(listaCuentaComision) < 1 {
			erro = fmt.Errorf("no se pudo encontrar una comision para la cuenta %s del cliente %s", pagosIntentos[i].Pago.PagosTipo.Cuenta.Cuenta, pagosIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Cliente)
			s.utilService.CreateNotificacionService(
				entities.Notificacione{
					Tipo:        entities.NotificacionCierreLote,
					Descripcion: erro.Error(),
				},
			)
			return
		}

		movimientoCierreLote.ListaPagosCalculado = pagos

		var importe entities.Monto
		importe = pagosIntentos[i].Amount

		// if pagosIntentos[i].Valorcupon != 0 {
		// 	importe = pagosIntentos[i].Valorcupon
		// }

		movimiento := entities.Movimientotemporale{}
		// monto_pagado = movimientoCierreLote.ListaPagoIntentos[i].Amount
		// if movimientoCierreLote.ListaPagoIntentos[i].Valorcupon > 0 {
		// 	monto_pagado = movimientoCierreLote.ListaPagoIntentos[i].Valorcupon
		// }
		movimiento.AddCredito(uint64(pagosIntentos[i].Pago.PagosTipo.CuentasID), uint64(pagosIntentos[i].ID), importe)

		s.utilService.BuildComisionesTemporales(&movimiento, &listaCuentaComision, pagosIntentos[i].Pago.PagosTipo.Cuenta.Cliente.Iva, importe)

		s.BuildRetencionesTemporales(&movimiento, importe, pagosIntentos[i], clientes)

		movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)

	}
	return
}

func (s *service) ConciliacionPagosReportesService(filtro filtros.PagoFiltro) (valoresNoEncontrados []string, erro error) {

	successPayments, erro := s.repository.GetSuccessPaymentsRepository(filtro)

	if erro != nil {
		erro = fmt.Errorf("error en consultar pagos exitosos para la conciliacion")
		return
	}

	// si no hay resultados en la busqueda de pagos exitosos, no se debe continuar
	if len(successPayments) == 0 {
		erro = fmt.Errorf("no existen pagos exitosos para ser conciliados")
		return
	}

	// necesario convertir fecha para consultar en tabla reportes
	filtro.Fecha[0] = s.commonsService.ConvertirFechaToDDMMYYYY(filtro.Fecha[0])

	reporte, erro := s.repository.GetReportesPagoRepository(filtro)

	if erro != nil {
		erro = fmt.Errorf("error en consultar reportes de pagos enviados para la conciliacion: %s", erro.Error())
		return
	}

	// si no hay resultados en la busqueda de reportes de pagos, no se debe continuar
	if reporte.ID == 0 {
		erro = fmt.Errorf("no existen reportes de pagos para ser conciliados")
		return
	}

	// comparar pagos exitosos con reportes
	valoresNoEncontrados = _conciliarPagosYReportes(successPayments, reporte)

	// si hay valores no encontrados o conciliados, se debe loguear y notificar email
	if len(valoresNoEncontrados) > 0 {

		// hacemos un log de ls valores no encontrados
		erro = fmt.Errorf("algunos pagos no se encontraron en la conciliacion del cliente %s: %v", reporte.Cliente, valoresNoEncontrados)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       erro.Error(),
			Funcionalidad: "ConciliacionPagosReportesService",
		}

		err := s.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), erro.Error())
			logs.Error(mensaje)
		}

		var email = []string{config.EMAIL_TELCO}
		fafafa := _armarMensajeVariable(valoresNoEncontrados)

		filtro := utildtos.RequestDatosMail{
			Email:            email,
			Asunto:           "Conciliacion Pagos con Reporte de Cobranzas",
			From:             "Wee.ar!",
			Nombre:           "Administrador",
			Mensaje:          "Las siguientes referencias de pagos para el cliente " + reporte.Cliente + " no se conciliaron: " + fafafa,
			CamposReemplazar: valoresNoEncontrados,
			AdjuntarEstado:   false,
			TipoEmail:        "template",
		}
		erro = s.utilService.EnviarMailService(filtro)
		logs.Info(erro)
	} // Fin de if se encuentran valores no conciliados

	/*  Si el cliente es DPEC, controlar los montos */
	filtroCuenta := filtros.CuentaFiltro{
		Id: uint(filtro.CuentaId),
	}
	cuenta, erro := s.repository.GetCuenta(filtroCuenta)
	var montosIguales bool
	if commons.ContainStrings([]string{cuenta.Cliente.Cliente}, "dpec") {
		// conciliar el total de monto de pago con el total cobrado del reporte
		montosIguales, erro = _conciliarByMontos(successPayments, reporte.Totalcobrado)

		if !montosIguales {

			erro = fmt.Errorf("los montos de pagos exitosos no coinciden con el total cobrado reportado")

			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       erro.Error(),
				Funcionalidad: "ConciliacionPagosReportesService",
			}

			err := s.utilService.CreateLogService(log)

			if err != nil {
				mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), erro.Error())
				logs.Error(mensaje)
			}
		}

	}
	/*  Fin del caso cliente DPEC */

	return
}

// compara las external references de los pagos exitosos de un determiando dia y cuenta, con los reportes de pagos enviados a ese cliente en ese dia
func _conciliarPagosYReportes(pagos []entities.Pago, reporte entities.Reporte) (valoresNoEncontrados []string) {

	var detalleExternalsIds []string
	externalReferences := _filtrarPorExternalReference(pagos)

	for _, detalle := range reporte.Reportedetalle {
		// guardar en un array para comparar de manera inversa los external con los detalles
		detalleExternalsIds = append(detalleExternalsIds, detalle.PagosId)

		// el pago_id del Reportedetalle contiene el valor del external reference de cada pago exitoso informado por el reporte
		if !commons.ContainStrings(externalReferences, detalle.PagosId) {
			valoresNoEncontrados = append(valoresNoEncontrados, detalle.PagosId)
		}

	} // end for detalle := range reporte.Reportedetalle

	for _, er := range externalReferences {

		if !commons.ContainStrings(detalleExternalsIds, er) {
			valoresNoEncontrados = append(valoresNoEncontrados, er)
		}

	} // end for er := range externalReferences

	return
}

func _filtrarPorExternalReference(pagos []entities.Pago) (externalReferences []string) {
	for _, pago := range pagos {
		externalReferences = append(externalReferences, pago.ExternalReference)
	}
	return
}

func _conciliarByMontos(pagos []entities.Pago, monto string) (result bool, erro error) {
	var amount entities.Monto
	// recorrer los pago
	for _, pago := range pagos {
		// recorrer los PagoIntentos
		for _, intento := range pago.PagoIntentos {
			amount += intento.Amount
		}
	} // fin de for de pagos
	// El monto del reporte esta en string, luego para comparar
	montoPagosString := util.Resolve().FormatNum(amount.Float64())

	if montoPagosString == monto {
		result = true
	}
	return
}

// recibe un array de string que se deben ocupar como campos a reemplazar en el mensaje del email
func _armarMensajeVariable(p_arrayString []string) (mensaje string) {
	mensaje = "<br>"
	for i := 0; i < len(p_arrayString); i++ {
		mensaje += fmt.Sprintf("<b>#%d</b><br>", i)
	}

	return
}

// Metodos Retenciones
func (s *service) GetRetencionesService(request retenciondtos.RentencionRequestDTO, getDTO bool) (response retenciondtos.RentencionesResponseDTO, erro error) {

	retenciones, count, erro := s.repository.GetRetencionesRepository(request)

	if erro != nil {
		return
	}
	if getDTO {
		var retencionesDTO []retenciondtos.RentencionResponseDTO
		// pasar entidades a DTO response
		for _, retencion := range retenciones {
			var tempRetencionResponseDTO retenciondtos.RentencionResponseDTO
			tempRetencionResponseDTO.FromEntity(retencion)
			retencionesDTO = append(retencionesDTO, tempRetencionResponseDTO)
		}
		response.RetencionesDTO = retencionesDTO
	}

	if !getDTO {
		response.Retenciones = retenciones
	}

	response.Count = count

	return
}

func (s *service) GetClienteRetencionesService(request retenciondtos.RentencionRequestDTO) (response retenciondtos.RentencionesResponseDTO, erro error) {
	var (
		cliente_retenciones []entities.ClienteRetencion
		count               int64
	)
	// validar datos de la requets
	erro = request.Validar()
	if erro != nil {
		return
	}

	cliente_retenciones, count, erro = s.repository.GetClienteRetencionesRepository(request)

	if erro != nil {
		return
	}

	var retencionesDTO []retenciondtos.RentencionResponseDTO
	// pasar entidades a DTO response
	for _, cr := range cliente_retenciones {
		var tempRetencionResponseDTO retenciondtos.RentencionResponseDTO
		tempRetencionResponseDTO.FromEntity(cr.Retencion)
		for _, cert := range cr.Certificados {
			var tempCertificadoResponseDTO retenciondtos.CertificadoResponseDTO
			tempCertificadoResponseDTO.FromEntity(cert)
			tempRetencionResponseDTO.Certificados = append(tempRetencionResponseDTO.Certificados, tempCertificadoResponseDTO)
		}
		retencionesDTO = append(retencionesDTO, tempRetencionResponseDTO)
	}

	response.RetencionesDTO = retencionesDTO
	response.Count = count

	return
}

func (s *service) GetClienteUnlinkedRetencionesService(request retenciondtos.RentencionRequestDTO) (response retenciondtos.RentencionesResponseDTO, erro error) {
	var (
		retenciones []entities.Retencion
		count       int64
	)
	// validar datos de la requets
	erro = request.Validar()
	if erro != nil {
		return
	}

	retenciones, count, erro = s.repository.GetClienteUnlinkedRetencionesRepository(request)

	if erro != nil {
		return
	}

	var retencionesDTO []retenciondtos.RentencionResponseDTO
	// pasar entidades a DTO response
	for _, retencion := range retenciones {
		var tempRetencionResponseDTO retenciondtos.RentencionResponseDTO
		tempRetencionResponseDTO.FromEntity(retencion)
		retencionesDTO = append(retencionesDTO, tempRetencionResponseDTO)
	}
	response.RetencionesDTO = retencionesDTO
	response.Count = count
	return
}

func (s *service) CreateClienteRetencionService(request retenciondtos.RentencionRequestDTO) (erro error) {

	// validar datos de la requets
	erro = request.ValidarPost()
	if erro != nil {
		return
	}

	erro = s.repository.CreateClienteRetencionRepository(request)

	if erro != nil {
		return
	}

	// response.FromEntity(cliente)

	return
}

func (s *service) GetClienteRetencionService(retencion_id uint, cliente_id uint) (cliente_retencion entities.ClienteRetencion, erro error) {

	cliente_retencion, erro = s.repository.GetClienteRetencionRepository(retencion_id, cliente_id)

	if erro != nil {
		return
	}

	return
}

func (s *service) BuildRetenciones(movimiento *entities.Movimiento, importe entities.Monto, pagointento entities.Pagointento, clientes []entities.Cliente) (erro error) {
	var (
		retenciones          []entities.Retencion
		CuentaId             uint64
		SujetoClienteRetener entities.Cliente
	)

	// Si es reversion se controla que el campo reverted_at sea posterior o igual a la fecha de inicio de retenciones.
	// Si es anterior, el pagointento no tiene movimiento_retencion para ser revertido, y se sale de la funcion sin calcular ninguna retencion
	if !pagointento.RevertedAt.IsZero() {
		if !pagointento.RevertirMovimientoRetencion() {
			return
		}
	}

	// se obtiene el id de la cuenta a partir del movimiento
	CuentaId = movimiento.CuentasId

	// se filtra el slice de Clientes buscando aquel a cual corresponde la CuentaId
	for _, cliente := range clientes {
		if len(*cliente.Cuentas) > 0 {
			for _, cuenta := range *cliente.Cuentas {
				if uint64(cuenta.ID) == CuentaId {
					SujetoClienteRetener = cliente
					break
				}
			}
		}
	}

	// detectar si el cliente es SujetoClienteRetener. Checkear campo sujeto_retencion
	// si no estas sujeto a retencion, retornar sin error para no calcular retencion alguna.
	if !SujetoClienteRetener.SujetoRetencion {
		return
	}

	// obtener las retenciones del cliente asociado al movimiento
	retenciones = SujetoClienteRetener.Retenciones

	// dado las combinaciones posibles de retenciones por canales de pago
	// tomar solo las que coincidan con el channel_id del pago asociado al movimiento en cuestion
	var retencionesVigentesFilter, retencionesChannelIdFilter []entities.Retencion

	// de las retenciones se toman las vigentes
	for _, r := range retenciones {
		if r.IsCurrent(pagointento.PaidAt) {
			retencionesVigentesFilter = append(retencionesVigentesFilter, r)
		}
	}

	// de las retenciones vigentes se toma la que coincide con el channel id del pagointento
	for _, r := range retencionesVigentesFilter {
		if r.ChannelsId == pagointento.Mediopagos.Channel.ID {
			retencionesChannelIdFilter = append(retencionesChannelIdFilter, r)
		}
	}

	var movimiento_retencions []entities.MovimientoRetencion
	// la estrategia para el calculo de cada tipo de retencion
	var estrategia = StrategyRetencion{}

	// si retencionesChannelIdFilter esta vacio significa que ninguna retencion asignada al cliente coincide con el channel id del pagointento recibido como parametro en BuildRetenciones
	if len(retencionesChannelIdFilter) > 0 {
		for _, retencion := range retencionesChannelIdFilter {
			gravamen := retencion.Condicion.Gravamen.Gravamen
			// decidir que tipo de retencion crear
			iRetencion, erro := RetencionFactory(gravamen, retencion)

			if erro != nil {
				s.utilService.BuildLog(erro, "BuildRetenciones")
				break
			}
			estrategia.setStrategy(iRetencion)
			result := estrategia.execStrategy(importe)

			// crear objeto MovimientoRetencion y guardar en array
			MovRet := entities.MovimientoRetencion{
				MovimientoId:    uint64(movimiento.ID),
				RetencionId:     retencion.ID,
				ClienteId:       uint64(SujetoClienteRetener.ID),
				Monto:           importe,
				ImporteRetenido: entities.Monto(result),
				Efectuada:       true,
			}

			// Restar del monto del movimiento, el importe de la retencion
			neto := movimiento.Monto - entities.Monto(result)
			// Asignar al Monto del Movimiento el valor calculado, neto de retencion
			movimiento.Monto = neto
			movimiento_retencions = append(movimiento_retencions, MovRet)
		}
	}

	// guardar las retenciones en el movimiento como un array
	movimiento.Movimientoretencions = movimiento_retencions

	return
}

func (s *service) BuildRetencionesTemporales(movimiento *entities.Movimientotemporale, importe entities.Monto, pagointento entities.Pagointento, clientes []entities.Cliente) (erro error) {
	var (
		retenciones          []entities.Retencion
		CuentaId             uint64
		SujetoClienteRetener entities.Cliente
	)

	// Si es reversion se controla que el campo reverted_at sea posterior o igual a la fecha de inicio de retenciones.
	// Si es anterior, el pagointento no tiene movimiento_retencion para ser revertido, y se sale de la funcion sin calcular ninguna retencion
	if !pagointento.RevertedAt.IsZero() {
		if !pagointento.RevertirMovimientoRetencion() {
			return
		}
	}

	// se obtiene el id de la cuenta a partir del movimiento
	CuentaId = movimiento.CuentasId

	// se filtra el slice de Clientes buscando aquel a cual corresponde la CuentaId
	for _, cliente := range clientes {
		if len(*cliente.Cuentas) > 0 {
			for _, cuenta := range *cliente.Cuentas {
				if uint64(cuenta.ID) == CuentaId {
					SujetoClienteRetener = cliente
					break
				}
			}
		}
	}

	// detectar si el cliente es SujetoClienteRetener. Checkear campo sujeto_retencion
	// si no estas sujeto a retencion, retornar sin error para no calcular retencion alguna.
	if !SujetoClienteRetener.SujetoRetencion {
		return
	}

	// obtener las retenciones del cliente asociado al movimiento
	retenciones = SujetoClienteRetener.Retenciones

	// dado las combinaciones posibles de retenciones por canales de pago
	// tomar solo las que coincidan con el channel_id del pago asociado al movimiento en cuestion
	var retencionesVigentesFilter, retencionesChannelIdFilter []entities.Retencion

	// de las retenciones se toman las vigentes
	for _, r := range retenciones {
		if r.IsCurrent(pagointento.PaidAt) {
			retencionesVigentesFilter = append(retencionesVigentesFilter, r)
		}
	}

	// de las retenciones vigentes se toma la que coincide con el channel id del pagointento
	for _, r := range retencionesVigentesFilter {
		if r.ChannelsId == pagointento.Mediopagos.Channel.ID {
			retencionesChannelIdFilter = append(retencionesChannelIdFilter, r)
		}
	}

	var movimiento_retencions []entities.MovimientoRetenciontemporale
	// la estrategia para el calculo de cada tipo de retencion
	var estrategia = StrategyRetencion{}

	// si retencionesChannelIdFilter esta vacio significa que ninguna retencion asignada al cliente coincide con el channel id del pagointento recibido como parametro en BuildRetenciones
	if len(retencionesChannelIdFilter) > 0 {
		for _, retencion := range retencionesChannelIdFilter {
			gravamen := retencion.Condicion.Gravamen.Gravamen
			// decidir que tipo de retencion crear
			iRetencion, erro := RetencionFactory(gravamen, retencion)

			if erro != nil {
				s.utilService.BuildLog(erro, "BuildRetencionesTemporales")
				break
			}
			estrategia.setStrategy(iRetencion)
			result := estrategia.execStrategy(importe)

			// crear objeto MovimientoRetencion y guardar en array
			MovRet := entities.MovimientoRetenciontemporale{
				MovimientotemporalesID: uint64(movimiento.ID),
				RetencionId:            retencion.ID,
				ClienteId:              uint64(SujetoClienteRetener.ID),
				Monto:                  importe,
				ImporteRetenido:        entities.Monto(result),
				Efectuada:              true,
			}

			// Restar del monto del movimiento, el importe de la retencion
			neto := movimiento.Monto - entities.Monto(result)
			// Asignar al Monto del Movimiento el valor calculado, neto de retencion
			movimiento.Monto = neto
			movimiento_retencions = append(movimiento_retencions, MovRet)
		}
	}

	// guardar las retenciones en el movimiento como un array
	movimiento.Movimientoretenciontemporales = movimiento_retencions

	return
}

func (s *service) ComprobarMinimoRetencion(retencion entities.Retencion, clienteId uint) (result bool, monto entities.Monto, erro error) {

	cliente_retencion, erro := s.GetClienteRetencionService(retencion.ID, clienteId)

	if erro != nil {
		return
	}

	if cliente_retencion.ID == 0 {
		erro = errors.New("la retencion que busca comprobar no esta asignada a ningun cliente")
		return
	}

	FechaInicio, FechaFin, erro := s.commonsService.GetFechaInicioActualMes()

	if erro != nil {
		return
	}

	// este filtro permite buscar retenciones por cliente, periodo y gravamen
	filtro := retenciondtos.RentencionRequestDTO{
		ClienteId:   clienteId,                       // el cliente del cual se busca las retenciones
		FechaInicio: FechaInicio,                     // fecha inicio del mes actual. Dia uno del mes
		FechaFin:    FechaFin,                        // dia actual del mes actual
		GravamensId: retencion.Condicion.GravamensId, // el impuesto que se busca evaluar
	}
	retencionAgrupada, erro := s.repository.GetCalcularRetencionesRepository(filtro)
	if erro != nil {
		return
	}

	if len(retencionAgrupada) != 1 {
		erro := errors.New("no se puede comprabar minimo de retencion. la agrupacion por gravamen no arroja un solo resultado")
		s.utilService.BuildLog(erro, "ComprobarMinimoRetencion")
	}

	if len(retencionAgrupada) > 0 {
		monto = retencionAgrupada[0].TotalRetencion
	}

	minimo := retencion.MontoMinimo
	if monto.Float64() >= minimo {
		result = true
	}

	return
}

func (s *service) PostRetencionesCertificadosService(request retenciondtos.RetencionCertificadoRequestDTO) (erro error) {

	fechasToTime := []string{request.Fecha_Presentacion, request.Fecha_Caducidad}
	fechasProcesadas, erro := procesarFechas(fechasToTime)

	if erro != nil {
		return
	}

	certificado := entities.Certificado{
		ClienteRetencionsId: request.ClienteRetencionId,
		Fecha_Presentacion:  fechasProcesadas[0],
		Fecha_Caducidad:     fechasProcesadas[1],
		Ruta_file:           request.RutaFile,
	}

	erro = s.repository.PostRetencionCertificadoRepository(certificado)

	return
}

func (s *service) ValidarRetencion(request retenciondtos.RetencionCertificadoRequestDTO) (cliente_retencions_Id uint, cliente_name string, erro error) {

	retencionDTO := retenciondtos.RentencionRequestDTO{
		ClienteId:   request.ClienteId,
		RetencionId: request.RetencionId,
	}

	clienteretencion, erro := s.repository.GetRetencionClienteRepository(retencionDTO)

	if erro != nil {
		return
	}

	cliente_name = clienteretencion.Cliente.Cliente
	cliente_retencions_Id = clienteretencion.ID

	return
}

func (s *service) GetCalcularRetencionesService(request retenciondtos.RentencionRequestDTO) (retencionesAgrupadas []retenciondtos.RetencionAgrupada, erro error) {
	erro = request.ValidarFechas()
	if erro != nil {
		return
	}
	erro = request.Validar()
	if erro != nil {
		return
	}
	// retencionesAgrupadas, erro = s.repository.GetCalcularRetencionesRepository(request)
	retencionesAgrupadas, erro = s.repository.GetCalcularRetencionesByTransferenciasRepository(request)

	if erro != nil {
		return
	}

	if len(retencionesAgrupadas) > 0 {
		for i := 0; i < len(retencionesAgrupadas); i++ {
			retencionesAgrupadas[i].FechaInicio = request.FechaInicio
			retencionesAgrupadas[i].FechaFin = request.FechaFin
		}
	}

	return
}

func procesarFechas(fechasIn []string) (fechasOut []time.Time, erro error) {

	format := "2006-01-02"
	for _, fecha := range fechasIn {
		fechaTime, err := time.Parse(format, fecha)
		if err != nil {
			erro = err
			return
		}
		fechasOut = append(fechasOut, fechaTime)
	}

	return
}

func (s *service) GetCertificadoService(certificadoId uint) (certificado entities.Certificado, erro error) {

	certificado, erro = s.repository.GetCertificadoRepository(certificadoId)

	return
}

func (s *service) GetCertificadoCloudService(ctx context.Context, nombreFile string) (err error) {

	splits := strings.Split(nombreFile, "/")
	directorio := strings.Join(splits[0:(len(splits)-1)], "/")
	// nombreFileOnly := splits[len(splits)]

	err = s.store.GetObjectSpecific(config.AWS_BUCKET, nombreFile, directorio)
	return err
}

func (s *service) LeerContenidoDirectorio(datos entities.Certificado) (file retenciondtos.CertificadoFileDTO, erro error) {

	file.Id = datos.ID
	file.ClienteRetencionId = datos.ClienteRetencionsId

	split := strings.Split(datos.Ruta_file, "/")
	filenameServidor := split[len(split)-1]
	// split2 := strings.Split(filenameServidor, "-")
	// file.FileName = split2[len(split2)-1]
	file.FileName = filenameServidor

	archivo := config.DIR_BASE + config.DIR_CERT + datos.Ruta_file

	bytes, err := ioutil.ReadFile(archivo)
	if err != nil {
		erro = errors.New("error al leer contenido del archivo - " + err.Error())
		return
	}
	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)
	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "text/plain":
		base64Encoding += "data:text/plain;base64,"
	default:
		base64Encoding += "data:application/octet-stream;base64,"
	}
	// Append the base64 encoded output
	base64Encoding += s.commonsService.ToBase64(bytes)
	file.Content = base64Encoding

	split = strings.Split(archivo, "/")
	rutaArchivos := strings.Join(split[0:(len(split)-1)], "/")
	nombreArchivo := split[len(split)-1]

	/* antes de borrar el archivo se verifica si:
	el archivo fue leido, si fue movido y si la informacion que contiene el archivo se inserto en la db */
	erro = s.commonsService.BorrarArchivo(rutaArchivos, nombreArchivo)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarArchivos",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}

	erro = s.commonsService.BorrarDirectorio(rutaArchivos)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}

	return
}

func (s *service) LeerContenidoComprobanteRetencion(datos retenciondtos.ComprobanteResponseDTO) (file retenciondtos.ComprobanteFileDTO, erro error) {

	file.Id = datos.Id
	file.ClienteId = datos.ClienteId

	split := strings.Split(datos.RutaFile, "/")
	file.FileName = split[len(split)-1]

	archivo := config.DIR_BASE + config.DIR_COMP_RETENCIONES + datos.RutaFile

	bytes, err := ioutil.ReadFile(archivo)
	if err != nil {
		erro = errors.New("error al leer contenido del archivo - " + err.Error())
		return
	}
	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)
	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "text/plain":
		base64Encoding += "data:text/plain;base64,"
	default:
		base64Encoding += "data:application/octet-stream;base64,"
	}
	// Append the base64 encoded output
	base64Encoding += s.commonsService.ToBase64(bytes)
	file.Content = base64Encoding

	split = strings.Split(archivo, "/")
	rutaArchivos := strings.Join(split[0:(len(split)-1)], "/")
	nombreArchivo := split[len(split)-1]

	/* antes de borrar el archivo se verifica si:
	el archivo fue leido, si fue movido y si la informacion que contiene el archivo se inserto en la db */
	erro = s.commonsService.BorrarArchivo(rutaArchivos, nombreArchivo)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarArchivos",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}

	erro = s.commonsService.BorrarDirectorio(rutaArchivos)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}

	return
}

func (s *service) LeerContenidoReporteRendicionMensual(request retenciondtos.RentencionRequestDTO) (file retenciondtos.ComprobanteFileDTO, erro error) {

	file.Id = request.ReporteId
	file.ClienteId = request.ClienteId

	split := strings.Split(request.RutaFile, "/")
	file.FileName = split[len(split)-1]

	archivo := config.DIR_BASE + config.DIR_COMP_RETENCIONES + request.RutaFile

	bytes, err := ioutil.ReadFile(archivo)
	if err != nil {
		erro = errors.New("error al leer contenido del archivo - " + err.Error())
		return
	}
	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)
	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "text/plain":
		base64Encoding += "data:text/plain;base64,"
	default:
		base64Encoding += "data:application/octet-stream;base64,"
	}
	// Append the base64 encoded output
	base64Encoding += s.commonsService.ToBase64(bytes)
	file.Content = base64Encoding

	split = strings.Split(archivo, "/")
	rutaArchivos := strings.Join(split[0:(len(split)-1)], "/")
	nombreArchivo := split[len(split)-1]

	/* antes de borrar el archivo se verifica si:
	el archivo fue leido, si fue movido y si la informacion que contiene el archivo se inserto en la db */
	erro = s.commonsService.BorrarArchivo(rutaArchivos, nombreArchivo)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarArchivos",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}

	erro = s.commonsService.BorrarDirectorio(rutaArchivos)
	if erro != nil {
		logs.Error(erro.Error())
		log := entities.Log{
			Tipo:          entities.EnumLog("Error"),
			Funcionalidad: "BorrarDirectorio",
			Mensaje:       erro.Error(),
		}
		erro = s.utilService.CreateLogService(log)
		if erro != nil {
			logs.Error("error: al crear logs: " + erro.Error())
			return
		}
	}

	return
}

func (s *service) GetCondicionesService(request retenciondtos.RentencionRequestDTO) (response []retenciondtos.CondicionResponseDTO, erro error) {

	condiciones, erro := s.repository.GetCondicionesRepository(request)

	if erro != nil {
		return
	}

	// pasar entidades a DTO response
	for _, condicion := range condiciones {
		var tempcondicionResponseDTO retenciondtos.CondicionResponseDTO
		tempcondicionResponseDTO.FromEntity(condicion)
		response = append(response, tempcondicionResponseDTO)
	}

	return
}

func (s *service) CreateRetencionService(request retenciondtos.PostRentencionRequestDTO, isUpdate bool) (response []entities.Retencion, erro error) {

	erro = request.ValidarUpSert(isUpdate)
	if erro != nil {
		return
	}

	retenciones := request.ToEntitiesByChannel(isUpdate)

	for _, retencion := range retenciones {
		var created_retencion entities.Retencion
		created_retencion, erro = s.repository.CreateRetencionRepository(retencion)
		if erro != nil {
			return
		}
		response = append(response, created_retencion)
	}

	if erro != nil {
		return
	}

	return
}

func (s *service) UpdateRetencionService(request retenciondtos.PostRentencionRequestDTO, isUpdate bool) (response entities.Retencion, erro error) {

	erro = request.ValidarUpSert(isUpdate)
	if erro != nil {
		return
	}

	retencion := request.ToEntity(isUpdate, request.ChannelsId[0])

	response, erro = s.repository.UpdateRetencionRepository(retencion)
	if erro != nil {
		return
	}

	return
}

func (s *service) DeleteClienteRetencionService(request retenciondtos.RentencionRequestDTO) (erro error) {

	erro = request.ValidarDelete()
	if erro != nil {
		return
	}

	erro = s.repository.DeleteClienteRetencionRepository(request)
	if erro != nil {
		return
	}

	return
}

func (s *service) UpSertCondicionService(request retenciondtos.CondicionRequestDTO, isUpdate bool) (erro error) {

	erro = request.ValidarUpSert(isUpdate)
	if erro != nil {
		return
	}
	// request to entitty condicion
	condicion := request.ToEntity(isUpdate)

	if !isUpdate {
		erro = s.repository.CreateCondicionRepository(condicion)
	}

	if isUpdate {
		erro = s.repository.UpdateCondicionRepository(condicion)
	}

	if erro != nil {
		return
	}

	return
}

func (s *service) GetGravamenesService(filtro retenciondtos.GravamenRequestDTO) (response []retenciondtos.GravamenResponseDTO, erro error) {

	gravamenes, erro := s.repository.GetGravamensRepository(filtro)

	if erro != nil {
		return
	}

	// pasar entidades a DTO response
	for _, gravamen := range gravamenes {
		var tempGravamenResponseDTO retenciondtos.GravamenResponseDTO
		tempGravamenResponseDTO.FromEntity(gravamen)
		response = append(response, tempGravamenResponseDTO)
	}

	return
}

func (s *service) EvaluarRetencionesByClienteService(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {

	erro = request.ValidarFechas()
	if erro != nil {
		return
	}
	erro = request.Validar()
	if erro != nil {
		return
	}

	resultado, erro = s.repository.EvaluarRetencionesByClienteRepository(request)
	if erro != nil {
		return
	}
	var filtro retenciondtos.GravamenRequestDTO
	gravamenes, erro := s.repository.GetGravamensRepository(filtro)
	if erro != nil {
		return
	}
	var gravamen_names []string
	if len(gravamenes) > 0 {
		for _, g := range gravamenes {
			gravamen_names = append(gravamen_names, g.Gravamen)
		}
	}

	// Inicializar el mapa vacío
	gravamensMap := make(map[string]uint64)

	// Recorrer el array y agregar cada elemento al mapa con valor nulo
	for _, name := range gravamen_names {
		gravamensMap[name] = 0
	}

	// Evaluar si corresponde retener segun el monto acumulado or gravamen comparado con el minimo por tipo de gravamen
	CompararMontoMinimo(resultado, gravamensMap)

	return
}

func (s *service) EvaluarRetencionesByMovimientoService(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {
	resultado, erro = s.repository.EvaluarRetencionesByMovimientosRepository(request)
	return
}

func (s *service) GenerarCertificacionService(request retenciondtos.RentencionRequestDTO) (erro error) {

	erro = request.ValidarFechas()
	if erro != nil {
		return
	}

	erro = request.Validar()
	if erro != nil {
		return
	}

	// solicitar retenciones por cliente
	// Existen tanto objetos RetencionAgrupada como combinaciones posibles de gravamen y codigo de regimen tenga asociado el sujeto de retencion
	// El atributo Retener es false en este punto
	retencionesAgrupadas, erro := s.GetCalcularRetencionesService(request)
	// capturar el error y crear log
	if erro != nil {
		s.utilService.BuildLog(erro, "GenerarCertificacionService. Error al obtener retenciones agrupadas")
		return
	}

	// obtener una lista de movimientos ids de la tabla que relaciona movimientos y retenciones, segun datos del request filtro
	// movimientos_ids, erro := s.GetMovimientosRetencionesService(request) // obtener movimientos solamente por fecha
	movimientos_ids, erro := s.repository.GetMovimientosIdsCalculoRetencionComprobante(request) // obtener movimientos desde transferencias
	if erro != nil {
		s.utilService.BuildLog(erro, "GenerarCertificacionService. Error al obtener ids de movimientos con retenciones")
		return
	}

	if len(movimientos_ids) == 0 {
		erro = errors.New("no se encontraron movimientos asociados a retenciones para generar certificado")
		s.utilService.BuildLog(erro, "GenerarCertificacionService")
		return
	}

	// // Dado una lista de movimientos_id, obtener la suma de los amounts (importe o monto base sujeto a retencion)
	// totalAmount, erro := s.repository.GetTotalAmountByMovimientoIdsRepository(movimientos_ids)

	// if erro != nil {
	// 	s.utilService.BuildLog(erro, "GenerarCertificacionService")
	// 	return
	// }

	var filtro retenciondtos.GravamenRequestDTO
	gravamenes, erro := s.repository.GetGravamensRepository(filtro)
	if erro != nil {
		s.utilService.BuildLog(erro, "GenerarCertificacionService")
		return
	}
	var gravamen_names []string
	if len(gravamenes) > 0 {
		for _, g := range gravamenes {
			gravamen_names = append(gravamen_names, g.Gravamen)
		}
	}

	// Inicializar el mapa vacío
	gravamensMap := make(map[string]uint64)

	// Recorrer el array y agregar cada elemento al mapa con valor nulo
	for _, name := range gravamen_names {
		gravamensMap[name] = 0
	}

	// determinar si el importe retenido de cada RetencionAgrupada supera el minimo por gravamen
	retencionesComprobadas, _ := CompararMontoMinimo(retencionesAgrupadas, gravamensMap)

	var detalles []entities.ComprobanteDetalle
	// cerar N comprobante detalles
	for _, ra := range retencionesComprobadas {
		comprobante_detalle := entities.ComprobanteDetalle{
			TotalRetencion: ra.TotalRetencion,
			TotalMonto:     ra.TotalMonto,
			CodigoRegimen:  ra.CodigoRegimen,
			Gravamen:       ra.Gravamen,
			Retener:        ra.Retener,
		}
		detalles = append(detalles, comprobante_detalle)
	}

	filtroCliente := filtros.ClienteFiltro{
		Id: request.ClienteId,
	}

	// datos del cliente para el comprobante de retencion
	cliente, erro := s.repository.GetCliente(filtroCliente)
	if erro != nil {
		return
	}

	// parsear los ids de los movimientos con retencion, a un string separado por coma
	movsidsString := s.commonsService.NumberSliceToString(movimientos_ids)

	// AgruparPorGravamen y crear un comprobante por gravamen
	var comprobantes []entities.Comprobante
	for _, name := range gravamen_names {
		// agrupar por nombre de gravamen
		comprobantesDetalles, montoImponibleTotalDetalles := AgruparPorGravamen(detalles, name)
		// crear comprobante por cada tipo de gravamen
		comprobante := entities.Comprobante{
			ClienteId:           request.ClienteId,
			Importe:             uint64(montoImponibleTotalDetalles.Int64()),
			Numero:              "", // el numero se crea en el repositorio con el id del comprobante creado
			RazonSocial:         cliente.Razonsocial,
			Domicilio:           cliente.Domicilio,
			Cuit:                cliente.Cuit,
			ComprobanteDetalles: comprobantesDetalles,
			Gravamen:            name,
			MovimientosId:       movsidsString,
			ReporteId:           request.NumeroReporteRrm,
		}
		// Pueden existir tantos comprobantes como gravamenes hay. Solo se computan los que tienen detalles
		if len(comprobante.ComprobanteDetalles) > 0 {
			comprobantes = append(comprobantes, comprobante)
		}
	}

	// Guardar comprobante de retencion y sus detalles
	erro = s.repository.GenerarCertificacionRepository(comprobantes)
	if erro != nil {
		s.utilService.BuildLog(erro, "GenerarCertificacionService")
		return
	}

	return
}

func (s *service) GetMovimientosRetencionesService(request retenciondtos.RentencionRequestDTO) (listaMovimientosId []uint, erro error) {

	listaMovimientosId, erro = s.repository.GetMovimientosRetencionesRepository(request)

	if erro != nil {
		return
	}

	return
}

// determina si el monto es objeto de retencion comparando con un importe de referencia
func CompararMontoMinimo(resultado []retenciondtos.RetencionAgrupada, gravamenes map[string]uint64) ([]retenciondtos.RetencionAgrupada, map[string]uint64) {

	// sumar y acumular los montos de retenciones por gravamen
	for _, r := range resultado {
		gravamenes[r.Gravamen] = gravamenes[r.Gravamen] + uint64(r.TotalRetencion)
	}

	for i := range resultado {
		// si el importe minimo es 0, implica que no tiene importe minimo. Se retiene el monto porque no esta sujeto a una condicion de importe minimo
		if resultado[i].Minimo == 0 {
			resultado[i].Retener = true
		}
		// si el importe minimo es distinto de 0, se compara el minimo con el monto total. Si el monto total de retenciones es igual o mayor el monto minimo, corresponde retener. Caso contrario, corresponde devolver
		if resultado[i].Minimo > 0 {
			// el monto acumulado del gravamen de la retencion actual del loop

			acumulado := entities.Monto(gravamenes[resultado[i].Gravamen]).Float64()
			if acumulado > resultado[i].Minimo {
				resultado[i].Retener = true
			}
		}
	}
	return resultado, gravamenes
}

func AgruparPorGravamen(detalles []entities.ComprobanteDetalle, gravamen_name string) (output []entities.ComprobanteDetalle, montoImponibleDetalle entities.Monto) {

	for _, detalle := range detalles {
		if gravamen_name == detalle.Gravamen {
			output = append(output, detalle)
		}
	}

	// para cada comprobante detalle se suman sus montos imponibles
	for _, cd := range output {
		montoImponibleDetalle += cd.TotalMonto
	}

	return
}

func (s *service) ComprobantesRetencionesDevolverService(request retenciondtos.RentencionRequestDTO) (comprobantesdto retenciondtos.DevolverRetencionesDTO, erro error) {
	var total uint64
	erro = request.ValidarFechas()
	if erro != nil {
		return
	}

	// erro = request.Validar()
	// if erro != nil {
	// 	return
	// }

	comprobantes, erro := s.repository.ComprobantesRetencionesDevolverRepository(request)
	if erro != nil {
		return
	}

	for _, c := range comprobantes {
		var tempComprobanteDTO retenciondtos.ComprobanteResponseDTO
		tempComprobanteDTO.FromEntity(c)
		comprobantesdto.ListaComprobantes = append(comprobantesdto.ListaComprobantes, tempComprobanteDTO)
		// acumular los montos de las retenciones de cada comprobante detalle
		for _, cd := range c.ComprobanteDetalles {
			total += uint64(cd.TotalRetencion)
		}
	}

	// el monto a devolver es la suma de lo retenido por todos los gravamenes
	comprobantesdto.MontoDevolver = total

	return
}

func (s *service) TotalizarRetencionesMovimientosService(listaMovimientoIds []uint) (totalRetenciones entities.Monto, erro error) {
	var resultado uint64
	if len(listaMovimientoIds) == 0 {
		erro = errors.New("lista de ids de movimientos par totalizar rtenciones está vacia")
		return
	}

	resultado, erro = s.repository.TotalizarRetencionesMovimientosRepository(listaMovimientoIds)
	if erro != nil {
		return
	}

	totalRetenciones = entities.Monto(resultado)
	return
}

func (s *service) GetComprobantesService(request retenciondtos.RentencionRequestDTO) (comprobantesdto []retenciondtos.ComprobanteResponseDTO, erro error) {
	var comprobantes []entities.Comprobante
	erro = request.Validar()
	if erro != nil {
		return
	}

	comprobantes, erro = s.repository.GetComprobantesRepository(request)
	if erro != nil {
		return
	}

	for _, c := range comprobantes {
		var tempComprobanteDTO retenciondtos.ComprobanteResponseDTO
		tempComprobanteDTO.FromEntity(c)
		comprobantesdto = append(comprobantesdto, tempComprobanteDTO)
	}

	return
}

func (s *service) NotificarVencimientoCertificadosService(request retenciondtos.CertificadoVencimientoDTO) (erro error) {

	certificados, erro := s.repository.GetCertificadosVencimientoRepository(request)

	if erro != nil {
		return
	}

	if len(certificados) > 0 {
		for _, certificado := range certificados {

			fecha_caducidad_ventana := certificado.Fecha_Caducidad.AddDate(0, 0, -3)
			hoy := time.Now()

			fecha_presentacion_string := certificado.Fecha_Presentacion.Format("2006-01-02")
			fecha_vencimiento_string := certificado.Fecha_Caducidad.Format("2006-01-02")

			if hoy.After(fecha_caducidad_ventana) && hoy.Before(certificado.Fecha_Caducidad) {

				reemplazos := []string{fecha_presentacion_string, certificado.ClienteRetencion.Retencion.Condicion.Gravamen.Gravamen, fecha_vencimiento_string}

				var emails []string
				// Establecer a que correo enviar (Email contacto o emails reportes)
				emails = append(emails, certificado.ClienteRetencion.Cliente.Emailcontacto)
				// for _, contactoReporte := range *certificado.ClienteRetencion.Cliente.Contactosreportes {
				// 	emails = append(emails, contactoReporte.Email)
				// }

				// if request.Administracion {
				// 	// Agregar Emails para Control de Administracion (hardcodeado o tabla)
				// 	// emails = append(emails, "sebasescobar2210@gmail.com")
				// }

				filtro := utildtos.RequestDatosMail{

					Email:            emails,
					Asunto:           "Certificado de Retencion próximo a caducar",
					From:             "Wee.ar!",
					Nombre:           certificado.ClienteRetencion.Cliente.Cliente,
					Mensaje:          "Se le informa que el certificado presentado el #0 con gravamen #1 tiene proxima fecha de caducidad, siendo la misma #2.",
					CamposReemplazar: reemplazos,
					AdjuntarEstado:   false,
					TipoEmail:        "template",
				}

				erro = s.utilService.EnviarMailService(filtro)

				if erro != nil {
					return
				}

			}

		}
	}

	return
}

func (s *service) CreateComisionManualService(request administraciondtos.RequestComisionManual) (err error) {
	//Envio el filtro necesario, los ids de los movimientos que recibo
	//Traigo los pagos intentos asociados a los movimientos con "AcumularPorPagoIntentos" y "CargarPagoIntentos"
	//Con "CargarComision" traigo la comision y los impuestos, si tienen debe devolver un error
	movimientoFiltro := filtros.MovimientoFiltro{
		Ids:                     request.MovimientoId,
		AcumularPorPagoIntentos: true,
		CargarComision:          true,
		CargarPagoIntentos:      true,
	}
	movimientos, _, err := s.repository.GetMovimientos(movimientoFiltro)
	if err != nil {
		return
	}
	if len(movimientos) == 0 {
		err = errors.New("no se encontro el movimiento asociado al pago")
		return
	}
	pagosID := []uint{}
	for _, movimiento := range movimientos {
		//Debo verificar que no tenga comisiones ni impuestos asociados
		if len(movimiento.Movimientocomisions) == 0 && len(movimiento.Movimientoimpuestos) == 0 {
			pagosID = append(pagosID, uint(movimiento.Pagointentos.PagosID))
		}
	}
	if len(pagosID) == 0 {
		return errors.New("los movimientos enviados ya tienen comisiones o impuestos asociados")
	}
	//Busco los pagoIntentos asociados para traer los siguientes datos:
	//"CargarPagoTipo", "CargarCuenta", "CargarCliente", "CargarImpuestos" para traer el impuesto asociado
	//"Channel" para traer los medios de pagos y el channel del medio de pago
	pagoIntentoFiltro := filtros.PagoIntentoFiltro{
		PagosId:         pagosID,
		CargarPagoTipo:  true,
		CargarCuenta:    true,
		CargarCliente:   true,
		CargarImpuestos: true,
		Channel:         true,
	}
	var listaCuentaComision []entities.Cuentacomision
	pagosIntentos, err := s.repository.GetPagosIntentos(pagoIntentoFiltro)
	if err != nil {
		return
	}
	if len(pagosIntentos) == 0 {
		return fmt.Errorf("no se encontro el pago asociado al movimiento %v ingresado", request.MovimientoId)
	}
	for _, pagoIntento := range pagosIntentos {
		//Busco la cuenta comision para traer el porcentaje de comision de la cuenta del cliente con el pago
		filtroComisionChannel := filtros.CuentaComisionFiltro{
			CargarCuenta:      true,
			ChannelId:         uint(pagoIntento.Mediopagos.ChannelsID),
			CuentaId:          pagoIntento.Pago.PagosTipo.Cuenta.ID,
			FechaPagoVigencia: pagoIntento.PaidAt,
			Channelarancel:    true,
		}
		cuentaComision, erro := s.repository.GetCuentaComision(filtroComisionChannel)
		if erro != nil {
			return erro
		}
		listaCuentaComision = append([]entities.Cuentacomision{}, cuentaComision)
		for _, movimiento := range movimientos {
			if movimiento.Pagointentos.ID == pagoIntento.ID {
				err = s.utilService.BuildComisionesRedondeoMenosExigente(&movimiento, &listaCuentaComision, pagoIntento.Pago.PagosTipo.Cuenta.Cliente.Iva, movimiento.Monto)
				if err != nil {
					return
				}
				movimiento.Pagointentos = nil
				ctxPrueba := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				//Actualizo el monto del movimiento. Al tener ya la asociacion con movimientoImpuesto y movimientoComision realizada, al actualizar el movimiento, me crea los registros en las tablas correspondientes de  movimientoImpuesto y movimientoComision
				err = s.repository.UpdateMovimientoMontoRepository(ctxPrueba, movimiento)
				if err != nil {
					return
				}
				break
			}
		}
	}
	return
}

func (s *service) GetClientesConfiguracionService(filtro filtros.ClienteConfiguracionFiltro) (response administraciondtos.ResponseClientesConfiguracion, erro error) {
	clientes, erro := s.repository.GetClientesConfiguracion(filtro)

	if erro != nil {
		return
	}

	for _, cliente := range clientes {
		cli := administraciondtos.ClienteConfiguracionInfo{
			Id:      cliente.ID,
			Cliente: cliente.Cliente,
		}

		response.Clientes = append(response.Clientes, cli)
	}

	return
}

func (s *service) NotificarPagosWebhookSinNotificarService() (err error) {

	// Obtengo el id del estado PAID
	filtroPagosEstado := filtros.PagoEstadoFiltro{
		Nombre: "PAID",
	}
	pagoEstado, erro := s.repository.GetPagosEstados(filtroPagosEstado)
	if erro != nil {
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionWebhook,
			Descripcion: fmt.Sprintf("No se pudo obtener los pagos estados para notificar"),
		}
		s.CreateNotificacionService(notificacion)
		return
	}
	if len(pagoEstado) == 0 {
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionWebhook,
			Descripcion: "no se pudo obtener los pagos estados para notificar, por que no existe un pago con nombre PAID",
		}
		s.CreateNotificacionService(notificacion)
		return
	}
	//Obtengo los pagos que estan en estado PAID, que no fueron notificados, desde las 00:00 del dia de hoy hasta el momento que se ejecute y con los medios de pago de tarjeta (hardcodeo en consulta)
	var pagosEstados []uint64 = []uint64{uint64(pagoEstado[0].ID)}
	now := time.Now()
	fechaInicio := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	filtroPagos := filtros.PagoTipoFiltro{
		FiltroEstadoInicialNotificado:      true,
		FiltroMediopagosID:                 true,
		PagoEstadosIds:                     pagosEstados,
		FechaPagoInicio:                    fechaInicio,
		FechaPagoFin:                       now,
		CargarPagosEstadoInicialNotificado: true,
	}
	//Obtengo los pagostipos juntos a sus pagos, por que ahi estan los url a notificar
	pagosTipo, _, err := s.repository.GetPagosTipo(filtroPagos)

	if err != nil {
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionWebhook,
			Descripcion: fmt.Sprintf("No se pudo obtener los pagos para notificar. %s", err),
		}
		s.CreateNotificacionService(notificacion)
		return
	}
	if len(pagosTipo) == 0 {
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionWebhook,
			Descripcion: fmt.Sprintln("No existen pagos por notificar"),
		}
		s.CreateNotificacionService(notificacion)
		return
	}
	//Creo la estructura que necesito para notificar
	pagosNotificar, err := s.CreateNotificacionPagosService(pagosTipo)
	if err != nil {
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionWebhook,
			Descripcion: fmt.Sprintf("No se pudo obtener los pagos para notificar. %s", err),
		}
		s.CreateNotificacionService(notificacion)
		return
	}
	if len(pagosNotificar) == 0 {
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionWebhook,
			Descripcion: fmt.Sprintln("No existen pagos por notificar"),
		}
		s.CreateNotificacionService(notificacion)
		return
	}
	//Notifico los pagos
	pagosupdate := s.NotificarPagos(pagosNotificar)
	if len(pagosupdate) == 0 {
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionWebhook,
			Descripcion: fmt.Sprintln("No se pudo notificar ningun pago"),
		}
		s.CreateNotificacionService(notificacion)
		return
	}
	err = s.UpdatePagosEstadoInicialNotificado(pagosupdate) /* actualzar estado de pagos a notificado */
	if err != nil {
		logs.Info(fmt.Sprintf("Los siguientes pagos que se notificaron al cliente no se actualizaron: %v", pagosupdate))
		logs.Error(err)
		notificacion := entities.Notificacione{
			Tipo:        entities.NotificacionWebhook,
			Descripcion: fmt.Sprintf("webhook: Error al actualizar estado de pagos a notificado .: %s", err),
		}
		s.CreateNotificacionService(notificacion)
	}
	return

}

func (s *service) GetMovimientosIdsCalculoRetencionComprobanteService(request retenciondtos.RentencionRequestDTO) (resultado []uint, erro error) {
	return s.repository.GetMovimientosIdsCalculoRetencionComprobante(request)
}

func (s *service) CalcularRetencionesByTransferenciasSinAgruparService(request retenciondtos.RentencionRequestDTO) (resultado []retenciondtos.RetencionAgrupada, erro error) {
	return s.repository.CalcularRetencionesByTransferenciasSinAgruparRepository(request)
}
