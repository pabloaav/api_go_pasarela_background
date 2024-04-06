package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"

type ResponseCuenta struct {
	Id                   uint   `json:"id"`
	ClientesID           int64  `json:"clientes_id"`
	RubrosID             uint   `json:"rubros_id"`
	Cuenta               string `json:"cuenta"`
	Cbu                  string `json:"cbu"`
	Cvu                  string `json:"cvu"`
	Apikey               string `json:"apikey"`
	DiasRetiroAutomatico int64  `json:"dias_retiro_automatico"`
}

func (r *ResponseCuenta) FromCuenta(c entities.Cuenta) {
	r.Id = c.ID
	r.ClientesID = c.ClientesID
	r.RubrosID = c.RubrosID
	r.Cuenta = c.Cuenta
	r.Cbu = c.Cbu
	r.Cvu = c.Cvu
	r.Apikey = c.Apikey
	r.DiasRetiroAutomatico = c.DiasRetiroAutomatico
}
