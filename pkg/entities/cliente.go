package entities

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Cliente struct {
	gorm.Model
	IvaID             int64               `json:"iva_id"`
	IibbID            int64               `json:"iibb_id"`
	Cliente           string              `json:"cliente"`
	Razonsocial       string              `json:"razonsocial"`
	Nombrefantasia    string              `json:"nombrefantasia"`
	Email             string              `json:"email"`
	Emailcontacto     string              `json:"emailcontacto"`
	Personeria        string              `json:"personeria"`
	RetiroAutomatico  bool                `json:"retiro_automatico"`
	Cuit              string              `json:"cuit"`
	ReporteBatch      bool                `json:"reporte_batch"`
	NombreReporte     string              `json:"nombre_reporte"`
	OrdenDiaria       bool                `json:"orden_diaria"`
	SplitCuentas      bool                `json:"split_cuentas"`
	Iva               *Impuesto           `json:"iva" gorm:"foreignKey:iva_id"`
	Iibb              *Impuesto           `json:"iibb" gorm:"foreignKey:iibb_id"`
	Cuentas           *[]Cuenta           `json:"cuentas" gorm:"foreignKey:ClientesID"`
	Clienteusers      *[]Clienteuser      `json:"cliente_users" gorm:"foreignKey:ClientesId"`
	Contactosreportes *[]Contactosreporte `json:"contactos_reportes" gorm:"foreignKey:ClientesID"`
	Retenciones       []Retencion         `gorm:"many2many:cliente_retencions;"`
	ClienteRetencions []ClienteRetencion  `json:"cliente_retencions" gorm:"foreignKey:cliente_id"`
	Comprobantes      []Comprobante       `gorm:"foreignkey:ClienteId"`
	Domicilio         string              `json:"domicilio"`
	SujetoRetencion   bool                `json:"sujeto_retencion"`
	Formulario8125    bool                `json:"formulario_8125"`
}

func (ct *Cliente) AfterSave(tx *gorm.DB) (err error) {
	var audit Auditoria
	ctxValue := tx.Statement.Context.Value(AuditUserKey{})
	if ctxValue == nil {
		return errors.New("no hay datos de usuario para la transacción indicada")
	}
	audit = ctxValue.(Auditoria)
	stmt := tx.Statement
	str := tx.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	audit.Fila = ct.ID
	audit.Query = str
	audit.Tabla = "cliente"
	newCtx := context.WithValue(tx.Statement.Context, AuditUserKey{}, audit)
	tx.Statement.Context = newCtx

	return nil
}

// dado el name de un gravamen, retorna si el cliente tiene asociadas retenciones de ese gravamen
func (ct *Cliente) HasRetencionByGravamenName(gravamenes []string) (res bool) {
	if len(ct.Retenciones) < 1 {
		return
	}
	for _, g := range gravamenes {
		for _, retencion := range ct.Retenciones {
			if retencion.Condicion.Gravamen.Gravamen == g {
				res = true
			}
		}
	}
	return
}
