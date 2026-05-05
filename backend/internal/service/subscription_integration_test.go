package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

// Тесты проверяют горячую логику subscription-сервиса:
//   - BuildContext (anchor map с дублями на тир, set anchor IDs, TiersDesc)
//   - resolveTierIDFromContext (priority desc, fail-stop при ошибке API)
//   - CheckAndSyncUser (grant, revoke с auto-kick=true/false, anchor-skip)
//   - DryRunCheckUser (никаких записей в БД и audit)
//
// Tier-ы (beginner/foreman/master/king) — из миграций; tier-маппинги и
// чаты сидим явно.

func subTablesTruncate(t *testing.T, db *gorm.DB) {
	testutil.TruncateAll(t, db,
		"subscription_user_chat_access",
		"subscription_audit_logs",
		"subscription_tier_chats",
		"subscription_users",
		"subscription_chats",
	)
}

func newTestSubService() *SubscriptionService {
	// nil redis — IsMember/InvalidateMemberCache nil-safe; тесты идут без
	// кэш-слоя, проверяя именно бизнес-логику.
	return NewSubscriptionService(nil)
}

func mustTier(t *testing.T, db *gorm.DB, slug string) *models.SubscriptionTier {
	t.Helper()
	var tier models.SubscriptionTier
	if err := db.Where("slug = ?", slug).First(&tier).Error; err != nil {
		t.Fatalf("tier %s not found: %v", slug, err)
	}
	return &tier
}

func seedSubChat(t *testing.T, db *gorm.DB, id int64, title string, anchorForTierID *uint) {
	t.Helper()
	if err := db.Create(&models.SubscriptionChat{
		ID:              id,
		Title:           title,
		ChatType:        "supergroup",
		AnchorForTierID: anchorForTierID,
	}).Error; err != nil {
		t.Fatalf("seed chat %d: %v", id, err)
	}
}

func linkChatToTier(t *testing.T, db *gorm.DB, chatID int64, tierID uint) {
	t.Helper()
	if err := db.Create(&models.SubscriptionTierChat{
		TierID: tierID,
		ChatID: chatID,
	}).Error; err != nil {
		t.Fatalf("link chat %d to tier %d: %v", chatID, tierID, err)
	}
}

func seedSubUser(t *testing.T, db *gorm.DB, id int64, resolved, manual *uint) {
	t.Helper()
	if err := db.Create(&models.SubscriptionUser{
		ID:             id,
		FullName:       "test",
		ResolvedTierID: resolved,
		ManualTierID:   manual,
		IsActive:       true,
	}).Error; err != nil {
		t.Fatalf("seed user %d: %v", id, err)
	}
}

func seedActiveAccess(t *testing.T, db *gorm.DB, userID, chatID int64) {
	t.Helper()
	if err := db.Create(&models.SubscriptionUserChatAccess{
		UserID:    userID,
		ChatID:    chatID,
		GrantedAt: time.Now(),
	}).Error; err != nil {
		t.Fatalf("seed access user=%d chat=%d: %v", userID, chatID, err)
	}
}

func seedManualAccess(t *testing.T, db *gorm.DB, userID, chatID int64) {
	t.Helper()
	if err := db.Create(&models.SubscriptionUserChatAccess{
		UserID:    userID,
		ChatID:    chatID,
		GrantedAt: time.Now(),
		IsManual:  true,
	}).Error; err != nil {
		t.Fatalf("seed manual access user=%d chat=%d: %v", userID, chatID, err)
	}
}

// staticChecker — мок MemberCheckFunc.
type staticChecker struct {
	members map[string]bool
	errs    map[string]error
	calls   []string
}

func keyOf(chatID, userID int64) string {
	return fmt.Sprintf("%d:%d", chatID, userID)
}

func (m *staticChecker) check(chatID, userID int64) (bool, error) {
	key := keyOf(chatID, userID)
	m.calls = append(m.calls, key)
	if err, ok := m.errs[key]; ok {
		return false, err
	}
	return m.members[key], nil
}

