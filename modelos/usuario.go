package modelos

type Usuario struct {
	Id    int
	Nome  string
	Email string
	Senha string
}

type UsuarioGetDTO struct {
	Id    int
	Nome  string
	Email string
}

type UsuarioCreateDTO struct {
	Nome  string
	Email string
	Senha string
}
