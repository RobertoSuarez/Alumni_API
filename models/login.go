package models

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PerfilUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RespuestaLogin struct {
	Token  string     `json:"token"`
	Perfil PerfilUser `json:"perfil"`
}