// --- BuildContext ---

func TestSubscriptionBuildContext(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	beginner := mustTier(t, db, "beginner")
	foreman := mustTier(t, db, "foreman")
	master := mustTier(t, db, "master")

	// Foreman имеет два anchor-чата — проверяем, что оба попадают в map.
	seedSubChat(t, db, -1001, "anchor-beginner", &beginner.ID)
	seedSubChat(t, db, -1002, "anchor-foreman-a", &foreman.ID)
	seedSubChat(t, db, -1003, "anchor-foreman-b", &foreman.ID)
	seedSubChat(t, db, -1004, "anchor-master", &master.ID)
	seedSubChat(t, db, -2001, "content-only", nil)

	svc := newTestSubService()
	ctx, err := svc.BuildContext()
	if err != nil {
		t.Fatalf("BuildContext: %v", err)
	}

	if len(ctx.AnchorChatIDs) != 4 {
		t.Errorf("AnchorChatIDs len = %d, want 4 (got %v)", len(ctx.AnchorChatIDs), ctx.AnchorChatIDs)
	}
	if ctx.AnchorChatIDs[-2001] {
		t.Errorf("content chat -2001 leaked into AnchorChatIDs")
	}

	foremanAnchors := ctx.AnchorChatsByTier[foreman.ID]
	if len(foremanAnchors) != 2 {
		t.Errorf("AnchorChatsByTier[foreman] len = %d, want 2 (got %v)",
			len(foremanAnchors), foremanAnchors)
	}

	if len(ctx.TiersDesc) < 3 {
		t.Fatalf("TiersDesc has %d tiers, want >=3", len(ctx.TiersDesc))
	}
	for i := 1; i < len(ctx.TiersDesc); i++ {
		if ctx.TiersDesc[i-1].Level < ctx.TiersDesc[i].Level {
			t.Errorf("TiersDesc not sorted desc at idx %d: %d < %d",
				i, ctx.TiersDesc[i-1].Level, ctx.TiersDesc[i].Level)
		}
	}
}

// --- resolveTierIDFromContext ---

func TestResolveTierFromContextPriorityDesc(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	beginner := mustTier(t, db, "beginner")
	master := mustTier(t, db, "master")

	seedSubChat(t, db, -100, "anchor-beginner", &beginner.ID)
	seedSubChat(t, db, -300, "anchor-master", &master.ID)

	const userID int64 = 42
	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-100, userID): true,
			keyOf(-300, userID): true, // юзер в обоих anchor'ах
		},
	}

	svc := newTestSubService()
	subCtx, err := svc.BuildContext()
	if err != nil {
		t.Fatalf("BuildContext: %v", err)
	}

	tierID, err := svc.resolveTierIDFromContext(userID, mock.check, subCtx)
	if err != nil {
		t.Fatalf("resolveTierIDFromContext: %v", err)
	}
	if tierID == nil || *tierID != master.ID {
		t.Errorf("expected master tier %d, got %v", master.ID, tierID)
	}
	for _, c := range mock.calls {
		if c == keyOf(-100, userID) {
			t.Errorf("beginner anchor should not be checked after master matched, calls=%v",
				mock.calls)
		}
	}
}

func TestResolveTierFromContextDuplicateAnchorsPerTier(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	foreman := mustTier(t, db, "foreman")

	seedSubChat(t, db, -100, "anchor-foreman-a", &foreman.ID)
	seedSubChat(t, db, -101, "anchor-foreman-b", &foreman.ID)

	const userID int64 = 42
	// Юзер только во втором anchor'е foreman'а — без поддержки массива
	// в anchor-map (что было в одной из прошлых регрессий) тир бы не
	// определился, потому что первый anchor вернул бы false и второй
	// игнорировался бы.
	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-100, userID): false,
			keyOf(-101, userID): true,
		},
	}

	svc := newTestSubService()
	subCtx, err := svc.BuildContext()
	if err != nil {
		t.Fatalf("BuildContext: %v", err)
	}

	tierID, err := svc.resolveTierIDFromContext(userID, mock.check, subCtx)
	if err != nil {
		t.Fatalf("resolveTierIDFromContext: %v", err)
	}
	if tierID == nil || *tierID != foreman.ID {
		t.Errorf("expected foreman tier %d, got %v", foreman.ID, tierID)
	}
}

