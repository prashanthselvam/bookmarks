# Bookmark App – Engineering Architecture

## High-Level System Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                            CLIENTS                                  │
├─────────────────┬─────────────────┬─────────────────────────────────┤
│   Web App       │ Chrome Extension│   Mobile Share Extension        │
│   (SPA)         │                 │   (iOS / Android)               │
└────────┬────────┴────────┬────────┴────────────────┬────────────────┘
         │                 │                         │
         └─────────────────┼─────────────────────────┘
                           │
                           ▼
                 ┌───────────────────┐
                 │   API Gateway /   │
                 │   Load Balancer   │
                 └─────────┬─────────┘
                           │
                           ▼
                 ┌───────────────────┐
                 │   Backend API     │
                 │   (Stateless)     │
                 └─────────┬─────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
┌─────────────────┐ ┌─────────────┐ ┌─────────────────┐
│  Primary DB     │ │ Search      │ │ Background      │
│  (Relational)   │ │ Index       │ │ Job Queue       │
└─────────────────┘ └─────────────┘ └─────────────────┘
                                            │
                                            ▼
                                   ┌─────────────────┐
                                   │ Metadata Fetch  │
                                   │ Worker          │
                                   └─────────────────┘
```

---

## System Components

### 1. Clients

**Web App (SPA)**

- Primary interface for feed browsing, search, and bookmark management
- Talks to backend via REST or GraphQL
- Handles OAuth flow initiation

**Chrome Extension**

- Popup UI for quick save + optional note
- Communicates with backend API (needs auth token storage)
- Keyboard shortcut listener

**Mobile Share Extension (iOS / Android)**

- Minimal UI triggered from system share sheet
- Sends URL + optional note to backend
- Needs to handle auth (likely stored token or OAuth handoff)

---

### 2. API Gateway / Load Balancer

- Single entry point for all client requests
- SSL termination
- Rate limiting (protect against abuse, even if single-user for now)
- Routes requests to backend

For a solo project this might just be your hosting provider's built-in routing, but it's worth calling out as a logical layer.

---

### 3. Backend API (Stateless)

Core of the system. Handles all business logic.

**Responsibilities:**

- Authentication (OAuth token validation, session management)
- Bookmark CRUD operations
- Duplicate detection
- Feed generation (recent + resurfacing logic)
- Search query handling
- Trash management (soft delete, 30-day purge)

**Key design choices:**

- Stateless – no server-side sessions, all state in DB or client tokens
- RESTful API (GraphQL is overkill for this scope)
- Multi-tenant data model from day one (user_id on all tables)

**Endpoints (rough sketch):**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /auth/google | OAuth callback |
| GET | /bookmarks/feed | Get feed (recent + resurfaced) |
| GET | /bookmarks/search | Full-text search with filters |
| POST | /bookmarks | Create bookmark |
| GET | /bookmarks/:id | Get single bookmark |
| PATCH | /bookmarks/:id | Update bookmark (e.g., add note) |
| DELETE | /bookmarks/:id | Soft delete (move to trash) |
| GET | /bookmarks/trash | List trashed bookmarks |
| POST | /bookmarks/:id/restore | Restore from trash |
| DELETE | /bookmarks/:id/permanent | Hard delete |

---

### 4. Primary Database (Relational)

Relational makes sense here because:

- Data is structured and predictable (bookmarks, users)
- We need strong consistency (duplicate detection, trash state)
- Multi-user support later means relational integrity matters

**Core tables:**

```
users
  - id (PK)
  - email
  - name
  - oauth_provider
  - oauth_id
  - created_at

bookmarks
  - id (PK)
  - user_id (FK)
  - url
  - url_hash (for fast duplicate lookup)
  - title
  - description
  - og_image_url
  - favicon_url
  - note
  - created_at
  - updated_at
  - last_opened_at (for resurfacing logic)
  - deleted_at (soft delete)
```

**Indexes needed:**

- `user_id` + `url_hash` (duplicate detection)
- `user_id` + `deleted_at` + `created_at` (feed queries)
- `user_id` + `deleted_at` + `last_opened_at` (resurfacing queries)

---

### 5. Search Index

Full-text search across title, URL, description, and notes.

**Options spectrum:**

- **Database built-in** – Most relational DBs have full-text search. Simpler ops, might be good enough for thousands of bookmarks.
- **Dedicated search engine** – More powerful (fuzzy matching, relevance tuning), but adds infrastructure complexity.

**Recommendation:** Start with database-native full-text search. If it feels slow or limited, migrate to a dedicated engine later. For a single user with even 10k bookmarks, DB-native will likely be fine.

**Search index contains:**

- Bookmark ID
- Title
- URL
- Description
- Note
- User ID (for filtering)

---

### 6. Background Job Queue + Worker

Some operations shouldn't block the save flow.

**Why we need it:**

- Fetching OG metadata can be slow (network request to the bookmarked URL)
- We want save to feel instant (<1 sec)

**Flow:**

1. User saves bookmark → API writes URL + timestamp to DB immediately, returns success
2. API enqueues "fetch metadata" job
3. Worker picks up job, fetches OG tags, updates bookmark record
4. Client can poll or receive push update (or just see it on next load)

**Job types:**

- `fetch_metadata` – Scrape OG tags from URL
- `purge_trash` – Scheduled job to hard-delete bookmarks in trash > 30 days
- (P1) `capture_screenshot` – Headless browser screenshot

---

### 7. File/Blob Storage (P1)

Only needed for screenshot capture feature.

- Store screenshot images
- Serve via CDN for fast loading
- Can defer this entirely until P1

---

## Data Flow Examples

**Saving a bookmark (happy path):**

```
1. User clicks extension button
2. Extension sends POST /bookmarks { url, title, note? }
3. API checks for duplicate (url_hash + user_id)
4. If duplicate → return warning
5. If new → insert bookmark row, enqueue metadata job, return 201
6. Worker fetches OG data, updates row
7. Next time user loads feed, bookmark appears with rich metadata
```

**Searching:**

```
1. User types in search bar
2. Client sends GET /bookmarks/search?q=react&period=last_month
3. API queries search index with user_id filter
4. Returns ranked results
5. Client renders as user types (debounced)
```

**Feed generation:**

```
1. Client requests GET /bookmarks/feed
2. API runs two queries:
   a. Recent: WHERE created_at > 7 days ago ORDER BY created_at DESC
   b. Resurface: WHERE last_opened_at < 30 days ago ORDER BY (has_note, last_opened_at)
3. Merge and return
```

---

## Auth Architecture

**OAuth flow (Google example):**

```
1. User clicks "Sign in with Google" on web app
2. Redirect to Google OAuth consent screen
3. Google redirects back with auth code
4. Backend exchanges code for tokens, extracts user info
5. Backend creates/updates user record, issues app JWT
6. JWT stored in client (httpOnly cookie for web, secure storage for extension/mobile)
7. All subsequent API calls include JWT
```

**Extension auth:**

- Extension opens a tab to web app login
- After login, web app passes token to extension via messaging or redirect
- Extension stores token securely

**Mobile share extension auth:**

- Share extension accesses token from shared app group storage (iOS) or account manager (Android)
- Requires companion app or initial web login

---

## Open Architecture Questions

1. **Real-time updates?** – If you save via extension, should the web app update live? Probably overkill for v1, but websockets or SSE could add this later.

2. **Search latency target?** – "Results as you type" implies <200ms. Need to validate DB full-text can hit this.

3. **Metadata fetch reliability** – Some sites block scrapers. Do we retry? Fall back to just URL/title? Show "metadata unavailable" state?
