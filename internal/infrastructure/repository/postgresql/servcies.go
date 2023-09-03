package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableServices = "services"
)

type servicesRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewServicesRepo(db *postgres.PostgresDB) repository.Services {
	return &servicesRepo{
		table: tableServices,
		db:    db,
	}
}

func (r servicesRepo) Create(ctx context.Context, req *entity.Services) error {
	queryBuilder := r.db.Sq.Builder.Insert(r.table).SetMap(
		map[string]interface{}{
			"guid":       req.GUID,
			"group_id":   req.GroupID,
			"name":       req.Name,
			"price":      req.Price,
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

func (r servicesRepo) List(ctx context.Context, filter map[string]string) ([]*entity.Services, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"group_id",
		"name",
		"price",
		"created_at",
	).From(r.table).OrderBy("created_at asc")

	for k, v := range filter {
		switch k {
		case "group_id":
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

	var services []*entity.Services
	for rows.Next() {
		var service entity.Services
		if err := rows.Scan(
			&service.GUID,
			&service.GroupID,
			&service.Name,
			&service.Price,
			&service.CreatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}

		services = append(services, &service)
	}

	return services, nil
}

func (r servicesRepo) Update(ctx context.Context, req *entity.Services) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"name":  req.Name,
			"price": req.Price,
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

func (r servicesRepo) Delete(ctx context.Context, filter map[string]string) error {
	queryBuilder := r.db.Sq.Builder.Delete(r.table)

	for k, v := range filter {
		switch k {
		case "guid":
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(k, v))
		case "group_id":
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
