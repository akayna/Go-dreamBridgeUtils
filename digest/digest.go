package digest

import (
	"crypto/sha256"
	"encoding/base64"
	"log"

	"golang.org/x/text/encoding/unicode"
)

// GenerateDigestSHA256 - Return the digest SHA-256 hashing in base64
func GenerateDigestSHA256(data string) (string, error) {

	decoder := unicode.UTF8.NewDecoder()

	dataUTF8, err := decoder.Bytes([]byte(data))

	if err != nil {
		log.Println("digest.GenerateDigestSHA256: Error converting string to UTF8 char array.")
		return "", err
	}

	h := sha256.New()
	h.Write(dataUTF8)

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
