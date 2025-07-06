package uploadkey

import (
	"context"
	"errors"
	awss3 "my-go-api/common-lib/aws-service/aws-s3"
	s3Model "my-go-api/common-lib/aws-service/aws-s3/model"
	errorsmodel "my-go-api/common-lib/errors_model"
	"my-go-api/key-rotation/internal/model"

	"golang.org/x/sync/errgroup"
)

func PutKeyToSharedS3(
	ctx context.Context,
	jwks []byte,
	s3Client *awss3.S3AccessImpl,
) error {
	err := s3Client.PutObject(ctx, s3Model.S3Object{
		BucketName:  "goback-end-shared-bucket",
		KeyName:     ".well-known/jwks.json",
		FileContent: jwks,
	})
	if err != nil {
		return err
	}
	return nil
}

func PutKeyToS3SecureStorage(
	ctx context.Context,
	generated model.Generated,
	fileName []string,
	s3Client *awss3.S3AccessImpl,
) error {
	g := errgroup.Group{}
	for i, file := range fileName {
		g.Go(func() error {
			var fileContent []byte
			if i == 0 {
				fileContent = generated.PrivateKey
			} else if i == 1 {
				fileContent = generated.PublicKey
			} else if i == 2 {
				fileContent = generated.JWK
			}
			return s3Client.PutObject(ctx, s3Model.S3Object{
				BucketName:  "go-api-bucket-v1-21-6-2025",
				KeyName:     file,
				FileContent: fileContent,
			})
		})
	}
	if err := g.Wait(); err != nil {
		var errModel *errorsmodel.ErrorModel[string]
		if errors.As(err, &errModel) {
			err := cleanUp(ctx, s3Client, fileName, errModel.Data)
			if err != nil {
				return err
			}
			return errModel.Err
		}
	}
	return nil
}

func cleanUp(
	ctx context.Context,
	s3Client *awss3.S3AccessImpl,
	fileName []string,
	ignoreKeyName string,
) error {
	for _, file := range fileName {
		if file == ignoreKeyName {
			continue
		}
		err := s3Client.DeleteObject(ctx, s3Model.S3Object{
			BucketName: "go-api-bucket-v1-21-6-2025",
			KeyName:    file,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
