package emailutils

const emailQueueMaxSize = 200

var textEmailQueue []TextEmail

// Email - Struct to send text email
type TextEmail struct {
	From     string   `json:"from"`
	Password string   `json:"password"`
	To       []string `json:"to"`
	Co       []string `json:"co"`
	Cco      []string `json:"cco"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
}
