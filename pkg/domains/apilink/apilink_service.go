package apilink

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linkdebin"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
	uuid "github.com/satori/go.uuid"
)

type AplinkService interface {
	GetDebinesApiLinkService(requerimientoId string, request linkdebin.RequestGetDebinesLink) (response linkdebin.ResponseGetDebinesLink, erro error)
	GetDebinesPendientesApiLinkService(requerimientoId string, cbu string) (response linkdebin.ResponseGetDebinesPendientesLink, erro error)
	GetDebinApiLinkService(requerimientoId string, request linkdebin.RequestGetDebinLink) (response linkdebin.ResponseGetDebinLink, erro error)
	DeleteDebinApiLinkService(requerimientoId string, request linkdebin.RequestDeleteDebinLink) (response bool, erro error)

	/*
		Este servicio fue creado para realizar transferencias de la cuenta de telco para sus clientes.
	*/
	CreateTransferenciaApiLinkService(requerimientoId, token string, request linktransferencia.RequestTransferenciaCreateLink) (response linktransferencia.ResponseTransferenciaCreateLink, erro error)
	GetTransferenciasApiLinkService(requerimientoId string, request linktransferencia.RequestGetTransferenciasLink) (response linktransferencia.ResponseGetTransferenciasLink, erro error)
	GetTransferenciaApiLinkService(requerimientoId string, request linktransferencia.RequestGetTransferenciaLink) (response linktransferencia.ResponseGetTransferenciaLink, erro error)
	/*
		Genera un UUid
	*/
	GenerarUUid() string
	EliminarPagosRepetidos(pagosPendientesDebin []entities.Pago) (pagosDistintos []entities.Pago)
	// Elimina los pagosintentos "erroneos" de un pago (este se verifica por el campo "paid_at" en null)
	EliminarPagoIntentosErroneos(pagoDebin *entities.Pago) error

	GetTokenApiLinkService(identificador string, scope []linkdtos.EnumScopeLink) (linkdtos.TokenLink, error)
}

// apilink variable que va a manejar la instancia del servicio
var apilink *aplinkService

type aplinkService struct {
	remoteRepository RemoteRepository
	repository       ApilinkRepository
}

func NewService(rm RemoteRepository, r ApilinkRepository) AplinkService {
	// al instanciar el servicio lo almaceno en la variable apilink
	apilink = &aplinkService{
		remoteRepository: rm,
		repository:       r,
	}
	return apilink
}

// Resolve devuelve la instancia antes creada
func Resolve() *aplinkService {
	return apilink
}

func (s *aplinkService) GenerarUUid() string {
	return uuid.NewV4().String()
}
func (s *aplinkService) EliminarPagoIntentosErroneos(pagoDebin *entities.Pago) error {
	for i := 0; i < len(pagoDebin.PagoIntentos); {
		if pagoDebin.PagoIntentos[i].PaidAt.IsZero() {
			pagoDebin.PagoIntentos = append(pagoDebin.PagoIntentos[:i], pagoDebin.PagoIntentos[i+1:]...)
		} else {
			i++
		}
	}
	return nil
}
func (s *aplinkService) EliminarPagosRepetidos(pagosPendientesDebin []entities.Pago) (pagosDistintos []entities.Pago) {
	pagosIndex := make(map[int64]int) // Mapa para rastrear el índice de cada pago por su ID

	for i, pago := range pagosPendientesDebin {
		pagosIndex[int64(pago.ID)] = i // Guardar el índice del pago por su ID
	}

	// Crear un nuevo slice para almacenar los pagos distintos
	pagosDistintos = make([]entities.Pago, 0, len(pagosIndex))

	for _, index := range pagosIndex {
		pagosDistintos = append(pagosDistintos, pagosPendientesDebin[index])
	}

	return pagosDistintos
}

func (s *aplinkService) GetTokenApiLinkService(identificador string, scope []linkdtos.EnumScopeLink) (linkdtos.TokenLink, error) {

	var response linkdtos.TokenLink
	var err error
	scopes := []linkdtos.EnumScopeLink{linkdtos.TransferenciasBancariasInmediatas}

	response, err = s.remoteRepository.GetTokenApiLink(identificador, scopes)

	if err != nil {
		return linkdtos.TokenLink{}, err
	}
	return response, nil
}
