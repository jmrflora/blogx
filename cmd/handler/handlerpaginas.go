package handler

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/jmrflora/blogx/db"
	"github.com/jmrflora/blogx/modelos"
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

		return c.Redirect(301, "/")
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

	// a, err := db.GetArtigoPorId(tx, "2fc67b56-5b62-4341-8d95-cf90435b906c")
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// }
	//
	// categs, err := db.GetCategorias(tx)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// }
	//
	// autor, err := db.GetUsuarioPorId(tx, a.IdAutor)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// }
	//

	artigos, err := db.GetArtigoPagina(tx, 1)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	arts := []modelos.ArtigoGetDtoCatgsUsuario{}

	for _, artigo := range artigos {
		ctgs, err := db.GetCategoriasDeArtigo(tx, artigo.Uuid)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		u, err := db.GetUsuarioPorId(tx, artigo.IdAutor)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		a := modelos.ArtigoGetDtoCatgsUsuario{
			ArtigoGetDTO:  artigo,
			Categs:        ctgs,
			UsuarioGetDTO: *u,
		}

		arts = append(arts, a)

	}

	for _, art := range arts {
		println(art.Titulo)
		println("dash")
		for _, catg := range art.Categs {
			println(catg.NomeCategoria)
		}
	}
	var cmp templ.Component
	sess, _ := session.Get("session", c)
	if sess.IsNew {
		cmp = paginas.ConteudosNaoLogado(arts)
	} else {
		cmp = paginas.Conteudos(arts)
	}

	return views.Renderizar(cmp, c)
}

func (h *Handler) HandleArtigo(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	id := c.Param("id")
	tx, err := h.Dbaccess.Beginx()
	if err != nil {
		return err
	}

	art, err := db.GetArtigoPorId(tx, id)
	if err != nil {
		return err
	}

	src := "/assets/markdowns/" + strconv.Itoa(art.IdAutor) + "/" + id

	cmp := paginas.PaginaArtigo(src)
	return views.Renderizar(cmp, c)
}
