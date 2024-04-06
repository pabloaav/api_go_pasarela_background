package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type RunEndpoint interface {
	RunEndpoint(typeEndpoint string, urlEndpoint string, mapConHeaders map[string]string, bodyEndpoint interface{}, queryParameters map[string]string, registarPeticion bool) (objeto interface{}, erro error)
}

type runendpoint struct {
	HTTPClient  *http.Client
	UtilService UtilService
}

func NewRunEndpoint(http *http.Client, u UtilService) RunEndpoint {
	return &runendpoint{
		HTTPClient:  http,
		UtilService: u,
	}
}

func (r *runendpoint) RunEndpoint(typeEndpoint string, urlEndpoint string, mapConHeaders map[string]string, bodyEndpoint interface{}, queryParameters map[string]string, registarPeticion bool) (objeto interface{}, erro error) {
	// Construye la URL base
	base, err := url.Parse(urlEndpoint)
	if err != nil {
		logs.Error("ERROR_URL" + err.Error())
		return objeto, err
	}

	// Agrega parámetros de consulta a la URL
	q := base.Query()
	for key, value := range queryParameters {
		q.Add(key, value)
	}
	base.RawQuery = q.Encode()

	json_data, err := json.Marshal(bodyEndpoint)
	if err != nil {
		fmt.Println("Error al convertir a JSON: ", err.Error())
		return objeto, err
	}

	// Crea la solicitud HTTP
	req, _ := http.NewRequest(typeEndpoint, base.String(), bytes.NewBuffer(json_data))

	buildHeaderUniversal(req, mapConHeaders)
	// Ejecuta la solicitud y maneja la respuesta
	err = executeRequest(r, req, ERROR_RUN_ENDPOINTS, &objeto)

	// // Registra la petición realizada
	// if registarPeticion {
	// 	peticionApiLink := dtos.RequestWebServicePeticion{
	// 		Operacion: urlEndpoint,
	// 		Vendor:    "Vendor",
	// 	}
	// 	err1 := r.UtilService.CrearPeticionesService(peticionApiLink)
	// 	if err1 != nil {
	// 		logs.Error("ERROR_CREAR_PETICION" + err1.Error())
	// 	}
	// }

	if err != nil {
		return objeto, err
	}
	return objeto, nil
}

func buildHeaderUniversal(request *http.Request, nameValueHeaders map[string]string) {
	_, hasAccept := nameValueHeaders["accept"]
	_, hasContentType := nameValueHeaders["content-type"]

	if !hasAccept {
		request.Header.Add("accept", "application/json")
	}

	if !hasContentType {
		request.Header.Add("content-type", "application/json")
	}

	// Agrega los demás headers del map
	for name, value := range nameValueHeaders {
		request.Header.Add(name, value)
	}

}

func executeRequest(r *runendpoint, req *http.Request, erro string, objeto interface{}) error {

	resp, err := r.HTTPClient.Do(req)

	if err != nil {
		logs.Error(err.Error())
		return errors.New(erro)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		logs.Error(fmt.Sprint(resp.Status))
		return nil
	}

	if resp.StatusCode != 200 && resp.StatusCode != 202 {

		log := entities.Log{
			Tipo:          entities.Error,
			Funcionalidad: "executeRequestCuentas",
		}

		apiError := linkdtos.ErrorApiLink{}

		if resp.StatusCode == 500 {

			apiError.Codigo = "500"

			apiError.Descripcion = "en este momento no podemos realizar la operacion intente nuevamente mas tarde"

			log.Mensaje = fmt.Sprint(resp.Status)

		}

		err := json.NewDecoder(resp.Body).Decode(&apiError)

		if err != nil {

			apiError.Codigo = strconv.Itoa(resp.StatusCode)

			apiError.Descripcion = "en este momento no podemos realizar la operacion intente nuevamente mas tarde"

			log.Mensaje = fmt.Sprintf("%s, %s", erro, resp.Status)

		}

		if resp.StatusCode == 401 {
			apiError.Codigo = "401"
			apiError.Descripcion = "Unauthorized"
		}

		r.UtilService.CreateLogService(log)

		return &apiError
	}

	err = json.NewDecoder(resp.Body).Decode(&objeto)

	if err != nil {
		return err
	}

	return nil

}
