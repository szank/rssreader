package filters

import "rssreader/internal/pkg/sources"

type Categories struct {
	categories map[string]struct{}
}

func NewCategories(categories []string) *Categories {
	returnValue := &Categories{
		categories: map[string]struct{}{},
	}

	for _, category := range categories {
		returnValue.categories[category] = struct{}{}
	}

	return returnValue
}

func (c *Categories) Apply(article sources.Article) (bool, error) {
	for _, category := range article.Categories {
		if _, ok := c.categories[category]; ok {
			return true, nil
		}
	}

	return false, nil
}
