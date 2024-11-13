package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	db        *sql.DB
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

// UserInput represents the input for creating a new user.
type UserInput struct {
	Name  string
	Email string
}

func (c *User) Create(input UserInput) (User, error) {
	var id int64
	now := time.Now()
	err := c.db.QueryRow("INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, $3, $3) RETURNING id", input.Name, input.Email, now).Scan(&id)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return User{ID: id, Name: input.Name, Email: input.Email, CreatedAt: now, UpdatedAt: now}, nil
}

func (c *User) FindAll() ([]User, error) {
	rows, err := c.db.Query("SELECT id, name, email, created_at, updated_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
