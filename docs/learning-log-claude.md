# Learning Log - Claude's Notes

*Key concepts covered during Milestone 0 setup*

---

## Go Project Structure

### `cmd/` Directory
- Contains application **entry points** (main.go files)
- Each subdirectory becomes a separate executable
- Example: `cmd/api/main.go` compiles to the API server binary, `cmd/worker/main.go` compiles to the worker binary
- These are thin wrappers that wire things together and start the program

### `internal/` Directory
- **Special directory enforced by Go compiler**
- Code here can ONLY be imported by code in the same module
- Go's way of marking code as "private"
- Prevents external projects from importing your internal packages
- Most of your actual business logic lives here

**Why this matters:** Allows you to refactor internal code freely without worrying about breaking external dependencies.

---

## Docker Fundamentals

### Key Components

**Dockerfile** (Recipe)
- Text file with build instructions
- Defines HOW to build an image

**Docker Engine** (Builder/Runner)
- Program running on your computer
- Reads Dockerfiles, builds images, runs containers

**Docker Image** (Blueprint/Snapshot)
- Saved filesystem + metadata
- Immutable - never changes once built
- Made of stacked layers

**Container** (Running Instance)
- A running process created from an image
- You can run multiple containers from the same image
- Changes while running, but disappear when stopped (unless using volumes)

**Layers** (Build Cache)
- Each Dockerfile instruction creates a new layer
- Layers are cached and reused if instruction hasn't changed
- If ANY layer changes, all layers after it must rebuild

---

## Docker Layer Caching

### How It Works

Each Dockerfile instruction is checked:
- If instruction + inputs are identical to previous build → **cache hit** (reuse layer, instant)
- If anything changed → **cache miss** (rebuild this layer + all layers after it)

### Why Order Matters

**Bad:**
```dockerfile
COPY . .                    # Copies everything (changes frequently)
RUN go mod download         # Runs every time ANY file changes
```

**Good:**
```dockerfile
COPY go.mod* go.sum* ./     # Only copy dependency files (change rarely)
RUN go mod download         # Only reruns when dependencies change
COPY . .                    # Copy source code last (changes frequently)
```

**Key principle:** Put slow, infrequent changes early; fast, frequent changes late.

---

## Multi-Stage Builds

### Why Use Them

**Single-stage build:**
- Final image includes: Go compiler, build tools, source code, binary
- Size: ~800MB
- Security risk: more software = more vulnerabilities

**Multi-stage build:**
- Stage 1: Build the binary (has all the tools)
- Stage 2: Copy ONLY the binary to a minimal image
- Size: ~15MB
- Secure: minimal attack surface

### How It Works

```dockerfile
FROM golang:1.25 AS builder    # Stage 1 - named "builder"
# ... build steps, creates /app/api binary

FROM alpine:latest             # Stage 2 - completely new container
COPY --from=builder /app/api . # Copy binary from Stage 1 layers
CMD ["./api"]
```

**Key insight:** Stage 1 container is deleted, but **layers remain on disk**. Stage 2 can read from those layers.

**Only the last `FROM` becomes the final image** - earlier stages are temporary/intermediate.

---

## Docker Compose

### Purpose
Orchestrates multiple containers that need to work together (API + Database).

### Key Concepts

**Services** - One service = one container (or group of identical containers)

**Volumes** - Persist data outside containers
- Without volumes: data deleted when container stops
- With volumes: data saved to disk, survives restarts
- Two declarations needed:
  - Top-level `volumes:` → declares volume exists
  - Service-level `volumes:` → mounts volume into container

**Networks** - Containers can find each other by service name
- API connects to database at hostname `db` (the service name)
- Docker Compose creates network automatically

**Environment Variables** - Pass configuration to containers

**depends_on** - Controls startup order (starts db before api)

### Volume Syntax

```yaml
services:
  db:
    volumes:
      - postgres_data:/var/lib/postgresql/data
        └─ volume name  └─ path inside container

volumes:
  postgres_data:  # Declares the volume
```

---

## Static vs Dynamic Linking

