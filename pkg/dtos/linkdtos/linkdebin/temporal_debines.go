package linkdebin

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/cierrelotedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type TemporalDebines struct {
	Debines *entities.Apilinkcierrelote
	Pagos   cierrelotedtos.ResponsePagosApilink // lista de pago a actualizar
}
