package sql

import (
	"errors"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSQLSchema() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSQLSchemaCreate,
		Read:          resourceSQLSchemaRead,
		Update:        resourceSQLSchemaUpdate,
		Delete:        resourceSQLSchemaDelete,
		Exists:        resourceSQLSchemaExists,
		CustomizeDiff: resourceSQLSchemaCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"database": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Run migrations against this database (driver://url). Only MySQL (mysql://url) and PostgreSQL (postgres://url) drivers available.",
			},
			"source": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "file://migrations",
				Description: "Location of the migrations (driver://url). Only Filesystem (file://url) driver available.",
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSQLSchemaCreate(d *schema.ResourceData, m interface{}) (err error) {
	mig, err := buildMigrate(d)
	if err != nil {
		return
	}
	err = mig.Up()
	if err != nil {
		return
	}
	d.SetId(d.Get("database").(string))
	return resourceSQLSchemaRead(d, m)
}

func resourceSQLSchemaRead(d *schema.ResourceData, m interface{}) (err error) {
	mig, err := buildMigrate(d)
	if err != nil {
		return
	}
	v, dirty, err := mig.Version()
	if err == migrate.ErrNilVersion {
		err = nil
		d.SetId("")
		return
	} else if dirty {
		if err != nil {
			err = errors.New("database is dirty")
		}
	}
	if err != nil {
		return
	}
	d.Set("version", strconv.FormatUint(uint64(v), 10))
	return
}

func resourceSQLSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSQLSchemaCreate(d, m)
}

func resourceSQLSchemaDelete(d *schema.ResourceData, m interface{}) (err error) {
	mig, err := buildMigrate(d)
	if err != nil {
		return
	}
	err = mig.Down()
	if os.IsNotExist(err) {
		err = nil
	}
	return
}

func resourceSQLSchemaExists(d *schema.ResourceData, m interface{}) (result bool, err error) {
	mig, err := buildMigrate(d)
	if err != nil {
		return
	}
	_, _, err = mig.Version()
	if err == migrate.ErrNilVersion {
		err = nil
		result = false
	} else {
		result = true
	}
	return
}

func resourceSQLSchemaCustomizeDiff(d *schema.ResourceDiff, m interface{}) (err error) {
	drv, err := source.Open(d.Get("source").(string))
	if err != nil {
		return
	}
	v, err := drv.First()
	if err != nil {
		return
	}
	var n uint
	for {
		n, err = drv.Next(v)
		if os.IsNotExist(err) {
			err = nil
			break
		} else if err != nil {
			return
		} else {
			v = n
		}
	}
	err = d.SetNew("version", strconv.FormatUint(uint64(v), 10))
	return
}

func buildMigrate(d *schema.ResourceData) (*migrate.Migrate, error) {
	return migrate.New(d.Get("source").(string), d.Get("database").(string))
}
