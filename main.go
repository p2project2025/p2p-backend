package main

import (
	"log"
	"p2p/config"
	"p2p/config/db"
	"p2p/routes"
	midleware "p2p/utils/midleWare"

	"github.com/gin-gonic/gin"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.LoadConfig()
	db.InitMongoClient(config.Cfg.DBConnectionString)

	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(midleware.CORSMiddleware()) // Apply CORS middleware

	// Import and set up routes
	routes.Routes(router)

	router.Run(":8080")
}
