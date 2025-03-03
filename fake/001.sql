-- sqlfluff:dialect:postgres
-- sqlfluff:max_line_length:1024
-- sqlfluff:rules:capitalisation.keywords:capitalisation_policy:upper

SELECT
    id,
    role,
    company,
    password
FROM
    users
WHERE
    -- We need a void condition parameter. So we can inject
    -- an instance of pgx.QueryRewriter that writes the dynamic condition
    -- provided by the end user.
    $1::void IS NULL
    -- :condition
    AND id::text >= $2::text
    AND ($3::text IS NULL OR company::text = $3::text)
ORDER BY
    id
LIMIT
    $5::int
    OFFSET
    $4::int;
