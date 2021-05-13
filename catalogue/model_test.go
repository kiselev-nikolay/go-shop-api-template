package catalogue_test

import (
	"testing"

	"github.com/kiselev-nikolay/go-shop-api-template/catalogue"
	"gotest.tools/assert"
)

func TestGorm(t *testing.T) {
	db, stop := MustCreateDB()
	defer stop()

	err := db.AutoMigrate(&catalogue.Creator{})
	assert.NilError(t, err)
	err = db.AutoMigrate(&catalogue.Category{})
	assert.NilError(t, err)
	err = db.AutoMigrate(&catalogue.Product{})
	assert.NilError(t, err)

	db.Create(&catalogue.Creator{
		Name: "test creator",
	})
	var creator catalogue.Creator
	db.Where("name = ?", "test creator").First(&creator)

	assert.Equal(t, "test creator", creator.Name)

	user := catalogue.Product{
		Code:      "abc",
		Price:     100,
		Name:      "test",
		CreatorID: creator.ID,
		Categories: []catalogue.Category{
			{Name: "test cat 1"},
			{Name: "test cat 2"},
			{Name: "test cat 3"},
		},
	}
	db.Create(&user)
	db.Save(&user)

	var product catalogue.Product
	db.First(&product, 1)
	assert.Equal(t, "abc", product.Code)

	err = db.Model(&product).Association("Categories").Find(&product.Categories)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(product.Categories))
}
