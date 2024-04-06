package routes

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos/retenciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/webhook"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	uuid "github.com/satori/go.uuid"

	apiresponder "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/responderdto"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	"github.com/gofiber/fiber/v2"
)

func AdministracionRoutes(app fiber.Router, middlewares middlewares.MiddlewareManager, service administracion.Service, util util.UtilService, banco banco.BancoService) {

	// WEBHOOK:EJECUTAR PROCESO NOTIFICACION DE PAGOS A CLIENTES
	app.Get("/notificacion-pagos", middlewares.ValidarPermiso("psp.herramienta"), NotificacionPagos(service))                    // psp.notificacion.pagos
	app.Get("/notificar-pagos-references", middlewares.ValidarPermiso("psp.herramienta"), NotificarPagosWithReferences(service)) // psp.notificacion.pagos

	app.Get("/get-pago", GetPago(service))
	// ? CIERRELOTE APILINK
	// ? PASO 1 Buscar en Apilink lote de pagos, actualzar solo con los llegaron a un estado final se crean en la tabla apilinkcierrelote y se actualiza el estado de los pagos
	// ? PASO 2 Informar a clientes el cambio de estado de los debines
	// ? PASO 3 Conciliacion con banco
	// ? PASO 4 Generacion de movimientos
	app.Get("/apilink-cierre-lote", middlewares.ValidarPermiso("psp.herramienta"), CierreLoteApilink(service))
	app.Get("/notificar-pagos-clapilink-actualizados", middlewares.ValidarPermiso("psp.herramienta"), NotificarPagosClApilink(service))
	app.Get("/apilink-cierre-lote-conciliar-banco", middlewares.ValidarPermiso("psp.herramienta"), CierreLoteApilinkBanco(service, banco))
	app.Get("/apilink-cierre-lote-generarmov", middlewares.ValidarPermiso("psp.herramienta"), CierreLoteApilinkMov(service))

	app.Get("/get-debin-id", middlewares.ValidarPermiso("psp.herramienta"), GetDebinById(service))

	// Crea registros en "apilinkcierrelotes" que estan en apilink y no se crearon por algun fallo. Luego de crear los registros en apilinkcierrelotes se puede notificar los mismos manualmente desde la herramienta (API-Link ---> Ejecutar proceso ----> Notificar Debines).
	app.Get("/apilink-debin-not-registered", middlewares.ValidarPermiso("psp.herramienta"), DebinNotRegisteredApilink(service))

	// ? TRANSFERENCIAS AUTOMATICAS - CLIENTES - COMISIONES
	app.Get("/transferencias-automaticas", middlewares.ValidarPermiso("psp.herramienta"), transferenciasAutomaticas(service)) // psp.transferencia.automatica
	app.Get("/transferencias-automaticas-subcuentas", middlewares.ValidarPermiso("psp.herramienta"), transferenciasAutomaticasSubcuentas(service))
	app.Post("/transferencia-comisiones-impuestos", middlewares.ValidarPermiso("psp.herramienta"), transferenciaComisionesImpuestos(service))

	// & CIERRELOTE RAPIPAGO
	//& 2 actualizar estado del pago con lo encontrado en cierrelote(archivo recibido)
	//& 3 Notificar al cliente el pago actualizado
	// & 4 conciliacion con banco
	// & 5 generar movimientos manuales
	app.Get("/actualizar-pagos-cl", middlewares.ValidarPermiso("psp.herramienta"), ActualizarEstadosPagosClRapipago(service, banco))
	app.Get("/notificar-pagos-cl-actualizados", middlewares.ValidarPermiso("psp.herramienta"), NotificarPagosClRapipago(service))
	app.Get("/rapipago-cierre-lote", middlewares.ValidarPermiso("psp.herramienta"), CierreLoteRapipago(service, banco))
	app.Get("/generar-movimiento-rapipago", middlewares.ValidarPermiso("psp.herramienta"), GenerarMovimientosRapipago(service, banco)) // psp.rapipago.cierre.lote                                            // psp.rapipago.cierre.lote

	// & CIERRELOTE MULTIPAGOS
	//& 2 actualizar estado del pago con lo encontrado en cierrelote(archivo recibido)
	//& 3 Notificar al cliente el pago actualizado
	// & 4 conciliacion con banco
	// & 5 generar movimientos manuales
	app.Get("/actualizar-pagos-cl-mp", ActualizarEstadosPagosClMultipagos(service, banco))
	app.Get("/notificar-pagos-cl-actualizados-mp", NotificarPagosClMultipagos(service))
	app.Get("/multipagos-cierre-lote", CierreLoteMultipagos(service, banco))
	app.Get("/generar-movimiento-multipagos", GenerarMovimientosMultipagos(service, banco)) // psp.rapipago.cierre.lote                                            // psp.rapipago.cierre.lote

	// ? EXPIRAR PAGOS OFFLINES
	app.Get("/caducar-pagosintentos-offline", middlewares.ValidarPermiso("psp.herramienta"), getCaducarOfflineIntentos(service))

	// consultar clientes
	app.Get("/clientes-herramienta", middlewares.ValidarPermiso("psp.herramienta"), getClientes(service))
	app.Post("/clientes-configuracion", middlewares.ValidarPermiso("psp.herramienta"), getClientesConfiguracion(service))

	// RETENCIONES IMPOSITIVAS
	app.Get("/retenciones", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetRetenciones(service))
	app.Get("/cliente-retenciones", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetClienteRetenciones(service))
	app.Post("/asignar-retencion-cliente", middlewares.ValidarPermiso("psp.consultar.impuestos"), CreateClienteRetencion(service))
	app.Post("/certificado", middlewares.ValidarPermiso("psp.consultar.impuestos"), PostRetencionCertificados(service))
	app.Get("/certificado", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetRetencionCertificados(service))
	app.Get("/calcular-retenciones", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetCalcularRetenciones(service))
	app.Get("/condiciones", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetCondiciones(service))
	app.Post("/condicion", middlewares.ValidarPermiso("psp.consultar.impuestos"), CreateCondicion(service))
	app.Put("/condicion", middlewares.ValidarPermiso("psp.consultar.impuestos"), UpdateCondicion(service))
	app.Post("/retencion", middlewares.ValidarPermiso("psp.consultar.impuestos"), CreateRetencion(service))
	app.Delete("/delete-cliente-retencion", middlewares.ValidarPermiso("psp.consultar.impuestos"), DeleteClienteRetencion(service))
	app.Put("/retencion", middlewares.ValidarPermiso("psp.consultar.impuestos"), UpdateRetencion(service))
	app.Get("/gravamenes", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetGravamenes(service))
	// app.Get("/comprobar-minimo-retencion", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetComprobarMinimoRetencion(service))
	app.Get("/evaluar-retencion-cliente", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetEvaluarRetencionCliente(service))
	app.Post("/evaluar-retencion-movimiento", middlewares.ValidarPermiso("psp.consultar.impuestos"), EvaluarRetencionMovimiento(service))
	app.Get("/generar-certificacion", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetGenerarCertificacion(service))
	/* app.Get("/pruebaMovimientosSubcuentas", getMovimientosSubcuentas(util, service)) */
	app.Get("/retenciones-devolver", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetComprobantesRetencionesDevolver(service))
	app.Get("/comprobantes", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetComprobantes(service))

	/* Inicio*/
	// ^ Proceso que permitira calcular las comisiones de pagos del dia
	app.Get("/calcular-movimientostemporales-pagos", middlewares.ValidarPermiso("psp.herramienta"), GenerarMovimientosTemporalesPagos(service)) //

	app.Get("informar-vencimientos-certificados", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetInformarVencimientosCertificados(service))
	// Crear comision manual

	app.Post("/comision-manual", middlewares.ValidarPermiso("psp.herramienta"), createComisionManual(service))

	// Cambiar estado pagos expirados
	app.Put("/cambiar-estado-pagos-expirados", middlewares.ValidarPermiso("psp.herramienta"), CambiarEstadoPagosExpirados(service))

	//app.Post("/create-auditoria", middlewares.ValidarPermiso("psp.herramienta"), createAuditoria(service))
	// Cambiar estado pagos expirados
	app.Get("/notificar-pagos", middlewares.ValidarPermiso("psp.herramienta"), NotificarPagos(service))

	// obtener un archivo de comprobante de retencion
	app.Get("/archivo-comprobante-retencion", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetArchivoComprobanteRetencion(service))
	// obtener un archivo de resumen rrm
	app.Get("/archivo-reporte-rendicion-mensual", middlewares.ValidarPermiso("psp.consultar.impuestos"), GetArchivoReporteRendicionMensual(service))
}

