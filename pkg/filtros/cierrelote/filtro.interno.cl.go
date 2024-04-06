package filtros

type FiltroInternoCl struct {
	DataFiltro []FiltroConsultaInterna
}

type FiltroConsultaInterna struct {
	CodigoAutorizacion string
	TicketNro          string
	IdOperacion        string
}

func (fci *FiltroConsultaInterna) CargarData(codigoAtorizacion, ticket, idOperacion string) {
	fci.CodigoAutorizacion = codigoAtorizacion
	fci.TicketNro = ticket
	fci.IdOperacion = idOperacion
}
