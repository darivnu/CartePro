# CartePro build artifact

This bundle is self-contained: it does not need the git repo, Node, or Go
installed to build anything — everything here is already compiled. You still
need a reachable Postgres instance and a static file server for the frontend.

## Contents

- `backend/cartepro` — statically linked Linux (amd64) backend binary.
- `frontend/dist/` — built frontend static assets.
- `env.example` — reference environment variables for the backend.

## Running

1. Copy `env.example` to `.env` and point `DB_HOST`/`DB_PORT`/`DB_USER`/
   `DB_PASSWORD`/`DB_NAME` at a running Postgres instance.

2. Start the backend (reads its config from `.env` in the working directory,
   or from exported environment variables):

   ```sh
   ./backend/cartepro
   ```

3. Serve the frontend static files with any static file server, e.g.:

   ```sh
   npx serve -s frontend/dist -l 3000
   ```
