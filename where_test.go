package pgxmql_test

import (
	"github.com/pgx-contrib/pgxmql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/pgx-contrib/pgxmql/fake"
)

var _ = Describe("WhereClause", func() {
	type User struct {
		ID       string `db:"id" json:"id"`
		Role     string `db:"role" json:"role"`
		Password string `db:"password" json:"password"`
		Company  string `db:"-" json:"company"`
	}

	var clause *pgxmql.WhereClause

	BeforeEach(func() {
		clause = &pgxmql.WhereClause{
			Condition: "role = 'admin'",
			Model:     &User{},
		}
	})

	Describe("RewriteQuery", func() {
		It("rewrites the query successfully", func(ctx SpecContext) {
			query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
			querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
			Expect(err).NotTo(HaveOccurred())
			Expect(querySQL).To(ContainSubstring("role=$5"))
			Expect(queryArgs).To(HaveLen(5))
			Expect(queryArgs).To(ContainElement("admin"))
		})

		When("the expression is empty", func() {
			BeforeEach(func() {
				clause.Condition = ""
			})

			It("rewrites the query successfully", func(ctx SpecContext) {
				query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
				querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).NotTo(HaveOccurred())
				Expect(querySQL).To(ContainSubstring("TRUE"))
				Expect(queryArgs).To(HaveLen(4))
			})
		})

		When("the expression is invalid", func() {
			BeforeEach(func() {
				clause.Condition = "first_name = 'John'"
			})

			It("returns an error", func(ctx SpecContext) {
				query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
				querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).To(MatchError(ContainSubstring(`invalid column "first_name"`)))
				Expect(querySQL).To(BeEmpty())
				Expect(queryArgs).To(BeEmpty())
			})
		})
	})
})
