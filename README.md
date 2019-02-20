# terraform-provider-sql

Terraform provider for managing SQL schemas using migrations.

This plugin uses [golang-migrate/migrate](https://github.com/golang-migrate/migrate),
it is recommended to go read how it works before using this provider.

## Usage

In your Terraform configuration:

```terraform
resource "sql_schema" "this" {
  database = "dialect://user@host/database"
  source   = "file://migrations" # optional
}
```
