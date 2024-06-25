package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/jmrflora/blogx/db"
	"github.com/jmrflora/blogx/views"
	"github.com/jmrflora/blogx/views/paginas"
	"github.com/jmrflora/blogx/views/partials"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func HandleIndex(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		println("aquiiii")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// println(sess.IsNew)
	// println(c.Request().Header.Get("HX-Request"))
	var cmp templ.Component
	if sess.IsNew {
		cmp = paginas.Index()
	} else {
		cmp = paginas.IndexLogado()
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return views.Renderizar(cmp, c)
}

func HandlePaginaLogin(c echo.Context) error {
	var cmp templ.Component
	if c.Request().Header.Get("HX-Request") == "true" {
		cmp = partials.LoginHyper("", true)
	} else {
		cmp = paginas.PaginaLogin()
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return views.Renderizar(cmp, c)
}

func (h *Handler) HandlePaginaUpload(c echo.Context) error {
	// check for cookie
	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	println(sess.IsNew)
	var cmp templ.Component

	if sess.IsNew {
		// handle não ser logado
		cmp = paginas.Index()
	} else {

		id := sess.Values["id"].(int)

		tx, err := h.Dbaccess.Beginx()
		defer tx.Rollback()
		if err != nil {
			return err
		}
		categs, err := db.GetCategorias(tx)
		if err != nil {
			return err
		}

		if c.Request().Header.Get("HX-Request") == "true" {
			cmp = partials.Uploadartigo(id, categs)
		} else {
			cmp = paginas.PaginaUpload(id, categs)
		}
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return views.Renderizar(cmp, c)
}

func (h *Handler) HandlePaginaRegistro(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	var cmp templ.Component
	if c.Request().Header.Get("HX-Request") == "true" {
		cmp = partials.RegistroPartial("", "")
	} else {
		cmp = paginas.PaginaRegistroUsuario()
	}

	return views.Renderizar(cmp, c)
}

func (h *Handler) HandleTesteCnteudo(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	tx, err := h.Dbaccess.Beginx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()

	a, err := db.GetArtigoPorId(tx, "2fc67b56-5b62-4341-8d95-cf90435b906c")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	categs, err := db.GetCategorias(tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	autor, err := db.GetUsuarioPorId(tx, a.IdAutor)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return views.Renderizar(paginas.ConteudoTeste(*a, categs, *autor), c)
}
