package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

// ProgramByID returns a Gin handler function searches a program by its ID
//
// @Summary Search program by ID
// @Description Search program by ID
// @Tags search-programs
// @ID search-program-by-id
// @Param id path string true "id"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/search/programs/{id} [get]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler searchHandler) ProgramByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		programID := c.Param("id")
		programs, err := handler.searchApi.ProgramByID(c, programID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, programs)
		return
	}
}

// Programs returns a Gin handler function that searches all programs
//
// @Summary Search program by ID
// @Description Search program by ID
// @Tags search-programs
// @ID search-all-programs
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/search/programs/ [get]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler searchHandler) Programs() gin.HandlerFunc {
	return func(c *gin.Context) {
		programs, err := handler.searchApi.Programs(c)
		if err != nil {
			log.Error().Err(err).Msg("could not search programs")
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, programs)
		return
	}
}
