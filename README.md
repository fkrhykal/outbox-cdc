# Outbox CDC Pattern with Go, Debezium, and Kafka

This project demonstrates a robust implementation of the **Transactional Outbox Pattern** combined with **Change Data Capture (CDC)** to achieve reliable event-driven microservices.

## Architecture Overview

The system consists of independent microservices sharing events through an asynchronous message broker (Kafka). To ensure data consistency between a service's database and Kafka, we use the Outbox pattern.

### The Flow

1.  **Service Action**: A microservice (e.g., Order Service) performs a database transaction.
2.  **Atomic Write**: Inside the same transaction, it writes the business data (e.g., an Order) and an event record into a specialized `outbox` table.
3.  **CDC (Debezium)**: A Debezium connector (running on Kafka Connect) monitors the `outbox` table's transaction log.
4.  **Event Routing**: Debezium captures the change and publishes it to a Kafka topic.
5.  **Consumption**: Other services or workers consume the event from Kafka to perform downstream actions (e.g., updating inventory).

## Visualizing the Flow

### Sequence Diagram

A detailed view of the interaction between components and the data payloads involved.

![Outbox CDC Sequence](out/diagrams/sequence_diagram/Outbox%20CDC%20Sequence.png)

[View Sequence Diagram Source](diagrams/sequence_diagram.puml)

## System Components

### 1. Microservices

- **Order Service** (Port `9000`): Handles order creation and persists events to its own outbox.
- **Inventory Service** (Port `8000`): Provides an API to query products and state.
- **Inventory Worker**: Processes events (like `OrderPlaced`) received from Kafka.

### 2. Infrastructure

- **PostgreSQL**:
  - Inventory DB: Port `5432`
  - Order DB: Port `5433`
- **Kafka**: Port `29092` (External), `9092` (Internal).
- **Kafka Connect**: Port `8083`. Runs Debezium connectors.

## Getting Started

### Prerequisites

- Go 1.25+
- Podman or Docker with `podman-compose`/`docker-compose`.

### 1. Start Infrastructure

Launch the databases, Kafka, and Kafka Connect:

```bash
podman-compose up -d
```

### 2. Run Services

You can run the services individually using `go run` or locally using a process manager with the provided `Procfile`:

```bash
# Using Hivemind/Overmind
hivemind Procfile

# Or individually
go run cmd/server/order/main.go
go run cmd/server/inventory/main.go
go run cmd/worker/inventory/main.go
```

## Implementation Details

### The Outbox Table

The `outbox` table in each database is defined as:

```sql
CREATE TABLE public.outbox (
    id UUID PRIMARY KEY,
    aggregateid VARCHAR(255) NOT NULL,
    aggregatetype VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    timestamp TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL
);
```

### Event Routing (Debezium)

We use the Debezium Outbox Event Router. It extracts fields from the `outbox` table and maps them to Kafka messages:

- `aggregatetype` -> Used as the Kafka topic name (e.g., `order`, `inventory`).
- `aggregateid` -> Used as the Kafka message key.
- `payload` -> Becomes the message value.

### Kafka Listener

The project includes a custom Go listener (`kafka/debezium_listener.go`) that handles the Debezium envelope format, making it easy to consume strongly-typed events.

## Features

- **Atomic Consistency**: No "dual writes". The database and Kafka are kept in sync via transaction logs.
- **Resilience**: If the service or Kafka stays down, Debezium resumes from where it left off.
- **Clean Architecture**: Business logic is decoupled from event publishing infra.
