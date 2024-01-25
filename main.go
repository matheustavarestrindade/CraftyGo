package main

import (
	"github.com/matheustavarestrindade/CraftyGo/internal/app/database"
	web_server "github.com/matheustavarestrindade/CraftyGo/internal/app/web"
)

func main() {
	db := database.New()
	db.CreateUser("admin", "email", "password", true)
	web_server.Start("3000")
}
