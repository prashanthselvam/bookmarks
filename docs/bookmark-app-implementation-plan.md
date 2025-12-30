# Bookmark App – Implementation Plan

## Principles

- **Vertical slices** – Each milestone delivers working functionality, not just a layer (e.g., "API done, frontend later")
- **Core loop first** – Get saving and viewing bookmarks working end-to-end before adding polish
- **Defer complexity** – Mobile, extensions, and advanced features come after the core is solid
- **Each milestone is demo-able** – You can show something working at each stage

---

## Milestone 0: Project Setup

**Goal:** Development environment ready, empty app deploys successfully.

**Tasks:**

- Initialize monorepo structure
- Set up Go module for backend
- Create basic Dockerfile for API
- Set up Docker Compose (Postgres + API)
- Create React + Vite frontend scaffold
- Verify local dev workflow works (API + frontend + DB)
- Deploy empty API to Fly.io
- Deploy empty frontend to Vercel
- Confirm end-to-end connectivity (frontend can hit deployed API)

**Deliverable:** "Hello world" response from deployed API, blank React app deployed.

**Estimated effort:** Half day

---

## Milestone 1: Core Data Model + Basic API

**Goal:** Can create and retrieve bookmarks via API (no auth, no UI yet).

**Tasks:**

- Design and create Postgres schema (users, bookmarks tables)
- Set up sqlc and generate Go types
- Implement endpoints:
  - `POST /bookmarks` – create bookmark (URL, title, note)
  - `GET /bookmarks` – list all bookmarks
  - `GET /bookmarks/{id}` – get single bookmark
- Add URL hash generation for duplicate detection
- Write basic request validation
- Test with curl or Postman
- Deploy and verify on Fly.io

**Deliverable:** Can create and list bookmarks via API calls.

**Estimated effort:** 1 day

---

## Milestone 2: Web App – View Bookmarks

**Goal:** See your bookmarks in a browser.

**Tasks:**

- Set up API client in React (fetch wrapper or lightweight library)
- Build bookmarks list component
- Display title, URL, note, timestamp
- Basic styling (minimal, clean)
- Handle loading and empty states
- Connect to deployed API (CORS setup)
- Deploy frontend

**Deliverable:** Open web app, see list of bookmarks (manually created via API).

**Estimated effort:** 1 day

---

## Milestone 3: Web App – Create Bookmarks

**Goal:** Save bookmarks from the web app.

**Tasks:**

- Build "add bookmark" form (URL, optional note)
- Implement form submission to API
- Add optimistic UI or loading state
- Implement duplicate detection (API returns warning, frontend shows message)
- Auto-fetch title from URL (client-side for now, simple fetch)
- Deploy

**Deliverable:** Can add bookmarks from the web UI and see them appear in the list.

**Estimated effort:** 1 day

---

## Milestone 4: Authentication

**Goal:** Bookmarks are tied to a user, login required.

**Tasks:**

- Set up Google OAuth app (Google Cloud Console)
- Implement OAuth endpoints in Go:
  - `GET /auth/google` – redirect to Google
  - `GET /auth/google/callback` – handle callback, issue JWT
- Create users table entries on first login
- Implement JWT middleware for protected routes
- Add user_id to all bookmark queries
- Build login UI in React
- Handle auth state (logged in / logged out)
- Secure cookie handling
- Deploy and test full flow

**Deliverable:** Must log in to see/create bookmarks. Each user only sees their own.

**Estimated effort:** 1.5 days

---

## Milestone 5: Background Metadata Fetching

**Goal:** Bookmarks automatically get rich metadata (OG image, description, favicon).

**Tasks:**

- Create jobs table in Postgres
- Implement job queue logic (enqueue, claim, complete)
- Build worker entrypoint (same codebase, `cmd/worker`)
- Implement metadata fetcher:
  - Fetch URL
  - Parse OG tags (title, description, image)
  - Extract favicon
  - Update bookmark record
- Enqueue job on bookmark creation
- Update Dockerfile / Docker Compose for worker
- Deploy worker to Fly.io as separate service
- Update frontend to display metadata (image, description)

**Deliverable:** Save a URL, see it appear with preview image and description after a few seconds.

**Estimated effort:** 1.5 days

---

## Milestone 6: Search

**Goal:** Find bookmarks by keyword.

**Tasks:**

