package tokenservice

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"my-go-api/auth/internal/model"
	awss3 "my-go-api/common-lib/aws-service/aws-s3"
	s3model "my-go-api/common-lib/aws-service/aws-s3/model"
	cachemodel "my-go-api/common-lib/cache_model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(userId string) (*string, error)
	ValidateToken(token string) (*string, error)
}

type TokenServiceImpl struct {
	jwks          cachemodel.CacheModel[[]model.JWK]
	privateKey    cachemodel.CacheModel[*rsa.PrivateKey]
	publicKeys    cachemodel.CacheModel[[]*rsa.PublicKey]
	tokenDuration time.Duration
	s3Client      awss3.S3Access
}

func NewTokenService(
	ctx context.Context,
	tokenDuration time.Duration,
	s3Client awss3.S3Access,
) TokenService {
	jwks, privateKey, publicKeys := fetchLatestRecord(ctx, s3Client)
	return &TokenServiceImpl{
		tokenDuration: tokenDuration,
		s3Client:      s3Client,
		jwks: cachemodel.CacheModel[[]model.JWK]{
			Data:           jwks,
			HardExpireTime: time.Now().Add(time.Second * 10),
		},
		privateKey: cachemodel.CacheModel[*rsa.PrivateKey]{
			Data:           privateKey,
			HardExpireTime: time.Now().Add(time.Second * 10),
		},
		publicKeys: cachemodel.CacheModel[[]*rsa.PublicKey]{
			Data:           publicKeys,
			HardExpireTime: time.Now().Add(time.Second * 10),
		},
	}
}

func (t *TokenServiceImpl) GenerateToken(userId string) (*string, error) {
	// TODO: Generate token with private key
	// (*t.S3Client)
	var privateKey = t.privateKey.Data
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(t.tokenDuration).Unix(),
		"kid": t.jwks.Data[0].KeyIdentifier,
	})
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println("Signed string error: ", err)
		return nil, err
	}
	go t.renewKey()
	return &tokenString, nil
}

func (t *TokenServiceImpl) ValidateToken(token string) (*string, error) {
	// TODO: Get public key from S3 bucket , return error if the key is not found

	claims := model.Claims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		for i, jwk := range t.jwks.Data {
			if jwk.KeyIdentifier == claims.Kid {
				return t.publicKeys.Data[i], nil
			}
		}
		return nil, errors.New("key not found")
	}, jwt.WithValidMethods([]string{"RS256"}))
	// fmt.Println("parsedToken: ", parsedToken.Valid)
	if err != nil {
		return nil, err
	}
	go t.renewKey()
	if sub, err := claims.GetSubject(); err != nil {
		return nil, err
	} else {
		return &sub, nil
	}

}

func (t *TokenServiceImpl) renewKey() {
	// TODO: Renew stored jwks and private key
	if !t.jwks.HardExpireTime.Before(time.Now()) {
		fmt.Println("Hard expire time: ", t.jwks.HardExpireTime)
		fmt.Println("Time now: ", time.Now())
		return
	}
	fmt.Println("renewing key")
	jwks, privateKey, publicKeys := fetchLatestRecord(context.Background(), t.s3Client)
	t.jwks = cachemodel.CacheModel[[]model.JWK]{
		Data:           jwks,
		HardExpireTime: time.Now().Add(time.Hour * 24),
	}
	t.privateKey = cachemodel.CacheModel[*rsa.PrivateKey]{
		Data:           privateKey,
		HardExpireTime: time.Now().Add(time.Hour * 24),
	}
	t.publicKeys = cachemodel.CacheModel[[]*rsa.PublicKey]{
		Data:           publicKeys,
		HardExpireTime: time.Now().Add(time.Hour * 24),
	}
}

func fetchLatestRecord(
	ctx context.Context,
	s3Client awss3.S3Access,
) ([]model.JWK, *rsa.PrivateKey, []*rsa.PublicKey) {
	object, _ := s3Client.GetObject(ctx, s3model.S3Object{
		BucketName: "goback-end-shared-bucket",
		KeyName:    ".well-known/jwks.json",
	})
	defer object.Body.Close()
	body, err := io.ReadAll(object.Body)
	if err != nil {
		panic(err)
	}
	data := []model.JWK{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	privateKey, err := getKey(ctx, s3Client, data[0].KeyIdentifier)
	if err != nil {
		panic(err)
	}

	publicKeys, err := getPublicKeys(ctx, s3Client, data)
	if err != nil {
		panic(err)
	}
	return data, privateKey, publicKeys
}

func getPublicKeys(
	ctx context.Context,
	s3Client awss3.S3Access,
	jwks []model.JWK,
) ([]*rsa.PublicKey, error) {
	var publicKeys []*rsa.PublicKey
	for _, jwk := range jwks {
		object, err := s3Client.GetObject(ctx, s3model.S3Object{
			BucketName: "go-api-bucket-v1-21-6-2025",
			KeyName:    fmt.Sprintf("public-%s.key", jwk.KeyIdentifier),
		})
		if err != nil {
			return nil, err
		}
		defer object.Body.Close()
		body, err := io.ReadAll(object.Body)
		if err != nil {
			return nil, err
		}
		block, _ := pem.Decode(body)
		if block.Type != "RSA PUBLIC KEY" {
			return nil, errors.New("invalid public key type")
		}
		key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		publicKeys = append(publicKeys, key)
	}
	return publicKeys, nil
}

func getKey(
	ctx context.Context,
	s3Client awss3.S3Access,
	kid string,
) (*rsa.PrivateKey, error) {
	privateKeyName := fmt.Sprintf("%s-%s.key", model.PrivateKey, kid)
	privateKey, err := s3Client.GetObject(ctx, s3model.S3Object{
		BucketName: "go-api-bucket-v1-21-6-2025",
		KeyName:    privateKeyName,
	})
	if err != nil {
		return nil, err
	}
	defer privateKey.Body.Close()
	privateKeyBody, err := io.ReadAll(privateKey.Body)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKeyBody)
	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key type")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}
