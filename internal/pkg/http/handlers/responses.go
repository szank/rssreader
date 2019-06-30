package handlers

type SourceInfo struct {
	SourceID   string   `json:"sourceID"`
	Categories []string `json:"categories,omitempty"`
	SourceURL  string   `json:"sourceURL"`
	Summary    string   `json:"summary"`
	Title      string   `json:"title"`
}

type SourceList struct {
	Sources []SourceInfo `json:"sources"`
}

type ArticleList struct {
	Articles []Article `json:"articles,omitempty"`
}

type Article struct {
	Title                    string   `json:"title"`
	PublicationTime          string   `json:"publication_time"`
	PublicationUnixTimestamp int64    `json:"publication_unix_timestamp"`
	Categories               []string `json:"categories,omitempty"`
	FeedTitle                string   `json:"feed_title"`
	Description              string   `json:"description"`
	Link                     string   `json:"link`
}
