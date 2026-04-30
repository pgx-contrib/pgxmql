# Changelog

## 1.0.0 (2026-04-30)


### Features

* add pgxfilter module with query rewriter ([10e44ea](https://github.com/pgx-contrib/pgxmql/commit/10e44ea35a0ea612f887cc0fcf4ed569bc40bab5))
* add WhereClause builder API with Where, In, And, Or ([72a2e21](https://github.com/pgx-contrib/pgxmql/commit/72a2e21b214dbf1148fc3865abf8ac23ac2aea27))
* enhance WhereClause with table and model mappings ([a8d5eed](https://github.com/pgx-contrib/pgxmql/commit/a8d5eed6b7869f5a22dd8c5172ce77246efd35a6))
* map JSON tags to field names in WhereClause ([b90cb07](https://github.com/pgx-contrib/pgxmql/commit/b90cb073069a1dba31cec76e0d998458695cccc1))
* refactor query rewriting methods ([cb61449](https://github.com/pgx-contrib/pgxmql/commit/cb61449e16579c8d6fdc56e38ae67c89ad4edd10))
* rename pgxfilter to pgxmql ([c68683c](https://github.com/pgx-contrib/pgxmql/commit/c68683cba85d0a4825670d2fa54261ef7085fd6a))
* replace QueryRewriter with WhereClause ([183e17d](https://github.com/pgx-contrib/pgxmql/commit/183e17d86ac0a1022764c80bf13ad05256948193))
* update README with reference to hashicorp/mql ([144ddd2](https://github.com/pgx-contrib/pgxmql/commit/144ddd2c1ea89d26abefe1d89f250d22d31440c7))
* update where clause to use exclude list ([3f804ea](https://github.com/pgx-contrib/pgxmql/commit/3f804ead69fd3bdeedda4a2b607fecf585cc343d))
* **where.go:** refactor options handling ([a7e43ae](https://github.com/pgx-contrib/pgxmql/commit/a7e43ae3fdbf0813c483d5d10cf2ec893f4d2057))
* **workflows:** split CI and PR workflows ([3acbbce](https://github.com/pgx-contrib/pgxmql/commit/3acbbce9d4a1667dae68992b65e5047be05cfa1f))


### Bug Fixes

* correct grammar in README.md ([2825dd9](https://github.com/pgx-contrib/pgxmql/commit/2825dd9143cf914d6eb09066539f21c9e2319cb1))
* **docs:** improve README formatting ([eb33146](https://github.com/pgx-contrib/pgxmql/commit/eb33146a30627591c4938a7f7e300aa4abddd008))
* remove debug print and add test cases ([2a40a8d](https://github.com/pgx-contrib/pgxmql/commit/2a40a8d6a4012c1f34cc31b5654737493a3b8289))
* rename filter to clause in WhereClause tests ([25183a3](https://github.com/pgx-contrib/pgxmql/commit/25183a32d296d1ab34bd8bfeda46d2632fb78802))
* update dependencies and refactor where.go ([19859de](https://github.com/pgx-contrib/pgxmql/commit/19859de016324026515103eeb36b4c8339a1d63b))
* **where.go:** correct RewriteQuery return types ([47fdc8f](https://github.com/pgx-contrib/pgxmql/commit/47fdc8f74079885d31f17a8b306a688209423e5a))
* **where.go:** improve error handling in RewriteQuery ([81131d1](https://github.com/pgx-contrib/pgxmql/commit/81131d1235639a2d2ec7a425e65cd442ef7a9501))
* **where.go:** replace condition placeholder with TRUE ([4eb9955](https://github.com/pgx-contrib/pgxmql/commit/4eb99551460ccb6d630174063f79faf40a8de549))
