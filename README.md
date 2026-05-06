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

## Roadmap
### Этап 1. Чистый internal/game
-[x] Шаг 1. Создать Go-проект и базовую структуру
-[x] Шаг 2. Описать основные типы
-[x] Шаг 3. Сделать NewGame
-[x] Шаг 4. Добавить реестр карт
-[x] Шаг 5. Добавить добор карт
-[x] Шаг 6. Добавить ману
-[x] Шаг 7. Добавить ApplyAction (ActionEndTurn)
-[x] Шаг 8. Добавить spell-карту “урон боссу”
-[x] Шаг 9. Добавить проверку победы
-[x] Шаг 10. Добавить spell-карту лечения
-[x] Шаг 11. Добавить ValidTargets
-[ ] Шаг 12. Добавить minion-карты
-[ ] Шаг 13. Добавить атаку minion’ом по боссу
-[ ] Шаг 14. Обновлять minions в начале хода
-[ ] Шаг 15. Добавить способности босса
-[ ] Шаг 16. Добавить перемещение босса
-[ ] Шаг 17. Добавить поражение
-[ ] Шаг 18. Добавить смерть minion’ов
-[ ] Шаг 19. Добавить полный сценарный тест
-[ ] Шаг 20. Добавить сценарий поражения

### Этап 2. Single-player CLI/dev runner
- Шаг 21. Сделать простой console runner
- Шаг 22. Добавить команды в console runner
- Шаг 23. Доиграть партию в консоли

### Этап 3. Web single-player на Go + HTMX
- Шаг 24. Поднять HTTP server
- Шаг 25. Отрендерить состояние игры
- Шаг 26. Добавить HTMX End Turn
- Шаг 27. Добавить HTMX play card без цели
- Шаг 28. Добавить play card с целью
- Шаг 29. Добавить атаку minion’ом
- Шаг 30. Добавить game over UI

### Этап 4. Полировка single-player MVP
- Шаг 31. Добавить restart
- Шаг 32. Добавить нормальный combat log
- Шаг 33. Добавить минимальный CSS
