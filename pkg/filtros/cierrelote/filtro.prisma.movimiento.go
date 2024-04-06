package filtros

import "time"

type FiltroPrismaMovimiento struct {
	IdMovimientoMxTotal          int64
	FechaPresentacion            time.Time
	EstablecimientoNro           string
	AutorizacionXlNro            string
	HashTarjeta                  string
	NroAutorizacionXl            string
	Lote                         int64
	NroCuota                     int64
	CuponNro                     int64
	Importe                      int64
	Match                        bool
	CargarDetalle                bool
	Contracargovisa              bool
	Contracargomaster            bool
	Tipooperacion                bool
	Rechazotransaccionprincipal  bool
	Rechazotransaccionsecundario bool
	Motivoajuste                 bool
	FechasPagos                  []string
	ContraCargo                  bool
	CodigosOperacion             []string
	TipoAplicacion               string
}
