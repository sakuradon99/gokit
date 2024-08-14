package web

import "github.com/gin-gonic/gin"

type ServerCustomEngineConfig interface {
	CustomEngine(engine *gin.Engine) error
}
