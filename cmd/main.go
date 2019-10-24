package main

import (
	"database/sql"
	"log"
	"my-utils-api/internal/model"
	"my-utils-api/internal/package/db"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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
// type Todo struct {
// 	Id        int    `db:"id" json:"id"`
// 	Title     string `db:"title" json:"title"`
// 	IsDone    int    `db:"is_done" json:"isDone"`
// 	Detail    string `db:"detail" json:"detail"`
// 	CreatedAt string `db:"created_at" json:"createdAt"`
// 	UpdatedAt string `db:"updated_at" json:"updatedAt"`
// }

type SctIsDone struct {
	ID     int  `db:"id" json:"id"`
	IsDone bool `db:"is_done" json:isDone`
}

type SctDeleteTodo struct {
	ID int `db:"id" json:"id"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env.")
	}
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
	todoID := c.Param("todoId")
	// println("todoId", todoId)

	db, err := sql.Open("mysql", "root@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT * FROM todo_list WHERE id = ?", todoID)
	if err != nil {
		log.Fatal(err)
	}
	todos := []model.Todo{}
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
		todo := model.Todo{
			ID:        id,
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
	database := db.Connect()
	defer database.Close()

	todosEx := []model.Todo{}
	database.Find(&todosEx)
	return c.JSON(http.StatusOK, todosEx)
}

func addTodo(c echo.Context) (err error) {
	t := new(model.Todo)

	if err = c.Bind(t); err != nil {
		return err
	}
	title := t.Title
	detail := t.Detail

	db := db.Connect()
	defer db.Close()
	// db, err := sql.Open("mysql", "root@/todo")
	// println("addtodo")
	// if err != nil {
	// 	println("opennotoko")
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	todoEx := model.Todo{}
	todoEx.Title = title
	todoEx.Detail = detail
	todoEx.CreatedAt = "2013-11-17 21:34:10"
	// TODO 未指定似できない
	// Error 1292: Incorrect datetime value: '' for column 'updated_at' at row 1
	// time.formatうんちゃらでYYYY-MM-DD HH:MM:SSにしたいんだけどもなぁ
	todoEx.UpdatedAt = "2013-11-17 21:34:10"
	result := db.Create(&todoEx)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	// ins, err := db.Prepare("INSERT INTO todo_list (title, detail) VALUES (?, ?)")
	// if err != nil {
	// 	println("dbnotoko")
	// 	log.Fatal(err)
	// }
	// ins.Exec(title, detail)
	return getTodos(c)
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
	return getTodos(c)
}

func updateTodoByID(c echo.Context) (err error) {
	t := new(model.Todo)
	if err = c.Bind(t); err != nil {
		return err
	}
	id := t.ID
	title := t.Title
	isDone := t.IsDone
	detail := t.Detail

	db := db.Connect()
	defer db.Close()

	// 更新したいレコードを選択
	todoExBefore := model.Todo{}
	todoExBefore.ID = id
	// 更新後のレコードを生成
	todoExAfter := todoExBefore
	db.First(&todoExAfter)
	todoExAfter.Title = title
	todoExAfter.IsDone = isDone
	todoExAfter.Detail = detail
	// 更新実行
	db.Model(&todoExBefore).Update(&todoExAfter)
	// db, err := sql.Open("mysql", "root@/todo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// upd, err := db.Prepare(`
	// 	UPDATE todo_list SET title = ?, is_done = ?, detail = ? where id = ?
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// upd.Exec(title, isDone, detail, id)

	return getTodos(c)
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
	return getTodos(c)
}
