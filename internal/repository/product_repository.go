package repository

import (
	"context"
	"rest_api_pks/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, "SELECT id, title, image_url, name, price, description, specifications, quantity, is_favorite, in_cart FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Title, &product.ImageURL, &product.Name, &product.Price, &product.Description, &product.Specifications, &product.Quantity, &product.IsFavorite, &product.InCart)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	err := r.db.QueryRow(ctx, "SELECT id, title, image_url, name, price, description, specifications, quantity, is_favorite, in_cart FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Title, &product.ImageURL, &product.Name, &product.Price, &product.Description, &product.Specifications, &product.Quantity, &product.IsFavorite, &product.InCart)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO products (title, image_url, name, price, description, specifications, quantity)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, product.Title, product.ImageURL, product.Name, product.Price, product.Description, product.Specifications, product.Quantity)
	return err
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	_, err := r.db.Exec(ctx, `
		UPDATE products SET title = $1, image_url = $2, name = $3, price = $4, description = $5, specifications = $6, quantity = $7
		WHERE id = $8
	`, product.Title, product.ImageURL, product.Name, product.Price, product.Description, product.Specifications, product.Quantity, product.ID)
	return err
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	return err
}

func (r *ProductRepository) UpdateProductQuantity(ctx context.Context, id int, quantity int) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET quantity = $1 WHERE id = $2", quantity, id)
	return err
}

func (r *ProductRepository) ToggleFavorite(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET is_favorite = NOT is_favorite WHERE id = $1", id)
	return err
}

func (r *ProductRepository) ToggleCart(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET in_cart = NOT in_cart WHERE id = $1", id)
	return err
}
