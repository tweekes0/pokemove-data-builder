package server

import (
	"github.com/gin-gonic/gin"
	"github.com/tweekes0/pokemonmoves-backend/internal/models"
)

type httpServer struct {
	*gin.Engine
}

func NewHttpServer() *httpServer {
	return &httpServer{
		gin.Default(),
	}
} 

func SetupRoutes(srv *httpServer, db *models.DBConn) {
	srv.GET("/", indexHandler())
	srv.GET("/pokemon", validateID(db), getPokemon(db))
}