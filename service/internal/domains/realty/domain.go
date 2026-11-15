package realty

import "database/sql"

type Domain struct {
	db *sql.DB
}

func NewDomain(db *sql.DB) *Domain {
	return &Domain{db: db}
}
