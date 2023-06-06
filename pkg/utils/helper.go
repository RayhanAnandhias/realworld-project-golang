package utils

import (
	"database/sql"
	"regexp"
	"strings"
)

func NewNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func GenerateSlug(s string) string {
	stringRegex := regexp.MustCompile(`\s`)
	trimmedString := strings.TrimSpace(s)
	loweredString := strings.ToLower(trimmedString)
	return stringRegex.ReplaceAllString(loweredString, "-")
}
