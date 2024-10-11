Database local proccess:
Step 1: Make sure to have sqlc and goose locally.
`sqlc version`

`goose -version`

Step 2: create the sql migrations.

Step 3: run migrations.
`goose postgres ${URL} up`
