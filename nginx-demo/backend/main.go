package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main(){
	engine:=gin.Default()
	engine.GET("/abc", func(context *gin.Context) {
		for k,v:=range context.Request.Header{
			fmt.Println(k,v)
		}
		context.String(http.StatusOK,"backend")
	})
	engine.Run(":9002")
}
