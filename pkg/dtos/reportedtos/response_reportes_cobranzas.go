package reportedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type ResponseCobranzasClientes struct {
	CantidadCobranzas int
	Total             uint
	Cobranzas         []DetallesCobranza
}

type DetallesCobranza struct {
	Fecha     string
	Nombre    string
	Registros uint
	Subtotal  uint
	Pagos     []DetallesPagosCobranza
}

type DetallesPagosCobranza struct {
	Id            int                  `json:"id"`
	Cliente       string               `json:"cliente"`
	Pagoestado    string               `json:"pagoestado"`
	Descripcion   string               `json:"descripcion"`
	Referencia    string               `json:"referencia"`
	PayerName     string               `json:"payer_name"`
	PayerEmail    string               `json:"payer_email"`
	TotalPago     uint                 `json:"total_pago"`
	MedioPago     string               `json:"medio_pago"`
	CanalPago     string               `json:"canal_pago"`
	Cuenta        string               `json:"cuenta"`
	FechaPago     time.Time            `json:"fecha_pago"`
	FechaCobro    time.Time            `json:"fecha_cobro"`
	Comision      uint                 `json:"comision"` // comision solo telco
	ComisionTotal uint                 `json:"comision_total"`
	Iva           uint                 `json:"iva"`
	Retencion     uint                 `json:"retencion"`
	Lote          uint                 `json:"lote"`
	PagoItems     []entities.Pagoitems `json:"pago_items" gorm:"foreignKey:PagosID"`
}
