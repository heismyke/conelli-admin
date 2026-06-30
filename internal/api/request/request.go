package request

import "github.com/gin-gonic/gin"

type Request struct {
	Context *gin.Context
	Body    interface{}
}
