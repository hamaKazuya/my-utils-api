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

// Todo struct
type Todo struct {
	Id        int    `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	IsDone    int    `db:"is_done" json:"isDone"`
	Detail    string `db:"detail" json:"detail"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
}

type SctIsDone struct {
	ID     int  `db:"id" json:"id"`
	IsDone bool `db:"is_done" json:isDone`
}

type SctDeleteTodo struct {
	ID int `db:"id" json:"id"`
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
	e.POST("/api/todo/deleteTodoByID", deleteTodoByID)
	e.POST("/api/todo/updateTodoByID", updateTodoByID)
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
	db, err := sql.Open("mysql", "root@/todo")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	println("open ok.")
	var id int
	var name string
	rows, err := db.Query("SELECT * FROM members")
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

func getTodoByID(c echo.Context) error {
	todoId := c.Param("todoId")
	// println("todoId", todoId)

	db, err := sql.Open("mysql", "root@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT * FROM todo_list WHERE id = ?", todoId)
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
	db, err := sql.Open("mysql", "root@/todo")
	// db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/GoLife")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM todo_list")
	if err != nil {
		println("getTodos 2")
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

func addTodo(c echo.Context) (err error) {
	t := new(Todo)

	if err = c.Bind(t); err != nil {
		println("bindnotoko")
		return err
	}
	title := t.Title
	detail := t.Detail

	db, err := sql.Open("mysql", "root@/todo")
	println("addtodo")
	if err != nil {
		println("opennotoko")
		log.Fatal(err)
	}
	defer db.Close()
	ins, err := db.Prepare("INSERT INTO todo_list (title, detail) VALUES (?, ?)")
	if err != nil {
		println("dbnotoko")
		log.Fatal(err)
	}
	ins.Exec(title, detail)

	// res := getTodos(c)
	return c.JSON(http.StatusOK, true)
}

func deleteTodoByID(c echo.Context) (err error) {
	println("start deleteTodoByID")
	t := new(SctDeleteTodo)
	if err = c.Bind(t); err != nil {
		return err
	}
	todoID := t.ID

	db, err := sql.Open("mysql", "root@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dlt, err := db.Prepare("DELETE FROM todo_list WHERE id = ?;")
	dlt.Exec(todoID)

	res := getTodos(c)
	println(res)
	return c.JSON(http.StatusOK, res)
}

func updateTodoByID(c echo.Context) (err error) {
	t := new(Todo)
	if err = c.Bind(t); err != nil {
		return err
	}
	id := t.Id
	title := t.Title
	isDone := t.IsDone
	detail := t.Detail

	db, err := sql.Open("mysql", "root@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	upd, err := db.Prepare(`
		UPDATE todo_list SET title = ?, is_done = ?, detail = ? where id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	upd.Exec(title, isDone, detail, id)

	return c.JSON(http.StatusOK, "success")
}

func updateTodoIsDone(c echo.Context) (err error) {
	s := new(SctIsDone)
	if err = c.Bind(s); err != nil {
		println(s)
		return err
	}
	todoID := s.ID
	isDone := s.IsDone

	db, err := sql.Open("mysql", "root@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	upd, err := db.Prepare("UPDATE todo_list SET is_done = ? where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	upd.Exec(isDone, todoID)

	return c.JSON(http.StatusOK, "success")
}
