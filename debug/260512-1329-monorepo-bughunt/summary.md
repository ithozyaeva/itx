# Debug 260512-1329 — Executive Summary

**Scope:** весь монорепо IT-X (backend Go + 3 фронта Vue + nginx/CI).
**Mode:** report only (без авто-фиксов).
**Iterations:** ~26 шагов (13 confirmed hypotheses + 13 disproven/уже исправленных).

## Score

```
Bugs found:        9 (0 Critical, 2 High, 3 Medium, 4 Low)
Hypotheses tested: 26 (9 confirmed, 13 disproven, 4 already fixed in earlier PRs)
Static analysis:   clean baseline (go vet/build, vue-tsc, eslint — 0 signals)
```

## Findings overview

| # | Severity | Location | Title |
|---|----------|----------|-------|
| H1 | HIGH | `_nginx/conf.d/*.conf*:client_max_body_size` | nginx 10MB режет marketplace uploads до backend BodyLimit 12MB |
| H2 | HIGH | `handler/auth_token.go:184-213` | Authenticate() не проверяет expiry — Logout не инвалидирует session |
| M1 | MEDIUM | `cmd/main.go:134,139` | `UTC().Weekday()` для еженедельного award — съезжает с MSK |
| M2 | MEDIUM | `repository/{feedback,kudos,chat_quest,chat_activity}.go` | `CURRENT_DATE` (UTC) для дневных лимитов — 3ч смещение в MSK |
| M3 | MEDIUM | `service/auth_token.go:75-84` | `CreateOrUpdateToken` молча глотает Create/Update ошибки |
| L1 | LOW | `handler/auth_token.go:192-204` | Authenticate goroutine race на `user.Roles` + нет recover |
| L2 | LOW | `admin-frontend/src/components/modals/*Modal.vue` | Дебаунс-таймауты не очищаются на unmount |
| L3 | LOW | `bot/subscription.go:1178,1315,1350` | Bot admin-command goroutines без `defer recover()` |
| L4 | LOW | `handler/*.go` (15+ сайтов) | Fire-and-forget goroutines без `defer recover()` — крах процесса при panic |

Подробно — `findings.md`. Disproven гипотезы — `eliminated.md`.

## Top 3 рекомендации (если фиксить только что-то)

1. **H1 — поднять nginx body limit до 12m в обоих конфигах.** Один LOC × 2 файла, мгновенный fix известной регрессии PR #342.
2. **H2 — добавить `CheckExpirationDate` в `Authenticate`.** Симметрично с `RefreshToken` и middleware; ~3 LOC. Закрывает logout-loop, имеет security-семантику.
3. **L4 — ввести `service.SafeGo(fn, label)` хелпер.** Системно укрепляет 15+ точек запуска goroutine от panic-induced restarts. Параллельно решает L1/L3.

## Паттерны, которые повторяются между debug-сессиями

Просмотрев PR #341 (260505), #342 (260506), #355 (260512 предыдущий), я вижу системные классы:
- **Таймзонные баги** (UTC vs MSK) — третий раз вылезает. Хорошая инвестиция: ввести `db_now_msk()` SQL-функцию ИЛИ единый Go-хелпер и пройтись по всем `CURRENT_DATE`/`time.Now().UTC()`.
- **Расширение `[]Role`-подобных слайсов через прямой override** — PR #341 чинил Roles, я искал такие же в Tags/Permissions — нашёл только non-bug инициализацию пустого слайса.
- **Goroutine + state-write без recover/lock** — H1 (Roles race), L1, L3, L4. Систематический.
- **«один слой защиты» в auth-flow** — H2: middleware ловит expiry, но `Authenticate` не ловит. Аналогично PR #342 правил отсутствие admin-check на роутере, хотя API уже 403-ил.

## Что я **не** сканировал (out of scope этой сессии)

- Migrations correctness (есть много `*.sql` файлов; не профильный домен этой сессии).
- Mentor.go manual transactions (большой файл с 10+ Commit/Rollback — отдельная сессия).
- Frontend a11y / WCAG.
- Performance (N+1, missing index) — это `/autoresearch:debug --technique=trace`-сессия.
- nginx CSP/security headers — design discussion, не баг.

## Note on tooling

Три параллельных Explore-агента запущены в фоне в Phase 2 для broad recon (Go backend / Vue platform / security). Все три упали с `API Error: ConnectionRefused` после ~12 минут. Findings ниже получены ручным поиском паттернов от main-агента, опираясь на git log последних 30 коммитов + структуру кода + знание трёх предыдущих debug-сессий, которые пользователь уже закрыл (260505/260506/260512-предыдущая).
