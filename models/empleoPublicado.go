package models

// esto sera como una vista, en mvc
// enviara el json de un empleo publicado, que es diferente a un
// empleo simple como lo definimos en el modelo.

type EmpleoPublicadoAPI struct {
	// datos del empleo
	Empleo struct {
		ID     uint64
		Titulo string
		Activo *bool
	}

	// cosas que no estan en la tabla
	Solicitudes int64  `json:"solicitudes"`
	ImgEmpresa  string `json:"imgEmpresa"`
}