func TestResolveTierFromContextErrorStopsCascade(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	beginner := mustTier(t, db, "beginner")
	master := mustTier(t, db, "master")

	seedSubChat(t, db, -100, "anchor-beginner", &beginner.ID)
	seedSubChat(t, db, -300, "anchor-master", &master.ID)

	const userID int64 = 42
	apiErr := errors.New("Telegram rate limit")
	mock := &staticChecker{
		errs: map[string]error{
			keyOf(-300, userID): apiErr,
		},
		members: map[string]bool{
			// Если бы fail-stop не работал, beginner вернул бы false
			// и тир «потерялся» бы в nil — мы вместо этого распространяем
			// ошибку, чтобы вызывающий слой пропустил юзера на этом проходе.
			keyOf(-100, userID): false,
		},
	}

	svc := newTestSubService()
	subCtx, err := svc.BuildContext()
	if err != nil {
		t.Fatalf("BuildContext: %v", err)
	}

	tierID, err := svc.resolveTierIDFromContext(userID, mock.check, subCtx)
	if err == nil {
		t.Fatalf("expected error from resolveTierIDFromContext, got nil (tierID=%v)", tierID)
	}
	if !errors.Is(err, apiErr) {
		t.Errorf("expected wrapped apiErr, got %v", err)
	}
	if tierID != nil {
		t.Errorf("tierID must be nil on error, got %v", *tierID)
	}
}

// --- CheckAndSyncUser ---

func TestCheckAndSyncUserGrantsContentForResolvedTier(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	master := mustTier(t, db, "master")

	seedSubChat(t, db, -300, "anchor-master", &master.ID)
	seedSubChat(t, db, -500, "content-1", nil)
	seedSubChat(t, db, -501, "content-2", nil)
	linkChatToTier(t, db, -500, master.ID)
	linkChatToTier(t, db, -501, master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, nil, nil)

	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-300, userID): true,
		},
	}
	inviteCalls := 0
	stubInvite := func(int64) (string, error) {
		inviteCalls++
		return "https://t.me/+x", nil
	}

	svc := newTestSubService()
	result, err := svc.CheckAndSyncUser(userID, mock.check, stubInvite,
		func(int64, int64) bool { return true })
	if err != nil {
		t.Fatalf("CheckAndSyncUser: %v", err)
	}

	if len(result.Granted) != 2 {
		t.Errorf("expected 2 grants, got %d (%+v)", len(result.Granted), result.Granted)
	}
	if inviteCalls != 2 {
		t.Errorf("expected 2 invite-link calls, got %d", inviteCalls)
	}

	var u models.SubscriptionUser
	if err := db.First(&u, userID).Error; err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if u.ResolvedTierID == nil || *u.ResolvedTierID != master.ID {
		t.Errorf("resolved_tier must be master=%d, got %v", master.ID, u.ResolvedTierID)
	}

	var grantAudits int64
	db.Model(&models.SubscriptionAuditLog{}).
		Where("user_id = ? AND action = ?", userID, "grant").
		Count(&grantAudits)
	if grantAudits != 2 {
		t.Errorf("expected 2 'grant' audits, got %d", grantAudits)
	}
}

