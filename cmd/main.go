package main

import (
	"database/sql"
	"fmt"
	"log"
	"my-utils-api/internal/model"
	"my-utils-api/internal/package/db"
	"net/http"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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

// SctIsDone is ~~
// TODO struct名が変
type SctIsDone struct {
	ID     int  `db:"id" json:"id"`
	IsDone bool `db:"is_done" json:isDone`
}

// SctDeleteTodo is ~~
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
	e.GET("/api/todo", getTodos)
	e.GET("/api/todo/:todoId", getTodoByID)
	e.POST("/api/todo/add", addTodo)
	e.POST("/api/todo/updateIsDone", updateTodoIsDone)
	e.POST("/api/todo/deleteTodoByID", deleteTodoByID)
	e.POST("/api/todo/updateTodoByID", updateTodoByID)
}

// FIXME errorの返り値これであってるのかな
func getTodos(c echo.Context) error {
	database := db.Connect()
	defer database.Close()

	todosEx := []model.Todo{}
	// gormはtableの指定がない場合、構造体の複数形の名称でtableを参照する
	database.Find(&todosEx)
	return c.JSON(http.StatusOK, todosEx)
}

func getTodoByID(c echo.Context) error {
	// c.Paramはstringらしい
	paramID := c.Param("todoId")

	database := db.Connect()
	defer database.Close()

	todoEx := model.Todo{}
	println("paramID: ", paramID, reflect.TypeOf(paramID))
	fmt.Printf("paramID: %T\n", paramID)

	// 構造体の情報を元にテーブルを参照するらしい
	todoEx.ID = 1 // FIXME paramIDがstringだけど構造体的にはintでcastしたいんだけどできない
	result := database.First(&todoEx)
	return c.JSON(http.StatusOK, result)
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

	todoEx := model.Todo{}
	todoEx.Title = title
	todoEx.Detail = detail
	todoEx.CreatedAt = "2013-11-17 21:34:10"
	// FIXME
	// 未指定にしてMySQL側で勝手にtimestampが入ってほしいんだけど、
	// 構造体的にnilがだめだからとりあえず適当なdateを入れている
	// Error 1292: Incorrect datetime value: '' for column 'updated_at' at row 1
	// time.formatうんちゃらでYYYY-MM-DD HH:MM:SSにしたいんだけどもなぁ
	todoEx.UpdatedAt = "2013-11-17 21:34:10"
	result := db.Create(&todoEx)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return getTodos(c)
}

func deleteTodoByID(c echo.Context) (err error) {
	params := SctDeleteTodo{}

	database := db.Connect()
	defer database.Close()

	todoEx := model.Todo{}
	todoEx.ID = params.ID

	// FIXME
	// .firstをやらないとmodelにprimary_key追加しても全件削除になってしまう
	// where句使ったほうが安心っていうしそうしたほうがいいのかな
	database.First(&todoEx)
	database.Delete(&todoEx)

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
