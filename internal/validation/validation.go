package validation

import "fmt"

const maxSQLLength = 10000

// ValidateNoControlChars rejects ASCII control characters (0x00-0x1F except tab/newline/CR, and 0x7F).
func ValidateNoControlChars(input string, fieldName string) error {
	for i, r := range input {
		if r == '\t' || r == '\n' || r == '\r' {
			continue
		}
		if (r >= 0x00 && r <= 0x1F) || r == 0x7F {
			return fmt.Errorf("%s contains invalid control character at position %d", fieldName, i)
		}
	}
	return nil
}

// ValidateSQL validates a SQL query string.
func ValidateSQL(sql string) error {
	if err := ValidateNoControlChars(sql, "sql"); err != nil {
		return err
	}
	if len(sql) > maxSQLLength {
		return fmt.Errorf("sql query is too long (maximum %d characters, got %d)", maxSQLLength, len(sql))
	}
	return nil
}

// ValidateSearchQuery validates a search query string.
func ValidateSearchQuery(query string) error {
	return ValidateNoControlChars(query, "search query")
}
