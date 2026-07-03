# HEALMATA

Backend service built with Go and Gin.

---

## Prerequisites

* Go (latest stable version recommended)
* PostgreSQL

---

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd <project-folder>
```

### 2. Create the environment file

Copy the example configuration:

```bash
cp .env.example .env
```

Then update the database configuration inside `.env`

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable
```

---

## Environment

### Development

```env
APP_ENV=development
GIN_MODE=debug
```

### Production

```env
APP_ENV=production
GIN_MODE=release
```

---

## Run

Start the server:

```bash
go run ./cmd/server
```

---

## Health Check

After the server starts successfully, open:

```
http://localhost:8080/auth/health
```

---
