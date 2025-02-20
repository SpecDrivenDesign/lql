package signing

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadPrivateKey reads a PEM file and parses it as an RSA private key.
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	// Clean and validate the filename.
	cleanFilename := filepath.Clean(filename)
	if !strings.HasSuffix(cleanFilename, ".pem") {
		return nil, fmt.Errorf("invalid private key file: expected a .pem file")
	}

	data, err := os.ReadFile(cleanFilename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", cleanFilename, err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing RSA private key")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSA private key: %v", err)
	}
	return privKey, nil
}

// LoadPublicKey reads a PEM file and returns an RSA public key.
func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	// Clean and validate the filename.
	cleanFilename := filepath.Clean(filename)
	if !strings.HasSuffix(cleanFilename, ".pem") {
		return nil, fmt.Errorf("invalid public key file: expected a .pem file")
	}

	data, err := os.ReadFile(cleanFilename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", cleanFilename, err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	switch block.Type {
	case "RSA PUBLIC KEY":
		pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err == nil {
			return pub, nil
		}
		// If parsing as PKCS1 fails, try PKIX.
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing RSA public key: %v", err)
		}
		rsaPub, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}
		return rsaPub, nil
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing public key: %v", err)
		}
		rsaPub, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}
		return rsaPub, nil
	default:
		return nil, fmt.Errorf("unsupported key type %q", block.Type)
	}
}
