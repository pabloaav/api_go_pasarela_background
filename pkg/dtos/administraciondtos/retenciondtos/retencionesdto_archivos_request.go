package retenciondtos

type RentencionArchivoRequest struct {
	CertificadoId uint ` json:"certificado_id"`
	ComprobanteId uint ` json:"comprobante_id"`
}
