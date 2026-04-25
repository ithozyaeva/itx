package service

import (
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

// seedSimpleEvent создаёт минимальное Event в БД, привязывает указанных
// host'ов и member'ов, возвращает событие с заполненными связями.
func seedSimpleEvent(t *testing.T, db *gorm.DB, hosts []*models.Member, attendees []*models.Member) *models.Event {
	t.Helper()
	ev := &models.Event{
		Title:     "T9 test event",
		EventType: "online",
		Date:      time.Now().Add(-time.Hour),
	}
	if err := db.Create(ev).Error; err != nil {
		t.Fatalf("create event: %v", err)
	}
	if len(hosts) > 0 {
		if err := db.Model(ev).Association("Hosts").Append(hosts); err != nil {
			t.Fatalf("append hosts: %v", err)
		}
	}
	if len(attendees) > 0 {
		if err := db.Model(ev).Association("Members").Append(attendees); err != nil {
			t.Fatalf("append members: %v", err)
		}
	}

	// Перезагружаем, чтобы Hosts/Members были корректно подгружены —
	// у AwardEventPoints итерация идёт по этим срезам.
	full := &models.Event{}
	if err := db.Preload("Hosts").Preload("Members").First(full, ev.Id).Error; err != nil {
		t.Fatalf("reload event: %v", err)
	}
	return full
}

func countTransactions(t *testing.T, db *gorm.DB) int {
	t.Helper()
	var n int64
	if err := db.Model(&models.PointTransaction{}).Count(&n).Error; err != nil {
		t.Fatalf("count point_transactions: %v", err)
	}
	return int(n)
}

// TestAwardEventPoints_Idempotent — повторный вызов не дублирует
// транзакции (защита через WHERE NOT EXISTS).
func TestAwardEventPoints_Idempotent(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "point_transactions", "event_members", "event_hosts", "events", "members")

	host := seedMember(t, db, 5001)
	attendee := seedMember(t, db, 5002)
	ev := seedSimpleEvent(t, db, []*models.Member{host}, []*models.Member{attendee})

	svc := NewPointsService()

	if err := svc.AwardEventPoints(ev); err != nil {
		t.Fatalf("первый AwardEventPoints: %v", err)
	}
	first := countTransactions(t, db)
	if first != 2 {
		t.Fatalf("ожидали 2 записи (1 host + 1 attendee), получили %d", first)
	}

	if err := svc.AwardEventPoints(ev); err != nil {
		t.Fatalf("повторный AwardEventPoints: %v", err)
	}
	if got := countTransactions(t, db); got != first {
		t.Errorf("повторный вызов задублировал записи: до=%d, после=%d", first, got)
	}
}

// TestAwardEventPoints_Concurrent — несколько одновременных вызовов
// для одного и того же события не приводят к дублям (race на уровне
// БД с уникальным предикатом WHERE NOT EXISTS).
func TestAwardEventPoints_Concurrent(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "point_transactions", "event_members", "event_hosts", "events", "members")

	host := seedMember(t, db, 6001)
	attendees := []*models.Member{
		seedMember(t, db, 6002),
		seedMember(t, db, 6003),
		seedMember(t, db, 6004),
	}
	ev := seedSimpleEvent(t, db, []*models.Member{host}, attendees)

	svc := NewPointsService()

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = svc.AwardEventPoints(ev)
		}()
	}
	wg.Wait()

	expected := 1 + len(attendees) // host + members
	if got := countTransactions(t, db); got != expected {
		t.Errorf("ожидали %d уникальных транзакций после %d параллельных вызовов, получили %d",
			expected, goroutines, got)
	}
}

// TestCheckProfileComplete_Idempotent — одноразовый бонус за полный
// профиль начисляется ровно один раз, даже при повторном вызове.
func TestCheckProfileComplete_Idempotent(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "point_transactions", "members")

	bday := models.DateOnly(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	m := &models.Member{
		TelegramID: 7001,
		Username:   "complete_profile",
		FirstName:  "Full",
		LastName:   "Profile",
		Bio:        "I exist",
		Birthday:   &bday,
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}

	svc := NewPointsService()

	for i := 0; i < 3; i++ {
		svc.CheckProfileComplete(m)
	}

	var rows []models.PointTransaction
	if err := db.Where("member_id = ? AND reason = ?", m.Id, models.PointReasonProfileComplete).
		Find(&rows).Error; err != nil {
		t.Fatalf("query transactions: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("ожидали ровно одну запись profile_complete, получили %d", len(rows))
	}
}

// TestGetBalance_SumsTransactions — баланс = сумма amounts.
func TestGetBalance_SumsTransactions(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "point_transactions", "members")

	m := seedMember(t, db, 8001)
	svc := NewPointsService()

	// Три ручных начисления через AdminAwardPoints — каждый INSERT.
	for _, amt := range []int{10, 25, 5} {
		if err := svc.AdminAwardPoints(m.Id, amt, "test"); err != nil {
			t.Fatalf("AdminAwardPoints: %v", err)
		}
	}

	got, err := svc.GetBalance(m.Id)
	if err != nil {
		t.Fatalf("GetBalance: %v", err)
	}
	if got != 40 {
		t.Errorf("balance = %d, want 40", got)
	}
}
