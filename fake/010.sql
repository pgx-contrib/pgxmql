-- sqlfluff:dialect:postgres
-- sqlfluff:max_line_length:1024
-- sqlfluff:rules:capitalisation.keywords:capitalisation_policy:upper

-- name: middle
SELECT name FROM users
WHERE
    id = $1
    -- OR :condition
    AND name = $2
    AND perm = $3;
