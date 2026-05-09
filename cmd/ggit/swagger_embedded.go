//go:build swagger

package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/nef734/ggit/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	registerSwagger = func(r *gin.Engine) {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