/* func getMovimientosSubcuentas(utilService util.UtilService, administrasionService administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		movimientoId := c.Query("movimientoId")
		MovimientoId, err := strconv.ParseUint(movimientoId, 10, 64)

		movimiento, err := administrasionService.GetMovimientosById(MovimientoId)
		if err != nil {
			return err
		}

		cuenta := entities.Cuenta{
			ID: uint(movimiento.CuentasId),
		}
		movimientoSubcuentas := []entities.Movimientosubcuenta{}
		err = utilService.BuildMovimientoSubcuentas(&movimiento, int64(movimiento.Monto), &movimientoSubcuentas, &cuenta)
		if err != nil {
			return err
		}
		var movimientosSubcuentasResponse []administraciondtos.MovimientosSubcuentas

		//Cargo un array para la respuesta
		for _, v := range movimientoSubcuentas {

			var movimientoSubcuenta administraciondtos.MovimientosSubcuentas
			movimientoSubcuenta.Monto = v.Monto
			movimientoSubcuenta.SubcuentasID = v.SubcuentasID
			movimientoSubcuenta.MovimientosID = v.MovimientosID

			movimientosSubcuentasResponse = append(movimientosSubcuentasResponse, movimientoSubcuenta)
		}

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"data":    movimientosSubcuentasResponse,
			"message": "Calculo realizado con exito.",
		})

	}
} */

// func GetComprobarMinimoRetencion(service administracion.Service) fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 		filtro := retenciondtos.RentencionRequestDTO{
// 			RetencionId: 18,
// 		}

// 		getDTO := false
// 		retenciones, err := service.GetRetencionesService(filtro, getDTO)
// 		if err != nil {
// 			r := apiresponder.NewResponse(404, nil, "Ocurrio un error: "+err.Error(), c)
// 			return r.Responder()
// 		}
// 		result, monto, err := service.ComprobarMinimoRetencion(retenciones.Retenciones[0], 9)
// 		if err != nil {
// 			r := apiresponder.NewResponse(404, nil, "Ocurrio un error: "+err.Error(), c)
// 			return r.Responder()
// 		}
// 		fmt.Println(result, monto)

// 		r := apiresponder.NewResponse(200, nil, "OK", c)
// 		return r.Responder()
// 	}
// }

// CONSULTAR POR UN DEBIN POR EL ID
func GetDebinById(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		//obtener lista pagos apilink(solo los conciliados con banco)
		var request linkdebin.RequestGetDebinLink
		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		request.Cbu = config.CBU_CUENTA_TELCO
		debin, err := service.GetDebinService(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, debin, "Debin obtenido con éxito.", c)
		return r.Responder()
	}
}

func GetPago(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type External struct {
			ExternalReference string `json:"external_reference"`
		}

		var req External
		err := c.QueryParser(&req)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		filtro := filtros.PagoFiltro{
			ExternalReference: req.ExternalReference,
			CargaPagoIntentos: true,
			CargarPagoEstado:  true,
		}
		pago, err := service.GetPaymentByExternalService(filtro)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			fmt.Print(pago)
			return r.Responder()
		}

		// caso de exito del proceso
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"status":  true,
			"data":    pago,
			"message": "se obtuvieron las datos exitosamente",
		})
	}
}

func EvaluarRetencionMovimiento(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request retenciondtos.RentencionRequestDTO

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.EvaluarRetencionesByMovimientoService(request)

		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Ocurrio un error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron las datos exitosamente", c)
		return r.Responder()
	}
}

func GetEvaluarRetencionCliente(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.EvaluarRetencionesByClienteService(request)

		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Ocurrio un error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron las datos exitosamente", c)
		return r.Responder()
	}
}

func DebinNotRegisteredApilink(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request cierrelotedtos.RequestDebinNotRegisteredApilink

		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		// obtener debines de apilink
		listas, err := service.BuildDebinNotRegisteredApiLinkService(request)
		logs.Info(listas)

		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"status":  false,
				"message": "Error: " + err.Error(),
				"code":    400,
			})
		}

		// caso no existen debines para procesar
		if len(listas.ListaPagos) == 0 || len(listas.ListaCLApiLink) == 0 {
			r := apiresponder.NewResponse(404, "tipoconciliacion: CierreLoteApilink", "el proceso fue ejecutado con exito. no existen debines para procesar", c)
			return r.Responder()
		}

		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
		err = service.CreateCLApilinkPagosService(ctx, listas)
		if err != nil {
			logs.Error(err)
			return c.Status(400).JSON(&fiber.Map{
				"status":  false,
				"message": "Error: " + err.Error(),
				"code":    400,
			})
		}

		// caso de exito del proceso
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": "Registros creados en apilinkcierrelotes exitosamente",
			"code":    200,
		})
	}
}

