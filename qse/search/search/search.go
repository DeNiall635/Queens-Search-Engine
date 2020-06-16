package search

import (
	"database/sql"
	"math"
)

// Searcher allows for searching for pages based on a query using pagination
type Searcher interface {
	Search(query string, page int) (*Result, error)
}

// DBSearch searches a database
type DBSearch struct {
	DB *sql.DB
}

const pageSize = 10

// Search takes in a query and a page, and produces a search result by querying the database
func (s *DBSearch) Search(query string, page int) (*Result, error) {
	offset := page * pageSize
	result, err := s.DB.Query("SELECT id, title, uri, index_time, ts_headline('english', content, plainto_tsquery($1)) AS highlight, count(*) OVER() AS count FROM page WHERE tokens @@ plainto_tsquery($2) OFFSET $3 LIMIT $4;", query, query, offset, pageSize)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return &Result{
			PageCount: 0,
			Pages:     []*Page{},
		}, nil
	}
	defer result.Close()

	var pages []*Page
	var rowCount int
	for result.Next() {
		page := Page{}
		err = result.Scan(&page.ID, &page.Title, &page.URI, &page.IndexTime, &page.Highlight, &rowCount)
		if err != nil {
			return nil, err
		}
		pages = append(pages, &page)
	}

	err = result.Err()
	if err != nil {
		return nil, err
	}

	return &Result{
		PageCount: int(math.Ceil(float64(rowCount) / float64(pageSize))),
		Pages:     pages,
	}, nil
}
