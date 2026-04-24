package main

import "github.com/gin-gonic/gin"

// registerSwagger is a function pointer that will be set by the included swagger implementation.
// By default, it does nothing.
var registerSwagger = func(_ *gin.Engine) {
	// Swagger is disabled by default
}
