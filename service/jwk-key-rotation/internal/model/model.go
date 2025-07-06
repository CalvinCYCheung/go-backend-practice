package model

import "fmt"

type Generated struct {
	PrivateKey    []byte
	PublicKey     []byte
	JWK           []byte
	KeyIdentifier string
}

func (g Generated) String() string {
	prvLen := len(g.PrivateKey)
	pubLen := len(g.PublicKey)
	return fmt.Sprintf(
		"PrivateKey: %d, PublicKey: %d, JWK: %s",
		prvLen,
		pubLen,
		g.JWK)
}

type JWK struct {
	Algorithm     string `json:"alg"`
	KeyIdentifier string `json:"kid"`
	KeyType       string `json:"kty"`
	Modulus       string `json:"n"`
	Exponent      string `json:"e"`
	KeyUse        string `json:"use"`
}
