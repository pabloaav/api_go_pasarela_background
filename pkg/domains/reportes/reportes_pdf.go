package reportes

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/reportedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type reportePDF struct {
	M pdf.Maroto
	s reportesService
}

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func (r *reportePDF) buildTitle(cliente administraciondtos.ResponseFacturacion, nroReporte string) (err error) {

	fecha_impresion_pdf := fmt.Sprintf(time.Now().Format("02-01-2006"))
	nroReporte_pdf := lpad(nroReporte, "0", 6)
	byteSlices, err := ioutil.ReadFile("assets/images/wee_reduce.png")
	if err != nil {
		fmt.Println("Got error while opening file:", err)
		return
	}
	base64image := base64.StdEncoding.EncodeToString(byteSlices)

	r.M.RegisterHeader(func() {
		r.M.Row(20, func() {

			r.M.Col(8, func() {

			})

			r.M.Col(4, func() {
				_ = r.M.Base64Image(base64image, consts.Png, props.Rect{
					Center:  false,
					Left:    40,
					Percent: 80,
				})
			})
		})
		r.M.Row(35, func() {

			r.M.Col(6, func() {
				r.M.Text(cliente.Cliente, props.Text{Size: 7, Style: consts.Bold})
				r.M.Text(("CUIT:" + cliente.Cuit), props.Text{Size: 7, Top: 5})
				r.M.Text(cliente.RazonSocial, props.Text{Size: 7, Top: 10})

			})

			r.M.Col(6, func() {
				r.M.Text("Corrientes Telecomunicaciones SAPEM", props.Text{Size: 7, Align: consts.Right})
				r.M.Text("Dr. Carrillo 444 5to Piso - 3400 - Corrientes", props.Text{Size: 7, Top: 5, Align: consts.Right})
				r.M.Text(("I.V.A: Responsable Inscripto"), props.Text{Size: 7, Top: 10, Align: consts.Right})
				r.M.Text("Fecha: "+fecha_impresion_pdf+" - Número: "+nroReporte_pdf, props.Text{Size: 7, Top: 15, Align: consts.Right})
				r.M.Text("CUIT: 30716550849", props.Text{Size: 7, Top: 20, Align: consts.Right})
				r.M.Text("II BB: 30716550849", props.Text{Size: 7, Top: 25, Align: consts.Right})
				r.M.Text("Inicio de Actividades: 09/2019", props.Text{Size: 7, Top: 30, Align: consts.Right})

			})
		})
	})
	return
}

func (r *reportePDF) buildBodyPagos(pagos reportedtos.ResponseClientesReportes) {

	header := []string{"Cuenta", "Referencia", "Fecha Cobro", "Medio de Pago", "Tipo", "Estado", "Monto"}

	contents := [][]string{}

	for _, pago := range pagos.Pagos {
		contents = append(contents, []string{pago.Cuenta, pago.Id, pago.FechaPago, pago.MedioPago, pago.Tipo, pago.Estado, pago.Monto})
	}

	r.M.TableList(header, contents, props.TableList{
		Align: consts.Center,
		ContentProp: props.TableListContent{
			GridSizes: []uint{3, 1, 1, 2, 1, 2, 1},
			Size:      8,
			// Color:     color.Color{100, 0, 0},
		},
		HeaderProp: props.TableListContent{
			GridSizes: []uint{3, 1, 1, 2, 1, 2, 1},
			Size:      8,
		},
		// AlternatedBackground: &color.Color{
		// 	Red:   200,
		// 	Green: 200,
		// 	Blue:  200,
		// },
		VerticalContentPadding: 2,
	})

	r.M.Row(20, func() {

		r.M.Col(8, func() {
			r.M.Text("Cantidad Operaciones: ", props.Text{Size: 10, Align: consts.Right})
			r.M.Text("Total Cobrado: ", props.Text{Size: 10, Top: 10, Align: consts.Right})
		})

		r.M.Col(4, func() {
			r.M.Text(pagos.CantOperaciones, props.Text{Size: 10, Align: consts.Right})
			r.M.Text("$"+pagos.TotalCobrado, props.Text{Size: 10, Top: 10, Align: consts.Right})

		})
	})

}

