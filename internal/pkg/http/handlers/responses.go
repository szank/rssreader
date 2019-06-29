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
