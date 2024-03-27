package modelos

type Blog struct {
	Artigo
	UsuarioGetDTO
	Categorias []Categoria
}
