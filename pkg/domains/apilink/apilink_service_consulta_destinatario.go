package apilink

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/linkdtos/linkconsultadestinatario"
)

func (s *aplinkService) GetConsultaDestinatarioService(requerimientoId string, request linkconsultadestinatario.RequestConsultaDestinatarioLink) (response linkconsultadestinatario.ResponseConsultaDestinatarioLink, erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	scopes := []linkdtos.EnumScopeLink{linkdtos.ConsultaDestinatario}

	token, erro := s.remoteRepository.GetTokenApiLink(requerimientoId, scopes)

	if erro != nil {
		return
	}

	return s.remoteRepository.GetConsultaDestinatario(requerimientoId, request, token.AccessToken)

}