func TestCheckAndSyncUserKickFalseKeepsAccess(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	master := mustTier(t, db, "master")
	foreman := mustTier(t, db, "foreman")

	seedSubChat(t, db, -100, "anchor-foreman", &foreman.ID)
	seedSubChat(t, db, -300, "anchor-master", &master.ID)
	seedSubChat(t, db, -500, "content-master", nil)
	linkChatToTier(t, db, -500, master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, &master.ID, nil)
	seedActiveAccess(t, db, userID, -500)

	// Юзер опустился с master до foreman.
	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-100, userID): true,
			keyOf(-300, userID): false,
		},
	}

	svc := newTestSubService()
	result, err := svc.CheckAndSyncUser(userID, mock.check,
		func(int64) (string, error) { return "", nil },
		func(int64, int64) bool { return false }, // SUBSCRIPTION_AUTO_KICK_ENABLED=false
	)
	if err != nil {
		t.Fatalf("CheckAndSyncUser: %v", err)
	}
	if len(result.Revoked) != 0 {
		t.Errorf("Revoked must be empty when kick=false: %v", result.Revoked)
	}

	var active int64
	db.Model(&models.SubscriptionUserChatAccess{}).
		Where("user_id = ? AND chat_id = ? AND revoked_at IS NULL",
			userID, int64(-500)).
		Count(&active)
	if active != 1 {
		t.Errorf("access must remain active when kick=false, got count=%d", active)
	}

	var revokeAudits int64
	db.Model(&models.SubscriptionAuditLog{}).
		Where("user_id = ? AND action = ?", userID, "revoke").
		Count(&revokeAudits)
	if revokeAudits != 0 {
		t.Errorf("no 'revoke' audit must be created when kick=false, got %d", revokeAudits)
	}
}

func TestCheckAndSyncUserKickTrueRevokesAccess(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	master := mustTier(t, db, "master")
	foreman := mustTier(t, db, "foreman")

	seedSubChat(t, db, -100, "anchor-foreman", &foreman.ID)
	seedSubChat(t, db, -300, "anchor-master", &master.ID)
	seedSubChat(t, db, -500, "content-master", nil)
	linkChatToTier(t, db, -500, master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, &master.ID, nil)
	seedActiveAccess(t, db, userID, -500)

	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-100, userID): true,
			keyOf(-300, userID): false,
		},
	}

	var kicked []int64
	stubKick := func(chatID, _ int64) bool {
		kicked = append(kicked, chatID)
		return true
	}

	svc := newTestSubService()
	result, err := svc.CheckAndSyncUser(userID, mock.check,
		func(int64) (string, error) { return "", nil }, stubKick)
	if err != nil {
		t.Fatalf("CheckAndSyncUser: %v", err)
	}

	if len(result.Revoked) != 1 || result.Revoked[0] != -500 {
		t.Errorf("expected revoke of -500, got %v", result.Revoked)
	}
	if len(kicked) != 1 || kicked[0] != -500 {
		t.Errorf("kick stub not called for -500, got %v", kicked)
	}

	var active int64
	db.Model(&models.SubscriptionUserChatAccess{}).
		Where("user_id = ? AND chat_id = ? AND revoked_at IS NULL",
			userID, int64(-500)).
		Count(&active)
	if active != 0 {
		t.Errorf("access must be revoked when kick=true, got count=%d", active)
	}
}

func TestCheckAndSyncUserAnchorAccessNotRevoked(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	master := mustTier(t, db, "master")
	seedSubChat(t, db, -300, "anchor-master", &master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, &master.ID, nil)
	// Симулируем legacy state: до фикса anchor мог попасть в access-таблицу.
	// После фикса каждый periodic должен skip'нуть его в revoke-loop.
	seedActiveAccess(t, db, userID, -300)

	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-300, userID): true,
		},
	}
	var kicked []int64
	stubKick := func(chatID, _ int64) bool {
		kicked = append(kicked, chatID)
		return true
	}

	svc := newTestSubService()
	result, err := svc.CheckAndSyncUser(userID, mock.check,
		func(int64) (string, error) { return "", nil }, stubKick)
	if err != nil {
		t.Fatalf("CheckAndSyncUser: %v", err)
	}

	if len(result.Revoked) != 0 {
		t.Errorf("anchor must be skipped from revoke, got %v", result.Revoked)
	}
	if len(kicked) != 0 {
		t.Errorf("kickUser must not be called for anchor, got %v", kicked)
	}
}

