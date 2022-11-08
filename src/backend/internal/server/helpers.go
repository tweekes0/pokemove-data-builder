package server

import (
	"errors"
	"unicode"

	"github.com/gin-gonic/gin"
)

var (
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidID = errors.New("invalid id, id must be int")
	ErrNotFound = errors.New("no record found with that id")
)

func isWord(word string) bool {
	for _, c := range word {
		if !unicode.IsLetter(c) {
			return false
		}
	}

	return true
}

func handleError(
	c *gin.Context, err error, status int, msg string, abort bool) {
		if err != nil {
			if abort {
				c.AbortWithStatusJSON(status, gin.H{
					"error" : msg,
				})

				return
			}

			c.JSON(status, gin.H{
				"error" : err.Error(),
			})

			return 
		}
}