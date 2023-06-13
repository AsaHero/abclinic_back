package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	errorspkg "github.com/AsaHero/abclinic/internal/errors"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDB struct {
	*pgxpool.Pool
	Sq *Squirrel
}

func NewPostgresDB(cfg config.Config) (*PostgresDB, error) {
	postgresDB := PostgresDB{Sq: NewSquirrel()}

	if err := postgresDB.connect(&cfg); err != nil {
		return nil, err
	}

	return &postgresDB, nil
}

func GetStrConfig(config *config.Config) string {
	var conn []string

	if len(config.DB.Host) != 0 {
		conn = append(conn, "host="+config.DB.Host)
	}

	if len(config.DB.Port) != 0 {
		conn = append(conn, "port="+config.DB.Port)
	}

	if len(config.DB.User) != 0 {
		conn = append(conn, "user="+config.DB.User)
	}

	if len(config.DB.Password) != 0 {
		conn = append(conn, "password="+config.DB.Password)
	}

	if len(config.DB.Name) != 0 {
		conn = append(conn, "dbname="+config.DB.Name)
	}

	if len(config.DB.SSLMode) != 0 {
		conn = append(conn, "sslmode="+config.DB.SSLMode)
	}

	return strings.Join(conn, " ")
}

func (p *PostgresDB) connect(cfg *config.Config) error {
	pgxpoolConfig, err := pgxpool.ParseConfig(GetStrConfig(cfg))
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), pgxpoolConfig)
	if err != nil {
		return fmt.Errorf("unable to connect database config: %w", err)
	}

	p.Pool = pool

	return nil
}

func (p *PostgresDB) Close() {
	p.Pool.Close()
}

func (p *PostgresDB) Error(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return errorspkg.ErrorConflict
		}
	}
	if err == pgx.ErrNoRows {
		return errorspkg.ErrorNotFound
	}
	return err
}

func (p *PostgresDB) ErrSQLBuild(err error, message string) error {
	return fmt.Errorf("error during sql build, %s: %w", message, err)
}
