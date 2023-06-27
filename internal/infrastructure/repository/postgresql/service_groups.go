package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableServiceGroups = "service_groups"
)

type serviceGroupsRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewServiceGroupsRepo(db *postgres.PostgresDB) repository.ServiceGroups {
	return &serviceGroupsRepo{
		table: tableServiceGroups,
		db:    db,
	}
}

func (r serviceGroupsRepo) Create(ctx context.Context, req *entity.ServiceGroups) error {
	queryBuilder := r.db.Sq.Builder.Insert(r.table).SetMap(
		map[string]interface{}{
			"guid":       req.GUID,
			"name":       req.Name,
			"created_at": time.Now().Local(),
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

func (r serviceGroupsRepo) List(ctx context.Context, filter map[string]string) ([]*entity.ServiceGroups, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"name",
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

	var groups []*entity.ServiceGroups
	for rows.Next() {
		var group entity.ServiceGroups
		if err := rows.Scan(
			&group.GUID,
			&group.Name,
			&group.CreatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}

		groups = append(groups, &group)
	}

	return groups, nil
}

func (r serviceGroupsRepo) Update(ctx context.Context, req *entity.ServiceGroups) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"name": req.Name,
		},
	).Where(r.db.Sq.Equal("id", req.GUID))

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

func (r serviceGroupsRepo) Delete(ctx context.Context, id string) error {
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
