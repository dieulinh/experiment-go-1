package repository

import (
	"restapi/internal/model"
)

func (r *Repository) GetById(email string) (*model.Product, error) {
	query := `SELECT id, email, password_digest FROM users WHERE email=$1`

	row := r.db.QueryRow(query, email)

	var product model.Product
	err := row.Scan(&product.ID, &product.Name, &product.Description)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetAllProducts retrieves all products from the database
func (r *Repository) GetAllProducts() ([]model.Product, error) {
	query := `SELECT id, name, description, price FROM products ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
