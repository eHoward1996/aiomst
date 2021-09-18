package db

import (
	"database/sql"
)

// integerQuery returns a single integer value from the input query
func (s *SqlBackend) integerQuery(query string, args ...interface{}) (int64, error) {
	// Perform query and fetch result
	result := struct {
		Int int64 `db:"int"`
	}{0}
	if err := s.db.Get(&result, query, args...); err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return result.Int, nil
}