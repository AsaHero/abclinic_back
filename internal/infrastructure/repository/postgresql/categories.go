package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableCategories = "categories"
)

type categoriesRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewCategoriesRepo(db *postgres.PostgresDB) repository.Categories {
	return &categoriesRepo{
		table: tableCategories,
		db:    db,
	}
}

func (r categoriesRepo) Create(ctx context.Context, req *entity.Categories) error {
	queryBuilder := r.db.Sq.Builder.Insert(r.table).SetMap(
		map[string]interface{}{
			"guid":        req.GUID,
			"title":       req.Title,
			"description": req.Description,
			"url":         req.URL,
			"created_at":  req.CreatedAt,
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

func (r categoriesRepo) List(ctx context.Context, filter map[string]string) ([]*entity.Categories, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"title",
		"description",
		"url",
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

	var categories []*entity.Categories
	for rows.Next() {
		var category entity.Categories
		if err := rows.Scan(
			&category.GUID,
			&category.Title,
			&category.Description,
			&category.URL,
			&category.CreatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (r categoriesRepo) Update(ctx context.Context, req *entity.Categories) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"title":       req.Title,
			"description": req.Description,
			"url":         req.URL,
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

func (r categoriesRepo) Delete(ctx context.Context, filter map[string]string) error {
	queryBuilder := r.db.Sq.Builder.Delete(r.table).Where(r.db.Sq.Equal("guid", filter["guid"]))

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
