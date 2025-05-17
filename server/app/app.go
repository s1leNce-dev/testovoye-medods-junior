package app

import (
	"fmt"
	"medods/db"
	"medods/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	db.InitDatabase()
	defer db.CloseDB()

	routes.InitRoutes(r, db.GetDB())

	r.Run(fmt.Sprintf("%s:%s",
		os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")))
}
