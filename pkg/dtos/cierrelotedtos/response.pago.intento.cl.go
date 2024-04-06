package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type ResponseLIstaPagoIntentoCl struct {
	TotalPagosIntento        int64                    `json:"total_pagos_intento"`
	TotalCierreLote          int64                    `json:"totla_cierre_lote"`
	TotalTransaccionFaltante int64                    `json:"total_transaccion_faltante"`
	DatosPagosIntentoCl      []ResponsePagoIntentosCl `json:"datos_pagos_intento_cl"`
}

type ResponsePagoIntentosCl struct {
	PagoIntentoId        int64                 `json:"pago_intento_id"`
	PagosID              int64                 `json:"pagos_id"`
	ExternalID           string                `json:"external_id"`
	PaidAt               time.Time             `json:"paid_at"`
	ReportAt             time.Time             `json:"report_at"`
	IsAvailable          bool                  `json:"is_available"`
	Amount               float64               `json:"amount"`
	Valorcupon           float64               `json:"valorcupon"`
	StateComment         string                `json:"state_comment"`
	AvailableAt          time.Time             `json:"available_at"`
	RevertedAt           time.Time             `json:"reverted_at"`
	AuthorizationCode    string                `json:"authorization_code"`
	TransactionID        string                `json:"transaction_id"`
	SiteID               int64                 `json:"site_id"`
	Calculado            bool                  `json:"calculado"`
	TicketNumber         string                `json:"ticket_number"`
	MediopagosID         int64                 `json:"mediopagos_id"`
	ChannelsID           int64                 `json:"channels_id"`
	Mediopago            string                `json:"mediopago"`
	TelcoExternalID      string                `json:"telco_external_id"`
	InstallmentdetailsID int64                 `json:"installmentdetails_id"`
	Cuota                int64                 `json:"cuota"`
	Tna                  float64               `json:"tna"`
	Tem                  float64               `json:"tem"`
	Coeficiente          float64               `json:"coeficiente"`
	CierreLote           CierreLotePagoIntento `json:"cierre_lote"`
}

func (rpicl *ResponsePagoIntentosCl) EntityToDtos(entityPagoIntento entities.Pagointento) {
	rpicl.PagoIntentoId = int64(entityPagoIntento.Model.ID)
	rpicl.PagosID = entityPagoIntento.PagosID
	rpicl.ExternalID = entityPagoIntento.ExternalID
	rpicl.PaidAt = entityPagoIntento.PaidAt
	rpicl.ReportAt = entityPagoIntento.ReportAt
	rpicl.IsAvailable = entityPagoIntento.IsAvailable
	rpicl.Amount = entityPagoIntento.Amount.Float64()
	rpicl.Valorcupon = entityPagoIntento.Valorcupon.Float64()
	rpicl.StateComment = entityPagoIntento.StateComment
	rpicl.AvailableAt = entityPagoIntento.AvailableAt
	rpicl.RevertedAt = entityPagoIntento.RevertedAt
	rpicl.AuthorizationCode = entityPagoIntento.AuthorizationCode
	rpicl.TransactionID = entityPagoIntento.TransactionID
	rpicl.SiteID = entityPagoIntento.SiteID
	rpicl.Calculado = entityPagoIntento.Calculado
	rpicl.TicketNumber = entityPagoIntento.TicketNumber
	rpicl.MediopagosID = entityPagoIntento.MediopagosID
	rpicl.ChannelsID = entityPagoIntento.Mediopagos.ChannelsID
	rpicl.Mediopago = entityPagoIntento.Mediopagos.Mediopago
	rpicl.TelcoExternalID = entityPagoIntento.Mediopagos.ExternalID
	rpicl.InstallmentdetailsID = entityPagoIntento.InstallmentdetailsID
	rpicl.Cuota = entityPagoIntento.Installmentdetail.Cuota
	rpicl.Tna = entityPagoIntento.Installmentdetail.Tna
	rpicl.Tem = entityPagoIntento.Installmentdetail.Tem
	rpicl.Coeficiente = entityPagoIntento.Installmentdetail.Coeficiente
}

