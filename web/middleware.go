package web

import "github.com/gin-gonic/gin"

type Middleware interface {
	Register(path string) int
	Handle(ctx *gin.Context)
}

type middlewareHandlerWithOrder struct {
	fn    gin.HandlerFunc
	order int
}
