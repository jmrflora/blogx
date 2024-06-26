package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/jmrflora/blogx/modelos"
	_ "github.com/mattn/go-sqlite3"
)

func GetUsuarioPorId(tx *sqlx.Tx, id int) (*modelos.UsuarioGetDTO, error) {
	u := modelos.UsuarioGetDTO{}

	err := tx.Get(&u, "select usuario.idusuario, usuario.nome, usuario.email from usuario where usuario.idusuario = $1 limit 1", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUsuarioComSenhaPorEmail(tx *sqlx.Tx, mail string) (*modelos.UsuarioSenhaGetDTO, error) {
	u := modelos.UsuarioSenhaGetDTO{}

	err := tx.Get(&u, "select usuario.idusuario, usuario.nome, usuario.email, usuario.senha from usuario where usuario.email = $1 limit 1", mail)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	return &u, nil
}

func CreateUsuario(tx *sqlx.Tx, u *modelos.UsuarioCreateDTO) (sql.Result, error) {
	nstmt, err := tx.PrepareNamed("INSERT INTO usuario (nome, email, senha) values (:nome, :email, :senha)")
	if err != nil {
		return nil, err
	}
	result, err := nstmt.Exec(u)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetArtigoPorId(tx *sqlx.Tx, id string) (*modelos.ArtigoGetDTO, error) {
	a := modelos.ArtigoGetDTO{}

	err := tx.Get(&a, "SELECT artigo.uuid, artigo.titulo, artigo.subtitulo, artigo.idautor, artigo.estrelas from artigo where artigo.uuid = $1 limit 1", id)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func GetArtigoPagina(tx *sqlx.Tx, pag int) ([]modelos.ArtigoGetDTO, error) {
	limit := 3
	offset := (pag - 1) * limit

	artgs := []modelos.ArtigoGetDTO{}

	err := tx.
		Select(&artgs, "SELECT artigo.uuid, artigo.titulo, artigo.subtitulo, artigo.idautor, artigo.estrelas from artigo limit $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}

	return artgs, nil
}

func CreateArtigo(tx *sqlx.Tx, a *modelos.ArtigoCreateDTO) (sql.Result, error) {
	nstmt, err := tx.PrepareNamed("INSERT INTO artigo (uuid, titulo, subtitulo, idautor, estrelas) values (:uuid, :titulo, :subtitulo, :idautor, 0)")
	if err != nil {
		println("primeiro erro")
		return nil, err
	}

	result, err := nstmt.Exec(a)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CreateCategoria(tx *sqlx.Tx, c modelos.CategoriaCreateDTO) (sql.Result, error) {
	nstmt, err := tx.PrepareNamed("INSERT INTO categoria (nomecateg) VALUES (:nomecategoria)")
	if err != nil {
		return nil, err
	}
	result, err := nstmt.Exec(c)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetCategorias(tx *sqlx.Tx) ([]modelos.Categoria, error) {
	mc := []modelos.Categoria{}

	err := tx.Select(&mc, "SELECT * FROM categoria")
	if err != nil {
		return nil, err
	}
	return mc, err
}

func GetCategoriasDeArtigo(tx *sqlx.Tx, id string) ([]modelos.Categoria, error) {
	mc := []modelos.Categoria{}

	err := tx.
		Select(&mc, "SELECT categoria.idcateg, categoria.nomecateg from categoria inner join categoriasdeartigo c on c.idcategoria = categoria.idcateg inner join artigo on c.idartigo = artigo.uuid WHERE artigo.uuid = $1", id)
	if err != nil {
		return nil, err
	}
	return mc, nil
}

func CreateCategoriasDeArtigo(tx *sqlx.Tx, idArtigo string, idCategoria int) (sql.Result, error) {
	result, err := tx.Exec("INSERT INTO categoriasdeartigo (idartigo, idcategoria) VALUES ($1, $2)", idArtigo, idCategoria)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CategorizarArtigo(tx *sqlx.Tx, idArtigo string, idsCategoria []int) (sql.Result, error) {
	var result sql.Result
	var err error
	for _, idCategoria := range idsCategoria {
		result, err = CreateCategoriasDeArtigo(tx, idArtigo, idCategoria)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func CreateBlog(tx *sqlx.Tx, b *modelos.BlogCreateDTO) (sql.Result, error) {
	_, err := CreateArtigo(tx, &b.ArtigoCreateDTO)
	if err != nil {
		println("aqui 2")
		return nil, err
	}
	result, err := CategorizarArtigo(tx, b.Uuid, b.CategoriasIds)
	if err != nil {
		return nil, err
	}
	return result, nil
}
