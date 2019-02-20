package sql

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	_sql "database/sql"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	_ "github.com/lib/pq"
)

func TestAccSqlSchema_defaultSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlSchemaCreateDefaultSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlSchemaExists("table_1", "table_2"),

					resource.TestCheckResourceAttr("sql_schema.this", "database", os.Getenv("PGCONN")),
					resource.TestCheckResourceAttr("sql_schema.this", "source", "file://migrations"),
				),
			},
		},
	})
}

func TestAccSqlSchema_explicitSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlSchemaCreateExplicitSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlSchemaExists("table_3", "table_4"),

					resource.TestCheckResourceAttr("sql_schema.this", "database", os.Getenv("PGCONN")),
					resource.TestCheckResourceAttr("sql_schema.this", "source", "file://migrations/explicit"),
				),
			},
		},
	})
}

func TestAccSqlSchema_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlSchemaUpdateStep1Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlSchemaExists("table_1", "table_2"),

					resource.TestCheckResourceAttr("sql_schema.this", "database", os.Getenv("PGCONN")),
					resource.TestCheckResourceAttr("sql_schema.this", "source", "file://migrations/update/step1"),
				),
			},
			{
				Config: testAccSqlSchemaUpdateStep2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlSchemaExists("table_1", "table_2", "table_3", "table_4"),

					resource.TestCheckResourceAttr("sql_schema.this", "database", os.Getenv("PGCONN")),
					resource.TestCheckResourceAttr("sql_schema.this", "source", "file://migrations/update/step2"),
				),
			},
		},
	})
}

func testAccCheckSqlSchemaDestroy(s *terraform.State) (err error) {
	var _rez int
	err = pgClient.QueryRow(`
SELECT 1
FROM information_schema.tables
WHERE table_schema = 'public' AND table_name != 'schema_migrations'
`).Scan(&_rez)
	if err == _sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		err = errors.New("Schema still exists after destroy")
	}
	return
}

func testAccCheckSqlSchemaExists(tables ...interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) (err error) {
		var rez int
		// Creating a list of variables according to amount of tables (e.g ["$1", "$2", "$3"])
		vars := make([]string, len(tables))
		for i, _ := range tables {
			vars[i] = fmt.Sprintf("$%d", i+1)
		}
		varsStr := strings.Join(vars, ", ")
		err = pgClient.QueryRow(fmt.Sprintf(`
SELECT COUNT(*)
FROM information_schema.tables
WHERE table_schema = 'public' AND table_name IN (%s)
`, varsStr), tables...).Scan(&rez)
		if err != nil {
			return
		}
		if rez != len(tables) {
			err = fmt.Errorf("Schema not complete: expected %d tables, found %d", len(tables), rez)
		}
		return
	}
}

var testAccSqlSchemaCreateDefaultSourceConfig = fmt.Sprintf(`
resource "sql_schema" "this" {
  database = "%s"
}`, os.Getenv("PGCONN"))

var testAccSqlSchemaCreateExplicitSourceConfig = fmt.Sprintf(`
resource "sql_schema" "this" {
  database = "%s"
  source   = "file://migrations/explicit"
}`, os.Getenv("PGCONN"))

var testAccSqlSchemaUpdateStep1Config = fmt.Sprintf(`
resource "sql_schema" "this" {
  database = "%s"
  source   = "file://migrations/update/step1"
}`, os.Getenv("PGCONN"))

var testAccSqlSchemaUpdateStep2Config = fmt.Sprintf(`
resource "sql_schema" "this" {
  database = "%s"
  source   = "file://migrations/update/step2"
}`, os.Getenv("PGCONN"))
