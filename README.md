# Bookmarks App

A personal bookmarking system that solves the "bookmark graveyard" problem by actively resurfacing saved content and making it effortlessly searchable.

**Design philosophy:** Minimal, fast, frictionless.

## Tech Stack

**Backend:**
- Go (stdlib + Chi if needed)
- PostgreSQL (with native full-text search)
- sqlc (type-safe SQL)
- Docker + Docker Compose
- Deployed on Fly.io

**Frontend:**
- React + TypeScript
- Vite (build tool)
- Deployed on Cloudflare Pages

**Future:**
- Chrome Extension (React + TypeScript)
- Mobile Share Extension (React Native + Expo)

## Project Structure

```
bookmarks/
├── backend/
│   ├── cmd/
│   │   ├── api/          # API server entrypoint
│   │   └── worker/       # Background worker entrypoint
│   ├── internal/         # Private Go packages
│   ├── migrations/       # SQL migration files
│   └── Dockerfile
├── web/                  # React frontend
├── extension/            # Chrome extension (future)
├── docs/                 # Documentation
└── docker-compose.yml    # Local dev environment
```

## Local Development

### Prerequisites

- Docker Desktop
- Go 1.25+
- Node.js 20+
- npm

### Setup

**1. Clone the repository**

```bash
git clone <repo-url>
cd bookmarks
```

**2. Start backend + database**

```bash
docker compose up -d
```

The API will be available at http://localhost:8080

**3. Start frontend**

```bash
cd web
npm install
npm run dev
```

The frontend will be available at http://localhost:5173

### Useful Commands

**Backend:**
```bash
# View logs
docker compose logs -f api

# Rebuild after code changes
docker compose up -d --build api

# Connect to database
docker compose exec db psql -U bookmarks bookmarks_dev

# Stop everything
docker compose down
```

**Frontend:**
```bash
cd web/

# Development server
npm run dev

# Production build
npm run build

# Preview production build
npm run preview
```

## Production Deployment

### Backend (Fly.io)

```bash
cd backend/

# First time setup
flyctl auth login
flyctl launch

# Deploy updates
flyctl deploy
```

**Live URL:** https://backend-empty-sun-8345.fly.dev

### Frontend (Cloudflare Pages)

```bash
cd web/

# Build
npm run build

# First time setup
wrangler login

# Deploy
npx wrangler pages deploy dist
```

**Live URL:** https://bookmarks-web-qbt.pages.dev

## Documentation

- [Product Requirements](docs/bookmark-app-prd.md)
- [Architecture](docs/bookmark-app-architecture.md)
- [Tech Stack Details](docs/bookmark-app-tech-stack.md)
- [Implementation Plan](docs/bookmark-app-implementation-plan.md)
- [Learning Log](docs/learning-log-claude.md)

## Current Status

**Milestone 0: Complete** ✅
- Development environment set up
- Backend and frontend deployed
- End-to-end connectivity verified

**Next:** Milestone 1 - Core data model and basic API endpoints
