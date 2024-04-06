package reportedtos

import "time"

type RequestGetReportes struct {
	FechaInicio, FechaFin string
	Number                uint            `json:"number"`
	Size                  uint            `json:"size"`
	TipoReporte           EnumTipoReporte `json:"tipo_reporte"`
	Cliente               string          `json:"cliente"`
	NroReporte            uint            `json:"nro_reporte"`
	Fechacobranza         string          `json:"fechacobranza"`
	Fecharendicion        string          `json:"fecharendicion"`
	PeriodoInicio         time.Time       `json:"periodo_inicio"`
	PeriodoFin            time.Time       `json:"periodo_fin"`
	LastRrm               bool            `json:"last_rrm"`
}
