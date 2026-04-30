package service

import (
	"errors"
	"strings"
	"testing"

	"gorm.io/gorm"

	"ithozyeva/config"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/testutil"
)

// seedAIMaterial — мини-helper для создания «корневой» сущности под тесты
// CommentService. Возвращает ID, чтобы тесты не зависели от структуры
// AIMaterial напрямую.
func seedAIMaterial(t *testing.T, author *models.Member) int64 {
	t.Helper()
	svc := NewAIMaterialService()
	m, err := svc.Create(&models.CreateAIMaterialRequest{
		Title:        "Title",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
	}, author.Id)
	if err != nil {
		t.Fatalf("create ai-material: %v", err)
	}
	return m.Id
}

// seedAdminMember — создаёт члена с ADMIN-ролью. Используется в тестах,
// где нужна проверка bypass'ов.
func seedAdminMember(t *testing.T, db *gorm.DB, telegramID int64) *models.Member {
	t.Helper()
	m := seedMember(t, db, telegramID)
	if err := db.Create(&models.MemberRole{MemberId: m.Id, Role: models.MemberRoleAdmin}).Error; err != nil {
		t.Fatalf("seed admin role: %v", err)
	}
	m.Roles = []models.Role{models.MemberRoleAdmin}
	return m
}

// commentSvcWithMockedAIVisibility — собирает CommentService для тестов,
// где hideous tier-проверка не нужна (admin/master tier симулируются
// через role/SubscriptionUser в БД отдельно). Делаем checker'ы in-place,
// чтобы тест не зависел от реального состояния subscription_users.
func commentSvcWithMockedAIVisibility(visible bool) *CommentService {
	checkers := map[models.CommentEntityType]EntityVisibilityChecker{
		models.CommentEntityAIMaterial: func(_ int64, _ *models.Member) error {
			if !visible {
				return ErrEntityNotFound
			}
			return nil
		},
		models.CommentEntityEvent: func(_ int64, _ *models.Member) error {
			if !visible {
				return ErrEntityNotFound
			}
			return nil
		},
	}
	return NewCommentService(checkers)
}

