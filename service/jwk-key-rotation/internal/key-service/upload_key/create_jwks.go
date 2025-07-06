package uploadkey

import (
	"encoding/json"
	"my-go-api/key-rotation/internal/model"
)

func CreateJWKS(generated model.Generated, previousKey model.JWK) ([]byte, error) {
	var newJWK model.JWK
	err := json.Unmarshal(generated.JWK, &newJWK)
	if err != nil {
		return nil, err
	}
	// log.Println("newJWK: ", newJWK)
	jwks := []model.JWK{
		newJWK,
		previousKey,
	}
	// jsonData := map[string]interface{}{
	// 	"keys": jwks,
	// }
	jwksJSON, err := json.Marshal(jwks)
	if err != nil {
		return nil, err
	}
	// log.Println("jwksJSON: ", string(jwksJSON))
	return jwksJSON, nil
}
