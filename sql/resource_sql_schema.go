package sql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	migrate "github.com/rubenv/sql-migrate"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var dialects = map[string]string{
	"mysql":            "mysql",
	"postgres":         "postgres",
	"cloudsql":         "mysql",
	"cloudsqlpostgres": "postgres",
}
var availableDialects = []string{"mysql", "postgres", "cloudsql", "cloudsqlpostgres"}

func resourceSQLSchema() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSQLSchemaCreate,
		Read:          resourceSQLSchemaRead,
		Update:        resourceSQLSchemaUpdate,
		Delete:        resourceSQLSchemaDelete,
		CustomizeDiff: resourceSQLSchemaCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"driver": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Database driver. Available drivers: %s", strings.Join(availableDialects, ", ")),
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if dialects[v.(string)] == "" {
						errors = append(errors, fmt.Errorf("Unknown database driver: %s, only know: %s", v, strings.Join(availableDialects, ", ")))
					}
					return
				},
			},
			"datasource": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Database connection string as compatible with sql.Open.",
			},
			"directory": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "migrations",
				Description: "Directory of the migrations.",
			},
			"table": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "schema_migrations",
				Description: "Name of the table to use to store applied migrations.",
			},
			"migrations": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSQLSchemaCreate(d *schema.ResourceData, m interface{}) (err error) {
	db, err := getDatabase(d)
	if err != nil {
		return
	}
	_, err = migrate.Exec(db, getDialect(d), getSource(d), migrate.Up)
	if err != nil {
		return
	}
	d.SetId(d.Get("datasource").(string))
	return resourceSQLSchemaRead(d, m)
}

func resourceSQLSchemaRead(d *schema.ResourceData, m interface{}) (err error) {
	db, err := getDatabase(d)
	if err != nil {
		return
	}
	databaseMigrations, err := migrate.GetMigrationRecords(db, getDialect(d))
	if err != nil {
		return
	}
	migrations := make([]string, len(databaseMigrations))
	for i, databaseMigration := range databaseMigrations {
		migrations[i] = databaseMigration.Id
	}
	d.Set("migrations", migrations)
	return
}

func resourceSQLSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSQLSchemaCreate(d, m)
}

func resourceSQLSchemaDelete(d *schema.ResourceData, m interface{}) (err error) {
	db, err := getDatabase(d)
	if err != nil {
		return
	}
	_, err = migrate.Exec(db, getDialect(d), getSource(d), migrate.Down)
	if err != nil {
		return
	}
	return
}

func resourceSQLSchemaCustomizeDiff(d *schema.ResourceDiff, m interface{}) (err error) {
	sourceMigrations, err := getSource(d).FindMigrations()
	if err != nil {
		return
	}
	migrations := make([]string, len(sourceMigrations))
	for i, sourceMigration := range sourceMigrations {
		migrations[i] = sourceMigration.Id
	}
	err = d.SetNew("migrations", migrations)
	return
}

type schemaResource interface {
	Get(key string) interface{}
}

func getSource(d schemaResource) migrate.MigrationSource {
	return migrate.FileMigrationSource{Dir: d.Get("directory").(string)}
}

func getDatabase(d schemaResource) (*sql.DB, error) {
	migrate.SetTable(d.Get("table").(string))
	driver := d.Get("driver").(string)
	dataSource := d.Get("datasource").(string)

	switch driver {
	case "mysql", "cloudsql":
		if strings.Contains(dataSource, "?") {
			dataSource += "&parseTime=true"
		} else {
			dataSource += "?parseTime=true"
		}
		break
	}

	return sql.Open(driver, dataSource)
}

func getDialect(d schemaResource) string {
	return dialects[d.Get("driver").(string)]
}
