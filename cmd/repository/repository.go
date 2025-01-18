package repository

import "waybackdownloader/cmd/data"

type Repository interface {
	Migrate() error
	InsertURL(url data.Link) (*data.Link, error)
	CollectionURL(websiteURL string) (*data.Link, error)
	UpdateURL(url data.Link) error
}
