package data

import (
	"database/sql"
	"errors"
)

type Models struct {
	Ads AdModel
	Publics PublicModel
}

var (
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
)


func NewModels(db *sql.DB) Models {
	return Models{
		Ads: AdModel{DB: db},
		Publics: PublicModel{DB: db},
	}
}
