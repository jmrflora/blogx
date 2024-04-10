package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/jmrflora/blogx/modelos"
	_ "github.com/mattn/go-sqlite3"
)

func GetUsuarioPorId(tx *sqlx.Tx, id int) (*modelos.UsuarioGetDTO, error) {
	u := modelos.UsuarioGetDTO{}

	err := tx.Get(&u, "select usuario.nome, usuario.email from usuario where usuario.idusuario = $1 limit 1", id)
	if err != nil {
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
