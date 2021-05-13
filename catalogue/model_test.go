package catalogue_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/kiselev-nikolay/go-shop-api-template/catalogue"
	"gotest.tools/assert"
)

// TestGorm is just Gorm test, nothing important
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

	product := catalogue.Product{
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
	db.Create(&product)
	db.Save(&product)

	product = catalogue.Product{}
	db.First(&product, 1)
	assert.Equal(t, "abc", product.Code)

	err = db.Model(&product).Association("Categories").Find(&product.Categories)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(product.Categories))
	assert.Equal(t, "test cat 1", product.Categories[0].Name)

	tableNames := make([]string, 0)
	db.Raw("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'").Scan(&tableNames)
	sort.Strings(tableNames)
	assert.Equal(t, "categories creators product_cat products", strings.Join(tableNames, " "))

	tableTest := map[string]int{
		"categories":  3,
		"creators":    1,
		"product_cat": 3,
		"products":    1,
	}
	for tname, expected := range tableTest {
		t.Run("Assert length of "+tname, func(t *testing.T) {
			tlength := 0
			db.Raw("SELECT count(*) FROM " + tname).Scan(&tlength)
			assert.Equal(t, expected, tlength)
		})
	}
}