func CierreLoteApilink(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// obtener debines de apilink
		listas, err := service.BuildCierreLoteApiLinkService()

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error: "+err.Error(), c)
			return r.Responder()
		}

		// caso no existen debines para procesar
		if len(listas.ListaPagos) == 0 || len(listas.ListaCLApiLink) == 0 {
			r := apiresponder.NewResponse(404, "tipoconciliacion: CierreLoteApilink", "el proceso fue ejecutado con exito. no existen debines para procesar", c)
			return r.Responder()
		}

		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
		err = service.CreateCLApilinkPagosService(ctx, listas)
		if err != nil {
			logs.Error(err)
			r := apiresponder.NewResponse(400, nil, "error: "+err.Error(), c)
			return r.Responder()
		}

		// caso de exito del proceso
		r := apiresponder.NewResponse(200, "tipoconciliacion: CierreLoteApilink", "el proceso fue ejecutado con exito", c)
		return r.Responder()
	}
}

func NotificarPagosClApilink(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")
		// SE CONSULTA EN LA TABLA APILINKCIERRE LOTE LOS DEBINES QUE AUN NO FUERON INFORMADOS
		filtro := linkdebin.RequestDebines{
			BancoExternalId: false,
			Pagoinformado:   true,
		}
		debines, erro := service.GetConsultarDebines(filtro)
		if erro != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+erro.Error(), c)
			return r.Responder()
		}
		//    buscar pagos encontrados en tabla apilinkcierrelote
		if len(debines) > 0 {
			// NOTE construir lote de pagos debin que se notificara al cliente
			pagos, debin, erro := service.BuildNotificacionPagosCLApilink(debines)
			if erro != nil {
				r := apiresponder.NewResponse(400, nil, "Error: "+erro.Error(), c)
				return r.Responder()
			}
			if len(pagos) > 0 && len(debin) > 0 {
				// NOTE enviar lote de pagos a clientes
				pagosNotificar := service.NotificarPagos(pagos)
				if len(pagosNotificar) > 0 {
					//NOTE Si se envian los pagos con exito se debe actualziar el campo pagoinformado en la tabla aplilinkcierrelote
					filtro := linkdebin.RequestListaUpdateDebines{
						DebinId: debin,
					}
					erro := service.UpdateCierreLoteApilink(filtro)
					if erro != nil {
						logs.Error(erro)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote apilink pagoinformado: %s", erro),
						}
						service.CreateNotificacionService(notificacion)
						r := apiresponder.NewResponse(404, "notificacion pagos apilink: NotificarPagosClApilink", "no se pudieron actualizar pagosinformados de clapilink", c)
						return r.Responder()
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionWebhook,
						Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos"),
					}
					service.CreateNotificacionService(notificacion)
					r := apiresponder.NewResponse(400, nil, "Error: "+erro.Error(), c)
					return r.Responder()
				}
			}

			data := map[string]interface{}{
				"pagos informados":           pagos,
				"debines actualizados":       debin,
				"notificacion pagos apilink": "NotificarPagosClApilink",
			}
			r := apiresponder.NewResponse(200, data, "el proceso fue ejecutado con exito", c)
			return r.Responder()
		} else {
			r := apiresponder.NewResponse(404, "notificacion pagos apilink: NotificarPagosClApilink", "no existen debines para informar", c)
			return r.Responder()
		}
	}
}

func CierreLoteApilinkBanco(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")
		// SE CONSULTA EN LA TABLA APILINKCIERRELOTE LOS DEBINES QUE AUN NO FUERON INFORMADOS
		filtro := linkdebin.RequestDebines{
			BancoExternalId:  false,
			CargarPagoEstado: true,
		}
		debines, err := service.GetDebines(filtro)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(debines) > 0 {
			if err == nil {
				request := bancodtos.RequestConciliacion{
					TipoConciliacion: 2,
					ListaApilink:     debines,
				}
				listaCierreApiLinkBanco, listaBancoId, err := banco.ConciliacionPasarelaBanco(request)
				// 1.1 conciliar lista de debines de apilink con los movimientos de banco
				// listaBancoId se utilizara para actualizar los movimientos de banco
				// listaCierreApiLinkBanco, listaBancoId, err := banco.ConciliacionBancoApliLInk(listaCierre)

				/*si no hay error guardar listaCierreloteapilink en la base de datos */
				if err != nil {
					logs.Error(err)
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("error al conciliar movimiento banco y cierre loteapilink: %s", err),
					}
					service.CreateNotificacionService(notificacion)
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				} else {

					// NOTE Actualizar lista de cierreloteapilink campo banco external_id, match y fecha de acreditacion
					if len(listaCierreApiLinkBanco.ListaApilink) > 0 || len(listaCierreApiLinkBanco.ListaApilinkNoAcreditados) > 0 {
						listas := linkdebin.RequestListaUpdateDebines{
							Debines:              listaCierreApiLinkBanco.ListaApilink,
							DebinesNoAcreditados: listaCierreApiLinkBanco.ListaApilinkNoAcreditados,
						}
						//Actualiza los "apilinkcierrelotes" y da baja los que no se acreditaron
						erro := service.UpdateCierreLoteApilink(listas)
						if erro != nil {
							logs.Error(erro)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote apilink: %s", erro),
							}
							service.CreateNotificacionService(notificacion)
							r := apiresponder.NewResponse(400, "conciliacion banco-apilink: CierreLoteApilinkBanco", "error al actualizar registros de cierrelote apilink", c)
							return r.Responder()
						}
					}

					if len(listaBancoId) > 0 {
						_, err := banco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes movimientos del banco no se actualizaron: %v", listaBancoId))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar registros del banco: %s", err),
							}
							service.CreateNotificacionService(notificacion)
							r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
							return r.Responder()
							// en le caso de este error y si el pago no se actualizo a estados finales no afecta el cierre de apilink
							// el estado del pago se actualiza a estado final y no tendra en cuenta al consultar a apilink
							// ACCION : se debe actualizar manualmente el campo check en la tabla de movimientos de banco(no es obligatorio)
						}
					}
				}

			}

		} else {
			r := apiresponder.NewResponse(404, "tipoconciliacion: listaCierreApiLinkBanco", "no existen pagos con debin para conciliar", c)
			return r.Responder()
		}

		// caso de exito
		r := apiresponder.NewResponse(200, "tipoconciliacion: listaCierreApiLinkBanco", "conciliacion banco-apilink se ejecuto con exito", c)
		return r.Responder()
	}
}

