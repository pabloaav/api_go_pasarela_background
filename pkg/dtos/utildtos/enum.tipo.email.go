package utildtos

import "errors"

type EnumTipoEmail string

const (
	Template  EnumTipoEmail = "template"
	Adjunto   EnumTipoEmail = "adjunto"
	Reporte   EnumTipoEmail = "reporte"
	Retencion EnumTipoEmail = "mixto"
)

func (e EnumTipoEmail) IsValid() (int, error) {
	switch e {
	case Template:
		return 1, nil
	case Adjunto:
		return 2, nil
	case Reporte:
		return 3, nil

	case Retencion:
		return 4, nil
	}

	return 0, errors.New("el tipo de correo no es valido")
}