func (r *reportePDF) buildBodyCobranzas(request []reportedtos.ResponseClientesReportes, totales reportedtos.TotalCobranzasDiarias) {

	for _, pago := range request {
		contents := [][]string{}
		for _, pg := range pago.Pagos {
			contents = append(contents, []string{pg.Cuenta, pg.Id, pg.FechaPago, pg.MedioPago, pg.Tipo, pg.Estado, pg.Monto})
		}
		header := []string{pago.Clientes, "Referencia pago", "Fecha Cobro", "Medio de Pago", "Tipo", "Estado", "Monto"}
		r.M.TableList(header, contents, props.TableList{
			Align: consts.Center,
			ContentProp: props.TableListContent{
				GridSizes: []uint{3, 2, 1, 2, 1, 2, 1},
				Size:      8,
				// Color:     color.Color{100, 0, 0},
			},
			HeaderProp: props.TableListContent{
				GridSizes: []uint{3, 2, 1, 2, 1, 2, 1},
				Size:      8,
			},
			LineProp: props.Line{Color: color.NewBlack(), Width: 2},
			// AlternatedBackground: &color.Color{
			// 	Red:   200,
			// 	Green: 200,
			// 	Blue:  200,
			// },
			VerticalContentPadding: 2,
		})

		r.M.Line(1.2)

		r.M.Row(20, func() {
			r.M.Col(12, func() {
				r.M.Text("Subtotal $"+pago.TotalCobrado, props.Text{Style: consts.Bold, Align: consts.Right})
				r.M.Text("Comision $"+pago.TotalComision, props.Text{Style: consts.Bold, Align: consts.Right, Top: 8})
			})
		})

	}

	// Final totales
	r.M.Row(30, func() {

		r.M.Col(12, func() {
			r.M.Text("Total cant. operaciones: "+totales.CantidadOperaciones, props.Text{Size: 10, Align: consts.Right})
			r.M.Text("Total comisiones $"+totales.ComisionTotal, props.Text{Size: 10, Top: 10, Align: consts.Right})
			r.M.Text("Total $"+totales.Total, props.Text{Size: 10, Top: 20, Align: consts.Right})

		})
	})

}

func (r *reportePDF) buildBodyRendiciones(pagos reportedtos.ResponseClientesReportes) {

	header := []string{"Cuenta", "Referencia", "Concepto", "Fecha Cobro", "Fecha Deposito", "Importe Cobrado", "Importe Deposito ", "Cant. Boletas Cobradas", "Comision", "IVA"}

	contents := [][]string{}

	for _, rendicion := range pagos.Rendiciones {
		contents = append(contents, []string{rendicion.Cuenta, rendicion.Id, rendicion.Concepto, rendicion.FechaCobro, rendicion.FechaDeposito, rendicion.ImporteCobrado, rendicion.ImporteDepositado, rendicion.CantidadBoletasCobradas, rendicion.Comision, rendicion.Iva})
	}

	r.M.TableList(header, contents, props.TableList{
		Align: consts.Center,
		ContentProp: props.TableListContent{
			GridSizes: []uint{2, 1, 2, 1, 1, 1, 1, 1, 1, 1},
			Size:      7,
			// Color:     color.Color{100, 0, 0},
		},
		HeaderProp: props.TableListContent{
			GridSizes: []uint{2, 1, 2, 1, 1, 1, 1, 1, 1, 1},
			Size:      7,
			Style:     consts.Bold,
		},
		// AlternatedBackground: &color.Color{
		// 	Red:   200,
		// 	Green: 200,
		// 	Blue:  200,
		// },
		VerticalContentPadding: 2,
		HeaderContentSpace:     0.5,
	})
	r.M.Row(5, func() {})

	// Totales acumulativos del reporte
	TotalesReporteRendicion(r, pagos)
}

