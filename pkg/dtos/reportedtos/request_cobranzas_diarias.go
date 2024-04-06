package reportedtos

import "time"

type RequestCobranzasDiarias struct {
	FechaInicio time.Time `json:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"`
	Email       []string    `json:"email"`
}
