# Mechazod Card Game

Cooperative card game backend prototype inspired by Hearthstone Tavern Brawl "Unite Against Mechazod!".

## Current stage

Stage 1, steps 1-11 completed.

Implemented in `internal/game`:

- game initialization with deterministic shuffle by seed
- turn flow (`EndTurn`) via `ApplyAction`
- playing spell cards (`strike`, `repair`)
- target validation (`ValidTargets`)
- win condition handling (`GameStatusWon`)

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
