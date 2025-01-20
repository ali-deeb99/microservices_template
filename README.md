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



Contributing

Contributions are welcome! Please open an issue or submit a pull request.

License

This project is licensed under the MIT License. See the LICENSE file for details.
