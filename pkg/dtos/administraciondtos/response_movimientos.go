package administraciondtos

type MovimientosSubcuentas struct {
	Id            uint   `json:"id"`
	SubcuentasID  uint64 `json:"subcuentas_id"`
	MovimientosID uint64 `json:"movimientos_id"`
	Transferido   bool   `json:"transferido"`
	Monto         int64  `json:"monto"`
}