type CierreLotePagoIntento struct {
	ClId                       int64     `json:"cl_id"`
	PagoestadoexternosId       int64     `json:"pagoestadoexternos_id"`
	ChannelarancelesId         int64     `json:"channelaranceles_id"`
	PrismamovimientodetallesId int64     `gorm:"default:(null),foreignkey:prismamovimientodetalles_id"`
	PrismatrdospagosId         int64     `gorm:"default:(null),foreignkey:prismatrdospagos_id"`
	BancoExternalId            int64     `json:"banco_external_id"`
	Tiporegistro               string    `json:"tiporegistro"`
	PagosUuid                  string    `json:"pagos_uuid"`
	ExternalmediopagoId        int64     `json:"externalmediopago"`
	Tipooperacion              string    `json:"tipooperacion"`
	Fechaoperacion             time.Time `json:"fechaoperacion"`
	Monto                      float64   `json:"monto"`
	Montofinal                 float64   `json:"montofinal"`
	Valorpresentado            float64   `json:"valorpresentado"`
	Diferenciaimporte          float64   `json:"difernciaimporte"`
	Coeficientecalculado       float64   `json:"coeficientecalculado"`
	Costototalporcentaje       float64   `json:"costototalporcentaje"`
	Importeivaarancel          float64   `json:"importeivaarancel"`
	ImportearancelCalculado    float64   `json:"importearancel_calculado"`
	ImporteivaArancelCalculado float64   `json:"importeiva_arancel_calculado"`
	ImporteCfPrisma            float64   `json:"importe_cf_prisma"`
	ImporteIvaCfCalculado      float64   `json:"importe_iva_cf_calculado"`
	Descripcion                string    `json:"descripcion"`
	Codigoautorizacion         string    `json:"codigoautorizacion"`
	Nroticket                  int64     `json:"nroticket"`
	SiteID                     int64     `json:"site_id"`
	ExternalloteId             int64     `json:"externallote_id"`
	Nrocuota                   int64     `json:"nrocuota"`
	FechaCierre                time.Time `json:"fecha_cierre"`
	Nroestablecimiento         int64     `json:"nroestablecimiento"`
	ExternalclienteID          string    `json:"externalcliente_id"`
	Nombrearchivolote          string    `json:"nombrearchivolote"`
	Match                      int       `json:"match"`
	FechaPago                  time.Time `json:"fecha_pago"`
	Disputa                    bool      `json:"disputa"`
	Cantdias                   int64     `json:"cantdias"`
	Enobservacion              bool      `json:"enobservacion"`
	// Descripcionpresentacion    string    `json:"descripcionpresentacion"`
	// Reversion                  bool      `json:"reversion"`
	// DetallemovimientoId        int64     `json:"detallemovimineto_id"`
	// DetallepagoId              int64     `json:"detallepago_id"`
	// Descripcioncontracargo     string    `json:"descripcioncontracargo"`
	// ExtbancoreversionId        int64     `json:"extbancoreversion_id"`
	// Conciliado                 bool      `json:"conciliacion"`
	// Estadomovimiento           bool      `json:"estadomovimineto"`
	// Descripcionbanco           string    `json:"descripcionbanco"`
	// Prismamovimientodetalle    *Prismamovimientodetalle `json:"prismamovimientodetalles_id" gorm:"foreignKey:PrismamovimientodetallesId"`
	// Prismatrdospagos           *Prismatrdospago         `json:"prismatrdospagos_id" gorm:"foreignKey:PrismatrdospagosId"`
	// Channelarancel             *Channelarancele         `json:"channelarancel" gorm:"foreignKey:channelaranceles_id"`
}

func (clpi *CierreLotePagoIntento) EntityToDtos(entityCl entities.Prismacierrelote) {
	clpi.ClId = int64(entityCl.Model.ID)
	clpi.PagoestadoexternosId = entityCl.PagoestadoexternosId
	clpi.ChannelarancelesId = entityCl.ChannelarancelesId
	clpi.PrismamovimientodetallesId = entityCl.PrismamovimientodetallesId
	clpi.PrismatrdospagosId = entityCl.PrismatrdospagosId
	clpi.BancoExternalId = entityCl.BancoExternalId
	clpi.Tiporegistro = entityCl.Tiporegistro
	clpi.PagosUuid = entityCl.PagosUuid
	clpi.ExternalmediopagoId = entityCl.ExternalmediopagoId
	clpi.Tipooperacion = string(entityCl.Tipooperacion)
	clpi.Fechaoperacion = entityCl.Fechaoperacion
	clpi.Monto = entityCl.Monto.Float64()
	clpi.Montofinal = entityCl.Montofinal.Float64()
	clpi.Valorpresentado = entityCl.Valorpresentado.Float64()
	clpi.Diferenciaimporte = entityCl.Diferenciaimporte.Float64()
	clpi.Coeficientecalculado = entityCl.Coeficientecalculado
	clpi.Costototalporcentaje = entityCl.Costototalporcentaje
	clpi.Importeivaarancel = entityCl.Importeivaarancel
	clpi.ImportearancelCalculado = entityCl.ImportearancelCalculado
	clpi.ImporteivaArancelCalculado = entityCl.ImporteivaArancelCalculado
	clpi.ImporteCfPrisma = entityCl.ImporteCfPrisma
	clpi.ImporteIvaCfCalculado = entityCl.ImporteIvaCfCalculado
	clpi.Descripcion = entityCl.Descripcion
	clpi.Codigoautorizacion = entityCl.Codigoautorizacion
	clpi.Nroticket = entityCl.Nroticket
	clpi.SiteID = entityCl.SiteID
	clpi.ExternalloteId = entityCl.ExternalloteId
	clpi.Nrocuota = entityCl.Nrocuota
	clpi.FechaCierre = entityCl.FechaCierre
	clpi.Nroestablecimiento = entityCl.Nroestablecimiento
	clpi.ExternalclienteID = entityCl.ExternalclienteID
	clpi.Nombrearchivolote = entityCl.Nombrearchivolote
	clpi.Match = entityCl.Match
	clpi.FechaPago = entityCl.FechaPago
	clpi.Disputa = entityCl.Disputa
	clpi.Cantdias = entityCl.Cantdias
	clpi.Enobservacion = entityCl.Enobservacion
}
