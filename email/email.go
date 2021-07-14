package email

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"

	gmail "google.golang.org/api/gmail/v1"

	gmailDB "github.com/rafaelcunha/Go-dreambridgeUtils/email/gmail"
)

func (email *Email) montaMensagemEmailTexto() ([]byte, error) {
	var messageStr string

	if len(email.DE) > 0 {
		messageStr += "From: " + email.DE + "\r\n"
	}

	if len(email.Para) <= 0 {
		mensagem := "email.Email.montaMensagemEmailTexto - Falta email de destino."
		log.Println(mensagem)
		return nil, errors.New(mensagem)
	}

	messageStr += "TO: "

	for index, value := range email.Para {
		messageStr += value
		if index < (len(email.Para) - 1) {
			messageStr += ","
		}
	}

	messageStr += "\r\n"

	if len(email.Copia) > 0 {
		messageStr += "CO: "

		for index, value := range email.Copia {
			messageStr += value
			if index < (len(email.Copia) - 1) {
				messageStr += ","
			}
		}

		messageStr += "\r\n"
	}

	if len(email.CopiaOculta) > 0 {
		messageStr += "CCO: "

		for index, value := range email.CopiaOculta {
			messageStr += value
			if index < (len(email.CopiaOculta) - 1) {
				messageStr += ","
			}
		}

		messageStr += "\r\n"
	}

	messageStr += "Subject: " + email.Assunto + "\r\n\r\n"

	messageStr += email.Corpo

	log.Println(messageStr)

	return []byte(messageStr), nil
}

// EnviaEmail - Envia um email de texto simples
func (email *Email) EnviaEmail() error {

	messageStr, err := email.montaMensagemEmailTexto()

	if err != nil {
		log.Println("email.Email.EnviaEmail - falha ao montar mensagem.")
		return err
	}

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println("email.Email.EnviaEmail - falha ao ler arquivo credentials.json.")
		return err
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Println("email.Email.EnviaEmail - Unable to parse client secret file to config.")
		return err
	}

	srv, err := gmail.New(gmailDB.GetClient(config))
	if err != nil {
		log.Println("email.Email.EnviaEmail - Unable to retrieve Gmail client.")
		return err
	}

	var message gmail.Message

	// Place messageStr into message.Raw in base64 encoded format
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	// Send the message
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Println("email.Email.EnviaEmail - Erro ao enviar email.")
		return err
	}

	return nil
}
