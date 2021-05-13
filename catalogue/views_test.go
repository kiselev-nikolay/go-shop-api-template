package catalogue_test

import (
	"log"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kiselev-nikolay/go-shop-api-template/catalogue"
	"github.com/kiselev-nikolay/go-test-docker-dependencies/testdep"
	"gotest.tools/assert"
)

func MustCreateDB() (*gorm.DB, func()) {
	dbPort, err := testdep.FindFreePort()
	if err != nil {
		log.Fatalln(err)
	}
	dockerPg := testdep.Postgres{
		Port:     dbPort,
		User:     "test",
		Password: "test",
		Database: "test",
	}
	stop, err := dockerPg.Run(10)
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	db, err := gorm.Open(postgres.Open(dockerPg.ConnString()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	return db, func() {
		err := stop()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func TestViews(t *testing.T) {
	db, stop := MustCreateDB()
	defer stop()

	err := db.AutoMigrate(&catalogue.Product{})
	assert.NilError(t, err)

	db.Create(&catalogue.Product{Code: "D42", Price: 100})

	var product catalogue.Product
	db.First(&product, 1)
	db.First(&product, "code = ?", "D42")

	db.Model(&product).Update("Price", 200)

	db.Model(&product).Updates(catalogue.Product{Price: 200, Code: "F42"})
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	db.Delete(&product, 1)
}
