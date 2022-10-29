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
	APILimit        = 2000
	ListenPort      = 8080
	Language        = "en"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func initialize() *models.DBConn {
	db, err := models.NewDBConn()
	handleError(err)

	populated, err := db.CheckDB()
	handleError(err)

	if !populated {
		ability := client.AbilityReceiver{Endpoint: AbilityEndpoint}
		moves := client.MovesReceiver{Endpoint: MoveEndpoint}
		pokemon := client.PokemonReceiver{Endpoint: PokemonEndpoint}

		// fetch api data
		err := client.FetchData(APILimit, Language, &ability, &moves, &pokemon)
		handleError(err)
		
		err = db.PopulateDB(&ability, &moves, &pokemon)
		handleError(err)
	}

	return db
}

func main() {
	db := initialize()

	gin.SetMode(gin.ReleaseMode)
	srv := server.NewHttpServer()
	server.SetupRoutes(srv, db)

	log.Printf("Server running on port: %v\n", ListenPort)
	srv.Run(fmt.Sprintf(":%v", ListenPort))
}
