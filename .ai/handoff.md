# HANDOFF

## Project overview

This project is a single-player prototype of a cooperative card game inspired by Hearthstone Tavern Brawl "Unite Against Mechazod!".

The current goal is to build a clean Go game engine first, then later add HTMX/web UI on top.

Current stack:

- Go
- Pure backend/game engine
- Single-player prototype where one user controls two players
- CLI/dev runner for manual testing
- No database
- No HTTP/web dependency inside `internal/game`
- No HTMX/templates/UI dependency inside `internal/game`

Core package:

```txt
internal/game
```

The package `internal/game` must stay pure domain/game-engine code.

It must not import:

- `net/http`
- `html/template`
- `text/template`
- HTMX/UI packages
- database packages

All gameplay actions are expected to go through:

```go
func ApplyAction(g *Game, action Action) ([]GameEvent, error)
```

Current validation command:

```bash
go test ./...
```

Current server/dev commands:

```bash
go run ./cmd/server
go run ./cmd/devgame
```

---

# Completed work

## Stage 1 - Pure game engine

Stage 1 is complete.

The Stage 1 roadmap in the attached plan covered steps 1-20, from initial Go project setup through full win/loss scenario tests. The plan includes creating the pure `internal/game` package, adding core types, `NewGame`, card registry, draw, mana, `ApplyAction`, playable cards, targeting, minions, attacks, boss behavior, win/loss conditions, minion cleanup, and scenario tests. This handoff reflects that Stage 1 has been completed. 

All Stage 1 steps are considered done, and:

```bash
go test ./...
```

passes.

## Stage 2 - Single-player CLI/dev runner

Stage 2 is complete.

The Stage 2 roadmap covered:

- Step 21: create a simple console runner
- Step 22: add console commands
- Step 23: make a full game playable from console until `won` or `lost`

All Stage 2 steps are considered done, and:

```bash
go test ./...
```

passes.

---

# Current project structure

Expected high-level structure:

```txt
cmd/
  server/
    main.go
  devgame/
    main.go

internal/
  game/
    action.go
    boss.go
    card.go
    card_registry.go
    card_validation.go
    combat.go
    deck.go
    engine.go
    event.go
    game.go
    gameover.go
    healing.go
    ids.go
    mana.go
    minion.go
    player.go
    player_config.go
    shuffle.go
    summon.go
    targeting.go
    ...
```

There are also test files under `internal/game`.

Exact file names may differ slightly if some logic was grouped differently, but the important architectural boundary is:

```txt
cmd/*       -> app entry points / runners
internal/game -> pure game engine
```

---

# Core domain model

## Game

The root aggregate is `Game`.

It contains concepts such as:

- game ID
- game status
- two players
- boss
- turn number
- active player ID
- event history
- seed/random state if implemented for deterministic behavior

Game statuses use the `GameStatus...` naming style:

```go
GameStatusCreated
GameStatusActive
GameStatusWon
GameStatusLost
```

Important invariant:

```txt
If Game.Status is not GameStatusActive, ApplyAction should reject new actions.
```

## Player

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

Starting player health:

```go
StartingPlayerHealth = 30
```

Mana starts at `0/0`.

Mana rules:

- `RefreshMana` increases `MaxMana` by 1.
- `MaxMana` is capped at 10.
- after refresh, `Mana == MaxMana`.
- `SpendMana` rejects negative amounts.
- `SpendMana` rejects spending more mana than available.
- error paths do not mutate player mana.

Implemented:

```go
func RefreshMana(g *Game, playerIndex int)
func SpendMana(g *Game, playerIndex int, amount int) error
```

## Boss

Boss has concepts such as:

- `ID`
- `Name`
- `Health`
- `MaxHealth`
- `Attack`
- `Armor`
- side/position if implemented in the boss movement step

Boss HP is clamped to 0 when damaged below zero.

Boss behavior has been added in later Stage 1 steps:

- boss abilities
- boss movement
- defeat condition through player death

---

# Cards

Cards are defined in `CardRegistry`.

`CardRegistry` is the source of truth for all card definitions.

Implemented base cards:

## Strike

```txt
ID: strike
Cost: 1
Type: spell
Effect: deal 3 damage to boss
Targeting: does not require manual target selection
```

