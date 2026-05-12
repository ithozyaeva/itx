# Eliminated hypotheses — debug 260512-1329

Гипотезы, которые проверены и **не подтверждены** или относятся к уже исправленным паттернам. Полезно для будущих сессий — не дублировать работу.

## E1. SSE `send on closed channel` panic в Publish — УЖЕ ИСПРАВЛЕНО
- Файл: `backend/internal/service/sse.go`
- Гипотеза: Race между Publish (RLock) и Unsubscribe (Lock).
- Результат: PR #341 (260505 debug) уже починил — Publish теперь делает send под RLock. Под `-race` стресс-тест `TestSSEHub_PublishWhileUnsubscribe` проходит.

## E2. `target="_blank"` в Mini App без openLink-обёртки — УЖЕ ИСПРАВЛЕНО
- Файлы: EventCard.vue:266,278, ContentGrid.vue:43, AIMaterialContentBlock.vue:54, MentorProfile.vue:104
- Гипотеза: остались места, где из Mini App юзер вываливается в браузер.
- Результат: все 4 точки имеют `@click.prevent="openLink(...)"` (PR #355). target="_blank" оставлен по делу — a11y, правый клик «открыть в новом окне» в обычном браузере. App.vue:62-83 ставит глобальный делегат для UGC-ссылок через v-html. Чисто.

## E3. Casino dice over off-by-one — УЖЕ ИСПРАВЛЕНО
- Файл: `backend/internal/service/casino.go`
- Результат: PR #341. roll >= target вместо strict >, winChance = 100-target. Симметрия с under восстановлена. EV=-3% на обе стороны.

## E4. Roles override `user.Roles = []models.Role{...}` — УЖЕ ИСПРАВЛЕНО
- Файл: `backend/internal/handler/auth_token.go`
- Результат: `mergeSubscriptionRole` ровно держит ADMIN/MENTOR/EVENT_MAKER. Слайсовый override-паттерн поискал — нашёл только в `repository/ai_material.go:105, 139` (`items[i].Tags = []string{}`), но это инициализация nil-tags для JSON-маршалинга, не баг.

## E5. Self-referral abuse (TrackConversion) — УЖЕ ИСПРАВЛЕНО
- Файл: `backend/internal/handler/referal_link.go`
- Результат: PR #342. authorId == member.id check перед grant.

## E6. Avatar S3 path с TG ID — УЖЕ ИСПРАВЛЕНО
- Файлы: `backend/internal/service/member.go`, `bot/telegram_bot.go`
- Результат: PR #342. UUID-only S3 keys + backfill миграция.

## E7. Username uniqueness / `/whois` takeover — УЖЕ ИСПРАВЛЕНО
- Файлы: миграции 20260506000000_dedupe_and_unique_username.sql, repository/member.go ErrUsernameTaken
- Результат: PR #342. Партиальный UNIQUE на LOWER(username), PATCH /me не пишет username, OAuth claim'ит.

## E8. JWT alg-confusion / отсутствие exp в Telegram Mini App auth
- Файл: `backend/internal/service/telegram_webapp.go`
- Гипотеза: `ValidateInitData` пропускает hash без проверки.
- Результат: реализация корректна — HMAC-SHA256(WebAppData→bot_token→data_check_string), `subtle.ConstantTimeCompare`, проверка `auth_date` против `WebAppInitDataMaxAge = 24h`. Бот-токен не утекает (только в Authorization header HMAC). Минор: `time.Since(authDate) > maxAge` для будущего authDate (clock skew) пропускает → допустим лёгкий look-ahead, но безопасно с т.зрения подмены.

## E9. Privilege escalation в admin-frontend — УЖЕ ИСПРАВЛЕНО
- Файл: `admin-frontend/src/router/index.ts`
- Результат: PR #342 добавил `ensureAdminAccess` через `/api/admin/me/permissions` перед каждым `requiresAuth` маршрутом.

## E10. nginx XSS / CSP headers
- Файл: `_nginx/conf.d/*.conf*`
- Проверено: нет Content-Security-Policy, нет X-Frame-Options, нет X-Content-Type-Options. Это design choice (не обязательно баг, особенно если фронты сами рендерят), но имеет смысл задизайнить в отдельной сессии. Не финдинг этой сессии.

## E11. landing-frontend addEventListener cleanup
- Файлы: HeroConstellation.vue, PromoteBackground.vue, useScrollProgress.ts, useMagneticHover.ts, PromoteSection.vue, UiPopover.vue
- Результат: все имеют onUnmounted с removeEventListener / clearTimeout / clearInterval. Чисто.

## E12. Vue type-check / lint
- Все три фронта: `vue-tsc --noEmit` + `eslint` — 0 errors, 0 warnings.

## E13. Go static analysis
- `go vet ./...` — clean.
- `go build ./...` — clean.

## E14. Subscription tier purchase atomicity
- Файл: `backend/internal/service/subscription.go:1110-1211`
- Гипотеза: Spend + SetManualTier в разных транзакциях → возможен баланс-без-тарифа.
- Результат: всё в одной `database.DB.Transaction(...)` — Spend (через SELECT FOR UPDATE) + EnsureUser + SetManualTier + AddAudit атомарно. Bessrochny grant защищён (ErrBessrochnyGrantExists). Downgrade защищён (ErrTierDowngrade). Чисто.

## E15. Raffle BuyTickets double-spend
- Файл: `backend/internal/service/raffle.go:50-115`
- Гипотеза: проверка balance отдельно от deduct.
- Результат: проверка balance перед транзакцией, но *деднут+INSERT тикетов внутри одной `database.DB.Transaction(...)`*. UNIQUE (raffle_id, member_id, source_type, source_id) для AwardTicketTx идемпотентен. Для BuyTicketsTx через `nextval('raffle_ticket_purchase_seq')` — каждый тикет уникален. Проверка `MaxBuyTicketsPerRequest=10_000` ограничивает single-batch. Чисто.
- Замечание: проверка balance ВНЕ транзакции (line 81) теоретически TOCTOU vs параллельная списанка через casino. На практике PlaceBet тоже в Transaction с проверкой `ErrInsufficientBalance` — двойного списания не будет, но юзер может увидеть «нет баллов» после успешного preflight. Минор.

## E16. Bot HMAC validation / shared secret
- Файл: `backend/internal/handler/auth_token.go:305` (HandleBotMessage)
- Результат: `subtle.ConstantTimeCompare` для X-Bot-Secret. Timing-safe.

## E17. mentor.go manual tx без defer commit/rollback в одной из веток
- Файл: `backend/internal/repository/mentor.go`
- Гипотеза: ручное управление tx часто упускает rollback на error path.
- Результат: проверено быстрым взглядом — везде `defer tx.Rollback()` либо явный `tx.Rollback()` перед return. `tx.Commit()` в конце успешных путей. Похоже корректно, но проверить полностью требует deep dive — не в этой сессии.

## E18. CreateOrUpdateToken token-reuse при concurrent login
- Файл: `backend/internal/service/auth_token.go:75-84`
- Гипотеза: race между двумя concurrent логинами одного юзера → один из них упадёт на UNIQUE, второй увидит ghost-token.
- Результат: основной баг тут — silent error swallow (см. M3 в findings.md). Race добавляет к нему второе измерение, но root cause тот же.
