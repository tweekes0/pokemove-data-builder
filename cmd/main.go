package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
	"github.com/tweekes0/pokemonmoves-backend/internal/models"
	"github.com/tweekes0/pokemonmoves-backend/internal/server"
)

const (
	AbilityEndpoint = "https://pokeapi.co/api/v2/ability"
	MoveEndpoint    = "https://pokeapi.co/api/v2/move"
	PokemonEndpoint = "https://pokeapi.co/api/v2/pokemon"
	ListenPort      = 8080
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	ability := client.AbilityReceiver{}
	moves := client.MovesReceiver{}
	pokemon := client.PokemonReceiver{}

	lang := "en"
	limit := 2000

	// fetch api data
	err := client.GetAPIData(&moves, limit, MoveEndpoint, lang)
	handleError(err)

	err = client.GetAPIData(&pokemon, limit, PokemonEndpoint, lang)
	handleError(err)

	err = client.GetAPIData(&ability, limit, AbilityEndpoint, lang)
	handleError(err)

	// Generate CSV files of fetched API data
	// client.generateCsvs(pokemon, moves, ability)

	db, err := models.NewDBConn()
	handleError(err)

	gin.SetMode(gin.ReleaseMode)

	srv := server.NewHttpServer()
	server.SetupRoutes(srv, db)

	log.Printf("Server running on port: %v\n", ListenPort)
	srv.Run(fmt.Sprintf(":%v", ListenPort))
}
