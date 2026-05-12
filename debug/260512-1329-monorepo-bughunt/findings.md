# Debug session 260512-1329 — Monorepo bug hunt

Scope: весь монорепо (backend + 3 фронтенда + infra/CI). Депт: unlimited.
Mode: report only (без авто-фиксов).

Static-analysis baseline чистый: `go vet ./...`, `go build ./...`, `vue-tsc --noEmit` × 3 фронта, `eslint` × 3 фронта — нулевые сигналы. Все findings ниже найдены ручным поиском паттернов + чтением кода.

---

## [HIGH] H1. Nginx `client_max_body_size 10m` блокирует marketplace-загрузки до того, как они дойдут до backend BodyLimit 12MB

**Location:**
- `_nginx/conf.d/nginx.conf.template:52` — `client_max_body_size 10m;`
- `_nginx/conf.d/localhost.conf:12` — `client_max_body_size 10m;`
- `backend/cmd/main.go:54-56` — `BodyLimit: 12 * 1024 * 1024`
- `backend/internal/handler/marketplace.go` — `fileHeader.Size > 10MB → 400`

**Hypothesis:** PR #342 (260506 security debug) поднял Fiber BodyLimit до 12MB специально под marketplace, чтобы 10MB-картинка + поля формы прошли. Но nginx сидит перед Fiber и режет тело по своему `client_max_body_size 10m` — то есть всё то же, что было.

**Evidence:** nginx-первый, Fiber-второй. Запрос с 9.8MB JPEG + multipart-полями (Content-Type, boundary, имя файла, описание, цена) лезет за 10MB суммарно — nginx 413 Request Entity Too Large прилетает раньше, чем Fiber вообще видит запрос. Backend-валидатор `fileHeader.Size > 10MB` никогда не получает шанс отработать на boundary-кейсах.

**Reproduction:**
1. На фронте marketplace кладёшь файл ~9.9MB (порог frontend).
2. POST /api/marketplace/items приходит на nginx с Content-Length ~10.1MB (файл + поля).
3. Nginx отвечает 413 — без `proxy_pass`.
4. Юзер видит generic upload error.

**Impact:** MEDIUM-HIGH UX-регрессия для marketplace и любой другой формы с файлом близким к 10MB.  Frontend-лимит и backend-лимит «защищают», но nginx разрывает прежде их обоих, причём с непрозрачным error response. Скорее всего, та же проблема и в avatar/resume загрузках, если они когда-нибудь подойдут к 10MB.

**Suggested fix:** в обоих nginx-конфигах поднять `client_max_body_size 12m` (или сколько решено).  Альтернатива (хуже): срезать backend BodyLimit до 10MB и понизить frontend-лимит, чтобы все слои согласованно отказывали ровно на одном пороге.

**Severity:** HIGH (функциональная регрессия, тихая ошибка, скрытая от юзера).

---

## [HIGH] H2. `Authenticate()` не проверяет expiry — Logout не инвалидирует session-token при повторной авторизации

**Location:** `backend/internal/handler/auth_token.go:175-214`

**Hypothesis:** Хендлер `POST /api/auth/telegram` берёт токен из тела и просто возвращает 200 + тот же токен, если строка существует в `auth_tokens`. В отличие от `RefreshToken` (line 260) и middleware (lines 47, 82) — нет вызова `utils.CheckExpirationDate(authToken.ExpiredAt)`.

**Evidence:**
```go
existingToken, existingUser, err := h.authService.GetByToken(req.Token)
if err != nil {
    return c.Status(fiber.StatusUnauthorized).JSON(...)
}
// ← здесь нет CheckExpirationDate(existingToken.ExpiredAt)
go func(user *models.Member) { ... CheckUserInChat ... }(existingUser)
c.Response().Header.Add("X-Telegram-User-Token", existingToken.Token)
return c.JSON(fiber.Map{"user": existingUser, "token": existingToken.Token})
```

