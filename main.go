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
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

func main() {

	// rows, err := db.Query("select * from users where id < ?", 6)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	err := rows.Scan(&id, &name)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(id, name)
	// }

	// rows, err := db.Query("SELECT id, name FROM users WHERE id = ?", 1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	err := rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	e := echo.New()
	initRouting(e)
	e.Logger.Fatal(e.Start(":1313"))
}

func initRouting(e *echo.Echo) {
	e.GET("/", hello)
	e.GET("/api/v1/member", getMember)
	e.GET("/api/v1/todo", getTodos)
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

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"hello": "world"})
}

func getMember(c echo.Context) error {
	// db, err := sql.Open("mysql", "admin:SrNOTeV3vJa9gY92kpbt@tcp(db1.cnkjsudgk7pu.ap-northeast-1.rds.amazonaws.com:3306)/todo")
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

// Todos struct
type Todo struct {
	id         int    `ison:id`
	title      string `json:title`
	isDone     int    `json:isDone`
	detail     string `json:detail`
	created_at string `json:created_at`
	updated_at string `json:updated_at`
}

func getTodos(c echo.Context) error {
	db, err := sql.Open("mysql", "root:waiting2@/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var (
		id         int
		title      string
		isDone     int
		detail     string
		created_at string
		updated_at string
	)
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	todos := []*Todo{}
	for rows.Next() {
		err := rows.Scan(
			&id,
			&title,
			&isDone,
			&detail,
			&created_at,
			&updated_at,
		)
		if err != nil {
			log.Fatal(err)
		}
		println(id, title, isDone, detail, created_at, updated_at)
		todos = append(todos, &Todo{
			id:         id,
			title:      title,
			isDone:     isDone,
			detail:     detail,
			created_at: created_at,
			updated_at: updated_at,
		})
		println(todos)
	}
	return c.JSON(http.StatusOK, todos)
}
