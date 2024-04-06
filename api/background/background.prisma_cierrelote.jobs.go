package background

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/cierrelote"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/administracion"
	filtroCl "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/cierrelote"
	"github.com/robfig/cron"
)

func GetArchivosTxtCierreLote(cronjob *cron.Cron, periodicidad string, cierrelote cierrelote.Service, service administracion.Service) {

	// definir job
	var getArchivosTxtCierreLote = func() {
		ctxaws := context.Background()
		/*
			PROCESO: Diariamente el proceso de cierre lote en backproud se divide en 4 pasos.
				paso 1:
					leer el directorio ftp de cierre de lote minio y obtener informacion de los archivos
					y se guarda en un directorio temporal los archivos txt existentes
					se obtiene todos los estados externos de prisma
				paso 2:
					se recorren uno a uno los archivos de cierre de lotes y se almacena a la bd
				paso 3:
					se mueven todos los archivos de la carpeta temporal al minio.
				paso 4:
					por ultimo se borran todos los archivos creados temporalmente y el directorio temporal
		*/
		///////////////////////////////////////PROCESO///////////////////////////////////////
		/* paso 1: */
		archivos, rutaArchivos, totalArchivos, err := cierrelote.LeerArchivoLoteExterno(ctxaws, config.DIR_KEY)
		if err == nil {
			if totalArchivos != 0 {
				/* se obtiene todos los estados externos de prisma */
				filtro := filtros.PagoEstadoExternoFiltro{
					Vendor:           strings.ToUpper("prisma"),
					CargarEstadosInt: true,
				}
				estadosPagoExterno, err := service.GetPagosEstadosExternoService(filtro)
				if err != nil {
					errObtenerEstados := errors.New("error al solicitar lista de estados de pago")
					err = errObtenerEstados
					logError := entities.Log{
						Tipo:          entities.EnumLog("error"),
						Funcionalidad: "GetPagosEstadosExternoService",
						Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
					}
					errCrearLog := service.CreateLogService(logError)
					if errCrearLog != nil {
						logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
					}
				} else {
					/* paso 2: */
					listaArchivo, err := cierrelote.LeerCierreLoteTxt(archivos, rutaArchivos, estadosPagoExterno)
					if err != nil {
						logs.Error(err)
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("No se pudo realizar el cierre de lote de prisma. %s", err),
						}
						service.CreateNotificacionService(notificacion)
					}
					if len(listaArchivo) != 0 {
						/* paso 3: */
						countArchivos, err := cierrelote.MoverArchivos(ctxaws, rutaArchivos, listaArchivo)
						if err != nil {
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al borrar los archivos temporales: %s", err),
							}
							service.CreateNotificacionService(notificacion)
						}
						/* paso 4: */

						err = cierrelote.BorrarArchivos(ctxaws, config.DIR_KEY, rutaArchivos, listaArchivo)
						if err != nil {
							logs.Error(err)
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("error al borrar los archivos temporales: %s", err),
							}
							service.CreateNotificacionService(notificacion)
						}
						var notificacion entities.Notificacione
						if countArchivos > 0 {
							notificacion = entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("fecha: %v - Se procesaron %v archivos de cierre de lote Prisma, recibido por: SFTP", time.Now().String(), countArchivos),
							}
						} else {
							notificacion = entities.Notificacione{
								Tipo:        entities.NotificacionCierreLote,
								Descripcion: fmt.Sprintf("fecha: %v - No existe movimientos de cierre de lote Prisma: %s FTP", time.Now().String(), err),
							}
						}
						service.CreateNotificacionService(notificacion)
					} else {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionCierreLote,
							Descripcion: fmt.Sprintf("fecha: %v - No existe archivos de cierre de lote Prisma: %s FTP", time.Now().String(), err),
						}
						service.CreateNotificacionService(notificacion)
					}
				}
			}
		} else {
			notificacion := entities.Notificacione{
				Tipo:        entities.NotificacionCierreLote,
				Descripcion: fmt.Sprintf("fecha: %v - No se pudo realizar cierre de lote Prisma: %s FTP", time.Now().String(), err),
			}
			service.CreateNotificacionService(notificacion)
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getArchivosTxtCierreLote)
}

