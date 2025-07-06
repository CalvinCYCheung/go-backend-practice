package main

import (
	"encoding/json"
	"fmt"
	"io"
	"my-go-api/common-lib/router"
	"net/http"

	docs "my-go-api/cat-image/docs"
	"my-go-api/cat-image/internal"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var num = 0

// @BasePath /api/v1

// PingExample godoc
// @Summary get cat images
// @Schemes
// @Description get cat images
// @Tags cat-images
// @Accept json
// @Produce json
// @Success 200 {object} internal.Response[internal.CatImage] "description"
// @Router /cat-images [GET]
func getCatImages(ctx *gin.Context) {
	fmt.Println("Request received")
	response, err := http.Get("https://api.thecatapi.com/v1/images/search?limit=10")
	fmt.Printf("response %v\n", response)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot get images"})
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}
	// var images []map[string]interface{}
	var images []internal.CatImage
	err = json.Unmarshal(body, &images)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot unmarshal body"})
		return
	}
	num++
	fmt.Printf("num %v\n", num)
	ctx.JSON(http.StatusOK, gin.H{
		"data": &internal.Response[internal.CatImage]{
			Data:    images,
			Result:  "success",
			Message: "",
		}})
}

func main() {
	router := router.InitRouter(func(router *gin.RouterGroup) {
		router.GET("/cat-images", getCatImages)
	},
		func(c *gin.Context) {
			c.Header("Content-Type", "application/json")
			c.Header("X-Cat", "Meow")
			// c.Header("Cache-Control", "max-age=100, stale-if-error=30")
			c.Next()
		})

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(":3001")
}
