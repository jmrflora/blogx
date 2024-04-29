package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/jmrflora/blogx/db"
	"github.com/jmrflora/blogx/modelos"
	"github.com/jmrflora/blogx/views"
	"github.com/jmrflora/blogx/views/partials"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Dbaccess sqlx.DB
}

func (h *Handler) HandleUpload(c echo.Context) error {

	var blog modelos.BlogRegistroDTO

	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = echo.FormFieldBinder(c).
		String("Titulo", &blog.Titulo).
		String("Subtitulo", &blog.Subtitulo).
		Ints("CategoriasIds", &blog.CategoriasIds).
		BindError()
	if err != nil {
		return err
	}

	blog.IdAutor = sess.Values["id"].(int)

	tx, err := h.Dbaccess.Beginx()
	defer tx.Rollback()
	if err != nil {
		return err
	}

	b := modelos.BlogCreateDTO{
		ArtigoCreateDTO: modelos.ArtigoCreateDTO{
			Uuid:      uuid.NewString(),
			Titulo:    blog.Titulo,
			Subtitulo: blog.Subtitulo,
			IdAutor:   blog.IdAutor,
			Estrelas:  0,
		},
		CategoriasIds: blog.CategoriasIds,
	}

	_, err = db.CreateBlog(tx, &b)
	if err != nil {
		println("aquiiiiiiii")
		return err
	}

	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination

	// Destination directory
	uploadDir := "internal/assets/markdowns/usuariologado"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		println("ola")
		// Create the directory if it doesn't exist
		if err := os.Mkdir(uploadDir, 0777); err != nil {
			return err
		}
	}
	dstPath := filepath.Join(uploadDir, b.Uuid)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	tx.Commit()
	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully</p>", file.Filename))
}

func (h *Handler) HandleLogin(c echo.Context) error {
	u := modelos.UsuarioLoginDTO{}

	if err := c.Bind(&u); err != nil {
		println("aqui 1")
		return views.Renderizar(partials.LoginHyper("", false), c)

	}

	//adicionar validação

	tx, err := h.Dbaccess.Beginx()
	defer tx.Rollback()
	if err != nil {
		return echo.ErrInternalServerError
	}

	usuarioDb, err := db.GetUsuarioComSenhaPorEmail(tx, u.Email)
	if err != nil {
		println("aqui 2")
		return views.Renderizar(partials.LoginHyper(u.Email, false), c)
	}
	if ok := CheckPasswordHash(u.Senha, usuarioDb.Senha); !ok {
		println("aqui 3")
		return views.Renderizar(partials.LoginHyper(u.Email, false), c)
	}
	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
	}

	sess.Values["email"] = usuarioDb.Email
	sess.Values["nome"] = usuarioDb.Nome
	sess.Values["id"] = usuarioDb.Id
	sess.Values["permissao"] = "normal"

	sess.Save(c.Request(), c.Response())
	// precisa ser assim pro js do bootstrap não bugar
	// c.Response().Header().Add("HX-Push-Url", "/")
	// c.Response().Header().Add("HX-Refresh", "true")
	c.Response().Header().Add("HX-Redirect", "/")
	// return c.Redirect(http.StatusSeeOther, "/")
	return echo.NewHTTPError(http.StatusOK, "ok")
}

func (h *Handler) HandleRegistroUsuario(c echo.Context) error {
	u := modelos.UsuarioRegistroDTO{}

	if err := c.Bind(&u); err != nil {
		return echo.ErrBadRequest
	}

	//adicionar validação

	if u.Senha != u.ConfSenha {
		return echo.ErrBadRequest
	}

	novaSenha, err := HashPassword(u.Senha)
	if err != nil {
		return echo.ErrInternalServerError
	}
	novoUsuario := modelos.UsuarioCreateDTO{
		Nome:  u.Nome,
		Email: u.Email,
		Senha: novaSenha,
	}

	tx, err := h.Dbaccess.Beginx()
	defer tx.Rollback()
	if err != nil {
		return echo.ErrInternalServerError
	}

	_, err = db.CreateUsuario(tx, &novoUsuario)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()

	return echo.NewHTTPError(http.StatusOK)
}

func (h *Handler) HandleRegistroUsuarioConfSenha(c echo.Context) error {
	u := modelos.UsuarioRegistroDTO{}

	if err := c.Bind(&u); err != nil {
		return views.Renderizar(partials.DivComSenhaDiferenteErro(), c)
	}
	if u.Senha != u.ConfSenha {
		return views.Renderizar(partials.DivComSenhaDiferenteErro(), c)
	}

	return views.Renderizar(partials.DivComSenha(u.ConfSenha), c)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
