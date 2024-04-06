package retenciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type RetencionAgrupada struct {
	ClienteId      uint           `json:"cliente_id"`
	TotalRetencion entities.Monto `json:"total_retencion"`
	TotalMonto     entities.Monto `json:"total_monto"`
	Gravamen       string         `json:"gravamen"`
	FechaInicio    string         `json:"fecha_inicio"`
	FechaFin       string         `json:"fecha_fin"`
	Minimo         float64        `json:"minimo"`
	Retener        bool           `json:"retener"` // si el importe debe considerarse como una retencion efectiva
	MovimientoId   uint           `json:"movimiento_id"`
	Efectuada      bool           `json:"efectuada"`
	CodigoRegimen  string         `json:"codigo_regimen"`
}

type RentencionesResponseDTO struct {
	RetencionesDTO []RentencionResponseDTO `json:"retenciones"`
	Count          int64                   `json:"count"`
	Retenciones    []entities.Retencion    `json:"entities_retenciones"`
	// Meta                  dtos.Meta               `json:"meta"`
}

type RentencionResponseDTO struct {
	Id                uint                     `json:"id"`
	CondicionsId      uint                     `json:"condicions_id"`
	Condicion         CondicionResponseDTO     `json:"condicion"`
	ChannelsId        uint                     `json:"channels_id"`
	CanalDePago       string                   `json:"canal_de_pago"`
	Alicuota          float64                  `json:"alicuota"`
	AlicuotaOpcional  float64                  `json:"alicuota_opcional"`
	Rg2854            bool                     `json:"rg2854"`
	Minorista         bool                     `json:"minorista"`
	Monto_Minimo      float64                  `json:"monto_minimo"`
	CodigoRegimen     string                   `json:"codigo_regimen"`
	Certificados      []CertificadoResponseDTO `json:"certificados"`
	FechaValidezDesde time.Time                `json:"fecha_validez_desde"`
	FechaValidezHasta time.Time                `json:"fecha_validez_hasta"`
}

type CalculoRentencionResponseDTO struct {
	Iva       float64 `json:"iva"`
	Iibb      float64 `json:"Iibb"`
	Ganancias float64 `json:"ganancias"`
}
type CondicionResponseDTO struct {
	Id          uint                `json:"id"`
	Condicion   string              `json:"condicion"`
	GravamensId uint                `json:"gravamens_id"`
	Gravamen    GravamenResponseDTO `json:"gravamen"`
	Descripcion string              `json:"descripcion"`
}

type GravamenResponseDTO struct {
	Id             uint   `json:"id"`
	Gravamen       string `json:"gravamen"`
	CodigoGravamen string `json:"codigo_gravamen"`
	Descripcion    string `json:"descripcion"`
}

type CertificadoResponseDTO struct {
	Id                  uint      `json:"id"`
	ClienteRetencionsId uint      `json:"cliente_retencions_id"`
	Fecha_Presentacion  time.Time `json:"fecha_presentacion,omitempty"`
	Fecha_Caducidad     time.Time `json:"fecha_caducidad,omitempty"`
	CreadoEl            time.Time `json:"creado_el,omitempty"`
	Ruta_file           string    `json:"ruta_file"`
	IsExpired           bool      `json:"expired"`
}

type ComprobanteResponseDTO struct {
	Id            uint                            `json:"id"`
	ClienteId     uint                            `json:"cliente_id"`
	Importe       uint64                          `json:"importe"`
	Numero        string                          `json:"numero"`
	RazonSocial   string                          `json:"razon_social"`
	Cuit          string                          `json:"cuit"`
	Gravamen      string                          `json:"gravamen"`
	EmitidoEl     time.Time                       `json:"emitido_el"`
	Detalles      []ComprobanteDetalleResponseDTO `json:"detalles"`
	MovimientosId []uint                          `json:"movimientos_id"`
	RutaFile      string                          `json:"ruta_file"`
	ReporteRrm    ReporteRrm                      `json:"reporte_rrm"`
}

type ReporteRrm struct {
	Id       uint   `json:"id"`
	RutaFile string `json:"ruta_file"`
}

type ComprobanteDetalleResponseDTO struct {
	Id             uint    `json:"id"`
	TotalRetencion float64 `json:"total_retencion"`
	CodigoRegimen  string  `json:"codigo_regimen"`
	Gravamen       string  `json:"gravamen"`
	Retener        bool    `json:"retener"`
}

type DevolverRetencionesDTO struct {
	ListaComprobantes []ComprobanteResponseDTO
	MontoDevolver     uint64
}

