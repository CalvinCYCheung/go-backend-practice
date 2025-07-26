package main

import (
	"context"
	"fmt"
	"my-go-api/auth/internal/route"
	tokenservice "my-go-api/auth/internal/token_service"
	awss3 "my-go-api/common-lib/aws-service/aws-s3"
	awsconfig "my-go-api/common-lib/aws-service/aws_config"
	"my-go-api/common-lib/aws-service/constant"
	"my-go-api/common-lib/router"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	awsConfig := awsconfig.InitAws(constant.ActivateRegionMap[constant.ApSouthEast1])
	var s3Client awss3.S3Access
	s3Client = &awss3.S3AccessImpl{
		Client: s3.NewFromConfig(awsConfig),
	}

	tokenService := tokenservice.NewTokenService(
		context.Background(),
		1*time.Hour,
		s3Client,
	)
	router := router.InitRouter(func(routerGroup *gin.RouterGroup) {
		routerGroup.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
		routerGroup.POST("/login", func(ctx *gin.Context) {
			token, err := tokenService.GenerateToken("user-123")
			if err != nil {
				fmt.Println(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"token": token})
		})
		routerGroup.POST("/validate", func(ctx *gin.Context) {
			token := ctx.GetHeader("Authorization")
			userId, err := tokenService.ValidateToken(token)
			if err != nil {
				fmt.Println(err)
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "Validated", "userId": userId})
		})
		// router.GET("/validate", func(ctx *gin.Context) {
		// 	gin.
		// 	tokenService.ValidateToken("token")
		// }),
	}, func(ctx *gin.Context) {
		fmt.Println(ctx.Request.URL.Path)
		if ctx.Request.URL.Path == route.Login {
			ctx.Next()
			return
		}
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			ctx.Abort()
			return
		}
		_, err := tokenService.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
		// ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	// 	router.POST("/login", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"login": "success"})
	// 	})
	// 	router.POST("/logout", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"logout": "success"})
	// 	})
	// 	router.POST("/register", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"register": "success"})
	// 	})
	// 	router.POST("/refresh", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"refresh": "success"})
	// 	})
	// 	router.POST("/forgot-password", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"forgot-password": "success"})
	// 	})
	// 	router.POST("/reset-password", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"reset-password": "success"})
	// 	})
	// 	router.POST("/verify-email", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"verify-email": "success"})
	// 	})
	// 	router.POST("/send-verification-email", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"send-verification-email": "success"})
	// 	})
	// 	router.POST("/send-reset-password-email", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"send-reset-password-email": "success"})
	// 	})
	// 	router.POST("/change-password", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"change-password": "success"})
	// 	})
	// 	router.POST("/change-email", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"change-email": "success"})
	// 	})
	// 	router.POST("/delete-account", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"delete-account": "success"})
	// 	})
	// })

	router.Run(":3002")
}
