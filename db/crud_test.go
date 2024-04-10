package db

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jmrflora/blogx/modelos"
	"github.com/stretchr/testify/assert"
)

func TestGetUsuario(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", "../mydb.db")

	tx, err := db.Beginx()
	defer tx.Rollback()

	assert.NoError(t, err)

	u, err := GetUsuarioPorId(tx, 1)

	assert.NoError(t, err)
	println(u.Nome)
}

func TestCreateUsuario(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", "../mydb.db")

	tx, err := db.Beginx()
	defer tx.Rollback()

	assert.NoError(t, err)

	u := modelos.UsuarioCreateDTO{
		Nome:  "jose teste",
		Email: "email@email",
		Senha: "1234",
	}
	_, err = CreateUsuario(tx, &u)

	assert.NoError(t, err)

}
