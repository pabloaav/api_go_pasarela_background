package reportedtos

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type RequestPagosClientes struct {
	Cliente        uint      `json:"cliente_id"`
	FechaInicio    time.Time `json:"fecha_inicio"`
	FechaFin       time.Time `json:"fecha_fin"`
	FechaAdicional time.Time `json:"fecha_adicional"`
	// FechaReversion time.Time `json:"fecha_reversion"`
	CargarFechaAdicional bool
	EnviarEmail          bool
	ClientesIds          []uint `json:"clientes_ids"`
	ClientesString       string `json:"clientes_string"`
	OrdenMayorCobranza   bool   `json:"orden_mayor_cobranza"`
	CuentaId             uint
}

type ClientesId struct {
	Id uint `json:"id"` // Id clientereporte"`
}

func (r *RequestPagosClientes) ValidarFechas() (estadoValidacion ValidacionesFiltro, erro error) {
	estadoValidacion.Cliente = true
	estadoValidacion.Inicio = true
	estadoValidacion.Fin = true
	if r.FechaInicio.IsZero() && r.FechaFin.IsZero() {
		erro = errors.New("por lo menos debe enviar una fecha de inicio")
		return
	}
	if !r.FechaFin.IsZero() && !r.FechaInicio.IsZero() {
		if r.FechaFin.After(r.FechaInicio) {
			erro = errors.New("la fecha de inicio no puede ser mayor que la fecha fin ")
			return
		}
	} else {
		/* if r.FechaInicio.IsZero() && !r.FechaFin.IsZero() {
			erro = errors.New("la fecha de inicio no puede estar vacía ")
			return
		} */
		if r.FechaFin.IsZero() {
			r.FechaFin = r.FechaInicio
		}
	}
	if r.Cliente == 0 {
		estadoValidacion.Cliente = false
	}
	return
}

type ValidacionesFiltro struct {
	Inicio  bool
	Fin     bool
	Cliente bool
}

func (r *RequestPagosClientes) ObtenerIdsClientes() {
	var arrayString []string
	//if len(r.ClientesString) > 0 {
	arrayString = strings.Split(r.ClientesString, ",")
	//}
	for _, value := range arrayString {
		result, _ := strconv.ParseUint(value, 10, 32)
		r.ClientesIds = append(r.ClientesIds, uint(result))
	}
}
