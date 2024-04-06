package routes

import (
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/api/middlewares/middlewareinterno"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/email"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/reportes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/reportedtos"
	dtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/reportedtos"
	apiresponder "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/responderdto"
	"github.com/gofiber/fiber/v2"
)

func ReporteRoutes(app fiber.Router, middlewares middlewares.MiddlewareManager, middlewaresInterno middlewareinterno.MiddlewareManagerPasarela, service reportes.ReportesService, email email.Emailservice, runEndpoint util.RunEndpoint) {

	app.Get("/send-pagos", middlewares.ValidarPermiso("psp.herramienta"), sendPagos(service))
	app.Get("/send-rendiciones", middlewares.ValidarPermiso("psp.herramienta"), sendRendiciones(service))
	app.Get("/send-reversiones", middlewares.ValidarPermiso("psp.herramienta"), sendReversiones(service))
	app.Get("/batch-pago-items", middlewares.ValidarPermiso("psp.herramienta"), batchPagoItems(service))
	/* endpoints de consulta de informacion */
	app.Get("/logs", middlewares.ValidarPermiso("psp.herramienta"), reportesLogs(service))
	app.Get("/notificaciones", middlewares.ValidarPermiso("psp.herramienta"), reportesNotificaciones(service))
	app.Get("/reportes-enviados", middlewares.ValidarPermiso("psp.herramienta"), reportesEnviados(service))
	app.Get("/reporte-cobranza-diaria", middlewares.ValidarPermiso("psp.herramienta"), sendCobranzasDiarias(service))
	app.Get("/reportes-pagos-mensuales", middlewares.ValidarPermiso("psp.herramienta.reporte"), consultaPagosMensuales(service))
	app.Get("/reportes-rendiciones-mensuales", middlewares.ValidarPermiso("psp.herramienta.reporte"), consultaRendicionesMensuales(service))
	app.Post("/reporte-retencion-comprobante", middlewares.ValidarPermiso("psp.consultar.impuestos"), reporteComprobanteRetencion(service))
	// reporte mensual pdf de rendiciones
	// app.Get("/send-reportes-rendiciones", middlewares.ValidarPermiso("psp.herramienta"), sendReportesRendiciones(service))
	app.Post("/send-archivo-email", middlewares.ValidarPermiso("psp.herramienta.reporte"), sendArchivoEmail(service, email))

	// Route Ejecuta todo el proceso de liquidacion de retenciones para un cliente
	app.Get("/liquidar-retenciones", middlewares.ValidarPermiso("psp.consultar.impuestos"), LiquidarRetenciones(service))
	// Route Ejecuta todo el proceso de creacion de archivos txt de retenciones
	app.Post("/create-txt-retenciones", middlewares.ValidarPermiso("psp.consultar.impuestos"), CreateTxtRetenciones(service))

	// Route Creacion de archivos txt de comisiones form8125
	app.Post("/create-txt-comisiones", CreateTxtComisiones(service))

	// batch pagos
	app.Get("/batch-pagos", batchPagos(service))

	// Verifica que los montos que devuelven distintos endpoints de cobranzas
	app.Get("/verificar-cobranzas", verificarCobranzasClienteCobranzas(service, runEndpoint))
	app.Post("/create-excel-retenciones", CreateExcelRetenciones(service))
}

func verificarCobranzasClienteCobranzas(service reportes.ReportesService, runEndpoint util.RunEndpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := c.Get("Authorization")

		var query reportedtos.RequestControlReportes

		var response []reportedtos.ResponseControlReporte

		_ = c.QueryParser(&query)

		apikeys := strings.Join(query.ApiKey, ",")

		response, err := service.MakeControlReportes(apikeys, query.FechaConsultar, token, runEndpoint)
		if err != nil {
			fmt.Println("Error al enviar email", err.Error())
		}

		err = service.SendControlReportes(response)
		if err != nil {
			fmt.Println("Error al enviar email", err.Error())
		}

		return c.Status(400).JSON(&fiber.Map{
			"status":  true,
			"data":    response,
			"message": "Reportes control enviado con éxito",
		})
	}
}

