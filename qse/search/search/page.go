package search

// Page is a page that has been indexed
type Page struct {
	ID        int64  `json:"id"`
	Highlight string `json:"highlight"`
	Title     string `json:"title"`
	URI       string `json:"uri"`
	IndexTime string `json:"index_time"`
}

// Result is the result of a search, with a list of pages
type Result struct {
	Pages     []*Page `json:"pages"`
	PageCount int     `json:"page_count"`
}