func GetTablaMovimientosMXCierreLote(cronjob *cron.Cron, periodicidad string, cierrelote cierrelote.Service, service administracion.Service) {

	// definir job
	var getTablaMovimientosMXCierreLote = func() {

		movimientoMx, movimientoMxEntity, err := cierrelote.ObtenerMxMoviminetosServices()
		if err != nil {
			errObtenerEstados := errors.New("error al obtener registros de la tablas movimientos mx")
			err = errObtenerEstados
			logError := entities.Log{
				Tipo:          entities.EnumLog("error"),
				Funcionalidad: "ObtenerMxMoviminetosServices",
				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
			}
			errCrearLog := service.CreateLogService(logError)
			if errCrearLog != nil {
				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
			}
		} else {
			tablasRelaciondas, err := cierrelote.ObtenerTablasRelacionadasServices()
			if err != nil {
				errObtenerEstados := errors.New("al obtener las tablas relacionadas" + err.Error())
				err = errObtenerEstados
				logError := entities.Log{
					Tipo:          entities.EnumLog("error"),
					Funcionalidad: "ObtenerTablasRelacionadasServices",
					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
				}
				errCrearLog := service.CreateLogService(logError)
				if errCrearLog != nil {
					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
				}
			} else {
				resultadoMovimientoMx := cierrelote.ProcesarMovimientoMxServices(movimientoMx, tablasRelaciondas)
				if len(resultadoMovimientoMx) <= 0 {
					errObtenerEstados := errors.New("error: procesar movimiento mx se encuentra vacia")
					err = errObtenerEstados
					logError := entities.Log{
						Tipo:          entities.EnumLog("error"),
						Funcionalidad: "ProcesarMovimientoMxServices",
						Mensaje:       errObtenerEstados.Error(),
					}
					errCrearLog := service.CreateLogService(logError)
					if errCrearLog != nil {
						logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
					}
				} else {
					err = cierrelote.SaveMovimientoMxServices(resultadoMovimientoMx, movimientoMxEntity)
					if err != nil {
						errObtenerEstados := errors.New("error al guardar los movimientos: " + err.Error())
						err = errObtenerEstados
						logError := entities.Log{
							Tipo:          entities.EnumLog("error"),
							Funcionalidad: "SaveMovimientoMxServices",
							Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
						}
						errCrearLog := service.CreateLogService(logError)
						if errCrearLog != nil {
							logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
						}
					} else {
						notificacion := entities.Notificacione{
							Tipo:        entities.NotificacionProcesoMx,
							Descripcion: fmt.Sprintf("fecha : %v - procesamiento de movimientos mx se realizo con exito", time.Now().String()),
						}
						service.CreateNotificacionService(notificacion)
					}
				}
			}
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getTablaMovimientosMXCierreLote)
}

func GetPagosPXCierreLote(cronjob *cron.Cron, periodicidad string, cierrelote cierrelote.Service, service administracion.Service) {

	// definir job
	var getPagosPXCierreLote = func() {

		pagosPx, entityPagoPxStr, err := cierrelote.ObtenerPxPagosServices()
		if err != nil {
			errObtenerEstados := errors.New("error al obtener registros de la tablas pagos px")
			err = errObtenerEstados
			logError := entities.Log{
				Tipo:          entities.EnumLog("error"),
				Funcionalidad: "ObtenerPxPagosServices",
				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
			}
			errCrearLog := service.CreateLogService(logError)
			if errCrearLog != nil {
				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
			}
		} else {
			err = cierrelote.SavePagoPxServices(pagosPx, entityPagoPxStr)
			if err != nil {
				errObtenerEstados := errors.New("error al guardar liquidacion de prisma")
				err = errObtenerEstados
				logError := entities.Log{
					Tipo:          entities.EnumLog("error"),
					Funcionalidad: "SavePagoPxServices",
					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
				}
				errCrearLog := service.CreateLogService(logError)
				if errCrearLog != nil {
					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
				}
			} else {
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionProcesoPx,
					Descripcion: fmt.Sprintf("fecha : %v - procesamiento de pagos px se realizo con exito", time.Now().String()),
				}
				service.CreateNotificacionService(notificacion)
			}
		}
	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getPagosPXCierreLote)
}