// --- DryRunCheckUser ---

func TestDryRunCheckUserNoWrites(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	master := mustTier(t, db, "master")
	foreman := mustTier(t, db, "foreman")

	seedSubChat(t, db, -100, "anchor-foreman", &foreman.ID)
	seedSubChat(t, db, -300, "anchor-master", &master.ID)
	seedSubChat(t, db, -500, "content-master", nil)
	linkChatToTier(t, db, -500, master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, &master.ID, nil)
	seedActiveAccess(t, db, userID, -500)

	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-100, userID): true,
			keyOf(-300, userID): false,
		},
	}

	svc := newTestSubService()
	res, err := svc.DryRunCheckUser(userID, mock.check)
	if err != nil {
		t.Fatalf("DryRunCheckUser: %v", err)
	}

	if len(res.WouldRevoke) != 1 || res.WouldRevoke[0] != -500 {
		t.Errorf("WouldRevoke = %v, want [-500]", res.WouldRevoke)
	}

	var u models.SubscriptionUser
	if err := db.First(&u, userID).Error; err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if u.ResolvedTierID == nil || *u.ResolvedTierID != master.ID {
		t.Errorf("DryRun changed resolved_tier from master=%d to %v", master.ID, u.ResolvedTierID)
	}

	var revoked int64
	db.Model(&models.SubscriptionUserChatAccess{}).
		Where("user_id = ? AND revoked_at IS NOT NULL", userID).
		Count(&revoked)
	if revoked != 0 {
		t.Errorf("DryRun revoked %d accesses, expected 0", revoked)
	}

	var audits int64
	db.Model(&models.SubscriptionAuditLog{}).
		Where("user_id = ?", userID).
		Count(&audits)
	if audits != 0 {
		t.Errorf("DryRun wrote %d audit logs, expected 0", audits)
	}
}

// --- PeriodicCheck skip-path ---

func TestPeriodicCheckSkipsUserOnTelegramError(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	master := mustTier(t, db, "master")
	seedSubChat(t, db, -300, "anchor-master", &master.ID)

	// Два юзера: первому Telegram отвечает ОК, второму — error.
	// Цель — убедиться, что error на одном не останавливает loop, и
	// проблемный юзер попадает в Skipped, а не молча теряется.
	const okUser int64 = 100
	const failUser int64 = 200
	seedSubUser(t, db, okUser, &master.ID, nil)
	seedSubUser(t, db, failUser, &master.ID, nil)

	apiErr := errors.New("Too Many Requests: retry after 30")
	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-300, okUser): true,
		},
		errs: map[string]error{
			keyOf(-300, failUser): apiErr,
		},
	}

	svc := newTestSubService()
	report := svc.PeriodicCheck(mock.check,
		func(int64) (string, error) { return "", nil },
		func(int64, int64) bool { return true },
		0,
	)

	if report.UsersTotal != 2 {
		t.Errorf("UsersTotal = %d, want 2", report.UsersTotal)
	}
	if len(report.Skipped) != 1 || report.Skipped[0] != failUser {
		t.Errorf("Skipped = %v, want [%d]", report.Skipped, failUser)
	}

	// Убедимся, что resolved_tier у failUser НЕ был обнулён —
	// скип означает «следующий тикер попробует снова», а не «понизить тир».
	var u models.SubscriptionUser
	if err := db.First(&u, failUser).Error; err != nil {
		t.Fatalf("reload failUser: %v", err)
	}
	if u.ResolvedTierID == nil || *u.ResolvedTierID != master.ID {
		t.Errorf("failUser resolved_tier должен остаться master=%d, got %v",
			master.ID, u.ResolvedTierID)
	}
}

