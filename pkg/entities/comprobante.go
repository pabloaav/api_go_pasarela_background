package entities

import (
	"time"

	"gorm.io/gorm"
)

type Comprobante struct {
	gorm.Model
	ClienteId           uint                 `json:"cliente_id"`
	Cliente             Cliente              `json:"cliente" gorm:"foreignKey:ClienteId"`
	Importe             uint64               `json:"importe"`
	Numero              string               `json:"numero"`
	RazonSocial         string               `json:"razon_social"`
	Domicilio           string               `json:"domicilio"`
	Cuit                string               `json:"cuit"`
	Gravamen            string               `json:"gravamen"`
	EmitidoEl           time.Time            `json:"emitido_el"`
	ComprobanteDetalles []ComprobanteDetalle `gorm:"foreignkey:ComprobanteId"`
	MovimientosId       string               `json:"movimientos_id" gorm:"type:blob"`
	ReporteId           uint                 `json:"reporte_id"`
	Reporte             Reporte              `json:"reporte" gorm:"foreignKey:ReporteId"`
	RutaFile            string               `json:"ruta_file"`
}

type ComprobanteDetalle struct {
	gorm.Model
	ComprobanteId  uint        `json:"comprobante_id"`
	Comprobante    Comprobante `json:"comprobante" gorm:"foreignKey:ComprobanteId"`
	TotalRetencion Monto       `json:"total_retencion"`
	TotalMonto     Monto       `json:"total_monto"`
	CodigoRegimen  string      `json:"codigo_regimen"`
	Gravamen       string      `json:"gravamen"`
	Retener        bool        `json:"retener"`
}

type Comprobantes []Comprobante

// obtener de un []Comprobante los que se evaluaron como retenibles por superar el monto minimo imponible 
func (cs Comprobantes) GetComprobantesARetener() (comprobantesOut []Comprobante){
  if len(cs) == 0 {
    return
  }
  var retener bool
  for _, comprobante := range cs {
    // recorrer cada detalle
    for _, cd := range comprobante.ComprobanteDetalles {
      if cd.Retener {
        retener = true
      }
    }
    // los detalles de un comprobante tienen el mismo valor en retener
    if retener {
      comprobantesOut = append(comprobantesOut, comprobante)
    }
    retener = false
  }

  return
}
