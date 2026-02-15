package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"wearlab_backend/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	currentTime := (time.Now())
	_, err := r.db.Exec(
		"INSERT INTO public.product(name, sku, price, stock, type, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7);",
		product.Name, product.SKU, product.Price, product.Stock, product.Type, currentTime, currentTime,
	)
	return err
}

func (r *ProductRepository) GetByID(id int) (domain.Product, error) {
	var p domain.Product

	row := r.db.QueryRow(`
		SELECT 
			id, name, sku, price, stock, type, created_at, updated_at
		FROM 
			product
		WHERE 
			id = $1;
	`, id)

	err := row.Scan(&p.ID, &p.Name, &p.SKU, &p.Price, &p.Stock, &p.Type, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Product{}, fmt.Errorf("no product found with id %d", id)
		}
		return domain.Product{}, err
	}

	return p, nil
}

func (r *ProductRepository) Update(id int, product *domain.Product) (domain.Product, error) {
	var p domain.Product
	currentTime := time.Now()

	row := r.db.QueryRow(
		`UPDATE public.product
		SET name = $1, sku = $2, price = $3, stock = $4, type = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, name, sku, price, stock, type, created_at, updated_at;`,
		product.Name, product.SKU, product.Price, product.Stock, product.Type, currentTime, id,
	)

	err := row.Scan(
		&p.ID, &p.Name, &p.SKU, &p.Price, &p.Stock, &p.Type, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return domain.Product{}, err
	}

	return p, nil
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec(
		"DELETE FROM public.product WHERE id = $1", id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no record found to delete")
	}

	return nil
}

func (r *ProductRepository) GetAll() ([]domain.Product, int, float64, error) {
	var totalStock int
	err := r.db.QueryRow("SELECT COALESCE(SUM(stock), 0) FROM product").Scan(&totalStock)
	if err != nil {
		return nil, 0, 0, err
	}

	// Get total value (price * stock for all products)
	var totalValue float64
	err = r.db.QueryRow("SELECT COALESCE(SUM(price * stock), 0) FROM product").Scan(&totalValue)
	if err != nil {
		return nil, 0, 0, err
	}

	// Get all products
	rows, err := r.db.Query(`
		SELECT 
			id, name, sku, price, stock, type, created_at, updated_at
		FROM 
			product
		ORDER BY id;
	`)

	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Price, &p.Stock, &p.Type, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, 0, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return products, totalStock, totalValue, nil
}

func (r *ProductRepository) GetWithFilter(category, keyword string) ([]domain.Product, int, float64, error) {
	var (
		products     []domain.Product
		args         []interface{}
		whereClauses []string
	)

	argID := 1

	if category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("type = $%d", argID))
		args = append(args, category)
		argID++
	}
	if keyword != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(name ILIKE $%d OR sku ILIKE $%d)", argID, argID))
		args = append(args, "%"+keyword+"%")
		argID++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Total stock query (sum of all stock quantities)
	totalStockQuery := "SELECT COALESCE(SUM(stock), 0) FROM product " + whereSQL
	var totalStock int
	err := r.db.QueryRow(totalStockQuery, args...).Scan(&totalStock)
	if err != nil {
		return nil, 0, 0, err
	}

	// Total value query with the same filters
	totalValueQuery := "SELECT COALESCE(SUM(price * stock), 0) FROM product " + whereSQL
	var totalValue float64
	err = r.db.QueryRow(totalValueQuery, args...).Scan(&totalValue)
	if err != nil {
		return nil, 0, 0, err
	}

	query := fmt.Sprintf(`
		SELECT 
			id, name, sku, price, stock, type, created_at, updated_at
		FROM 
			product
		%s
		ORDER BY id
	`, whereSQL)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Price, &p.Stock, &p.Type, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, 0, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return products, totalStock, totalValue, nil
}

func (r *ProductRepository) IsSKUDuplicate(sku string) (bool, error) {
	var count int

	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM public.product WHERE sku = $1",
		sku,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ProductRepository) UpdateStock(productID int, quantity int) error {
	result, err := r.db.Exec(
		`UPDATE public.product 
		 SET stock = stock - $1 
		 WHERE id = $2 AND stock >= $1`,
		quantity,
		productID,
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("Insufficient stock or product not found")
	}

	return nil
}

func (r *ProductRepository) UpdatePrice(productID int, newPrice float64) (int64, error) {

	result, err := r.db.Exec(
		"UPDATE public.product SET price = $1 WHERE id = $2",
		newPrice,
		productID,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}