// generar mov pagos debin conciliados con banco
func CierreLoteApilinkMov(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		//obtener lista pagos apilink(solo los conciliados con banco)
		filtro := linkdebin.RequestDebines{
			BancoExternalId:  true,
			CargarPagoEstado: true,
		}
		debines, err := service.GetDebines(filtro)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(debines) > 0 {
			responseCierreLote, err := service.BuildMovimientoApiLink(debines)
			if err == nil {
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				err = service.CreateMovimientosService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)

					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				} else {
					// caso de exito
					r := apiresponder.NewResponse(200, "tipoconciliacion: listaCierreApiLinkBanco", "proceso ejecutado con exito", c)
					return r.Responder()
				}
			}
		}
		r := apiresponder.NewResponse(404, "tipoconciliacion: listaCierreApiLinkBanco", "no existen debines para generar movimientos", c)
		return r.Responder()
	}
}

func NotificacionPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request webhook.RequestWebhook

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		pagos, err := service.BuildNotificacionPagosService(request)

		if err == nil {
			pagosNotificar, err := service.CreateNotificacionPagosService(pagos)
			if err == nil {
				if len(pagosNotificar) > 0 {
					// si inicia el proceso de notifiacar al cliente
					pagosupdate := service.NotificarPagos(pagosNotificar)
					// NOTE se debe actualizar solo si el pago llegi a un estado final
					if len(pagosupdate) > 0 && request.EstadoFinalPagos { /* actualzar estado de pagos a notificado */
						err = service.UpdatePagosNoticados(pagosupdate)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes pagos que se notificaron al cliente no se actualizaron: %v", pagosupdate))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionWebhook,
								Descripcion: fmt.Sprintf("webhook: Error al actualizar estado de pagos a notificado .: %s", err),
							}
							service.CreateNotificacionService(notificacion)
							return fiber.NewError(400, "Error: "+err.Error())
						}
					}
					if len(pagosupdate) > 0 {
						return c.JSON(&fiber.Map{
							"data":    pagosupdate,
							"message": "se notifico con exito los siguientes pagos",
						})
					}
					if len(pagosupdate) == 0 {
						return c.JSON(&fiber.Map{
							"message": "error al notificar pagos a clientes",
						})
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionWebhook,
						Descripcion: fmt.Sprintf("webhook: No existen pagos por notificar. %s", err),
					}
					service.CreateNotificacionService(notificacion)
					return c.JSON(&fiber.Map{
						"error":            "no existen pagos por notificar",
						"tipoconciliacion": "notificacionPagos",
					})
				}
			}
		}

		return c.JSON(&fiber.Map{
			"error":            err,
			"tipoconciliacion": "notificacionPagos",
		})

	}
}

func NotificarPagosWithReferences(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request webhook.RequestWebhookReferences

		if err := c.QueryParser(&request); err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		pagos, err := service.BuildNotificacionPagosWithReferences(request)

		if err == nil {
			pagosNotificar, err := service.CreateNotificacionPagosService(pagos)
			if err == nil {
				if len(pagosNotificar) > 0 {
					// si inicia el proceso de notifiacar al cliente
					pagosupdate := service.NotificarPagos(pagosNotificar)
					// NOTE se debe actualizar solo si el pago llegi a un estado final
					if len(pagosupdate) > 0 && request.EstadoFinalPagos { /* actualzar estado de pagos a notificado */
						err = service.UpdatePagosNoticados(pagosupdate)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes pagos que se notificaron al cliente no se actualizaron: %v", pagosupdate))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionWebhook,
								Descripcion: fmt.Sprintf("webhook: Error al actualizar estado de pagos a notificado .: %s", err),
							}
							service.CreateNotificacionService(notificacion)
							return fiber.NewError(400, "Error: "+err.Error())
						}
					}
					if len(pagosupdate) > 0 {
						return c.JSON(&fiber.Map{
							"data":    pagosupdate,
							"message": "se notifico con exito los siguientes pagos",
						})
					}
					if len(pagosupdate) == 0 {
						return c.JSON(&fiber.Map{
							"message": "error al notificar pagos a clientes",
						})
					}
				} else {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionWebhook,
						Descripcion: fmt.Sprintf("webhook: No existen pagos por notificar. %s", err),
					}
					service.CreateNotificacionService(notificacion)
					return c.JSON(&fiber.Map{
						"error":            "no existen pagos por notificar",
						"tipoconciliacion": "notificacionPagos",
					})
				}
			}
		}

		return c.JSON(&fiber.Map{
			"error":            err,
			"tipoconciliacion": "notificacionPagos",
		})

	}
}

func transferenciasAutomaticas(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request administraciondtos.TransferenciasClienteId
		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})

		response, erro := service.RetiroAutomaticoClientes(ctx, request)

		if erro != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionPagoExpirado,
				Descripcion: fmt.Sprintf("No se pudo realizar la transferencia automatica para los clientes. %s", erro),
			}
			err := service.CreateNotificacionService(notificacion)
			if err != nil {
				return c.JSON(&fiber.Map{
					"error":            err,
					"tipoconciliacion": "transferenciasAutomaticas",
				})
			}
		}
		var respuesta string
		if len(response.MovimientosId) > 0 {
			respuesta = "se ejecuto con éxito transferencias automaticas de clientes"
		} else {
			respuesta = "no existen movimientos para transferir"
		}
		return c.JSON(&fiber.Map{
			"message":          respuesta,
			"tipoconciliacion": "transferenciasAutomaticas",
		})

	}
}

func transferenciasAutomaticasSubcuentas(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})

		response, erro := service.RetiroAutomaticoClientesSubcuentas(ctx)

		if erro != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionPagoExpirado,
				Descripcion: fmt.Sprintf("No se pudo realizar la transferencia automatica para los clientes. %s", erro),
			}
			err := service.CreateNotificacionService(notificacion)
			if err != nil {
				return c.JSON(&fiber.Map{
					"error":            err,
					"tipoconciliacion": "transferenciasAutomaticas",
				})
			}
		}
		var respuesta string
		if len(response.MovimientosId) > 0 {
			respuesta = "se ejecuto con éxito transferencias automaticas de clientes"
		} else {
			respuesta = "no existen movimientos para transferir"
		}
		return c.JSON(&fiber.Map{
			"message":          respuesta,
			"tipoconciliacion": "transferenciasAutomaticas",
		})

	}
}

