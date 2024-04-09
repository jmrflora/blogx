package main

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jmrflora/blogx/cmd/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// db, err := sqlx.Connect("sqlite3", "../mydb.db")
	// if err != nil {
	// 	println(err.Error())
	// }
	db := sqlx.MustConnect("sqlite3", "../mydb.db")

	db.Ping()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/assets", "internal/assets")

	e.POST("/upload", handler.HandleUpload)

	e.GET("/", handler.HandleIndex)
	e.GET("/paginalogin", handler.HandlePaginaLogin)

	e.POST("/login", func(c echo.Context) error {
		println("i tried")
		return c.String(http.StatusOK, "there was an attempt")
	})

	e.Logger.Fatal(e.Start(":1323"))

}
