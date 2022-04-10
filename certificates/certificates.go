package certificates

import (
	"encoding/pem"
	"log"

	"golang.org/x/crypto/pkcs12"
)

// Pkcs12ToPemBytes - Converts a pkcs file data to PEM
func Pkcs12ToPemBytes(pfxData []byte, password string) ([][]byte, error) {

	pems, err := pkcs12.ToPEM(pfxData, password)
	if err != nil {
		log.Println("certificates.Pkcs12ToPemBytes - Error converting pfx data to PEM.")
		return nil, err
	}

	var pubsBytes [][]byte

	for _, pemBytes := range pems {
		pubBytes := pem.EncodeToMemory(pemBytes)

		pubsBytes = append(pubsBytes, pubBytes)
	}

	return pubsBytes, nil
}
