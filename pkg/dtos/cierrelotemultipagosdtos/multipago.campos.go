package cierrelotemultipagosdtos

type MultipagosCamposDescripcion struct {
	NombreVariable string
	Cantidad       int
	Desde          int
	Hasta          int
}

func (rcd *MultipagosCamposDescripcion) RapipagoDescripcionHeader() (multipagoHeader []MultipagosCamposDescripcion) {
	multipagoHeader = append(multipagoHeader,
		MultipagosCamposDescripcion{NombreVariable: "Id_header", Cantidad: 8, Desde: 0, Hasta: 9},
		MultipagosCamposDescripcion{NombreVariable: "Nombre de empresa", Cantidad: 20, Desde: 9, Hasta: 30},
		MultipagosCamposDescripcion{NombreVariable: "Fh_proceso", Cantidad: 8, Desde: 30, Hasta: 39},
		MultipagosCamposDescripcion{NombreVariable: "Id_archivo", Cantidad: 20, Desde: 39, Hasta: 59},
		MultipagosCamposDescripcion{NombreVariable: "Filler", Cantidad: 17, Desde: 59, Hasta: 76},
	)
	return multipagoHeader
}

func (rcd *MultipagosCamposDescripcion) RapipagoDescripcionDetalle() (multipagoDetalle []MultipagosCamposDescripcion) {
	multipagoDetalle = append(multipagoDetalle,
		MultipagosCamposDescripcion{NombreVariable: "Fecha de Cobro", Cantidad: 8, Desde: 0, Hasta: 9},
		MultipagosCamposDescripcion{NombreVariable: "Importe Cobrado", Cantidad: 15, Desde: 9, Hasta: 29},
		MultipagosCamposDescripcion{NombreVariable: "Código de barras", Cantidad: 0, Desde: 29, Hasta: 0},
	)
	return multipagoDetalle
}

func (rcd *MultipagosCamposDescripcion) RapipagoDescripcionTrailer() (multipagoTrailer []MultipagosCamposDescripcion) {
	multipagoTrailer = append(multipagoTrailer,
		MultipagosCamposDescripcion{NombreVariable: "Id_trailer", Cantidad: 8, Desde: 0, Hasta: 9},
		MultipagosCamposDescripcion{NombreVariable: "Cant_reg", Cantidad: 8, Desde: 9, Hasta: 18},
		MultipagosCamposDescripcion{NombreVariable: "Importe_tot", Cantidad: 18, Desde: 18, Hasta: 37},
		MultipagosCamposDescripcion{NombreVariable: "Filler", Cantidad: 39, Desde: 37, Hasta: 76},
	)
	return multipagoTrailer
}
