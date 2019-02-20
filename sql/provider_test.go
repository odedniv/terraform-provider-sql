package sql

import (
	"os"
	"testing"

	_sql "database/sql"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var pgClient *_sql.DB

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sql": testAccProvider,
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
	pgConnStr := os.Getenv("PGCONN")
	if pgConnStr == "" {
		t.Fatal("PGCONN must be set for acceptance tests")
	}
	if pgClient == nil {
		var err error
		pgClient, err = _sql.Open("postgres", pgConnStr)
		if err != nil {
			t.Fatal(err)
		}
	}
}
