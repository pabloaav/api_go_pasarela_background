package reportedtos

type RequestRRComprobante struct {
	FechaInicio, FechaFin string
	Id                    int64 `json:"id"`
	Cliente_id            int64 `json:"cliente_id"`
	ComprobanteId         int64
	ReporteId             uint
	GravamenesIn          []string
}
