package keypair

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
)

// GenerateKeyPair - Generate one private and public key pair with bits length
func GenerateKeyPair(bits int) (*rsa.PrivateKey, error) {
	// generate key
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Println("keypair.GenerateKeyPair - Error generating key pair.")
		return nil, err
	}

	return key, err
}

// EncondePrivateKeyToPKCS1 - Returns the PKCS #1, ASN.1 DER form of a RSA private key.
func EncondePrivateKeyToPKCS1(privateKey *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return privBytes
}

// EncondePublicKeyToPKCS1 - Returns the PKCS #1, ASN.1 DER form of a RSA public key.
func EncondePublicKeyToPKCS1(publicKey *rsa.PublicKey) []byte {

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})

	return pubBytes
}

// EncondePrivateKeyToPKCS8 - Returns the PKCS #8, ASN.1 DER form of a RSA private key.
func EncondePrivateKeyToPKCS8(privateKey *rsa.PrivateKey) ([]byte, error) {
	privatePKCS8, err := x509.MarshalPKCS8PrivateKey(privateKey)

	if err != nil {
		log.Println("keypair.EncondePrivateKeyToPKCS8 - Error converting private key to PKCS8.")
		return nil, err
	}

	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privatePKCS8,
		},
	)

	return privBytes, nil
}

// EncondePublicKeyToPKCS8 - Returns the PKCS #8, ASN.1 DER form of a RSA public key.
func EncondePublicKeyToPKCS8(publicKey *rsa.PublicKey) ([]byte, error) {

	publicPKCS8, err := x509.MarshalPKIXPublicKey(publicKey)

	if err != nil {
		log.Println("keypair.EncondePublicKeyToPKCS8 - Error converting public key to PKCS8.")
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicPKCS8,
	})

	return pubBytes, nil
}

// SignDataSHA256PKCS1v15 - Sing the hash calculated with SHA256, using the PKCS #1 v15
func SignDataSHA256PKCS1v15(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHashSum := sha256.Sum256(data)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, msgHashSum[:])
	if err != nil {
		log.Println("keypair.SignData - Error signing the data.")
		return nil, err
	}

	return signature, err
}

// SignDataSHA256PKCS1v15Base64 - Sing the hash calculated with SHA256, using the PKCS #1 v15 and returning it in a base64 format.
func SignDataSHA256PKCS1v15Base64(data []byte, privateKey *rsa.PrivateKey) (string, error) {
	signaturePKCS, err := SignDataSHA256PKCS1v15(data, privateKey)
	if err != nil {
		log.Println("keypair.SignDataSHA256PKCS1v15Base64 - Error gernerating signature.")
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signaturePKCS), nil
}

// SignDataSHA256PSS - Sing the hash calculated with SHA256, using PSS
func SignDataSHA256PSS(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHashSum := sha256.Sum256(data)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum[:], nil)
	if err != nil {
		log.Println("keypair.SignData - Error signing the data.")
		return nil, err
	}

	return signature, err
}

// SignDataSHA256PSSBase64 - Sing the hash calculated with SHA256, using PSS and returning it in a base64 format.
func SignDataSHA256PSSBase64(data []byte, privateKey *rsa.PrivateKey) (string, error) {

	signaturePSS, err := SignDataSHA256PSS(data, privateKey)
	if err != nil {
		log.Println("keypair.SignDataSHA256PSSBase64 - Error generating signature.")
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signaturePSS), nil
}

// VerifyPKCS1Signature - Verifies if the message's signature matches its origin, returning an error when it doesn't.
func VerifyPKCS1Signature(data string, signature []byte, publicKey *rsa.PublicKey) error {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHashSum := sha256.Sum256([]byte(data))

	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, msgHashSum[:], signature)
	if err != nil {
		log.Println("keypair.VerifySignature - Error verifying signature.")
		return err
	}

	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	return nil
}

// VerifyBase64PKCS1Signature - Verify verifies that the message's signature matches its origin, returning an error when it doesn't.
func VerifyBase64PKCS1Signature(data string, signatureBase64 string, publicKey *rsa.PublicKey) error {

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		log.Println("keypair.VerifyBase64Signature - Error converting signature from base64.")
		return err
	}

	err = VerifyPKCS1Signature(data, signature, publicKey)

	if err != nil {
		log.Println("keypair.VerifyBase64Signature - Error verifying signature.")
		return err
	}
	return nil
}

