package main

import (
	config2 "applichic.com/chic_secret/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func main() {
	// LoadConfiguration configurations
	config2.LoadConfiguration()

	// Init database
	db, err := config2.InitDB()

	if err != nil {
		panic(err)
	}

	router := InitRouter()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://chic-secret.com"}
	router.Use(cors.New(config))

	//err = router.Run(":3000")
	err = autotls.Run(router, "chic-secret.com")
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}
}