func TotalesReporteRendicion(r *reportePDF, pagos reportedtos.ResponseClientesReportes) {

	gananciasMayorQueCero := CompararNumeroString(pagos.TotalRetencionGanancias)
	IVAMayorQueCero := CompararNumeroString(pagos.TotalRetencionIva)
	IIBBMayorQueCero := CompararNumeroString(pagos.TotalRetencionIIBB)

	// Totales Pie de Pagina
	r.M.Row(60, func() {

		r.M.Col(8, func() {
			r.M.Text("Total Cobrado: ", props.Text{Size: 10, Align: consts.Right})
			r.M.Text("Total Comision: ", props.Text{Size: 10, Top: 10, Align: consts.Right})
			r.M.Text("Total IVA 21%: ", props.Text{Size: 10, Top: 20, Align: consts.Right})
			r.M.Text("Total Rendido: ", props.Text{Size: 10, Top: 30, Align: consts.Right})
			if gananciasMayorQueCero {
				r.M.Text("Ret. Gcia. RG 4622/2019: ", props.Text{Size: 10, Top: 40, Align: consts.Right})
			}
			if IVAMayorQueCero {
				r.M.Text("Ret. IVA RG 4622/2019: ", props.Text{Size: 10, Top: 50, Align: consts.Right})
			}
			if IIBBMayorQueCero {
				r.M.Text("Ret. IIBB RG 202/2020: ", props.Text{Size: 10, Top: 60, Align: consts.Right})
			}

		})

		r.M.Col(4, func() {
			r.M.Text("$"+pagos.TotalCobrado, props.Text{Size: 10, Align: consts.Right})
			r.M.Text("$"+pagos.TotalComision, props.Text{Size: 10, Top: 10, Align: consts.Right})
			r.M.Text("$"+pagos.TotalIva, props.Text{Size: 10, Top: 20, Align: consts.Right})
			r.M.Text("$"+pagos.RendicionTotal, props.Text{Size: 10, Top: 30, Align: consts.Right})
			if gananciasMayorQueCero {
				r.M.Text("$"+pagos.TotalRetencionGanancias, props.Text{Size: 10, Top: 40, Align: consts.Right})
			}
			if IVAMayorQueCero {
				r.M.Text("$"+pagos.TotalRetencionIva, props.Text{Size: 10, Top: 50, Align: consts.Right})
			}
			if IIBBMayorQueCero {
				r.M.Text("$"+pagos.TotalRetencionIIBB, props.Text{Size: 10, Top: 60, Align: consts.Right})
			}
		})

	})
}

// Parsea el numero string a float64.
// Retorna true en caso que el numero sea mayor que cero
// Uso: imprimir valores totales de retenciones mayores que cero.
func CompararNumeroString(strNumero string) (mayorQueCero bool) {

	// Elimina los puntos como separadores de miles
	strNumero = strings.ReplaceAll(strNumero, ".", "")

	// Reemplaza la coma por un punto como separador decimal
	strNumero = strings.Replace(strNumero, ",", ".", 1)

	// Convierte el string a un float64
	numero, err := strconv.ParseFloat(strNumero, 64)
	if err != nil {
		fmt.Println("Error al convertir el número en reporte de rendicion pdf:", err)
	}

	if numero > 0 {
		mayorQueCero = true
	}
	return
}

func GetPagosPdf(pagos reportedtos.ResponseClientesReportes, cliente administraciondtos.ResponseFacturacion, fecha string, nroReporte string) error {

	reversionPdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	reversionPdf.SetPageMargins(10, 10, 10)

	var rep reportePDF
	rep.M = reversionPdf

	err := rep.buildTitle(cliente, nroReporte)
	if err != nil {
		return err
	}

	rep.M.Row(5, func() {})

	rep.M.Row(10, func() {
		rep.M.Col(6, func() {
			rep.M.Text("REPORTE DE COBRANZAS", props.Text{Style: consts.Bold})
		})
		rep.M.Col(6, func() {
			rep.M.Text(("Fecha: " + fecha), props.Text{Align: consts.Right})
		})
	})
	rep.M.Row(5, func() {

	})

	// Las cabeceras de cada pago revertidos y los items de cada pago
	rep.buildBodyPagos(pagos)

	// Se crea la carpeta en caso de que no exista
	tempFolder := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE)
	if _, err := os.Stat(tempFolder); os.IsNotExist(err) {
		err = os.MkdirAll(tempFolder, 0755)
		if err != nil {
			return err
		}
	}

	err = rep.M.OutputFileAndClose(tempFolder + "/" + (cliente.Cliente + "-" + fecha) + ".pdf")
	if err != nil {
		return err
	}

	return nil
}

