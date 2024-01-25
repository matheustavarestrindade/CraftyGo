package models

import "gorm.io/gorm"

type GameServer struct {
	gorm.Model
	Name string `json:"name"`
	Port int    `json:"port"`
}

type MinecraftServer struct {
	GameServer
}
