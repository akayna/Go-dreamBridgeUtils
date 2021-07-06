package jsonfile

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// ReadJSONFile2 - Le um arquivo Json em um diretório específico
func ReadJSONFile2(caminho string, nome string, estrutura interface{}) error {
	return ReadJSONFile(caminho+nome, estrutura)
}

// ReadJSONFile - Le um arquivo Json
func ReadJSONFile(nome string, estrutura interface{}) error {

	// Open our jsonFile
	jsonFile, err := os.Open(nome)

	if err != nil {
		log.Println("jsonfile - ReadJsonFile - Erro ao ler o arquivo " + nome)
		log.Println("Erro: ", err)

		return err
	}

	defer jsonFile.Close()

	//log.Println("jsonfile - ReadJsonFile - Arquivo aberto com sucesso.")

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Println("jsonfile - ReadJsonFile - Erro ao abrir byte array.")
		log.Println("Erro: ", err)
		return err
	}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, estrutura)

	if err != nil {
		log.Println("jsonfile - ReadJsonFile - Erro ao converter json na estrutura.")
		log.Println("Erro: ", err)
		return err
	}

	return nil
}
