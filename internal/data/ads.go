package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Ad struct {
	ID        	int64     `json:"id"`               
	UserId    	int64	  `json:"user_id"`         
	Link 		string 	  `json:"link"`
	Msg     	string    `json:"msg"`            
	Version  	int32     `json:"version"`
	Paid 		bool `json:"paid"`
	CreatedAt time.Time `json:"created_at"`          
}


type AdModel struct {
	DB *sql.DB
}

func (m AdModel) Insert(ad *Ad) error {
	query := `
		INSERT INTO ads (user_id,link,msg,paid)
		VALUES ($1,$2,$3,$4)
		RETURNING id,created_at,version`

	args := []interface{}{ad.UserId, ad.Link, ad.Msg, ad.Paid}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&ad.ID, &ad.CreatedAt, &ad.Version)
}


func (m AdModel) Get(id int64) (*Ad, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, user_id, msg, link, paid, version
		FROM ads
		WHERE user_id = $1
		ORDER BY created_at DESC 
		`

	var ad Ad

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&ad.ID,
		&ad.CreatedAt,
		&ad.UserId,
		&ad.Msg,
		&ad.Link,
		&ad.Paid,
		&ad.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &ad, nil
}


func (m AdModel) GetByMessage(message string) (*Ad, error) {

	query := `
		SELECT msg, created_at, user_id, id, paid, link, version
		FROM ads
		WHERE msg = $1
		`

	var ad Ad
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, message).Scan(
		&ad.Msg,
		&ad.CreatedAt,
		&ad.UserId,
		&ad.ID,
		&ad.Paid,
		&ad.Link,
		&ad.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &ad, nil
}



func (m AdModel) Update(ad *Ad) error {
	query := `
		UPDATE ads
		SET user_id = $1, msg = $2, link = $3, paid = $4, version = version + 1, created_at = $5
		WHERE id = $6 AND version = $7
		RETURNING version`

	args := []interface{}{
		ad.UserId,
		ad.Msg,
		ad.Link,
		ad.Paid,
		ad.CreatedAt,
		ad.ID,
		ad.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&ad.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
