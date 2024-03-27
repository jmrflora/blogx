package main

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// db, err := sqlx.Connect("sqlite3", "../mydb.db")
	// if err != nil {
	// 	println(err.Error())
	// }
	db := sqlx.MustConnect("sqlite3", "../mydb.db")

	db.Ping()

	// exactly the same as the built-in

	u := uuid.New()
	println(u.String())
	e := echo.New()

	e.Logger.Fatal(e.Start(":1323"))

}