// VerifyPSSSignature - Verifies if the message's signature matches its origin, using the PSS algorithm and returning an error when it doesn't.
func VerifyPSSSignature(data string, signature []byte, publicKey *rsa.PublicKey) error {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHashSum := sha256.Sum256([]byte(data))

	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSum[:], signature, nil)
	if err != nil {
		log.Println("keypair.VerifyPSSSignature - Error verifying signature.")
		return err
	}

	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	return nil
}

// VerifyBase64PSSSignature - Verifies if the message's signature matches its origin, using the PSS algorithm and returning an error when it doesn't.
func VerifyBase64PSSSignature(data string, signatureBase64 string, publicKey *rsa.PublicKey) error {

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		log.Println("keypair.VerifyBase64PSSSignature - Error converting signature from base64.")
		return err
	}

	err = VerifyPSSSignature(data, signature, publicKey)

	if err != nil {
		log.Println("keypair.VerifyBase64PSSSignature - Error verifying signature.")
		return err
	}
	return nil
}

// EncryptDataOAEP - Retunrs encrypted data using OAEP.
func EncryptDataOAEP(data, label string, publicKey *rsa.PublicKey) ([]byte, error) {

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(data), []byte(label))
	if err != nil {
		log.Println("keypair.EncryptDataOAEP - Error encrypting data.")
		return nil, err
	}

	return ciphertext, nil
}

// EncryptDataOAEPBase64 - Retunrs base64 encrypted data using OAEP.
func EncryptDataOAEPBase64(data, label string, publicKey *rsa.PublicKey) (string, error) {

	encryptDataOAEP, err := EncryptDataOAEP(data, label, publicKey)
	if err != nil {
		log.Println("keypair.EncryptDataOAEPBase64 - Error encrypting data.")
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptDataOAEP), nil
}

// DecryptDataOAEP - Retunrs decrypted data using OAEP.
func DecryptDataOAEP(data, label string, privateKey *rsa.PrivateKey) ([]byte, error) {

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, []byte(data), []byte(label))
	if err != nil {
		log.Println("keypair.DecryptDataOAEP - Error decrypting data.")
		return nil, err
	}

	return plaintext, nil
}

// DecryptBase64DataOAEP - Retunrs decrypted base64 data using OAEP.
func DecryptBase64DataOAEP(dataBase64, label string, privateKey *rsa.PrivateKey) (string, error) {

	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		log.Println("keypair.DecryptBase64DataOAEP - Error converting signature from base64.")
		return "", err
	}

	plaintext, err := DecryptDataOAEP(string(data), label, privateKey)
	if err != nil {
		log.Println("keypair.DecryptBase64DataOAEP - Error decrypting data.")
		return "", err
	}

	return string(plaintext), nil
}

// EncryptDataPKCS1v15 - Retunrs encrypted data using PKCS #01 v15.
func EncryptDataPKCS1v15(data string, publicKey *rsa.PublicKey) ([]byte, error) {

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data))
	if err != nil {
		log.Println("keypair.EncryptDataPKCS1v15 - Error encrypting data.")
		return nil, err
	}

	return ciphertext, nil
}

// EncryptDataPKCS1v15Base64 - Retunrs base64 encrypted data using PKCS #01 v15.
func EncryptDataPKCS1v15Base64(data string, publicKey *rsa.PublicKey) (string, error) {

	encryptDataPKCS1v15, err := EncryptDataPKCS1v15(data, publicKey)
	if err != nil {
		log.Println("keypair.EncryptDataPKCS1v15Base64 - Error encrypting data.")
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptDataPKCS1v15), nil
}

// DecryptDataPKS1v15 - Retunrs decrypted data using PKS #1 v15.
func DecryptDataPKS1v15(data string, privateKey *rsa.PrivateKey) ([]byte, error) {

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(data))
	if err != nil {
		log.Println("keypair.DecryptDataPKS1v15 - Error decrypting data.")
		return nil, err
	}

	return plaintext, nil
}

// DecryptBase64DataPKS1v15 - Retunrs decrypted base64 data using PKS #1 v15.
func DecryptBase64DataPKS1v15(dataBase64 string, privateKey *rsa.PrivateKey) (string, error) {

	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		log.Println("keypair.DecryptBase64DataPKS1v15 - Error converting signature from base64.")
		return "", err
	}

	plaintext, err := DecryptDataPKS1v15(string(data), privateKey)
	if err != nil {
		log.Println("keypair.DecryptBase64DataPKS1v15 - Error decrypting data.")
		return "", err
	}

	return string(plaintext), nil
}
