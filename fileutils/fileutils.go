package fileutils

import (
	"bufio"
	"io"
	"log"
	"os"
)

// OpenFile - opens a file and returns the reference to it
func OpenFile(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Println("fileutils.openCSVFile - Erro ao ler arquivo: " + fileName)
		return nil, err
	}

	return file, err
}

// ReadLine - Return one line from a file
func ReadLine(file *os.File) string {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return scanner.Text()
}

// CreateNewFile - Create a new file and return the reference to it
func CreateNewFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)

	if err != nil {
		log.Println("fileutils.CreateNewFile - Error creating file: " + fileName)
		return nil, err
	}

	return file, nil
}

// WriteLineToFile - Writes a line to a file and keep it openned
func WriteLineToFile(file *os.File, line string) (int, error) {

	datawriter := bufio.NewWriter(file)
	bytes, err := datawriter.WriteString(line + "\n")

	if err != nil {
		log.Println("fileutils.WriteLineToFile - Error writting file.")
		return 0, err
	}

	datawriter.Flush()

	return bytes, err
}

// WriteLinesToFile - Write an array of line to a file and keep it openned
func WriteLinesToFile(file *os.File, lines *[]string) (int, error) {

	datawriter := bufio.NewWriter(file)

	totalBytes := 0

	for _, each_ln := range *lines {
		bytes, err := datawriter.WriteString(each_ln + "\n")

		if err != nil {
			log.Println("fileutils.WriteLineToFile - Error writting file.")
			return 0, err
		}

		totalBytes += bytes + 1
	}

	datawriter.Flush()

	return totalBytes, nil
}

// LineCounter - Counts the total line numbers of a file and resets the file pointer
func LineCounter(file *os.File) (int, error) {

	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}

	log.Println("fileutils.LineCounter - Total number of lines: ", lineCount)

	file.Seek(0, io.SeekStart)

	return lineCount, nil
	/*

		The method bellow is faster but don´t count the last line withous a "\n"

		r := bufio.NewReader(file)

		buf := make([]byte, 32*1024)
		count := 0
		lineSep := []byte{'\n'}

		for {
			c, err := r.Read(buf)
			count += bytes.Count(buf[:c], lineSep)

			switch {
			case err == io.EOF:
				return count, nil

			case err != nil:
				return count, err
			}
		}
	*/
}
