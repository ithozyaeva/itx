package handler

import (
	"reflect"
	"sort"
	"testing"

	"ithozyeva/internal/models"
)

// TestMergeSubscriptionRole_PreservesOtherRoles — регрессия Bug #1:
// при флипе chat-membership ADMIN/MENTOR/EVENT_MAKER не должны затираться.
// Раньше handler.Authenticate делал user.Roles = []Role{Subscriber} —
// и админ-подписчик после кратковременного выхода из anchor-чата терял ADMIN
// при следующем /authenticate.
func TestMergeSubscriptionRole_PreservesOtherRoles(t *testing.T) {
	cases := []struct {
		name         string
		input        []models.Role
		isSubscriber bool
		want         []models.Role
		wantChanged  bool
	}{
		{
			name:         "admin+subscriber stays in chat → no-op",
			input:        []models.Role{models.MemberRoleAdmin, models.MemberRoleSubscriber},
			isSubscriber: true,
			want:         []models.Role{models.MemberRoleAdmin, models.MemberRoleSubscriber},
			wantChanged:  false,
		},
		{
			name:         "admin+subscriber leaves chat → admin preserved, sub→unsub",
			input:        []models.Role{models.MemberRoleAdmin, models.MemberRoleSubscriber},
			isSubscriber: false,
			want:         []models.Role{models.MemberRoleAdmin, models.MemberRoleUnsubscriber},
			wantChanged:  true,
		},
		{
			name:         "mentor+unsubscriber joins chat → mentor preserved, unsub→sub",
			input:        []models.Role{models.MemberRoleMentor, models.MemberRoleUnsubscriber},
			isSubscriber: true,
			want:         []models.Role{models.MemberRoleMentor, models.MemberRoleSubscriber},
			wantChanged:  true,
		},
		{
			name:         "admin+mentor+event_maker+sub leaves chat → all extras preserved",
			input:        []models.Role{models.MemberRoleAdmin, models.MemberRoleMentor, models.MemberRoleEventMaker, models.MemberRoleSubscriber},
			isSubscriber: false,
			want:         []models.Role{models.MemberRoleAdmin, models.MemberRoleMentor, models.MemberRoleEventMaker, models.MemberRoleUnsubscriber},
			wantChanged:  true,
		},
		{
			name:         "missing flag is added",
			input:        []models.Role{models.MemberRoleAdmin},
			isSubscriber: true,
			want:         []models.Role{models.MemberRoleAdmin, models.MemberRoleSubscriber},
			wantChanged:  true,
		},
		{
			name:         "duplicate desired flag is deduplicated",
			input:        []models.Role{models.MemberRoleSubscriber, models.MemberRoleSubscriber, models.MemberRoleAdmin},
			isSubscriber: true,
			want:         []models.Role{models.MemberRoleSubscriber, models.MemberRoleAdmin},
			wantChanged:  true,
		},
		{
			name:         "corrupted state with both sub and unsub → cleaned up",
			input:        []models.Role{models.MemberRoleSubscriber, models.MemberRoleUnsubscriber, models.MemberRoleMentor},
			isSubscriber: true,
			want:         []models.Role{models.MemberRoleSubscriber, models.MemberRoleMentor},
			wantChanged:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, changed := mergeSubscriptionRole(tc.input, tc.isSubscriber)
			if changed != tc.wantChanged {
				t.Errorf("changed = %v, want %v", changed, tc.wantChanged)
			}
			// Порядок выходных ролей детерминированно зависит от входа,
			// но семантика — это set membership: сравниваем как sorted slice.
			if !sameRoleSet(got, tc.want) {
				gs := sortedRoleStrings(got)
				ws := sortedRoleStrings(tc.want)
				t.Errorf("roles = %v, want %v (got=%#v want=%#v)", gs, ws, got, tc.want)
			}
		})
	}
}

func sameRoleSet(a, b []models.Role) bool {
	if len(a) != len(b) {
		return false
	}
	as := sortedRoleStrings(a)
	bs := sortedRoleStrings(b)
	return reflect.DeepEqual(as, bs)
}

func sortedRoleStrings(rs []models.Role) []string {
	out := make([]string, len(rs))
	for i, r := range rs {
		out[i] = string(r)
	}
	sort.Strings(out)
	return out
}
