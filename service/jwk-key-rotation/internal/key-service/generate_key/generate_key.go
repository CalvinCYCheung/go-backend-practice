package generatekey

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"my-go-api/key-rotation/internal/model"
	"os"
)

func Generate(ctx context.Context) (model.Generated, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return model.Generated{}, errors.New("failed to generate rsa key")
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPem := pem.EncodeToMemory(privateKeyBlock)
	if privateKeyPem == nil {
		return model.Generated{}, errors.New("Failed to encode private key")
	}
	// os.WriteFile("./private.key", privateKeyPem, 0644)

	publicKey := &privateKey.PublicKey
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyPem := pem.EncodeToMemory(publicKeyBlock)
	if publicKeyPem == nil {
		return model.Generated{}, errors.New("Failed to encode public key")
	}
	// os.WriteFile("./public.key", publicKeyPem, 0644)

	jwk := model.JWK{
		Algorithm: "RS256",
		KeyType:   "RSA",
		Modulus:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		Exponent:  base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
		KeyUse:    "sig",
	}
	jwkBytes, err := json.Marshal(jwk)
	if err != nil {
		return model.Generated{}, errors.New("Failed to marshal jwk")
	}
	jwkHash := sha256.Sum256(jwkBytes)
	hash := base64.RawURLEncoding.EncodeToString(jwkHash[:])
	// fmt.Println(hash)
	jwk.KeyIdentifier = hash
	jwkBytes, err = json.Marshal(jwk)
	if err != nil {
		return model.Generated{}, errors.New("Failed to marshal jwk")
	}
	// os.WriteFile("./jwk.json", jwkBytes, 0644)
	return model.Generated{
		PrivateKey:    privateKeyPem,
		PublicKey:     publicKeyPem,
		JWK:           jwkBytes,
		KeyIdentifier: jwk.KeyIdentifier,
	}, nil
}

func GenerateToLocal() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println(err)
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPem := pem.EncodeToMemory(privateKeyBlock)
	if privateKeyPem == nil {
		fmt.Println("Failed to encode private key")
		return
	}
	os.WriteFile("./private.key", privateKeyPem, 0644)

	publicKey := &privateKey.PublicKey
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyPem := pem.EncodeToMemory(publicKeyBlock)
	if publicKeyPem == nil {
		fmt.Println("Failed to encode public key")
		return
	}
	os.WriteFile("./public.key", publicKeyPem, 0644)

	jwk := model.JWK{
		Algorithm: "RS256",
		KeyType:   "RSA",
		Modulus:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		Exponent:  base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
		KeyUse:    "sig",
	}

	jwkBytes, err := json.Marshal(jwk)
	if err != nil {
		fmt.Println(err)
		return
	}

	jwkHash := sha256.Sum256(jwkBytes)
	hash := base64.RawURLEncoding.EncodeToString(jwkHash[:])
	fmt.Println(hash)
	jwk.KeyIdentifier = hash
	jwkBytes, err = json.Marshal(jwk)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.WriteFile("./jwk.json", jwkBytes, 0644)
}
