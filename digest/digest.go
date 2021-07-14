package digest

import (
	"crypto/hmac"
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

// GenerateSignature - Generate the CyberSource signature using the key and data
func GenerateSignature(key, data string) (string, error) {

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Println("digest - GenerateSignature: Error decoding key from Base64.")
		log.Println("error:", err)
		return "", err
	}

	hasher := hmac.New(sha256.New, decodedKey)
	hasher.Write([]byte(data))

	signatureEncodedSTD := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	return signatureEncodedSTD, nil
}
