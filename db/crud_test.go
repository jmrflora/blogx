package db

import (
	"testing"

	"github.com/google/uuid"
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

func TestGetArtigo(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", "../mydb.db")

	tx, err := db.Beginx()
	defer tx.Rollback()

	assert.NoError(t, err)

	_, err = GetArtigoPorId(tx, "63f7e639-2d32-46e7-94a4-328217d81487")

	assert.NoError(t, err)
	// println(a.Titulo)

}

func TestCreateArtigo(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", "../mydb.db")

	tx, err := db.Beginx()
	defer tx.Rollback()

	assert.NoError(t, err)
	id := uuid.New()
	a := modelos.ArtigoCreateDTO{
		Uuid:      id.String(),
		Titulo:    "aaaaaaa",
		Subtitulo: "aaaaaaaaas",
		IdAutor:   1,
		Estrelas:  0,
	}

	_, err = CreateArtigo(tx, &a)

	assert.NoError(t, err)

}