## Repair

```txt
ID: repair
Cost: 2
Type: spell
Effect: restore 5 health to chosen hero
Targeting: required
Valid targets: hero:0, hero:1
Invalid target: boss
```

## Drone

```txt
ID: drone
Cost: 2
Type: minion
Stats: 2/3
Effect: summon
```

Later Stage 1 steps may also include extra cards or test fixtures to make full win/loss scenarios possible.

## CardDefinition

The card definition model supports:

- card ID
- name
- type
- cost
- description
- effect
- targeting rules
- minion attack
- minion health

Expected fields include concepts like:

```go
ID
Name
Type
Cost
Description
Effect
Targeting
Attack
Health
```

## CardInstance

A card instance represents a concrete copy of a card in a specific game.

It has:

- instance ID
- definition ID
- owner ID

`Action.CardID` refers to a `CardInstanceID`, not the registry `CardID`.

---

# Minions

Minions are board entities created from minion card definitions.

Minions have concepts such as:

- `ID`
- `DefinitionID`
- `OwnerID`
- `Name`
- `Attack`
- `Health`
- `MaxHealth`
- `CanAttack`
- possibly `Exhausted`, if still kept

Rules implemented:

- playing `drone` creates a minion on the active player's board
- summoned minions start with `CanAttack = false`
- max board size is 7
- full board prevents minion play
- minions can attack the boss through `ActionAttack`
- minions cannot attack twice in a turn
- minions refresh at the start of their owner's turn
- dead minions can be cleaned up

---

# Game setup

Implemented:

```go
func NewGame(id string, p1 PlayerConfig, p2 PlayerConfig, seed int64) *Game
```

`NewGame` creates:

- game ID
- two players
- boss
- starting decks
- starting hands
- active player
- active game status
- deterministic shuffle by seed

Acceptance conditions already satisfied:

- game has ID
- there are 2 players
- each player starts with 30 HP
- boss has positive HP
- each player has deck and hand
- game starts active
- active player is set
- same seed produces predictable setup
- `go test ./...` passes

---

# Decks and draw

Starting decks use card IDs that must exist in `CardRegistry`.

`StartingDeckSize` is still kept.

Implemented:

```go
func NewStartingDeck(ownerID PlayerID) []CardInstance
func DrawCard(g *Game, playerIndex int) []GameEvent
```

Draw rules:

- draws the top card from `Deck[0]`
- appends the card to `Hand`
- deck size decreases by 1
- hand size increases by 1
- creates `EventTypeCardDrawn`
- does not panic for empty deck, nil game, or invalid player index

---

# Actions

All actions go through:

```go
func ApplyAction(g *Game, action Action) ([]GameEvent, error)
```

Implemented action types include:

```go
ActionTypeEndTurn
ActionTypePlayCard
ActionTypeAttack
```

Aliases may exist:

```go
ActionEndTurn
ActionPlayCard
ActionAttack
```

`Action` includes concepts like:

- action type
- player ID
- card instance ID
- source minion ID or source index
- target ID
- target object

Target IDs include:

```txt
hero:0
hero:1
boss
```

---

# End turn

`ActionTypeEndTurn` is implemented.

Rules:

- only works when game status is active
- only active player can end turn
- switches active player
- increments turn
- refreshes mana for new active player
- draws a card for new active player
- refreshes minions for new active player
- moves boss if boss movement is enabled
- applies boss ability if boss abilities are enabled
- creates turn/boss-related events
- rejects actions after game is won or lost

---

# Play card

`ActionTypePlayCard` is implemented.

Rules:

- only active player can play cards
- card must be in active player's hand
- card definition must exist in `CardRegistry`
- player must have enough mana
- valid target is required when card targeting says so
- on success:
  - mana is spent
  - card is removed from hand
  - effect is applied
  - events are appended to `g.Events`
  - events are returned from `ApplyAction`
- on error:
  - state must not partially mutate

Implemented play behavior:

## Strike

- costs 1 mana
- removes the card from hand
- deals 3 damage to boss
- creates card played event
- creates damage event
- checks game over after damage

## Repair

- costs 2 mana
- requires target
- valid targets are `hero:0` and `hero:1`
- invalid targets include `boss`, unknown target, and empty target
- heals selected hero for 5
- clamps hero HP to `MaxHealth`
- creates card played event
- creates heal event