func GetRendicionesPdf(pagos reportedtos.ResponseClientesReportes, cliente administraciondtos.ResponseFacturacion, fecha string, nroReporte string) error {

	rendicionPdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	rendicionPdf.SetPageMargins(10, 10, 10)

	var rep reportePDF
	rep.M = rendicionPdf

	err := rep.buildTitle(cliente, nroReporte)
	if err != nil {
		return err
	}
	rep.M.Row(5, func() {})

	rep.M.Row(10, func() {
		rep.M.Col(6, func() {
			rep.M.Text("REPORTE DE RENDICIONES", props.Text{Style: consts.Bold})
		})
		rep.M.Col(6, func() {
			rep.M.Text(("Fecha: " + fecha), props.Text{Align: consts.Right})
		})
	})
	rep.M.Row(5, func() {})

	// Las cabeceras de cada pago revertidos y los items de cada pago
	rep.buildBodyRendiciones(pagos)

	// Se crea la carpeta en caso de que no exista
	tempFolder := fmt.Sprintf(config.DIR_BASE + config.DIR_REPORTE)
	if _, err := os.Stat(tempFolder); os.IsNotExist(err) {
		err = os.MkdirAll(tempFolder, 0755)
		if err != nil {
			return err
		}
	}

	err = rep.M.OutputFileAndClose(tempFolder + "/" + (cliente.Cliente + "-" + fecha) + ".pdf")
	if err != nil {
		return err
	}

	return nil
}

func GetCobranzasPdf(pagos []reportedtos.ResponseClientesReportes, cliente administraciondtos.ResponseFacturacion, fecha string, totales reportedtos.TotalCobranzasDiarias) error {

	reversionPdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	reversionPdf.SetPageMargins(10, 10, 10)

	var rep reportePDF
	rep.M = reversionPdf

	err := rep.buildTitleCobranzas(cliente)
	if err != nil {
		return err
	}

	rep.M.Row(10, func() {
		rep.M.Col(6, func() {
			rep.M.Text("Cobranzas del "+fecha, props.Text{Style: consts.Bold})
		})
		// rep.M.Col(2, func() {
		// 	rep.M.Text(("Fecha: " + fecha), props.Text{Style: consts.Bold})
		// })
	})

	// Las cabeceras de cada pago revertidos y los items de cada pago
	rep.buildBodyCobranzas(pagos, totales)

	// Se crea la carpeta en caso de que no exista
	tempFolder := fmt.Sprintf(config.DIR_BASE + config.DOC_CL + "/reportes")
	if _, err := os.Stat(tempFolder); os.IsNotExist(err) {
		err = os.MkdirAll(tempFolder, 0755)
		if err != nil {
			return err
		}
	}

	err = rep.M.OutputFileAndClose(tempFolder + "/" + (cliente.Cliente + "-" + fecha) + ".pdf")
	if err != nil {
		return err
	}

	return nil
}

func (r *reportePDF) buildTitleCobranzas(cliente administraciondtos.ResponseFacturacion) (err error) {

	// fecha_impresion_pdf := fmt.Sprintf(time.Now().Format("02-01-2006"))
	byteSlices, err := ioutil.ReadFile("assets/images/wee_reduce.png")
	if err != nil {
		fmt.Println("Got error while opening file:", err)
		return
	}
	base64image := base64.StdEncoding.EncodeToString(byteSlices)

	r.M.RegisterHeader(func() {
		r.M.Row(30, func() {

			r.M.Col(8, func() {
				r.M.Text("Corrientes Telecomunicaciones SAPEM", props.Text{Size: 8})
				r.M.Text(("Dr. R. Carrillo 444 5to Piso - 3400 - Corrientes"), props.Text{Size: 8, Top: 4})
				r.M.Text("I.V.A. Responsable Inscripto", props.Text{Size: 8, Top: 8})
				r.M.Text("CUIT: 30716550849", props.Text{Size: 8, Top: 12})
				r.M.Text("Inicio de Actividades: 09/2019", props.Text{Size: 8, Top: 16})
			})

			r.M.Col(4, func() {
				_ = r.M.Base64Image(base64image, consts.Png, props.Rect{
					Center:  true,
					Percent: 80,
				})
			})
		})
	})
	return
}

