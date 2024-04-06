package administracion

import (
	"math"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

// Strategy Interface
type IStrategyRetencion interface {
	CalcularAlicuota(importe entities.Monto) (result float64)
}

// Strategy Struct
type StrategyRetencion struct {
	retentionStrategy IStrategyRetencion
}

// asigna una estrartegia concreta al patron
func (sr *StrategyRetencion) setStrategy(concreteStrategy IStrategyRetencion) {
	sr.retentionStrategy = concreteStrategy
}

// ejecutar la estrategia
func (sr *StrategyRetencion) execStrategy(importe entities.Monto) (result float64) {
	if sr.retentionStrategy != nil {
		return sr.retentionStrategy.CalcularAlicuota(importe)
	}
	return
}

/* *********** Estrategias Concretas *********** */
type RetencionIva struct {
	Retencion entities.Retencion
}

func (ri *RetencionIva) CalcularAlicuota(importe entities.Monto) (result float64) {
	if ri.Retencion.Condicion.Exento{
		return
	}
	porcentaje := ri.Retencion.Alicuota / 100
	importe_retencion := porcentaje * importe.Float64()
	result = math.Round(importe_retencion * 100)
	return 
}

type RetencionGanancias struct {
	Retencion entities.Retencion
}

func (rg *RetencionGanancias) CalcularAlicuota(importe entities.Monto) (result float64) {
	if rg.Retencion.Condicion.Exento{
		return
	}
	porcentaje := rg.Retencion.Alicuota / 100
	importe_retencion := porcentaje * importe.Float64()
	result = math.Round(importe_retencion * 100)
	return
}

type RetencionIibb struct {
	Retencion entities.Retencion
}

// El descuento esta sujeto a un monto minimo en el caso de IIBB.
func (rib *RetencionIibb) CalcularAlicuota(importe entities.Monto) (result float64) {
	if rib.Retencion.Condicion.Exento{
		return
	}
	if importe.Float64() >= rib.Retencion.MontoMinimo {
		porcentaje := rib.Retencion.Alicuota / 100
		importe_retencion := porcentaje * importe.Float64()
		result = math.Round(importe_retencion * 100)
		return
	}

	return
}
