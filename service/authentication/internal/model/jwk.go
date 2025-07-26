package model

type JWK struct {
	Algorithm     string `json:"alg"`
	KeyIdentifier string `json:"kid"`
	KeyType       string `json:"kty"`
	Modulus       string `json:"n"`
	Exponent      string `json:"e"`
	KeyUse        string `json:"use"`
}
