package repository

import "github.com/jackc/pgx/v5/pgxpool"

type ProductsRepository struct {
	pool *pgxpool.Pool
}

/// я жив и это ок
