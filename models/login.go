package models

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RespuestaLogin struct {
	Token   string `json:"token"`
	Usuario `json:"usuario"`
}
