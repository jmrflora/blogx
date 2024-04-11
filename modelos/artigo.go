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

type ArtigoCreateDTO struct {
	Uuid      string
	Titulo    string
	Subtitulo string
	IdAutor   int
	Estrelas  int
}

type ArtigoRegistroDTO struct {
	Titulo    string
	Subtitulo string
	IdAutor   int
	Estrelas  int
}
