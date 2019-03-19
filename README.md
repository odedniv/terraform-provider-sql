# terraform-provider-sql

Terraform provider for managing SQL schemas using migrations.

This plugin uses [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate),
it is recommended to go read how it works before using this provider.

## Usage

### Installation

Build the provider and put it in Terraform's third-party providers directory in `~/.terraform.d/plugins`:

```bash
go get github.com/odedniv/terraform-provider-sql
mkdir -p ~/.terraform.d/plugins
go build -o ~/.terraform.d/plugins/terraform-provider-sql github.com/odedniv/terraform-provider-sql
```

I recommend using [Go modules](https://github.com/golang/go/wiki/Modules) to ensure
using the same version in development and production.

### Configuration

In your Terraform configuration:

```terraform
resource "sql_schema" "this" {
  driver     = "<database driver>" # mysql/postgres/cloudsql/cloudsqlpostgres
  datasource = "<database connection string>"
  directory  = "migrations" # optional
  table      = "schema_migrations" # optional
}
```
