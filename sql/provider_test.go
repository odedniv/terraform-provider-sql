package sql

import (
	"errors"
	"os"
	"testing"

	_sql "database/sql"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var pgClient *_sql.DB
var mysqlClient *_sql.DB
var dbErr error

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sql": testAccProvider,
	}

	pgDataSource := os.Getenv("POSTGRES_DATA_SOURCE")
	if pgDataSource == "" {
		dbErr = errors.New("POSTGRES_DATA_SOURCE must be set for acceptance tests")
		return
	}
	pgClient, dbErr = _sql.Open("postgres", pgDataSource)
	if dbErr != nil {
		return
	}
	mysqlDataSource := os.Getenv("MYSQL_DATA_SOURCE")
	if mysqlDataSource == "" {
		dbErr = errors.New("MYSQL_DATA_SOURCE must be set for acceptance tests")
		return
	}
	mysqlClient, dbErr = _sql.Open("mysql", mysqlDataSource)
	if dbErr != nil {
		return
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if dbErr != nil {
		t.Fatal(dbErr)
	}
}