## Drone

- costs 2 mana
- creates 2/3 minion on active player's board
- summoned minion starts with `CanAttack = false`
- board limit is 7
- full board rejects the play without spending mana or removing card
- creates card played event
- creates minion summoned event

---

# Targeting

Implemented:

```go
func ValidTargets(g *Game, playerID string, cardInstanceID string) ([]Target, error)
```

Rules:

- finds player by `playerID`
- finds card instance in that player's hand
- finds card definition in `CardRegistry`
- if card does not require a target, returns an empty list
- for `repair`, returns `hero:0` and `hero:1`
- for `repair`, does not return `boss`
- unknown player returns error
- card not in hand returns error
- unknown card definition returns error

Target contains UI-friendly metadata such as:

- `ID`
- `Type`
- `Kind`
- `PlayerID`
- `BossID`
- `MinionID`
- `OwnerID`
- `DisplayName`

Targeting helpers include:

```go
ResolveTarget
ValidateTargetForCard
```

Important invariant:

```txt
ValidTargets and ActionPlayCard targeting validation must use the same CardDefinition.Targeting rules.
```

---

# Combat

Implemented combat concepts include:

- damage to boss
- healing heroes
- minion attack against boss
- minion attack availability
- minion refresh at start of owner's turn
- dead minion cleanup

Damage rules:

- boss HP is clamped to 0
- damage creates a damage event
- after boss damage, game over is checked

Healing rules:

- healing targets heroes
- healing cannot exceed `MaxHealth`
- repair cannot heal the boss
- healing creates a heal event

Minion attack rules:

- only active player's minion can attack
- minion must have `CanAttack == true`
- target currently supported: boss
- after attack, minion's `CanAttack` becomes false
- boss receives damage equal to minion attack
- game over is checked after boss damage

---

# Game over

Implemented:

```go
func CheckGameOver(g *Game) []GameEvent
```

Victory rule:

```txt
if Boss.Health <= 0 -> GameStatusWon
```

Defeat rule:

```txt
if any player's Health <= 0 -> GameStatusLost
```

Game over rules:

- status becomes `GameStatusWon` or `GameStatusLost`
- creates corresponding event
- appends event to `g.Events`
- repeated checks do not duplicate game over events
- after terminal status, `ApplyAction` rejects new actions

---

# Boss behavior

Stage 1 included boss abilities and boss movement.

Implemented boss ability examples from the roadmap:

```txt
Zap Heroes: deal 2 damage to both heroes
Bomb Salvo: deal 2 damage to a random target
Overclock: boss gets +1 Attack
```

Rules:

- boss applies one ability at the start of turns, if implemented as planned
- creates a `boss_ability` event
- ability mutates game state
- random behavior is controlled by seed
- seeded tests should be deterministic

Boss movement:

- boss side/position changes with active player, or boss is considered to stand on active player's side
- movement creates a `boss_moved` event
- expected side behavior:
  - active Player 1 -> boss side 0
  - active Player 2 -> boss side 1

---

# Events

The event system is used to describe state changes.

Known event concepts include:

```go
EventTypeCardDrawn
EventTypeTurnStarted
EventTypeCardPlayed
EventTypeDamageDealt
EventTypeHeal
EventTypeMinionSummoned
EventTypeBossAbility
EventTypeBossMoved
EventTypeMinionDied
EventTypeGameWon
EventTypeGameLost
```

Aliases may exist, such as:

```go
EventCardPlayed
EventDamage
EventHeal
EventMinionSummoned
EventGameWon
EventGameLost
```

Events are expected to be:

- returned from `ApplyAction`
- appended to `g.Events`
- useful later for UI, logs, and replay/debugging

---

# Errors

