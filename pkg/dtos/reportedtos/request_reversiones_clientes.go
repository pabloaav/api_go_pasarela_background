package reportedtos

import (
	"errors"
	"regexp"
)

type RequestReversionesClientes struct {
	Cliente     uint   `json:"cliente_id"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
	//FechaAdicional       time.Time `json:"fecha_adicional"`
	// FechaReversion time.Time `json:"fecha_reversion"`
	//CargarFechaAdicional bool
	EnviarEmail    bool
	ClientesIds    []uint `json:"clientes_ids"`
	ClientesString string `json:"clientes_string"`
	//OrdenMayorCobranza   bool   `json:"orden_mayor_cobranza"`
	CuentaId uint
}

/*
func (r *RequestReversionesClientes) ValidarFechas() (estadoValidacion ValidacionesFiltro, erro error) {
	estadoValidacion.Cliente = true
	estadoValidacion.Inicio = true
	estadoValidacion.Fin = true
	if r.FechaInicio.IsZero() && r.FechaFin.IsZero() {
		erro = errors.New("por lomenos debe enviar una fecha de inicio")
		return
	}
	if !r.FechaInicio.IsZero() && !r.FechaFin.IsZero() {
		if r.FechaInicio.After(r.FechaFin) {
			erro = errors.New("la fecha de inicio no puede ser mayor que la fecha fin ")
			return
		}
	}
	if r.Cliente == 0 {
		estadoValidacion.Cliente = false
	}
	return
} */

func (r *RequestReversionesClientes) ValidarFechaString() error {

	/* expresion regular para velidar fecha -> formato: año/mes/dia (20210330)*/
	regularCheckFecha := regexp.MustCompile(`([0-2][0-9]|3[0-1])(-)(0[1-9]|1[0-2])(-)(\d{4})$`)
	// regexp.MustCompile(`(\d{4})(-)(0[1-9]|1[0-2])(-)([0-2][0-9]|3[0-1])$`)
	// regularCheckHora := regexp.MustCompile(`([0-2][0-9])(:)([0-5][0-9])(:)([0-5][0-9])$`)

	if len(r.FechaInicio) <= 0 && len(r.FechaFin) <= 0 {
		return errors.New("se debe indicar una fecha")
	}
	/* FECHA */
	if len(r.FechaInicio) > 10 || !regularCheckFecha.MatchString(r.FechaInicio) {
		return errors.New("error en el formato de fecha enviado")
	}
	return nil
}
