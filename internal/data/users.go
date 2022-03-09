package data

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        	int64     `json:"id"`            
	UserId        int64     `json:"user_id"`                  
	Version  	int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`          
}

type UserModel struct {
	DB *sql.DB
}


func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (user_id)
		VALUES ($1)
		RETURNING id,created_at,version`

	args := []interface{}{user.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UserId, &user.CreatedAt, &user.Version)
}


func (m AdModel) GetAll() ([]*User, error) {

	query := `
		SELECT user_id
		FROM users
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
		
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	
		defer rows.Close()
	
	

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.UserId,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil

}




func (m UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM users
		WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}


