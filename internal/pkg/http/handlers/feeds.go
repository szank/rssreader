package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"rssreader/internal/pkg/filters"
	"rssreader/internal/pkg/sources"
)

type Articles struct {
	provider *sources.Sources
}

func NewArticles(provider *sources.Sources) *Articles {
	return &Articles{
		provider: provider,
	}
}

func (a *Articles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	count, err := getQueryCount(r)
	if err != nil {
		render.Render(w, r, invalidRequestError(err))
		return
	}

	maxTimestamp, err := getMaxTimestamp(r)
	if err != nil {
		render.Render(w, r, invalidRequestError(err))
		return
	}

	categories, err := getCategories(r)
	if err != nil {
		render.Render(w, r, invalidRequestError(err))
		return
	}

	sources := a.provider.All()
	chronologicalIterator, err := filters.NewChronologicalSource(maxTimestamp, sources)
	if err != nil {
		fmt.Printf("error retrieving article list: %v\n", err)
		render.Render(w, r, internalServerError(err))
		return
	}

	activeFilters := []filters.Filter{}
	if len(categories) > 0 {
		activeFilters = append(activeFilters, filters.NewCategories(categories))
	}
	pipeline := filters.NewPipeline(chronologicalIterator, activeFilters)

	returnValue := ArticleList{}
	articleCount := 0
	for {
		article, done, err := pipeline.Next()
		if err != nil {
			fmt.Printf("error processing article: %v\n", err)
			render.Render(w, r, internalServerError(err))
			return
		}
		if done {
			break
		}

		returnValue.Articles = append(returnValue.Articles, Article{
			Title:                    article.Title,
			PublicationTime:          article.PublicationDate.Format(time.RFC3339),
			PublicationUnixTimestamp: article.PublicationDate.Unix(),
			Categories:               article.Categories,
			FeedTitle:                article.FeedName,
		})
		articleCount++
		if articleCount == count {
			break
		}
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(returnValue); err != nil {
		fmt.Printf("Error encoding the JSON response for /feeds: %v\n", err)
		render.Render(w, r, renderError(err))
	}

}

func getQueryCount(r *http.Request) (int, error) {
	countString := r.URL.Query().Get("count")
	if countString == "" {
		return 10, nil
	}

	parsedCount, err := strconv.ParseInt(countString, 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "invalid count query parameter")
	}

	if parsedCount <= 0 {
		return 0, errors.New("count query parameter must be > 0")
	}

	return int(parsedCount), nil
}

func getMaxTimestamp(r *http.Request) (time.Time, error) {
	unixTimeString := r.URL.Query().Get("maxtimestamp")
	if unixTimeString == "" {
		//return maximum time
		return time.Unix(math.MaxInt64, 0), nil
	}

	parsedUnixTimestamp, err := strconv.ParseInt(unixTimeString, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "invalid maxtimestamp query parameter")
	}

	if parsedUnixTimestamp <= 0 {
		return time.Time{}, errors.New("maxtimestamp query parameter must be > 0")
	}
	timestamp := time.Unix(parsedUnixTimestamp, 0)

	return timestamp, nil
}

func getCategories(r *http.Request) ([]string, error) {
	categoriesString := r.URL.Query().Get("categories")
	if categoriesString == "" {
		return nil, nil
	}

	categories := strings.Split(categoriesString, ",")
	for _, category := range categories {
		if category == "" {
			return nil, errors.New("empty category")
		}
		if strings.TrimSpace(category) != category {
			return nil, errors.Errorf("category %q contains leading or trailing whitespaces", category)
		}
	}

	return categories, nil
}
