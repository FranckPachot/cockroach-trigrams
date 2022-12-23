package trgm

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
)

type Querier interface {
	Description() string
	Query(context.Context, string) ([]string, error)
}

type Naive struct {
	Conn                    *pgx.Conn
	ForceIndex              bool
	AnalyzedField           bool
	AnalyzedQuery           bool
	TuneSimilarityThreshold bool
}

func (s *Naive) Description() string {
	d := ""
	if s.AnalyzedQuery {
		d += "Analyzed query "
	} else {
		d += "Raw query "
	}
	d += "% operator on "
	if s.AnalyzedField {
		d += "analyzed text "
	} else {
		d += "raw text "
	}
	if s.ForceIndex {
		d += "with forced GIN index usage "
	}
	return strings.TrimRight(d, " ")
}

func (s *Naive) Query(ctx context.Context, query string) ([]string, error) {
	idx := ""
	field := "name"
	if s.AnalyzedField {
		field = "analyzed"
	}
	if s.ForceIndex {
		if s.AnalyzedField {
			idx += `@foods_analyzed_idx`
		} else {
			idx += `@foods_name_idx`
		}
	}
	sql := fmt.Sprintf(
		`SELECT name FROM foods%s WHERE %s %% $1 ORDER BY similarity(%s, $1) LIMIT 10`,
		idx,
		field,
		field,
	)

	if s.AnalyzedQuery {
		query = analyze(query)
	}

	if s.TuneSimilarityThreshold {
		s.Conn.Query(ctx, `SET pg_trgm.similarity_threshold = $1`, 0.2)
	}

	rows, err := s.Conn.Query(ctx, sql, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	i := 0
	results := make([]string, 10)
	for ; rows.Next(); i++ {
		if err := rows.Scan(&results[i]); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return results[:i], nil
}

type Tokenized struct {
	Conn *pgx.Conn
}

func (s *Tokenized) Description() string {
	return "% operator query on (food_id, token) table with join"
}

func (s *Tokenized) Query(ctx context.Context, query string) ([]string, error) {
	inner := `SELECT food_id, COUNT(*) FROM food_tokens WHERE`

	stemmed := Stemmed(query)
	for i, stem := range stemmed {
		if i > 0 {
			inner += ` OR`
		}
		inner += fmt.Sprintf(` token %% '%s'`, stem)
	}

	inner += ` GROUP BY food_id HAVING COUNT(*) > 1`
	sql := fmt.Sprintf(`
		SELECT foods.name
		FROM foods JOIN (%s)
		results ON results.food_id = foods.id
		ORDER BY results.count DESC
		LIMIT 10`,
		inner,
	)

	rows, err := s.Conn.Query(ctx, sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	i := 0
	results := make([]string, 10)
	for ; rows.Next(); i++ {
		if err := rows.Scan(&results[i]); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return results[:i], nil
}

type ILike struct {
	Conn *pgx.Conn
}

func (s *ILike) Description() string {
	return "ILIKE on analyzed ordered by similarity"
}

func (s *ILike) Query(ctx context.Context, query string) ([]string, error) {
	sql := `SELECT name FROM foods WHERE`
	stemmed := Stemmed(query)
	for i, stem := range stemmed {
		if i > 0 {
			sql += ` AND`
		}
		sql += fmt.Sprintf(` analyzed ILIKE '%%%s%%'`, stem)
	}
	sql += ` ORDER BY similarity(analyzed, $1) DESC LIMIT 10`

	rows, err := s.Conn.Query(ctx, sql, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	i := 0
	results := make([]string, 10)
	for ; rows.Next(); i++ {
		if err := rows.Scan(&results[i]); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return results[:i], nil
}