func sendPagos(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		clientes, err := service.GetClientes(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente
				listaPagosClientes, err := service.GetPagosClientes(clientes, request)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}

				// enviar los pagos a clientes
				if len(listaPagosClientes) > 0 {
					listaErro, err := service.SendPagosClientes(listaPagosClientes)
					if err != nil {
						r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
						return r.Responder()
					} else if len(listaErro) > 0 {
						r := apiresponder.NewResponse(400, listaErro, "error al intentar enviar reportes a clientes", c)
						return r.Responder()
					} else {
						// caso de exito del proceso
						r := apiresponder.NewResponse(200, nil, "el proceso se ejecuto con éxito", c)
						return r.Responder()
					}
				} else {
					// caso de no existen pagos para reportar
					r := apiresponder.NewResponse(404, "tipoconciliacion: notificacionPagos", "no existen pagos por enviar", c)
					return r.Responder()
				}
			}

		}
		// CASO DE QUE NO SE ENCUENTRAN LOS CLIENTES
		r := apiresponder.NewResponse(404, "tipoconciliacion: ReportesPagosEnviados", "no existen clientes: verifique datos enviados", c)
		return r.Responder()
	}
}

func sendRendiciones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		// 1 obtener lista de cliente
		clientes, err := service.GetClientes(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener transferencias para rendicion por cliente. Creacion de reporte en PDF
				listaRendicionClientes, err := service.GetRendicionClientes(clientes, request)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}
				// enviar los reporte de rendiciones a los clientes
				if len(listaRendicionClientes) > 0 {
					listaErro, err := service.SendPagosClientes(listaRendicionClientes)
					if err != nil {
						r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
						return r.Responder()
					} else if len(listaErro) > 0 {
						r := apiresponder.NewResponse(400, listaErro, "error al intentar enviar reportes a clientes", c)
						return r.Responder()
					}
				} else {
					// caso sin datos que enviar
					r := apiresponder.NewResponse(404, "tipoconciliacion: notificacionPagos", "no existen rendiciones por enviar", c)
					return r.Responder()
				}
			}

		}
		// caso de exito
		r := apiresponder.NewResponse(200, "tipoconciliacion: ReportesRendicionesEnviados", "el proceso se ejecuto con éxito", c)
		return r.Responder()
	}
}

func sendReversiones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var req struct {
			Cliente     uint
			FechaInicio string
		}
		err := c.BodyParser(&req)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: parametros recibidos no son validos", c)
			return r.Responder()
		}

		var fechaInicio time.Time

		if len(req.FechaInicio) > 0 {
			fechaInicio, err = time.Parse("2006-01-02", req.FechaInicio)
			if err != nil {
				r := apiresponder.NewResponse(400, nil, "Error: formato de fecha inválido", c)
				return r.Responder()
			}
		}

		request := dtos.RequestPagosClientes{
			Cliente:     req.Cliente,
			FechaInicio: fechaInicio,
			FechaFin:    fechaInicio,
		}

		// se valida que las fechas no sean nulas y que la fecha este antes de la fecha fin o sean iguales
		filtro, err := request.ValidarFechas()
		if err != nil {
			r := apiresponder.NewResponse(400, nil, fmt.Sprintf("Error %v", err.Error()), c)
			return r.Responder()
		}
		logs.Info(filtro)
		// 1 obtener lista de cliente
		clientes, err := service.GetClientes(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente
				listaReversionesClientes, err := service.GetReversionesClientes(clientes, request)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}
				if len(listaReversionesClientes) <= 0 {
					data := map[string]interface{}{
						"status":  true,
						"message": "no exiten reversiones",
					}
					// caso sin resultados
					r := apiresponder.NewResponse(404, data, "no exiten reversiones", c)
					return r.Responder()
				}

				// enviar los reverciones a clientes
				listaErro, err := service.SendPagosClientes(listaReversionesClientes)
				logs.Info(listaErro)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "error al intentar enviar reportes a clientes", c)
					return r.Responder()
				}
				if len(listaErro) > 0 {
					r := apiresponder.NewResponse(400, listaErro, "error al intentar enviar reportes a clientes", c)
					return r.Responder()
				}
			}

		}

		data := map[string]interface{}{
			"status":           true,
			"message":          "success",
			"tipoconciliacion": "ReportesReversionesEnviados",
		}
		// caso de exito
		r := apiresponder.NewResponse(200, data, "el proceso se ejecuto con exito", c)
		return r.Responder()
	}
}

// Crea, guarda y envia por email, el reporte de rendiciones mensual (rrm)
// func sendReportesRendiciones(service reportes.ReportesService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		c.Accepts("application/json")
// 		var request reportedtos.RequestReportesEnviados

