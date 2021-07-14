package csvfile

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

// ReadCSVFile - Le um arquivo CSV e o retorna em uma slice de linhas e colunas de string
func ReadCSVFile(arquivo string) ([][]string, error) {
	csvFile, err := os.Open(arquivo)

	if err != nil {
		log.Println("csvfile.ReadCSVFile - Erro ao ler arquivo: " + arquivo)
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	var arquivoCSV [][]string

	for {
		linha, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("csvfile.ReadCSVFile - Erro ao ler arquivo: " + arquivo)
			return nil, err
		}

		arquivoCSV = append(arquivoCSV, linha)
	}

	return arquivoCSV, nil
}

// EscreveArquivoCSV - Escreve um arquivo CSV unsado , como separador
func EscreveArquivoCSV(dados *[][]string, arquivo string) error {
	file, err := os.Create(arquivo)

	if err != nil {
		log.Println("csvfile.EscreveArquivoCSV - Erro ao criar arquivo: " + arquivo)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, linha := range *dados {
		err := writer.Write(linha)

		if err != nil {
			log.Println("csvfile.EscreveArquivoCSV - Erro ao escrever no arquivo: " + arquivo)
			return err
		}
	}

	return nil
}