### The Problem We Hit

Built Go binary on Debian (golang:1.25 image) → expects glibc libraries
Ran binary on Alpine → has musl libraries, not glibc
Binary couldn't find libraries → crashed with "no such file or directory"

**Misleading error:** It found the binary file, but couldn't find the **libraries** it needed.

### The Solution: `CGO_ENABLED=0`

```dockerfile
RUN CGO_ENABLED=0 go build -o api ./cmd/api
```

**What this does:**
- Disables C interop (cgo)
- Forces a **statically-linked** binary (all dependencies included)
- Doesn't need glibc, musl, or any external C libraries
- Works on any Linux distro

**Trade-off:** Can't use packages that require cgo (like some SQLite drivers), but pure Go code works perfectly.

---

## Go HTTP Server Basics

### Minimal Server

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from API!")
    })

    http.ListenAndServe(":8080", nil)
}
```

**Key parts:**
- `http.HandleFunc` - registers a handler for a URL path
- `w http.ResponseWriter` - write response back to client
- `r *http.Request` - incoming request data
- `http.ListenAndServe(":8080", nil)` - start server on port 8080
  - **Must include the colon!** `:8080` not `8080`
  - `:8080` means "listen on all network interfaces"

---

## Common Docker Commands

### Docker Compose Workflow

```bash
# Start in background (detached)
docker compose up -d

# Start + rebuild after code changes
docker compose up -d --build

# Stop everything (keeps volumes/data)
docker compose down

# Stop + delete volumes (wipes database!)
docker compose down -v

# Watch logs
docker compose logs -f
docker compose logs -f api    # specific service

# Check status
docker compose ps

# Restart service
docker compose restart api

# Shell into container
docker compose exec api sh
docker compose exec db psql -U bookmarks bookmarks_dev
```

### Cleanup Commands

```bash
# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Nuclear option: remove everything unused
docker system prune -a --volumes
```

---

## YAML Syntax Quick Reference

**Maps (key-value pairs)** - use colons
```yaml
key: value
another_key: another_value
```

**Lists** - use dashes
```yaml
- item1
- item2
```

**Maps with list values** - combine them
```yaml
ports:
  - "8080:8080"
  - "9090:9090"
```

**Common mistake:** Forgetting the dash for list items or colon for map keys.

---

## Port Mapping

### Syntax
```yaml
ports:
  - "host_port:container_port"
```

**Example:** `"8080:8080"`
- Application listens on port 8080 **inside container**
- Docker forwards port 8080 on **your Mac** to the container
- Access via `localhost:8080` from your Mac

**Different ports:** `"3000:8080"`
- Inside container: still 8080
- From your Mac: access via `localhost:3000`

---

## Frontend Scaffolding with Vite

### What is Vite?

**Vite** is a modern build tool for frontend projects:
- Extremely fast dev server (starts instantly)
- Hot Module Replacement (HMR) - changes appear immediately
- Modern alternative to Create React App
- Built by Evan You (creator of Vue.js)

### Creating a Vite Project

```bash
npm create vite@latest
```

**Prompts:**
- Project name: `web`
- Framework: React
- Variant: TypeScript

**What gets created:**
- `src/` - React components
- `dist/` - production build output (after `npm run build`)
- `vite.config.ts` - configuration
- `package.json` - dependencies

### Key Commands

```bash
npm run dev      # Start dev server (usually port 5173)
npm run build    # Production build → creates dist/
npm run preview  # Preview production build locally
```

### Why Vite over Create React App

- ✅ Much faster (instant dev server start)
- ✅ Better DX (developer experience)
- ✅ Actively maintained (CRA is deprecated)
- ✅ Works great with TypeScript out of the box

---

## CORS (Cross-Origin Resource Sharing)

### What is CORS?

**Browser security rule:** JavaScript on one website cannot access data from another website (different origin) without permission.

**Origin = protocol + domain + port**

Examples:
- `http://localhost:5173` (frontend)
- `http://localhost:8080` (backend)
- These are **different origins** (different ports!)

