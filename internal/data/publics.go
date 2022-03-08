package data

import (
	"context"
	"database/sql"
	"time"
)

type Public struct {
	ID        			int64     `json:"id"`
	Name     	 		string    `json:"name"`
	Photo     	 		string    `json:"photo"`
	TelegraphLink      	string    `json:"telegraph_link"`
	Username      		string    `json:"username"`
	LinkToUser 			string 	  `json:"link_to_user"`
	LinkToPublic 		string	  `json:"link_to_public"`
	Version   			int32     `json:"version"`
	CreatedAt 			time.Time `json:"created_at"`
}

type PublicModel struct {
	DB *sql.DB
}

func (m PublicModel) GetAll() ([]Public, error) {

	query := `SELECT id, created_at, name, photo, telegraph_link, username, link_to_user, link_to_public 
	FROM publics`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	publics := []Public{}

	for rows.Next() {
		var public Public

		err := rows.Scan(
			&public.ID,
			&public.CreatedAt,
			&public.Name,
			&public.Photo,
			&public.TelegraphLink,
			&public.Username,
			&public.LinkToUser,
			&public.LinkToPublic,
		)
		if err != nil {
			return nil, err
		}

		publics = append(publics, public)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return publics, nil
}
