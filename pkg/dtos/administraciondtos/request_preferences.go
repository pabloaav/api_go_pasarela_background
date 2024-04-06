package administraciondtos

import (
	"fmt"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
)

type RequestPreferences struct {
	ClientId       string `json:"clientId"`
	MainColor      string `json:"mainColor"`
	SecondaryColor string `json:"secondaryColor"`
	RutaLogo       string `json:"ruta_logo"`
}

func (request *RequestPreferences) Validar() (erro error) {
	mensaje := "los parametros enviados no son válidos"
	clienteId, err := strconv.Atoi(request.ClientId)
	if err != nil {
		return fmt.Errorf(mensaje)
	}
	if clienteId < 1 {
		return fmt.Errorf(mensaje)
	}

	if commons.StringIsEmpity(request.MainColor) {
		return fmt.Errorf(mensaje)
	}
	if commons.StringIsEmpity(request.SecondaryColor) {
		return fmt.Errorf(mensaje)
	}

	return
}
