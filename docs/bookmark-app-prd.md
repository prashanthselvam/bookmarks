# Bookmark App – Product Requirements Document

## Overview

A personal bookmarking system designed to solve the "bookmark graveyard" problem. Unlike browser bookmarks or read-later apps that become cluttered and forgotten, this app actively resurfaces saved content and makes it effortlessly searchable.

**Design philosophy:** Minimal, fast, frictionless. This app removes clutter from your life, not adds it.

---

## Problem Statement

Current bookmarking solutions fail because:

1. **Bookmarks disappear** – saved content goes into folders and is never seen again
2. **Search is broken** – browser bookmark search is basic and doesn't cover content/notes
3. **No context** – you can't attach notes or see why you saved something
4. **Browser-only** – bookmarking is tied to desktop browsers, not mobile apps
5. **Organization overhead** – folders and tags require upfront effort that doesn't pay off

---

## User

Initially a single user (the developer). Architecture should support multiple users for future expansion.

---

## Core Use Cases

### 1. Save a bookmark (capture)

User encounters something worth saving and captures it with minimal friction.

**Flow:**

- Trigger save via browser extension button, keyboard shortcut, or mobile share sheet
- App auto-captures: URL, page title, favicon, timestamp, OG metadata (description, preview image)
- Optional: user adds a note before confirming
- If URL already exists, user is warned and can choose to update or cancel
- Save completes in under 1 second

**Success criteria:** Saving a bookmark should feel instant and require no more than 2 interactions (trigger + confirm).

### 2. Browse the feed (rediscovery)

User opens the app and is presented with a feed that surfaces bookmarks worth revisiting.

**Feed logic:**

- Recent saves (last 7 days) appear at the top
- Below that, resurface candidates: bookmarks not opened in 30+ days, weighted toward those with notes
- Simple chronological + staleness heuristic (no ML/recommendations in v1)

**Success criteria:** Every time you open the app, you see something you'd forgotten about.

### 3. Search bookmarks (retrieval)

User is looking for something specific they saved.

**Flow:**

- Single search bar, prominent placement
- Full-text search across: title, URL, description, user notes
- Results appear as you type (no submit button)
- Filter by date range (last week / last month / custom)

**Success criteria:** If you remember any word from the bookmark or your notes, you can find it in under 5 seconds.

---

## Features by Priority

### P0 – Launch

| Feature | Description |
|---------|-------------|
| Web app | Primary interface for browsing/searching bookmarks |
| Chrome extension | One-click or keyboard shortcut to save current tab, optional note input |
| iOS share sheet | Save URLs from any iOS app via system share menu |
| OAuth authentication | Google and/or Apple sign-in |
| Auto-metadata capture | URL, title, favicon, timestamp, OG image/description |
| Notes on save | Optional free-text note at capture time |
| Feed view | Blended recent + resurfaced bookmarks |
| Full-text search | Across all captured fields and notes |
| Recency filter | Filter search results by time period |
| Duplicate warning | Alert user if URL already saved, option to update or cancel |
| Trash with 30-day retention | Deleted bookmarks recoverable for 30 days |
| Clean minimal UI | Sparse, fast, no visual clutter |

### P1 – Fast Follow

| Feature | Description |
|---------|-------------|
| Android share sheet | Parity with iOS |
| Screenshot capture | Manual trigger to attach screenshot at save time |
| Offline access | Browse and search bookmarks without connectivity |
| Keyboard navigation | Full keyboard support in web app |

### P2 – Future

| Feature | Description |
|---------|-------------|
| Semantic search | Find bookmarks by concept, not just keyword |
| Smart resurfacing | Topic-aware feed based on current interests |
| Multi-user support | Accounts, sharing, collaboration |
| Other browser extensions | Firefox, Safari, Edge |

---

## Non-Goals (v1)

- Folders or tags
- Social features
- Public bookmark sharing
- Read-it-later / article parsing
- Browser history import (revisit later)

---

## Technical Considerations

*(To be detailed in tech spec, but noting for PRD awareness)*

- **Multi-user architecture from day one** – user IDs on all data, even if only one user exists
- **Search infrastructure** – needs to support fast full-text search; evaluate options during tech spec
- **Extension ↔ backend auth** – secure token handling for browser extension
- **Mobile share extension** – iOS and Android have different implementation patterns

---

## Success Metrics

Since this is a personal learning project, success is qualitative:

1. **You actually use it** – this replaces your current bookmarking habit
2. **You rediscover things** – the feed surfaces content you'd forgotten
3. **You can find anything** – search feels instant and complete
4. **It feels good** – the UI is calm, fast, and doesn't add cognitive load

---

## Open Questions

1. What OG metadata fields do we trust? (Some sites have garbage OG tags)
2. How do we handle URLs that require authentication to view? (e.g., paywalled content)
3. Do we want any onboarding, or is the app self-explanatory?
