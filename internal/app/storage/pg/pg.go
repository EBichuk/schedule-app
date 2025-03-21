package pg

import (
	"fmt"
	"log"
	"schedule-app/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func InitDB(cnfg config.DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cnfg.Host, cnfg.Username, cnfg.Password, cnfg.Dbname, cnfg.Port, cnfg.Sslmodel)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключить к бд %v", err)
	}
	return db, err
}
