package webhook

type RequestWebhook struct {
	DiasPago         int64  `json:"dias_pago"`
	PagosNotificado  bool   `json:"pagos_notificado"`
	EstadoFinalPagos bool   `json:"estado_final_pagos"`
	CuentaId         uint64 `json:"cuenta_id"`
}

type RequestWebhookReferences struct {
	CuentaId           uint64   `query:"cuenta_id"`
	ExternalReferences []string `query:"external_references"`
	EstadoFinalPagos   bool     `query:"estado_final"`
	PagosId            []int    `query:"pagos_id"`
}
