# SpendSense

## Migrations

```sh
atlas schema apply \
    --url "postgres://postgres:test1234@localhost:5432/expense_tracker?sslmode=disable" \
    --dev-url "docker://postgres/15" \
    --to "file://migrations"
```
