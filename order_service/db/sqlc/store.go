package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	ConnPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		Queries:  New(connPool),
		ConnPool: connPool,
	}
}

func NewPgxPool(dbSource string) (*pgxpool.Pool, error) {
	// Set up a context with a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new connection pool
	pool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	log.Println("Connected to the database successfully")
	return pool, nil
}