func transferenciaComisionesImpuestos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request administraciondtos.RequestComisiones

		err := c.BodyParser(&request)

		if err != nil {
			return fiber.NewError(404, "Error en los parámetros enviados: "+err.Error())
		}

		ctx := getContextAuditable(c)

		uuid := uuid.NewV4()
		result, err := service.SendTransferenciasComisiones(ctx, uuid.String(), request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.JSON(&fiber.Map{
			"resultado": result,
		})
	}
}

func CierreLoteRapipago(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      true,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote - los que no fueron conciliados  */
		listaCierreRapipago, err := service.GetCierreLoteRapipagoService(filtroMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}
		if len(listaCierreRapipago) > 0 {
			if err == nil {

				request := bancodtos.RequestConciliacion{
					TipoConciliacion: 1,
					ListaRapipago:    listaCierreRapipago,
				}
				// aqui hay retornar la lista de id de repipagocierre lote y los id del banco
				listaCierreRapipago, listaBancoId, err := banco.ConciliacionPasarelaBanco(request)

				if len(listaBancoId) == 0 {
					logs.Info("no existen movimientos en banco para conciliar con pagos rapipago")
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("no existen movimientos en banco para conciliar con pagos rapipago: %s", err),
					}
					service.CreateNotificacionService(notificacion)
					return c.JSON(&fiber.Map{
						"error":            notificacion.Descripcion,
						"tipoconciliacion": "CierreLoteRapipago",
					})
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
						// son los registros que coincidieron en el cierrerapipago y banco
						// si no se actualiza los registros del banco se debera actualizar manualmente
						_, err := banco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
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
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}

		return c.JSON(&fiber.Map{
			"tipoconciliacion": "CierreLoteRapipago se ejecuto con exito",
		})

	}
}

func ActualizarEstadosPagosClRapipago(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      false,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote  */
		listaPagoaRapipago, err := service.GetCierreLoteRapipagoService(filtroMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clrapipago. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}
		if len(listaPagoaRapipago) > 0 {
			if err == nil {
				listaPagosClRapipago, err := service.BuildPagosClRapipago(listaPagoaRapipago)
				if err == nil {
					// Actualizar estados del pago y cierrelote
					err = service.ActualizarPagosClRapipagoService(listaPagosClRapipago)
					if err != nil {
						logs.Error(err)
						return fiber.NewError(400, "Error: "+err.Error())
					}

				} else {
					logs.Error(err)
					return fiber.NewError(400, "Error: "+err.Error())
				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("No existen pagos de rapipago para actualizar"),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}

		return c.JSON(&fiber.Map{
			"tipoconciliacion": "CierreLoteRapipago se ejecuto con exito",
		})

	}
}

func NotificarPagosClRapipago(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request filtros.PagoEstadoFiltro

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		pagos, barcode, err := service.BuildNotificacionPagosCLRapipago(request)

		// control de error al llamar al servicio BuildNotificacionPagosCLRapipago
		if err != nil {
			return c.JSON(&fiber.Map{
				"error":            "error al notificar los pagos: " + err.Error(),
				"tipoconciliacion": "notificacionPagosClRapipago",
			})
		}

		if len(barcode) > 0 {
			if len(pagos) > 0 {
				pagosNotificar := service.NotificarPagos(pagos)
				if len(pagosNotificar) == 0 {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionWebhook,
						Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos"),
					}
					service.CreateNotificacionService(notificacion)
					return fiber.NewError(400, "Error: "+err.Error())
				}

			}

			//si la notificacion se realiza con exito se debera actualizar el campo pagonotificado en repiapagocierrolote
			err := service.ActualizarPagosClRapipagoDetallesService(barcode)
			if err != nil {
				return c.JSON(&fiber.Map{
					"error":            "error al actualizar pagos en clrapipago",
					"tipoconciliacion": "notificacionPagosClRapipago",
				})
			}
			mensaje := fmt.Sprintf("los pagos con los siguientes codigos de barra se actualizaron correctamente. %v", barcode)
			return c.JSON(&fiber.Map{
				"pagos actualizados": mensaje,
				"notificacion":       "exitosa",
				"tipoconciliacion":   "notificacionPagosClRapipago",
			})

		} else {
			return c.JSON(&fiber.Map{
				"error":            "no existen pagos por notificar",
				"tipoconciliacion": "notificacionPagosClRapipago",
			})
		}

	}
}

func GenerarMovimientosRapipago(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovRapipago := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: true,
			PagosNotificado:      true,
		}

		listaCierreRapipago, err := service.GetCierreLoteRapipagoService(filtroMovRapipago)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}
		// Si no se guarda ningún cierre no hace falta seguir el proce
		if len(listaCierreRapipago) > 0 {
			// 2 - Contruye los movimientos y hace la modificaciones necesarias para modificar los
			// pagos y demás datos necesarios en caso de error se repetira el día siguiente
			responseCierreLote, err := service.BuildRapipagoMovimiento(listaCierreRapipago)

			if err == nil {

				// 3 - Guarda los movimientos en la base de datos en caso de error se
				// repetira en el día siguiente
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				err = service.CreateMovimientosService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)
					return fiber.NewError(400, "Error: "+err.Error())
				} else {
					return c.JSON(&fiber.Map{
						"mensaje":          "conciliacion exitosa",
						"tipoconciliacion": "CierreLoteRapipago",
					})
				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("no existen pagos para generar movimientos"),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}

		return nil

	}
}

func getCaducarOfflineIntentos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		response, err := service.GetCaducarOfflineIntentos()

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&response)
	}
}

func getContextAuditable(c *fiber.Ctx) context.Context {
	userid := string(c.Response().Header.Peek("user_id"))
	intUserID, _ := strconv.Atoi(userid)
	userctx := entities.Auditoria{
		UserID: uint(intUserID),
		IP:     c.IP(),
	}
	ctx := context.WithValue(c.Context(), entities.AuditUserKey{}, userctx)
	return ctx
}

func getClientes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ClienteFiltro

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		clientes, err := service.GetClientesService(request)

		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&clientes)
	}
}

// Endpoints Retenciones

func GetRetenciones(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}
		getDTO := true
		response, err := service.GetRetencionesService(request, getDTO)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron los resultados exitosamente", c)
		return r.Responder()
	}
}

