package repository

import (
	"database/sql"

	"wearlab_backend/internal/domain"
)

type CommonRepository struct {
	db *sql.DB
}

func NewCommonRepository(db *sql.DB) *CommonRepository {
	return &CommonRepository{db: db}
}

func (r *CommonRepository) GetTypes() ([]domain.Type, error) {
	rows, err := r.db.Query("SELECT id, name FROM type")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var types []domain.Type

	for rows.Next() {
		var t domain.Type
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return types, nil
}