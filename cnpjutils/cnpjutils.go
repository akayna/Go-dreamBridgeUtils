package cnpjutils

import (
	"errors"
	"log"
	"regexp"
	"strconv"
)

// NumDigitosCNPJ - Quantidade de digitos de um CNPJ
const NumDigitosCNPJ = 14

// CNPJ - Estrutura para armazenar um CNPJ
type CNPJ struct {
	CnpjDigito [NumDigitosCNPJ]uint
}

// FormataCNPJInt - Retorna o cnpj formatado
func FormataCNPJInt(cnpj int64) (string, error) {
	cpfStr, err := Int64ToCNPJ(cnpj)

	if err != nil {
		log.Println("cnpjutils.FormataCNPJInt - Erro ao converte cnpj para estrutura CNPJ.")
		return "", err
	}

	return CNPJtoStringFormatada(cpfStr), nil
}

// CNPJToString - Converte o CPF numérico em string
func CNPJToString(cnpj *CNPJ) string {
	var cnpjString string
	var flagZeroEsquerda = true

	cnpjString = ""
	for i := 0; i < NumDigitosCNPJ; i++ {
		valDigito := int(cnpj.CnpjDigito[i])

		if flagZeroEsquerda {
			if valDigito != 0 {
				flagZeroEsquerda = false
				cnpjString += strconv.Itoa(valDigito)
			}
		} else {
			cnpjString += strconv.Itoa(valDigito)
		}
	}

	return cnpjString
}

// CNPJtoStringFormatada - Converte o CPF numérico em string formatada XX.XXX.XXX/XXXX-XX
func CNPJtoStringFormatada(cnpj *CNPJ) string {
	var cnpjString string
	var flagZeroEsquerda = true

	cnpjString = ""
	for i := 0; i < NumDigitosCNPJ; i++ {
		valDigito := int(cnpj.CnpjDigito[i])

		if flagZeroEsquerda {
			if valDigito != 0 {
				flagZeroEsquerda = false
				cnpjString += strconv.Itoa(valDigito)
			}
		} else {
			cnpjString += strconv.Itoa(valDigito)
		}

		if i == 1 || i == 4 {
			cnpjString += "."
		}

		if i == 7 {
			cnpjString += "/"
		}

		if i == 11 {
			cnpjString += "-"
		}
	}

	return cnpjString
}

// Int64ToCNPJ - Recebe um CNPJ em forma de inteiro e transforma em números
func Int64ToCNPJ(cnpj int64) (*CNPJ, error) {
	cnpjString := strconv.FormatInt(cnpj, 10)

	return StringToCNPJ(cnpjString)
}

// StringToCNPJ - Recebe um CNPJ, transforma em números
func StringToCNPJ(cnpjStr string) (*CNPJ, error) {
	var cnpj CNPJ

	// Elimina todos os caracteres diferentes de números da string
	cnpjStrNum, err := LimpaStringCNPJ(cnpjStr)
	if err != nil {
		log.Println("cnpjutils.StringToCNPJ - Erro ao limpar string de cnpj.")
		return nil, err
	}

	if len(cnpjStrNum) > NumDigitosCNPJ {
		return nil, errors.New("cnpjutils.StringToCNPJ - Número de dígitos maior do que o máximo permitido para um CNPJ")
	}

	for index, cnpjRune := range cnpjStrNum {
		cnpj.CnpjDigito[index] = uint(cnpjRune - '0')
	}

	return &cnpj, nil
}

// LimpaStringCNPJ - Retorna uma string contendo apenas os números do CNPJ
func LimpaStringCNPJ(cnpj string) (string, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Println("cnpjutils.LimpaStringCNPJ - Erro ao limpar string de cnpj.")
		return "", err
	}
	apenasNumeros := reg.ReplaceAllString(cnpj, "")

	return apenasNumeros, nil
}
