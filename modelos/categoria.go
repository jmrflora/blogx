package modelos

type Categoria struct {
	Id            int    `db:"idcateg"`
	NomeCategoria string `db:"nomecateg"`
}

type CategoriaCreateDTO struct {
	NomeCategoria string
}
