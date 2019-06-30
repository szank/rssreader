package filters

import (
	"sort"
	"time"

	"github.com/pkg/errors"

	"rssreader/internal/pkg/sources"
)

type Source interface {
	Next() (sources.Article, bool, error)
}

type Filter interface {
	Apply(sources.Article) (bool, error)
}

type ChronologicalSource struct {
	sources [][]sources.Article
}

func NewChronologicalSource(timestampsBefore time.Time, sources []sources.Source) (*ChronologicalSource, error) {
	iterator := &ChronologicalSource{}
	for _, singleSource := range sources {
		articles, err := singleSource.AllArticles()
		if err != nil {
			return nil, errors.Wrapf(err, "could not read article list from source %s", singleSource.Title)
		}

		sort.Slice(articles, func(i, j int) bool {
			if articles[i].PublicationDate == nil {
				return false
			}
			if articles[j].PublicationDate == nil {
				return true
			}

			return articles[i].PublicationDate.After(*articles[j].PublicationDate)
		})

		// find the first article that is before the given timestmap
		articlePosition := 0

		for i := range articles {
			if articles[i].PublicationDate == nil {
				return nil, errors.New("article has missing publication date")
			}

			// find first article published before the given timestamp
			// should be binary search
			if articles[i].PublicationDate.Before(timestampsBefore) ||
				articles[i].PublicationDate.Equal(timestampsBefore) {
				articlePosition = i
				break
			}
		}

		// cut off head of the slice until the point that all the articles are older than the provided minTimestamp
		source := articles[articlePosition:]
		iterator.sources = append(iterator.sources, source)
	}

	return iterator, nil
}

func (s *ChronologicalSource) Next() (sources.Article, bool, error) {
	// then get the latest article from the heads of all the available sources
	latestArticleSourceIndex := 0
	latestTimestamp := time.Time{}
	emptySourceCount := 0
	for i, source := range s.sources {
		// if all sources are empty, return finished flag
		if len(source) == 0 {
			emptySourceCount++
			continue
		}
		if source[0].PublicationDate == nil {
			return sources.Article{}, true, errors.New("article has missing publication date")
		}

		if source[0].PublicationDate.After(latestTimestamp) {
			latestArticleSourceIndex = i
			latestTimestamp = *source[0].PublicationDate
		}
	}

	// return true means there is are no articles available
	if emptySourceCount == len(s.sources) {
		return sources.Article{}, true, nil
	}

	// pop the latest article
	article := s.sources[latestArticleSourceIndex][0]
	s.sources[latestArticleSourceIndex] = s.sources[latestArticleSourceIndex][1:]

	return article, false, nil
}
