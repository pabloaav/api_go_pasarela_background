package utilfake

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"
)

func EstructuraVerificarCbu() (tableDriverTestPeyment TableDriverTestConsultarMoviento) {
	tableDriverTestPeyment = TableDriverTestConsultarMoviento{
		TituloPrueba: "el tipo de movimiento no es valido, los valores correctos son debin, prisma, transferencia",
		WantTable:    true,
		Cbu:          "56477491421212212121212",
	}
	return
}

// // Construir el texto html del mensaje del email
// mensaje := "<p style='box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif,'Apple Color Emoji','Segoe UI Emoji','Segoe UI Symbol';font-size:16px;line-height:1.5em;margin-top:0;text-align:center'><h2 style='text-align:center'>Operación de pago exitosa</h2> El pago de la referencia <b>#4</b> fue aprobado. <ul><li> Importe: <b>#0</b></li><li> Identificador de la transacción: <b>#1</b></li><li> Medio de pago: <b>#2</b></li><li> Concepto: <b>#3</b></li></ul></p>"
// /* enviar mail al usuario pagador */
// var arrayEmail []string
// var email string
// email = request.HolderEmail
// if request.HolderEmail == "" {
// 	email = pago.PayerEmail
// }
// arrayEmail = append(arrayEmail, email)
// params := utildtos.RequestDatosMail{
// 	Email:            arrayEmail,
// 	Asunto:           "Información de Pago",
// 	Nombre:           pago.PayerName,
// 	Mensaje:          mensaje,
// 	CamposReemplazar: []string{fmt.Sprintf("$%v", response.ImportePagado), pago.Uuid, medio.Mediopago, response.Description, response.ExternalReference},
// 	From:             "Wee.ar!",
// 	TipoEmail:        "template",
// }

func EstructuraEmail() (tableDriverTestPeyment TableDriverTestEmailSend) {
	tableDriverTestPeyment = TableDriverTestEmailSend{
		TituloPrueba: "envios de email en pagos exitosos",
		WantTable:    "",
		Request: utildtos.RequestDatosMail{
			Email:            []string{"jose.alarcon@telco.com.ar"},
			Asunto:           "Información de Pago",
			Nombre:           "jose",
			Mensaje:          "<p style='box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif,'Apple Color Emoji','Segoe UI Emoji','Segoe UI Symbol';font-size:16px;line-height:1.5em;margin-top:0;text-align:center'><h2 style='text-align:center'>Operación de pago exitosa</h2> El pago de la referencia <b>#4</b> fue aprobado. <ul><li> Importe: <b>#0</b></li><li> Identificador de la transacción: <b>#1</b></li><li> Medio de pago: <b>#2</b></li><li> Concepto: <b>#3</b></li></ul></p>",
			CamposReemplazar: []string{"djasjds", "dsadsad", "dasdasdasdasd"},
			From:             "Wee.ar",
			TipoEmail:        "template",
		},
	}
	return
}

func EstructuraValidarCbu() (tableDriverTest TableDriverTestConsultarMoviento) {
	tableDriverTest = TableDriverTestConsultarMoviento{
		TituloPrueba: "validar cbu",
		WantTable:    true,
		Cbu:          "0940099372007393130021",
	}
	return
}

func EstructuraFormatNum() (tableDriverTest TableDriverTestMoneda) {
	tableDriverTest = TableDriverTestMoneda{
		TituloPrueba: "Formatear valor",
		WantTable:    true,
		Importe:      []float64{475.10, 500.99, 11900.87, 11900, 34954.39, 10, 5, 89.41, 74.986541, 1000000, 425.01},
	}
	return
}

