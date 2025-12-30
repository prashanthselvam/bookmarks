# Bookmark App – Tech Stack

## Overview

This document captures the technology choices for the bookmark app, based on the PRD and architecture decisions.

---

## Stack Summary

| Component | Technology |
|-----------|------------|
| Backend API | Go (stdlib, Chi if needed) |
| Database | PostgreSQL |
| DB access | sqlc |
| Background jobs | Go worker process (same codebase, separate entrypoint) |
| Job queue | Postgres-backed |
| Full-text search | Postgres native (pg_trgm + tsvector) |
| Frontend | React + TypeScript + Vite |
| Chrome extension | React + TypeScript |
| Mobile app | React Native + Expo |
| Containerization | Docker |
| App hosting | Fly.io |
| DB hosting | Fly.io Postgres or Neon |
| Frontend hosting | Vercel |
| Auth | OAuth (Google) – self-implemented flow |

---

## Component Details

### Backend: Go

**Why Go:**
- Learning goal – used at work, want deeper fluency
- Excellent for stateless APIs with background workers
- Single binary deployment, fast cold starts
- Strong stdlib for HTTP and JSON

**Framework approach:**
- Start with Go stdlib (`net/http`) – Go 1.22+ has capable routing with method matching and path parameters
- Add Chi only if friction emerges (middleware chaining, route grouping)
- This approach maximizes learning and minimizes magic

**Database access:**
- sqlc – write SQL, generate type-safe Go code
- Keeps you close to SQL while avoiding manual row scanning
- Good for learning how Go interacts with databases

### Database: PostgreSQL

**Why Postgres:**
- Relational model fits structured bookmark data
- Strong consistency for duplicate detection, soft deletes
- Native full-text search (tsvector + pg_trgm) avoids needing a separate search engine
- Excellent tooling and hosting options

**Full-text search approach:**
- Use Postgres tsvector for title, description, notes
- Use pg_trgm for fuzzy/partial matching on URL
- Revisit dedicated search engine only if this becomes a bottleneck (unlikely at single-user scale)

### Background Jobs: Go Worker

**Architecture:**
- Same codebase as API, different entrypoint
- Polls a Postgres-backed job queue (simple `jobs` table)
- Avoid Redis/external queue for now – Postgres can handle this scale

**Job types:**
- `fetch_metadata` – scrape OG tags from bookmarked URLs
- `purge_trash` – scheduled cleanup of soft-deleted bookmarks > 30 days

### Frontend: React + TypeScript

**Why React over HTMX:**
- "Results as you type" search is awkward in HTMX
- Chrome extension will use React anyway – shared mental model
- Feed interactions benefit from client-side state
- Keeps learning budget focused on Go backend

**Build tooling:**
- Vite – fast dev server, simple config, good TypeScript support

**Styling:**
- TBD – Tailwind, CSS modules, or vanilla CSS
- Keep it minimal per design philosophy

### Chrome Extension: React + TypeScript

- Popup UI for quick save + optional note
- Shares types/API client with web app where possible
- Keyboard shortcut registration via manifest

### Mobile: React Native + Expo

**Why Expo:**
- Simplifies build/deploy pipeline
- Good support for share extensions (with some native config)
- Matches React knowledge from web

**Note:** iOS share extensions are fiddly. Will require native configuration even with Expo. Tackle this after web + extension are solid.

---

## Hosting Architecture

```
┌─────────────────────────────────────────────┐
│                  Fly.io                     │
├─────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐        │
│  │   API       │    │   Worker    │        │
│  │   (Go)      │    │   (Go)      │        │
│  │   Docker    │    │   Docker    │        │
│  └──────┬──────┘    └──────┬──────┘        │
│         │                  │               │
│         └────────┬─────────┘               │
│                  ▼                         │
│         ┌─────────────┐                    │
│         │  Postgres   │                    │
│         │  (Fly)      │                    │
│         └─────────────┘                    │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│                  Vercel                     │
│  ┌─────────────────────────────────────┐   │
│  └─────────────────────────────────────┘   │
│         Web App (React static build)       │
└─────────────────────────────────────────────┘
```

**Fly.io** (API, Worker, Database):
- Deploy via Dockerfile
- Run API and worker as separate services
- Managed Postgres with automatic backups
- Generous free tier

**Vercel** (Frontend):
- Static React build
- Free tier covers personal projects
- Automatic deploys from Git

### Why this split?

- Go backend needs a container runtime – Fly.io is great for this
- React frontend is static files – Vercel is optimized for this and free
- Keeps concerns separated, each hosted where it fits best

---

## Local Development

**Docker Compose setup:**
- Postgres container
- API container (hot reload via air or similar)
- Worker container
- Frontend runs outside Docker (Vite dev server) for fast HMR

**Environment management:**
- `.env` files for local config
- Fly.io secrets for production

---

## Authentication

**OAuth flow (Google):**
1. Frontend redirects to Google consent screen
2. Google redirects back with auth code
3. Backend exchanges code for tokens
4. Backend creates/updates user, issues JWT
5. JWT stored in httpOnly cookie (web) or secure storage (extension/mobile)

**Self-implemented vs library:**
- Start with self-implemented – OAuth code grant is straightforward
- Adds learning value for Go
- Can add Apple sign-in later with same pattern

---

## Open Decisions

1. **CSS approach** – Tailwind vs CSS modules vs vanilla. Decide when starting frontend.
2. **API documentation** – OpenAPI spec? Useful if we want generated clients.
3. **Monorepo vs separate repos** – Leaning monorepo for simplicity (API, worker, web, extension in one place).

---

## Out of Scope for v1

- Redis or external job queue
- Dedicated search engine (Elasticsearch, Typesense, Meilisearch)
- Kubernetes or complex orchestration
- CI/CD pipeline (manual deploys initially, automate later)
