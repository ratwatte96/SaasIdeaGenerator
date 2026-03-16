# SaaS Idea Miner MVP — Implementation-Ready Checklist

This checklist turns the SPEC into a concrete build plan with completion criteria.

---

## 0) Project bootstrap

- [ ] Initialize repository structure:
  - [ ] `/backend`
  - [ ] `/backend/api`
  - [ ] `/backend/models`
  - [ ] `/backend/queries`
  - [ ] `/cmd/collector`
  - [ ] `/frontend`
  - [ ] `/docker`
- [ ] Add root-level `.env.example` with all required configuration keys.
- [ ] Add `README.md` with local setup and run instructions.
- [ ] Add Makefile targets:
  - [ ] `make dev`
  - [ ] `make test`
  - [ ] `make lint`
  - [ ] `make migrate-up`
  - [ ] `make migrate-down`

**Definition of done:** New developer can clone repo and understand how to run app end-to-end.

---

## 1) Database design and migrations

### 1.1 PostgreSQL setup

- [ ] Define DB name, user, password in docker compose and env files.
- [ ] Add migration tooling (Go migrate or equivalent).

### 1.2 Schema

- [ ] Create `products` table:
  - [ ] `id UUID PRIMARY KEY`
  - [ ] `name TEXT NOT NULL`
  - [ ] `description TEXT`
  - [ ] `category TEXT`
  - [ ] `source TEXT NOT NULL`
  - [ ] `source_external_id TEXT` (for upstream dedup)
  - [ ] `upvotes INT DEFAULT 0`
  - [ ] `competitor_count INT DEFAULT 0`
  - [ ] `created_at TIMESTAMPTZ DEFAULT now()`
  - [ ] `updated_at TIMESTAMPTZ DEFAULT now()`
- [ ] Create `ideas` table:
  - [ ] `id UUID PRIMARY KEY`
  - [ ] `idea_text TEXT NOT NULL`
  - [ ] `source_product_id UUID REFERENCES products(id) ON DELETE CASCADE`
  - [ ] `demand_score NUMERIC(10,2) NOT NULL`
  - [ ] `competition_level TEXT NOT NULL`
  - [ ] `created_at TIMESTAMPTZ DEFAULT now()`
- [ ] Add constraints and indexes:
  - [ ] Unique index on `(source, source_external_id)` for products.
  - [ ] Unique index on `(source_product_id, idea_text)` for ideas.
  - [ ] Index on `products(category)`.
  - [ ] Index on `ideas(demand_score DESC)`.
  - [ ] Index on `ideas(competition_level)`.

### 1.3 Competition level thresholds

- [ ] Define deterministic thresholds for `competition_level` from `competitor_count`, e.g.:
  - [ ] `low`: 0–10
  - [ ] `medium`: 11–30
  - [ ] `high`: 31+
- [ ] Document thresholds in code constants and README.

**Definition of done:** Migrations run cleanly on empty DB and produce all tables, constraints, and indexes.

---

## 2) SQLC and query layer

- [ ] Configure SQLC in `/backend/queries/sqlc.yaml`.
- [ ] Add SQL query files for:
  - [ ] Insert/upsert product.
  - [ ] Update competitor counts.
  - [ ] List products with filters.
  - [ ] Insert idea with dedup handling.
  - [ ] List ideas with filters + pagination + sorting.
  - [ ] Get idea by ID with joined source product.
  - [ ] Get related ideas by same category/source product.
- [ ] Generate typed query code.

**Definition of done:** Application code performs DB access only through generated SQLC query methods.

---

## 3) Collector service (`cmd/collector/main.go`)

### 3.1 Runtime and config

- [ ] Build collector as one-shot CLI executable.
- [ ] Add env config:
  - [ ] `DATABASE_URL`
  - [ ] `PRODUCT_HUNT_API_KEY`
  - [ ] request timeout settings
  - [ ] rate limit settings for AlternativeTo (`1 req/sec`)

### 3.2 Product Hunt ingestion

- [ ] Fetch products with fields:
  - [ ] `product_name`
  - [ ] `tagline`
  - [ ] `category`
  - [ ] `launch_date`
  - [ ] `upvotes`
- [ ] Map API payload to `products` table (`description` from tagline).
- [ ] Upsert products using `(source, source_external_id)`.
- [ ] Add logging of counts: fetched, inserted, updated, skipped.

### 3.3 AlternativeTo enrichment

- [ ] Implement scraping client with strict rate limit of 1 request/sec.
- [ ] Scrape/store metadata only (no review text).
- [ ] Match products to AlternativeTo entries:
  - [ ] exact normalized name match
  - [ ] fallback heuristics
  - [ ] unmatched list logging
- [ ] Update `competitor_count` in products.

### 3.4 Demand score + idea generation

- [ ] Compute demand score per product:
  - [ ] `demand_score = upvotes*0.6 + competitor_count*0.4`
- [ ] Generate idea variants from category/product phrase using niche list:
  - [ ] recruiters
  - [ ] sales teams
  - [ ] lawyers
  - [ ] consultants
  - [ ] real estate agents
  - [ ] content creators
- [ ] Compute `competition_level` and store with each idea.
- [ ] Upsert ideas with `(source_product_id, idea_text)` uniqueness.

