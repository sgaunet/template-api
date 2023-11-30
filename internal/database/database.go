package database

import (
	"database/sql"
	"embed"
	"fmt"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/sgaunet/dsn/v2/pkg/dsn"
)

//go:embed db/migrations/*.sql
var fs embed.FS

type Postgres struct {
	DB               *sql.DB
	pgDataSourceName dsn.DSN
}

func NewPostgres(pgdsn string) (*Postgres, error) {
	d, err := dsn.New(pgdsn)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", d.GetPostgresUri())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	p := &Postgres{DB: db, pgDataSourceName: d}
	err = p.InitDB()
	return p, err
}

func (p *Postgres) InitDB() error {
	u, err := url.Parse(genDbmateUri(p.pgDataSourceName))
	if err != nil {
		return err
	}
	db := dbmate.New(u)
	db.FS = fs

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

func (p *Postgres) Close() error {
	return p.DB.Close()
}

func (p *Postgres) GetDB() *sql.DB {
	return p.DB
}

func genDbmateUri(d dsn.DSN) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.GetUser(),
		d.GetPassword(),
		d.GetHost(),
		d.GetPort("5432"),
		d.GetDBName())
}
