package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"waybackdownloader/cmd/data"
)

type SelectString string

var (
	errorUpdate = errors.New("failed to update")
)

const (
	SELECT_ONE SelectString = "1"
	SELECT_ALL SelectString = ""
)

type SQLiteRepository struct {
	Conn *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

func (repo *SQLiteRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS links(
			urlkey TEXT PRIMARY KEY,
			timestamp TEXT NOT NULL,
			original TEXT NOT NULL,
			mimetype TEXT NOT NULL,
			statuscode TEXT NOT NULL,
			websiteurl TEXT NOT NULL,
			downloaded BOOLEAN NOT NULL DEFAULT false);
	`

	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) InsertURLs(urls []data.Link) (*[]data.Link, error) {
	if len(urls) == 0 {
		return nil, nil
	}

	const maxParams = 999
	const paramsPerRow = 7
	maxRows := maxParams / paramsPerRow

	for start := 0; start < len(urls); start += maxRows {
		end := start + maxRows
		if end > len(urls) {
			end = len(urls)
		}

		bindValues := []any{}
		queryValues := []string{}

		for _, url := range urls[start:end] {
			bindValues = append(bindValues, url.Urlkey, url.Timestamp, url.Original, url.Mimetype, url.Statuscode, url.WebsiteURL, url.Downloaded)
			queryValues = append(queryValues, "(?, ?, ?, ?, ?, ?, ?)")
		}

		stmt := fmt.Sprintf(`INSERT INTO links (urlkey, timestamp, original, mimetype, statuscode, websiteurl, downloaded)
		VALUES %s ON CONFLICT(urlkey) DO UPDATE SET downloaded = false`, strings.Join(queryValues, ", "))

		_, err := repo.Conn.Exec(stmt, bindValues...)
		if err != nil {
			return nil, err
		}
	}

	return &urls, nil
}

func (repo *SQLiteRepository) CollectionURL(websiteURL string) ([]data.Link, error) {
	query := fmt.Sprintf(`SELECT urlkey, timestamp, original, mimetype, statuscode, websiteurl, downloaded FROM links WHERE websiteurl = '%s'`, websiteURL)

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []data.Link

	for rows.Next() {
		var h data.Link
		err := rows.Scan(
			&h.Urlkey,
			&h.Timestamp,
			&h.Original,
			&h.Mimetype,
			&h.Statuscode,
			&h.WebsiteURL,
			&h.Downloaded,
		)
		if err != nil {
			return nil, err
		}

		links = append(links, h)
	}

	return links, nil
}

func (repo *SQLiteRepository) SelectOne(websiteURL string) (bool, error) {
	query := `SELECT 1 FROM links WHERE websiteurl = ?`

	var exists int
	err := repo.Conn.QueryRow(query, websiteURL).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (repo *SQLiteRepository) UpdateURL(url data.Link) error {
	stmt := `UPDATE links SET downloaded = ? WHERE urlkey = ?`
	res, err := repo.Conn.Exec(stmt, true, url.Urlkey)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorUpdate
	}

	return nil
}
