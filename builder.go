package pgxmql

import (
	"fmt"
	"strings"
)

// Where creates a WhereClause with the given MQL condition for model type T.
func Where[T any](condition string) *WhereClause {
	var model T
	return &WhereClause{
		Condition: condition,
		Model:     model,
	}
}

// In creates a WhereClause that matches any of the given values for the
// specified column, equivalent to column = "v1" or column = "v2" or ...
// Values are escaped to prevent MQL injection. An empty values slice produces
// a WhereClause with an empty condition.
func In[T any](column string, values ...string) *WhereClause {
	if len(values) == 0 {
		return Where[T]("")
	}

	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts,
			fmt.Sprintf("%s = \"%s\"", column, escapeValue(v)),
		)
	}

	return Where[T](strings.Join(parts, " or "))
}

// And combines the receiver with one or more additional WhereClause values
// using the MQL "and" logical operator. Nil clauses and clauses with empty
// conditions are silently skipped. Returns a new WhereClause; the receiver
// is not mutated.
func (x *WhereClause) And(clauses ...*WhereClause) *WhereClause {
	return x.combine("and", clauses)
}

// Or combines the receiver with one or more additional WhereClause values
// using the MQL "or" logical operator. Nil clauses and clauses with empty
// conditions are silently skipped. Returns a new WhereClause; the receiver
// is not mutated.
func (x *WhereClause) Or(clauses ...*WhereClause) *WhereClause {
	return x.combine("or", clauses)
}

func (x *WhereClause) combine(op string, clauses []*WhereClause) *WhereClause {
	parts := make([]string, 0, 1+len(clauses))

	if cond := strings.TrimSpace(x.Condition); cond != "" {
		parts = append(parts, "("+cond+")")
	}

	for _, c := range clauses {
		if c != nil {
			if cond := strings.TrimSpace(c.Condition); cond != "" {
				parts = append(parts, "("+cond+")")
			}
		}
	}

	var condition string

	switch len(parts) {
	case 0:
		// all empty
	case 1:
		// unwrap the parentheses we just added
		condition = parts[0][1 : len(parts[0])-1]
	default:
		condition = strings.Join(parts, " "+op+" ")
	}

	return &WhereClause{
		Condition: condition,
		Model:     x.Model,
	}
}

func escapeValue(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
