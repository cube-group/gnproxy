package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main(){
	engine:=gin.Default()
	engine.Any("/", func(context *gin.Context) {
		for k,v:=range context.Request.Header{
			fmt.Println(k,v)
		}
		context.Header("X-LINYANG","TRUE")
		context.String(http.StatusOK,"ok")
	})
	engine.Any("/_auth", func(context *gin.Context) {
		context.Header("X-LINYANG","TRUE")
		context.String(http.StatusOK,"ok")
	})
	engine.Run(":9001")
}
