package apilink

import (
	"net/http"
	"net/url"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linkconsultadestinatario"
)

func (r *remoteRepository) GetConsultaDestinatario(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink, token string) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error) {

	base, erro := url.Parse(config.APILINKCONSULTADESTINATARIOHOST)

	if erro != nil {
		return
	}

	base.Path += "destinatarios"

	params := url.Values{}
	params.Add("cbu", request.Cbu)
	params.Add("alias", request.Alias)

	base.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", base.String(), nil)

	buildHeaderAutorizacion(req, requerimientoId, token)

	erro = executeRequest(r, req, ERROR_GET_CONSULTA_DESTINATARIOS, &response)
	/*
		 se registra la peticion realizada a la api de apilink
			-	armo el request para registrar la peticion realizada
			-	registro la peticion realizada
	*/
	peticionApiLink := dtos.RequestWebServicePeticion{
		Operacion: "GetConsultaDestinatario",
		Vendor:    "ApiLink",
	}
	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	if err1 != nil {
		logs.Error(ERROR_CREAR_PETICION + err1.Error())
	}
	return

}
