package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Todo struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Desc string `json:"desc"`
	Done bool   `json:"done"`
}

type CreateTodoInput struct {
	Desc string `json:"desc" binding:"required"`
	Done bool   `json:"done" binding:"required"`
}

type UpdateTodoInput struct {
	Desc string `json:"desc"`
	Done bool   `json:"done"`
}

var DB *gorm.DB

func createDefaultTodos() {
	todo1 := Todo{Desc: "host usrvtodo on distro.watch", Done: true}
	todo2 := Todo{Desc: "deploy usrvtodo on k3s", Done: false}
	DB.Create(&todo1)
	DB.Create(&todo2)
}

func ConnectDB() {
	dbPath := "todo.db"
	if os.Getenv("DB_PATH") != "" {
		dbPath = os.Getenv("DB_PATH")
	}
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("DB error:", err)
	}
	log.Println("connected to database:", db)
	db.AutoMigrate(&Todo{})
	DB = db
	createDefaultTodos()
}

func GetTodos(c *gin.Context) {
	var todos []Todo
	DB.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

func GetTodo(c *gin.Context) {
	var todo Todo
	err := DB.Where("id = ?", c.Param("id")).First(&todo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notfound"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func CreateTodo(c *gin.Context) {
	var input CreateTodoInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo := Todo{Desc: input.Desc, Done: input.Done}
	DB.Create(&todo)
	c.JSON(http.StatusOK, todo)
}

func UpdateTodo(c *gin.Context) {
	var todo Todo
	err := DB.Where("id = ?", c.Param("id")).First(&todo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notfound"})
		return
	}
	var input UpdateTodoInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Model(&todo).Updates(input)
	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	var todo Todo
	err := DB.Where("id = ?", c.Param("id")).First(&todo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notfound"})
		return
	}
	DB.Delete(&todo)
	c.JSON(http.StatusOK, todo)
}

func GetTodos2(c *gin.Context) {
	var todos []Todo
	DB.Find(&todos)
	c.HTML(http.StatusOK, "index.tmpl", todos)
}

func GetTodo2(c *gin.Context) { // actually implement this
	formID := c.PostForm("id")
	id, err := strconv.ParseInt(formID, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "badparse"})
		return
	}
	var todo Todo
	err = DB.Where("id = ?", id).First(&todo).Error
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "notfound"})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", todo)
}

func CreateTodo2(c *gin.Context) {
	desc := c.DefaultPostForm("desc", "")
	formDone := c.DefaultPostForm("done", "false")
	done, err := strconv.ParseBool(formDone)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "badparse"})
		return
	}

	todo := Todo{Desc: desc, Done: done}
	DB.Create(&todo)
	c.Redirect(http.StatusFound, "/todo")
}

func UpdateTodo2(c *gin.Context) {
	formID := c.PostForm("id")
	id, err := strconv.ParseInt(formID, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "badparse"})
		return
	}
	desc := c.PostForm("desc")
	formDone := c.DefaultPostForm("done", "false")
	log.Println("formDone:", formDone)
	done, err := strconv.ParseBool(formDone)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "badparse"})
		return
	}

	var todo Todo
	err = DB.Where("id = ?", id).First(&todo).Error
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "notfound"})
		return
	}
	if desc != "" {
		DB.Model(&todo).Updates(map[string]interface{}{"desc": desc})
	}
	if done != todo.Done {
		DB.Model(&todo).Updates(map[string]interface{}{"done": done})
	}
	c.Redirect(http.StatusFound, "/todo")
}

func DeleteTodo2(c *gin.Context) {
	formID := c.PostForm("idboi")
	id, err := strconv.ParseInt(formID, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "badparse"})
		return
	}
	var todo Todo
	err = DB.Where("id = ?", id).First(&todo).Error
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "notfound"})
		return
	}
	DB.Delete(&todo)
	c.Redirect(http.StatusFound, "/todo")
}

func main() {
	ConnectDB()
	r := gin.Default()
	// standard api routes
	r.GET("/api/todo", GetTodos)
	r.GET("/api/todo/:id", GetTodo)
	r.POST("/api/todo", CreateTodo)
	r.PATCH("/api/todo/:id", UpdateTodo)
	r.DELETE("/api/todo/:id", DeleteTodo)
	// browser api routes
	r.LoadHTMLGlob("./*.tmpl")
	r.GET("/todo", GetTodos2)
	r.GET("/todo/get", GetTodo2) // fix to get id some way
	r.POST("/todo/new", CreateTodo2)
	r.POST("/todo/edit", UpdateTodo2)
	r.POST("/todo/delete", DeleteTodo2)
	// if env var PORT is set, it will use that
	err := r.Run()
	log.Fatal("error starting HTTP server:", err)
}
