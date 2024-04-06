package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/entities"

type RequestAuditoria struct {
	UserID    uint
	CuentaID  uint
	IP        string
	Tabla     string
	Fila      uint
	Operacion string
	Query     string
	Resultado string
	Origen    string
}

func (r *RequestAuditoria) ToAuditoria() (response entities.Auditoria) {

	response.UserID = r.UserID
	response.CuentaID = r.CuentaID
	response.IP = r.IP
	response.Tabla = r.Tabla
	response.Fila = r.Fila
	response.Operacion = "insert"
	response.Query = r.Query
	response.Resultado = r.Resultado
	response.Origen = r.Origen

	return
}
