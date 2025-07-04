package db

type Option int

const (
	OPTION_COPY_FULL_PATH Option = iota
)

var OptionsMap = map[Option]string{
	OPTION_COPY_FULL_PATH: "COPY_FULL_PATH",
}

func (repo *SQLiteRepository) GetOptionsFromRepo() (map[string]bool, error) {
	query := `SELECT key, value FROM options`
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	options := make(map[string]bool)
	for rows.Next() {
		var key string
		var value bool
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		options[key] = value
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return options, nil
}

func (repo *SQLiteRepository) SetOptionRepo(key string, value bool) error {
	stmt := `INSERT INTO options (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?`
	_, err := repo.Conn.Exec(stmt, key, value, value)
	return err
}