func GetClienteRetenciones(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO
		var response retenciondtos.RentencionesResponseDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		if !request.Unlinked {
			response, err = service.GetClienteRetencionesService(request)
		}

		if request.Unlinked {
			response, err = service.GetClienteUnlinkedRetencionesService(request)
		}

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron los resultados exitosamente", c)
		return r.Responder()
	}
}

func CreateClienteRetencion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		err = service.CreateClienteRetencionService(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, nil, "el proceso fue ejecutado exitosamente", c)
		return r.Responder()
	}
}

func PostRetencionCertificados(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request retenciondtos.RetencionCertificadoRequestDTO

		retencionString := c.FormValue("retencion_id")
		valorRetencion, erro := strconv.ParseInt(retencionString, 10, 64)
		if erro != nil {
			return fiber.NewError(400, "Error: no se pudo convertir los datos "+erro.Error())
		}
		valorClienteId, erro := strconv.ParseInt(c.FormValue("cliente_id", "0"), 10, 64)
		if erro != nil {
			return fiber.NewError(400, "Error: no se pudo convertir los datos "+erro.Error())
		}

		request.RetencionId = uint(valorRetencion)
		request.ClienteId = uint(valorClienteId)
		request.Fecha_Presentacion = c.FormValue("fecha_presentacion", "")
		request.Fecha_Caducidad = c.FormValue("fecha_caducidad", "")

		C_R_id, cliente_name, erro := service.ValidarRetencion(request)
		if erro != nil {
			return fiber.NewError(400, "Error: no se pudo validar los datos "+erro.Error())
		}
		request.ClienteRetencionId = C_R_id

		form, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(400, "Error: no se pudo obtener el archivo. "+err.Error())
		}

		ruta := config.DIR_BASE + config.DIR_CERT + "/" + cliente_name
		if _, err := os.Stat(ruta); os.IsNotExist(err) {
			err = os.MkdirAll(ruta, 0755)
			if err != nil {
				return fiber.NewError(400, "Error: no se pudo crear el directorio. "+err.Error())
			}
		}

		retencionsTRING := c.FormValue("tipo_retencion", "")
		split := strings.Split(form.Filename, ".")
		if retencionsTRING != "" {
			split[0] = retencionsTRING
		}
		tipo_retencion := strings.Join(split[0:(len(split))], ".")

		// erro = service.PostRetencionFile(ruta, form)
		milissconds := time.Now().Unix()
		fechaUnix := fmt.Sprintf("%v", milissconds)
		rutaDir := fmt.Sprintf("%s/%s", ruta, (fechaUnix + "-" + tipo_retencion))
		request.RutaFile = "/" + cliente_name + "/" + (fechaUnix + "-" + tipo_retencion)
		erro = c.SaveFile(form, fmt.Sprintf(rutaDir))
		if erro != nil {
			return fiber.NewError(400, "Error: no se pudo crear el directorio. "+err.Error())
		}
		nombreFile := administraciondtos.ArchivoResponse{
			NombreArchivo: (fechaUnix + "-" + tipo_retencion),
		}
		nombreFiles := []administraciondtos.ArchivoResponse{nombreFile}

		_, erro = service.SubirArchivosCloud(c.Context(), ruta, nombreFiles, (config.DIR_CERT + "/" + cliente_name))
		if erro != nil {
			return fiber.NewError(400, "Error: no se pudo crear el directorio. "+err.Error())
		}

		erro = service.PostRetencionesCertificadosService(request)
		if erro != nil {
			return fiber.NewError(400, "Error: no se pudo procesar el archivo. "+erro.Error())
		}
		r := apiresponder.NewResponse(201, nil, "Se proceso correctamente el certificado", c)
		return r.Responder()

	}
}

func GetRetencionCertificados(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionArchivoRequest

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		certificado, err := service.GetCertificadoService(request.CertificadoId)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		err = service.GetCertificadoCloudService(c.Context(), config.DIR_CERT+certificado.Ruta_file)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.LeerContenidoDirectorio(certificado)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron los resultados exitosamente", c)
		return r.Responder()
	}
}

func GetCalcularRetenciones(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.GetCalcularRetencionesService(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron los resultados exitosamente", c)
		return r.Responder()
	}
}

func GetCondiciones(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.GetCondicionesService(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron los resultados exitosamente", c)
		return r.Responder()
	}
}

func DeleteClienteRetencion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		// service
		err = service.DeleteClienteRetencionService(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, nil, "la retención fue eliminada exitosamente", c)
		return r.Responder()
	}
}

func CreateRetencion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request retenciondtos.PostRentencionRequestDTO

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		isUpdate := false

		// service
		response, err := service.CreateRetencionService(request, isUpdate)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		var data []map[string]interface{}
		for _, res := range response {
			res_data := map[string]interface{}{
				"id":                  res.Model.ID,
				"alicuota":            res.Alicuota,
				"alicuota_opcional":   res.AlicuotaOpcional,
				"rg2854":              res.Rg2854,
				"minorista":           res.Minorista,
				"monto_minimo":        res.MontoMinimo,
				"descripcion":         res.Descripcion,
				"codigo_regimen":      res.CodigoRegimen,
				"fecha_validez_desde": res.FechaValidezDesde,
				"fecha_validez_hasta": res.FechaValidezHasta,
			}
			data = append(data, res_data)
		}

		r := apiresponder.NewResponse(201, data, "la retención fue creada exitosamente", c)
		return r.Responder()
	}
}

func UpdateRetencion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request retenciondtos.PostRentencionRequestDTO

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}
		isUpdate := true
		// service
		response, err := service.UpdateRetencionService(request, isUpdate)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		data := map[string]interface{}{
			"id":                  response.Model.ID,
			"alicuota":            response.Alicuota,
			"alicuota_opcional":   response.AlicuotaOpcional,
			"rg2854":              response.Rg2854,
			"minorista":           response.Minorista,
			"monto_minimo":        response.MontoMinimo,
			"descripcion":         response.Descripcion,
			"fecha_validez_desde": response.FechaValidezDesde.Format("2006-01-02"),
			"fecha_validez_hasta": response.FechaValidezHasta.Format("2006-01-02"),
		}

		r := apiresponder.NewResponse(201, data, "la retención fue modificada exitosamente", c)
		return r.Responder()
	}
}

func GetGravamenes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		// filtro vacio para obtener todos los gravamenes
		var filtro retenciondtos.GravamenRequestDTO
		response, err := service.GetGravamenesService(filtro)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron los resultados exitosamente", c)
		return r.Responder()
	}
}

