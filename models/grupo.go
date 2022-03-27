package models

type Grupo struct {
	ID       uint64    `json:"id" gorm:"primary_key"`
	Name     string    `json:"name" gorm:"not null;unique"`
	Usuarios []Usuario `json:"usuarios,omitempty" gorm:"many2many:usuario_grupos;"`
}

func (Grupo) TableName() string {
	return "grupo"
}

// Crear registra un grupo en la base de datos
// el grupo creado se llenara con el todo los datos y el id
// como el objeto es un puntero se llena automaticamente
func (g *Grupo) Crear() error {
	tx := DB.Begin()

	if err := tx.Create(&g).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Lista grupos
func (g Grupo) ObtenerGrupos() (grupos []Grupo, err error) {
	grupos = []Grupo{}

	fieldsStr := "id,name"
	if err = DB.Select(fieldsStr).Find(&grupos).Error; err != nil {
		return grupos, err
	}

	return grupos, nil
}

// Agregar usuario a un grupo
// se pasa el id de usuario y se agrega al grupo
func (g *Grupo) AgregarUsuario(u Usuario) error {
	tx := DB.Begin()

	// En la tabla usuario_grupos se agregara esta asociación
	err := tx.Model(&g).Association("Usuarios").Append(u)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