И `repository/auth_token.go:76` (GetByToken) — `Where("token = ?", token).First(&user)` без фильтра по `expired_at`. То есть Logout, который двигает `expired_at` в прошлое, помечает строку «недействительной» только для middleware и Refresh.

**Reproduction:**
1. Юзер открывает платформу через бот-ссылку `?token=T`. Это вызывает `POST /api/auth/telegram {token: T}` → 200.
2. Юзер кликает «Выйти» (`POST /api/auth/telegram/logout`). `InvalidateToken(T)` ставит `expired_at = now-1h`.
3. Юзер снова кликает старую бот-ссылку с `?token=T` (телега DM хранит её бессрочно).
4. `POST /api/auth/telegram {token: T}` → находит строку, **возвращает 200 с тем же expired T в `X-Telegram-User-Token`**.
5. Клиент кладёт T в localStorage и зовёт `GET /api/profile/me` → middleware `CheckExpirationDate` → 401.
6. UI скатывается в «сессия истекла» (#332), хотя сервер только что сказал «логин успешен».

**Impact:**
- Логично-семантическая дыра: серверный logout (PR #341 — закрывал «sniff-токен из nginx access-log SSE-стрима») должен инвалидировать сессию. Сейчас инвалидация **частичная**: middleware ловит, но `/api/auth/telegram` (точка входа) — нет.
- Side-effect через goroutine (`bot.CheckUserInChat` + `memberService.Update`) выполняется на «вышедшей» сессии: DB-запись через invalid token. Малый impact, но семантическая утечка.
- Не «полный» privilege escalation, потому что middleware всё равно бьёт по голове — но это анти-паттерн «один слой защиты» и UX-яма: после logout юзер с одним кликом снова получает X-Telegram-User-Token успешно и видит платформу мерцающую → лендинг.

**Suggested fix:**
```go
existingToken, existingUser, err := h.authService.GetByToken(req.Token)
if err != nil { return 401 }
if utils.CheckExpirationDate(existingToken.ExpiredAt) {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
}
```

Симметрично с уже исправленным RefreshToken (auth_token.go:260).

**Severity:** HIGH (логика logout, security-adjacent + UX-loop).

---

## [MEDIUM] M1. `time.Now().UTC().Weekday() == time.Monday` — еженедельный award «Чаттер недели» по UTC, ticker уплывает в воскресенье/вторник для MSK-юзера

**Location:** `backend/cmd/main.go:134, 139`

**Hypothesis:** 24-часовой ticker запускается в момент старта процесса. На каждом тике (= +24h от старта) проверяется `time.Now().UTC().Weekday() == time.Monday` и выдаётся «Чаттер недели». Между MSK и UTC 3 часа дельты, поэтому если контейнер стартанул в окне 21:00–23:59 MSK (= 18:00–20:59 UTC), tick фиксируется в это же время суток, и в день, когда **MSK-понедельник, но UTC-воскресенье** (с 21:00 до 23:59 MSK воскресенья) — `Weekday() == Sunday`, award не выдан. И наоборот: вторник MSK в 00:00–02:59 MSK = понедельник 21:00–23:59 UTC, и award выдан на «не тот» календарный день в восприятии юзера.

**Evidence:** там же — рядом, line 196, `nowMSK := time.Now().In(utils.MSKLocation())` для morning/evening push явно использует MSK. Значит, разработчик уже знает про MSK; weekly-award просто забыли.

**Reproduction:** старт контейнера в Wed 22:00 MSK → ticker fires Thu/Fri/Sat/Sun/Mon at 22:00 MSK. На Mon 22:00 MSK UTC weekday = Monday (19:00 UTC) — фиксируется корректно. Но: старт в Tue 01:00 MSK → ticker fires Wed/Thu/Fri/Sat/Sun/Mon/Tue/... at 01:00 MSK. Mon 01:00 MSK = Sun 22:00 UTC → Weekday = Sunday → award пропущен; Tue 01:00 MSK = Mon 22:00 UTC → Weekday = Monday → award выдан **во вторник** в MSK-восприятии.

**Impact:** «Чаттер недели» иногда выдаётся не в понедельник, иногда не выдаётся за неделю вовсе (если идемпотентность keyed по неделе и в одну ISO-неделю случай попал в Sunday-окно). Для рестарта-сразу-после-полуночи MSK — стабильно мимо.

**Suggested fix:** `time.Now().In(utils.MSKLocation()).Weekday() == time.Monday` в обеих точках. Симметрично birthday-checker (`backend/internal/bot/telegram_bot.go startBirthdayChecker`), который уже переписан на MSK в PR #341.

**Severity:** MEDIUM (геймификация-корректность, low business impact, повторяемо).

---

## [MEDIUM] M2. `CURRENT_DATE` (UTC session TZ) для «сегодня»-лимитов — feedback / kudos / chat_quest / chat_activity

**Locations:**
- `backend/internal/repository/feedback.go:21` — `created_at >= CURRENT_DATE` (дневной лимит фидбэков на юзера)
- `backend/internal/repository/kudos.go:44` — `created_at >= CURRENT_DATE` (дневной лимит kudos)
- `backend/internal/repository/chat_quest.go:195` — `day = CURRENT_DATE - (rn-1) * INTERVAL '1 day'` (расчёт текущего streak)
- `backend/internal/repository/chat_activity.go:97,98,104,114,115,121,190,206,207,213,234,246` — несколько мест, окна активности «N days back from today»

**Hypothesis:** PR #341 (260505 fix) уже исправил birthday — `(CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow')::date`. Здесь тот же класс багов остался: DSN ставит session TZ='UTC', `CURRENT_DATE` отдаёт UTC-дату, но «сегодня» в восприятии юзера — MSK-дата. Дельта 3 часа.

**Evidence (kudos example):**
```go
Where("from_id = ? AND created_at >= CURRENT_DATE", fromId).Count(&count)
```
Юзер шлёт 1-ый kudos в 23:50 MSK (20:50 UTC May 11) и 2-ой в 00:30 MSK May 12 (= 21:30 UTC May 11) — для CURRENT_DATE оба попадают в **«May 11 UTC»**, дневной счётчик считает 2 за этот «UTC-день». Юзер: «я только что после полуночи MSK, должен был обнулиться». Симметрично в 03:00 MSK обнуляется CURRENT_DATE, юзер видит обнулённый счётчик «случайно».

**Impact:** дневные ограничения в восприятии юзера съезжают на 3 часа. Для feedback/kudos — слабо заметно. Для **chat_quest streak** хуже: streak ломается, если активность была в окне 00:00–03:00 MSK (UTC-вчера) или, наоборот, фантомный «лишний» день streak'а.

**Suggested fix:** заменить все `CURRENT_DATE` на `(CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow')::date`. Лучше — единый SQL-хелпер `msk_today()` (либо migration `CREATE FUNCTION`, либо Go-константа со строкой подстановки).

**Severity:** MEDIUM (correctness, особенно streak; user-visible).

---

## [MEDIUM] M3. `CreateOrUpdateToken` молча глотает ошибки Create/Update — silent token-write failure

**Location:** `backend/internal/service/auth_token.go:75-84`

**Hypothesis:** Возвращает `(authToken, nil)` всегда, при этом ошибки `s.authRepo.Create(...)` и `s.authRepo.Update(...)` ничем не обрабатываются и не пробрасываются. `authToken` в return — результат `GetByTelegramID`, **до** Create/Update. Если Update упал (deadlock, connection drop, constraint violation) — серверу всё равно: возвращаем «успех».

**Evidence:**
```go
func (s *AuthTokenService) CreateOrUpdateToken(telegramID int64, token string) (*models.AuthToken, error) {
    authToken, err := s.authRepo.GetByTelegramID(telegramID)
    if err != nil {
        s.authRepo.Create(&models.AuthToken{TelegramID: telegramID, Token: token, ExpiredAt: ...})
    } else {
        s.authRepo.Update(&models.AuthToken{ID: authToken.ID, TelegramID: telegramID, Token: token, ExpiredAt: ...})
    }
    return authToken, nil  // ← всегда nil, даже если Create/Update упал
}
```

Зовётся из: `AuthenticateWebApp` (line 153), `HandleBotMessage` (line 378), `RefreshToken` (line 266). Все три точки в случае tomato ошибки Create/Update НЕ скажут юзеру 500 — отдадут 200 со свежим token-стрингом, который в БД не сохранён (либо сохранён предыдущая запись со старым token).

**Impact:** при сетевом/DB-сбое юзер получает токен, который НЕ работает. Не security: middleware просто 401 на следующий запрос. Но это превращает временный DB-flap в постоянное «не залогиниться» (пока не сделаешь refresh, который попробует обновить тот же ghost-token).

**Suggested fix:**
```go
if err != nil {
    created, cerr := s.authRepo.Create(&models.AuthToken{...})
    if cerr != nil { return nil, cerr }
    return created, nil
} else {
    updated, uerr := s.authRepo.Update(&models.AuthToken{...})
    if uerr != nil { return nil, uerr }
    return updated, nil
}
```

И хендлеры обязаны проверять err и возвращать 500.

**Severity:** MEDIUM (надёжность, silent failure под нагрузкой; не security).

---

## [LOW] L1. Authenticate goroutine race — `user.Roles = newRoles` параллельно с JSON-сериализацией

**Location:** `backend/internal/handler/auth_token.go:192-204`

**Hypothesis:** Хендлер запускает goroutine, которая **мутирует тот же `*models.Member`**, который параллельно сериализуется в JSON-ответ (`return c.JSON(fiber.Map{"user": existingUser})`). `user.Roles = newRoles` (line 202) — write; маршалинг ответа — read. Гонка по полю.

**Evidence:** идентичный pattern уже всплыл в PR #341 (mergeSubscriptionRole vынесли в чистую функцию, но goroutine-вызов остался на том же указателе). `existingUser` передан в goroutine по pointer-захвату (хоть и `func(user *models.Member)`), потом возвращён в response без копии. `go test -race` нашёл бы.

**Impact:**
- При маршалинге ответа Vue может получить inconsistent role-list: либо старый, либо новый, либо комбинацию (если slice realloc'нулся в середине read).
- Кроме того, **goroutine не имеет `defer recover()`**: если `bot.CheckUserInChat` или `memberService.Update` запаникуют — крашится весь backend-процесс.

**Suggested fix:**
1. Передать в goroutine копию `*user` либо только нужные поля (`id`, `roles`).
2. Добавить `defer func() { if r := recover(); r != nil { log.Printf("auth role-sync panic: %v", r) } }()`.

**Severity:** LOW (узкое окно гонки, последствия — нестабильность под -race; panic-recover важнее).

---

## [LOW] L2. Дебаунс-таймауты в admin-modal'ях не очищаются на unmount — stale fetch

**Locations:**
- `admin-frontend/src/components/modals/PointsAwardModal.vue:27`
- `admin-frontend/src/components/modals/CreditsAwardModal.vue:27`
- `admin-frontend/src/components/modals/SubscriptionChatModal.vue:29`

**Hypothesis:** Во всех трёх — `let searchTimeout: ReturnType<typeof setTimeout> | null = null` + `setTimeout(async () => api.get(...), 300)`. На `onUnmounted` или закрытии модалки `clearTimeout(searchTimeout)` не вызывается.

**Evidence:** в PointsAwardModal.vue нет `onUnmounted` хука вовсе. Watch(`props.isOpen`) сбрасывает state при reopen, но не отменяет pending timeout от прошлого открытия.

**Reproduction:** Юзер открывает модалку, печатает «aleks», закрывает за 200ms до того, как debounce-timeout сработает. Через 100ms timeout уходит в `api.get('members', { username: 'aleks' })`. Ответ приходит — пишется в `memberResults.value` уже отрендеренной модалки? Нет, компонент unmounted. Vue 3 refs не падают на set после unmount, но сетевой запрос состоялся впустую, и в дев-тулзах будет видно «hanging» request.

**Impact:** wasted network calls, лёгкая утечка. Не блокер.

**Suggested fix:**
```ts
onUnmounted(() => {
  if (searchTimeout) clearTimeout(searchTimeout)
})
```

**Severity:** LOW (UX-/perf-hygiene).

---

## [LOW] L3. Admin bot-команды (PeriodicCheck / Sweep / DryRun) — goroutine без `defer recover()`, panic кладёт бот

**Locations:**
- `backend/internal/bot/subscription.go:1178` (handleSubCheckAllCommand)
- `backend/internal/bot/subscription.go:1315` (handleSubMemberSweepCommand)
- `backend/internal/bot/subscription.go:1350` (handleSubKickDryCommand)

**Hypothesis:** Долгие админ-команды запускаются в фоне без recover. В `service/gamification_hook.go` уже есть пример правильного паттерна — `defer func(){ if r := recover(); r != nil { log.Printf(...) } }()`. В bot/subscription.go — нет.

**Impact:** один nil-deref в `subscriptionService.PeriodicCheck` или `SweepRealMembership` (например, при изменении схемы или редком race с deleted member) убивает весь bot-процесс. Бот на NL-сервере перезапустится по docker restart-policy, но в середине sweep оставит неконсистентный state, и админ потеряет отчёт.

**Suggested fix:** обернуть тело каждой goroutine в `defer recover()`-блок, как в `gamification_hook.go`.

**Severity:** LOW (admin-only attack surface, но реальный crash риск).

---

## [LOW] L4. Fire-and-forget handler goroutines без `defer recover()`

**Locations (выборка):**
- `backend/internal/handler/auth_token.go:192` (role-sync, см. L1)
- `backend/internal/handler/referal_link.go:170` (TrackConversion points-grant)
- `backend/internal/handler/events.go:235, 284, 323` (event-* нотификации)
- `backend/internal/handler/marketplace.go:202`
- `backend/internal/handler/bulk.go:43, 64, 85, 106, 136, 170, 197` (bulk-нотификации)
- `backend/internal/handler/task_exchange.go:226, 250, 295`

**Hypothesis:** Все эти фоновые goroutines обслуживают побочные эффекты (points, notifications, SSE-publish). Любой nil-pointer / index out of range в любой из них кладёт **весь backend-процесс**, а не только этот запрос.

**Evidence:** `grep -rn "recover()" backend/internal` находит recover ТОЛЬКО в `mentor.go` (2 места) и `gamification_hook.go` (3 места). 15+ других точек запуска goroutines — голые.

**Impact:** хрупкость прода. Один редкий код-путь (например, разлогиненный member, deleted чат) → крэш бэка → 502 для всех юзеров на ~5-10 сек, пока docker рестартует. Сейчас риск низкий, потому что большинство этих handlers зрелые и устаканенные, но любой новый feature-PR может это сломать.

**Suggested fix:** ввести helper `service.SafeGo(fn, label string)` (или middleware-обёртку), которая всегда оборачивает в recover + log.Printf. Заменить голые `go func()` на `service.SafeGo`.

**Severity:** LOW (системно, проактивное укрепление).

---

## Eliminated hypotheses (см. eliminated.md)

См. отдельный файл с гипотезами, которые я проверил и не подтвердил — important для будущих сессий, чтобы не перепроверять то же самое.

## Out of scope / design questions (НЕ баги)

- **PurchaseTierWithCredits** триггерит `AwardForFirstPurchase` + `AwardForRecurringPurchase` на покупке за credits, хотя cash inflow = 0. Если бизнес-логика говорит «реферальные выплаты — только за реальный кэш», то это баг M-severity. Если кредиты считаются эквивалентом покупки — design.
- Telegram **WebAppInitDataMaxAge = 24h** — задокументированное design решение «safety net при скриншот-атаке». Не баг.
- `?token=` в query для SSE — известный trade-off (EventSource не поддерживает custom headers); user уже знает (PR #341 коммит-сообщение). Acceptable, но повод посмотреть на secure-WebSocket-альтернативу позднее.
