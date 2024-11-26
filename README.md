# Order Info

This service processes orders received from Kafka, stores them in a database, caches them for fast retrieval, and provides an API to search and display order details.

## Setup and Run

1. Start the services using Docker Compose: `docker-compose up -d`

2. Send sample messages to Kafka: `./producer/send_message.sh`

3. Access the order search page in your browser: http://localhost:8080/static/templates/order.html

- Enter the order UID in the input field to search for an order.

4. Alternatively, you can directly query the server for JSON data: http://localhost:8080/order/:id

- Replace `:id` with the desired order UID.

## Data Flow

1. Orders are received as messages from Kafka.
2. The service processes these messages, storing the order data in both:
- **Database**: Persistent storage
- **Cache**: In-memory storage with a 5-minute expiry
3. On server restart:
- All orders are retrieved from the database and reloaded into the cache.