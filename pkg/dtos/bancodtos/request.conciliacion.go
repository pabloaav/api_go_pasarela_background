package bancodtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

type RequestConciliacion struct {
	Transferencias   administraciondtos.TransferenciaResponsePaginado
	ListaApilink     []*entities.Apilinkcierrelote
	ListaRapipago    []*entities.Rapipagocierrelote
	ListaMultipagos  []*entities.Multipagoscierrelote
	TipoConciliacion int64
}
