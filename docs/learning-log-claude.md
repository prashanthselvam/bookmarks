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

## What We Built

✅ Monorepo structure with backend/, web/, extension/ directories
✅ Go module initialized (go.mod)
✅ Multi-stage Dockerfile (builder + alpine runtime)
✅ Docker Compose setup (Postgres + Go API)
✅ Minimal working API endpoint
✅ Database with persistent volume

**Next:** Frontend scaffold, deployments, then Milestone 1 (actual features!)
