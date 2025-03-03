package pgxfilter_test

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgx-contrib/pgxfilter"
)

func ExampleWhereClause() {
	type User struct {
		ID    int    `db:"id"`
		Perm  int    `db:"perm"`
		Name  string `db:"name"`
		Role  string `db:"role"`
		Group string `db:"-"`
	}

	config, err := pgxpool.ParseConfig(os.Getenv("PGX_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()
	// Create a new pgxpool with the config
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}
	// close the pool
	defer pool.Close()

	rows, err := pool.Query(ctx,
		"SELECT * from user WHERE $1:void IS NULL",
		&pgxfilter.WhereClause{
			Condition: "role = 'admin'",
			Model:     User{},
		},
	)
	if err != nil {
		panic(err)
	}
	// close the rows
	defer rows.Close()

	for rows.Next() {
		entity, err := pgx.RowToStructByName[User](rows)
		if err != nil {
			panic(err)
		}

		fmt.Println(entity.Name)
	}
}
