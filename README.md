Order and Track Users System

Overview

This project demonstrates a distributed system consisting of two microservices: order_service and track_users. The services communicate asynchronously using Apache Kafka as the message broker.

Services

order_service: Handles user requests and stores orders in its database.

track_users: Listens to events from Kafka and manages user tracking data.

Full Scenario

A user sends a request to the order_service API.

Payload: { "name": "<user_name>", "note": "<user_note>" }

The order_service:

Creates a new record in its database table orders with the provided data.

Sends an event containing the name field to the Kafka message broker.

The track_users service:

Listens for events on the Kafka topic.

Checks if the user already exists in its database table track_user.

If the user does not exist, it creates a new record.

If the user exists, it increments their counter field.

System Architecture

Components

order_service:

Database: orders table

Fields: id, name, note, created_at

Kafka Producer: Sends name to Kafka.

track_users:

Database: track_user table

Fields: id, name, counter, created_at

Kafka Consumer: Listens to events from Kafka.

Apache Kafka:

Topic: user_events

Workflow

[User] --> [order_service API] --> [order_service DB]
                                 --> [Kafka Producer] --> [Kafka Broker]
[Kafka Broker] --> [Kafka Consumer] --> [track_users DB]

Technologies Used

Programming Language: Go (Golang)

Message Broker: Apache Kafka

Databases: PostgreSQL (for both services)

Libraries:

Sarama: Kafka client for Go

Setup Instructions

Clone the Repository:

git clone <repository-url>
cd <repository-folder>

Run Kafka Broker:
Install and configure Apache Kafka.

Start Services:

order_service

cd order_service
go run main.go

track_users

cd track_users
go run main.go

Test the System:
Use tools like Postman or curl to send requests to order_service API and observe updates in track_users.

Example API Usage

Request to order_service

Endpoint: POST /api/v1/orders

Payload:

{
  "name": "John Doe",
  "note": "First order"
}

Response:

{
  "status": "success",
  "order_id": 123
}

Kafka Event Example

Produced Event:

{
  "name": "John Doe"
}

Record in track_user Table

Before:

id

name

counter

created_at

1

John Doe

1

2025-01-01 12:00

After Another Event:

id

name

counter

created_at

1

John Doe

2

2025-01-01 12:00

Future Enhancements

Add authentication and authorization.

Implement retries for Kafka consumer.

Introduce monitoring tools for Kafka and the microservices.

Contributing

Contributions are welcome! Please open an issue or submit a pull request.

License

This project is licensed under the MIT License. See the LICENSE file for details.