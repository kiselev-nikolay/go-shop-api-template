package testtools

import (
	"log"

	"github.com/kiselev-nikolay/go-test-docker-dependencies/testdep"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	if err != nil {
		log.Fatalln(err)
	}
	return db, func() {
		err := stop()
		if err != nil {
			log.Fatalln(err)
		}
	}
}
