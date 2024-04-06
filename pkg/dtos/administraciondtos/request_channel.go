package administraciondtos

import (
	"fmt"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/tools"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type RequestChannel struct {
	Id         uint
	Channel    string
	Nombre     string
	CodigoBcra int32
}

func (r *RequestChannel) ToChannel(cargarId bool) (response entities.Channel) {
	if cargarId {
		response.ID = r.Id
	}
	response.Channel = r.Channel
	response.Nombre = r.Nombre
	response.CodigoBcra = r.CodigoBcra

	return
}

func (r *RequestChannel) IsVAlid(isUpdate bool) (erro error) {

	if isUpdate && r.Id < 1 {
		return fmt.Errorf(tools.ERROR_ID)
	}

	if commons.StringIsEmpity(r.Channel) {
		return fmt.Errorf(tools.ERROR_CHANNEL)
	}

	if commons.StringIsEmpity(r.Nombre) {
		return fmt.Errorf(tools.ERROR_NOMBRE_CHANNEL)
	}

	if r.CodigoBcra == 0 {
		return fmt.Errorf(tools.ERROR_CODIGO_BCRA)
	}

	r.Channel = strings.ToUpper(r.Channel)

	return
}
