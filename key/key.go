// Package key handles the generation and other actions regarding RSA keypairs. These are required for signing of requests on behalf of users, and one will exist for every user in the system.
//
// This is essentially listed wholesale from go-fed/apcore: https://github.com/go-fed/apcore/blob/master/keys.go
package key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	// KeySize is the size of the keypair to use.
	KeySize = 1024
)

// New creates an RSA private key with the default size.
func New() (k *rsa.PrivateKey, err error) {
	k, err = rsa.GenerateKey(rand.Reader, KeySize)
	return
}

// NewOfSIze creates an RSA private key of a specified size.
func NewOfSize(n int) (k *rsa.PrivateKey, err error) {
	k, err = rsa.GenerateKey(rand.Reader, n)
	return
}

// MarshalPublicKey takes a public key and returns its string representation.
func MarshalPublicKey(p crypto.PublicKey) (string, error) {
	pkix, err := x509.MarshalPKIXPublicKey(p)
	if err != nil {
		return "", err
	}
	pb := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkix,
	})
	return string(pb), nil
}

// SerializeRSAPrivateKey takes a private key and returns its byte representation.
func SerializeRSAPrivateKey(k *rsa.PrivateKey) ([]byte, error) {
	return x509.MarshalPKCS8PrivateKey(k)
}

// DeserializeRSAPrivateKey takes byte representation of a private key and returns its PrivateKey representation.
func DeserializeRSAPrivateKey(b []byte) (crypto.PrivateKey, error) {
	return x509.ParsePKCS8PrivateKey(b)
}
