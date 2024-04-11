package modelos

type Categoria struct {
	Id            int
	NomeCategoria string
}

type CategoriaCreateDTO struct {
	NomeCategoria string
}
