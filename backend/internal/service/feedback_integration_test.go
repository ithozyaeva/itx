package service

import (
	"fmt"
	"strings"
	"testing"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

// seedMember создаёт минимально валидного member для тестов и возвращает указатель.
// Username включает telegramID — без этого тест с двумя seedMember-вызовами
// внутри одного t.Name() ловил бы дубли по UNIQUE-индексу members_username_unique
// (см. миграцию 20260506000000_dedupe_and_unique_username.sql).
func seedMember(t *testing.T, db *gorm.DB, telegramID int64) *models.Member {
	t.Helper()
	m := &models.Member{
		TelegramID: telegramID,
		Username:   fmt.Sprintf("tester_%d_%s", telegramID, strings.ReplaceAll(t.Name(), "/", "_")),
		FirstName:  "Test",
		LastName:   "Member",
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}
	return m
}

func TestFeedbackService_Create_HappyPath(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "feedbacks", "members")

	member := seedMember(t, db, 1001)
	svc := NewFeedbackService()

	got, err := svc.Create(member, models.CreateFeedbackRequest{
		Score:   9,
		Comment: ptr("отлично"),
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.Score != 9 {
		t.Errorf("Score = %d, want 9", got.Score)
	}
	if got.Comment == nil || *got.Comment != "отлично" {
		t.Errorf("Comment = %v, want %q", got.Comment, "отлично")
	}
	if got.UserId == nil || *got.UserId != member.Id {
		t.Errorf("UserId = %v, want %d", got.UserId, member.Id)
	}
}

func TestFeedbackService_Create_RateLimitDailyOnePerUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "feedbacks", "members")

	member := seedMember(t, db, 1002)
	svc := NewFeedbackService()

	if _, err := svc.Create(member, models.CreateFeedbackRequest{Score: 5}); err != nil {
		t.Fatalf("первый Create должен пройти: %v", err)
	}

	_, err := svc.Create(member, models.CreateFeedbackRequest{Score: 8})
	if err == nil {
		t.Fatal("второй Create в сутки должен вернуть ошибку")
	}
	if !strings.Contains(err.Error(), "не более одного отзыва") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFeedbackService_Create_ScoreBoundaries(t *testing.T) {
	db := testutil.SetupTestDB(t)

	// Граничные значения 0 и 10 валидны (лимита 1/сутки на разных юзеров — нет).
	for i, score := range []int{0, 10} {
		testutil.TruncateAll(t, db, "feedbacks", "members")
		member := seedMember(t, db, int64(2000+i))
		_, err := NewFeedbackService().Create(member, models.CreateFeedbackRequest{Score: score})
		if err != nil {
			t.Errorf("score=%d должно проходить, error: %v", score, err)
		}
	}

	// Вне диапазона — отказ.
	testutil.TruncateAll(t, db, "feedbacks", "members")
	member := seedMember(t, db, 2099)
	for _, score := range []int{-1, 11, 100} {
		if _, err := NewFeedbackService().Create(member, models.CreateFeedbackRequest{Score: score}); err == nil {
			t.Errorf("score=%d должно отклоняться", score)
		}
	}
}

func TestFeedbackService_Create_CommentTooLong(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "feedbacks", "members")

	member := seedMember(t, db, 3001)
	long := strings.Repeat("a", 2001)
	_, err := NewFeedbackService().Create(member, models.CreateFeedbackRequest{
		Score:   5,
		Comment: &long,
	})
	if err == nil {
		t.Fatal("ожидали ошибку слишком длинного коммента")
	}
	if !strings.Contains(err.Error(), "слишком длинный") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFeedbackService_Create_EmptyCommentStoredAsNil(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "feedbacks", "members")

	member := seedMember(t, db, 3002)
	got, err := NewFeedbackService().Create(member, models.CreateFeedbackRequest{
		Score:   7,
		Comment: ptr("   "),
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.Comment != nil {
		t.Errorf("Comment = %v, want nil (whitespace must be nil)", got.Comment)
	}
}

func TestFeedbackService_List_ReturnsRecentFirst(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "feedbacks", "members")

	svc := NewFeedbackService()
	for i := 0; i < 3; i++ {
		// Each gets its own member to bypass the daily rate limit.
		m := seedMember(t, db, int64(4000+i))
		if _, err := svc.Create(m, models.CreateFeedbackRequest{Score: i + 1}); err != nil {
			t.Fatalf("seed feedback %d: %v", i, err)
		}
	}

	items, total, err := svc.List(20, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(items) != 3 {
		t.Fatalf("items = %d, want 3", len(items))
	}
	// Самый свежий — первый, score=3.
	if items[0].Score != 3 {
		t.Errorf("expected newest first (score=3), got %d", items[0].Score)
	}
}

func ptr[T any](v T) *T { return &v }
