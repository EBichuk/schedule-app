package integration

import (
	"fmt"
	"log"
	"schedule-app/internal/config"
	"schedule-app/internal/domain/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func InitDB(cnfg config.DbConfig) (*gorm.DB, error) {

	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cnfg.Username,
		cnfg.Password,
		cnfg.Host,
		cnfg.Port,
		cnfg.Dbname)
	//dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cnfg.Host, cnfg.Username, cnfg.Password, cnfg.Dbname, cnfg.Port, cnfg.Sslmode)
	log.Println("Connection DSN:", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключить к бд %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)

	return db, err
}

func MigrationSchedule(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.Schedule{})
	if err != nil {
		return err
	}
	return nil
}
