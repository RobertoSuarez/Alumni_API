package models

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RespuestaLogin struct {
	Token   string `json:"token"`
	Usuario `json:"usuario"`
}
