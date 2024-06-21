package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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
	println(h.Dbaccess.Stats().InUse)
	var blog modelos.BlogRegistroDTO

	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// criar sub handler para fazer arse parcial do md para descobrir titulo e Subtitulo
	var texto string
	err = echo.FormFieldBinder(c).
		String("texto", &texto).BindError()
	if err != nil {
		return views.Renderizar(partials.ModalComErro(err), c)
	}
	fmt.Printf("%q\n", texto)

	err, titulo, sub := Parse(texto)
	if err != nil {
		return views.Renderizar(partials.ModalComErro(err), c)
	}

	blog.Titulo = titulo
	blog.Subtitulo = sub

	err = echo.FormFieldBinder(c).
		Ints("CategoriasIds", &blog.CategoriasIds).
		BindError()
	if err != nil {
		println("aqui aaaaaaaa" + err.Error())

		return views.Renderizar(partials.ModalComErro(err), c)
	}

	blog.IdAutor = sess.Values["id"].(int)

	err = h.Dbaccess.Ping()
	if err != nil {
		println(err.Error())
	}

	tx, err := h.Dbaccess.Beginx()
	defer tx.Rollback()
	if err != nil {
		return views.Renderizar(partials.ModalComErro(err), c)
	}

	b := modelos.BlogCreateDTO{
		ArtigoCreateDTO: modelos.ArtigoCreateDTO{
			Uuid:      uuid.NewString(),
			Titulo:    blog.Titulo,
			Subtitulo: blog.Subtitulo,
			IdAutor:   blog.IdAutor,
		},
		CategoriasIds: blog.CategoriasIds,
	}

	fmt.Printf("b: %v\n", b)

	// _, err = db.CreateBlog(tx, &b)
	// if err != nil {
	// 	println(err.Error())
	// 	println("erro aqui aa erro")
	// 	return views.Renderizar(partials.ModalComErro(err), c)
	// }
	//
	_, err = db.CreateArtigo(tx, &b.ArtigoCreateDTO)
	if err != nil {
		println(err.Error())
	}
	println(h.Dbaccess.Stats().MaxOpenConnections)
	println(h.Dbaccess.Stats().InUse)
	err = tx.Commit()
	println(err.Error())
	//-----------
	// Read file
	//-----------

	// Source

	// Destination

	// Destination directory

	uploadDir := "internal/assets/markdowns/" + strconv.Itoa(sess.Values["id"].(int))

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		println("ola")
		// Create the directory if it doesn't exist
		if err := os.Mkdir(uploadDir, 0777); err != nil {
			return views.Renderizar(partials.ModalComErro(err), c)
		}
	}
	dstPath := filepath.Join(uploadDir, b.Uuid)
	dst, err := os.Create(dstPath)
	if err != nil {
		return views.Renderizar(partials.ModalComErro(err), c)
	}
	defer dst.Close()

	// Copy

	_, err = dst.WriteString(texto)
	if err != nil {
		println(err.Error())

		return views.Renderizar(partials.ModalComErro(err), c)
	}

	tx.Commit()
	println("final do upload")
	c.Response().Header().Add("HX-Trigger", "myEvent")

	return views.Renderizar(partials.Modal(), c)
}

func (h *Handler) HandleUploadParse(c echo.Context) error {
	var texto string
	err := echo.FormFieldBinder(c).
		String("texto", &texto).BindError()
	if err != nil {
		return echo.ErrBadRequest
	}
	fmt.Printf("%q\n", texto)

	regex, err := regexp.Compile(`^# .*(?:\r?\n){1,2}## .*`)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if !regex.MatchString(texto) {
		return echo.NewHTTPError(http.StatusBadRequest, "coloque um titulo e um sub")
	}

	return echo.NewHTTPError(http.StatusOK, texto)
}

func (h *Handler) HandleLogin(c echo.Context) error {
	u := modelos.UsuarioLoginDTO{}

	if err := c.Bind(&u); err != nil {
		println("aqui 1")
		return views.Renderizar(partials.LoginHyper("", false), c)

	}

	// adicionar validação

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
		Secure:   false,
	}

	sess.Values["email"] = usuarioDb.Email
	sess.Values["nome"] = usuarioDb.Nome
	sess.Values["id"] = usuarioDb.Id
	sess.Values["permissao"] = "normal"

	sess.Save(c.Request(), c.Response())
	// precisa ser assim pro js do bootstrap não bugar
	// c.Response().Header().Add("HX-Push-Url", "/")
	// c.Response().Header().Add("HX-Refresh", "true")
	// return c.Redirect(http.StatusSeeOther, "/")
	c.Response().Header().Add("HX-Redirect", "/")
	return echo.NewHTTPError(http.StatusOK, "ok")
}

func (h *Handler) HandleRegistroUsuario(c echo.Context) error {
	u := modelos.UsuarioRegistroDTO{}

	if err := c.Bind(&u); err != nil {
		return echo.ErrBadRequest
	}

	// adicionar validação

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

func Parse(texto string) (error, string, string) {
	regex := regexp.MustCompile(`^# .*(?:\r?\n){1,2}## .*`)

	if !regex.MatchString(texto) {
		return errors.New("no match"), "", ""
	}

	fds := regex.FindString(texto)
	println(fds)
	linhas := splitLines(fds)

	fmt.Printf("linhas[0]: %v\n", linhas[0])
	fmt.Printf("linhas[1]: %v\n", linhas[1])

	return nil, linhas[0], linhas[1]
}

func splitLines(s string) []string {
	// Split the string into lines using FieldsFunc
	return strings.FieldsFunc(s, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
}
