package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type ParseBody struct {
	Language string `json:"language"`
	Body     string `json:"body"`
}

func main() {

	r := gin.Default()
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/ascii-table", func(context *gin.Context) {
		var body ParseBody

		err := context.BindJSON(&body)
		if err != nil {
			return
		}

		s := make(map[string][][]string)
		var c CVar
		cStr := body.Body
		c.parse(cStr)
		getTable(&c, s)
		var Ascii string
		for _, v := range s {
			Ascii += fmt.Sprintf("%s\n\n\n", MultilineComment(getTableFormatString(v)))
		}
		context.Data(http.StatusOK, binding.MIMEHTML, []byte(Ascii))
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
