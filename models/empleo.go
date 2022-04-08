package models

import (
	"errors"
	"time"
)

// Esto es el modelo de empleo, el cual son las ofertas de empleos
// que las empresas publican

type Empleo struct {
	ID                       uint64 `json:"id" gorm:"primary_key"`
	CreatedAt                time.Time
	UpdatedAt                time.Time
	Titulo                   string    `json:"titulo" gorm:"size:200"`
	Descripcion              string    `json:"descripcion"`
	Profesion                string    `json:"profesion" gorm:"size:200"`
	Puesto                   string    `json:"puesto" gorm:"size:200"`
	TipoEmplo                string    `json:"tipoEmpleo" gorm:"size:200"` //Modalidad de trabajo
	SubareaID                uint64    `json:"subareaid"`
	Subarea                  Subarea   `json:"subarea" gorm:"foreignKey:SubareaID"`
	Sueldo                   string    `json:"sueldo" gorm:"size:200"`
	TiempoExperiencia        string    `json:"tiempoExperiencia" gorm:"size:200"` // Los años de experiencia
	Jornada                  string    `json:"jornada" gorm:"size:200"`
	TipoContrato             string    `json:"tipoContrato" gorm:"size:200"`
	ConocimientosAdicionales string    `json:"conocimientosAdicionales" gorm:"size:200"`
	Ciudad                   string    `json:"ciudad" gorm:"size:200"`
	PostulanteDiscapacidad   *bool     `json:"postulanteDiscapacidad" gorm:"default:false"` // si el trabajo es para personas con capacidades limitadas.
	Publicado                time.Time `json:"publicado"`
	Borrador                 *bool     `json:"borrador" gorm:"default:false"`
	EmpresaID                uint64    `json:"empresaid"`
	//Area                     string    `json:"area" gorm:"size:200"`       // Categoria
}

func (Empleo) TableName() string {
	return "empleo"
}

// Hacer: Publicar empleo
func (e *Empleo) Crear() error {
	if e.EmpresaID < 1 {
		return errors.New("fatal no existe el id de la empresa")
	}
	e.Publicado = time.Now()

	result := DB.Create(&e)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Actulizar empleo
func (e *Empleo) Actualizar() error {
	tx := DB.Begin()

	result := tx.Model(&e).Updates(Empleo{
		Titulo:                   e.Titulo,
		Descripcion:              e.Descripcion,
		Profesion:                e.Profesion,
		Puesto:                   e.Puesto,
		TipoEmplo:                e.TipoEmplo,
		SubareaID:                e.SubareaID,
		Sueldo:                   e.Sueldo,
		TiempoExperiencia:        e.TiempoExperiencia,
		Jornada:                  e.Jornada,
		TipoContrato:             e.TipoContrato,
		ConocimientosAdicionales: e.ConocimientosAdicionales,
		Ciudad:                   e.Ciudad,
		PostulanteDiscapacidad:   e.PostulanteDiscapacidad,
		Borrador:                 e.Borrador,
		//Area:                     e.Area,
	})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Preload("Subarea.Area").First(&e, e.ID)

	tx.Commit()
	return nil
}

// Cambiar estado del empleo

// Listar los empleos
func (Empleo) ObtenerTodos() (empleos []Empleo, err error) {

	result := DB.Where("borrador = false").Preload("Subarea.Area").Find(&empleos)
	if result.Error != nil {
		return empleos, result.Error
	}

	return empleos, nil
}

// Obtener empleo por el id
func (e *Empleo) ObtenerEmpleoByID() error {

	result := DB.Model(&e).Preload("Subarea.Area").First(&e)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
