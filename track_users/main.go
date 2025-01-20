package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "track_users/db/sqlc"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://neondb_owner:XVxkNaf82vPg@ep-twilight-smoke-a8hq9d76.eastus2.azure.neon.tech/neondb?sslmode=require"
	serverAddress = "0.0.0.0:8082"
)

func main() {
	conn, err := db.NewPgxPool(dbSource)

	store := NewStore(conn)
	if err != nil {
		log.Fatal("cannot connect to track user db:", err)
	}

	server := NewServer(store)

	err = ConsumeMessagesFromKafka(server)

}

func ConsumeMessagesFromKafka(server *Server) error {

	topic := "order"
	brokers := []string{"kafka:9092"}

	consumer, err := connectConsumer(brokers)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Consumer started")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)

	go func() {
		for {
			select {
			case err := <-partitionConsumer.Errors():
				fmt.Printf("Error: %s\n", err.Error())
			case msg := <-partitionConsumer.Messages():
				err := track(string(msg.Value), server)
				fmt.Println(err)
			case <-signals:
				fmt.Println("Interrupt detected")
				done <- true
				return
			}
		}
	}()
	<-done

	return nil
}

func track(value string, server *Server) error {

	_, err := server.store.GetCounterUser(context.Background(), value)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = server.store.CreateTrackUser(context.Background(), db.CreateTrackUserParams{Name: value, Counter: pgtype.Int4{Int32: 1, Valid: true}})
		} else {
			return err
		}
		return nil
	}
	err = server.store.UpdateUserCounter(context.Background(), value)

	if err != nil {
		return err
	}

	return nil
}

type ConsumerGroupHandler struct {
}

type Server struct {
	store *Store
}

func NewServer(store *Store) *Server {
	return &Server{
		store: store,
	}
}

type Store struct {
	*db.Queries
	ConnPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		Queries:  db.New(connPool),
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

func connectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
