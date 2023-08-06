package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

var (
	tableUsers = "users"
)

type usersRepo struct {
	table string
	db    *postgres.PostgresDB
}

func NewUsersRepo(db *postgres.PostgresDB) repository.Users {
	return &usersRepo{
		table: tableUsers,
		db:    db,
	}
}

func (r usersRepo) Get(ctx context.Context, filter map[string]string) (*entity.Users, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"role",
		"firstname",
		"lastname",
		"username",
		"password",
		"created_at",
		"updated_at",
	).From(r.table)

	for k, v := range filter {
		switch k {
		case "guid", "firstname", "lastname", "username":
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(k, v))
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, r.table+" Get")
	}

	var user entity.Users
	err = r.db.QueryRow(ctx, query, args...).Scan(&user)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &user, nil
}

func (r usersRepo) Create(ctx context.Context, req *entity.Users) error {
	queryBuilder := r.db.Sq.Builder.Insert(r.table).SetMap(
		map[string]interface{}{
			"guid":       req.GUID,
			"role":       req.Role,
			"firstname":  req.Firstname,
			"lastname":   req.Lastname,
			"username":   req.Username,
			"password":   req.Password,
			"created_at": req.CreatedAt,
			"updated_at": req.UpdatedAt,
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

func (r usersRepo) List(ctx context.Context, filter map[string]string) ([]*entity.Users, error) {
	queryBuilder := r.db.Sq.Builder.Select(
		"guid",
		"role",
		"firstname",
		"lastname",
		"username",
		"password",
		"created_at",
		"updated_at",
	).From(r.table)

	for k, v := range filter {
		switch k {
		case "role":
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

	var users []*entity.Users
	for rows.Next() {
		var user entity.Users
		if err := rows.Scan(
			&user.GUID,
			&user.Role,
			&user.Firstname,
			&user.Lastname,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r usersRepo) Update(ctx context.Context, req *entity.Users) error {
	queryBuilder := r.db.Sq.Builder.Update(r.table).SetMap(
		map[string]interface{}{
			"role":       req.Role,
			"firstname":  req.Firstname,
			"lastname":   req.Lastname,
			"username":   req.Username,
			"password":   req.Password,
			"updated_at": req.UpdatedAt,
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

func (r usersRepo) Delete(ctx context.Context, filter map[string]string) error {
	queryBuilder := r.db.Sq.Builder.Delete(r.table)

	for k, v := range filter {
		switch k {
		case "guid":
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(k, v))
		case "username":
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
