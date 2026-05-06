# PROJECT_STATE

## Project

Single-player prototype of a cooperative card game inspired by Hearthstone Tavern Brawl “Unite Against Mechazod!”.

Current stack:

- Go
- Pure backend/game engine for now
- No HTTP
- No HTMX
- No templates
- No database
- No UI dependencies inside `internal/game`

## Current stage

Stage 1, steps 1–12 completed.

Current status:

```bash
go test ./...
```

passes.

## Architecture

Core package:

```txt
internal/game
```

This package contains pure domain/game engine logic.

`internal/game` must not import:

- `net/http`
- `html/template`
- `text/template`
- HTMX/UI packages
- database packages

All player actions go through:

```go
func ApplyAction(g *Game, action Action) ([]GameEvent, error)
```

## Implemented features

### Game creation

Implemented:

```go
func NewGame(id string, p1 PlayerConfig, p2 PlayerConfig, seed int64) *Game
```

Creates:

- game ID
- 2 players
- boss
- starting decks
- starting hands
- active player
- active game status
- deterministic shuffle by seed

Game status constants use the `GameStatus...` naming style:

```go
GameStatusCreated
GameStatusActive
GameStatusWon
GameStatusLost
```

## Players

Players have:

- `ID`
- `Name`
- `Health`
- `MaxHealth`
- `Mana`
- `MaxMana`
- `Deck`
- `Hand`
- `Board`
- `Discard`
- `IsCurrent`

Starting health:

```go
StartingPlayerHealth = 30
```

Mana:

- starts at `0/0`
- refresh increases `MaxMana` by 1
- max mana cap is `10`

Implemented:

```go
func RefreshMana(g *Game, playerIndex int)
func SpendMana(g *Game, playerIndex int, amount int) error
```

Mana-related errors include:

```go
ErrNilGame
ErrInvalidPlayerIndex
ErrNotEnoughMana
ErrNegativeManaAmount
```

## Boss

Boss has:

- `ID`
- `Name`
- `Health`
- `MaxHealth`
- `Attack`
- `Armor`

Boss HP is clamped to `0` when damaged below zero.

## Cards

Cards are defined in:

```go
CardRegistry
```

`CardRegistry` is the source of truth for card definitions.

Implemented cards:

### Strike

```txt
ID: strike
Cost: 1
Type: spell
Effect: deal 3 damage to boss
Targeting: does not require manual target selection
```

### Repair

```txt
ID: repair
Cost: 2
Type: spell
Effect: restore 5 health to chosen hero
Targeting: required
Valid targets: hero:0, hero:1
Invalid target: boss
```

### Drone

```txt
ID: drone
Cost: 2
Type: minion
Stats: 2/3
Effect: summon
```

`drone` is playable via `ActionTypePlayCard` and summons a minion.

## Decks

Starting decks use card IDs from `CardRegistry`.

`StartingDeckSize` is still kept.

All starting deck cards must exist in `CardRegistry`.

Implemented:

```go
func NewStartingDeck(ownerID PlayerID) []CardInstance
```

Starting hands are dealt during `NewGame`.

## Draw

Implemented:

```go
func DrawCard(g *Game, playerIndex int) []GameEvent
```

Rules:

- draws top card from `Deck[0]`
- appends it to `Hand`
- creates `EventTypeCardDrawn`
- does not panic on empty deck, invalid player index, or nil game

## Actions

Implemented action types:

```go
ActionTypeEndTurn
ActionTypePlayCard
```

Aliases exist:

```go
ActionEndTurn
ActionPlayCard
```

`Action` includes:

- `Type`
- `PlayerID`
- `CardID`
- `SourceID`
- `TargetID`
- `Target`

`TargetID` uses string IDs:

```txt
hero:0
hero:1
boss
```

## End turn

`ActionTypeEndTurn`:

- only works in `GameStatusActive`
- only active player can end turn
- switches active player
- increments turn
- refreshes mana for the new active player
- draws a card for the new active player
- creates `EventTypeTurnStarted`
- returns and stores generated events

## Play card

`ActionTypePlayCard`:

- only works in `GameStatusActive`
- only active player can play cards
- card must be in active player's hand
- card definition must exist in `CardRegistry`
- player must have enough mana
- for minion cards, board must have free slots (`MaxBoardSize = 7`)
- on success:
  - mana is spent
  - card is removed from hand
  - card effect is applied
  - events are returned
  - events are appended to `g.Events`

