package retenciondtos

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type CondicionRequestDTO struct {
	Id          uint
	Condicion   string
	GravamensId uint
	Descripcion string
	Exento      bool
}

func (cr *CondicionRequestDTO) ToEntity(isUpdate bool) (e entities.Condicion) {
	if isUpdate {
		e.ID = cr.Id
	}
	e.Condicion = cr.Condicion
	e.GravamensId = cr.GravamensId
	e.Descripcion = cr.Descripcion
	e.Exento = cr.Exento
	return
}

func (crdto *CondicionRequestDTO) ValidarUpSert(isUpdate bool) (erro error) {

	if isUpdate && crdto.Id < 1 {
		erro = fmt.Errorf(ERROR_UPDATE)
		return
	}

	if len(crdto.Condicion) == 0 {
		erro = fmt.Errorf(ERROR_CAMPO, "Condicion")
		return
	}

	if len(crdto.Descripcion) == 0 {
		erro = fmt.Errorf(ERROR_CAMPO, "Descripcion")
		return
	}

	if crdto.GravamensId <= 0 {
		erro = fmt.Errorf(ERROR_CAMPO, "GravamensId")
		return
	}

	return
}
