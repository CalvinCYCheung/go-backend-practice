package model

import (
	"crypto/rsa"
)

type KeyStore struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type KeyType string

const (
	PrivateKey KeyType = "private"
	PublicKey  KeyType = "public"
)
