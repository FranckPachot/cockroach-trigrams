package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/chrisseto/cockroach-trigrams/trgm"
	"github.com/jackc/pgx/v5"
)

func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
}

func MustT[T any](rv T, err error) T {
	Must(err)
	return rv
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conn := MustT(pgx.Connect(ctx, "postgres://root@localhost:26257/defaultdb"))
	defer conn.Close(context.Background())

	switch os.Args[1] {
	case "load":
		Must(Load(ctx, conn))
	case "query":
		Must(Query(ctx, conn, os.Args[2]))
	default:
		fmt.Printf("Unknown command %q\n", os.Args[1])
		os.Exit(1)
	}

}

func Query(ctx context.Context, conn *pgx.Conn, query string) error {
	queriers := []trgm.Querier{
		&trgm.Naive{Conn: conn},
		&trgm.Naive{Conn: conn, AnalyzedQuery: true},
		&trgm.Naive{Conn: conn, AnalyzedQuery: true, AnalyzedField: true},
	}

	for i, querier := range queriers {
		if i > 0 {
			fmt.Print("\n\n")
		}

		start := time.Now()
		results, err := querier.Query(ctx, query)
		if err != nil {
			return err
		}
		elapsed := time.Since(start)

		fmt.Printf("Searched %q via %q (%s):\n", query, querier.Description(), elapsed)

		if len(results) == 0 {
			fmt.Printf("\tNo Results Found\n")
			continue
		}

		for j, result := range results {
			fmt.Printf("\t%d: %q\n", j+1, result)
		}
	}

	return nil
}

func Load(ctx context.Context, conn *pgx.Conn) error {
	indexer := &trgm.Indexer{Conn: conn}

	if err := indexer.Setup(ctx); err != nil {
		return err
	}

	for _, path := range []string{"data/foundation-foods.json", "data/branded-foods.json"} {
		f := MustT(os.Open(path))
		defer f.Close()

		if err := LoadChunks(f, 250, func(chunk []*trgm.Food) error {
			if err := indexer.LoadChunk(ctx, chunk); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if err := indexer.PostLoad(ctx); err != nil {
		return err
	}

	return nil
}

func LoadChunks(src io.Reader, chunkSize int, fn func([]*trgm.Food) error) error {
	decoder := json.NewDecoder(src)

	for {
		chunk := make([]*trgm.Food, chunkSize)
		for i := 0; i < chunkSize; i++ {
			var food trgm.Food
			if err := decoder.Decode(&food); err != nil {
				if errors.Is(err, io.EOF) {
					chunk = chunk[:i]
					break
				}
				return err
			}
			chunk[i] = &food
		}

		if err := fn(chunk); err != nil {
			return err
		}

		if len(chunk) < chunkSize {
			return nil
		}
	}
}
