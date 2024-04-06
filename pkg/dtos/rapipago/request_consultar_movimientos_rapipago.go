package rapipago

type RequestConsultarMovimientosRapipago struct {
	CargarMovConciliados bool
	PagosNotificado      bool
}

type RequestConsultarMovimientosRapipagoDetalles struct {
	PagosInformados bool
}
