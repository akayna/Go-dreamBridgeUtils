package email

// Email - Estrutura para o envio de emails
type Email struct {
	DE          string   `json:"de"`
	Para        []string `json:"para"`
	Copia       []string `json:"cc"`
	CopiaOculta []string `json:"cco"`
	Assunto     string   `json:"assunto"`
	Corpo       string   `json:"corpo"`
}
