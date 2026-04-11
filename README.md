# pgxmql

[![CI](https://github.com/pgx-contrib/pgxmql/actions/workflows/ci.yml/badge.svg)](https://github.com/pgx-contrib/pgxmql/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/pgx-contrib/pgxmql?include_prereleases)](https://github.com/pgx-contrib/pgxmql/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/pgx-contrib/pgxmql.svg)](https://pkg.go.dev/github.com/pgx-contrib/pgxmql)
[![License](https://img.shields.io/github/license/pgx-contrib/pgxmql)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![pgx](https://img.shields.io/badge/pgx-v5-blue)](https://github.com/jackc/pgx)
[![hashicorp/mql](https://img.shields.io/badge/hashicorp%2Fmql-enabled-blueviolet)](https://github.com/hashicorp/mql)

A [pgx](https://github.com/jackc/pgx) `QueryRewriter` adapter for
[hashicorp/mql](https://github.com/hashicorp/mql) dynamic filtering. Write
your SQL with a void placeholder and let `WhereClause` inject a safe,
parameterised filter at query time.

## Installation

```bash
go get github.com/pgx-contrib/pgxmql
```

## Usage

Define your model with `db` struct tags to map fields to table columns. Use
`$1::void IS NULL` as the filter placeholder in your SQL:

```go
type User struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
    Role string `db:"role"`
}

rows, err := pool.Query(ctx,
    "SELECT * FROM users WHERE $1::void IS NULL",
    &pgxmql.WhereClause{
        Condition: "role = 'admin'",
        Model:     User{},
    },
)
```

The placeholder is replaced by the parsed MQL condition before the query
reaches PostgreSQL. When `Condition` is empty the placeholder becomes `TRUE`,
returning all rows.

## Development

### DevContainer

Open in VS Code with the Dev Containers extension. The environment provides Go,
PostgreSQL 18, and Nix automatically.

```
PGX_DATABASE_URL=postgres://vscode@postgres:5432/pgxmql?sslmode=disable
```

### Nix

```bash
nix develop          # enter shell with Go
go tool ginkgo run -r
```

### Run tests

```bash
# Unit tests only (no database required)
go tool ginkgo run -r

# With integration tests
export PGX_DATABASE_URL="postgres://localhost/pgxmql?sslmode=disable"
go tool ginkgo run -r
```

## License

[MIT](LICENSE)
