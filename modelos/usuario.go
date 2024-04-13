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

type UsuarioSenhaGetDTO struct {
	Id    int
	Nome  string
	Email string
	Senha string
}

type UsuarioCreateDTO struct {
	Nome  string
	Email string
	Senha string
}

type UsuarioRegistroDTO struct {
	Nome      string
	Email     string
	Senha     string
	ConfSenha string
}

type UsuarioLoginDTO struct {
	Email string
	Senha string
}
