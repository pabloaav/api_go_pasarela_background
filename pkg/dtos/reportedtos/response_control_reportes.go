package reportedtos

type ResponseControlReporte struct {
	Estado               string `json:"estado"`
	CuentaId             string `json:"cuenta_id"`
	Cuenta               string `json:"cuenta"`
	MontoCobranza        string `json:"monto_cobranza"`
	MontoCobranzaCliente string `json:"monto_cobranza_cliente"`
	Diferencia           string `json:"diferencia"`
}

type ResponseData struct {
	Data    []ResponseControlReporte `json:"data"`
	Message string                   `json:"message"`
	Status  bool                     `json:"status"`
}
