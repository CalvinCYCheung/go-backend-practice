package keyservice

import (
	"context"
	awss3 "my-go-api/common-lib/aws-service/aws-s3"
	generatekey "my-go-api/key-rotation/internal/key-service/generate_key"
	readkey "my-go-api/key-rotation/internal/key-service/read_key"
	uploadkey "my-go-api/key-rotation/internal/key-service/upload_key"
	"my-go-api/key-rotation/internal/model"
	"time"
)

type KeyService interface {
	GenrateKey() (model.Generated, error)
	UploadKey(ctx context.Context, generated model.Generated, previousKey *model.Generated) error
	GetPreviousKey(ctx context.Context) (*model.Generated, error)
}

type KeyServiceImpl struct {
	S3Client *awss3.S3AccessImpl
}

func (k *KeyServiceImpl) GenrateKey(ctx context.Context) (model.Generated, error) {
	generated, err := generatekey.Generate(ctx)
	if err != nil {
		// log.Panicln("error: ", err)
		return model.Generated{}, err
	}
	return generated, nil
}
func (k *KeyServiceImpl) UploadKey(ctx context.Context, generated model.Generated, previousKey *model.JWK) error {
	var fileName = []string{
		"private-" + generated.KeyIdentifier + ".key",
		"public-" + generated.KeyIdentifier + ".key",
		"jwk-" + generated.KeyIdentifier + ".json",
	}
	err := uploadkey.PutKeyToS3SecureStorage(ctx, generated, fileName, k.S3Client)
	if err != nil {
		return err
	}
	jwks, err := uploadkey.CreateJWKS(generated, *previousKey)
	if err != nil {
		return err
	}
	err = uploadkey.PutKeyToSharedS3(ctx, jwks, k.S3Client)
	if err != nil {
		return err
	}
	return nil
}
func (k *KeyServiceImpl) GetPreviousKey(ctx context.Context) (*model.JWK, error) {
	now := time.Now()
	var preKeyName string
	objs, err := k.S3Client.ListObjects(ctx, "jwk-", "go-api-bucket-v1-21-6-2025")
	if err != nil {
		return nil, err
	}
	// log.Println("objs: ", objs)
	var diff time.Duration = 0
	for _, obj := range objs.Contents {
		if diff == 0 {
			preKeyName = *obj.Key
			diff = now.Sub(*obj.LastModified)
			continue
		}
		if now.Sub(*obj.LastModified) < diff {
			diff = now.Sub(*obj.LastModified)
			preKeyName = *obj.Key
		}
	}
	previousKey, err := readkey.GetPreviousKey(ctx, k.S3Client, preKeyName)
	if err != nil {
		return nil, err
	}
	// log.Println("previousKey: ", previousKey)
	return previousKey, nil
}
