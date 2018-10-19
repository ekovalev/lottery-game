package main

import (
	"go-server/handler"
	"go-server/monitor"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	api := router.Group("/api/v1")

	h := handler.New(luckyPairs, jackpot, bonusGames)
	api.Use(h.ParseRequestMiddleware)
	api.POST("/bets", h.Play)

	return router
}

func main() {
	log.Info("Initializing app")
	Init()
	router := SetupRouter()

	// Start monitoring routines
	m := monitor.New(luckyPairs)
	go m.RemoveStalePair()
	go m.RefillStack()

	// Start and run the server
	log.Info("Starting server on port 3000")
	router.Run(":3000")
}
