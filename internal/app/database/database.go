package database

import (
	"github.com/matheustavarestrindade/CraftyGo/internal/app/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	instance *gorm.DB
}

func New() *Database {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.MinecraftServer{})

	database := &Database{
		instance: db,
	}

	return database
}

func (database *Database) CreateUser(username string, email string, password string, isAdmin bool) {
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
		IsAdmin:  isAdmin,
	}
	database.instance.Create(user)
}

func (database *Database) CreateMinecraftServer(name string, port int) {
	minecraftServer := &models.MinecraftServer{
        GameServer: models.GameServer{
            Name: name,
            Port: port,
        },
	}
	database.instance.Create(minecraftServer)
}
