# terraform-provider-sql

Terraform provider for managing SQL schemas using migrations.

This plugin uses [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate),
it is recommended to go read how it works before using this provider.

## Usage

In your Terraform configuration:

```terraform
resource "sql_schema" "this" {
  driver     = "<database driver>" # mysql/postgres/cloudsql/cloudsqlpostgres
  datasource = "<database connection string>"
  directory  = "migrations" # optional
  table      = "schema_migrations" # optional
}
```