### Why CORS Exists

**The threat it prevents:**
```
User visits evil.com
→ evil.com's JavaScript tries to call yourbank.com/transfer
→ Browser automatically includes user's cookies
→ Without CORS: evil.com could steal money
→ With CORS: Browser blocks the response ✅
```

### How CORS Works

**Frontend makes request:**
```javascript
fetch('http://localhost:8080')
```

**Backend must respond with header:**
```
Access-Control-Allow-Origin: http://localhost:5173
```

**Browser checks:**
- Response has correct CORS header? → Allow JavaScript to see response ✅
- Missing or wrong header? → Block response, show CORS error ❌

### Implementing CORS in Go

```go
w.Header().Set("Access-Control-Allow-Origin", origin)
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
w.Header().Add("Vary", "Origin")
```

**Important headers:**
- `Access-Control-Allow-Origin` - which origins can access
- `Vary: Origin` - tells CDNs to cache differently per origin (prevents cache poisoning)

### Why curl/Postman Don't Need CORS

**CORS is enforced by browsers, not servers.**

- curl/Postman: No origin concept, just making direct HTTP requests
- Browser: Has origin (the webpage JavaScript runs on), enforces CORS for security
- Server sends same response to both - browser decides if JavaScript can see it

### Can JavaScript Fake the Origin Header?

**No.** The `Origin` header is set by the browser (C++ code), not JavaScript.

**Forbidden headers** JavaScript cannot modify:
- `Origin`
- `Host`
- `Cookie`
- `Referer`

JavaScript runs in a sandbox - the browser has ultimate control.

---

## Pattern Matching for CORS Origins

### The Problem

Cloudflare Pages generates new preview URLs for each deployment:
- `https://c108ff7e.bookmarks-web-qbt.pages.dev`
- `https://a1b2c3d4.bookmarks-web-qbt.pages.dev`

Can't hardcode all preview URLs!

### The Solution

Check if origin **matches a pattern:**

```go
func isOriginAllowed(origin string) bool {
    // Check exact matches
    for _, allowed := range allowedOrigins {
        if origin == allowed {
            return true
        }
    }

    // Allow any Cloudflare Pages preview URL
    if strings.HasSuffix(origin, ".bookmarks-web-qbt.pages.dev") &&
       strings.HasPrefix(origin, "https://") {
        return true
    }

    return false
}
```

**This allows:**
- ✅ Production: `https://bookmarks-web-qbt.pages.dev`
- ✅ Any preview: `https://*.bookmarks-web-qbt.pages.dev`
- ✅ Local dev: `http://localhost:5173`

**Security:** Safe because only YOU can deploy to your Cloudflare Pages project.

---

## Deploying to Fly.io

### What is Fly.io?

Platform for running Docker containers globally:
- Deploys containers to edge locations worldwide
- Production-ready (not just for side projects!)
- Simple CLI-based workflow
- Generous free tier

### Deployment Steps

**1. Install CLI:**
```bash
brew install flyctl
```

**2. Login:**
```bash
flyctl auth login
```

**3. Initialize app:**
```bash
cd backend/
flyctl launch
```

Prompts:
- App name (becomes your URL)
- Region (pick closest to you)
- PostgreSQL? No (for now)
- Deploy now? No (configure first)

Creates `fly.toml` config file.

**4. Deploy:**
```bash
flyctl deploy
```

**What happens:**
1. Builds Docker image (using your Dockerfile)
2. Pushes to Fly.io registry
3. Deploys as container
4. Gives you URL: `https://your-app.fly.dev`

### Key fly.toml Settings

```toml
[http_service]
  internal_port = 8080      # Must match your Go app's port
  force_https = true        # Redirect HTTP → HTTPS
```

### Redeploying After Changes

```bash
flyctl deploy
```

That's it! Fly.io rebuilds and redeploys automatically.

---

## Deploying to Cloudflare Pages

### What is Cloudflare Pages?

Static site hosting on Cloudflare's global CDN:
- Optimized for frontend frameworks (React, Vue, etc.)
- Unlimited bandwidth on free tier
- Global CDN (200+ locations)
- Built-in preview deployments