func (crdto *CertificadoResponseDTO) FromEntity(cert entities.Certificado) {
	crdto.Id = cert.ID
	crdto.ClienteRetencionsId = cert.ClienteRetencionsId
	crdto.Fecha_Presentacion = cert.Fecha_Presentacion
	crdto.Fecha_Caducidad = cert.Fecha_Caducidad
	crdto.CreadoEl = cert.CreatedAt
	crdto.Ruta_file = cert.Ruta_file
	crdto.IsExpired = cert.IsExpired()
}

func (grdto *GravamenResponseDTO) FromEntity(g entities.Gravamen) {
	grdto.Id = g.Model.ID
	grdto.CodigoGravamen = g.CodigoGravamen
	grdto.Gravamen = g.Gravamen
	grdto.Descripcion = g.Descripcion
}

func (crdto *CondicionResponseDTO) FromEntity(c entities.Condicion) {
	crdto.Id = c.Model.ID
	crdto.Condicion = c.Condicion
	crdto.GravamensId = c.GravamensId
	crdto.Descripcion = c.Descripcion
	// GravamenResponseDTO FromEntity
	if c.Gravamen.ID != 0 {
		var tempGravamenDTO GravamenResponseDTO
		tempGravamenDTO.FromEntity(c.Gravamen)
		crdto.Gravamen = tempGravamenDTO
	}
}

func (rrdto *RentencionResponseDTO) FromEntity(r entities.Retencion) {
	rrdto.Id = r.Model.ID
	rrdto.Alicuota = r.Alicuota
	rrdto.AlicuotaOpcional = r.AlicuotaOpcional
	rrdto.ChannelsId = r.ChannelsId
	rrdto.CanalDePago = r.Channel.Nombre
	rrdto.CondicionsId = r.CondicionsId
	// CondicionResponseDTO FromEntity
	if r.Condicion.ID != 0 {
		var tempCondicionDTO CondicionResponseDTO
		tempCondicionDTO.FromEntity(r.Condicion)
		rrdto.Condicion = tempCondicionDTO
	}
	rrdto.Rg2854 = r.Rg2854
	rrdto.Minorista = r.Minorista
	rrdto.Monto_Minimo = r.MontoMinimo
	// debe haber una sola ClienteRetencion por cada Retencion
	if len(r.ClienteRetencions) > 0 && len(r.ClienteRetencions[0].Certificados) > 0 {
		certificates := r.ClienteRetencions[0].Certificados
		// para cada certificado
		for _, cert := range certificates {
			var tempCertificadoDTO CertificadoResponseDTO
			tempCertificadoDTO.FromEntity(cert)
			rrdto.Certificados = append(rrdto.Certificados, tempCertificadoDTO)
		}
	}
	rrdto.CodigoRegimen = r.CodigoRegimen
	rrdto.FechaValidezDesde = r.FechaValidezDesde
	rrdto.FechaValidezHasta = r.FechaValidezHasta
}

func (crdto *ComprobanteResponseDTO) FromEntity(c entities.Comprobante) {
	crdto.Id = c.ID
	crdto.ClienteId = c.ClienteId
	crdto.Importe = c.Importe
	crdto.Numero = c.Numero
	crdto.RazonSocial = c.RazonSocial
	crdto.Cuit = c.Cuit
	crdto.Gravamen = c.Gravamen
	crdto.EmitidoEl = c.EmitidoEl
	crdto.RutaFile = c.RutaFile

	movs, err := commons.StringToUintSliceNumber(c.MovimientosId)
	if err != nil {
		logs.Error("en ComprobanteResponseDTO.FromEntity se produjo un error: " + err.Error())
	}
	crdto.MovimientosId = movs
	if len(c.ComprobanteDetalles) > 0 {
		for _, cd := range c.ComprobanteDetalles {
			var tempComprobanteDetallesDTO ComprobanteDetalleResponseDTO
			tempComprobanteDetallesDTO.FromEntity(cd)
			crdto.Detalles = append(crdto.Detalles, tempComprobanteDetallesDTO)
		}
	}
	if c.Reporte.ID != 0 {
		crdto.ReporteRrm.Id = c.Reporte.ID
		crdto.ReporteRrm.RutaFile = c.Reporte.RutaFile
	}
}

func (cdrdto *ComprobanteDetalleResponseDTO) FromEntity(cd entities.ComprobanteDetalle) {
	cdrdto.Id = cd.ID
	cdrdto.TotalRetencion = cd.TotalRetencion.Float64()
	cdrdto.CodigoRegimen = cd.CodigoRegimen
	cdrdto.Gravamen = cd.Gravamen
	cdrdto.Retener = cd.Retener
}
