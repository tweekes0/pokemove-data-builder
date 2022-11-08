package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/tweekes0/pokemoves/src/backend/internal/client"
	"github.com/tweekes0/pokemoves/src/backend/internal/models"
	"github.com/tweekes0/pokemoves/src/backend/internal/server"
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

func initializeDB() *models.DBConn {
	db, err := models.NewDBConn()
	handleError(err)

	populated, err := db.CheckDB()
	handleError(err)

	if !populated {
		ability := client.AbilityReceiver{Endpoint: AbilityEndpoint}
		moves := client.MovesReceiver{Endpoint: MoveEndpoint}
		pokemon := client.PokemonReceiver{Endpoint: PokemonEndpoint}

		// fetch api data
		log.Println("Fetching API data")
		err := client.FetchData(APILimit, Language, &ability, &moves, &pokemon)
		handleError(err)

		log.Println("Populating database")
		err = db.PopulateDB(&ability, &moves, &pokemon)
		handleError(err)
	}

	return db
}

func main() {
	db := initializeDB()

	gin.SetMode(gin.ReleaseMode)
	srv := server.NewHttpServer(db)

	log.Printf("Server running on port: %v\n", ListenPort)
	srv.Run(fmt.Sprintf(":%v", ListenPort))
}
