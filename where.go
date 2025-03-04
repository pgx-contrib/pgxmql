package pgxfilter

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/mql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ pgx.QueryRewriter = &WhereClause{}

// WhereClause is a pgx.QueryRewriter that rewrites queries to include a filter.
type WhereClause struct {
	// Condition is the where clause condition
	Condition string
	// Model is the database model
	Model any
}

// RewriteQuery implements pgx.QueryRewriter.
func (x *WhereClause) RewriteQuery(ctx context.Context, _ *pgx.Conn, query string, args []any) (_ string, _ []any, err error) {
	// prepare the query
	query = x.condition(query)
	// if we don't have a placeholder, we need to decrement the query parameters
	if !strings.Contains(query, "$1") {
		query = x.parameters(query, -1)
	}

	if x.Condition != "" {
		// parse the filter expression
		clause, err := mql.Parse(x.Condition, x.Model, x.ignore(), x.placeholder())
		if err != nil {
			cause := fmt.Errorf("cause: %w query: %v filter: %v", err, query, x.Condition)

			return "", nil, &pgconn.PgError{
				Severity:      "ERROR",
				Code:          "42601",
				Where:         x.Condition,
				Message:       cause.Error(),
				InternalQuery: query,
			}
		}

		if strings.Contains(query, "$1") {
			clause.Condition = x.parameters(clause.Condition, len(args))
		}
		// append the filter to the query
		args = append(args, clause.Args...)
		// inject the clause into the query
		query = x.inject(query, clause.Condition)
	}

	return query, args, nil
}

func (x *WhereClause) condition(query string) string {
	return strings.Replace(query, "$1::void IS NULL", "TRUE", 1)
}

func (x *WhereClause) parameters(query string, value int) string {
	// Regular expression to match $ followed by a number
	re := regexp.MustCompile(`\$(\d+)`)

	// Replace function to decrement the numbers
	return re.ReplaceAllStringFunc(query, func(match string) string {
		// Extract the number from the match
		num, _ := strconv.Atoi(match[1:])
		// Decrement and return the new parameter
		return fmt.Sprintf("$%d", num+value)
	})
}

func (x *WhereClause) inject(query string, condition string) string {
	return strings.Replace(query, "-- :condition", condition, 1)
}

func (x *WhereClause) ignore() mql.Option {
	include := make(map[string]string)
	columns := []string{}

	// Get the type of the struct
	t := reflect.TypeOf(x.Model)
	// obtain the underlying type if it's a pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Ensure it's a struct
	if t.Kind() == reflect.Struct {
		// Iterate over struct fields
		for i := range t.NumField() {
			field := t.Field(i)

			if tag := field.Tag.Get("db"); tag != "" && tag != "-" {
				include[field.Name] = tag
			} else {
				columns = append(columns, field.Name)
			}
		}
	}

	return mql.WithIgnoredFields(columns...)
}

func (x *WhereClause) placeholder() mql.Option {
	return mql.WithPgPlaceholders()
}