func GetConciliarCLMX(cronjob *cron.Cron, periodicidad string, cierrelote cierrelote.Service, service administracion.Service) {
	var getConciliarCLMX = func() {

		filtro_clmx_compras := filtroCl.FiltroCierreLote{
			MatchCl:       true,
			MovimientosMX: true,
			PagosPx:       true,
			Banco:         true,
			Compras:       true,
			Devolucion:    false,
			Anulacion:     false,
			ContraCargo:   false,
			ContraCargoMx: false,
			ContraCargoPx: false,
			Reversion:     false,
		}

		filtro_clmx_devolucion := filtroCl.FiltroCierreLote{
			MatchCl:       true,
			MovimientosMX: true,
			PagosPx:       true,
			Banco:         true,
			Compras:       false,
			Devolucion:    true,
			Anulacion:     false,
			ContraCargo:   false,
			ContraCargoMx: false,
			ContraCargoPx: false,
			Reversion:     false,
		}

		filtro_clmx_anulacion := filtroCl.FiltroCierreLote{
			MatchCl:       true,
			MovimientosMX: true,
			PagosPx:       true,
			Banco:         true,
			Compras:       false,
			Devolucion:    false,
			Anulacion:     true,
			ContraCargo:   false,
			ContraCargoMx: false,
			ContraCargoPx: false,
			Reversion:     false,
		}

		filtro_clmx_cc := filtroCl.FiltroCierreLote{
			MatchCl:       false,
			MovimientosMX: false,
			PagosPx:       false,
			Banco:         false,
			Compras:       true,
			Devolucion:    false,
			Anulacion:     false,
			ContraCargo:   true,
			ContraCargoMx: false,
			ContraCargoPx: false,
			Reversion:     false,
		}

		var filtros_automaticos []filtroCl.FiltroCierreLote

		filtros_automaticos = append(filtros_automaticos, filtro_clmx_compras, filtro_clmx_devolucion, filtro_clmx_anulacion, filtro_clmx_cc)
		filtros_nombres := []string{"filtro compras", "filtro devolucion", "filtro anulacion", "filtro contracargo"}

		for i_loop, filtro_loop := range filtros_automaticos {
			filtro_aplicado := filtro_loop

			var filtro filtroCl.FiltroPrismaMovimiento
			filtro = filtroCl.FiltroPrismaMovimiento{
				Match:                        false,
				CargarDetalle:                true,
				Contracargovisa:              true,
				Contracargomaster:            true,
				Tipooperacion:                true,
				Rechazotransaccionprincipal:  true,
				Rechazotransaccionsecundario: true,
				Motivoajuste:                 true,
				ContraCargo:                  filtro_aplicado.ContraCargo,
				CodigosOperacion:             []string{"0005"},
				TipoAplicacion:               "+",
			}
			if filtro_aplicado.ContraCargo {
				filtro.Match = false
				filtro.CodigosOperacion = []string{"1507", "6000", "1517"}
				filtro.TipoAplicacion = "-"
			}
			listaMovimientos, codigoautorizacion, err := cierrelote.ObtenerPrismaMovimientosServices(filtro)
			if err != nil {
				errObtenerEstados := errors.New(err.Error())
				err = errObtenerEstados
				logError := entities.Log{
					Tipo:          entities.EnumLog("error"),
					Funcionalidad: "ObtenerPrismaMovimientosServices",
					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
				}
				errCrearLog := service.CreateLogService(logError)
				if errCrearLog != nil {
					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
				}
			} else {

				listaCierreLote, err := cierrelote.ObtenerCierreloteServices(filtro_aplicado, codigoautorizacion)
				if err != nil {
					errObtenerEstados := errors.New(err.Error())
					err = errObtenerEstados
					logError := entities.Log{
						Tipo:          entities.EnumLog("error"),
						Funcionalidad: "ObtenerCierreloteServices",
						Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
					}
					errCrearLog := service.CreateLogService(logError)
					if errCrearLog != nil {
						logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
					}
				} else {

					listaCierreLoteProcesado, listaIdsDetalle, listaIdsCabecera, err := cierrelote.ConciliarCierreLotePrismaMovimientoServices(listaCierreLote, listaMovimientos, false)
					if err != nil {
						errObtenerEstados := errors.New(err.Error())
						err = errObtenerEstados
						logError := entities.Log{
							Tipo:          entities.EnumLog("error"),
							Funcionalidad: "ConciliarCierreLotePrismaMovimientoServices",
							Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
						}
						errCrearLog := service.CreateLogService(logError)
						if errCrearLog != nil {
							logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
						}
					} else {

						if len(listaCierreLoteProcesado) <= 0 && len(listaIdsDetalle) <= 0 && len(listaIdsCabecera) <= 0 {
							notificacion := entities.Notificacione{
								Tipo:        entities.NotificacionConciliacionCLMx,
								Descripcion: fmt.Sprintf("fecha : %v - no existe cierre de lotes para conciliar con movimientos %v", time.Now().String(), filtros_nombres[i_loop]),
							}
							service.CreateNotificacionService(notificacion)
						} else {

							logs.Info("en end-point")
							logs.Info(listaCierreLoteProcesado)
							err = cierrelote.ActualizarCierreloteMoviminetosServices(listaCierreLoteProcesado, listaIdsCabecera, listaIdsDetalle)
							if err != nil {
								errObtenerEstados := errors.New(err.Error())
								err = errObtenerEstados
								logError := entities.Log{
									Tipo:          entities.EnumLog("error"),
									Funcionalidad: "ActualizarCierreloteMoviminetosServices",
									Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
								}
								errCrearLog := service.CreateLogService(logError)
								if errCrearLog != nil {
									logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
								}
							} else {
								notificacion := entities.Notificacione{
									Tipo:        entities.NotificacionConciliacionCLMx,
									Descripcion: fmt.Sprintf("fecha : %v - proceso de conciliacion movimientos con cierre lote exito  %v", time.Now().String(), filtros_nombres[i_loop]),
								}
								service.CreateNotificacionService(notificacion)
							}
						}
					}
				}
			}
		}

	}

	// add job to cron
	cronjob.AddFunc(periodicidad, getConciliarCLMX)
}