### Why Cloudflare Pages?

Compared to alternatives:
- **vs Vercel:** No ethical concerns, truly unlimited free tier, great learning opportunity
- **vs Netlify:** Better free tier (no bandwidth limits)
- **vs GitHub Pages:** Automatic builds, better DX

### Deployment Steps

**1. Build your app:**
```bash
cd web/
npm run build
```

Creates `dist/` folder with static files.

**2. Install Wrangler:**
```bash
npm install -g wrangler
```

**3. Login:**
```bash
wrangler login
```

**4. Deploy:**
```bash
npx wrangler pages deploy dist
```

Prompts:
- Project name (becomes URL)
- Branch name (usually `main`)

**Your URL:** `https://your-project.pages.dev`

### Preview Deployments

Each time you deploy, Cloudflare creates:
- **Production URL:** `https://your-project.pages.dev` (stable)
- **Preview URL:** `https://hash.your-project.pages.dev` (unique per deployment)

**Why previews are useful:**
- Test changes before making them production
- Share specific deployment with others
- Each preview is a full deployment (not just a link)

### Redeploying

```bash
npm run build
npx wrangler pages deploy dist
```

New preview URL created, production URL unchanged (unless you specifically deploy to production).

---

## Frontend Tooling Landscape

Understanding the layers:

### Layer 1: JavaScript Runtime
**What executes your code**
- **Node.js** (standard, what we're using)
- Bun (new, faster alternative)
- Deno (another alternative)

### Layer 2: Build Tool
**Bundles code, runs dev server**
- **Vite** (modern, fast - what we're using)
- Webpack (older, slower)
- esbuild (very fast, low-level)

### Layer 3: Framework
**How you structure your app**
- **React** (what we're using)
- Next.js (React + SSR + routing)
- Vue, Svelte (alternatives to React)

**Our stack:** Node.js + Vite + React + TypeScript

**Why this works:** Vite for fast builds, React for UI, separate Go backend for API logic. Clean separation of concerns.

---

## Platform Comparison

### Fly.io (Backend)
- Docker-native
- Global deployment
- Good for Go/any containerized app
- $0-5/month for small projects

### Cloudflare Pages (Frontend)
- Static sites only
- Unlimited bandwidth (free)
- Global CDN
- Preview deployments

### Why Not One Platform?
- Each optimized for different use cases
- Backend needs: Database, long-running processes, Docker
- Frontend needs: Fast global delivery, CDN, static hosting
- Separation = better performance + lower cost

---

## Development vs Production Environments

### Local Development

**Frontend:**
- Runs natively on your Mac (not in Docker)
- `npm run dev` on port 5173
- Fast hot reload

**Backend:**
- Runs in Docker Compose
- Postgres + API together
- Mirrors production architecture

**Why not Docker for frontend dev?**
- Faster file watching
- Better HMR (hot module replacement)
- Simpler workflow

### Production

**Frontend:**
- Static build (`npm run build`)
- Deployed to Cloudflare Pages
- Served from global CDN

**Backend:**
- Docker container on Fly.io
- Postgres database (to be added)
- Same Dockerfile as local

**Environment variables:**
- Local: `http://localhost:8080`
- Production: `https://your-app.fly.dev`

---

## What We Built (Milestone 0 Complete!)

✅ **Local Development:**
- Monorepo structure (backend, web, extension)
- Go module + multi-stage Dockerfile
- Docker Compose (Postgres + API)
- React + Vite frontend
- CORS configured for local + production

✅ **Production Deployment:**
- Backend: https://backend-empty-sun-8345.fly.dev
- Frontend: https://bookmarks-web-qbt.pages.dev
- End-to-end connectivity working
- Pattern-based CORS for preview URLs

✅ **Key Learnings:**
- Docker fundamentals
- CORS and browser security
- Modern deployment platforms
- Development workflow (deploy early and often!)

**Next:** Milestone 1 - Build actual features (database schema, API endpoints, CRUD operations)