Implemented effects:

```go
EffectDamageBoss
EffectHealHero
EffectSummon
```

Unsupported effects return:

```go
ErrUnsupportedCardEffect
```

## Strike behavior

Playing `strike`:

- costs 1 mana
- removes the card from hand
- deals 3 damage to the boss
- creates `EventCardPlayed`
- creates `EventDamage`
- checks game over after damage

## Repair behavior

Playing `repair`:

- costs 2 mana
- requires a target
- valid targets:
  - `hero:0`
  - `hero:1`
- invalid targets:
  - `boss`
  - unknown target IDs
  - empty target
- heals the selected hero for 5
- does not allow hero HP to exceed `MaxHealth`
- removes the card from hand
- creates `EventCardPlayed`
- creates `EventHeal`

Repair-related errors:

```go
ErrTargetRequired
ErrInvalidTarget
```

## Drone behavior

Playing `drone`:

- costs 2 mana
- removes the card from hand
- summons a minion on the active player's board
- creates `EventCardPlayed`
- creates `EventMinionSummoned`
- summoned minion is created from card stats:
  - `Attack = 2`
  - `Health = 3`
  - `MaxHealth = 3`
  - `CanAttack = false` on summon

Board constraints:

- maximum board size is `7` (`MaxBoardSize`)
- when board is full, play is rejected with `ErrBoardFull`
- on `ErrBoardFull`, game state must not mutate (no mana spend, card stays in hand, board unchanged)

## Victory

Implemented:

```go
func CheckGameOver(g *Game) []GameEvent
```

Rules:

- if boss HP `<= 0`, status becomes `GameStatusWon`
- creates `EventGameWon`
- appends event to `g.Events`
- returns the event
- repeated calls do not duplicate the win event
- after `GameStatusWon`, `ApplyAction` rejects new actions with `ErrGameNotActive`

## Targeting

Implemented:

```go
func ValidTargets(g *Game, playerID string, cardInstanceID string) ([]Target, error)
```

Rules:

- finds player by `playerID`
- finds card instance in that player's hand
- finds card definition in `CardRegistry`
- if card does not require target, returns an empty list
- for `repair`, returns `hero:0` and `hero:1`
- does not return `boss` for `repair`

Target contains UI-friendly fields:

- `ID`
- `Type`
- `Kind`
- `PlayerID`
- `BossID`
- `MinionID`
- `OwnerID`
- `DisplayName`

Targeting helpers:

```go
ResolveTarget
ValidateTargetForCard
```

Targeting should use the same source of truth as `ActionPlayCard`:

```go
CardDefinition.Targeting
```

## Events

Current event constants use the `EventType...` naming style, with aliases for convenience.

Important event types:

```go
EventTypeCardDrawn
EventTypeTurnStarted
EventTypeCardPlayed
EventTypeMinionSummoned
EventTypeDamageDealt
EventTypeHeal
EventTypeGameWon
```

Aliases:

```go
EventCardPlayed
EventMinionSummoned
EventDamage
EventHeal
EventGameWon
```

## Errors

Current domain errors include:

```go
ErrNilGame
ErrInvalidPlayerIndex
ErrNotEnoughMana
ErrNegativeManaAmount
ErrGameNotActive
ErrNotYourTurn
ErrUnknownAction
ErrCardNotInHand
ErrUnknownCard
ErrBoardFull
ErrUnsupportedCardEffect
ErrTargetRequired
ErrInvalidTarget
```

## Important invariants

- `internal/game` must stay UI-agnostic.
- All actions must go through `ApplyAction`.
- Error paths must not partially mutate game state.
- `go test ./...` must pass after every step.
- `CardRegistry` is the source of truth for card definitions.
- Starting decks must not reference cards missing from `CardRegistry`.
- `ValidTargets` and `ActionPlayCard` targeting validation should use the same targeting rules.
- Boss HP should not display below `0`.
- Player HP should not exceed `MaxHealth`.
- Summoned minions should not be able to attack on the same turn (`CanAttack = false` on summon).
- Player board size should never exceed `MaxBoardSize`.
- After `GameStatusWon`, no new action should be accepted by `ApplyAction`.
