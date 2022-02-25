package emailutils

import (
	"errors"
	"log"
	"net/smtp"

	"github.com/akayna/Go-dreamBridgeUtils/stringutils"
)

func (email *TextEmail) montaMensagemEmailTexto() ([]byte, error) {
	var messageStr string

	if len(email.From) > 0 {
		messageStr += "From: " + email.From + "\r\n"
	}

	if len(email.To) <= 0 {
		mensagem := "email.Email.montaMensagemEmailTexto - Falta email de destino"
		log.Println(mensagem)
		return nil, errors.New(mensagem)
	}

	messageStr += "TO: " + stringutils.VectorStringToStringLine(email.To) + "\r\n"

	if len(email.Co) > 0 {
		messageStr += "CO: " + stringutils.VectorStringToStringLine(email.Co) + "\r\n"
	}

	if len(email.Cco) > 0 {
		messageStr += "CCO: "

		messageStr += "CCO: " + stringutils.VectorStringToStringLine(email.Cco) + "\r\n"
	}

	messageStr += "Subject: " + email.Subject + "\r\n\r\n"

	messageStr += email.Body

	//log.Println(messageStr)

	return []byte(messageStr), nil
}

// EnviaEmailSMTP - Envia um email de texto simples
func (email *TextEmail) EnviaEmailSMTP() error {

	messageStr, err := email.montaMensagemEmailTexto()

	if err != nil {
		log.Println("email.Email.EnviaEmail - falha ao montar mensagem.")
		return err
	}

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", email.From, email.Password, "smtp.gmail.com"),
		email.From, email.To, messageStr)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	return nil
}
