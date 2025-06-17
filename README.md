# Insider Message Sending System

This project is an automated system designed to fetch messages from a PostgreSQL database, send them via a webhook, and track their status, leveraging Redis for caching.

---

## 🚀 Key Features

- **Automated Dispatch:** Sends 2 unsent messages every 2 minutes.
- **Webhook Integration:** Dispatches messages to a configurable external webhook.
- **Redis Caching:** Caches `messageId` and send time for dispatched messages.
- **API Control:** Provides endpoints to start/stop the scheduler and retrieve sent messages.
- **Containerized:** Built with Go, PostgreSQL, and Redis, all managed via Docker Compose.

---

## 🛠️ Getting Started

To set up and run the project, ensure you have Docker and Docker Compose installed.

### 1. Clone the Repository

```sh
git clone <repository-url>
cd insider-project
```
> Replace `<repository-url>` with the actual URL of Git repository.

---

### 2. Configure Webhook URL

The application sends messages to a webhook. Please create a `.env` file in the project root (next to `docker-compose.yml`) and set your specific `WEBHOOK_URL` and `WEBHOOK_AUTH_KEY` that you wish to use for testing. For consistency, I recommend using the `WEBHOOK_URL` and `WEBHOOK_AUTH_KEY` that I will provide separately, which were used during my development and testing.

```env
# .env example
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_DB=mydb
REDIS_ADDR=redis:6379
WEBHOOK_URL=YOUR_PROVIDED_WEBHOOK_URL_HERE
WEBHOOK_AUTH_KEY=YOUR_PROVIDED_AUTH_KEY_HERE
PORT=8080
```

---

### 3. One-Command Setup

Execute the following single command. This will build the application, start all services (app, db, redis), run database migrations, and seed initial messages into the database.

```sh
make all
```

---

## ⚙️ How It Works

Upon `make all` completion, the app container starts, connects to PostgreSQL and Redis, and the message scheduler automatically begins running in the background. Every 2 minutes, it fetches unsent messages, dispatches them to the configured webhook, updates their status in the DB, and caches the webhook's `messageId` in Redis.

---

## 📡 API Endpoints

The API is accessible at `http://localhost:8080/api/v1`. For detailed API specifications, refer to `api/docs/swagger.yml`.

- `POST /start`: Manually start the scheduler.
- `POST /stop`: Manually stop the scheduler.
- `GET /sent?limit=10&offset=0`: Retrieve a paginated list of sent messages.

---

## 🔍 Verifying Data

### Verifying Redis Data

After the system has been running for a few minutes and messages have been sent, you can verify cached data in Redis:

```sh
make redis
```

This command will connect you to the Redis CLI, where you can then use:

- `KEYS message:*` to list cached message keys
- `GET <key_name>` (e.g., `GET message:1`) to view the cached data

---

### Verifying Database Data

You can connect to the PostgreSQL database to inspect the `messages` table:

```sh
make psql
```

Once connected to the `psql` prompt, you can use the following commands:

- `\dt`: To list all tables in the current database.
- `SELECT * FROM messages;`: To view all records in the `messages` table.

---

## ℹ️ Explanation: 'sent' column ('f' vs 't')

In your PostgreSQL database, the `sent` column in the `messages` table is a BOOLEAN type. PostgreSQL represents boolean FALSE values as `f` and TRUE values as `t` when displayed in `psql` or other text-based clients.

- **'f' (False):** This indicates that the message has not yet been sent.
  - When messages are initially inserted into the database (either via seeding or through CREATE TABLE defaults), their sent status is set to FALSE.
  - The scheduler specifically retrieves messages where `sent = FALSE`.
- **'t' (True):** This indicates that the message has been successfully sent.
  - After the scheduler successfully dispatches a message via the webhook, the `MarkMessageSent` function in `internal/db/message.go` is called. This function updates the `sent` column to TRUE and records the `sent_at` timestamp.

So, the values `f` and `t` accurately reflect the current sending status of each message in the database.

---