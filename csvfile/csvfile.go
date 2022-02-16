package csvfile

import (
	"Go-dreamBridgeUtils/fileutils"
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
)

// ReadCSVHeaderFile - Reade a CSV file and returns an array with the header and another one with the file linas and columns
func ReadCSVHeaderFile(file string, headerLines int) ([][]string, [][]string, error) {
	csvFile, err := fileutils.OpenFile(file)

	if err != nil {
		log.Println("csvfile.ReadCSVHeaderFile - Error reading the file: " + file)
		return nil, nil, err
	}

	defer csvFile.Close()

	// Verifies if the number headers line is ok with the total file lines number
	numLines, err := fileutils.LineCounter(csvFile)

	if err != nil {
		log.Println("csvfile.ReadCSVHeaderFile - Error counting total line number from file.")
		return nil, nil, err
	}

	if numLines <= headerLines {
		log.Println("csvfile.ReadCSVHeaderFile - The headers lines is grater then the total file lines.")
		log.Println("Headers lines: ", headerLines)
		log.Println("Total file lines: ", numLines)
		return nil, nil, errors.New("csvfile.ReadCSVHeaderFile - The headers lines is grater then the total file lines")
	}

	// Reads the header and data from the file.
	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	csvReader.FieldsPerRecord = -1

	var headerCSV [][]string
	var dataCSV [][]string

	var counter = 0

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("csvfile.ReadCSVHeaderFile - Error reading file.")
			return nil, nil, err
		}

		if counter >= headerLines {
			dataCSV = append(dataCSV, line)
		} else {
			headerCSV = append(headerCSV, line)
			counter++
		}
	}

	return headerCSV, dataCSV, nil
}

// WriteCSV - Write a CSV file using "," as separator
func WriteCSV(header, data *[][]string, fileName string) error {
	file, err := os.Create(fileName)

	if err != nil {
		log.Println("csvfile.WriteCSV - Erro ao criar arquivo: " + fileName)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	for _, line := range *header {
		err := writer.Write(line)

		if err != nil {
			log.Println("csvfile.WriteCSV - Error writing file: " + fileName)
			return err
		}
	}

	for _, linha := range *data {
		err := writer.Write(linha)

		if err != nil {
			log.Println("csvfile.EscreveArquivoCSV - Error writing file: " + fileName)
			return err
		}
	}

	return nil
}
