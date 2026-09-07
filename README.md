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
   - Database (Postgres, for local tools/scripts): `localhost:5433`

### Stopping

If you ran `docker compose up --build` in the foreground, `Ctrl+C` stops the containers. They aren't removed though, so also run:

```sh
docker compose down

```

(add `-v` to also delete the `postgres_data` volume, i.e. wipe the database)

To avoid tying up your terminal in the first place, run it detached instead:

```sh
docker compose up --build -d
docker compose logs -f   # tail logs when you want them
docker compose down      # stop everything
```

### Resetting the database

`backend/database/nuke_database.sh` truncates all tables. It reads DB connection settings from `.env`, but `.env`'s `DB_HOST`/`DB_PORT` (`db`/`5432`) only resolve *inside* the Docker network. To run it from the host, override them to point at the published port:

```sh
DB_HOST=localhost DB_PORT=5433 ./backend/database/nuke_database.sh
```
