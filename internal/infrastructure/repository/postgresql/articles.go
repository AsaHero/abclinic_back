package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableArticles = "articles"
)

type articlesRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewArticlesRepo(db *postgres.PostgresDB) repository.Articles {
	return &articlesRepo{
		table: tableArticles,
		db:    db,
	}
}

func (r articlesRepo) Create(ctx context.Context, req *entity.Articles) error {
	queryBuilder := r.db.Sq.Builder.Insert(r.table).SetMap(
		map[string]interface{}{
			"guid":       req.GUID,
			"chapter_id": req.ChapterID,
			"info":       req.Info,
			"img":        req.Img,
			"side":       req.Side,
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

func (r articlesRepo) List(ctx context.Context, filter map[string]string) ([]*entity.Articles, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"chapter_id",
		"info",
		"img",
		"side",
		"created_at",
	).From(r.table).OrderBy("created_at asc")

	for k, v := range filter {
		switch k {
		case "chapter_id":
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

	var services []*entity.Articles
	for rows.Next() {
		var service entity.Articles
		if err := rows.Scan(
			&service.GUID,
			&service.ChapterID,
			&service.Info,
			&service.Img,
			&service.Side,
			&service.CreatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}

		services = append(services, &service)
	}

	return services, nil
}

func (r articlesRepo) Update(ctx context.Context, req *entity.Articles) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"info": req.Info,
			"img":  req.Img,
			"side": req.Side,
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

func (r articlesRepo) Delete(ctx context.Context, filter map[string]string) error {
	queryBuilder := r.db.Sq.Builder.Delete(r.table)

	for k, v := range filter {
		switch k {
		case "guid":
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(k, v))
		case "chapter_id":
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
