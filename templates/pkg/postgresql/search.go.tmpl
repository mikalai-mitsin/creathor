package postgresql

import (
	"fmt"
	"strings"
)

type Search struct {
	Lang   string
	Fields []string
	Query  string
}

// nolint:stylecheck
func (s Search) ToSql() (sql string, args []interface{}, err error) {
	if s.Lang == "" {
		s.Lang = "russian"
	}
	vector := "to_tsvector('%s', %s)"
	vector = fmt.Sprintf(vector, s.Lang, strings.Join(s.Fields, " || ' ' || "))
	query := "plainto_tsquery('%s', '%s')"
	query = fmt.Sprintf(query, s.Lang, s.Query)
	return fmt.Sprintf("%s @@ %s", vector, query), nil, nil
}
