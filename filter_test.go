package pgxfilter_test

import (
	"github.com/jackc/pgx/v5"
	"github.com/pgx-contrib/pgxfilter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/pgx-contrib/pgxfilter/fake"
)

var _ = Describe("QueryRewriter", func() {
	var rewriter *pgxfilter.QueryRewriter

	type User struct {
		ID    int    `db:"id"`
		Perm  int    `db:"perm"`
		Name  string `db:"name"`
		Role  string `db:"role"`
		Group string `db:"-"`
	}

	BeforeEach(func() {
		rewriter = pgxfilter.New[User]("role = 'admin'")
	})

	Describe("RewriteQuery", func() {
		It("should rewrite the query with the filter", func(ctx SpecContext) {
			query := NewFakeQuery("000.sql")

			querySQL, queryArgs, err := rewriter.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
			Expect(err).NotTo(HaveOccurred())
			Expect(querySQL).To(ContainSubstring("role = $1"))
			Expect(queryArgs).To(HaveLen(1))
			Expect(queryArgs).To(ContainElement("admin"))
		})

		When("the column is not found", func() {
			It("returns an error", func(ctx SpecContext) {
				query := NewFakeQuery("000.sql")

				rewriter = pgxfilter.New[User]("company = 'IBM'")
				querySQL, queryArgs, err := rewriter.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).To(HaveOccurred())
				Expect(querySQL).To(BeEmpty())
				Expect(queryArgs).To(BeEmpty())
			})
		})

		When("the filter is at the beginning", func() {
			It("should rewrite the query with the filter", func(ctx SpecContext) {
				query := NewFakeQuery("001.sql", 0, "root", 0)

				querySQL, queryArgs, err := rewriter.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).NotTo(HaveOccurred())
				Expect(querySQL).To(ContainSubstring("role = $4"))
				Expect(queryArgs).To(HaveLen(4))
				Expect(queryArgs).To(ContainElement("admin"))
			})
		})

		When("the filter is at the middle", func() {
			It("should rewrite the query with the filter", func(ctx SpecContext) {
				query := NewFakeQuery("010.sql", 0, "root", 0)

				querySQL, queryArgs, err := rewriter.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).NotTo(HaveOccurred())
				Expect(querySQL).To(ContainSubstring("role = $4"))
				Expect(queryArgs).To(HaveLen(4))
				Expect(queryArgs).To(ContainElement("admin"))
			})
		})

		When("the filter is at the end", func() {
			It("should rewrite the query with the filter", func(ctx SpecContext) {
				query := NewFakeQuery("100.sql", 0, "root", 0)

				querySQL, queryArgs, err := rewriter.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).NotTo(HaveOccurred())
				Expect(querySQL).To(ContainSubstring("role = $4"))
				Expect(queryArgs).To(HaveLen(4))
				Expect(queryArgs).To(ContainElement("admin"))
			})
		})

		When("there is no match", func() {
			It("should return the original query and args", func(ctx SpecContext) {
				query := &pgx.QueuedQuery{
					SQL: "SELECT * FROM users",
				}
				querySQL, queryArgs, err := rewriter.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).NotTo(HaveOccurred())
				Expect(querySQL).To(Equal(query.SQL))
				Expect(queryArgs).To(BeEmpty())
			})
		})
	})
})