func RRComprobantePDF(comprobante entities.Comprobante, codigoImpuesto, nro_reporte_rrm, fechaFinPeriodo string) (ruta string, erro error) {
	var concepto_pago = "Usuarios de sistemas de pago electrónico"
	RRComprobantePdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	RRComprobantePdf.SetPageMargins(10, 10, 10)

	var rep reportePDF
	rep.M = RRComprobantePdf

	rep.M.RegisterHeader(func() {
		rep.M.Row(20, func() {

			rep.M.Col(8, func() {
				rep.M.Text("Corrientes Telecomunicaciones SAPEM", props.Text{Size: 10})
				rep.M.Text(("Dr. R. Carrillo 444 5to Piso - 3400 - Corrientes"), props.Text{Size: 10, Top: 7})
				rep.M.Text("CUIT: 30-71655084-9", props.Text{Size: 10, Top: 14})
				rep.M.Text("Comprobante de Retención: "+strings.ToUpper(comprobante.Gravamen), props.Text{Size: 10, Top: 21, Style: consts.Bold})
			})
		})
	}) // fin rep.M.RegisterHeader

	rep.M.Row(7, func() {})

	rep.M.Row(7, func() {
		rep.M.Col(12, func() {
			rep.M.Text(("Numero: " + comprobante.Numero), props.Text{Align: consts.Right, Size: 10})
		})
	})
	rep.M.Row(7, func() {
		rep.M.Col(12, func() {
			rep.M.Text(("Número de Rendición: " + nro_reporte_rrm), props.Text{Align: consts.Right, Size: 10})
		})
	})
	rep.M.Row(7, func() {
		rep.M.Col(12, func() {
			rep.M.Text(("Fecha: " + fechaFinPeriodo), props.Text{Align: consts.Right, Size: 10})
		})
	})
	rep.M.Row(7, func() {})

	rep.M.Row(7, func() {
		rep.M.Col(12, func() {
			rep.M.Text(("Razon Social: " + comprobante.Cliente.Razonsocial), props.Text{Size: 10})
		})
	})
	rep.M.Row(7, func() {
		rep.M.Col(12, func() {
			rep.M.Text(("Domicilio: " + comprobante.Cliente.Domicilio), props.Text{Size: 10})
		})
	})
	rep.M.Row(7, func() {
		rep.M.Col(12, func() {
			rep.M.Text(("Nro. de C.U.I.T. : " + comprobante.Cliente.Cuit), props.Text{Size: 10})
		})
	})
	rep.M.Row(7, func() {})

	// para cada comprobante detalle
	for _, detalle := range comprobante.ComprobanteDetalles {
		rep.M.Row(7, func() {
			rep.M.Col(12, func() {
				rep.M.Text(("Concepto Pago: " + concepto_pago), props.Text{Size: 10})
			})
		})
		rep.M.Row(7, func() {
			rep.M.Col(12, func() {
				rep.M.Text(("Código de Impuesto: " + codigoImpuesto), props.Text{Size: 10})
			})
		})

		rep.M.Row(7, func() {
			rep.M.Col(12, func() {
				rep.M.Text(("Código de Regimen: " + detalle.CodigoRegimen), props.Text{Size: 10})
			})
		})

		// importe := float64(comprobante.Importe) / 100

		rep.M.Row(7, func() {
			rep.M.Col(12, func() {
				rep.M.Text(("Importe Pagado Sujeto a Retencion: " + strconv.FormatFloat(detalle.TotalMonto.Float64(), 'f', 2, 64)), props.Text{Size: 10})
			})
		})
		rep.M.Row(7, func() {
			rep.M.Col(12, func() {
				rep.M.Text(("Importe Retenido: " + strconv.FormatFloat(detalle.TotalRetencion.Float64(), 'f', 2, 64)), props.Text{Size: 10})
			})
		})

		rep.M.Row(7, func() {})
	} // Fin del for _, detalle := range comprobante.ComprobanteDetalles

	rep.M.Row(7, func() {
		rep.M.Col(12, func() {
			rep.M.Text(("Por Corrientes Telecomunicaciones SAPEM"), props.Text{Size: 10})
		})
	})

	rep.M.Row(7, func() {})

	rep.M.Line(3.0, props.Line{Style: consts.Dashed, Color: color.NewBlack()})

	rep.M.Row(7, func() {})

	rep.M.Row(7, func() {
		rep.M.Col(12, func() {
			rep.M.Text(("Declaro que los datos consignados en este Formulario son correctos y completos sin omitir ni falsear dato alguno que deba contener, siendo fiel expresión de la verdad."), props.Text{Size: 10})
		})
	})

	// Se crea la carpeta en caso de que no exista
	tempFolder := fmt.Sprintf((config.DIR_BASE + config.DIR_COMP_RETENCIONES))
	if _, err := os.Stat(tempFolder); os.IsNotExist(err) {
		err = os.MkdirAll(tempFolder, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	ruta = (comprobante.Cliente.Cliente + "-" + ("RRC" + comprobante.Numero) + fechaFinPeriodo)

	err := rep.M.OutputFileAndClose(tempFolder + "/" + ruta + ".pdf")
	if err != nil {
		erro = err
		return
	}

	return
}

func (s reportesService) ReportesRendicionesPDF(reportesRendiciones []reportedtos.DetallesRendicion, cliente entities.Cliente, cuenta entities.Cuenta, sigNumero uint, FechaInicio string, FechaFin string) (ruta string, totales reportedtos.ReporteMensual, erro error) {

	var totalComision, totalIva int64

	for _, item := range reportesRendiciones {
		totalComision += int64(s.util.ToFixed(item.TotalComision, 2) * 100) 
		totalIva += int64(s.util.ToFixed(item.TotalIva, 2) * 100) 
	}

	// for _, r := range rendiciones  {
	// 	totalComision += r.TotalComision
	// 	totalIva += r.TotalIva
	// }

	RRPdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	RRPdf.SetPageMargins(10, 10, 10)

	var rep reportePDF
	rep.M = RRPdf
	rep.s = s

	// fechaNow := time.Now().Format("02-01-2006")
	fechaNow := s.commons.ConvertirFechaToDDMMYYYY(FechaFin[:10])

	byteSlices, err := ioutil.ReadFile("assets/images/wee_reduce.png")

	if err != nil {
		erro = err
		s.util.BuildLog(erro, "ReportesRendicionesPDF")
		return
	}

	base64image := base64.StdEncoding.EncodeToString(byteSlices)

	rep.M.RegisterHeader(func() {
		rep.M.Row(20, func() {

			rep.M.Col(8, func() {

			})

			rep.M.Col(4, func() {
				_ = rep.M.Base64Image(base64image, consts.Png, props.Rect{
					Center:  false,
					Left:    40,
					Percent: 80,
				})
			})
		})
		rep.M.Row(35, func() {
			// Datos del cliente y comprobante
			rep.M.Col(6, func() {
				rep.M.Text(cliente.Cliente + " - Cuenta: " + cuenta.Cuenta, props.Text{Size: 7, Style: consts.Bold})
				rep.M.Text(("CUIT:" + cliente.Cuit), props.Text{Size: 7, Top: 5})
				rep.M.Text(cliente.Razonsocial, props.Text{Size: 7, Top: 10})
				rep.M.Text("Fecha: "+fechaNow, props.Text{Size: 7, Top: 15})
				rep.M.Text("Número de comprobante: "+fmt.Sprint(s.util.GenerarNumeroComprobante2(sigNumero)), props.Text{Size: 7, Top: 20})
			})
			// Datos del agente de retencion
			rep.M.Col(6, func() {
				rep.M.Text("Corrientes Telecomunicaciones SAPEM", props.Text{Size: 7, Align: consts.Right})
				rep.M.Text("Dr. Carrillo 444 5to Piso - 3400 - Corrientes", props.Text{Size: 7, Top: 5, Align: consts.Right})
				rep.M.Text(("I.V.A: Responsable Inscripto"), props.Text{Size: 7, Top: 10, Align: consts.Right})
				rep.M.Text("CUIT: 30716550849", props.Text{Size: 7, Top: 15, Align: consts.Right})
				rep.M.Text("II BB: 30716550849", props.Text{Size: 7, Top: 20, Align: consts.Right})
				rep.M.Text("Inicio de Actividades: 09/2019", props.Text{Size: 7, Top: 25, Align: consts.Right})

			})
		})
	})

	rep.M.Row(7, func() {})

	totales, erro = rep.buildBodyReportes(reportesRendiciones, totalComision, totalIva)
	if erro != nil {
		s.util.BuildLog(erro, "ReportesRendicionesPDF")
		return
	}

	rep.M.Row(7, func() {})

	// Se crea la carpeta en caso de que no exista
	tempFolder := fmt.Sprintf((config.DIR_BASE + config.DIR_REPORTE))
	if _, err := os.Stat(tempFolder); os.IsNotExist(err) {
		err = os.MkdirAll(tempFolder, 0755)
		if err != nil {
			erro = err
			return
		}
	}

	ruta = (cliente.Cliente + "-" + "RRM" + "-" + "CTA" + fmt.Sprint(cuenta.ID) + "-" + fechaNow)

	err = rep.M.OutputFileAndClose(tempFolder + "/" + ruta + ".pdf")
	if err != nil {
		erro = err
		return
	}

	return
}

func (r *reportePDF) buildBodyReportes(reportesRendiciones []reportedtos.DetallesRendicion, totalComision, totalIva int64) (totales reportedtos.ReporteMensual, erro error) {

	contents := [][]string{}
	var (
		totalRetenido, totalCobranzas, totalRendido, totalRetGanancias, totalRetIva, totalRetIIBB float64
	)

	// por cada reporte de rendiciones se totalizan las filas y acumulan importes
	for _, reporte := range reportesRendiciones {
		var (
			floatGanancia, floatIva, floatIIBB float64
		)
		// Ret Ganancias
		if reporte.TotalRetGanancias != "" {
			floatGanancia, _ = strconv.ParseFloat(intercambiarCommasDots(reporte.TotalRetGanancias), 64)
		} else {
			reporte.TotalRetGanancias = "0"
		}
		// Ret IVA
		if reporte.TotalRetIVA != "" {
			floatIva, _ = strconv.ParseFloat(intercambiarCommasDots(reporte.TotalRetIVA), 64)
		} else {
			reporte.TotalRetIVA = "0"
		}
		// Ret IIBB
		if reporte.TotalRetIIBB != "" {
			floatIIBB, _ = strconv.ParseFloat(intercambiarCommasDots(reporte.TotalRetIIBB), 64)
		} else {
			reporte.TotalRetIIBB = "0"
		}
		// total retenido por fila
		retenido := floatGanancia + floatIva + floatIIBB

		floatTCobrado:= reporte.TotalCobrado
		
		floatTRendido :=reporte.TotalRendido
		
		/* Acumuladores de final de reporte */
		totalRetGanancias += floatGanancia
		totalRetIva += floatIva
		totalRetIIBB += floatIIBB
		totalRetenido += retenido
		totalRendido += floatTRendido
		totalCobranzas += floatTCobrado

		stringTotalRetenidoRendicion := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(retenido, 2)))
		contents = append(contents, []string{reporte.Fecha, fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(reporte.TotalCobrado, 2))), fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(reporte.TotalRendido, 2))), reporte.TotalRetGanancias, reporte.TotalRetIVA, reporte.TotalRetIIBB, stringTotalRetenidoRendicion})
	} // fin for _, reporte := range reportes

	header := []string{"Fecha", "Importe Cobrado", "Importe Rendido", "Total Ret. Gcia.", "Total Ret. IVA", "Total Ret. IIBB", "Total Retenido"}

	// crear la tabla pasandole los valores calculados
	r.M.TableList(header, contents, props.TableList{
		Align: consts.Center,
		ContentProp: props.TableListContent{
			GridSizes: []uint{1, 2, 2, 2, 1, 2, 2},
			Size:      8,
		},
		HeaderProp: props.TableListContent{
			GridSizes: []uint{1, 2, 2, 2, 1, 2, 2},
			Size:      8,
		},
		LineProp:               props.Line{Color: color.NewBlack(), Width: 2},
		VerticalContentPadding: 2,
	})

	// if len(reportes) == 0 {

	// 	contentsAux := [][]string{}
	// 	contentsAux = append(contentsAux, []string{"", "0", "0", "0", "0", "0", "0"})
	// 	headerAux := []string{"Fecha", "Importe Cobrado", "Importe Rendido", "Total Ret. Gcia.", "Total Ret. IVA", "Total Ret. IIBB", "Total Retenido"}

	// 	r.M.TableList(headerAux, contentsAux, props.TableList{
	// 		Align: consts.Center,
	// 		ContentProp: props.TableListContent{
	// 			GridSizes: []uint{1, 2, 2, 2, 1, 2, 2},
	// 			Size:      8,
	// 		},
	// 		HeaderProp: props.TableListContent{
	// 			GridSizes: []uint{1, 2, 2, 2, 1, 2, 2},
	// 			Size:      8,
	// 		},
	// 		LineProp:               props.Line{Color: color.NewBlack(), Width: 2},
	// 		VerticalContentPadding: 2,
	// 	})
	// }

	r.M.Line(1.2)

	r.M.Row(7, func() {})

	// Final totales
	stringTotalCobrado := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(totalCobranzas, 2)))
	stringTotalRendido := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(totalRendido, 2)))
	stringTotalRetenido := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(totalRetenido, 2)))
	stringTotalRetGanancias := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(totalRetGanancias, 2)))
	stringTotalRetIva := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(totalRetIva, 2)))
	stringTotalRetIIBB := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(totalRetIIBB, 2)))

	// final Comision e Iva
	stringTotalComision := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(float64(totalComision) / 100, 2))) 
	stringTotalIva := fmt.Sprintf("%v", r.s.util.FormatNum(r.s.util.ToFixed(float64(totalIva) / 100, 2))) 

	totales.TotalCobranza = stringTotalCobrado
	totales.TotalRendicion = stringTotalRendido
	totales.TotalRetenido = stringTotalRetenido
	totales.TotalRetencionGanancias = stringTotalRetGanancias
	totales.TotalRetencionIVA = stringTotalRetIva
	totales.TotalRetencionIngresosBrutos = stringTotalRetIIBB
	totales.TotalOperaciones = int64(len(reportesRendiciones))

	// Totales Pie de Pagina
	r.M.Row(70, func() {
		// Titulos de totales
		r.M.Col(10, func() {
			r.M.Text("Total Cobrado: ", props.Text{Size: 10, Align: consts.Right})
			r.M.Text("Total Rendido: ", props.Text{Size: 10, Top: 10, Align: consts.Right})
			r.M.Text("Total Retenido: ", props.Text{Size: 10, Top: 20, Align: consts.Right})
			r.M.Text("Total Comision: ", props.Text{Size: 10, Top: 30, Align: consts.Right})
			r.M.Text("Total Iva: ", props.Text{Size: 10, Top: 40, Align: consts.Right})
			if totalRetGanancias > 0 {
				r.M.Text("Total Ret. Gcia.: ", props.Text{Size: 10, Top: 50, Align: consts.Right})
			}
			if totalRetIva > 0 {
				r.M.Text("Total Ret. IVA: ", props.Text{Size: 10, Top: 60, Align: consts.Right})
			}
			if totalRetIIBB > 0 {
				r.M.Text("Total Ret. IIBB: ", props.Text{Size: 10, Top: 70, Align: consts.Right})
			}
		})
		// valores de totales
		r.M.Col(2, func() {
			r.M.Text("$"+stringTotalCobrado, props.Text{Size: 10, Align: consts.Right})
			r.M.Text("$"+stringTotalRendido, props.Text{Size: 10, Top: 10, Align: consts.Right})
			r.M.Text("$"+stringTotalRetenido, props.Text{Size: 10, Top: 20, Align: consts.Right})
			r.M.Text("$"+stringTotalComision, props.Text{Size: 10, Top: 30, Align: consts.Right})
			r.M.Text("$"+stringTotalIva, props.Text{Size: 10, Top: 40, Align: consts.Right})
			if totalRetGanancias > 0 {
				r.M.Text("$"+stringTotalRetGanancias, props.Text{Size: 10, Top: 50, Align: consts.Right})
			}
			if totalRetIva > 0 {
				r.M.Text("$"+stringTotalRetIva, props.Text{Size: 10, Top: 60, Align: consts.Right})
			}
			if totalRetIIBB > 0 {
				r.M.Text("$"+stringTotalRetIIBB, props.Text{Size: 10, Top: 70, Align: consts.Right})
			}
		})
	})

	return
}

func StringFloatToInteger(valorString string) (valorEntero int){

	// Eliminar la coma
	valorStringSinComa := strings.ReplaceAll(valorString, ",", "")

		// Eliminar los puntos decimales
		valorCadenaFinal := strings.ReplaceAll(valorStringSinComa, ".", "")

	// Convertir a entero
	valorEntero, err := strconv.Atoi(valorCadenaFinal)
	if err != nil {
		return
	}

	return
}