- Add tsvector column to bookmarks table
- Create GIN index for full-text search
- Add pg_trgm for fuzzy URL matching
- Implement `GET /bookmarks/search?q=...` endpoint
- Build search UI:
  - Search input with debounce
  - Results as you type
  - Highlight matching terms (optional)
- Add recency filter (last week / month / all time)
- Deploy

**Deliverable:** Type in search box, see matching bookmarks instantly.

**Estimated effort:** 1 day

---

## Milestone 7: Feed + Resurfacing

**Goal:** Homepage shows a feed that resurfaces old bookmarks.

**Tasks:**

- Add `last_opened_at` column (updated when user clicks a bookmark)
- Implement feed endpoint `GET /bookmarks/feed`:
  - Recent (last 7 days)
  - Resurface candidates (not opened in 30+ days, prioritize those with notes)
- Track "opened" event when user clicks bookmark link
- Build feed UI (distinct from search results)
- Deploy

**Deliverable:** Open app, see mix of recent and forgotten bookmarks.

**Estimated effort:** 1 day

---

## Milestone 8: Soft Delete + Trash

**Goal:** Deleted bookmarks are recoverable.

**Tasks:**

- Add `deleted_at` column
- Implement soft delete on `DELETE /bookmarks/{id}`
- Implement `GET /bookmarks/trash` endpoint
- Implement `POST /bookmarks/{id}/restore`
- Implement `DELETE /bookmarks/{id}/permanent`
- Add scheduled job to purge trash > 30 days
- Build trash UI (view, restore, permanent delete)
- Deploy

**Deliverable:** Delete a bookmark, find it in trash, restore it.

**Estimated effort:** 0.5 day

---

## Milestone 9: Chrome Extension

**Goal:** Save bookmarks with one click from the browser.

**Tasks:**

- Set up extension project structure (manifest v3)
- Build popup UI (React):
  - Shows current tab URL/title
  - Optional note field
  - Save button
- Implement extension ↔ API auth:
  - Extension opens login page if not authenticated
  - Token stored in extension storage
- Handle duplicate warning in popup
- Keyboard shortcut for quick save
- Test locally (load unpacked extension)
- Package for distribution

**Deliverable:** Click extension icon, save current page, see it in web app.

**Estimated effort:** 1.5 days

---

## Milestone 10: iOS Share Extension

**Goal:** Save bookmarks from any iOS app.

**Tasks:**

- Set up Expo project for mobile
- Configure iOS share extension
- Build share UI (minimal – URL, note, save button)
- Handle auth (shared keychain with main app, or web login flow)
- Implement API call from share extension
- Test on physical device
- Build main app shell (even if minimal – needed for auth)

**Deliverable:** Share a link from Safari/Instagram, save to bookmarks.

**Estimated effort:** 2 days (share extensions are fiddly)

---

## Post-MVP Milestones (P1)

These come after the core is solid:

| Milestone | Description | Effort |
|-----------|-------------|--------|
| Android share extension | Parity with iOS | 1.5 days |
| Offline support | Cache bookmarks locally, sync when online | 2 days |
| Screenshot capture | Manual screenshot attachment on save | 1 day |
| Keyboard navigation | Full keyboard support in web app | 0.5 day |
| Import existing bookmarks | Chrome bookmark import | 1 day |

---

## Milestone Dependency Graph

```
M0 (Setup)
 │
 ▼
M1 (API)
 │
 ├──────────────┐
 ▼              ▼
M2 (View)      M5 (Metadata Worker) ←── can happen in parallel
 │              │
 ▼              │
M3 (Create) ◄───┘
 │
 ▼
M4 (Auth)
 │
 ├──────────────┬──────────────┐
 ▼              ▼              ▼
M6 (Search)   M7 (Feed)      M8 (Trash) ←── can happen in parallel
 │              │              │
 └──────────────┴──────────────┘
                │
                ▼
        M9 (Chrome Extension)
                │
                ▼
        M10 (iOS Share Extension)
```

---

## Suggested Approach with Claude Code

For each milestone:

1. **Share the milestone scope** – Copy the tasks from this plan
2. **Work incrementally** – One task at a time, test before moving on
3. **Commit at checkpoints** – Working code gets committed, even if rough
4. **Ask for explanations** – Since you're learning Go, have Claude Code explain patterns
