# CartePro build artifact

This bundle is self-contained: it needs only Docker (with the `docker
compose` plugin) installed. It does not need the git repo, Node, Go, or a
pre-existing Postgres instance — the artifact starts its own database
container.

## Contents

- `images.tar` — pre-built `cartepro-backend`, `cartepro-frontend`, and
  `postgres` Docker images.
- `docker-compose.yml` — wires the three images together (frontend, backend,
  database).
- `env.example` — reference environment variables, copied to `.env` on first
  run.
- `RUN_ME.sh` — loads the images and starts the stack.

## Running

From an empty folder containing just this artifact's contents:

```sh
./RUN_ME.sh
```

Then:

- Frontend: http://localhost:3000
- Backend: http://localhost:4242

Stop everything with `docker compose down` (run from the same folder).
