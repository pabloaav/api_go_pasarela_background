package util_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	util "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/test/domains/util/utilfake"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/test/mocks/mockrepository"
	"github.com/stretchr/testify/assert"
)

var (
	mockutils = new(mockrepository.MockRepositoryUtil)
	// mockFactory = new(mockservice.MockCrearMensajeServiceFactory)
	//mockServiceAdministracion  = new(mockservice.)
	service = util.NewUtilService(mockutils)
)

func TestRequestValidConsultarMovimientos(t *testing.T) {
	TableDriverTest := utilfake.EstructuraVerificarCbu()
	t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
		want := TableDriverTest.WantTable
		logs.Info(want)
		// _, got := service.ConsultarMovimientos(TableDriverTest.Request)
		// assert.Equal(t, got.Error(), want)
	})
}

// test para enviar email

func TestBuildEmailSend(t *testing.T) {
	TableDriverTest := utilfake.EstructuraValidarCbu()
	t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
		want := TableDriverTest.WantTable
		logs.Info(want)
		got, _ := service.ValidarCBU(TableDriverTest.Cbu)
		assert.Equal(t, got, want)
	})
}

// test construir movimientos , caluclo de comisiones e impuestos
func TestBuildComisiones(t *testing.T) {
	TableDriverTest := utilfake.EstructuraBuildComisiones()
	for _, test := range TableDriverTest {
		t.Run(test.TituloPrueba, func(t *testing.T) {
			want := test.WantTable
			logs.Info(test.TituloPrueba)
			got := service.BuildComisiones(test.RequestMovimiento, test.RequestCuentaComision, test.RequestIva, test.ImporteSolicitado)
			assert.Equal(t, got, want)
		})
	}
}

// test para calcular movimientos de subcuentas
func TestBuildMovimientosSubcuentas(t *testing.T) {
	TableDriverTest := utilfake.EstructuraBuildMovimientosSubcuentas() // crear estructura para orueba
	for _, test := range TableDriverTest {
		t.Run(test.TituloPrueba, func(t *testing.T) {
			want := test.WantTable
			logs.Info(test.TituloPrueba)
			got := service.BuildMovimientoSubcuentas(test.RequestMovimiento, &test.RequestCuenta)
			assert.Equal(t, got, want)
		})
	}
}

func TestFormatNum(t *testing.T) {
	TableDriverTest := utilfake.EstructuraFormatNum()
	for _, test := range TableDriverTest.Importe {
		t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
			want := TableDriverTest.WantTable
			logs.Info(TableDriverTest.TituloPrueba)
			logs.Info(test)
			got := service.FormatNum(test)
			logs.Info(got)
			assert.Equal(t, got, want)
		})
	}
}
