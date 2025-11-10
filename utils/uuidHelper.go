package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUIDToString converts pgtype.UUID → string
func UUIDToString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return id.String()
}

// StringToUUID converts string → pgtype.UUID
func StringToUUID(id string) pgtype.UUID {
	if id == "" {
		return pgtype.UUID{Valid: false}
	}

	parsed, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{
		Bytes: parsed,
		Valid: true,
	}
}
