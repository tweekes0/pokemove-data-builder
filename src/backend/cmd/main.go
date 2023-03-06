package main

import (
	"fmt"
	"log"
	"time"

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

	Retries = 5
	Timeout = time.Millisecond * 1500
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func initializeDB() (*models.DBConn, error) {
	db, err := models.NewDBConn()
	if err != nil {
		return nil, err
	}

	populated, err := db.CheckDB()
	if err != nil {
		return nil, err
	}

	if !populated {
		ability := client.AbilityReceiver{Endpoint: AbilityEndpoint}
		moves := client.MovesReceiver{Endpoint: MoveEndpoint}
		pokemon := client.PokemonReceiver{Endpoint: PokemonEndpoint}

		// fetch api data
		log.Println("Fetching API data")
		err := client.FetchData(APILimit, Language, &ability, &moves, &pokemon)
		if err != nil {
			return nil, err
		}

		log.Println("Populating database")
		if err = db.PopulateDB(&ability, &moves, &pokemon); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func main() {
	var i int 
	var db *models.DBConn
	var err error

	for i = 0; i < Retries; i++ {
		db, err = initializeDB()
		if err != nil {
			log.Println("Retrying.....")
			time.Sleep(Timeout)
		} else {
			break
		}
	}

	if i == Retries {
		log.Println("Failed to connect to DB")
		handleError(err)
	}

	gin.SetMode(gin.ReleaseMode)
	srv := server.NewHttpServer(db)

	log.Printf("Server running on port: %v\n", ListenPort)
	srv.Run(fmt.Sprintf(":%v", ListenPort))
}
