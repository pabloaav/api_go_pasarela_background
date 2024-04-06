package linkdtos

import (
	"errors"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/tools"
)

type PagoInformado string

const (
	Informado   PagoInformado = "Informado"
	NoInformado PagoInformado = "No informado"
)

func (e PagoInformado) IsValid() error {
	switch e {
	case Informado, NoInformado:
		return nil
	}
	return errors.New(tools.ERROR_PAGO_INFORMADO)
}
