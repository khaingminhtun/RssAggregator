package utils

import "github.com/jackc/pgx/v5/pgtype"

func Text(v string) pgtype.Text {
	if v == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: v,
		Valid:  true,
	}
}
