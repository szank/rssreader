package sources

import (
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

type Source struct {
	ID          string
	URL         string
	Title       string
	Categories  []string
	Description string
}

type Sources struct {
	sourceMap  map[string]Source
	sourceList []Source
}

func New() (*Sources, error) {
	sources := &Sources{}

	// Setup the feeds
	// ----- BBC news
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL("http://feeds.bbci.co.uk/news/uk/rss.xml")
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the BBC news feed")
	}

	sources.sourceList = append(sources.sourceList, Source{
		ID:          "1",
		URL:         "http://feeds.bbci.co.uk/news/uk/rss.xml",
		Title:       feed.Title,
		Categories:  feed.Categories,
		Description: feed.Description,
	})

	// ------ BBC technology
	feed, err = parser.ParseURL("http://feeds.bbci.co.uk/news/technology/rss.xml")
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the BBC technology feed")
	}

	sources.sourceList = append(sources.sourceList, Source{
		ID:          "2",
		URL:         "http://feeds.bbci.co.uk/news/technology/rss.xml",
		Title:       feed.Title,
		Categories:  feed.Categories,
		Description: feed.Description,
	})

	// ---- Reuters technology
	feed, err = parser.ParseURL("http://feeds.reuters.com/reuters/UKdomesticNews?format=xml")
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the reuters news feed")
	}

	sources.sourceList = append(sources.sourceList, Source{
		ID:          "3",
		URL:         "http://feeds.reuters.com/reuters/UKdomesticNews?format=xml",
		Title:       feed.Title,
		Categories:  feed.Categories,
		Description: feed.Description,
	})

	// ---- Reuters technology
	feed, err = parser.ParseURL("http://feeds.reuters.com/reuters/technologyNews?format=xml")
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the reuters technology feed")
	}

	sources.sourceList = append(sources.sourceList, Source{
		ID:          "4",
		URL:         "http://feeds.reuters.com/reuters/technologyNews?format=xml",
		Title:       feed.Title,
		Categories:  feed.Categories,
		Description: feed.Description,
	})

	sources.sourceMap = map[string]Source{
		"1": sources.sourceList[0],
		"2": sources.sourceList[1],
		"3": sources.sourceList[2],
		"4": sources.sourceList[3],
	}

	return sources, nil
}

func (s *Sources) Get(id string) (Source, error) {
	source, ok := s.sourceMap[id]
	if !ok {
		return Source{}, errors.Errorf("could not find source with ID %s", id)
	}

	return source, nil
}

func (s *Sources) All() []Source {
	return s.sourceList
}
