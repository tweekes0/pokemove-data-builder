package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tweekes0/pokemonmoves-backend/internal/models"
)

func indexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "index page =D",
		})
	}
}

func getPokemon(db *models.DBConn) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ok := c.MustGet("all").(bool)

		if ok {
			p, err := db.PokemonGetAll()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"data": p,
			})
		}

		id := c.MustGet("id").(int)
		p, err := db.PokemonGet(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pokemon": p,
		})
	}

	return gin.HandlerFunc(fn)
}

func validateID(db *models.DBConn) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.Query("id") == "" {
			c.Set("all", true)
			return
		}

		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid id, id must be int",
			})

			return

		}

		ok, err := db.PokemonExists(id)
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
		c.Set("all", false)
		c.Next()
	}

	return gin.HandlerFunc(fn)
}
