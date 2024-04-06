package cierrelotemultipagosdtos

type Multipagos struct {
	MultipagosHeader   HeaderTrailer
	MultipagosDetalles []Detalles
}

type HeaderTrailer struct {
	Header  Header
	Trailer Trailler
}
type Header struct {
	IdHeader      string
	NombreEmpresa string
	FechaProceso  string
	IdArchivo     string
	FillerHeader  string
}

type Trailler struct {
	IdTrailler    string
	CantDetalles  string
	ImporteTotal  string
	FillerTrailer string
}

type Detalles struct {
	FechaCobro     string
	ImporteCobrado string
	CodigoBarras   string
	Clearing       string
}
