package filtros

import (
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/cierrelotedtos"
)

const (
	layoutISO = "2006-01-02"
)

type OpcionConciliados uint32

type ApilinkCierreloteFiltro struct {
	FechaInicio     string            // query:FechaInicio db:created_at
	FechaFin        string            // query:FechaFin db:created_at
	Number          uint32            // query:Number
	Size            uint32            // query:Size
	ReferenciaBanco string            // query:ReferenciaBanco
	Conciliados     OpcionConciliados // query:Conciliados, 0 no se busca por conciliados o no conciliados, 1 para buscar los que estan conciliados, 2 para buscar los que no estan conciliados
}

func (aclf *ApilinkCierreloteFiltro) Validar() error {
	if len(strings.TrimSpace(aclf.FechaInicio)) <= 0 {
		return fmt.Errorf("parametros enviados la fecha de inicio no puede ser vacío")
	}
	if len(strings.TrimSpace(aclf.FechaFin)) <= 0 {
		return fmt.Errorf("parametros enviados la fecha de fin no puede ser vacío")
	}
	switch aclf.Conciliados {
	case 0, 1, 2:
		return nil
	default:
		return fmt.Errorf("campo conciliados inválido")
	}
}

// Realiza el formateo de las fechas que recibe y carga los datos recibidos al objeto que devuelve.
func (aclf *ApilinkCierreloteFiltro) ToFiltroRequest() (apilinkRequest cierrelotedtos.ApilinkRequest) {
	apilinkRequest.FechaInicio, _ = time.Parse(layoutISO, aclf.FechaInicio)
	apilinkRequest.FechaFin, _ = time.Parse(layoutISO, aclf.FechaFin)
	apilinkRequest.FechaFin = commons.GetDateLastMomentTime(apilinkRequest.FechaFin)
	apilinkRequest.Number = aclf.Number
	apilinkRequest.Size = aclf.Size
	apilinkRequest.Conciliados = uint32(aclf.Conciliados)
	apilinkRequest.ReferenciaBanco = aclf.ReferenciaBanco
	return
}
