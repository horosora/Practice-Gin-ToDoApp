package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func DBInit() {
	db ,err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master")
	if err != nil {
		panic(err)
	}

	makeDBTableFlag := true
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err)
		}
		if tableName  == "todo" {
			makeDBTableFlag = false
		}
	}

	rows.Close()

	if makeDBTableFlag {
		_, err = db.Exec("CREATE TABLE todo(contents TEXT)")
		if err != nil {
			panic(err)
		}
	}
}

func addToDo(contents string) {
	db ,err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO todo(contents) VALUES(?)", contents)
	if err != nil {
		panic(err)
	}
}

func getAll() []string {
	db ,err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT contents FROM todo")
	if err != nil {
		log.Println(err)
	}

	var contents []string
	for rows.Next() {
		var content string
		err := rows.Scan(&content)
		if err != nil {
			log.Println(err)
		}
		contents = append(contents, content)
	}

	rows.Close()

	return contents
}

func rmAll() {
	db ,err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM todo")
	if err != nil {
		log.Println(err)
	}
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
