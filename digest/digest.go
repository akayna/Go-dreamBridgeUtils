package digest

import (
	"crypto/sha256"
	"encoding/base64"
	"log"

	"golang.org/x/text/encoding/unicode"
)

// GenerateDigest - Return the digest SHA-256 hashing
func GenerateDigest(data string) (string, error) {

	decoder := unicode.UTF8.NewDecoder()

	dataUTF8, err := decoder.Bytes([]byte(data))

	if err != nil {
		log.Println("digest - GenerateDigest: Error converting string to UTF8 char array.")
		log.Println(err)
		return "", err
	}

	h := sha256.New()
	h.Write(dataUTF8)

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