func GetConciliarCLPX(cronjob *cron.Cron, periodicidad string, cierrelote cierrelote.Service, service administracion.Service) {
	var getConciliarCLPX = func() {
		filtro_clpx_compra := filtroCl.FiltroCierreLote{
			MatchCl:         true,
			MovimientosMX:   false,
			PagosPx:         true,
			Banco:           true,
			EstadoFechaPago: true,
			FechaPago:       "0000-00-00",
			Compras:         true,
			Devolucion:      false,
			Anulacion:       false,
			ContraCargo:     false,
			ContraCargoMx:   false,
			ContraCargoPx:   false,
			Reversion:       false,
		}

		filtro_clpx_devolucion := filtroCl.FiltroCierreLote{
			MatchCl:         true,
			MovimientosMX:   false,
			PagosPx:         true,
			Banco:           true,
			EstadoFechaPago: true,
			FechaPago:       "0000-00-00",
			Compras:         false,
			Devolucion:      true,
			Anulacion:       false,
			ContraCargo:     false,
			ContraCargoMx:   false,
			ContraCargoPx:   false,
			Reversion:       false,
		}

		filtro_clpx_anulacion := filtroCl.FiltroCierreLote{
			MatchCl:         true,
			MovimientosMX:   false,
			PagosPx:         true,
			Banco:           true,
			EstadoFechaPago: true,
			FechaPago:       "0000-00-00",
			Compras:         false,
			Devolucion:      false,
			Anulacion:       true,
			ContraCargo:     false,
			ContraCargoMx:   false,
			ContraCargoPx:   false,
			Reversion:       false,
		}

		filtro_clpx_cc := filtroCl.FiltroCierreLote{
			MatchCl:         true,
			MovimientosMX:   false,
			PagosPx:         true,
			Banco:           true,
			EstadoFechaPago: true,
			FechaPago:       "0000-00-00",
			Compras:         false,
			Devolucion:      false,
			Anulacion:       false,
			ContraCargo:     true,
			ContraCargoMx:   false,
			ContraCargoPx:   false,
			Reversion:       false,
		}

		var filtros_clpx_automaticos []filtroCl.FiltroCierreLote

		filtros_clpx_automaticos = append(filtros_clpx_automaticos, filtro_clpx_compra, filtro_clpx_devolucion, filtro_clpx_anulacion, filtro_clpx_cc)
		filtros_clpx_nombres := []string{"filtro compras", "filtro devolucion", "filtro anulacion", "filtro contracargo"}

		for i_loop, filtro_loop := range filtros_clpx_automaticos {
			filtro_clpx_aplicado := filtro_loop

			// var filtro filtroCl.FiltroPrismaMovimiento

			// filtro = filtroCl.FiltroPrismaMovimiento{
			// 	Match:                        false,
			// 	CargarDetalle:                true,
			// 	Contracargovisa:              true,
			// 	Contracargomaster:            true,
			// 	Tipooperacion:                true,
			// 	Rechazotransaccionprincipal:  true,
			// 	Rechazotransaccionsecundario: true,
			// 	Motivoajuste:                 true,
			// 	ContraCargo:                  filtro_clpx_aplicado.ContraCargo,
			// 	CodigosOperacion:             []string{"0005"},
			// 	TipoAplicacion:               "+",
			// }
			// if filtro_clpx_aplicado.ContraCargo {
			// 	filtro.Match = false
			// 	filtro.CodigosOperacion = []string{"1507", "6000", , "1517"}
			// 	filtro.TipoAplicacion = "-"
			// }

			var codigo []string
			listaCierreLote, err := cierrelote.ObtenerCierreloteServices(filtro_clpx_aplicado, codigo)
			if err != nil {
				errObtenerEstados := errors.New(err.Error())
				err = errObtenerEstados
				logError := entities.Log{
					Tipo:          entities.EnumLog("error"),
					Funcionalidad: "ObtenerCierreloteServices",
					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
				}
				errCrearLog := service.CreateLogService(logError)
				if errCrearLog != nil {
					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
				}
			} else {

				if len(listaCierreLote) == 0 {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionConciliacionCLPx,
						Descripcion: fmt.Sprintf("fecha : %v - no existe cierre de lotes para conciliar con pagos %v", time.Now().String(), filtros_clpx_nombres[i_loop]),
					}
					service.CreateNotificacionService(notificacion)
				} else {

					filtroCabecera := filtroCl.FiltroPrismaMovimiento{
						ContraCargo: filtro_clpx_aplicado.ContraCargo,
					}
					listaCierreLoteMovimientos, err := cierrelote.ObtenerPrismaMovimientoConciliadosServices(listaCierreLote, filtroCabecera)
					if err != nil {
						errObtenerEstados := errors.New(err.Error())
						err = errObtenerEstados
						logError := entities.Log{
							Tipo:          entities.EnumLog("error"),
							Funcionalidad: "ObtenerPrismaMovimientoConciliadosServices",
							Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
						}
						errCrearLog := service.CreateLogService(logError)
						if errCrearLog != nil {
							logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
						}
					} else {

						var listaFechaPagos []string
						for _, value := range listaCierreLoteMovimientos {
							fechaString := value.MovimientoCabecer.FechaPago.Format("2006-01-02")
							listaFechaPagos = append(listaFechaPagos, fechaString)
						}
						filtro_prisma := filtroCl.FiltroPrismaTrPagos{
							Match:         false,
							CargarDetalle: true,
							Devolucion:    filtro_clpx_aplicado.Devolucion,
							FechaPagos:    listaFechaPagos,
						}
						listaPrismaPago, err := cierrelote.ObtenerPrismaPagosServices(filtro_prisma)
						if err != nil {
							errObtenerEstados := errors.New(err.Error())
							err = errObtenerEstados
							logError := entities.Log{
								Tipo:          entities.EnumLog("error"),
								Funcionalidad: "ObtenerPrismaPagosServices",
								Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
							}
							errCrearLog := service.CreateLogService(logError)
							if errCrearLog != nil {
								logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
							}
						} else {
							listaCierreLoteProcesado, listaIdsDetalle, listaIdsCabecera, err := cierrelote.ConciliarCierreLotePrismaPagoServices(listaCierreLoteMovimientos, listaPrismaPago)
							if err != nil {
								errObtenerEstados := errors.New(err.Error())
								err = errObtenerEstados
								logError := entities.Log{
									Tipo:          entities.EnumLog("error"),
									Funcionalidad: "ConciliarCierreLotePrismaPagoServices",
									Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
								}
								errCrearLog := service.CreateLogService(logError)
								if errCrearLog != nil {
									logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
								}
							} else {

								if len(listaCierreLoteProcesado) <= 0 && len(listaIdsDetalle) <= 0 && len(listaIdsCabecera) <= 0 {
									notificacion := entities.Notificacione{
										Tipo:        entities.NotificacionConciliacionCLPx,
										Descripcion: fmt.Sprintf("fecha : %v - no existe cierre de lotes para conciliar con pagos %v", time.Now().String(), filtros_clpx_nombres[i_loop]),
									}
									service.CreateNotificacionService(notificacion)
								}
								err = cierrelote.ActualizarCierrelotePagosServices(listaCierreLoteProcesado, listaIdsCabecera, listaIdsDetalle)
								if err != nil {
									errObtenerEstados := errors.New(err.Error())
									err = errObtenerEstados
									logError := entities.Log{
										Tipo:          entities.EnumLog("error"),
										Funcionalidad: "ActualizarCierrelotePagosServices",
										Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
									}
									errCrearLog := service.CreateLogService(logError)
									if errCrearLog != nil {
										logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
									}
								} else {
									notificacion := entities.Notificacione{
										Tipo:        entities.NotificacionConciliacionCLPx,
										Descripcion: fmt.Sprintf("fecha : %v - proceso de conciliacion pagos con cierre lote exito %v", time.Now().String(), filtros_clpx_nombres[i_loop]),
									}
									service.CreateNotificacionService(notificacion)
								}
							}

						}

					}
				}
			}

		}

	}
	// add job to cron
	cronjob.AddFunc(periodicidad, getConciliarCLPX)
}

