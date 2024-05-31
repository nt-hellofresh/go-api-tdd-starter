package testdata

import (
	"database/sql"
	"log"
	"playground/migrations"
	"sync"
	"testing"
)

type FixturesTestRunner interface {
	Run(name string, fn func(t *testing.T)) bool
}

type dataFixture struct {
	db *sql.DB
	t  *testing.T
	mu *sync.Mutex
}

func NewTestFixture(t *testing.T, db *sql.DB) FixturesTestRunner {
	return &dataFixture{
		db: db,
		t:  t,
		mu: &sync.Mutex{},
	}
}

func (fixture *dataFixture) Run(name string, fn func(t *testing.T)) bool {
	fixture.mu.Lock()
	defer fixture.mu.Unlock()

	fixture.setup()
	defer fixture.teardown()

	return fixture.t.Run(name, fn)
}

func (fixture *dataFixture) setup() {
	fixture.t.Helper()
	log.Println("setting up")

	if err := migrations.MigrateUp(fixture.db); err != nil {
		fixture.t.Error(err)
	}

	fixture.addProduct("123", "Noise cancelling headphones", 295.00)
	fixture.addProduct("456", "Mechanical keyboard", 150.00)
	fixture.addProduct("789", "Wireless mouse", 79.95)

	log.Println("db setup complete")
}

func (fixture *dataFixture) addProduct(ID string, name string, price float64) {
	query := "INSERT INTO product (id, name, price) VALUES ($1, $2, $3)"
	if _, err := fixture.db.Exec(query, ID, name, price); err != nil {
		fixture.t.Error(err)
	}
}

func (fixture *dataFixture) teardown() {
	fixture.t.Helper()
	log.Println("tearing down")

	if err := migrations.ResetMigration(fixture.db); err != nil {
		fixture.t.Error(err)
	}

	log.Println("db teardown complete")
}
