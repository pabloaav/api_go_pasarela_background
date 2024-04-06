package multipagos

type RequestConsultarMovimientosMultipagos struct {
	CargarMovConciliados bool
	PagosNotificado      bool
}

type RequestConsultarMovimientosMultipagosDetalles struct {
	PagosInformados bool
}
