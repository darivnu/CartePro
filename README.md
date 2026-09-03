# CartePro

## Backend

### Language
- Go

### Database
- PostgreSQL
- GORM

## Frontend

## Getting Started

### Prerequisites
- Docker
- Docker Compose

### Setup

1. Copy the example environment file and adjust values if needed:

   ```sh
   cp .env.example .env
   ```

2. Build and start the containers:

   ```sh
   docker compose up --build
   ```

3. The app will be available at:
   - Frontend: [http://localhost:3000](http://localhost:3000)
   - Backend: [http://localhost:4242](http://localhost:4242)

To stop the containers:

```sh
docker compose down
```
