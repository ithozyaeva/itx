package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

func eventTablesTruncate(t *testing.T, db *gorm.DB) {
	testutil.TruncateAll(t, db,
		"event_members", "event_hosts", "event_event_tags",
		"events", "event_tags", "members",
	)
}

func seedEvent(t *testing.T, db *gorm.DB, ev *models.Event) *models.Event {
	t.Helper()
	if ev.EventType == "" {
		ev.EventType = "online"
	}
	if err := db.Create(ev).Error; err != nil {
		t.Fatalf("create event: %v", err)
	}
	return ev
}

func TestEventsService_AddMember_HappyPath(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	eventTablesTruncate(t, db)

	member := seedMemberWithRoles(t, db, 11001, "attendee", nil)
	ev := seedEvent(t, db, &models.Event{
		Title: "Test event",
		Date:  time.Now().Add(24 * time.Hour),
	})

	updated, err := NewEventsService().AddMember(int(ev.Id), int(member.Id))
	if err != nil {
		t.Fatalf("AddMember: %v", err)
	}
	if len(updated.Members) != 1 || updated.Members[0].Id != member.Id {
		t.Errorf("ожидали один member.Id=%d, got %+v", member.Id, updated.Members)
	}
}

func TestEventsService_AddMember_ParticipantLimitReached(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	eventTablesTruncate(t, db)

	first := seedMemberWithRoles(t, db, 11101, "first", nil)
	second := seedMemberWithRoles(t, db, 11102, "second", nil)
	ev := seedEvent(t, db, &models.Event{
		Title:           "Limited event",
		Date:            time.Now().Add(24 * time.Hour),
		MaxParticipants: 1,
	})

	svc := NewEventsService()
	if _, err := svc.AddMember(int(ev.Id), int(first.Id)); err != nil {
		t.Fatalf("первого должны добавить: %v", err)
	}

	_, err := svc.AddMember(int(ev.Id), int(second.Id))
	if !errors.Is(err, ErrParticipantLimitReached) {
		t.Errorf("ожидали ErrParticipantLimitReached, got %v", err)
	}
}

func TestEventsService_RemoveMember(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	eventTablesTruncate(t, db)

	member := seedMemberWithRoles(t, db, 11201, "tobe_removed", nil)
	ev := seedEvent(t, db, &models.Event{Title: "Removable", Date: time.Now().Add(24 * time.Hour)})

	svc := NewEventsService()
	if _, err := svc.AddMember(int(ev.Id), int(member.Id)); err != nil {
		t.Fatalf("AddMember: %v", err)
	}

	updated, err := svc.RemoveMember(int(ev.Id), int(member.Id))
	if err != nil {
		t.Fatalf("RemoveMember: %v", err)
	}
	if len(updated.Members) != 0 {
		t.Errorf("ожидали пустой список members, got %+v", updated.Members)
	}
}

func TestEventsService_ResolveEventTags(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "event_event_tags", "events", "event_tags")

	// Засеиваем существующий тег.
	existing := models.EventTag{Name: "existing"}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("seed existing tag: %v", err)
	}

	svc := NewEventsService()

	// 1. Тег с Id > 0 — оставляется как есть.
	// 2. Тег без Id с уникальным именем — создаётся.
	// 3. Дубликат по имени — игнорируется.
	// 4. Дубликат по Id — игнорируется.
	// 5. Пустое/whitespace имя без Id — пропускается.
	resolved, err := svc.ResolveEventTags([]models.EventTag{
		{Id: existing.Id, Name: "existing"},
		{Name: "new-tag"},
		{Name: "new-tag"},                // дубликат → пропуск
		{Id: existing.Id, Name: "other"}, // дубликат по id → пропуск
		{Name: "   "},                    // пустое → пропуск
	})
	if err != nil {
		t.Fatalf("ResolveEventTags: %v", err)
	}

	if len(resolved) != 2 {
		t.Fatalf("ожидали 2 тега, got %d (%+v)", len(resolved), resolved)
	}

	gotNames := []string{resolved[0].Name, resolved[1].Name}
	hasExisting := contains(gotNames, "existing")
	hasNew := contains(gotNames, "new-tag")
	if !hasExisting || !hasNew {
		t.Errorf("ожидали существующий + new-tag; got %v", gotNames)
	}

	// Убедимся что new-tag действительно создан в БД.
	var count int64
	if err := db.Model(&models.EventTag{}).Where("name = ?", "new-tag").Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("new-tag должен быть создан один раз, got count=%d", count)
	}
}

func TestEventsService_GetFutureEvents_FiltersPastAndBrokenRepeating(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	eventTablesTruncate(t, db)

	now := time.Now()

	// 1. Прошедший — отбрасывается.
	seedEvent(t, db, &models.Event{
		Title: "Past",
		Date:  now.Add(-2 * time.Hour),
	})
	// 2. Будущий — попадает.
	seedEvent(t, db, &models.Event{
		Title: "Future",
		Date:  now.Add(2 * time.Hour),
	})
	// 3. Повторяющийся без RepeatPeriod — отбрасывается фильтром в Go.
	repeatPeriodNil := &models.Event{
		Title:       "Broken repeating",
		Date:        now.Add(-24 * time.Hour),
		IsRepeating: true,
	}
	seedEvent(t, db, repeatPeriodNil)
	// 4. Повторяющийся с RepeatPeriod — попадает.
	period := "weekly"
	endDate := now.Add(30 * 24 * time.Hour)
	seedEvent(t, db, &models.Event{
		Title:         "Healthy repeating",
		Date:          now.Add(-24 * time.Hour),
		IsRepeating:   true,
		RepeatPeriod:  &period,
		RepeatEndDate: &endDate,
	})

	got, err := NewEventsService().GetFutureEvents(now)
	if err != nil {
		t.Fatalf("GetFutureEvents: %v", err)
	}

	titles := make([]string, 0, len(got))
	for _, e := range got {
		titles = append(titles, e.Title)
	}

	if contains(titles, "Past") {
		t.Errorf("прошедшее событие не должно попадать в Future; got %v", titles)
	}
	if contains(titles, "Broken repeating") {
		t.Errorf("повторяющееся без RepeatPeriod должно отфильтровываться; got %v", titles)
	}
	if !contains(titles, "Future") || !contains(titles, "Healthy repeating") {
		t.Errorf("ожидали Future и Healthy repeating; got %v", titles)
	}
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle || strings.EqualFold(s, needle) {
			return true
		}
	}
	return false
}
