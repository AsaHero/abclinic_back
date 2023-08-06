package postgresql

import (
	"context"
	"fmt"

	"github.com/AsaHero/abclinic/internal/entity"
	"github.com/AsaHero/abclinic/internal/infrastructure/repository"
	"github.com/AsaHero/abclinic/internal/pkg/postgres"
)

type refreshTokenRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewRefreshTokenRepo(db *postgres.PostgresDB) repository.RefreshTokenRepo {
	return &refreshTokenRepo{
		tableName: "refresh_tokens",
		db:        db,
	}
}

func (r *refreshTokenRepo) Get(ctx context.Context, refreshToken string) (*entity.RefreshToken, error) {
	query := r.db.Sq.Builder.
		Select(
			"guid",
			"refresh_token",
			"expiry_date",
			"created_at",
		).
		From(r.tableName).
		Where(r.db.Sq.Equal("refresh_token", refreshToken))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, r.tableName+" read")
	}

	var m entity.RefreshToken
	err = r.db.QueryRow(ctx, sqlStr, args...).Scan(
		&m.GUID,
		&m.RefreshToken,
		&m.ExpiryDate,
		&m.CreatedAt,
	)

	if err != nil {
		return nil, r.db.Error(err)
	}

	return &m, nil
}

func (r *refreshTokenRepo) Create(ctx context.Context, m *entity.RefreshToken) error {
	clauses := map[string]interface{}{
		"guid":          m.GUID,
		"refresh_token": m.RefreshToken,
		"expiry_date":   m.ExpiryDate,
		"created_at":    m.CreatedAt,
	}

	sqlStr, args, err := r.db.Sq.Builder.Insert(r.tableName).SetMap(clauses).ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.tableName+" create")
	}

	if _, err = r.db.Exec(ctx, sqlStr, args...); err != nil {
		return r.db.Error(err)
	}
	return nil
}

func (r *refreshTokenRepo) Delete(ctx context.Context, refreshToken string) error {
	sqlStr, args, err := r.db.Sq.Builder.
		Delete(r.tableName).
		Where(r.db.Sq.Equal("refresh_token", refreshToken)).
		ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.tableName+" delete")
	}

	commandTag, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return r.db.Error(err)
	}
	if commandTag.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("failed token delete: no such token"))
	}
	return nil
}
