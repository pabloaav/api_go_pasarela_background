package reportes

import (
	"fmt"
	"strings"
)

// estructura d ela linea de txt SICORE
// Los valores fijos son estaticos e iguales para cualquier linea
type LineData struct {
	Codigo                     string // Fijo 06
	FechaRrm                   string // Formato dd/mm/yyyy
	NumeroRrm                  string // 0000100000001
	ImporteRrm                 string // sin puntos, solo coma para separador decimal
	CodigoGravamen             string
	CodigoRegimen              string
	CodigoEsRetencion          string  // Fijo 1
	ImporteComprobanteCabecera float64 // sin puntos, solo coma para separador decimal
	FechaComprobanteCabecera   string  // Formato dd/mm/yyyy
	CodigoCondicion            string  // Fijo 01
	TotalRetenido              float64 // sin puntos, solo coma para separador decimal
	PorcentajeExclusion        string  // Fijo 000,00
	TipoDocumento              string  // Fijo 80
	CuitCliente                string  // 11 digitos
	NumeroCertificadoOriginal  string  // Fijo 00000000000000
}

// estructura d ela linea de txt SICAR
// Los valores fijos son estaticos e iguales para cualquier linea
type LineDataIIBB struct {
	NroRenglon        string  // Numerico (5)
	OrigenComprobante string  // Numerico (1) Siempre 1 ( Comprobante generado por software propio) -> 2 por sistema SIRCAR
	TipoComprobante   string  // Numerico (1) Siempre 1 ( Comprobante de Retencion) -> 2 Para anulado
	NroComprobante    string  // Numerico (12) Siempre 0 para nosotros
	CuitContribuyente string  // Numerico (11)
	FechaRetencion    string  // Formato dd/mm/yyyy
	Monto             float64 // sin puntos, solo coma para separador decimal
	Alicouta          float64 // sin puntos, solo coma para separador decimal
	MontoRetenido     float64 // sin puntos, solo coma para separador decimal
	TipoRegimen       string  // Numerico (3) siempre 02
	Jurisdiccion      string  // Numerico (3) siempre 905
}

// estructura de la linea Cabecera de txt Comisiones
// Los valores fijos son estaticos e iguales para cualquier linea. Entre parentesis la longitud del campo.
type CabeceraData struct {
	TipoRegistro      string // (2) Fijo 01
	CuitInformante    string // (11) CUIT sin guiones Fijo "30716550849"
	PeriodoInformado  string // (6) Formato YYYYMM
	Secuencia         string // (2) Fijo 00 = Original . 01 = rectificativa
	Denominacion      string // (200) Fijo "CORRIENTES TELECOMUNICACIONES SAPEM". Autocompletar con espacios a la derecha.
	Hora              string // (6) Formato HHMMSS. Fijo 000000 por conveniencia.
	CodigoImpuesto    string // (4) Fijo "0103".
	CodigoConcepto    string // (3) Fijo "830".
	NumeroVerificador string // (6) Numero Consecutivo al txt anterior
	NumeroFormulario  string // (4) Fijo "8125"
	NumeroVersion     string // (5) Fijo "00100"
	Establecimiento   string // (2) Fijo "00"
	CantidadRegistros string // (10) Cantidad de registros (lineas del txt).
}

// estructura de la linea Vendedor o prestador de servicios de txt Comisiones
// Los valores fijos son estaticos e iguales para cualquier linea. Entre parentesis la longitud del campo.
type VendedorData struct {
	TipoRegistro           string // (2) Fijo 02
	TipoIdentificacion     string // (2) Fijo 80 = "CUIT"
	IdentificacionVendedor string // (11) CUIT sin guiones
	CodigoRubro            string // (2) Por ahora fijo 07 "Servicios"
	SignoTotal             string // (1) 0 = "Positivo" o 1 = "Negativo"
	MontoTotal             string // (12) Monto total de operaciones del mes en pesos sin decimales
	ImporteComision        string // (12) Monto importe sin decimales
}

// estructura de la linea "Detalle de operaciones" de servicios de txt Comisiones
// Los valores fijos son estaticos e iguales para cualquier linea. Entre parentesis la longitud del campo.
type DetalleData struct {
	TipoRegistro            string // (2) Fijo 03
	MetodologiaAcreditacion string // (2) Fijo 01 = "CBU"
	TipoCuenta              string // (2) Por ahora Fijo 01 = "Caja de Ahorro"
	NumeroIdentificacion    string // (22) CBU donde se acredita
	SignoMonto              string // (1) 0 = Positivo o 1 = Negativo
	Monto                   string // (12) Sin decimales
}