func EstructuraBuildComisiones() (tableDriverTest []TableDriverBuildComisiones) {
	// var cuentaComision *[]entities.Cuentacomision
	var tableDriver []TableDriverBuildComisiones
	var err error
	cuentaComisionDebit := append([]entities.Cuentacomision{}, entities.Cuentacomision{
		CuentasID:          5,
		ChannelsId:         3,
		ChannelarancelesId: 15,
		Cuentacomision:     "DGR debito",
		Comision:           0.0025,
		Mediopagoid:        0,
		Importeminimo:      0,
		Importemaximo:      0,
		ChannelArancel: entities.Channelarancele{
			ChannelsId:    2,
			RubrosId:      1,
			Importe:       0.0035,
			Tipocalculo:   "PORCENTAJE",
			Importeminimo: 0,
			Importemaximo: 0,
			Mediopagoid:   0,
			Pagocuota:     false,
		},
	})

	cuentaComisionDebin1 := append([]entities.Cuentacomision{}, entities.Cuentacomision{
		CuentasID:          5,
		ChannelsId:         3,
		ChannelarancelesId: 15,
		Cuentacomision:     "DGR debin",
		Comision:           0.006,
		Mediopagoid:        0,
		Importeminimo:      9,
		Importemaximo:      105,
		ChannelArancel: entities.Channelarancele{
			ChannelsId:    4,
			RubrosId:      1,
			Importe:       5.32,
			Tipocalculo:   "FIJO",
			Importeminimo: 0,
			Importemaximo: 0,
			Mediopagoid:   0,
			Pagocuota:     false,
		},
	})

	cuentaComisionDebin := append([]entities.Cuentacomision{}, entities.Cuentacomision{
		CuentasID:          5,
		ChannelsId:         3,
		ChannelarancelesId: 15,
		Cuentacomision:     "DGR debin",
		Comision:           0.006,
		Mediopagoid:        0,
		Importeminimo:      9,
		Importemaximo:      105,
		ChannelArancel: entities.Channelarancele{
			ChannelsId:    1,
			RubrosId:      4,
			Importe:       5.32,
			Tipocalculo:   "FIJO",
			Importeminimo: 0,
			Importemaximo: 0,
			Mediopagoid:   0,
			Pagocuota:     false,
		},
	})

	tableDriverTest = append(tableDriver, TableDriverBuildComisiones{
		// &PRUEBA 1 - debito
		TituloPrueba: "Calculo de comisioones DEBITO : Telco minimo y Proveedor minimo",
		WantTable:    err,
		RequestMovimiento: &entities.Movimiento{
			CuentasId:      4,
			PagointentosId: 187,
			Tipo:           "C",
			Monto:          571807,
			MotivoBaja:     "",
			Reversion:      false,
			Enobservacion:  false,
		},
		RequestCuentaComision: &cuentaComisionDebit,
		RequestIva: &entities.Impuesto{
			Impuesto:   "IVA",
			Porcentaje: 0.21,
		},
		ImporteSolicitado: 571807},
		// &PRUEBA 2 - DEBIn 1
		TableDriverBuildComisiones{
			TituloPrueba: "Calculo de comisiones DEBIN: Telco minimo y Proveedor sin minimo",
			WantTable:    err,
			RequestMovimiento: &entities.Movimiento{
				CuentasId:      4,
				PagointentosId: 187,
				Tipo:           "C",
				Monto:          1760000,
				MotivoBaja:     "",
				Reversion:      false,
				Enobservacion:  false,
			},
			RequestCuentaComision: &cuentaComisionDebin1,
			RequestIva: &entities.Impuesto{
				Impuesto:   "IVA",
				Porcentaje: 0.21,
			},
			ImporteSolicitado: 1760000},
		// &PRUEBA 3 - DEBIN 2
		TableDriverBuildComisiones{
			TituloPrueba: "Calculo de comisiones DEBIN 1: Telco sin minimo y Proveedor sin minimo",
			WantTable:    err,
			RequestMovimiento: &entities.Movimiento{
				CuentasId:      4,
				PagointentosId: 187,
				Tipo:           "C",
				Monto:          10000,
				MotivoBaja:     "",
				Reversion:      false,
				Enobservacion:  false,
			},
			RequestCuentaComision: &cuentaComisionDebin,
			RequestIva: &entities.Impuesto{
				Impuesto:   "IVA",
				Porcentaje: 0.21,
			},
			ImporteSolicitado: 10000},
	)

	return
}

func EstructuraBuildMovimientosSubcuentas() (tableDriverTest []TableDriverBuildMovimientosSubcuentas) {
	// var cuentaComision *[]entities.Cuentacomision
	var tableDriver []TableDriverBuildMovimientosSubcuentas
	var err error

	var subcuentas []entities.Subcuenta
	subcuenta1 := entities.Subcuenta{
		Tipo:       "principal",
		CuentasID:  1,
		Cbu:        "12123123123123123123",
		Porcentaje: 0.7,
	}
	subcuenta2 := entities.Subcuenta{
		Tipo:       "secundaria",
		CuentasID:  1,
		Cbu:        "12123123123123123123",
		Porcentaje: 0.3,
	}

	subcuentas = append(subcuentas, subcuenta1, subcuenta2)
	tableDriver = append(tableDriver, TableDriverBuildMovimientosSubcuentas{
		// &PRUEBA 1
		TituloPrueba: "calculo movimientos 1",
		WantTable:    err,
		RequestMovimiento: &entities.Movimiento{
			CuentasId:      4,
			PagointentosId: 187,
			Tipo:           "C",
			Monto:          315204,
			MotivoBaja:     "",
			Reversion:      false,
			Enobservacion:  false,
		},
		RequestCuenta: entities.Cuenta{
			ClientesID: 1,
			RubrosID:   1,
			Cuenta:     "Cuenta prueba",
			Cbu:        "324234234234234234",
			Subcuentas: subcuentas,
		}},
		TableDriverBuildMovimientosSubcuentas{
			// &PRUEBA 2
			TituloPrueba: "calculo movimientos 2",
			WantTable:    err,
			RequestMovimiento: &entities.Movimiento{
				CuentasId:      4,
				PagointentosId: 187,
				Tipo:           "C",
				Monto:          3134786,
				MotivoBaja:     "",
				Reversion:      false,
				Enobservacion:  false,
			},
			RequestCuenta: entities.Cuenta{
				ClientesID: 1,
				RubrosID:   1,
				Cuenta:     "Cuenta prueba",
				Cbu:        "324234234234234234",
				Subcuentas: subcuentas,
			}},
	)

	tableDriverTest = tableDriver

	return
}
