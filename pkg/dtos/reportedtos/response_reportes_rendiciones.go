package reportedtos

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ResponseReportesRendiciones struct {
	PagoIntentoId           uint64
	Cuenta                  string // Nombre de la cuenta del cliente
	Id                      string // external_reference enviada por el cliente
	FechaCobro              string // fecha que el pagador realizo el pago
	FechaDeposito           string // fecha que se le envio el dinero al cliente(transferencia)
	ImporteCobrado          string // importe solicitud de pago
	ImporteDepositado       string // importe depositado al cliente
	CantidadBoletasCobradas string // pago items
	// ComisionPorcentaje      string // comision de telco cobrada al cliente
	// ComisionIva             string // iva Cobrado al cliente
	Comision           string // comision de telco cobrada al cliente
	Iva                string // iva Cobrado al cliente
	Concepto           string
	CBUOrigen          string
	CBUDestino         string
	ReferenciaBancaria string
	Retenciones        string // retenciones
}

type Totales struct {
	CantidadOperaciones    string
	TotalCobrado           string
	TotalRendido           string
	TotalIva               string
	TotalComision          string
	TotalRetGanancias      string
	TotalRetIva            string
	TotalRetIngresosBrutos string
	TotalRetencion         string
}

type ResponseTotales struct {
	Totales  Totales
	Detalles []*ResponseReportesRendiciones
}

// Seccion Para Reportes visuales y excel agrupados por fecha

type ResponseRendicionesClientes struct {
	CantidadRegistros   int
	Total               float64
	DetallesRendiciones []DetallesRendicion
}

type DetallesRendicion struct {
	Fecha                                        string
	Nombre                                       string
	CantidadOperaciones                          uint
	TotalCobrado                                 float64
	TotalRendido                                 float64
	TotalReversion                               float64
	TotalComision                                float64
	TotalIva                                     float64
	Rendiciones                                  []ResponseReportesRendiciones
	NroReporte                                   string
	TotalRetGanancias, TotalRetIVA, TotalRetIIBB string
}

func (dr DetallesRendicion) ToReporteRendicion() (reporte entities.Reporte) {
		numero, _ := strconv.ParseUint(dr.NroReporte, 10, 0)
	
		reporte.Cliente = dr.Nombre                  
		reporte.Tiporeporte= "rendiciones"            
		reporte.Totalcobrado= fmt.Sprintf("%v", FormatNum(ToFixed(dr.TotalCobrado, 2)))              
		reporte.Totalrendido= fmt.Sprintf("%v", FormatNum(ToFixed(dr.TotalRendido, 2)))                     
		reporte.Fecharendicion= dr.Fecha          
		reporte.Nro_reporte = uint(numero)             
		reporte.TotalRetencionGanancias =  dr.TotalRetGanancias
		reporte.TotalRetencionIva= dr.TotalRetIVA        
		reporte.TotalRetencionIibb= dr.TotalRetIIBB         
	        
	return
}

func FormatNum(num float64) string {
	p := message.NewPrinter(language.Spanish)
	valor := p.Sprintf("%.2f", num)
	return valor
}

func  ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}