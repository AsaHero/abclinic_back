package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tablePublications = "publications"
)

type publicationsRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewPublicationsRepo(db *postgres.PostgresDB) repository.Publications {
	return &publicationsRepo{
		table: tablePublications,
		db:    db,
	}
}

func (r publicationsRepo) Create(ctx context.Context, req *entity.Publications) error {
	queryBuilder := r.db.Sq.Builder.Insert(r.table).SetMap(
		map[string]interface{}{
			"guid":        req.GUID,
			"category_id": req.CategoryID,
			"author_id":   req.AuthorID,
			"title":       req.Title,
			"description": req.Description,
			"type":        req.Type,
			"content":     req.Content,
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

func (r publicationsRepo) List(ctx context.Context, filter map[string]string) ([]*entity.Publications, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"category_id",
		"author_id",
		"title",
		"description",
		"type",
		"content",
		"created_at",
	).From(r.table).OrderBy("created_at asc")

	for k, v := range filter {
		switch k {
		case "category_id":
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

	var publications []*entity.Publications
	for rows.Next() {
		var publication entity.Publications
		if err := rows.Scan(
			&publication.GUID,
			&publication.CategoryID,
			&publication.AuthorID,
			&publication.Title,
			&publication.Description,
			&publication.Type,
			&publication.Content,
			&publication.CreatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}

		publications = append(publications, &publication)
	}

	return publications, nil
}

func (r publicationsRepo) Update(ctx context.Context, req *entity.Publications) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"title":       req.Title,
			"description": req.Description,
			"content":     req.Content,
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

func (r publicationsRepo) Delete(ctx context.Context, filter map[string]string) error {
	queryBuilder := r.db.Sq.Builder.Delete(r.table)

	for k, v := range filter {
		switch k {
		case "guid":
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(k, v))
		case "category_id":
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(k, v))
		case "author_id":
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