func CreateCondicion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request retenciondtos.CondicionRequestDTO

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		isUpdate := false

		// service
		err = service.UpSertCondicionService(request, isUpdate)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(201, nil, "la condicion fue creada exitosamente", c)
		return r.Responder()
	}
}

func UpdateCondicion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request retenciondtos.CondicionRequestDTO

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		isUpdate := true
		// service
		err = service.UpSertCondicionService(request, isUpdate)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(201, nil, "la condicion fue modificada exitosamente", c)
		return r.Responder()
	}
}

func GetGenerarCertificacion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		// Hacer calculo de retenciones y guardar en la base de datos los registros de comprobantes de retencion y sus detalles
		err = service.GenerarCertificacionService(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, nil, "la operacion se realizó con exito", c)
		return r.Responder()
	}
}

func GetComprobantesRetencionesDevolver(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.ComprobantesRetencionesDevolverService(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "la operacion se realizó con exito", c)
		return r.Responder()
	}
}

func GetComprobantes(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.GetComprobantesService(request)
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Ocurrio un error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "la operacion se realizó con exito", c)
		return r.Responder()
	}
}

func GenerarMovimientosTemporalesPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request filtros.PagoIntentoFiltros
		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(404, "Error: "+err.Error())
		}
		// 1 Se debe consultar los pagos en estado aprobado y procesando
		request = filtros.PagoIntentoFiltros{
			PagoEstadosIds:      []uint64{4, 7},
			CargarPago:          true,
			CargarPagoTipo:      true,
			CargarPagoEstado:    true,
			CargarCuenta:        true,
			PagoIntentoAprobado: true,
			FechaPagoInicio:     request.FechaPagoInicio,
			FechaPagoFin:        request.FechaPagoFin,
		}

		pagos, err := service.GetPagosCalculoMovTemporalesService(request)
		if len(pagos) > 0 {
			// se deben aplicar las comisiones a los pagos: aplicar al ultimo pago intento
			responseCierreLote, err := service.BuildPagosCalculoTemporales(pagos)
			if err == nil {
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				// crear los movimientos temoorales y actualziar campo calculado en pago intento
				// esto inidica que el pago ya fue calculado y guardado en movimientostemporales
				err = service.CreateMovimientosTemporalesService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)
					return fiber.NewError(400, "Error: "+err.Error())
				}

			}

		}

		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo realizar el proceso de cierre de lote de apilink. %s", err),
			}
			service.CreateNotificacionService(notificacion)
		} else if len(pagos) < 1 {
			return c.JSON(&fiber.Map{
				"error":            "no existen pagos por procesar",
				"tipoconciliacion": "GenerarMovimientosTemporalesPagos",
			})
		}

		return c.JSON(&fiber.Map{
			"error":            err,
			"tipoconciliacion": "GenerarMovimientosTemporalesPagos",
		})

	}
}

func GetInformarVencimientosCertificados(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.CertificadoVencimientoDTO

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		err = service.NotificarVencimientoCertificadosService(request)
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Ocurrio un error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, nil, "la operacion se realizó con exito", c)
		return r.Responder()
	}
}
func createComisionManual(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var request administraciondtos.RequestComisionManual
		err := c.BodyParser(&request)
		//Valido que existan los parametros
		if err != nil {
			return fiber.NewError(400, "Error: con los parametros enviados")
		}
		err = service.CreateComisionManualService(request)
		if err != nil {
			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Comision creada correctamente",
		})
	}
}

func CambiarEstadoPagosExpirados(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		err := service.ModificarEstadoPagosExpirados()
		if err != nil {
			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Pagos actualizados correctamente",
		})
	}
}

// esta funcion permite consultar los clientes segun el campo configuracion_retiro_automatico
func getClientesConfiguracion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request filtros.ClienteConfiguracionFiltro

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		clientes, err := service.GetClientesConfiguracionService(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		return c.JSON(&fiber.Map{
			"data":    clientes,
			"status":  true,
			"message": "Clientes cargados correctamente",
		})
	}
}

// func createAuditoria(service administracion.Service) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		c.Accepts("application/json")
// 		var request administraciondtos.RequestAuditoria
// 		err := c.BodyParser(&request)
// 		//Valido que existan los parametros
// 		if err != nil {
// 			return fiber.NewError(400, "Error: con los parametros enviados")
// 		}
// 		err = service.CreateAuditoriaService(request)
// 		if err != nil {
// 			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
// 		}
// 		return c.JSON(&fiber.Map{
// 			"status":        true,
// 			"statusMessage": "datos creados correctamente",
// 		})
// 	}
// }

func NotificarPagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		err := service.NotificarPagosWebhookSinNotificarService()
		if err != nil {
			return fiber.NewError(400, fmt.Sprintf("Error: %v", err.Error()))
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "Pagos actualizados correctamente",
		})
	}
}

// Cierre Lote Multipagos

func ActualizarEstadosPagosClMultipagos(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovMultipagos := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      false,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote  */
		listaPagosMultipago, err := service.GetCierreLoteMultipagosService(filtroMovMultipagos)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos clrapipago. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteRapipago",
			})
		}
		if len(listaPagosMultipago) > 0 {
			if err == nil {
				listaPagosClMultipagos, err := service.BuildPagosClMultipagos(listaPagosMultipago)
				if err == nil {
					// Actualizar estados del pago y cierrelote
					err = service.ActualizarPagosClMultipagosService(listaPagosClMultipagos)
					if err != nil {
						logs.Error(err)
						return fiber.NewError(400, "Error: "+err.Error())
					}

				} else {
					logs.Error(err)
					return fiber.NewError(400, "Error: "+err.Error())
				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("No existen pagos de multipagos para actualizar"),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteMultipagos",
			})
		}

		return c.JSON(&fiber.Map{
			"tipoconciliacion": "CierreLoteMultipagos se ejecuto con exito",
		})

	}
}

