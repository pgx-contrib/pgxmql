package pgxmql_test

import (
	"github.com/pgx-contrib/pgxmql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/pgx-contrib/pgxmql/fake"
)

var _ = Describe("Builder", func() {
	type User struct {
		ID      string `db:"id_db"`
		Role    string `db:"role_db"`
		Company string `db:"company_db" json:"company_json"`
	}

	Describe("Where", func() {
		It("creates a WhereClause with condition and model", func(ctx SpecContext) {
			clause := pgxmql.Where[User]("role = 'root'")
			Expect(clause.Condition).To(Equal("role = 'root'"))

			query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
			querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
			Expect(err).NotTo(HaveOccurred())
			Expect(querySQL).To(ContainSubstring("role_db=$5"))
			Expect(queryArgs).To(HaveLen(5))
		})

		It("creates a WhereClause with empty condition", func(ctx SpecContext) {
			clause := pgxmql.Where[User]("")
			Expect(clause.Condition).To(BeEmpty())

			query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
			querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
			Expect(err).NotTo(HaveOccurred())
			Expect(querySQL).To(ContainSubstring("TRUE"))
			Expect(queryArgs).To(HaveLen(4))
		})
	})

	Describe("And", func() {
		It("combines two conditions with AND", func() {
			clause := pgxmql.Where[User]("role = 'admin'").
				And(pgxmql.Where[User]("company_json = 'ACME'"))

			Expect(clause.Condition).To(Equal("(role = 'admin') and (company_json = 'ACME')"))
		})

		It("combines three conditions with AND", func() {
			clause := pgxmql.Where[User]("role = 'admin'").
				And(
					pgxmql.Where[User]("company_json = 'ACME'"),
					pgxmql.Where[User]("id = '007'"),
				)

			Expect(clause.Condition).To(Equal("(role = 'admin') and (company_json = 'ACME') and (id = '007')"))
		})

		It("skips nil clauses", func() {
			clause := pgxmql.Where[User]("role = 'admin'").And(nil)
			Expect(clause.Condition).To(Equal("role = 'admin'"))
		})

		It("skips clauses with empty conditions", func() {
			clause := pgxmql.Where[User]("role = 'admin'").
				And(pgxmql.Where[User](""))

			Expect(clause.Condition).To(Equal("role = 'admin'"))
		})

		It("returns empty condition when all conditions are empty", func() {
			clause := pgxmql.Where[User]("").And(pgxmql.Where[User](""))
			Expect(clause.Condition).To(BeEmpty())
		})

		It("returns receiver condition when no clauses provided", func() {
			clause := pgxmql.Where[User]("role = 'admin'").And()
			Expect(clause.Condition).To(Equal("role = 'admin'"))
		})

		When("combined with RewriteQuery", func() {
			It("produces correct parameterized SQL", func(ctx SpecContext) {
				clause := pgxmql.Where[User]("role = 'root'").
					And(pgxmql.Where[User]("company_json = 'TSLA'"))

				query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
				querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).NotTo(HaveOccurred())
				Expect(querySQL).To(ContainSubstring("role_db=$5"))
				Expect(querySQL).To(ContainSubstring("company_db=$6"))
				Expect(queryArgs).To(HaveLen(6))
			})
		})
	})

	Describe("Or", func() {
		It("combines two conditions with OR", func() {
			clause := pgxmql.Where[User]("role = 'admin'").
				Or(pgxmql.Where[User]("role = 'root'"))

			Expect(clause.Condition).To(Equal("(role = 'admin') or (role = 'root')"))
		})

		It("combines three conditions with OR", func() {
			clause := pgxmql.Where[User]("role = 'admin'").
				Or(
					pgxmql.Where[User]("role = 'root'"),
					pgxmql.Where[User]("role = 'superuser'"),
				)

			Expect(clause.Condition).To(Equal("(role = 'admin') or (role = 'root') or (role = 'superuser')"))
		})

		It("skips nil clauses", func() {
			clause := pgxmql.Where[User]("role = 'admin'").Or(nil)
			Expect(clause.Condition).To(Equal("role = 'admin'"))
		})

		It("skips clauses with empty conditions", func() {
			clause := pgxmql.Where[User]("role = 'admin'").
				Or(pgxmql.Where[User](""))

			Expect(clause.Condition).To(Equal("role = 'admin'"))
		})

		It("returns empty condition when all conditions are empty", func() {
			clause := pgxmql.Where[User]("").Or(pgxmql.Where[User](""))
			Expect(clause.Condition).To(BeEmpty())
		})

		It("returns receiver condition when no clauses provided", func() {
			clause := pgxmql.Where[User]("role = 'admin'").Or()
			Expect(clause.Condition).To(Equal("role = 'admin'"))
		})

		When("combined with RewriteQuery", func() {
			It("produces correct parameterized SQL", func(ctx SpecContext) {
				clause := pgxmql.Where[User]("role = 'root'").
					Or(pgxmql.Where[User]("role = 'admin'"))

				query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
				querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).NotTo(HaveOccurred())
				Expect(querySQL).To(ContainSubstring("role_db=$5"))
				Expect(querySQL).To(ContainSubstring("role_db=$6"))
				Expect(queryArgs).To(HaveLen(6))
			})
		})
	})

	Describe("In", func() {
		It("creates OR condition for multiple values", func() {
			clause := pgxmql.In[User]("role", "admin", "root", "superuser")
			Expect(clause.Condition).To(Equal(`role = "admin" or role = "root" or role = "superuser"`))
		})

		It("creates single condition for one value", func() {
			clause := pgxmql.In[User]("role", "admin")
			Expect(clause.Condition).To(Equal(`role = "admin"`))
		})

		It("returns empty condition for no values", func() {
			clause := pgxmql.In[User]("role")
			Expect(clause.Condition).To(BeEmpty())
		})

		It("escapes double quotes in values", func() {
			clause := pgxmql.In[User]("role", `say "hello"`)
			Expect(clause.Condition).To(Equal(`role = "say \"hello\""`))
		})

		It("escapes backslashes in values", func() {
			clause := pgxmql.In[User]("role", `back\slash`)
			Expect(clause.Condition).To(Equal(`role = "back\\slash"`))
		})

		When("combined with RewriteQuery", func() {
			It("produces correct parameterized SQL", func(ctx SpecContext) {
				clause := pgxmql.In[User]("role", "root", "admin")

				query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
				querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
				Expect(err).NotTo(HaveOccurred())
				Expect(querySQL).To(ContainSubstring("role_db=$5"))
				Expect(querySQL).To(ContainSubstring("role_db=$6"))
				Expect(queryArgs).To(HaveLen(6))
			})
		})
	})

	Describe("Chaining", func() {
		It("supports And then Or", func() {
			clause := pgxmql.Where[User]("role = 'admin'").
				And(pgxmql.Where[User]("company_json = 'ACME'")).
				Or(pgxmql.Where[User]("role = 'root'"))

			Expect(clause.Condition).To(Equal(
				"((role = 'admin') and (company_json = 'ACME')) or (role = 'root')",
			))
		})

		It("supports Or then And", func() {
			clause := pgxmql.Where[User]("role = 'admin'").
				Or(pgxmql.Where[User]("role = 'root'")).
				And(pgxmql.Where[User]("company_json = 'ACME'"))

			Expect(clause.Condition).To(Equal(
				"((role = 'admin') or (role = 'root')) and (company_json = 'ACME')",
			))
		})

		It("supports In combined with Where via And", func(ctx SpecContext) {
			clause := pgxmql.In[User]("role", "admin", "root").
				And(pgxmql.Where[User]("company_json = 'ACME'"))

			Expect(clause.Condition).To(Equal(
				`(role = "admin" or role = "root") and (company_json = 'ACME')`,
			))

			query := NewFakeQuery("001.sql", "007", "Google", nil, nil)
			querySQL, queryArgs, err := clause.RewriteQuery(ctx, nil, query.SQL, query.Arguments)
			Expect(err).NotTo(HaveOccurred())
			Expect(querySQL).To(ContainSubstring("role_db=$5"))
			Expect(querySQL).To(ContainSubstring("role_db=$6"))
			Expect(querySQL).To(ContainSubstring("company_db=$7"))
			Expect(queryArgs).To(HaveLen(7))
		})

		It("does not mutate the receiver", func() {
			original := pgxmql.Where[User]("role = 'admin'")
			_ = original.And(pgxmql.Where[User]("company_json = 'ACME'"))
			Expect(original.Condition).To(Equal("role = 'admin'"))
		})
	})
})
