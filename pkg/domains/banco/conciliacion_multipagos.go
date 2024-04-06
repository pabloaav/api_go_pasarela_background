package banco

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/bancodtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/filtros/banco"
)

type multipagosConciliacion struct {
	utilService util.UtilService
}

func NewMultipagosConciliacion(util util.UtilService) MetodoConciliacionPagos {
	return &multipagosConciliacion{
		utilService: util,
	}
}

// rutaArchivos string, estadosPagoExterno []entities.Pagoestadoexterno
func (cl *multipagosConciliacion) FiltroRequestConsultaBanco(request bancodtos.RequestConciliacion) (response filtros.MovimientosBancoFiltro) {

	logs.Info("--------- Ejecutando proceso conciliacion multipagos --------------- ")
	var fechas []string
	var tipo filtros.EnumTipoOperacion
	if request.ListaMultipagos != nil {
		for _, referencia := range request.ListaMultipagos {
			fecha_acreditacion := referencia.Fechaacreditacion.Format("2006-01-02")
			fechas = append(fechas, fecha_acreditacion)
		}
		tipo = "multipagos"
	}
	response = filtros.MovimientosBancoFiltro{
		SubCuenta: config.COD_SUBCUENTA,
		Tipo:      tipo,
		Fechas:    fechas, // este filtro debo aplicar cuando tenga las fechas precisas (cuando el dinero ingresara al banco)
	}
	logs.Info(fechas)

	return
}

func (cl *multipagosConciliacion) ConciliacionBanco(request bancodtos.RequestConciliacion, listaMovimientosBanco []bancodtos.ResponseMovimientosBanco) (lista bancodtos.ResponseConciliacion, movimientosIds []uint) {

	// se debe verificar la fecha y el monto que ingresaron al banco
	// si el monto no coinicide este pago pasara a estar en observacion , se debe verificar el el arancel (proveedor)
	for _, clmultipagos := range request.ListaMultipagos {

		logs.Info(clmultipagos.Fechaacreditacion.Format("2006-01-02T00:00:00Z"))
		// 2022-08-02 00:00:00 +0000 UTC format fecha para comparar
		for _, mv := range listaMovimientosBanco {
			logs.Info(mv.Fecha)

			if clmultipagos.Fechaacreditacion.Format("2006-01-02T00:00:00Z") == mv.Fecha {
				// los id de banco para luego actualizar (match en banco)
				movimientosIds = append(movimientosIds, mv.Id)
				// cabecera cierrelote
				clmultipagos.BancoExternalId = int64(mv.Id)
				// & calcula la diferencia en monto del dinero que ingreso al banco con el monto calculado en el cierrelote
				clmultipagos.Difbancocl = float64(mv.Importe) - clmultipagos.ImporteTotalCalculado
				if entities.Monto(clmultipagos.ImporteTotalCalculado) != entities.Monto(mv.Importe) {
					// en el caso de que no coincidan los importes, se debera analizar este caso
					// arancel mal cobrado o calculado
					clmultipagos.Enobservacion = true
				}

				for _, cldetalle := range clmultipagos.MultipagosDetalle {
					cldetalle.Match = true
					if clmultipagos.Enobservacion {
						cldetalle.Enobservacion = true
					}
				}
			}
		}
	}
	lista.ListaMultipagos = request.ListaMultipagos

	return
}
