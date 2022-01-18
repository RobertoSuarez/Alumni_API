package models

type Admin struct {
	ID               uint   `json:"id" gorm:"primary_key"`
	PromoSeguimiento string `json:"promoSeguimiento"`
	Estado           `json:"-"`
}
