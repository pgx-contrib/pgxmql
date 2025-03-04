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
	query = x.replaceVoid(query)
	// if we don't have a placeholder, we need to decrement the query parameters
	if !strings.Contains(query, "$1") {
		query = x.replaceArgs(query, -1)
	}

	if x.Condition != "" {
		// parse the filter expression
		clause, err := mql.Parse(x.Condition, x.Model, x.options()...)
		if err != nil {
			return "", nil, &pgconn.PgError{
				Severity:      "ERROR",
				Code:          "42601",
				Where:         x.Condition,
				Message:       err.Error(),
				InternalQuery: query,
			}
		}

		if strings.Contains(query, "$1") {
			clause.Condition = x.replaceArgs(clause.Condition, len(args))
		}
		// append the filter to the query
		args = append(args, clause.Args...)
		// inject the clause into the query
		query = x.replaceCond(query, clause.Condition)
	} else {
		query = x.replaceCond(query, "TRUE")
	}

	// done!
	return query, args, nil
}

func (x *WhereClause) options() []mql.Option {
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
			if value := field.Tag.Get("json"); value != "" && value != "-" {
				field.Name = value
			}

			if value := field.Tag.Get("db"); value != "" && value != "-" {
				include[field.Name] = value
			} else {
				columns = append(columns, field.Name)
			}
		}
	}

	return []mql.Option{
		mql.WithPgPlaceholders(),
		mql.WithTableColumnMap(include),
		mql.WithIgnoredFields(columns...),
	}
}

func (x *WhereClause) replaceVoid(query string) string {
	return strings.Replace(query, "$1::void IS NULL", "-- :condition", 1)
}

func (x *WhereClause) replaceArgs(query string, delta int) string {
	// Regular expression to match $ followed by a number
	re := regexp.MustCompile(`\$(\d+)`)

	// Replace function to decrement the numbers
	return re.ReplaceAllStringFunc(query, func(match string) string {
		// Extract the number from the match
		position, _ := strconv.Atoi(match[1:])
		// next value
		value := position + delta
		// Decrement and return the new parameter
		return fmt.Sprintf("$%d", value)
	})
}

func (x *WhereClause) replaceCond(query string, condition string) string {
	return strings.Replace(query, "-- :condition", condition, 1)
}
