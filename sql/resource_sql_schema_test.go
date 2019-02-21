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
)

func TestAccSqlSchema_postgres(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "sql_schema" "this" {
  driver     = "postgres"
  datasource = "%s"
  directory  = "migrations/default"
}`, os.Getenv("POSTGRES_DATA_SOURCE")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlSchemaExists("postgres", pgClient, "table_1", "table_2"),

					resource.TestCheckResourceAttr("sql_schema.this", "driver", "postgres"),
					resource.TestCheckResourceAttr("sql_schema.this", "datasource", os.Getenv("POSTGRES_DATA_SOURCE")),
					resource.TestCheckResourceAttr("sql_schema.this", "directory", "migrations/default"),
				),
			},
		},
	})
}

func TestAccSqlSchema_mysql(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "sql_schema" "this" {
  driver     = "mysql"
  datasource = "%s"
  directory  = "migrations/default"
}`, os.Getenv("MYSQL_DATA_SOURCE")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlSchemaExists("mysql", mysqlClient, "table_1", "table_2"),

					resource.TestCheckResourceAttr("sql_schema.this", "driver", "mysql"),
					resource.TestCheckResourceAttr("sql_schema.this", "datasource", os.Getenv("MYSQL_DATA_SOURCE")),
					resource.TestCheckResourceAttr("sql_schema.this", "directory", "migrations/default"),
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
				Config: fmt.Sprintf(`
resource "sql_schema" "this" {
  driver     = "postgres"
  datasource = "%s"
  directory  = "migrations/update/step1"
}`, os.Getenv("POSTGRES_DATA_SOURCE")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlSchemaExists("postgres", pgClient, "table_1", "table_2"),

					resource.TestCheckResourceAttr("sql_schema.this", "driver", "postgres"),
					resource.TestCheckResourceAttr("sql_schema.this", "datasource", os.Getenv("POSTGRES_DATA_SOURCE")),
					resource.TestCheckResourceAttr("sql_schema.this", "directory", "migrations/update/step1"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "sql_schema" "this" {
  driver     = "postgres"
  datasource = "%s"
  directory  = "migrations/update/step2"
}`, os.Getenv("POSTGRES_DATA_SOURCE")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlSchemaExists("postgres", pgClient, "table_1", "table_2", "table_3", "table_4"),

					resource.TestCheckResourceAttr("sql_schema.this", "driver", "postgres"),
					resource.TestCheckResourceAttr("sql_schema.this", "datasource", os.Getenv("POSTGRES_DATA_SOURCE")),
					resource.TestCheckResourceAttr("sql_schema.this", "directory", "migrations/update/step2"),
				),
			},
		},
	})
}

func testAccCheckSqlSchemaDestroy(s *terraform.State) (err error) {
	var rez int
	rez, err = getTableCount("postgres", pgClient, "table1", "table2", "table3", "table4")
	if err != nil {
		return
	}
	if rez != 0 {
		err = errors.New("Schema still exists after destroy")
	}
	rez, err = getTableCount("mysql", mysqlClient, "table1", "table2", "table3", "table4")
	if err != nil {
		return
	}
	if rez != 0 {
		err = errors.New("Schema still exists after destroy")
	}
	return
}

func testAccCheckSqlSchemaExists(dialect string, client *_sql.DB, tables ...interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) (err error) {
		rez, err := getTableCount(dialect, client, tables...)
		if err != nil {
			return
		}
		if rez != len(tables) {
			err = fmt.Errorf("Schema not complete: expected %d tables, found %d", len(tables), rez)
		}
		return
	}
}

func getTableCount(dialect string, client *_sql.DB, tables ...interface{}) (rez int, err error) {
	// Creating a list of variables according to amount of tables (e.g ["$1", "$2", "$3"])
	vars := make([]string, len(tables))
	for i, _ := range tables {
		switch dialect {
		case "postgres":
			vars[i] = fmt.Sprintf("$%d", i+1)
			break
		case "mysql":
			vars[i] = "?"
			break
		}
	}
	varsStr := strings.Join(vars, ", ")
	err = client.QueryRow(fmt.Sprintf(`
SELECT COUNT(*)
FROM information_schema.tables
WHERE table_name IN (%s)
`, varsStr), tables...).Scan(&rez)
	return
}
