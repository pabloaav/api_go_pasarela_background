package mockservice

import (
	prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/cierrelotedtos"
	"github.com/stretchr/testify/mock"
)

type MockCierreLoteService struct {
	mock.Mock
}

func (mock *MockCierreLoteService) ArchivoLoteExterno() (totalArchivos int, err error) {
	args := mock.Called()
	resultado := args.Int(0)
	return resultado, args.Error(1)
}
func (mock *MockCierreLoteService) LeerCierreLote() (listaArchivo []prismaCierreLote.PrismaLogArchivoResponse, err error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]prismaCierreLote.PrismaLogArchivoResponse), args.Error(1)
}

