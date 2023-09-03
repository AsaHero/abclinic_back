package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableAuthors = "authors"
)

type authorsRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewAuthorsRepo(db *postgres.PostgresDB) repository.Authors {
	return &authorsRepo{
		table: tableAuthors,
		db:    db,
	}
}

func (r authorsRepo) Create(ctx context.Context, req *entity.Authors) error {
	queryBuilder := r.db.Sq.Builder.Insert(r.table).SetMap(
		map[string]interface{}{
			"guid":       req.GUID,
			"name":       req.Name,
			"img":        req.Img,
			"created_at": req.CreatedAt,
		},
	)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.table+" Create")
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return r.db.Error(err)
	}

	return nil
}

func (r authorsRepo) Get(ctx context.Context, guid string) (*entity.Authors, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"name",
		"img",
		"created_at",
	).From(r.table).Where(r.db.Sq.Equal("guid", guid))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, r.table+" Get")
	}

	var author entity.Authors
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&author.GUID,
		&author.Name,
		&author.Img,
		&author.CreatedAt)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &author, nil
}

func (r authorsRepo) List(ctx context.Context, filter map[string]string) ([]*entity.Authors, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"name",
		"img",
		"created_at",
	).From(r.table).OrderBy("created_at asc")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, r.table+" List")
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}

	var authors []*entity.Authors
	for rows.Next() {
		var author entity.Authors
		if err := rows.Scan(
			&author.GUID,
			&author.Name,
			&author.Img,
			&author.CreatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}

		authors = append(authors, &author)
	}

	return authors, nil
}

func (r authorsRepo) Update(ctx context.Context, req *entity.Authors) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"name": req.Name,
			"img":  req.Img,
		},
	).Where(r.db.Sq.Equal("guid", req.GUID))

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

func (r authorsRepo) Delete(ctx context.Context, filter map[string]string) error {
	queryBuilder := r.db.Sq.Builder.Delete(r.table)

	for k, v := range filter {
		switch k {
		case "guid":
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(k, v))
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.table+" Delete")
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
