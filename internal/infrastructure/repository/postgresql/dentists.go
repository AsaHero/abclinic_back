package postgresql

import (
	"context"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableName = "dentists"
)

type dentistsRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewDentistsRepo(db *postgres.PostgresDB) repository.Denstists {
	return &dentistsRepo{
		table: tableName,
		db:    db,
	}
}

func (r dentistsRepo) Get(ctx context.Context, id int64) (*entity.Dentists, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"id",
		"clone_name",
		"name",
		"info",
		"url",
		"side",
		"priority",
	).From(r.table).Where(r.db.Sq.Equal("id", id))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, r.table+" Get")
	}

	var d entity.Dentists

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&d.ID,
		&d.CloneName,
		&d.Name,
		&d.Info,
		&d.URL,
		&d.Side,
		&d.Priority,
	)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &d, nil
}
