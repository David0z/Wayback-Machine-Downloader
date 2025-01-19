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
	query := fmt.Sprintf(`SELECT urlkey, timestamp, original, mimetype, statuscode, websiteurl, downloaded FROM links WHERE websiteurl = '%s' AND statuscode = '200'`, websiteURL)

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

type MimetypeQuantity struct {
	Mimetype   string `json:"mimetype"`
	RowCount   int    `json:"row_count"`
	Downloaded int    `json:"downloaded_count"`
}

func (repo *SQLiteRepository) MimetypeQuantity(websiteURL string) ([]MimetypeQuantity, error) {
	query := fmt.Sprintf(`SELECT 
		mimetype, 
    COUNT(*) AS row_count,
    SUM(CASE WHEN downloaded = true THEN 1 ELSE 0 END) AS downloaded_count 
		FROM links 
		WHERE websiteurl = '%s' AND statuscode = '200' 
		GROUP BY mimetype;`, websiteURL)

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mimetypes []MimetypeQuantity

	for rows.Next() {
		var h MimetypeQuantity
		err := rows.Scan(
			&h.Mimetype,
			&h.RowCount,
			&h.Downloaded,
		)
		if err != nil {
			return nil, err
		}

		mimetypes = append(mimetypes, h)
	}

	return mimetypes, nil
}

func (repo *SQLiteRepository) HasAny(websiteURL string) (bool, error) {
	query := `SELECT 1 FROM links WHERE websiteurl = ? AND statuscode = '200'`

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

func (repo *SQLiteRepository) GetOne(websiteURL string, mimetypes []string, offset int) (*data.Link, error) {
	placeholders := strings.Repeat("?,", len(mimetypes))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`SELECT urlkey, timestamp, original, mimetype, statuscode, websiteurl, downloaded FROM links WHERE websiteurl = ? AND statuscode = '200' AND downloaded = false AND mimetype IN (%s) ORDER BY urlkey LIMIT 1 OFFSET ?`, placeholders)

	args := append([]interface{}{}, websiteURL)
	args = append(args, convertToInterfaceSlice(mimetypes)...)
	args = append(args, offset)

	var link data.Link
	row := repo.Conn.QueryRow(query, args...)
	err := row.Scan(
		&link.Urlkey,
		&link.Timestamp,
		&link.Original,
		&link.Mimetype,
		&link.Statuscode,
		&link.WebsiteURL,
		&link.Downloaded,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &link, nil
}

func convertToInterfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
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
