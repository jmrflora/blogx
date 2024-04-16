package modelos

type Usuario struct {
	Id    int
	Nome  string
	Email string
	Senha string
}

type UsuarioGetDTO struct {
	Id    int `db:"idusuario"`
	Nome  string
	Email string
}

type UsuarioSenhaGetDTO struct {
	Id    int `db:"idusuario"`
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
	Nome      string `form:"nome"`
	Email     string `form:"email"`
	Senha     string `form:"senha"`
	ConfSenha string `form:"confsenha"`
}

type UsuarioLoginDTO struct {
	Email string `form:"email"`
	Senha string `form:"senha"`
}