func GetGenerarMovimientosPrisma(cronjob *cron.Cron, periodicidad string, reversion bool, cierrelote cierrelote.Service, service administracion.Service) {
	var getGenerarMovimientosPrisma = func() {

		responseCierreLote, err := service.BuildPrismaMovimiento(reversion)
		if err != nil {
			errObtenerEstados := errors.New(err.Error())
			err = errObtenerEstados
			logError := entities.Log{
				Tipo:          entities.EnumLog("error"),
				Funcionalidad: "BuildPrismaMovimiento",
				Mensaje:       errObtenerEstados.Error(), // + "-" + err.Error(),
			}
			errCrearLog := service.CreateLogService(logError)
			if errCrearLog != nil {
				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
			}
		} else {

			ctxPrueba := context.WithValue(context.Background(), entities.AuditUserKey{}, entities.Auditoria{UserID: 1})
			err = service.CreateMovimientosService(ctxPrueba, responseCierreLote)
			if err != nil {
				errObtenerEstados := errors.New(err.Error())
				err = errObtenerEstados
				logError := entities.Log{
					Tipo:          entities.EnumLog("error"),
					Funcionalidad: "CreateMovimientosService",
					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
				}
				errCrearLog := service.CreateLogService(logError)
				if errCrearLog != nil {
					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
				}
			} else {
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionConciliacionCLPx,
					Descripcion: fmt.Sprintf("fecha : %v - proceso de conciliacion pagos con cierre lote exito ", time.Now().String()),
				}
				service.CreateNotificacionService(notificacion)
			}
		}
	}
	cronjob.AddFunc(periodicidad, getGenerarMovimientosPrisma)
}

