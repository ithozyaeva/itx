package service

import (
	"sort"
	"testing"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

// seedMemberWithRoles создаёт member и привязывает к нему набор ролей.
// Member-таблица имеет legacy-колонку `role` (NOT NULL без default) — задаём
// её через основной insert, ролевые отношения идут отдельной таблицей.
func seedMemberWithRoles(t *testing.T, db *gorm.DB, telegramID int64, username string, roles []models.Role) *models.Member {
	t.Helper()
	m := &models.Member{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  "Test",
		LastName:   "Member",
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}
	for _, r := range roles {
		mr := &models.MemberRole{MemberId: m.Id, Role: r}
		if err := db.Create(mr).Error; err != nil {
			t.Fatalf("create member_role: %v", err)
		}
	}
	return m
}

func TestMemberService_GetByTelegramID(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "member_roles", "members")

	m := seedMemberWithRoles(t, db, 9001, "lookup_user", nil)
	svc := NewMemberService()

	got, err := svc.GetByTelegramID(9001)
	if err != nil {
		t.Fatalf("GetByTelegramID: %v", err)
	}
	if got.Id != m.Id {
		t.Errorf("Id = %d, want %d", got.Id, m.Id)
	}
	if got.Username != "lookup_user" {
		t.Errorf("Username = %q, want lookup_user", got.Username)
	}
}

func TestMemberService_GetByTelegramID_NotFound(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "member_roles", "members")

	if _, err := NewMemberService().GetByTelegramID(424242); err == nil {
		t.Error("ожидали ошибку 'member not found' для несуществующего id")
	}
}

func TestMemberService_GetByUsername(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "member_roles", "members")

	seedMemberWithRoles(t, db, 9101, "alice", nil)
	got, err := NewMemberService().GetByUsername("alice")
	if err != nil {
		t.Fatalf("GetByUsername: %v", err)
	}
	if got.TelegramID != 9101 {
		t.Errorf("TelegramID = %d, want 9101", got.TelegramID)
	}
}

func TestMemberService_IsAdminByTelegramID(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "member_roles", "members")

	seedMemberWithRoles(t, db, 9201, "admin_user", []models.Role{models.MemberRoleAdmin})
	seedMemberWithRoles(t, db, 9202, "regular_user", []models.Role{models.MemberRoleSubscriber})

	svc := NewMemberService()

	if !svc.IsAdminByTelegramID(9201) {
		t.Error("admin_user должен распознаваться как ADMIN")
	}
	if svc.IsAdminByTelegramID(9202) {
		t.Error("regular_user не должен распознаваться как ADMIN")
	}
	if svc.IsAdminByTelegramID(0) {
		t.Error("несуществующий ID не должен распознаваться как ADMIN")
	}
}

func TestMemberService_GetAdminTelegramIDs(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "member_roles", "members")

	seedMemberWithRoles(t, db, 9301, "admin1", []models.Role{models.MemberRoleAdmin})
	seedMemberWithRoles(t, db, 9302, "admin2", []models.Role{models.MemberRoleAdmin, models.MemberRoleSubscriber})
	seedMemberWithRoles(t, db, 9303, "subscriber", []models.Role{models.MemberRoleSubscriber})
	seedMemberWithRoles(t, db, 0, "noTg_admin", []models.Role{models.MemberRoleAdmin})

	got := NewMemberService().GetAdminTelegramIDs()
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })

	want := []int64{9301, 9302}
	if len(got) != len(want) {
		t.Fatalf("ids = %v, want %v", got, want)
	}
	for i, id := range want {
		if got[i] != id {
			t.Errorf("ids[%d] = %d, want %d", i, got[i], id)
		}
	}
}

func TestMemberService_GetPermissions_AdminInherits(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "member_roles", "members")

	m := seedMemberWithRoles(t, db, 9401, "perms_admin", []models.Role{models.MemberRoleAdmin})

	got, err := NewMemberService().GetPermissions(m.Id)
	if err != nil {
		t.Fatalf("GetPermissions: %v", err)
	}
	if !containsPermission(got, models.PermissionCanViewAdminPanel) {
		t.Errorf("ADMIN-роль должна включать can_view_admin_panel; got = %v", got)
	}
	// Проверяем недавно добавленный пермишн из миграции #288
	if !containsPermission(got, models.PermissionCanViewAdminFeedback) {
		t.Errorf("ADMIN-роль должна включать can_view_admin_feedback; got = %v", got)
	}
}

func TestMemberService_GetPermissions_NoRolesEmpty(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "member_roles", "members")

	m := seedMemberWithRoles(t, db, 9501, "no_roles", nil)

	got, err := NewMemberService().GetPermissions(m.Id)
	if err != nil {
		t.Fatalf("GetPermissions: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("без ролей должен быть пустой набор пермишенов; got = %v", got)
	}
}

func containsPermission(perms []models.Permission, p models.Permission) bool {
	for _, x := range perms {
		if x == p {
			return true
		}
	}
	return false
}
