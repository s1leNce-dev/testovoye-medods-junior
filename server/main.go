package main

import (
	"log"
	"medods/app"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("[FATAL] %s", err.Error())
	}
}

var (
	isContainer = true
	timeToWait  = time.Second * 5
)

// @title           Тестовое задание MEDODS junior
// @version         1.0
// @description     часть сервиса аутентификации
// @host            localhost:8000
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if isContainer {
		time.Sleep(timeToWait)
	}
	app.Start()
}
