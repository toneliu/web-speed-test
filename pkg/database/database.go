package database

import (
	"fmt"
	"log"
	"speedtest/pkg/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func Init(dbPath string) error {
	var err error
	// 使用 modernc.org/sqlite 驱动 (纯 Go 实现, 不需要 CGO)
	DB, err = gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dbPath,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.AutoMigrate(&models.User{}, &models.Unit{}, &models.SpeedTest{}, &models.TopologyLink{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err := seedData(); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func seedData() error {
	var adminCount int64
	DB.Model(&models.User{}).Where("is_admin = ?", true).Count(&adminCount)
	if adminCount > 0 {
		return nil
	}

	adminPwd, err := hashPassword("admin123")
	if err != nil {
		return err
	}
	admin := models.User{
		Username: "admin",
		Password: adminPwd,
		IsAdmin:  true,
	}
	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	testUnits := []struct {
		name        string
		username    string
		password    string
	}{
		{"单位A", "unita", "unita123"},
		{"单位B", "unitb", "unitb123"},
		{"单位C", "unitc", "unitc123"},
	}

	for _, tu := range testUnits {
		unit := models.Unit{Name: tu.name}
		if err := DB.Create(&unit).Error; err != nil {
			return err
		}

		pwd, err := hashPassword(tu.password)
		if err != nil {
			return err
		}

		user := models.User{
			Username: tu.username,
			Password: pwd,
			IsAdmin:  false,
			UnitID:   &unit.ID,
		}
		if err := DB.Create(&user).Error; err != nil {
			return err
		}
	}

	log.Println("Default data seeded successfully")
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
