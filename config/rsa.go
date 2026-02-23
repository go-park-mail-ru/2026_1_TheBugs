package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("неверный формат PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("неверный тип ключа")
	}

	return rsaKey, nil
}

func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("неверный формат PEM")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("неверный тип ключа")
	}

	return rsaKey, nil
}
