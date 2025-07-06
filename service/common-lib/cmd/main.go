package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	systemLogger "my-go-api/common-lib/system_logger"
	"sync"
	"time"
)

func main() {
	logger := systemLogger.InitLogger(slog.LevelDebug)
	logger.Info("test", slog.String("test3", "test3"))

	// logger := slog.Default()
	// logger.Info("test", slog.String("test2", "test2"))
	// logger.Error("test")
	// logger.Warn("test")
	// logger.Debug("test")
	// logger.With("test", "test")
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// 	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
	// 		fmt.Println("groups", groups)
	// 		// fmt.Println("a", a)
	// 		// if a.Key == slog.LevelKey {

	// 		// }
	// 		return slog.String(slog.LevelKey, "test")
	// 		return a
	// 	},
	// }))
	// // fmt.Println("logger", logger)
	// logger.Info("test", slog.String("test2", "test2"))
	// logger.Error("test")
	// logger.Warn("test")
	// logger.Debug("test")
	// logger.With("test", "test")

	// awsConfig := awsconfig.InitAws("ap-southeast-1")
	// s3Client := s3.NewFromConfig(awsConfig)
	// object, err := s3Client.ListObjectsV2(
	// 	context.Background(),
	// 	&s3.ListObjectsV2Input{
	// 		Bucket:  aws.String("go-api-bucket-v1-21-6-2025"),
	// 		Prefix:  aws.String("jwk"),
	// 		MaxKeys: aws.Int32(5),
	// 		// StartAfter: aws.String("2025-06-28"),
	// 	},
	// )

	// if err != nil {
	// 	log.Panicln("error", err)
	// }
	// log.Println("object", *object.KeyCount)
	// var nearJwkKey string
	// var diff time.Duration = 0
	// now := time.Now()
	// log.Println("b4 diff", diff)
	// for _, obj := range object.Contents {
	// 	if diff == 0 {
	// 		diff = now.Sub(*obj.LastModified)
	// 		nearJwkKey = *obj.Key
	// 		continue
	// 	}
	// 	if now.Sub(*obj.LastModified) < diff {
	// 		diff = now.Sub(*obj.LastModified)
	// 		nearJwkKey = *obj.Key
	// 	}
	// }
	// log.Println("af diff", diff)
	// log.Println("nearJwkKey", nearJwkKey)
	// o, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
	// 	Bucket: aws.String("go-api-bucket-v1-21-6-2025"),
	// 	Key:    aws.String(nearJwkKey),
	// })
	// if err != nil {
	// 	log.Panicln("error", err)
	// }
	// defer o.Body.Close()
	// body, err := io.ReadAll(o.Body)
	// if err != nil {
	// 	log.Panicln("error", err)
	// }
	// log.Println("o", string(body))
	// errChan := make(chan error, 3)
	// waitGroup := sync.WaitGroup{}
	// waitGroup.Add(3)
	// go func(waitGroup *sync.WaitGroup, errChan chan error) {
	// 	select {
	// 	case <-ctx.Done():
	// 		errChan <- ctx.Err() // Timeout or cancellation
	// 		waitGroup.Done()
	// 		return
	// 	default:
	// 		if err := errors.New("error1"); err != nil {
	// 			// time.Sleep(1 * time.Second)
	// 			log.Println("error1", time.Now())
	// 			errChan <- err
	// 			waitGroup.Done()
	// 			return
	// 		}
	// 		errChan <- nil
	// 		waitGroup.Done()
	// 	}
	// }(&waitGroup, errChan)
	// go func(waitGroup *sync.WaitGroup, errChan chan error) {
	// 	select {
	// 	case <-ctx.Done():
	// 		errChan <- ctx.Err() // Timeout or cancellation
	// 		waitGroup.Done()
	// 		return
	// 	default:
	// 		if err := errors.New("error test"); err != nil {
	// 			log.Println("error test", time.Now())
	// 			errChan <- err
	// 			waitGroup.Done()
	// 			return
	// 		}
	// 		errChan <- nil
	// 		waitGroup.Done()
	// 	}
	// }(&waitGroup, errChan)
	// go func(waitGroup *sync.WaitGroup, errChan chan error) {
	// 	select {
	// 	case <-ctx.Done():
	// 		errChan <- ctx.Err() // Timeout or cancellation
	// 		waitGroup.Done()
	// 		return
	// 	default:
	// 		// time.Sleep(2 * time.Second)
	// 		log.Println("not error", time.Now())
	// 		errChan <- nil
	// 		waitGroup.Done()
	// 	}
	// }(&waitGroup, errChan)
	// waitGroup.Wait()
	// close(errChan)

	// select {
	// case err := <-errChan:
	// 	fmt.Println("err", err)
	// 	if err != nil {
	// 		fmt.Printf("Goroutine error: %v\n", err)
	// 	}
	// case <-ctx.Done():
	// 	fmt.Printf("Operation timed out: %v\n", ctx.Err())
	// }

	// defer func() {
	// 	// recover from panic
	// 	if r := recover(); r != nil {
	// 		log.Println("Recovered from panic:", r)
	// 	}
	// }()
	// ctx := context.Background()
	// ctx, cancel := context.WithCancelCause(ctx)
	// waitGroup := sync.WaitGroup{}
	// waitGroup.Add(1)
	// if err := action1(ctx, true); err != nil {
	// 	cancel(err)
	// }

	// go someAction(ctx, &waitGroup)

	// waitGroup.Wait()

	// log.Println("All done:", ctx.Err())
	// if ctx.Err() != nil {
	// 	panic(context.Cause(ctx))
	// }

	// log.Println("ctx err:", ctx.Err())
	// log.Println("cause:", context.Cause(ctx))
	// log.Println("err", err)
	// err := someAction(ctx)
	// log.Println("err", err)
	// testString := "test"
	// testStringPointer := &testString

	// testString2 := testStringPointer

	// *testString2 = "test2"
	// fmt.Println("testString", testString)
	// fmt.Println("testString2", *testString2)

	// awss3.AccessS3("test.txt", content)
	// fmt.Println("hi!")
	// router := router.InitRouter(func(router *gin.RouterGroup) {
	// 	router.GET("/health", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	// 	})
	// })

	// router.Run(":3001")
}

func action1(ctx context.Context, doError bool) error {
	// time.Sleep(6 * time.Second)
	if doError {
		return errors.New("SomeErr")
	}
	return nil
}

func action2(ctx context.Context, doError bool) error {
	if doError {
		return errors.New("SomeErr")
	}
	return nil
}

func someAction(ctx context.Context, waitGroup *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		log.Println("canceled cause:", context.Cause(ctx))
		waitGroup.Done()
	default:
		log.Println("Context not done, working")
		time.Sleep(1 * time.Second)
		waitGroup.Done()
	}
	// log.Println("Context done")
}
