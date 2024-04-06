package reportedtos

type RequestRetencionEmail struct {
	Archivos []string `json:"archivos"`
	Emails   []string `json:"emails"`
}
