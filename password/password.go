package password

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

// Retirado de: https://github.com/cmiceli/password-generator-go/blob/master/gen.go

//var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")
// StdChars - Caracteres utilizado para gerar o password
var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// NewPassword - Gera o password
func NewPassword(length int) string {
	return randChar(length, StdChars)
}

func randChar(length int, chars []byte) string {
	newPword := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
			panic(err)
		}
		for _, c := range randomData {
			if c >= maxrb {
				continue
			}
			newPword[i] = chars[c%clen]
			i++
			if i == length {
				return string(newPword)
			}
		}
	}
}

// Retirado de: https://www.systutorials.com/how-to-generate-rsa-private-and-public-key-pair-in-go-lang/

// GeneratePrivateKey - Generate one provate key with bits length
func GeneratePrivateKey(bits int) (*rsa.PrivateKey, error) {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Printf("Cannot generate RSA key\n")
		os.Exit(1)
	}

	return privatekey, err
}

// GetPublicKey - Get the public part of a private key
func GetPublicKey(privateKey *rsa.PrivateKey) *rsa.PublicKey {
	return &privateKey.PublicKey
}

// DumpPrivateKeyToFile - Save one private key into a file named private.pem
func DumpPrivateKeyToFile(privateKey *rsa.PrivateKey) error {
	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create("private.pem")
	if err != nil {
		fmt.Printf("error when create private.pem: %s \n", err)
		return err
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		fmt.Printf("error when encode private pem: %s \n", err)
		return err
	}

	return nil
}

// DumpPublicKeyToFile - Save one public key into a file named public.pem
func DumpPublicKeyToFile(publicKey *rsa.PublicKey) error {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Printf("error when dumping publickey: %s \n", err)
		return err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create("public.pem")
	if err != nil {
		fmt.Printf("error when create public.pem: %s \n", err)
		return err
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		fmt.Printf("error when encode public pem: %s \n", err)
		return err
	}

	return nil
}