type LineDataExcel struct {
	COD_PROV   string
	NOM_PROV   string
	IDENTIFTRI string
	INGR_PROV  string
	FEC_RET    string
	COD_RET    string
	T_COMP     string
	N_COMP     string
	N_CERTIFIC string
	IMP_PAGO   string
	IMP_RETEN  string
	COD_PROVE  string
	DESC_PROVE string
	T_RETEN    string
}

// LineBuilder es el constructor para crear líneas de archivo.
type LineBuilder struct {
	line string
}

// Constructor NewLineBuilder crea una nueva instancia de LineBuilder.
func NewLineBuilder() *LineBuilder {
	return &LineBuilder{}
}

// SetString agrega string en la línea.
func (lb *LineBuilder) SetString(dato string) *LineBuilder {
	lb.line += dato
	return lb
}

// SetSpaces establece espacios en la línea.
func (lb *LineBuilder) SetSpaces(count int) *LineBuilder {
	lb.line += strings.Repeat(" ", count)
	return lb
}

// SetComma establece comas en la línea.
func (lb *LineBuilder) SetComma(count int) *LineBuilder {
	lb.line += strings.Repeat(",", count)
	return lb
}

// SetZeros establece ceros en la línea.
func (lb *LineBuilder) SetZeros(count int) *LineBuilder {
	lb.line += strings.Repeat("0", count)
	return lb
}

// SetValueFloat agrega un valor de tipo float64 a la línea con relleno de ceros a la izquierda si es necesario.
func (lb *LineBuilder) SetValueFloat(value float64, length int) *LineBuilder {
	valueStr := fmt.Sprintf("%.2f", value)
	if length > 0 {
		// Rellena con ceros a la izquierda si es necesario.
		remainingZeros := length - len(valueStr)
		if remainingZeros > 0 {
			lb.line += strings.Repeat("0", remainingZeros) + strings.Replace(valueStr, ".", ",", -1)
		} else {
			lb.line += strings.Replace(valueStr, ".", ",", -1)
		}
	} else {
		// Sin longitud especificada, simplemente agrega el valor.
		lb.line += strings.Replace(valueStr, ".", ",", -1)
	}
	return lb
}

// SetValueString agrega un valor de cadena a la línea con relleno de ceros a la izquierda si es necesario.
func (lb *LineBuilder) SetValueString(value string, maxLength int) *LineBuilder {
	if maxLength > 0 && len(value) > maxLength {
		value = value[:maxLength]
	}
	remainingZeros := maxLength - len(value)
	if remainingZeros > 0 {
		lb.line += strings.Repeat("0", remainingZeros) + value
	} else {
		lb.line += value
	}
	return lb
}

// SetStringSpaced agrega un valor de cadena a la línea con relleno de espacios a la derecha si es necesario.
func (lb *LineBuilder) SetStringSpaced(value string, maxLength int) *LineBuilder {
	if maxLength > 0 && len(value) > maxLength {
		value = value[:maxLength]
	}
	remainingSpaces := maxLength - len(value)
	if remainingSpaces > 0 {
		lb.line += value + strings.Repeat(" ", remainingSpaces)
	} else {
		lb.line += value
	}
	return lb
}

// Build devuelve la línea construida.
func (lb *LineBuilder) Build() string {
	return lb.line
}

func CreateLine(ld LineData) (line string) {
	// crear una linea del archivo txt SICORE
	line = NewLineBuilder().
		SetValueString(ld.Codigo, 0).
		SetValueString(ld.FechaRrm, 10).
		SetSpaces(3).
		SetValueString(ld.NumeroRrm, 13).
		SetValueString(ld.ImporteRrm, 16).
		SetSpaces(1).
		SetValueString(ld.CodigoGravamen, 0).
		SetSpaces(1).
		SetValueString(ld.CodigoRegimen, 0).
		SetValueString(ld.CodigoEsRetencion, 0).
		SetValueFloat(ld.ImporteComprobanteCabecera, 14).
		SetValueString(ld.FechaComprobanteCabecera, 10).
		SetValueString(ld.CodigoCondicion, 2).
		SetSpaces(1).
		SetValueFloat(ld.TotalRetenido, 14).
		SetValueString(ld.PorcentajeExclusion, 6).
		SetSpaces(10).
		SetValueString(ld.TipoDocumento, 2).
		SetValueString(ld.CuitCliente, 11).
		SetSpaces(9).
		SetValueString(ld.NumeroCertificadoOriginal, 14).
		Build()

	return
}
