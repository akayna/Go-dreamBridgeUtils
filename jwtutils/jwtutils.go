package jwtutils

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"math/big"
)

// JSONKey represents a raw key inside a JWKS.
type JSONKey struct {
	Curve       string `json:"crv"`
	Exponent    string `json:"e"`
	ID          string `json:"kid"`
	Modulus     string `json:"n"`
	X           string `json:"x"`
	Y           string `json:"y"`
	precomputed interface{}
}

// RSA - Parses a JSONKey and turns it into an RSA public key.
func (j *JSONKey) RSA() (publicKey *rsa.PublicKey, err error) {

	// Check if the key has already been computed.
	if j.precomputed != nil {
		return j.precomputed.(*rsa.PublicKey), nil
	}

	// Confirm everything needed is present.
	if j.Exponent == "" || j.Modulus == "" {
		return nil, errors.New("missing assets")
	}

	// Decode the exponent from Base64.
	//
	// According to RFC 7518, this is a Base64 URL unsigned integer.
	// https://tools.ietf.org/html/rfc7518#section-6.3
	var exponent []byte
	if exponent, err = base64.RawURLEncoding.DecodeString(j.Exponent); err != nil {
		return nil, err
	}

	// Decode the modulus from Base64.
	var modulus []byte
	if modulus, err = base64.RawURLEncoding.DecodeString(j.Modulus); err != nil {
		return nil, err
	}

	// Create the RSA public key.
	publicKey = &rsa.PublicKey{}

	// Turn the exponent into an integer.
	//
	// According to RFC 7517, these numbers are in big-endian format.
	// https://tools.ietf.org/html/rfc7517#appendix-A.1
	publicKey.E = int(big.NewInt(0).SetBytes(exponent).Uint64())

	// Turn the modulus into a *big.Int.
	publicKey.N = big.NewInt(0).SetBytes(modulus)

	// Keep the public key so it won't have to be computed every time.
	j.precomputed = publicKey

	return publicKey, nil
}

// Populate - Populate the struct with the given json
func (j *JSONKey) Populate(jwkJson string) error {
	// Converts the json jwk to rsa key
	err := json.Unmarshal([]byte(jwkJson), &j)

	if err != nil {
		log.Println("microform - getPublicKey - Error parsing json to jwk.")
		return err
	}

	return nil
}