func NotificarPagosClMultipagos(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request filtros.PagoEstadoFiltro

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		pagos, barcode, err := service.BuildNotificacionPagosCLMultipagos(request)

		// control de error al llamar al servicio BuildNotificacionPagosCLMultipagos
		if err != nil {
			return c.JSON(&fiber.Map{
				"error":            "error al notificar los pagos: " + err.Error(),
				"tipoconciliacion": "notificacionPagosClMultipagos",
			})
		}

		if len(barcode) > 0 {
			if len(pagos) > 0 {
				pagosNotificar := service.NotificarPagos(pagos)
				if len(pagosNotificar) == 0 {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionWebhook,
						Descripcion: fmt.Sprintln("webhook: no se pudieron notificar los pagos"),
					}
					service.CreateNotificacionService(notificacion)
					return fiber.NewError(400, "Error: "+err.Error())
				}

			}

			//si la notificacion se realiza con exito se debera actualizar el campo pagonotificado en multipagoscierrelote
			err := service.ActualizarPagosClMultipagosDetallesService(barcode)
			if err != nil {
				return c.JSON(&fiber.Map{
					"error":            "error al actualizar pagos en clmultipagos",
					"tipoconciliacion": "notificacionPagosClMultipagos",
				})
			}
			mensaje := fmt.Sprintf("los pagos con los siguientes codigos de barra se actualizaron correctamente. %v", barcode)
			return c.JSON(&fiber.Map{
				"pagos actualizados": mensaje,
				"notificacion":       "exitosa",
				"tipoconciliacion":   "notificacionPagosClMultipagos",
			})

		} else {
			return c.JSON(&fiber.Map{
				"error":            "no existen pagos por notificar",
				"tipoconciliacion": "notificacionPagosClMultipagos",
			})
		}

	}
}

func CierreLoteMultipagos(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovMultipagos := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: false,
			PagosNotificado:      true,
		}
		/* obtener lista pagos rapipago encontrados en el tabla rapipagoscierrelote - los que no fueron conciliados  */
		listaCierreMultipagos, err := service.GetCierreLoteMultipagosService(filtroMovMultipagos)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteMultipagos",
			})
		}
		if len(listaCierreMultipagos) > 0 {
			if err == nil {

				request := bancodtos.RequestConciliacion{
					TipoConciliacion: 4,
					ListaMultipagos:  listaCierreMultipagos,
				}
				// aqui hay retornar la lista de id de repipagocierre lote y los id del banco
				listaCierreMultipagos, listaBancoId, err := banco.ConciliacionPasarelaBanco(request)

				if len(listaBancoId) == 0 {
					logs.Info("no existen movimientos en banco para conciliar con pagos multipagos")
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionCierreLote,
						Descripcion: fmt.Sprintf("no existen movimientos en banco para conciliar con pagos multipagos: %s", err),
					}
					service.CreateNotificacionService(notificacion)
					return c.JSON(&fiber.Map{
						"error":            notificacion.Descripcion,
						"tipoconciliacion": "CierreLoteMultipagos",
					})
				} else {
					/*en el caso de error a actualizar la tabla multipagocierrelote el proceso termina */
					err := service.UpdateCierreLoteMultipagos(listaCierreMultipagos.ListaMultipagos)
					if err != nil {
						logs.Error(err)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("error al actualizar registros de cierrelote multipagos (volver a ejecutar proceso): %s", err),
						}
						service.CreateNotificacionService(notificacion)
					} else {
						// son los registros que coincidieron en el cierrerapipago y banco
						// si no se actualiza los registros del banco se debera actualizar manualmente
						_, err := banco.ActualizarRegistrosMatchBancoService(listaBancoId, true)
						if err != nil {
							logs.Info(fmt.Sprintf("Los siguientes movimientos del banco no se actualizaron: %v", listaBancoId))
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al actualizar movimientos del banco - conciliacion multipagos(actualizar manualmente los siguientes movimientos): %s", err),
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
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteMultipagos",
			})
		}

		return c.JSON(&fiber.Map{
			"tipoconciliacion": "CierreLoteMultipagos se ejecuto con exito",
		})

	}
}

func GenerarMovimientosMultipagos(service administracion.Service, banco banco.BancoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		filtroMovMultipagos := rapipago.RequestConsultarMovimientosRapipago{
			CargarMovConciliados: true,
			PagosNotificado:      true,
		}

		listaCierreMultipagos, err := service.GetCierreLoteMultipagosService(filtroMovMultipagos)
		if err != nil {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("No se pudo obtener los pagos pago conciliar. %s", err),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteMultipagos",
			})
		}
		// Si no se guarda ningún cierre no hace falta seguir el proce
		if len(listaCierreMultipagos) > 0 {
			// 2 - Contruye los movimientos y hace la modificaciones necesarias para modificar los
			// pagos y demás datos necesarios en caso de error se repetira el día siguiente
			responseCierreLote, err := service.BuildMultipagosMovimiento(listaCierreMultipagos)

			if err == nil {

				// 3 - Guarda los movimientos en la base de datos en caso de error se
				// repetira en el día siguiente
				ctx := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
				err = service.CreateMovimientosService(ctx, responseCierreLote)
				if err != nil {
					logs.Error(err)
					return fiber.NewError(400, "Error: "+err.Error())
				} else {
					return c.JSON(&fiber.Map{
						"mensaje":          "conciliacion exitosa",
						"tipoconciliacion": "CierreLoteMultipagos",
					})
				}

			}

		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintln("no existen pagos para generar movimientos"),
			}
			service.CreateNotificacionService(notificacion)
			return c.JSON(&fiber.Map{
				"error":            notificacion.Descripcion,
				"tipoconciliacion": "CierreLoteMultipagos",
			})
		}

		return nil

	}
}

func GetArchivoComprobanteRetencion(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		comprobantes, err := service.GetComprobantesService(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		if len(comprobantes) < 1 {
			r := apiresponder.NewResponse(404, nil, "error: no se encuentran comprobantes con el id requerido", c)
			return r.Responder()
		}
		comprobante := comprobantes[0]
		err = service.GetCertificadoCloudService(c.Context(), config.DIR_COMP_RETENCIONES+comprobante.RutaFile)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.LeerContenidoComprobanteRetencion(comprobante)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron los resultados exitosamente", c)
		return r.Responder()
	}
}

func GetArchivoReporteRendicionMensual(service administracion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request retenciondtos.RentencionRequestDTO

		err := c.QueryParser(&request)

		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		if len(request.RutaFile) < 1 {
			r := apiresponder.NewResponse(400, nil, "error: debe enviar una ruta del archivo", c)
			return r.Responder()
		}

		if request.ReporteId < 1 {
			r := apiresponder.NewResponse(400, nil, "error: debe enviar un id de reporte", c)
			return r.Responder()
		}

		err = service.GetCertificadoCloudService(c.Context(), config.DIR_COMP_RETENCIONES+request.RutaFile)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		response, err := service.LeerContenidoReporteRendicionMensual(request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, response, "se obtuvieron los resultados exitosamente", c)
		return r.Responder()
	}
}
