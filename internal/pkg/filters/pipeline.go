package filters

import "rssreader/internal/pkg/sources"

type Pipeline struct {
	source  Source
	filters []Filter
}

func NewPipeline(source Source, filters []Filter) *Pipeline {
	return &Pipeline{
		source:  source,
		filters: filters,
	}
}

func (p *Pipeline) Next() (sources.Article, bool, error) {
	// read the articles from the source untill one passes all the filters
	for {
		article, done, err := p.source.Next()
		if err != nil {
			return sources.Article{}, true, err
		}

		if done {
			return sources.Article{}, true, nil
		}

		// nothing to filter out
		if len(p.filters) == 0 {
			return article, false, nil
		}

		ok := false
		for _, filter := range p.filters {
			ok, err = filter.Apply(article)
			if err != nil {
				return sources.Article{}, true, err
			}
			if !ok {
				break
			}
		}

		// all the filters agreed to pass the article through
		if ok {
			return article, false, nil
		}
	}
}