func GetConciliacionBancoPrisma(cronjob *cron.Cron, periodicidad string, reversion bool, cierrelote cierrelote.Service, service administracion.Service) {
	var getConciliacionBancoPrisma = func() {
		fechaActual, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		logs.Info(fmt.Sprintf("inicio proceso conciliacion banco fecha actual %v", fechaActual))
		if err != nil {
			errObtenerEstados := errors.New("error al parsear fecha " + err.Error())
			err = errObtenerEstados
			logError := entities.Log{
				Tipo:          entities.EnumLog("error"),
				Funcionalidad: "parseo de fecha ",
				Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
			}
			errCrearLog := service.CreateLogService(logError)
			if errCrearLog != nil {
				logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
			}
		} else {
			fechatemporal := fechaActual.Add(24 * -1)
			fechaPagoProcesar := fechatemporal.Format("2006-01-02")
			logs.Info(fmt.Sprintf("obteniendo tipo conciliacion reversion %v", reversion))
			logs.Info(fmt.Sprintf("obteniendo una fecha antes de la actual %v", fechaPagoProcesar))
			filtro := filtroCl.FiltroTablasConciliadas{
				FechaPago: fechaPagoProcesar,
				Match:     true,
				Reversion: reversion,
			}
			responseListprismaTrPagos, err := cierrelote.ObtenerRepoPagosPrisma(filtro)
			if err != nil {
				errObtenerEstados := errors.New(err.Error())
				err = errObtenerEstados
				logError := entities.Log{
					Tipo:          entities.EnumLog("error"),
					Funcionalidad: "ActualizarCierreloteMoviminetosServices",
					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
				}
				errCrearLog := service.CreateLogService(logError)
				if errCrearLog != nil {
					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
				}
			}
			/* obtener configuracion periodo de acreditacion */
			movimientoBanco, erro := cierrelote.ConciliacionBancoPrisma(fechaPagoProcesar, filtro.Reversion, responseListprismaTrPagos)
			if erro != nil {
				errObtenerEstados := errors.New(err.Error())
				err = errObtenerEstados
				logError := entities.Log{
					Tipo:          entities.EnumLog("error"),
					Funcionalidad: "ActualizarCierreloteMoviminetosServices",
					Mensaje:       errObtenerEstados.Error() + "-" + err.Error(),
				}
				errCrearLog := service.CreateLogService(logError)
				if errCrearLog != nil {
					logs.Error("error al intentar crear un log - " + errCrearLog.Error() + " - " + logError.Mensaje)
				}
			}
			if movimientoBanco == nil {
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionConciliacionBancoCL,
					Descripcion: fmt.Sprintf("fecha : %v - no existe movimientos en banco para conciliar con los pagos ", time.Now().String()),
				}
				service.CreateNotificacionService(notificacion)

			} else {
				notificacion := entities.Notificacione{
					Tipo:        entities.NotificacionConciliacionBancoCL,
					Descripcion: fmt.Sprintf("fecha : %v -  proceso de conciliacion pagos con banco exitoso ", time.Now().String()),
				}
				service.CreateNotificacionService(notificacion)
			}
		}

	}
	cronjob.AddFunc(periodicidad, getConciliacionBancoPrisma)
}
