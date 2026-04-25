package service

import (
	"strings"
	"testing"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

func kudosTablesTruncate(t *testing.T, db *gorm.DB) {
	testutil.TruncateAll(t, db, "kudos", "point_transactions", "members")
}

func TestKudosService_Send_HappyPath(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	kudosTablesTruncate(t, db)

	from := seedMemberWithRoles(t, db, 12001, "kudos_from", nil)
	to := seedMemberWithRoles(t, db, 12002, "kudos_to", nil)

	got, err := NewKudosService().Send(from.Id, to.Id, "спасибо")
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got.FromId != from.Id || got.ToId != to.Id {
		t.Errorf("from/to = %d/%d, want %d/%d", got.FromId, got.ToId, from.Id, to.Id)
	}
	if got.Message != "спасибо" {
		t.Errorf("Message = %q, want спасибо", got.Message)
	}
}

func TestKudosService_Send_SelfRejected(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	kudosTablesTruncate(t, db)

	m := seedMemberWithRoles(t, db, 12101, "self_kudos", nil)
	_, err := NewKudosService().Send(m.Id, m.Id, "себе")
	if err == nil {
		t.Fatal("ожидали ошибку отправки самому себе")
	}
	if !strings.Contains(err.Error(), "самому себе") {
		t.Errorf("неожиданное сообщение: %v", err)
	}
}

func TestKudosService_Send_EmptyMessageRejected(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	kudosTablesTruncate(t, db)

	from := seedMemberWithRoles(t, db, 12201, "empty_from", nil)
	to := seedMemberWithRoles(t, db, 12202, "empty_to", nil)
	_, err := NewKudosService().Send(from.Id, to.Id, "")
	if err == nil {
		t.Fatal("ожидали ошибку пустого сообщения")
	}
}

func TestKudosService_Send_DailyLimit(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	kudosTablesTruncate(t, db)

	from := seedMemberWithRoles(t, db, 12301, "limit_from", nil)
	tos := []*models.Member{
		seedMemberWithRoles(t, db, 12302, "to1", nil),
		seedMemberWithRoles(t, db, 12303, "to2", nil),
		seedMemberWithRoles(t, db, 12304, "to3", nil),
		seedMemberWithRoles(t, db, 12305, "to4", nil),
	}

	svc := NewKudosService()
	for i := 0; i < 3; i++ {
		if _, err := svc.Send(from.Id, tos[i].Id, "msg"); err != nil {
			t.Fatalf("send #%d: %v", i+1, err)
		}
	}

	_, err := svc.Send(from.Id, tos[3].Id, "msg")
	if err == nil {
		t.Fatal("ожидали ошибку лимита 3/сутки")
	}
	if !strings.Contains(err.Error(), "не более 3") {
		t.Errorf("неожиданное сообщение: %v", err)
	}
}

func TestKudosService_GetRecent_OrderedByCreatedDesc(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	kudosTablesTruncate(t, db)

	a := seedMemberWithRoles(t, db, 12401, "kudos_a", nil)
	b := seedMemberWithRoles(t, db, 12402, "kudos_b", nil)

	svc := NewKudosService()
	if _, err := svc.Send(a.Id, b.Id, "first"); err != nil {
		t.Fatalf("first send: %v", err)
	}
	if _, err := svc.Send(b.Id, a.Id, "second"); err != nil {
		t.Fatalf("second send: %v", err)
	}

	items, total, err := svc.GetRecent(20, 0)
	if err != nil {
		t.Fatalf("GetRecent: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(items) != 2 {
		t.Fatalf("items len = %d, want 2", len(items))
	}
	// Самый свежий — second.
	if items[0].Message != "second" {
		t.Errorf("первый item должен быть 'second', got %q", items[0].Message)
	}
}

func TestKudosService_GetRecent_LimitNormalization(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	kudosTablesTruncate(t, db)

	svc := NewKudosService()

	// limit=0 нормализуется в 20; пустые данные — total=0, items=[]
	items, total, err := svc.GetRecent(0, 0)
	if err != nil {
		t.Fatalf("GetRecent: %v", err)
	}
	if total != 0 || len(items) != 0 {
		t.Errorf("ожидали пустой результат для свежей БД; got total=%d items=%v", total, items)
	}
}
