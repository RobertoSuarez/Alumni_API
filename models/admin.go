package models

type Admin struct {
	ID               uint   `json:"id" gorm:"primary_key"`
	PromoSeguimiento string `json:"promoSeguimiento"`
	Estado           `json:"-"`
}

func (admin *Admin) SetNil() *Admin {
	admin.Estado = Estado{Usando: false}
	admin.PromoSeguimiento = "-- nil"
	return admin
}
