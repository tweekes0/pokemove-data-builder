package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tweekes0/pokemonmoves-backend/internal/models"
)

func (s *httpServer) indexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "index page =D",
		})
	}
}

func (s *httpServer) getPokemon() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.MustGet("id").(int)

		p, err := s.DBConn.PokemonGet(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})

			return
		}

		var mv []*models.MovesJoinRow
		if p.OriginGen != 0 {
			mv, err = s.DBConn.PokemonMovesJoinByGen(id, p.OriginGen)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
	
				return
			}
		}

		ab, err := s.DBConn.PokemonAbilitiesJoin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pokemon": p,
			"moves": mv,
			"abilities": ab,
		})
	}

	return gin.HandlerFunc(fn)
}

func (s *httpServer) getAllPokemon() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		p, err := s.DBConn.PokemonGetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H {
			"data": p,
		})
	}
		
	return gin.HandlerFunc(fn)
}

func (s *httpServer) validateID() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid id, id must be int",
			})

			return
		}

		ok, err := s.DBConn.PokemonExists(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})

			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "no pokemon with that id",
			})

			return
		}

		c.Set("id", id)
		c.Next()
	}

	return gin.HandlerFunc(fn)
}
