package ad

import (
	"database/sql"
	"errors"
)

// Advert is an advert that has been stored
type Advert struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	URI     string `json:"uri"`
	Content string `json:"content"`
	Keyword string `json:"keyword"`
}

// Server allows for searching for ads based on a keyword and creating new ads
type Server interface {
	Get(keywords string) (*Advert, error)
	Create(title string, uri string, content string, keyword string) error
}

// DBAdServer searches a database for ads
type DBAdServer struct {
	DB *sql.DB
}

// Get takes in keywords and returns an ad with relevant keywords
func (r *DBAdServer) Get(keywords string) (*Advert, error) {
	var ad Advert
	err := r.DB.QueryRow("SELECT id, title, uri, content, keyword FROM advert WHERE $1 LIKE '%' || keyword || '%' LIMIT 1;", keywords).Scan(&ad.ID, &ad.Title, &ad.URI, &ad.Content, &ad.Keyword)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		err = r.DB.QueryRow("SELECT id, title, uri, content, keyword FROM advert LIMIT 1;").Scan(&ad.ID, &ad.Title, &ad.URI, &ad.Content, &ad.Keyword)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, errors.New("No adverts in the advert database")
		}
	}
	return &ad, nil
}

// Create creates a new advert
func (r *DBAdServer) Create(title string, uri string, content string, keyword string) error {
	_, err := r.DB.Exec("INSERT INTO advert (title, uri, content, keyword) VALUES ($1, $2, $3, $4)", title, uri, content, keyword)
	return err
}
