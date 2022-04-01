package models

type ErrorAPI struct {
	Mensaje string `json:"mensaje"`
}

func NewError(msg string) ErrorAPI {
	return ErrorAPI{
		Mensaje: msg,
	}
}
