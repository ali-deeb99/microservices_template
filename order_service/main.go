package main

import (
	"log"
	"order_service/api"
	db "order_service/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://order_db_owner:MWn8FUumJdG9@ep-cool-star-a8ohjzvn.eastus2.azure.neon.tech/order_db?sslmode=require"
	serverAddress = "0.0.0.0:8081"
)

func main() {
	conn, err := db.NewPgxPool(dbSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
