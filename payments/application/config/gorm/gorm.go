package gorm

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/sqlite"
)

func init() {
	_, b, _, _ := runtime.Caller(0)

	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func ConnectDB(env string) *gorm.DB {
	var dns string
	var db *gorm.DB
	var err error

	if env != "test" {
		dns = os.Getenv("DNS")
		db, err = gorm.Open(os.Getenv("DB_TYPE"), dns)
	} else {
		dns = os.Getenv("DNS_TEST")
		db, err = gorm.Open(os.Getenv("DB_TYPE_TEST"), dns)
	}

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}

	if os.Getenv("AUTO_MIGRATE_DB") == "true" {
		db.AutoMigrate(&entity.Account{}, &entity.Service{}, &entity.ServicePrice{}, &entity.Store{}, &entity.Transaction{})
	}

	return db
}