// 		err := c.QueryParser(&request)
// 		if err != nil {
// 			r := apiresponder.NewResponse(400, nil, "Error: parametros recibidos no son validos", c)
// 			return r.Responder()
// 		}

// 		listaReportesCliente, err := service.GetReportesPdfService(request)
// 		if err != nil {
// 			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
// 			return r.Responder()
// 		}
// 		if len(listaReportesCliente) <= 0 {
// 			data := map[string]interface{}{
// 				"status":  true,
// 				"message": "no exiten reportes de rendiciones",
// 			}
// 			// caso sin resultados
// 			r := apiresponder.NewResponse(404, data, "no exiten reversiones", c)
// 			return r.Responder()
// 		}
// 		// enviar el reporte
// 		// listaErro, _, err := service.SendReporteRendiciones(listaReportesCliente)
// 		// logs.Info(listaErro)
// 		if err != nil {
// 			r := apiresponder.NewResponse(400, nil, "error al intentar enviar reportes a clientes", c)
// 			return r.Responder()
// 		}
// 		// if len(listaErro) > 0 {
// 		// 	r := apiresponder.NewResponse(400, listaErro, "error al intentar enviar reportes a clientes", c)
// 		// 	return r.Responder()
// 		// }

// 		data := map[string]interface{}{
// 			"status":         true,
// 			"message":        "success",
// 			"reporteEnviado": "ReporteMensualRendiciones",
// 		}
// 		// caso de exito
// 		r := apiresponder.NewResponse(200, data, "el proceso se ejecuto con exito", c)
// 		return r.Responder()
// 	}
// }

func batchPagoItems(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		ctx := getContextAuditable(c)
		// 1 obtener lista de cliente
		clientes, err := service.GetClientes(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {
				// obtener los pagos/pagoitems por cliente
				listaPagosItems, err := service.GetPagoItems(clientes, request)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}
				// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo
				if len(listaPagosItems) > 0 {
					resultpagositems := service.BuildPagosItems(listaPagosItems)
					if len(resultpagositems) > 0 {
						err := service.ValidarEsctucturaPagosItems(resultpagositems) // validar estructura antes de crear el archivo
						if err != nil {
							r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
							return r.Responder()

						} else {
							err := service.SendPagosItems(ctx, resultpagositems, request)
							if err != nil {
								r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
								return r.Responder()
							}
							// caso de exito
							r := apiresponder.NewResponse(200, "tipoconciliacion: batchPagoItems", "Reportes batch se ejecuto con éxito", c)
							return r.Responder()
						}
					}
				} else {
					// caso sin resultados
					r := apiresponder.NewResponse(404, "tipoconciliacion: batchPagoItems", "no existen pagos para informar", c)
					return r.Responder()
				}
			}

		}
		// caso error en la busqueda de clientes
		r := apiresponder.NewResponse(404, "tipoconciliacion: batchPagoItems", "no existen clientes para informar reporte", c)
		return r.Responder()
	}
}

/* endpoints de consulta de informacion */

func reportesLogs(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestLogs

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		_, err = request.ValidarFechas()
		if err != nil {
			return fiber.NewError(400, "Error en la validación de los parámetros enviados: "+err.Error())
		}

		res, err := service.GetLogs(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": "Sin resultados encontrados para la fecha",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud logs generada",
		})

	}
}

func reportesNotificaciones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestNotificaciones

		err := c.QueryParser(&request)
		if err != nil {
			return fiber.NewError(400, "Error en los parámetros enviados: "+err.Error())
		}

		_, err = request.ValidarFechas()
		if err != nil {
			return fiber.NewError(400, "Error en la validación de los parámetros enviados: "+err.Error())
		}

		res, err := service.GetNotificaciones(request)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}

		if len(res.Data) == 0 {
			return c.Status(200).JSON(&fiber.Map{
				"status": true,
				"result": "Sin resultados encontrados para la fecha",
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"result":  res,
			"message": "Solicitud notificaciones generada",
		})

	}
}

