package utils

import (
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

func Text(v string) pgtype.Text {
	if v == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: v,
		Valid:  true,
	}
}

func NullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
