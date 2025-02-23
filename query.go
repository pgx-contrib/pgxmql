package pgxfilter

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/mql"
	"github.com/jackc/pgx/v5"
)

var _ pgx.QueryRewriter = &QueryRewriter{}

// QueryRewriter is a pgx.QueryRewriter that rewrites queries to include a filter.
type QueryRewriter struct {
	expr string
	typ  any
}

// New creates a new QueryRewriter.
func New[T any](expr string) *QueryRewriter {
	var typ T

	return &QueryRewriter{
		expr: expr,
		typ:  typ,
	}
}

var (
	condition = regexp.MustCompile(`(?im)--(.*?)\s*:condition\s*(.*)`)
	exclusion = regexp.MustCompile(`(?im)^\s*--\s*exclude\s*\n[\s\S]*?\n\s*--\s*exclude\s*\n?`)
)

// RewriteQuery implements pgx.QueryRewriter.
func (x *QueryRewriter) RewriteQuery(ctx context.Context, _ *pgx.Conn, query string, args []any) (string, []any, error) {
	query = x.exclude(query)
	// Find the first match only
	match := condition.FindStringSubmatchIndex(query)
	if match == nil {
		return query, args, nil
	}

	// parse the filter expression
	where, err := mql.Parse(x.expr, x.typ, x.options()...)
	if err != nil {
		return "", nil, err
	}

	// create a PG argument
	param := func(i int) string {
		return fmt.Sprintf("$%d", i)
	}

	offet := len(args)
	// the original arguments should be moved by the offset
	for i := 1; i <= len(where.Args); i++ {
		prev := i
		next := i + offet

		where.Condition = strings.ReplaceAll(where.Condition,
			param(prev),
			param(next),
		)
	}

	// Apply the replacement with correct placement
	query = x.replace(query, match, where.Condition)
	// the remaining args are the placeholders
	args = append(args, where.Args...)
	// done!
	return query, args, nil
}

func (x *QueryRewriter) options() []mql.Option {
	include := make(map[string]string)
	ignore := []string{}

	// Get the type of the struct
	t := reflect.TypeOf(x.typ)

	// Ensure it's a struct
	if t.Kind() == reflect.Struct {
		// Iterate over struct fields
		for i := range t.NumField() {
			field := t.Field(i)

			if tag := field.Tag.Get("db"); tag != "" && tag != "-" {
				include[field.Name] = tag
			} else {
				ignore = append(ignore, field.Name)
			}
		}
	}

	return []mql.Option{
		mql.WithPgPlaceholders(),
		mql.WithIgnoredFields(ignore...),
	}
}

func (x *QueryRewriter) replace(query string, match []int, condition string) string {
	// Format the argument placeholders
	condition = strings.ReplaceAll(condition, "=$", " = $")

	start := match[0]
	end := match[1]

	// Extract prefix and suffix (if present)
	var prefix, suffix string

	if len(match) > 2 && match[2] != -1 {
		prefix = strings.TrimSpace(query[match[2]:match[3]])
	}
	if len(match) > 4 && match[4] != -1 {
		suffix = strings.TrimSpace(query[match[4]:match[5]])
	}

	// Trim condition to remove extra spaces
	condition = strings.TrimSpace(condition)

	// Prepend prefix if available
	if prefix != "" {
		condition = prefix + " " + condition
	}
	// Append suffix if available
	if suffix != "" {
		condition = condition + " " + suffix
	}

	prefix = query[:start]
	suffix = query[end:]

	// Replace the matched section with the formatted condition
	return x.trim(prefix) + x.trim(condition) + x.trim(suffix)
}

func (x *QueryRewriter) trim(query string) string {
	return strings.ReplaceAll(query, ";", "")
}

func (x *QueryRewriter) exclude(query string) string {
	return exclusion.ReplaceAllString(query, "")
}
