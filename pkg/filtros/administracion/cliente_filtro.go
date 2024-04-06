package filtros

type ClienteFiltro struct {
	Paginacion
	Id                   uint
	DistintoId           uint
	Cuit                 string
	UserId               uint
	RetiroAutomatico     bool
	CargarImpuestos      bool
	CargarCuentas        bool
	CargarRubros         bool
	CargarCuentaComision bool
	CargarTiposPago      bool
	CargarContactos      bool
	ClientesIds          []uint
	OrdenDiaria          bool // FIXME se debe cambiar por split cuentas
	CargarSubcuentas     bool
	SplitCuentas         bool
	Nombre               string
	SujetoRetencion      bool // determina si al cliente le corresponde retenciones
	Formulario8125       bool
}