func TestCommentService_AIMaterial_CRUD_AndCounter(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"comment_likes", "comments", "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 10001)
	commenter := seedMember(t, db, 10002)
	stranger := seedMember(t, db, 10003)
	admin := seedAdminMember(t, db, 10004)
	materialID := seedAIMaterial(t, author)
	svc := commentSvcWithMockedAIVisibility(true)

	c1, err := svc.Create(models.CommentEntityAIMaterial, materialID, commenter, "  Полезно  ")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if c1.Body != "Полезно" {
		t.Errorf("Body не trim: %q", c1.Body)
	}

	// Триггер обновил comments_count в ai_materials.
	aiSvc := NewAIMaterialService()
	got, err := aiSvc.GetByID(materialID, 0, true)
	if err != nil {
		t.Fatalf("get material: %v", err)
	}
	if got.CommentsCount != 1 {
		t.Errorf("CommentsCount = %d, want 1", got.CommentsCount)
	}

	// List
	list, total, err := svc.List(models.CommentEntityAIMaterial, materialID, commenter, false, 20, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Id != c1.Id {
		t.Errorf("list mismatch: total=%d, items=%+v", total, list)
	}

	// Update — чужой не может, автор может, админ может
	if _, err := svc.Update(c1.Id, stranger.Id, "Hacked", false); !errors.Is(err, ErrCommentForbidden) {
		t.Errorf("stranger update: want ErrForbidden, got %v", err)
	}
	if updated, err := svc.Update(c1.Id, commenter.Id, "Обновлённый", false); err != nil {
		t.Errorf("author update: %v", err)
	} else if updated.Body != "Обновлённый" {
		t.Errorf("Body не обновился: %q", updated.Body)
	}
	if _, err := svc.Update(c1.Id, stranger.Id, "Admin edit", true); err != nil {
		t.Errorf("admin update: %v", err)
	}

	// Delete — чужой не может, админ может (не автор)
	if err := svc.Delete(c1.Id, stranger.Id, false); !errors.Is(err, ErrCommentForbidden) {
		t.Errorf("stranger delete: want ErrForbidden, got %v", err)
	}
	if err := svc.Delete(c1.Id, admin.Id, true); err != nil {
		t.Errorf("admin delete: %v", err)
	}

	// Триггер обнулил счётчик.
	got, _ = aiSvc.GetByID(materialID, 0, true)
	if got.CommentsCount != 0 {
		t.Errorf("CommentsCount после delete = %d, want 0", got.CommentsCount)
	}
}

func TestCommentService_ToggleLike_RoundTrip_AndPropagatesViewerLiked(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"comment_likes", "comments", "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 10101)
	commenter := seedMember(t, db, 10102)
	liker := seedMember(t, db, 10103)
	materialID := seedAIMaterial(t, author)
	svc := commentSvcWithMockedAIVisibility(true)

	c, err := svc.Create(models.CommentEntityAIMaterial, materialID, commenter, "Полезно")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	liked, count, err := svc.ToggleLike(c.Id, liker)
	if err != nil {
		t.Fatalf("toggle 1: %v", err)
	}
	if !liked || count != 1 {
		t.Errorf("first toggle: liked=%v count=%d, want true 1", liked, count)
	}

	// liker видит Liked=true в листинге, автор — false.
	listLiker, _, _ := svc.List(models.CommentEntityAIMaterial, materialID, liker, false, 20, 0)
	if !listLiker[0].Liked {
		t.Error("liker should see Liked=true")
	}
	listAuthor, _, _ := svc.List(models.CommentEntityAIMaterial, materialID, author, false, 20, 0)
	if listAuthor[0].Liked {
		t.Error("author should see Liked=false")
	}

	liked, count, err = svc.ToggleLike(c.Id, liker)
	if err != nil {
		t.Fatalf("toggle 2: %v", err)
	}
	if liked || count != 0 {
		t.Errorf("second toggle: liked=%v count=%d, want false 0", liked, count)
	}

	if _, _, err := svc.ToggleLike(99999, liker); !errors.Is(err, ErrCommentNotFound) {
		t.Errorf("missing comment: want ErrCommentNotFound, got %v", err)
	}
}

func TestCommentService_ToggleLike_DeniedWhenEntityHidden(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"comment_likes", "comments", "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 10201)
	stranger := seedMember(t, db, 10202)
	materialID := seedAIMaterial(t, author)
	visibleSvc := commentSvcWithMockedAIVisibility(true)
	hiddenSvc := commentSvcWithMockedAIVisibility(false)

	c, err := visibleSvc.Create(models.CommentEntityAIMaterial, materialID, author, "сам себе")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, _, err := hiddenSvc.ToggleLike(c.Id, stranger); !errors.Is(err, ErrEntityNotFound) {
		t.Errorf("stranger like on hidden entity: want ErrEntityNotFound, got %v", err)
	}
}

// withSubscriptionGate временно переключает глобальный
// config.CFG.SubscriptionGateEnabled и восстанавливает его в Cleanup.
// Используем напрямую вместо config.LoadConfig — последний требует
// JWT_SECRET и других env-переменных, которых нет в тестовой среде.
// Тесты в одном пакете идут последовательно, поэтому глобал безопасен.
func withSubscriptionGate(t *testing.T, enabled bool) {
	t.Helper()
	if config.CFG == nil {
		// CFG nil в тестах, где LoadConfig не вызывался; ставим минимально
		// валидный, чтобы наши флаги имели куда писаться.
		config.CFG = &config.Config{}
		t.Cleanup(func() { config.CFG = nil })
	}
	prev := config.CFG.SubscriptionGateEnabled
	config.CFG.SubscriptionGateEnabled = enabled
	t.Cleanup(func() { config.CFG.SubscriptionGateEnabled = prev })
}

func TestAIMaterialVisibilityChecker_ForBidsLowerTier(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"comment_likes", "comments", "subscription_user_chat_access",
		"subscription_users", "ai_material_tags", "ai_materials", "members")
	withSubscriptionGate(t, true)

	author := seedMember(t, db, 10301)
	low := seedMember(t, db, 10302)
	materialID := seedAIMaterial(t, author)

	// foreman = level 2, не master+. Создаём прямо в БД, чтобы checker
	// видел реальный tier через GetUserEffectiveTierLevel.
	foremanTierID := uint(2)
	if err := db.Create(&models.SubscriptionUser{
		ID:             low.TelegramID,
		FullName:       "Low Tier",
		ManualTierID:   &foremanTierID,
		IsActive:       true,
	}).Error; err != nil {
		t.Fatalf("seed sub user: %v", err)
	}

	gate := NewSubscriptionTierGate(repository.NewSubscriptionRepository())
	aiSvc := NewAIMaterialService()
	checker := AIMaterialVisibilityChecker(aiSvc, gate)
	if err := checker(materialID, low); !errors.Is(err, ErrEntityNotFound) {
		t.Errorf("foreman: want ErrEntityNotFound, got %v", err)
	}

	// Поднимаем до master (level 3) — теперь должен пройти.
	masterTierID := uint(3)
	if err := db.Model(&models.SubscriptionUser{}).
		Where("id = ?", low.TelegramID).
		Update("manual_tier_id", masterTierID).Error; err != nil {
		t.Fatalf("upgrade tier: %v", err)
	}
	if err := checker(materialID, low); err != nil {
		t.Errorf("master: should pass, got %v", err)
	}
}