Domain errors include concepts like:

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
ErrUnsupportedCardEffect
ErrTargetRequired
ErrInvalidTarget
ErrBoardFull
ErrMinionNotFound
ErrMinionCannotAttack
ErrInvalidAttackTarget
```

Exact names should be confirmed against code before adding new features.

Important invariant:

```txt
Error paths must not partially mutate game state.
```

---

# CLI/dev runner

Stage 2 added a single-player CLI/dev runner.

Entry point:

```txt
cmd/devgame/main.go
```

Run command:

```bash
go run ./cmd/devgame
```

The CLI/dev runner prints:

- HP of both players
- HP of boss
- current turn / active player
- active player's hand
- board/minion state
- available actions or command hints
- game status

Implemented commands from Stage 2:

```txt
hand
play <handIndex> [targetID]
attack <minionIndex> boss
end
state
quit
```

Expected behavior:

- `hand` prints active player's hand
- `play <handIndex> [targetID]` plays a card from active player's hand
- `attack <minionIndex> boss` attacks the boss with active player's minion
- `end` ends current player's turn
- `state` prints current game state
- `quit` exits runner

The runner should:

- call `ApplyAction`
- print events returned by actions
- print errors without panicking
- allow a full game to be played to `won` or `lost`
- keep state consistent after every command

---

# Scenario tests

Stage 1 included scenario-level tests.

## Win scenario

Expected test concept:

```go
func TestGameCanBePlayedUntilWin(t *testing.T)
```

This test should:

- create game
- advance turns
- play cards
- use minion attacks
- let boss abilities happen
- eventually reach `GameStatusWon`

## Loss scenario

Expected test concept:

```go
func TestGameCanBeLost(t *testing.T)
```

This test should:

- create game
- let boss deal enough damage
- reduce at least one player to 0 HP
- set status to `GameStatusLost`
- reject further actions

---

# Important invariants

- `internal/game` must stay UI-agnostic.
- All gameplay actions must go through `ApplyAction`.
- `go test ./...` must pass after every step.
- `CardRegistry` is the source of truth for card definitions.
- Starting decks must not reference cards missing from `CardRegistry`.
- `ValidTargets` and `ActionPlayCard` targeting validation should use the same targeting rules.
- Boss HP should not display below 0.
- Player HP should not exceed `MaxHealth`.
- A full board prevents minion play.
- Summoned minions should not attack immediately.
- Minions that attacked should refresh only on their owner's turn.
- After `GameStatusWon` or `GameStatusLost`, no new actions should be accepted by `ApplyAction`.
- CLI runner must never panic on invalid user input.
- Error paths should not mutate mana, hand, board, HP, events, turn, or active player.

---

# Useful commands

Run all tests:

```bash
go test ./...
```

Run tests verbosely:

```bash
go test -v ./...
```

Run only game package tests:

```bash
go test -v ./internal/game
```

Run server stub:

```bash
go run ./cmd/server
```

Run CLI/dev runner:

```bash
go run ./cmd/devgame
```

Check that `internal/game` has no forbidden dependencies:

```bash
grep -R "net/http\\|html/template\\|text/template\\|htmx\\|database/sql" internal/game || true
```

---

# Suggested next directions

Since Stage 1 and Stage 2 are complete, possible next directions are:

## Option A - Refactor and stabilize before web

- split large files if needed
- centralize errors
- centralize fixtures/test helpers
- improve event messages
- document public functions
- add more scenario tests
- add command examples to README

## Option B - Start web/HTMX layer

Potential next stage:

```txt
Stage 3 - Minimal web shell
```

Possible first steps:

1. create `cmd/web/main.go`
2. add HTTP server outside `internal/game`
3. render game state read-only
4. add HTMX endpoint for `end`
5. add HTMX endpoint for `play`
6. add HTMX endpoint for `attack`
7. keep all game mutations inside `ApplyAction`

Important: the web layer should depend on `internal/game`, but `internal/game` must not depend on web.

## Option C - Improve CLI/dev runner

- better command parser
- command help
- target listing for playable cards
- hand indexes with card details
- board indexes with attack status
- event log printing
- seed selection from command line

---

# Known uncertainty to confirm against actual code

Before making the next code change, confirm these facts in the repository:

1. Exact event constant names and aliases.
2. Exact `Target` struct fields.
3. Exact `CardDefinition` fields.
4. Exact error names.
5. Exact boss ability implementation details.
6. Exact boss side/position field name.
7. Exact minion attack action shape.
8. Exact CLI command parsing structure.
9. Whether `Discard` is already used when playing cards.
10. Whether scenario tests use real cards or special test fixtures.

This handoff describes the completed architecture and behavior at the end of Stage 1 and Stage 2.
