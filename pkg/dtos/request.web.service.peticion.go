package dtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/enumsdtos"

type RequestWebServicePeticion struct {
	Operacion string
	Vendor    enumsdtos.EnumVendor
}