// --- Sweep ChecksFailed counter ---

func TestSweepRealMembershipCountsFailures(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	master := mustTier(t, db, "master")

	seedSubChat(t, db, -300, "anchor-master", &master.ID) // anchor — sweep пропустит
	seedSubChat(t, db, -500, "content-1", nil)
	seedSubChat(t, db, -501, "content-2", nil)
	linkChatToTier(t, db, -500, master.ID)
	linkChatToTier(t, db, -501, master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, &master.ID, nil)

	apiErr := errors.New("rate limit")
	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-500, userID): true, // в content-1 — member, выдадим access
		},
		errs: map[string]error{
			keyOf(-501, userID): apiErr, // content-2 упал → ChecksFailed++
		},
	}

	svc := newTestSubService()
	stats, err := svc.SweepRealMembership(mock.check, 0)
	if err != nil {
		t.Fatalf("SweepRealMembership: %v", err)
	}

	// 2 content-чата × 1 юзер = 2 попытки; одна upalа.
	if stats.ChecksPerformed != 2 {
		t.Errorf("ChecksPerformed = %d, want 2", stats.ChecksPerformed)
	}
	if stats.ChecksFailed != 1 {
		t.Errorf("ChecksFailed = %d, want 1", stats.ChecksFailed)
	}
	if stats.AccessGranted != 1 {
		t.Errorf("AccessGranted = %d, want 1 (только content-1)", stats.AccessGranted)
	}

	// Убедимся, что чат с упавшим check'ом НЕ был revoked (хотя и не member).
	// Раньше err трактовался как «не member» и сносил access честных юзеров
	// при rate-limit'е.
	var revokedCount int64
	db.Model(&models.SubscriptionUserChatAccess{}).
		Where("user_id = ? AND chat_id = ? AND revoked_at IS NOT NULL",
			userID, int64(-501)).
		Count(&revokedCount)
	if revokedCount != 0 {
		t.Errorf("content-2 access не должен быть revoked при error, got revoked count=%d",
			revokedCount)
	}
}

// --- DryRunCheckUser anchor skip (existing) ---

func TestDryRunCheckUserSkipsAnchorFromRevoke(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	master := mustTier(t, db, "master")
	seedSubChat(t, db, -300, "anchor-master", &master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, &master.ID, nil)
	seedActiveAccess(t, db, userID, -300)

	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-300, userID): true,
		},
	}

	svc := newTestSubService()
	res, err := svc.DryRunCheckUser(userID, mock.check)
	if err != nil {
		t.Fatalf("DryRunCheckUser: %v", err)
	}

	for _, cid := range res.WouldRevoke {
		if cid == -300 {
			t.Errorf("anchor -300 must not appear in WouldRevoke (got %v)", res.WouldRevoke)
		}
	}
}

// --- Manual grants ---

func TestCheckAndSyncUserManualAccessNotRevoked(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	beginner := mustTier(t, db, "beginner")
	master := mustTier(t, db, "master")

	// Юзер-foreman, в content-master он попасть не должен по тиру.
	// Но админ его туда добавил вручную — не должен быть кикнут.
	seedSubChat(t, db, -100, "anchor-beginner", &beginner.ID)
	seedSubChat(t, db, -500, "content-master", nil)
	linkChatToTier(t, db, -500, master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, &beginner.ID, nil)
	seedManualAccess(t, db, userID, -500)

	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-100, userID): true, // beginner anchor — определяет тир
		},
	}

	var kicked []int64
	stubKick := func(chatID, _ int64) bool {
		kicked = append(kicked, chatID)
		return true
	}

	svc := newTestSubService()
	result, err := svc.CheckAndSyncUser(userID, mock.check,
		func(int64) (string, error) { return "", nil }, stubKick)
	if err != nil {
		t.Fatalf("CheckAndSyncUser: %v", err)
	}

	if len(result.Revoked) != 0 {
		t.Errorf("manual-access -500 не должен попасть в Revoked, got %v", result.Revoked)
	}
	if len(kicked) != 0 {
		t.Errorf("kickUser не должен вызываться для manual-чата, got %v", kicked)
	}

	// Запись остаётся active.
	var active int64
	db.Model(&models.SubscriptionUserChatAccess{}).
		Where("user_id = ? AND chat_id = ? AND revoked_at IS NULL",
			userID, int64(-500)).
		Count(&active)
	if active != 1 {
		t.Errorf("manual-access должен остаться active, got count=%d", active)
	}
}

