package retenciondtos

import (
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type RentencionRequestDTO struct {
	ClienteId             uint
	RetencionId           uint
	Unlinked              bool // determina si deben devolverse las retenciones no vinculadas a un cliente
	FechaInicio, FechaFin string
	GravamensId           uint
	CargarGravamenes      bool
	ListaMovimientosId    []uint
	NumeroReporteRrm      uint
	ComprobanteId         uint
	RutaFile              string
	ReporteId             uint
}

type PostRentencionRequestDTO struct {
	Id                                   uint
	Rg2854                               bool
	Minorista                            bool
	Alicuota                             float64
	AlicuotaOpcional                     float64
	MontoMinimo                          float64
	Descripcion                          string
	CondicionsId                         uint
	CodigoRegimen                        string
	ChannelsId                           []uint
	FechaValidezDesde, FechaValidezHasta string
}

func (rrdto *RentencionRequestDTO) Validar() error {

	if rrdto.ClienteId <= 0 {
		return fmt.Errorf(ERROR_CAMPO, "ClienteId")
	}

	return nil
}

func (rrdto *RentencionRequestDTO) ValidarPost() error {

	erro := rrdto.Validar()
	if erro != nil {
		return erro
	}

	if rrdto.RetencionId <= 0 {
		return fmt.Errorf(ERROR_CAMPO, "RetencionId")
	}

	return nil
}

func (rrdto *RentencionRequestDTO) ValidarDelete() error {
	return rrdto.ValidarPost()
}

func (rrdto *RentencionRequestDTO) ValidarFechas() error {

	var (
		fechaInicioTime, fechaFinTime time.Time
	)

	if len(rrdto.FechaInicio) == 0 {
		return fmt.Errorf(ERROR_CAMPO, "FechaInicio")
	}

	if len(rrdto.FechaFin) == 0 {
		return fmt.Errorf(ERROR_CAMPO, "FechaFin")
	}

	longLayout := "2006-01-02T15:04:05Z"

	_, err := time.Parse(longLayout, rrdto.FechaInicio)
	// si hay error, es porque la fecha tiene el formato corto
	if err != nil {

		fechaInicioTime, err = time.Parse("2006-01-02", rrdto.FechaInicio)
		if err != nil {
			erro := fmt.Errorf(ERROR_CONVERTIR_FECHA, "FechaInicio")
			return erro
		}
		rrdto.FechaInicio = commons.GetDateFirstMoment(fechaInicioTime)

	}

	_, err = time.Parse(longLayout, rrdto.FechaFin)
	// si hay error, es porque la fecha tiene el formato corto
	if err != nil {

		fechaFinTime, err = time.Parse("2006-01-02", rrdto.FechaFin)
		if err != nil {
			erro := fmt.Errorf(ERROR_CONVERTIR_FECHA, "FechaFin")
			return erro
		}
		rrdto.FechaFin = commons.GetDateLastMoment(fechaFinTime)

	}

	return nil
}

func (prr *PostRentencionRequestDTO) ValidarUpSert(isUpdate bool) (erro error) {

	if isUpdate && prr.Id < 1 {
		erro = fmt.Errorf(ERROR_UPDATE)
		return
	}

	if prr.CondicionsId < 1 {
		erro = fmt.Errorf(ERROR_CAMPO, "CondicionsId")
		return
	}

	if len(prr.ChannelsId) == 0 {
		erro = fmt.Errorf(ERROR_CAMPO, "ChannelsId")
		return
	}

	if prr.Alicuota < 0 {
		erro = fmt.Errorf(ERROR_ALICUOTA, "Alicuota")
		return
	}

	if len(prr.FechaValidezDesde) > 0 && len(prr.FechaValidezDesde) != 10 {
		erro = fmt.Errorf(ERROR_CONVERTIR_FECHA, "FechaValidezDesde")
		return
	}

	if len(prr.FechaValidezHasta) > 0 && len(prr.FechaValidezHasta) != 10 {
		erro = fmt.Errorf(ERROR_CONVERTIR_FECHA, "FechaValidezHasta")
		return
	}

	return
}

func (prr *PostRentencionRequestDTO) ToEntity(isUpdate bool, channel uint) (e entities.Retencion) {
	if isUpdate {
		e.ID = prr.Id
	}
	e.CondicionsId = prr.CondicionsId
	e.ChannelsId = channel
	e.Alicuota = prr.Alicuota
	e.AlicuotaOpcional = prr.AlicuotaOpcional
	e.Rg2854 = prr.Rg2854
	e.Minorista = prr.Minorista
	e.MontoMinimo = prr.MontoMinimo
	e.Descripcion = prr.Descripcion
	e.CodigoRegimen = prr.CodigoRegimen
	if len(prr.FechaValidezDesde) > 0 {
		e.FechaValidezDesde, _ = time.Parse("2006-01-02", prr.FechaValidezDesde)
	}
	if len(prr.FechaValidezDesde) == 0 {
		e.FechaValidezDesde = commons.GetDateFirstMomentTime(time.Now())
	}
	if len(prr.FechaValidezHasta) > 0 {
		e.FechaValidezHasta, _ = time.Parse("2006-01-02", prr.FechaValidezHasta)
		e.FechaValidezHasta = commons.GetDateLastMomentTime(e.FechaValidezHasta)
	}

	return
}

func (prr *PostRentencionRequestDTO) ToEntitiesByChannel(isUpdate bool) (es []entities.Retencion) {

	for _, channel := range prr.ChannelsId {
		e := prr.ToEntity(isUpdate, channel)
		es = append(es, e)
	}

	return
}