### 3.5 Scheduling

- [ ] Run collector once daily via container scheduler (cron/host scheduler).
- [ ] Ensure reruns are idempotent (no duplicate products or ideas).

**Definition of done:** One collector run ingests products, enriches competitor counts, computes scores, and writes ideas without duplicates.

---

## 4) API service (Go + Gin)

### 4.1 Service bootstrap

- [ ] Create `/backend/main.go` with Gin router and DB wiring.
- [ ] Mount base path `/api`.
- [ ] Add health endpoint `GET /api/health`.

### 4.2 Endpoints

- [ ] `GET /api/ideas`
  - [ ] Filters: `category`, `min_demand_score`, `competition_level`
  - [ ] Pagination: `limit`, `offset`
  - [ ] Sorting: default `demand_score DESC`
  - [ ] Response includes source product summary
- [ ] `GET /api/ideas/{id}`
  - [ ] Returns idea detail, source product, demand score, competition level
  - [ ] Include related ideas section
- [ ] `GET /api/products`
  - [ ] Returns stored products list

### 4.3 API contracts and error handling

- [ ] Define JSON response schema structs.
- [ ] Validate query params and return 400 for invalid input.
- [ ] Return 404 for unknown idea IDs.
- [ ] Add consistent error response format.

**Definition of done:** Endpoints return expected JSON and support required filtering.

---

## 5) Frontend (Next.js dashboard)

### 5.1 App shell

- [ ] Initialize Next.js app with Tailwind CSS.
- [ ] Create layout:
  - [ ] Sidebar (filters)
  - [ ] Main panel (ideas table/detail)

### 5.2 Idea list page

- [ ] Build ideas table with columns:
  - [ ] Idea
  - [ ] Demand Score
  - [ ] Competition Level
  - [ ] Source Product
- [ ] Add filters:
  - [ ] category
  - [ ] min demand score
  - [ ] competition level
- [ ] Wire filters to API query params.
- [ ] Add loading, empty, and error states.

### 5.3 Idea detail page

- [ ] Create route for idea details.
- [ ] Show:
  - [ ] idea description
  - [ ] source product
  - [ ] demand signals
  - [ ] related ideas

**Definition of done:** User can browse ideas, filter them, and open details from the UI.

---

## 6) Docker and local runtime

- [ ] Create `docker/docker-compose.yml` (or root `docker-compose.yml`) with services:
  - [ ] `database` (PostgreSQL)
  - [ ] `backend` (API)
  - [ ] `collector`
  - [ ] `frontend`
- [ ] Ensure service dependencies and startup order are defined.
- [ ] Add DB migration execution in startup flow.
- [ ] Expose ports:
  - [ ] frontend `3000`
  - [ ] backend `8080` (or documented port)

**Definition of done:** `docker compose up` starts full stack and frontend can query backend.

---

## 7) Testing and validation checklist

### 7.1 Collector tests

- [ ] Unit tests for demand score function.
- [ ] Unit tests for competition level classification.
- [ ] Unit tests for idea template generation and dedup.
- [ ] Integration test for upsert/idempotent rerun behavior.

### 7.2 API tests

- [ ] Endpoint tests for `GET /ideas` filters.
- [ ] Endpoint test for `GET /ideas/{id}` happy path + 404.
- [ ] Endpoint test for `GET /products`.

### 7.3 Frontend tests

- [ ] Render test for idea table.
- [ ] Filter interaction test.
- [ ] Detail page load test.

### 7.4 End-to-end verification

- [ ] Run `docker compose up`.
- [ ] Run one collector cycle.
- [ ] Verify idea rows appear in DB.
- [ ] Verify UI shows list, filters work, detail page works.

**Definition of done:** Acceptance criteria all pass with reproducible local demo.

---

## 8) Acceptance criteria mapping (traceability)

- [ ] **AC1: Products ingested from Product Hunt** → Collector ingestion + DB writes verified.
- [ ] **AC2: Competitor counts scraped from AlternativeTo** → Enrichment and product updates verified.
- [ ] **AC3: Ideas generated automatically** → Post-collector run ideas exist for products.
- [ ] **AC4: API returns idea list** → `GET /api/ideas` works with filters.
- [ ] **AC5: Dashboard displays and filters ideas** → Frontend list + filters verified.

---

## 9) Suggested execution order (sprint plan)

### Sprint 1 (Foundation)
- [ ] Repo bootstrap
- [ ] DB migrations
- [ ] SQLC setup
- [ ] API skeleton + `/health`

### Sprint 2 (Data pipeline)
- [ ] Product Hunt ingestion
- [ ] AlternativeTo enrichment
- [ ] Demand scoring
- [ ] Idea generation

### Sprint 3 (User-facing)
- [ ] Ideas/product API endpoints
- [ ] Next.js list page + filters
- [ ] Next.js detail page

### Sprint 4 (Hardening)
- [ ] Docker compose full stack
- [ ] Tests and bug fixes
- [ ] Demo script and README polishing

---

## 10) MVP non-goals (explicitly defer)

- [ ] Advanced ML idea ranking
- [ ] Authentication/authorization
- [ ] Multi-tenant workspaces
- [ ] Real-time updates/websockets
- [ ] Additional data sources beyond Product Hunt + AlternativeTo

