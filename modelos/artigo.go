package modelos

import "github.com/google/uuid"

type Artigo struct {
	Uuid      uuid.UUID
	Titulo    string
	Subtitulo string
	IdAutor   int
	Estrelas  int
}
