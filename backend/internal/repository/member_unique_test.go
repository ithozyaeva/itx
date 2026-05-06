package repository

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsUsernameUniqueViolation(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"random error", errors.New("oops"), false},
		{
			"different unique (telegram_id)",
			&pgconn.PgError{Code: "23505", ConstraintName: "members_telegram_id_unique"},
			false,
		},
		{
			"username unique violation",
			&pgconn.PgError{Code: "23505", ConstraintName: usernameUniqueIndex},
			true,
		},
		{
			"different SQLSTATE on same constraint",
			&pgconn.PgError{Code: "23502", ConstraintName: usernameUniqueIndex},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUsernameUniqueViolation(tt.err); got != tt.want {
				t.Errorf("isUsernameUniqueViolation(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
