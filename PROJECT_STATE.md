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

Stage 1, steps 1–15 and 17–20 completed (step 16 intentionally skipped).
Stage 2 started: steps 21-23 completed.

Step 14 ("refresh minions at turn start") was re-verified on 2026-05-08:

- attacking minions recover `CanAttack = true` on their owner's next turn
- ally minions are not refreshed during the other player's turn

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
ActionTypeAttack
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

## Attack behavior

`ActionTypeAttack`:

- only works in `GameStatusActive`
- only active player can attack
- only minions on active player's board can attack
- for current step, the only valid target is `boss`
- attacking minion must have `CanAttack = true`
- boss takes damage equal to minion `Attack`
- after attack, minion gets `CanAttack = false`
- same minion cannot attack twice in the same turn
- events are returned and appended to `g.Events`
- game-over check runs after damage

Attack-related errors:

```go
ErrMinionNotFound
ErrMinionCantAttack
ErrInvalidTarget
ErrNotYourTurn
ErrGameNotActive
```

## Turn-start minion refresh

At the start of a player's turn:

- only that player's minions are refreshed
- refreshed minions get `CanAttack = true`
- the other player's minions are not changed
- implementation is part of the existing `ActionTypeEndTurn` flow

## Boss abilities

At the start of each player turn, the boss applies exactly one ability.
This is integrated into the existing `ActionTypeEndTurn` turn-start flow.

Implemented boss abilities:

- `Zap Heroes`
  - deals `2` damage to both heroes
- `Bomb Salvo`
  - deals `2` damage to one random valid target
  - valid targets: living heroes and living minions
  - invalid target: boss
- `Overclock`
  - increases boss `Attack` by `1`

Boss ability behavior:

- emits `EventTypeBossAbility` (`boss_ability`)
- mutates game state according to selected ability
- random ability selection and random `Bomb Salvo` target are deterministic for the same game seed/state
- minion damage reduces minion health (dead-minion cleanup remains part of later steps)

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

## Defeat

Implemented:

- if any hero has `Health <= 0`, game status becomes `GameStatusLost`
- loss detection is part of `CheckGameOver`
- `CheckGameOver` emits and stores `EventTypeGameLost` (`game_lost`)
- repeated calls after terminal status do not duplicate `game_lost`
- after `GameStatusLost`, `ApplyAction` rejects new actions with `ErrGameNotActive`
- precedence rule: loss has priority over win when both conditions are true in the same resolution

## Dead minion cleanup

Implemented:

```go
func CleanupDeadMinions(g *Game) []GameEvent
```

Rules:

- iterates through both players' boards
- removes every minion with `Health <= 0`
- keeps living minions on board
- emits one `EventTypeMinionDied` (`minion_died`) per removed minion
- repeated cleanup calls do not duplicate events for already removed minions

Integration:

- cleanup runs after boss abilities (including minion damage from `Bomb Salvo`)
- dead minions are excluded from minion target lists
- dead minions are rejected by target validation

## Full scenario test

Implemented:

- `TestGameCanBePlayedUntilWin`

Scenario guarantees:

- game is created via `NewGame` with fixed seed
- uses `ApplyAction` for player actions
- includes successful `ActionPlayCard` actions
- includes `ActionEndTurn` (turn alternation)
- boss ability logic is triggered in turn-start flow
- boss HP reaches `0` (with clamp behavior)
- final status becomes `GameStatusWon`

Additional scenario:

- `TestGameCanBeLost`
  - deterministic loss setup via fixed seed
  - triggers boss damage through normal `ActionEndTurn` turn-start flow
  - reduces hero HP to `0` (clamped)
  - asserts `GameStatusLost`
  - asserts `EventTypeGameLost` is stored in `g.Events`
  - asserts further actions are rejected with `ErrGameNotActive`

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
EventTypeAttack
EventTypeMinionSummoned
EventTypeMinionDied
EventTypeDamageDealt
EventTypeHeal
EventTypeBossAbility
EventTypeGameWon
EventTypeGameLost
```

Aliases:

```go
EventCardPlayed
EventAttack
EventMinionSummoned
EventMinionDied
EventDamage
EventHeal
EventBossAbility
EventGameWon
EventGameLost
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
ErrMinionNotFound
ErrMinionCantAttack
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
- Minions should be refreshed at the start of their owner's turn (`CanAttack = true`).
- At the start of each turn, the boss should apply exactly one deterministic ability and emit `boss_ability`.
- Minions with `Health <= 0` should be removed and emit `minion_died`.
- Player board size should never exceed `MaxBoardSize`.
- After `GameStatusWon`, no new action should be accepted by `ApplyAction`.
- After any hero reaches `0` health, game should transition to `GameStatusLost` and reject new actions.
- A deterministic end-to-end scenario test should validate that the game can be played until win.
- A deterministic end-to-end scenario test should validate that the game can be lost and blocks further actions.

## Dev console runner

Implemented:

- `cmd/devgame/main.go`

Runner behavior:

- creates a new game with fixed seed
- prints readable current state:
  - player 1 HP
  - player 2 HP
  - boss HP
  - active player / turn
  - active player mana
  - active player hand
  - active player board
  - available actions
  - recent events
- supports commands:
  - `hand` (prints active hand with index, name, cost, type, target hint)
  - `play <handIndex> [targetID]` (calls `ApplyAction` with `ActionPlayCard`)
  - `attack <minionIndex> boss` (calls `ApplyAction` with `ActionAttack`)
  - `end` / `end turn` (calls `ApplyAction` with `ActionEndTurn`)
  - `state` (reprints full state)
  - `quit` / `exit` (clean shutdown)
- prints action errors returned by `ApplyAction`
- validates command arguments and indexes without panics
- prints returned events after successful actions
- uses a dev-playable setup:
  - reduced boss HP for shorter sessions
  - starting mana baseline so cards can be played from early turns
- clearly prints terminal outcome:
  - `GAME OVER: PLAYERS WON`
  - `GAME OVER: PLAYERS LOST`
