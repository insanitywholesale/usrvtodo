package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	//"log"
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
	&Todo{
		ID:   3,
		Desc: "host on todo.distro.watch",
		Done: true,
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
	//fix logical problem here
	//only works for id==0
	//for _, todo := range todoList {
	//	if todo.ID == id {
	//		break
	//	} else {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "badrequest"})
	//		return
	//	}
	//}
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

//for html frontend

func GetTodos2(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", todoList)
}

func GetTodo2(c *gin.Context) {//actually implement this
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "badparse"})//make it return html
	}
	var todo Todo
	for _, todo := range todoList {
		if todo.ID == id {
			c.HTML(http.StatusOK, "index.tmpl", []*Todo{todo})
			return
		}
	}
	if (todo == Todo{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "notfound"})//make it return html
	}
}

func CreateTodo2(c *gin.Context) {//adjust to use form
	id := int64(len(todoList))
	desc := c.DefaultPostForm("desc", "")
	formDone := c.DefaultPostForm("done", "false")
	done, err := strconv.ParseBool(formDone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "badparse"})//make it return html
		return
	}
	var input Todo = Todo{ID: id, Desc: desc, Done: done}
	todoList = append(todoList, &input)
	c.Redirect(http.StatusFound, "/todo")
}

func UpdateTodo2(c *gin.Context) {//adjust to use form
	//fix logical problem here
	//only works for id==0
	//for _, todo := range todoList {
	//	if todo.ID == id {
	//		break
	//	} else {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "badrequest"})//make it return html
	//		return
	//	}
	//}
	formID := c.PostForm("id")
	id, err := strconv.ParseInt(formID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "badparse"})//make it return html
	}
	desc := c.DefaultPostForm("desc", "")
	formDone := c.DefaultPostForm("done", "false")
	done, err := strconv.ParseBool(formDone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "badparse"})//make it return html
		return
	}
	for _, todo := range todoList {
		if todo.ID == id {
			//avoid impacting the other variable if only one needs to be changed
			if desc != "" {
				todoList[id].Desc = desc
				c.Redirect(http.StatusFound, "/todo")
			}
			if done != todoList[id].Done {
				todoList[id].Done = done
				c.Redirect(http.StatusFound, "/todo")

			}
		}
	}
}

func DeleteTodo2(c *gin.Context) {//adjust to use form
	formID := c.PostForm("id")
	id, err := strconv.ParseInt(formID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "badparse"})//make it return html
	}
	todoList = append(todoList[:id], todoList[id+1:]...)
	c.Redirect(http.StatusFound, "/todo")
}

func main() {
	r := gin.Default()
	//standard api routes
	r.GET("/api/todo", GetTodos)
	r.GET("/api/todo/:id", GetTodo)
	r.POST("/api/todo", CreateTodo)
	r.PATCH("/api/todo/:id", UpdateTodo)
	r.DELETE("/api/todo/:id", DeleteTodo)
	//browser api routes
	r.LoadHTMLGlob("./index.tmpl")
	r.GET("/todo", GetTodos2)
	r.GET("/todo/get", GetTodo2) //fix to get id some way
	r.POST("/todo/new", CreateTodo2)
	r.POST("/todo/edit", UpdateTodo2)
	r.POST("/todo/delete", DeleteTodo2)
	//if env var PORT is set, it will use that
	r.Run()
}
