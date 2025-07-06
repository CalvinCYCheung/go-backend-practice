package readkey

import (
	"context"
	"encoding/json"
	"io"
	awss3 "my-go-api/common-lib/aws-service/aws-s3"
	s3Model "my-go-api/common-lib/aws-service/aws-s3/model"
	"my-go-api/key-rotation/internal/model"
)

func GetPreviousKey(ctx context.Context, s3Client *awss3.S3AccessImpl, keyName string) (*model.JWK, error) {
	obj, err := s3Client.GetObject(ctx, s3Model.S3Object{
		BucketName: "go-api-bucket-v1-21-6-2025",
		KeyName:    keyName,
	})
	if err != nil {
		return nil, err
	}
	// log.Println("obj: ", obj)
	defer obj.Body.Close()
	body, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}
	// log.Println("body: ", string(body))
	var previousKey model.JWK
	err = json.Unmarshal(body, &previousKey)
	if err != nil {
		return nil, err
	}

	return &previousKey, nil
}
