package main

import (
	"crypto/rand"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/jmrflora/blogx/cmd/handler"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// db, err := sqlx.Connect("sqlite3", "../mydb.db")
	// if err != nil {
	// 	println(err.Error())
	// }

	id := uuid.New()

	println(id.String())

	db := sqlx.MustConnect("sqlite3", "mydb.db")

	h := handler.Handler{
		Dbaccess: *db,
	}
	key := make([]byte, 16)
	rand.Read(key)
	println(key)

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"), nil)))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use()

	e.Static("/assets", "internal/assets")

	e.GET("/", handler.HandleIndex)

	e.GET("/paginalogin", handler.HandlePaginaLogin)

	e.POST("/login", h.HandleLogin)

	e.GET("/paginaregistro", h.HandlePaginaRegistro)

	e.POST("/registro", h.HandleRegistroUsuario)

	e.POST("/registro/confsenha", h.HandleRegistroUsuarioConfSenha)

	p := e.Group("", CookieAuthMiddleware)

	p.POST("/upload", h.HandleUpload)

	p.POST("/upload/parse", h.HandleUploadParse)

	p.GET("/paginaupload", h.HandlePaginaUpload)

	println("olaaaaaa")

	e.Logger.Fatal(e.Start(":1323"))
}

func CookieAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return echo.ErrUnauthorized
		}
		if sess.IsNew {
			return echo.ErrUnauthorized
		}
		if sess.Values["permissao"].(string) != "normal" {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
