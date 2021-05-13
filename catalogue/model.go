package catalogue

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Code  string
	Price uint
	Name  string

	CreatorID uint

	Categories []Category `gorm:"many2many:product_cat;"`
}

type Creator struct {
	gorm.Model

	Name     string
	Products []Product
}

type Category struct {
	gorm.Model

	Name string
}
