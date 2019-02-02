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

func DBInit() {
	db, err := gorm.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&ToDo{})
}

func addToDo(contents string) {
	db, err := gorm.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	db.Create(&ToDo{Contents: contents})
}

func getAll() []ToDo {
	db, err := gorm.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	var todo []ToDo
	db.Find(&todo)
	return todo
}

func rmAll() {
	db, err := gorm.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	db.Delete(ToDo{})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	DBInit()

	router.GET("/", func(context *gin.Context) {
		todo := getAll()
		context.HTML(http.StatusOK, "index.tmpl", gin.H{"todo": todo})
	})

	router.POST("/add", func(context *gin.Context) {
		content := context.PostForm("contents")
		if content != "" {
			addToDo(content)
		}
		context.Redirect(http.StatusFound, "/")
	})

	router.GET("/rm", func(context *gin.Context) {
		rmAll()
		context.Redirect(http.StatusFound, "/")
	})

	router.Run(":8080")
}
