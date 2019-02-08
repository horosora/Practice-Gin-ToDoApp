package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type ToDo struct {
	Id int
	Content string
}

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
		_, err = db.Exec("CREATE TABLE todo(id INTEGER PRIMARY KEY , contents TEXT NOT NULL)")
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

func getAllToDo() []ToDo {
	db ,err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM todo")
	if err != nil {
		log.Println(err)
	}

	var todos []ToDo
	for rows.Next() {
		var id int
		var content string
		var todo ToDo
		rows.Scan(&id, &content)
		todo.Id = id
		todo.Content = content
		todos = append(todos, todo)
	}

	rows.Close()

	return todos
}

func rmToDo(strId []string) {
	db ,err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for i := 0; i < len(strId); i++ {
		id, _ := strconv.Atoi(strId[i])
		_, err = db.Exec("DELETE FROM todo WHERE id = ?", id)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	DBInit()

	router.GET("/", func(context *gin.Context) {
		todo := getAllToDo()
		context.HTML(http.StatusOK, "index.tmpl", gin.H{"todo": todo})
	})

	router.POST("/add", func(context *gin.Context) {
		content := context.PostForm("contents")
		if content != "" {
			addToDo(content)
		}
		context.Redirect(http.StatusFound, "/")
	})

	router.POST("/rm", func(context *gin.Context) {
		rmToDo(context.PostFormArray("id"))
		context.Redirect(http.StatusFound, "/")
	})

	router.Run(":8080")
}
