package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tweekes0/pokemoves/src/backend/internal/client"
	"github.com/tweekes0/pokemoves/src/backend/internal/models"
)

func (s *httpServer) indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "index page",
	})
}

func (s *httpServer) getPokemon(gen... int) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var g int
		if len(gen) == 0 {
			g = client.CurrentGen
		} else {
			g = gen[0]
		}

		byName := c.MustGet("name").(bool)

		var p *client.Pokemon
		var err error

		if byName {
			p, err = s.DBConn.PokemonGetByName(c.Param("query"), g)
		} else {
			p, err = s.DBConn.PokemonGet(c.MustGet("id").(int), g)
		}

		handleError(c, err, http.StatusInternalServerError, ErrInternalServer.Error(),
			false)

		var mv []*models.MovesJoinRow
		if p.OriginGen != 0 {
			mv, err = s.DBConn.PokemonMovesJoinByGen(p.PokeID, p.OriginGen)
			handleError(c, err, http.StatusInternalServerError, ErrInternalServer.Error(),
				false)
		}

		ab, err := s.DBConn.PokemonAbilitiesJoin(p.PokeID)
		handleError(c, err, http.StatusInternalServerError, ErrInternalServer.Error(),
			false)

		c.JSON(http.StatusOK, gin.H{
			"pokemon":   p,
			"moves":     mv,
			"abilities": ab,
		})
	}

	return gin.HandlerFunc(fn)
}

func (s *httpServer) getAllPokemon(c *gin.Context) {
	p, err := s.DBConn.PokemonGetAllBrief()
	handleError(c, err, http.StatusInternalServerError, ErrInternalServer.Error(),
		false)

	c.JSON(http.StatusOK, gin.H{
		"data": p,
	})
}

func (s *httpServer) validateParam(c *gin.Context) {
	q := c.Param("query")

	if isWord(q) {
		c.Set("name", true)
		c.Next()
		return
	} else {
		c.Set("name", false)
	}

	id, err := strconv.Atoi(q)
	handleError(c, err, http.StatusInternalServerError, ErrInvalidID.Error(), true)

	ok, err := s.DBConn.PokemonExists(id)
	handleError(c, err, http.StatusInternalServerError, ErrInternalServer.Error(), true)

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": ErrNotFound.Error(),
		})

		return
	}

	c.Set("id", id)
	c.Next()
}
