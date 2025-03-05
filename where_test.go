package pgxmql_test

import (
	"github.com/pgx-contrib/pgxmql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/pgx-contrib/pgxmql/fake"
)

var _ = Describe("WhereClause", func() {
	type User struct {
		ID       string `db:"id_db"`
		Role     string `db:"role_db"`
		Password string `db:"-"`
		Company  string `db:"company_db" json:"company_json"`
	}

	var clause *pgxmql.WhereClause

	BeforeEach(func() {
		clause = &pgxmql.WhereClause{
			Condition: "role = 'root' and company_json % 'TSLA'",
			Model:     &User{},
		}
	})

	Describe("RewriteQuery", func() {
		It("rewrites the query successfully", func(ctx SpecContext) {
			query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
			querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)

			Expect(err).NotTo(HaveOccurred())
			Expect(querySQL).To(ContainSubstring("role_db=$5"))
			Expect(querySQL).To(ContainSubstring("company_db like $6"))
			Expect(queryArgs).To(HaveLen(6))
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
