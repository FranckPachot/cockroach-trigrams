package trgm

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type Indexer struct {
	Conn *pgx.Conn

	tokenMap map[string]uuid.UUID
}

func (i *Indexer) Search(ctx context.Context, query string) ([]string, error) {
	rows, err := i.Conn.Query(ctx, `SELECT name FROM raw_foods WHERE name % $1 ORDER BY similarity(name, $1) DESC LIMIT 10`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	j := 0
	results := make([]string, 10)
	for ; rows.Next(); j++ {
		if err := rows.Scan(&results[j]); err != nil {
			return nil, err
		}
	}
	return results[:j], nil
}

func (i *Indexer) Setup(ctx context.Context) error {
	_, err := i.Conn.Exec(ctx, `
DROP TABLE IF EXISTS foods;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS food_tokens;
DROP TABLE IF EXISTS food_to_token;

create extension if not exists pgcrypto;
create extension if not exists pg_trgm;

CREATE TABLE IF NOT EXISTS foods(
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name text NOT NULL,
	analyzed text NOT NULL,
	weight float8 NOT NULL
);

CREATE TABLE IF NOT EXISTS food_tokens(
	food_id uuid NOT NULL,
	token text NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens(
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	token text NOT NULL
);

CREATE TABLE IF NOT EXISTS food_to_token(
	token_id uuid NOT NULL,
	food_id uuid NOT NULL
);
	`)
	return errors.WithStack(err)
}

func (i *Indexer) PostLoad(ctx context.Context) error {
	j := 0
	rows := make([][]interface{}, len(i.tokenMap))
	for tok, id := range i.tokenMap {
		rows[j] = []interface{}{id, tok}
		j++
	}

	_, err := i.Conn.CopyFrom(ctx, pgx.Identifier{"tokens"}, []string{"id", "token"}, pgx.CopyFromRows(rows))
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = i.Conn.Exec(ctx, `CREATE INDEX ON foods USING GIN (name gin_trgm_ops);`)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = i.Conn.Exec(ctx, `CREATE INDEX ON food_tokens USING GIN (token gin_trgm_ops);`)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = i.Conn.Exec(ctx, `CREATE INDEX ON tokens USING GIN (token gin_trgm_ops);`)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = i.Conn.Exec(ctx, `CREATE INDEX ON foods USING GIN (analyzed gin_trgm_ops);`)
	return errors.WithStack(err)
}

func (i *Indexer) LoadChunk(ctx context.Context, chunk []*Food) error {
	if err := i.loadFoods(ctx, chunk); err != nil {
		return err
	}
	if err := i.loadFoodTokens(ctx, chunk); err != nil {
		return err
	}
	return i.loadTokens(ctx, chunk)
}

func (i *Indexer) loadFoods(ctx context.Context, chunk []*Food) error {
	rows := make([][]interface{}, len(chunk))
	for i, food := range chunk {
		rows[i] = []interface{}{
			food.ID(),
			food.FullName(),
			food.AnalyzedFullName(),
			food.Weight(),
		}
	}

	_, err := i.Conn.CopyFrom(ctx, pgx.Identifier{"foods"}, []string{"id", "name", "analyzed", "weight"}, pgx.CopyFromRows(rows))
	return errors.WithStack(err)
}

func (i *Indexer) loadFoodTokens(ctx context.Context, chunk []*Food) error {
	rows := make([][]interface{}, 0)

	for _, food := range chunk {
		tokens := make(map[string]struct{}, 15)
		for _, tok := range food.StemmedFullName() {
			tokens[tok] = struct{}{}
		}

		i := 0
		rs := make([][]interface{}, len(tokens))
		for tok := range tokens {
			rs[i] = []interface{}{food.ID(), tok}
			i++
		}

		rows = append(rows, rs...)
	}

	_, err := i.Conn.CopyFrom(ctx, pgx.Identifier{"food_tokens"}, []string{"food_id", "token"}, pgx.CopyFromRows(rows))
	return errors.WithStack(err)
}

func (i *Indexer) loadTokens(ctx context.Context, chunk []*Food) error {
	if i.tokenMap == nil {
		i.tokenMap = make(map[string]uuid.UUID)
	}

	rows := make([][]interface{}, 0)

	for _, food := range chunk {
		tokens := make(map[string]struct{}, 15)
		for _, tok := range food.StemmedFullName() {
			tokens[tok] = struct{}{}
		}

		j := 0
		rs := make([][]interface{}, len(tokens))
		for tok := range tokens {
			if _, ok := i.tokenMap[tok]; !ok {
				i.tokenMap[tok] = uuid.Must(uuid.NewV4())
			}
			rs[j] = []interface{}{i.tokenMap[tok], food.ID()}
			j++
		}

		rows = append(rows, rs...)
	}

	_, err := i.Conn.CopyFrom(ctx, pgx.Identifier{"food_to_token"}, []string{"token_id", "food_id"}, pgx.CopyFromRows(rows))
	return errors.WithStack(err)
}
