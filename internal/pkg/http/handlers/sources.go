package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"rssreader/internal/pkg/sources"
)

func NewSources(provider *sources.Sources) *Sources {
	return &Sources{
		provider: provider,
	}
}

type Sources struct {
	provider *sources.Sources
}

func (s *Sources) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnValue := SourceList{}
	sources := s.provider.All()
	for _, source := range sources {
		returnValue.Sources = append(returnValue.Sources, SourceInfo{
			SourceID:   source.ID,
			Categories: source.Categories,
			SourceURL:  source.URL,
			Summary:    source.Description,
			Title:      source.Title,
		})
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(returnValue); err != nil {
		fmt.Printf("Error encoding the JSON response: %v\n", err)
	}
}