func TestAIMaterialVisibilityChecker_GateDisabled_AllowsAnyTier(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"comment_likes", "comments", "subscription_user_chat_access",
		"subscription_users", "ai_material_tags", "ai_materials", "members")
	// Главное условие теста — gate выключен. Раньше checker всё равно
	// проверял tier и расходился с middleware при выключенном флаге;
	// теперь оба идут через SubscriptionTierGate и поведение совпадает.
	withSubscriptionGate(t, false)

	author := seedMember(t, db, 10401)
	low := seedMember(t, db, 10402)
	materialID := seedAIMaterial(t, author)

	gate := NewSubscriptionTierGate(repository.NewSubscriptionRepository())
	checker := AIMaterialVisibilityChecker(NewAIMaterialService(), gate)
	if err := checker(materialID, low); err != nil {
		t.Errorf("gate disabled — low-tier should pass, got %v", err)
	}
}

func TestEventVisibilityChecker_OnlyChecksExistence(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"comment_likes", "comments", "event_event_tags", "event_hosts",
		"event_members", "events", "members")

	member := seedMember(t, db, 10501)

	evt := &models.Event{
		Title:       "Test event",
		Description: "test",
		PlaceType:   models.EventOnline,
		EventType:   "MEETUP",
	}
	if err := db.Create(evt).Error; err != nil {
		t.Fatalf("seed event: %v", err)
	}

	checker := EventVisibilityChecker(NewEventsService())
	if err := checker(evt.Id, member); err != nil {
		t.Errorf("existing event: should pass, got %v", err)
	}
	if err := checker(99999, member); !errors.Is(err, ErrEntityNotFound) {
		t.Errorf("missing event: want ErrEntityNotFound, got %v", err)
	}
}

func TestSharedComments_PreservesIDsAcrossEntities(t *testing.T) {
	// Регрессия: после миграции 20260430030000_create_shared_comments.sql
	// id комментов сохраняются (sequence подкручен через setval) и
	// FK comment_likes.comment_id → comments.id остаётся валидным.
	// Тест проверяет это через свежие insert'ы — а не через ручную
	// эмуляцию старой схемы (которая в testcontainer уже не существует).
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"comment_likes", "comments", "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 10601)
	liker := seedMember(t, db, 10602)
	materialID := seedAIMaterial(t, author)
	svc := commentSvcWithMockedAIVisibility(true)

	c1, err := svc.Create(models.CommentEntityAIMaterial, materialID, author, "first")
	if err != nil {
		t.Fatalf("create c1: %v", err)
	}
	c2, err := svc.Create(models.CommentEntityAIMaterial, materialID, author, "second")
	if err != nil {
		t.Fatalf("create c2: %v", err)
	}
	if c2.Id <= c1.Id {
		t.Errorf("sequence не работает: c1.id=%d c2.id=%d", c1.Id, c2.Id)
	}

	// Лайк на c1 — FK comment_likes.comment_id → comments.id должен сработать
	if _, _, err := svc.ToggleLike(c1.Id, liker); err != nil {
		t.Errorf("FK через comment_likes сломан: %v", err)
	}
}
