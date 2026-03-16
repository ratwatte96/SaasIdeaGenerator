# SaaS Idea Generator

MVP stack for collecting product signals and generating SaaS ideas.

## Services
- PostgreSQL database
- Go API (`/api`)
- Go collector CLI (daily one-shot)
- Next.js frontend dashboard

## Quick start

```bash
cp .env.example .env
make dev
```

Open http://localhost:3000.

## Useful commands

```bash
make migrate-up
make collector-run
make test
make lint
```

## API endpoints
- `GET /api/health`
- `GET /api/ideas?category=&min_demand_score=&competition_level=&limit=&offset=`
- `GET /api/ideas/{id}`
- `GET /api/products?category=&limit=&offset=`
