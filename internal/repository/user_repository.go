package repository

import (
	"golang.org/x/crypto/bcrypt"

	"database/sql"
	"fmt"
	"restapi/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `SELECT id, email, password_digest FROM users WHERE email=$1`

	row := r.db.QueryRow(query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordDigest)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func HashPassword(password string) (string, error) {
	// Cost factor: 10-14 is typical (higher = slower = more secure)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // cost 10 hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes), err
}

func (r *UserRepository) CreateUser(email, password, name string) (*model.User, error) {
	fmt.Printf("Create user with email %v", email)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // cost 10
	if err != nil {
		panic(err)

	}

	query := "INSERT INTO users (email, password_digest,name) VALUES ($1, $2, $3) RETURNING id, email"

	user := &model.User{}
	err = r.db.QueryRow(query, email, string(bytes), name).Scan(
		&user.ID,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}

	return user, nil

}
