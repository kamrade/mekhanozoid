# Mechazod Card Game

Cooperative card game backend prototype inspired by Hearthstone Tavern Brawl "Unite Against Mechazod!".

## Current stage

Stage 1, Step 1: basic Go project structure.

## Run

```bash
go run ./cmd/server
```

## Test
```bash
go test ./...
```

## Architecture rule
The internal/game package contains pure game/domain logic.

It must not import:

- net/http
- HTML/template packages
- UI/HTMX-related packages
- database packages

