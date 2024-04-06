package retenciondtos

import (
	"fmt"
)

type CertificadoFileDTO struct {
	Id                 uint   ` json:"id"`
	FileName           string `json:"file_name"`
	Content            string `json:"content"`
	ClienteRetencionId uint   `json:"cliente_retencion_id"`
}
type ComprobanteFileDTO struct {
	Id        uint   ` json:"id"`
	ClienteId uint   `json:"cliente_id"`
	FileName  string `json:"file_name"`
	Content   string `json:"content"`
}

type RetencionCertificadoRequestDTO struct {
	Id                 uint   `json:"id"`
	RetencionId        uint   `json:"retencion_id"`
	ClienteId          uint   `json:"cliente_id"`
	ClienteRetencionId uint   `json:"cxr_id"`
	Fecha_Presentacion string `json:"fecha_presentacion"`
	Fecha_Caducidad    string `json:"fecha_caducidad"`
	RutaFile           string `json:"ruta_file"`
}

func (rcdto *RetencionCertificadoRequestDTO) Validar() error {

	if rcdto.RetencionId <= 0 || rcdto.ClienteId <= 0 {
		return fmt.Errorf(ERROR_CAMPO, "RetencionId")
	}

	return nil
}
