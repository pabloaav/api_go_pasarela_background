package administracion

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

// se devuelve algun objeto que implementa la IStrategyRetencion
func RetencionFactory(s string, retencion entities.Retencion) (IStrategyRetencion, error) {
	switch s {
	case "iva":
		retIva := RetencionIva{Retencion: retencion}
		return &retIva, nil
	case "iibb":
		retIibb := RetencionIibb{Retencion: retencion}
		return &retIibb, nil
	case "ganancias":
		retGanancias := RetencionGanancias{Retencion: retencion}
		return &retGanancias, nil
	default:
		return nil, fmt.Errorf("%s", "error: lo ingresado no corresponde con ningun tipo de retencion")
	}
}
