package pgxfilter_test

import (
	"github.com/pgx-contrib/pgxfilter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/pgx-contrib/pgxfilter/fake"
)

var _ = Describe("WhereClause", func() {
	type User struct {
		ID       string `db:"id"`
		Role     string `db:"role"`
		Password string `db:"password"`
		Company  string `db:"-"`
	}

	var filter *pgxfilter.WhereClause

	BeforeEach(func() {
		filter = &pgxfilter.WhereClause{
			Condition: "role = 'admin'",
			Model:     &User{},
		}
	})

	Describe("RewriteQuery", func() {
		It("rewrites the query successfully", func(ctx SpecContext) {
			query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
			querySQL, queryArgs, err := filter.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
			Expect(err).NotTo(HaveOccurred())
			Expect(querySQL).To(ContainSubstring("role=$5"))
			Expect(queryArgs).To(HaveLen(5))
			Expect(queryArgs).To(ContainElement("admin"))
		})

		When("the expression is invalid", func() {
			BeforeEach(func() {
				filter.Condition = "first_name = 'John'"
			})

			It("returns an error", func(ctx SpecContext) {
				query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
				querySQL, queryArgs, err := filter.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).To(MatchError(ContainSubstring(`invalid column "first_name"`)))
				Expect(querySQL).To(BeEmpty())
				Expect(queryArgs).To(BeEmpty())
			})
		})
	})
})
