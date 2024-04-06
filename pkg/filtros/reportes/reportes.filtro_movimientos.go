package filtros

type MovimientosComisionesFiltro struct {
	FechaInicio string
	FechaFin    string
	Number      int `json:"number"`
	Size        int `json:"size"`
	ClienteId   int `json:"cliente_id"`
	CuentaId    int `json:"cuenta_id"`
}
