package awss3

import (
	"bytes"
	"context"
	"errors"
	"log"
	"my-go-api/common-lib/aws-service/aws-s3/model"
	errorsmodel "my-go-api/common-lib/errors_model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Access interface {
	PutObject(ctx context.Context, s3Object model.S3Object) error
	DeleteObject(ctx context.Context, s3Object model.S3Object) error
	GetObject(ctx context.Context, s3Object model.S3Object) (*s3.GetObjectOutput, error)
	ListObjects(ctx context.Context, prefix string, bucketName string) (*s3.ListObjectsV2Output, error)
}

type S3AccessImpl struct {
	Client *s3.Client
	IsMock bool
}

func (s *S3AccessImpl) PutObject(ctx context.Context, s3Object model.S3Object) error {
	if s.IsMock {
		return &errorsmodel.ErrorModel[string]{
			Err:  errors.New("error test"),
			Data: s3Object.KeyName,
		}
		// return errors.New("error test")
	}
	// log.Println("PutObject: ", s3Object.KeyName, time.Now())
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3Object.BucketName),
		Key:    aws.String(s3Object.KeyName),
		Body:   bytes.NewReader(s3Object.FileContent),
	})
	if err != nil {
		// log.Println("PutObject Error", err)
		// return errorsmodel.NewErrorModel[string](err, s3Object.KeyName)
		return err
	}
	// log.Println(obj.)
	return nil
}

func (s *S3AccessImpl) DeleteObject(ctx context.Context, s3Object model.S3Object) error {
	if s.IsMock {
		return errors.New("test errors")
	}
	// log.Println("DeleteObject: ", s3Object.KeyName, time.Now())
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s3Object.BucketName),
		Key:    aws.String(s3Object.KeyName),
	})
	return err
}

func (s *S3AccessImpl) GetObject(ctx context.Context, s3Object model.S3Object) (*s3.GetObjectOutput, error) {
	if s.IsMock {
		return nil, errors.New("test errors")
	}
	// log.Println("GetObject: ", s3Object.KeyName, time.Now())
	object, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s3Object.BucketName),
		Key:    aws.String(s3Object.KeyName),
	})
	if err != nil {
		return nil, err
	}
	return object, nil
}
func (s *S3AccessImpl) ListObjects(ctx context.Context, prefix string, bucketName string) (*s3.ListObjectsV2Output, error) {
	if s.IsMock {
		return nil, errors.New("test errors")
	}
	// log.Println("ListObjects with prefix: ", prefix, time.Now())
	objs, err := s.Client.ListObjectsV2(
		ctx,
		&s3.ListObjectsV2Input{
			Bucket:  aws.String(bucketName),
			Prefix:  aws.String(prefix),
			MaxKeys: aws.Int32(5),
		},
	)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

func AccessS3(fileKeyName string, fileContent []byte) {

	// file, err := os.Open("./privatekey.key")
	// if err != nil {
	// 	log.Panicln("Error reading private key file:", err)
	// }
	// defer file.Close()
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-southeast-1"),
	)
	if err != nil {
		log.Panicln(err)
	}
	client := s3.NewFromConfig(cfg)
	obj, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("go-api-bucket-v1-21-6-2025"),
		Key:    aws.String(fileKeyName),
		Body:   bytes.NewReader(fileContent),
	})
	if err != nil {
		log.Panicln(err)
	}
	log.Println(obj.Size)
}
