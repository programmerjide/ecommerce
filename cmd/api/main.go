package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/programmerjide/ecommerce/internal/config"
	"github.com/programmerjide/ecommerce/internal/database"
	"github.com/programmerjide/ecommerce/internal/logger"
)

func main() {

	log := logger.NewLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	defer func(mainDB *sql.DB) {
		err := mainDB.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}(mainDB)

	gin.SetMode(cfg.Server.GinMode)

	log.Info().Msg("Starting server")
}
