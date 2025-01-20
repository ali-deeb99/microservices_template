package api

import (
	"fmt"
	"log"
	"net/http"
	db "order_service/db/sqlc"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateOrderRequest struct {
	Name string `form:"name" binding:"required"`
	Note string `form:"note"`
}

func (server *Server) CreateOrder(c *gin.Context) {

	var req CreateOrderRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	tx, err := server.store.ConnPool.Begin(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("issue encountered while trying to create an order"))
		return
	}

	defer tx.Rollback(c)

	q := db.New(tx)

	name, err := q.CreateOrder(c, db.CreateOrderParams{
		Name:   req.Name,
		Note:   pgtype.Text{String: req.Note, Valid: true},
		Status: 1, // 1 pending 2 processing 3 done
	})

	if err != nil {
		// in here we can add a middleware to add the errors to track it later , also we can use tool like sentry
		c.JSON(http.StatusInternalServerError, errorResponse("issue encountered while trying to create an order"))
		return
	}

	producer, err := KafkaProducerConnection([]string{"kakfa:9092"})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("issue encountered while trying to create an order"))
		return
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic:     "order", // replace with your topic name
		Value:     sarama.StringEncoder(name),
		Partition: 0,
	}

	// Send the message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("issue encountered while trying to create an order"))
		return
	}

	fmt.Printf("Message sent to partition %d with offset %d\n", partition, offset)

	c.JSON(http.StatusCreated, successResponse("order added  successfully"))

	tx.Commit(c)
}

func KafkaProducerConnection(brokerUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	// Create a new producer
	producer, err := sarama.NewSyncProducer(brokerUrl, config) // replace with your Kafka brokers
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
		return nil, err
	}

	return producer, nil

}

func errorResponse(err string) gin.H {
	return gin.H{"error": err}
}

func successResponse(data interface{}) gin.H {
	return gin.H{"message": "success", "data": data}
}
