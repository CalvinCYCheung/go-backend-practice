package main

import (
	"context"
	"errors"
	awss3 "my-go-api/common-lib/aws-service/aws-s3"
	awsconfig "my-go-api/common-lib/aws-service/aws_config"
	"my-go-api/common-lib/aws-service/constant"
	systemLogger "my-go-api/common-lib/system_logger"
	keyservice "my-go-api/key-rotation/internal/key-service"

	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) error {
	// now := time.Now().Format(time.DateOnly)
	// var fileName = []string{"private-" + now + ".key", "public-" + now + ".key", "jwk-" + now + ".json"}
	// ctx, cancel := context.WithCancelCause(ctx)
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// ctx := context.Background()
	ctx = context.WithValue(ctx, "test", "test")

	logger := systemLogger.InitLogger(slog.LevelDebug)
	awsConfig := awsconfig.InitAws(constant.ActivateRegionMap[constant.ApSouthEast1])
	// var s3Client awss3.S3Access
	s3Client := awss3.S3AccessImpl{
		Client: s3.NewFromConfig(awsConfig),
		IsMock: false,
	}
	keyService := keyservice.KeyServiceImpl{
		S3Client: &s3Client,
	}
	// log.Panicln("tests", keyService)
	previousKey, err := keyService.GetPreviousKey(ctx)
	if err != nil {
		logger.Error("Read error", slog.String("error", err.Error()))
		return errors.New(err.Error())
		// log.Panicln("Read error: ", err)
	}
	// log.Println("previousKey: ", previousKey)
	generated, err := keyService.GenrateKey(ctx)
	if err != nil {
		logger.Error("Generate error", slog.String("error", err.Error()))
		return errors.New(err.Error())
		// log.Panicln("Generate error: ", err)
	}
	err = keyService.UploadKey(ctx, generated, previousKey)
	if err != nil {
		logger.Error("Upload error", slog.String("error", err.Error()))
		return errors.New(err.Error())
		// log.Panicln("Upload error: ", err)
	}
	logger.Info("All Done")
	return nil
}
