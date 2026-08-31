# CartePro — Frontend

Frontend for CartePro

React + TypeScript, built with Vite. Styling with Tailwind CSS. Routing with React Router. Session-based auth against the Go backend (`/api/v1`), via an HttpOnly cookie.

## Setup

```bash
npm install
cp .env.example .env
```

Set `VITE_API_BASE_URL` in `.env` to wherever the backend is running.

## Run

```bash
npm run dev
```
