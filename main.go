package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Todo struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Desc string `json:"desc"`
	Done bool   `json:"done"`
}

type Todos []*Todo

var todoList Todos = []*Todo{
	&Todo{
		ID:   0,
		Desc: "finish api",
		Done: true,
	},
	&Todo{
		ID:   1,
		Desc: "add data store",
		Done: false,
	},
	&Todo{
		ID:   2,
		Desc: "add tests",
		Done: false,
	},
}

func GetTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "badparse"})
	}
	var todo Todo
	for _, todo := range todoList {
		if todo.ID == id {
			c.JSON(http.StatusOK, todo)
			return
		}
	}
	if (todo == Todo{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "notfound"})
	}
}

func GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, todoList)
}

func CreateTodo(c *gin.Context) {
	id := int64(len(todoList))
	var input Todo = Todo{ID: id, Desc: "", Done: false}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "badrequest"})
		return
	}
	todoList = append(todoList, &input)
	c.JSON(http.StatusOK, input)
}

func UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "badparse"})
	}
	for _, todo := range todoList {
		if todo.ID == id {
			break
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "badrequest"})
			return
		}
	}
	var input Todo
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notfound"})
		return
	}
	for _, todo := range todoList {
		if todo.ID == id {
			//avoid impacting the other variable if only one needs to be changed
			if input.Desc != "" {
				todoList[id].Desc = input.Desc
			}
			if input.Done != todoList[id].Done {
				todoList[id].Done = input.Done
			}
		}
	}
}

func DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "badparse"})
	}
	tobedeleted := todoList[id]
	todoList = append(todoList[:id], todoList[id+1:]...)
	c.JSON(http.StatusOK, tobedeleted)
}

func main() {
	r := gin.Default()
	r.GET("/todo", GetTodos)
	r.GET("/todo/:id", GetTodo)
	r.POST("/todo", CreateTodo)
	r.PUT("/todo/:id", UpdateTodo)
	r.PATCH("/todo/:id", UpdateTodo)
	r.DELETE("/todo/:id", DeleteTodo)
	//if env var PORT is set, it will use that
	r.Run()
}
