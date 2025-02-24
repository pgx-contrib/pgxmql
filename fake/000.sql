-- sqlfluff:dialect:postgres
-- sqlfluff:max_line_length:1024
-- sqlfluff:rules:capitalisation.keywords:capitalisation_policy:upper

-- name: first
SELECT name FROM users;
-- WHERE :condition;
