package retenciondtos

import "time"

type CertificadoVencimientoDTO struct {
	Ids               []uint    ` json:"ids"`
	Fecha_comparacion time.Time ` json:"fecha_comparacion"`
	Administracion    bool      ` json:"administracion"`
}
