package db

import "github.com/jackc/pgx/v5/pgxpool"

// Store represents an interface for interacting with the database.
type Store interface {
	Querier
}

// SQLStore represents an implementation of the Store interface that uses SQL queries to interact with the database.
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new SQLStore instance.
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
