package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// User aaa
type User struct {
	ID      int    `json:"id"`
	GroupID int    `json:"group_id"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
}

// Member struct
type Member struct {
	ID   int
	Name string
}

func main() {
	e := echo.New()
	initRouting(e)
	e.Logger.Fatal(e.Start(":1313"))
}

func initRouting(e *echo.Echo) {
	e.GET("/api/v1/member", getMember)
	e.GET("/api/todo", getTodos)
	e.GET("/api/todo/:todoId", getTodoByID)
	e.POST("/api/todo/add", addTodo)
	e.POST("/api/todo/updateIsDone", updateTodoIsDone)
}

func getUsers(c echo.Context) error {
	groupIDstr := c.Param("group_id")
	groupID, err := strconv.Atoi(groupIDstr)
	if err != nil {
		return errors.Wrapf(err, "errors when group id convert to int: %s", groupIDstr)
	}
	gender := c.QueryParam("gender")
	users := []*User{}
	if gender == "" || gender == "man" {
		users = append(users, &User{ID: 1, GroupID: groupID, Name: "Taro", Gender: "man"})
		users = append(users, &User{ID: 2, GroupID: groupID, Name: "Jiro", Gender: "man"})
	}
	if gender == "" || gender == "woman" {
		users = append(users, &User{ID: 3, GroupID: groupID, Name: "hanako", Gender: "woman"})
		users = append(users, &User{ID: 4, GroupID: groupID, Name: "Yoshiko", Gender: "woman"})
	}
	return c.JSON(http.StatusOK, users)
}

func getMember(c echo.Context) error {
	db, err := sql.Open("mysql", "root:waiting2@/todo")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	println("open ok.")
	var id int
	var name string
	rows, err := db.Query("SELECT * FROM member")
	if err != nil {
		log.Fatal(err)
	}

	println(rows)

	members := []*Member{}
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		println(id, name)
		members = append(members, &Member{
			ID:   id,
			Name: name,
		})
	}
	return c.JSON(http.StatusOK, members)
}

// Todo struct
type Todo struct {
	Id        int
	Title     string
	IsDone    int
	Detail    string
	CreatedAt string
	UpdatedAt string
}

type SctIsDone struct {
	id     int
	isDone bool
}

func getTodoByID(c echo.Context) error {
	todoId := c.Param("todoId")
	// println("todoId", todoId)

	db, err := sql.Open("mysql", "root:waiting2@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT * FROM todos WHERE id = ?", todoId)
	if err != nil {
		log.Fatal(err)
	}
	todos := []Todo{}
	var (
		id        int
		title     string
		isDone    int
		detail    string
		createdAt string
		updatedAt string
	)
	for row.Next() {
		err = row.Scan(
			&id,
			&title,
			&isDone,
			&detail,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}
		todo := Todo{
			Id:        id,
			Title:     title,
			IsDone:    isDone,
			Detail:    detail,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		todos = append(todos, todo)
	}
	return c.JSON(http.StatusOK, todos)
}

// FIXME errorの返り値これであってるのかな
func getTodos(c echo.Context) error {
	// TODO 環境ごとに.envに持たせる(localなのでこれは現状大丈夫かなと・・)
	db, err := sql.Open("mysql", "root:waiting2@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	todos := []Todo{}
	var (
		id        int
		title     string
		isDone    int
		detail    string
		createdAt string
		updatedAt string
	)
	for rows.Next() {
		err := rows.Scan(
			&id,
			&title,
			&isDone,
			&detail,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}
		todo := Todo{
			Id:        id,
			Title:     title,
			IsDone:    isDone,
			Detail:    detail,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		todos = append(todos, todo)
	}
	return c.JSON(http.StatusOK, todos)
}

func addTodo(c echo.Context) error {
	db, err := sql.Open("mysql", "root:waiting2@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ins, err := db.Prepare("INSERT INTO todos (title, isDone, detail) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	ins.Exec("a", 0, "b")

	res := getTodos(c)
	return c.JSON(http.StatusOK, res)
}

func updateTodoIsDone(c echo.Context) error {
	todoID := c.Param("id")
	isDone := c.Param("isDone")

	db, err := sql.Open("mysql", "root:waiting2@/todo")
	if err != nil {
		log.Fatal(err)
	}
	println("todoID: ", todoID, "isDone: ", isDone)
	defer db.Close()
	upd, err := db.Prepare("UPDATE todos SET isDone = ? where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	upd.Exec(isDone, todoID)
	return c.JSON(http.StatusOK, "success")
}
