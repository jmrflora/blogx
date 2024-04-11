package modelos

type Blog struct {
	Artigo
	Categorias []Categoria
}

type BlogCreateDTO struct {
	ArtigoCreateDTO
	CategoriasIds []int
}
