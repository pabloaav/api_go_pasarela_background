package filtros

type CobranzasClienteFiltro struct {
	FechaInicio string
	FechaFin    string
	ClienteId   int `json:"cliente_id"`
}
