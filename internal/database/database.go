package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres" // dbmate postgres driver
	"github.com/sgaunet/dsn/v2/pkg/dsn"
)

//go:embed db/migrations/*.sql
var fs embed.FS

// Postgres is the database connection
type Postgres struct {
	DB               *sql.DB
	pgDataSourceName dsn.DSN
}

// NewPostgres creates a new Postgres database connection
func NewPostgres(pgdsn string) (*Postgres, error) {
	d, err := dsn.New(pgdsn)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", d.GetPostgresUri())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	p := &Postgres{DB: db, pgDataSourceName: d}
	return p, err
}

// InitDB initializes the database
func (p *Postgres) InitDB() error {
	u, err := url.Parse(genDbmateURI(p.pgDataSourceName))
	if err != nil {
		return err
	}
	db := dbmate.New(u)
	db.FS = fs
	db.AutoDumpSchema = false

	err = db.Wait()
	if err != nil {
		return err
	}
	fmt.Println("Migrations:")
	migrations, err := db.FindMigrations()
	if err != nil {
		return err
	}
	for _, m := range migrations {
		fmt.Println(m.Version, m.FilePath)
	}
	fmt.Println("\nApplying...")
	return db.CreateAndMigrate()
}

// Close closes the database connection
func (p *Postgres) Close() error {
	return p.DB.Close()
}

// GetDB returns the database connection
func (p *Postgres) GetDB() *sql.DB {
	return p.DB
}

// genDbmateURI generates a dbmate URI
func genDbmateURI(d dsn.DSN) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.GetUser(),
		d.GetPassword(),
		d.GetHost(),
		d.GetPort("5432"),
		d.GetDBName())
}

// WaitForDB waits for the database to be ready
func WaitForDB(ctx context.Context, pgdsn string) error {
	chDBReady := make(chan struct{})
	go func() {
		for {
			pg, err := NewPostgres(pgdsn)
			select {
			case <-ctx.Done():
				return
			default:
				if err == nil {
					pg.Close()
					close(chDBReady)
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-chDBReady:
		return nil
	}
}
