package database

import (
	"fmt"
	"log"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/configs"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/gormrepo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectDB(cfg *configs.Config) *gorm.DB {

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUsername, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate — เพิ่ม Model ใหม่ตรงนี้เมื่อมีตารางใหม่
	if err := db.AutoMigrate(
		&model.UserGormModel{},
		&model.SessionGormModel{},
		&model.OAuthAccountGormModel{},
		&model.PasswordResetTokenGormModel{},
		&model.PhoneOTPGormModel{},
	); err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}

	return db
}
