package stringutils

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// CriptografaString - Criptografa uma senha
func CriptografaString(senha string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), cost)

	fmt.Println("Hash: " + string(bytes))
	return string(bytes), err
}

// checkStringHash - Verifica se a criptografia corresponde à string passada
func checkStringHash(str, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))

	if err != nil {
		log.Println("checkStringHash - Erro ao validar string.")
		return false, err
	}

	return true, nil
}

// SeparaNomeSobrenome - Separa uma string de nome completo em nome e sobrenome. OBS: Não funciona para nome composto.
func SeparaNomeSobrenome(nomeCompleto string) (string, string, error) {
	var nome string
	var sobrenome string

	i := strings.Index(nomeCompleto, " ")

	if i < 0 {
		log.Println("stringutils.SeparaNomeSobrenome - O nome não contém sobrenome.")
		nome = nomeCompleto
		return nome, sobrenome, errors.New("stringutils.SeparaNomeSobrenome - O nome não contém sobrenome")
	} else if i == 0 {
		log.Println("stringutils.SeparaNomeSobrenome - Existe um espeço na primeira letra? A string está vazia?")
		nome = nomeCompleto
		return nome, sobrenome, errors.New("stringutils.SeparaNomeSobrenome - Existe um espeço na primeira letra? A string está vazia?")
	} else {
		nome = nomeCompleto[0:i]
		sobrenome = nomeCompleto[(i + 1):len(nomeCompleto)]
	}

	return nome, sobrenome, nil
}