func TestDryRunCheckUserSkipsManualFromRevoke(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	beginner := mustTier(t, db, "beginner")
	master := mustTier(t, db, "master")

	seedSubChat(t, db, -100, "anchor-beginner", &beginner.ID)
	seedSubChat(t, db, -500, "content-master", nil)
	linkChatToTier(t, db, -500, master.ID)

	const userID int64 = 42
	seedSubUser(t, db, userID, &beginner.ID, nil)
	seedManualAccess(t, db, userID, -500)

	mock := &staticChecker{
		members: map[string]bool{
			keyOf(-100, userID): true,
		},
	}

	svc := newTestSubService()
	res, err := svc.DryRunCheckUser(userID, mock.check)
	if err != nil {
		t.Fatalf("DryRunCheckUser: %v", err)
	}

	for _, cid := range res.WouldRevoke {
		if cid == -500 {
			t.Errorf("manual -500 не должен попасть в WouldRevoke (got %v)", res.WouldRevoke)
		}
	}
}

func TestSyncContentJoinPropagatesIsManual(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	seedSubChat(t, db, -500, "content-1", nil)

	const userID int64 = 42
	svc := newTestSubService()

	if err := svc.SyncContentJoin(userID, -500, nil, "test", true); err != nil {
		t.Fatalf("SyncContentJoin manual=true: %v", err)
	}

	var rec models.SubscriptionUserChatAccess
	if err := db.Where("user_id = ? AND chat_id = ?", userID, int64(-500)).
		First(&rec).Error; err != nil {
		t.Fatalf("reload access: %v", err)
	}
	if !rec.IsManual {
		t.Errorf("IsManual должен быть true после manual SyncContentJoin")
	}

	// Повторный auto-grant не должен понизить флаг (sweep-ситуация).
	if err := svc.GrantAccess(userID, -500, false); err != nil {
		t.Fatalf("GrantAccess auto=false поверх manual: %v", err)
	}
	if err := db.Where("user_id = ? AND chat_id = ?", userID, int64(-500)).
		First(&rec).Error; err != nil {
		t.Fatalf("reload access after auto: %v", err)
	}
	if !rec.IsManual {
		t.Errorf("auto-grant поверх manual=true не должен понижать флаг")
	}
}

func TestGrantAccessManualUpgradeOnly(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	subTablesTruncate(t, db)

	seedSubChat(t, db, -500, "content-1", nil)
	const userID int64 = 42
	seedSubUser(t, db, userID, nil, nil)

	svc := newTestSubService()

	// Auto-grant первым.
	if err := svc.GrantAccess(userID, -500, false); err != nil {
		t.Fatalf("first auto-grant: %v", err)
	}

	// Затем юзера добавил админ вручную (chat_member event from!=user) —
	// upsert переходит в manual=true.
	if err := svc.GrantAccess(userID, -500, true); err != nil {
		t.Fatalf("manual-grant on existing auto: %v", err)
	}

	var rec models.SubscriptionUserChatAccess
	if err := db.Where("user_id = ? AND chat_id = ?", userID, int64(-500)).
		First(&rec).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !rec.IsManual {
		t.Errorf("manual-grant поверх auto должен поднять флаг")
	}
}
