package filtros

import (
	"errors"
	"time"
)

type FiltroPagosCl struct {
	FechaInicio time.Time
	FechaFin    time.Time
	Pagos       bool
	PagoIntento bool
	CierreLote  bool
}

func (fpcl *FiltroPagosCl) Validar() (erro error) {

	if fpcl.FechaInicio.IsZero() {
		erro = errors.New("no es una fecha valida el parametro fecha de inicio")
	}
	if fpcl.FechaFin.IsZero() {
		erro = errors.New("no es una fecha valida el parametro fecha de inicio")
	}

	return
}