func reportesEnviados(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		// filtro de la request
		var filtroReportesEnviados dtos.RequestReportesEnviados

		// parse de los parametros de la request al filtro CierreLoteFiltro
		err := c.QueryParser(&filtroReportesEnviados)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en los parametros recibidos "+err.Error(), c)
			return r.Responder()
		}

		// Validar los parametros
		err = filtroReportesEnviados.Validar()

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "error en validación de parámetros recibidos: "+err.Error(), c)
			return r.Responder()
		}

		// Enviar la consulta al servicio correspondiente
		result, err := service.GetReportesEnviadosService(filtroReportesEnviados)

		// si hubo un error devolver un map con mensaje de error y nil en data
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Error "+err.Error(), c)
			return r.Responder()
		}

		// si no hubo resultados en la consulta, pero tampoco errores, devolver en data un string vacio
		if len(result.Reportes) == 0 {

			r := apiresponder.NewResponse(200, []string{}, "Datos de consulta enviados, sin resultados", c)
			return r.Responder()

		}

		r := apiresponder.NewResponse(200, result, "Datos de consulta enviados", c)
		return r.Responder()
	}
}

func sendCobranzasDiarias(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestCobranzasDiarias
		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}
		request_cliente := dtos.RequestPagosClientes{}

		clientes, err := service.GetClientes(request_cliente)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente
				listaPagosClientes, err := service.GetPagos(clientes, request)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}

				// enviar los pagos a clientes
				if len(listaPagosClientes) > 0 {
					listaErro, err := service.SendPagosClientes(listaPagosClientes)
					if err != nil {
						r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
						return r.Responder()
					} else if len(listaErro) > 0 {
						r := apiresponder.NewResponse(400, listaErro, "error al intentar enviar reportes a clientes", c)
						return r.Responder()
					} else {
						// caso de exito del proceso
						r := apiresponder.NewResponse(200, nil, "el proceso se ejecuto con éxito", c)
						return r.Responder()
					}
				} else {
					// caso de no existen pagos para reportar
					r := apiresponder.NewResponse(404, "tipoconciliacion: notificacionPagos", "no existen pagos por enviar", c)
					return r.Responder()
				}
			}

		}
		// CASO DE QUE NO SE ENCUENTRAN LOS CLIENTES
		r := apiresponder.NewResponse(404, "tipoconciliacion: ReportesPagosEnviados", "no existen clientes: verifique datos enviados", c)
		return r.Responder()
	}
}

func consultaPagosMensuales(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		clientes, err := service.GetClientes(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente
				listaPagosClientes, err := service.GetPagosClientesMensual(clientes, request)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}

				reporteRespuesta, err := service.TratamientoReporteMensualPagos(listaPagosClientes, request.OrdenMayorCobranza)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}

				// caso de exito del proceso
				r := apiresponder.NewResponse(200, reporteRespuesta, "el proceso se ejecuto con éxito", c)
				return r.Responder()
			}

		}
		// CASO DE QUE NO SE ENCUENTRAN LOS CLIENTES
		r := apiresponder.NewResponse(404, "tipoconciliacion: ReportesPagosEnviados", "no existen clientes: verifique datos enviados", c)
		return r.Responder()
	}
}

func consultaRendicionesMensuales(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		clientes, err := service.GetClientes(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {

				// obtener los pagos por cliente
				listaRendicionesClientes, err := service.GetRendicionesClientesMensual(clientes, request)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}

				reporteRespuesta, err := service.TratamientoReporteMensualRendiciones(listaRendicionesClientes, request.OrdenMayorCobranza)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}

				// caso de exito del proceso
				r := apiresponder.NewResponse(200, reporteRespuesta, "el proceso se ejecuto con éxito", c)
				return r.Responder()
			}

		}
		// CASO DE QUE NO SE ENCUENTRAN LOS CLIENTES
		r := apiresponder.NewResponse(404, "tipoconciliacion: ReportesRendicionesEnviados", "no existen clientes: verifique datos enviados", c)
		return r.Responder()
	}
}

func reporteComprobanteRetencion(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request reportedtos.RequestRRComprobante

		err := c.BodyParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		res, err := service.SendReporteRetencionComprobante(request)
		if err != nil {
			r := apiresponder.NewResponse(404, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		if !res {
			r := apiresponder.NewResponse(200, true, "Sin resultados encontrados", c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, true, "Reporte comprobante de retenciones enviado", c)
		return r.Responder()
	}
}

func sendArchivoEmail(service reportes.ReportesService, email email.Emailservice) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body, contentType, tempDirectoryPath, erro := email.CreateEmailService(c)
		if erro != nil {
			r := apiresponder.NewResponse(400, nil, "Ocurrio un error al crear el email: "+erro.Error(), c)
			return r.Responder()
		}
		erro = email.SendEmailService(body, contentType, tempDirectoryPath)
		if erro != nil {
			r := apiresponder.NewResponse(400, nil, "Ocurrio un error al enviar el email: "+erro.Error(), c)
			return r.Responder()
		}
		r := apiresponder.NewResponse(200, nil, "El correo fue enviado con exito", c)
		return r.Responder()
	}
}

