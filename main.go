package main

import (
	config2 "applichic.com/chic_secret/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// LoadConfiguration configurations
	config2.LoadConfiguration()

	// Init database
	db, err := config2.InitDB()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	router := InitRouter()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://chic-secret.com"}
	router.Use(cors.New(config))

	err = router.Run(":3000")

	if err != nil {
		panic(err)
	}
}
