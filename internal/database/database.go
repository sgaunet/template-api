// Package database manages database connections and migrations.
package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres" // dbmate postgres driver
	"github.com/sgaunet/dsn/v2/pkg/dsn"
)

//go:embed db/migrations/*.sql
var fs embed.FS

// Postgres is the database connection.
type Postgres struct {
	DB               *sql.DB
	pgDataSourceName dsn.DSN
}

// NewPostgres creates a new Postgres database connection.
func NewPostgres(pgdsn string) (*Postgres, error) {
	d, err := dsn.New(pgdsn)
	if err != nil {
		return nil, fmt.Errorf("could not parse dsn: %w", err)
	}
	db, err := sql.Open("postgres", d.GetPostgresUri())
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}
	p := &Postgres{DB: db, pgDataSourceName: d}
	return p, nil
}

// InitDB initializes the database.
func (p *Postgres) InitDB() error {
	u, err := url.Parse(genDbmateURI(p.pgDataSourceName))
	if err != nil {
		return fmt.Errorf("could not parse dbmate uri: %w", err)
	}
	db := dbmate.New(u)
	db.FS = fs
	db.AutoDumpSchema = false

	if err = db.Wait(); err != nil {
		return fmt.Errorf("could not wait for database: %w", err)
	}
	fmt.Println("Migrations:")
	migrations, err := db.FindMigrations()
	if err != nil {
		return fmt.Errorf("could not find migrations: %w", err)
	}
	for _, m := range migrations {
		fmt.Println(m.Version, m.FilePath)
	}
	fmt.Println("\nApplying...")
	if err := db.CreateAndMigrate(); err != nil {
		return fmt.Errorf("could not create and migrate database: %w", err)
	}
	return nil
}

// Close closes the database connection.
func (p *Postgres) Close() error {
	if err := p.DB.Close(); err != nil {
		return fmt.Errorf("could not close database connection: %w", err)
	}
	return nil
}

// GetDB returns the database connection.
func (p *Postgres) GetDB() *sql.DB {
	return p.DB
}

// genDbmateURI generates a dbmate URI.
func genDbmateURI(d dsn.DSN) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		d.GetUser(),
		d.GetPassword(),
		net.JoinHostPort(d.GetHost(), d.GetPort("5432")),
		d.GetDBName())
}

// WaitForDB waits for the database to be ready.
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
					if err := pg.Close(); err != nil {
						fmt.Printf("could not close pg connection in wait loop: %v\n", err)
					}
					close(chDBReady)
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while waiting for db: %w", ctx.Err())
	case <-chDBReady:
		return nil
	}
}
