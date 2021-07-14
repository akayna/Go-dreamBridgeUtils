package gmailutils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

// EmailTexto - estrutura para envio de emails de texto simples
type EmailTexto struct {
	From    string
	TO      []string
	CCO     []string
	Subject string
	Body    string
}

func montaMensagemEmailTexto(email *EmailTexto) ([]byte, error) {

	if len(email.TO) <= 0 {
		return nil, errors.New("Falta email de destino")
	}

	var messageStr = "From: " + email.From + "\r\n"

	messageStr += "TO: "

	for index, value := range email.TO {
		messageStr += value
		if index < (len(email.TO) - 1) {
			messageStr += ","
		}
	}

	messageStr += "\r\n"

	if len(email.CCO) > 0 {
		messageStr += "CCO: "

		for index, value := range email.CCO {
			messageStr += value
			if index < (len(email.CCO) - 1) {
				messageStr += ","
			}
		}
	
		messageStr += "\r\n"
	}

	messageStr += "Subject: " + email.Subject + "\r\n\r\n"

	messageStr += email.Body

	return []byte(messageStr), nil
}

// EnviaEmailHTML - Envia um email de texto simples
/*func EnviaEmailHTML(email EmailTexto) {

	messageStr, err := montaMensagemEmailTexto(email)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	srv, err := gmail.New(getClient(config))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	var message gmail.Message

	// Place messageStr into message.Raw in base64 encoded format
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	// Send the message
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}
}*/

// EnviaEmailTexto - Envia um email de texto simples
func EnviaEmailTexto(email *EmailTexto) error {

	messageStr, err := montaMensagemEmailTexto(email)

	if err != nil {
		log.Println("gmailutils.EnviaEmailTexto - falha ao montar mensagem.")
		return err
	}

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println("gmailutils.EnviaEmailTexto - falha ao ler arquivo dcredentials.json.")
		return err
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Println("gmailutils.EnviaEmailTexto - Unable to parse client secret file to config.")
		return err
	}

	srv, err := gmail.New(getClient(config))
	if err != nil {
		log.Println("gmailutils.EnviaEmailTexto - Unable to retrieve Gmail client.")
		return err
	}

	var message gmail.Message

	// Place messageStr into message.Raw in base64 encoded format
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	// Send the message
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Println("gmailutils.EnviaEmailTexto - Erro ao enviar email.")
		return err
	}

	return nil
}

// GetClient - Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}
