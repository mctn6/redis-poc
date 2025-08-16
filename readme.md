# Redis POC (Go + PostgreSQL)

Simple POC showing how to use **Redis** for caching database queries in Go.

## Setup

Clone the repo and start services:

```bash
git clone https://github.com/mctn6/redis-poc.git
cd redis-poc
docker-compose up --build
```

You’ll see logs like:

```
✅ Connected to PostgreSQL!
✅ Connected to Redis!
🚀 Server running on :8080
💾 Cache miss for user:1, querying DB...
🔁 Cache hit for user:1
```

## Test API

```
curl http://localhost:8080/user/1
```


First request: Cache miss → loads from DB → caches in Redis

Next requests: Cache hit → served from Redis (faster, less DB load)

## Benefits
⚡ Fast reads from in-memory Redis

🛠 Reduces DB queries for frequent requests

📈 Easy to scale read-heavy apps
