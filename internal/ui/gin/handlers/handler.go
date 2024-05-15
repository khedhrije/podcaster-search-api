// Package handlers provides HTTP request handlers for managing categories.
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-search-api/internal/domain/api"
)

// Search represents the interface for managing search.
type Search interface {
	Programs() gin.HandlerFunc
	ProgramByID() gin.HandlerFunc
}

// searchHandler is an implementation of the Search interface.
type searchHandler struct {
	searchApi api.Search
}

// NewSearchHandler creates a new instance of Indexation interface.
func NewSearchHandler(searchApi api.Search) Search {
	return &searchHandler{
		searchApi: searchApi,
	}
}
