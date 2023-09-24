package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableDentists = "dentists"
)

type dentistsRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewDentistsRepo(db *postgres.PostgresDB) repository.Denstists {
	return &dentistsRepo{
		table: tableDentists,
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

func (r dentistsRepo) List(ctx context.Context, filter map[string]string) ([]*entity.Dentists, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"id",
		"clone_name",
		"name",
		"info",
		"url",
		"side",
		"priority",
	).From(r.table)

	for k, v := range filter {
		switch k {
		case "side":
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(k, v))
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, r.table+" List")
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}

	dentists := []*entity.Dentists{}
	for rows.Next() {
		var dentist entity.Dentists
		rows.Scan(
			&dentist.ID,
			&dentist.CloneName,
			&dentist.Name,
			&dentist.Info,
			&dentist.URL,
			&dentist.Side,
			&dentist.Priority,
		)
		dentists = append(dentists, &dentist)
	}

	return dentists, nil
}
func (r dentistsRepo) Update(ctx context.Context, req *entity.Dentists) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"info": req.Info,
			"name": req.Name,
			"url":  req.URL,
		},
	).Where(r.db.Sq.Equal("id", req.ID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.table+" Update")
	}

	commandTag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return r.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}
