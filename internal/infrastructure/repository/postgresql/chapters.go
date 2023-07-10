package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableChapters = "chapters"
)

type chaptersRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewChaptersRepo(db *postgres.PostgresDB) repository.Chapters {
	return &chaptersRepo{
		table: tableChapters,
		db:    db,
	}
}

func (r chaptersRepo) Create(ctx context.Context, req *entity.Chapters) error {
	queryBuilder := r.db.Sq.Builder.Insert(r.table).SetMap(
		map[string]interface{}{
			"guid":       req.GUID,
			"title":      req.Title,
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

func (r chaptersRepo) List(ctx context.Context, filter map[string]string) ([]*entity.Chapters, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"title",
		"created_at",
	).From(r.table)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, r.table+" List")
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}

	var groups []*entity.Chapters
	for rows.Next() {
		var group entity.Chapters
		if err := rows.Scan(
			&group.GUID,
			&group.Title,
			&group.CreatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}

		groups = append(groups, &group)
	}

	return groups, nil
}

func (r chaptersRepo) Update(ctx context.Context, req *entity.Chapters) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"title": req.Title,
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

func (r chaptersRepo) Delete(ctx context.Context, id string) error {
	queryBuilder := r.db.Sq.Builder.Delete(r.table).Where(r.db.Sq.Equal("guid", id))

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
