package config

import (
	"database/sql"
	"path"
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/repository/db"
	"waybackdownloader/cmd/util"

	"github.com/rivo/tview"
)

type Config struct {
	App     *tview.Application
	DB      *db.SQLiteRepository
	Options *Options
}

func (c *Config) Init() {
	util.CreatePathIfNotExists(data.MAIN_PATH)
	sqlDB, err := c.connectSQL()
	if err != nil {
		panic("failed to initiate database")
	}

	c.setupDB(sqlDB)
}

func (c *Config) connectSQL() (*sql.DB, error) {
	path := path.Join(data.MAIN_PATH, "sql.db")

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *Config) setupDB(sqlDB *sql.DB) {
	c.DB = db.NewSQLiteRepository(sqlDB)

	optionsMap, err := c.DB.Migrate()
	if err != nil {
		panic(err)
	}

	c.Options = NewOptions(optionsMap)
}