func LiquidarRetenciones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request reportedtos.RequestReportesEnviados

		err := c.QueryParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		err = service.LiquidarRetencionesService(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		r := apiresponder.NewResponse(200, nil, "la operacion liquidacion de retenciones se realizo con exito", c)
		return r.Responder()
	}
}

func CreateTxtRetenciones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request reportedtos.RequestReportesEnviados

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		// Service
		err = service.CreateTxtRetencionesService(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		// Respuesta
		r := apiresponder.NewResponse(200, nil, "la operacion se realizó con exito", c)
		return r.Responder()
	}
}

func CreateTxtComisiones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request reportedtos.RequestReportesEnviados

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		// Service
		_, err = service.CreateTxtForm8125Service(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		// Respuesta
		r := apiresponder.NewResponse(200, nil, "la operacion se realizó con exito", c)
		return r.Responder()
	}
}

func batchPagos(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		var request dtos.RequestPagosClientes
		err := c.QueryParser(&request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		ctx := getContextAuditable(c)
		// 1 obtener lista de cliente que tenga habilitado envio de batch pago
		// request.ReporteBatchPagos = true
		clientes, err := service.GetClientes(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}
		if len(clientes.Clientes) > 0 {
			if err == nil {
				// obtener los pagos/pagoitems por cliente
				listaPagosItems, err := service.GetPagoItemsAlternativo(clientes, request)
				if err != nil {
					r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
					return r.Responder()
				}
				// si la lista de pagos para informar es mayor a 0 se genera se sigue con el proceso de construir la estructura correspondiente al archivo

				// if len(listaPagosItems) == 0 {

				// 	fechaI := time.Now()
				// 	fechaI = fechaI.AddDate(0, 0, int(-1))

				// 	fecha := commons.ConvertFechaString(fechaI) // fecha de creacion del archivo
				// 	clienteInforme := clientes.Clientes[0]
				// 	clienteRes := dtos.ClientesResponse{
				// 		Id:          clienteInforme.Id,
				// 		Cliente:     clienteInforme.Cliente,
				// 		RazonSocial: clienteInforme.NombreFantasia,
				// 		Email:       clienteInforme.Email,
				// 	}
				// 	pagosItemsVacio := dtos.ResponsePagosItems{
				// 		Clientes: clienteRes,
				// 		Fecha:    fecha,
				// 	}

				// 	listaPagosItems = append(listaPagosItems, pagosItemsVacio)
				// }
				if len(listaPagosItems) > 0 {
					resultpagositems := service.BuildPagosArchivo(listaPagosItems)
					if len(resultpagositems) > 0 {
						err := service.ValidarEsctucturaPagosBatch(resultpagositems) // validar estructura antes de crear el archivo
						if err != nil {
							r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
							return r.Responder()

						} else {
							err := service.SendPagosBatch(ctx, resultpagositems, request)
							if err != nil {
								r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
								return r.Responder()
							}
							// caso de exito
							r := apiresponder.NewResponse(200, "tipoconciliacion: batchPagoItems", "Reportes batch se ejecuto con éxito", c)
							return r.Responder()
						}
					}
				} else {
					// caso sin resultados
					r := apiresponder.NewResponse(404, "tipoconciliacion: batchPagoItems", "no existen pagos para informar", c)
					return r.Responder()
				}
			}

		}
		// caso error en la busqueda de clientes
		r := apiresponder.NewResponse(404, "tipoconciliacion: batchPagoItems", "no existen clientes para informar reporte", c)
		return r.Responder()
	}
}

func CreateExcelRetenciones(service reportes.ReportesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request reportedtos.RequestReportesEnviados

		err := c.BodyParser(&request)

		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error en los parámetros enviados: "+err.Error(), c)
			return r.Responder()
		}

		// Service
		_, err = service.CreateExcelRetencionesService(request)
		if err != nil {
			r := apiresponder.NewResponse(400, nil, "Error: "+err.Error(), c)
			return r.Responder()
		}

		// Respuesta
		r := apiresponder.NewResponse(200, nil, "la operacion se realizó con exito", c)
		return r.Responder()
	}
}
