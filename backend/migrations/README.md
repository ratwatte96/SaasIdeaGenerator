# Database migrations

This project uses SQL migration files in this directory and is compatible with
[`golang-migrate`](https://github.com/golang-migrate/migrate).

## Files

- `0001_init.up.sql`: creates the MVP schema and indexes.
- `0001_init.down.sql`: drops the MVP schema and indexes.

## Competition level thresholds

The initial deterministic thresholds are:

- `low`: 0-10 competitors
- `medium`: 11-30 competitors
- `high`: 31+ competitors

These thresholds are mirrored in code at `backend/models/competition.go`.
