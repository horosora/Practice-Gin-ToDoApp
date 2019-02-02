package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

type ToDo struct {
	gorm.Model
	Contents string
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.tmpl", gin.H{
			"msg": "Hello",
		})
	})

	router.Run(":8080")
}
