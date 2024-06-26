package modelos

type Artigo struct {
	Uuid      string
	Titulo    string
	Subtitulo string
	IdAutor   int
	Estrelas  int
}

type ArtigoGetDTO struct {
	Uuid      string
	Titulo    string
	Subtitulo string
	IdAutor   int
	Estrelas  int
}

type ArtigoGetDtoCatgs struct {
	categs []Categoria
	ArtigoGetDTO
}

type ArtigoGetDtoCatgsUsuario struct {
	UsuarioGetDTO
	Categs []Categoria
	ArtigoGetDTO
}

type ArtigoCreateDTO struct {
	Uuid      string
	Titulo    string
	Subtitulo string
	IdAutor   int
}

type ArtigoRegistroDTO struct {
	Titulo    string
	Subtitulo string
	IdAutor   int
}
