package functions

import (
	"CatsSocial/db/models"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Match struct {
	dbPool *pgxpool.Pool
}

func NewMatch(dbPool *pgxpool.Pool) *Match {
	return &Match{
		dbPool: dbPool,
	}
}

func (m *Match) Create(ctx context.Context, match models.Match) error {
	conn, err := m.dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed acquire connection from db pool: %v", err)
	}

	defer conn.Release()

	_, err = conn.Exec(ctx, `INSERT INTO matches (user_id, match_user_id, match_cat_id, user_cat_id, message, status) values($1, $2, $3, $4, $5, $6)`,
		match.UserId, match.MatchUserId, match.MatchCatId, match.UserCatId, match.Message, match.Status,
	)

	return err
}

func (m *Match) Get(ctx context.Context, userId string) ([]models.Match, error) {
	result := []models.Match{}

	conn, err := m.dbPool.Acquire(ctx)
	if err != nil {
		return result, fmt.Errorf("failed acquire connection from db pool: %v", err)
	}

	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT id, user_id, match_user_id, match_cat_id, user_cat_id, message, status, created_at FROM matches WHERE user_id = $1 AND status != 'removed'`, userId)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var match models.Match

		err := rows.Scan(&match.Id, &match.MatchCatId, &match.UserCatId, &match.Message, &match.Status, &match.CreatedAt)
		if err != nil {
			return []models.Match{}, err
		}

		result = append(result, match)
	}

	return result, nil
}

func (m *Match) GetMatchById(ctx context.Context, matchId string) (models.Match, error) {
	result := models.Match{}

	conn, err := m.dbPool.Acquire(ctx)
	if err != nil {
		return result, fmt.Errorf("failed acquire connection from db pool: %v", err)
	}

	defer conn.Release()

	error_row := conn.QueryRow(ctx, `SELECT id, user_id, match_user_id, match_cat_id, user_cat_id, message, status, created_at FROM matches WHERE id = $1`, matchId).Scan(&result.Id, &result.UserId, &result.MatchUserId, &result.MatchCatId, &result.UserCatId, &result.Message, &result.Status, &result.CreatedAt)
	if error_row != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, ErrNoRow
		}
		return result, fmt.Errorf("failed get match: %v", err)
	}

	return result, nil
}

func (m *Match) GetRelatedMatches(ctx context.Context, userId string) ([]models.Match, error) {
	result := []models.Match{}

	conn, err := m.dbPool.Acquire(ctx)
	if err != nil {
		return result, fmt.Errorf("failed acquire connection from db pool: %v", err)
	}

	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT id, user_id, match_user_id, match_cat_id, user_cat_id, message, status, created_at FROM matches WHERE (user_id = $1 OR match_user_id = $2) AND status != 'removed'`, userId, userId)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var match models.Match

		err := rows.Scan(&match.Id, &match.UserId, &match.MatchUserId, &match.MatchCatId, &match.UserCatId, &match.Message, &match.Status, &match.CreatedAt)
		if err != nil {
			return []models.Match{}, err
		}

		result = append(result, match)
	}

	return result, nil
}

func (m *Match) GetRelatedCatMatches(ctx context.Context, catId int) ([]models.Match, error) {
	result := []models.Match{}

	conn, err := m.dbPool.Acquire(ctx)
	if err != nil {
		return result, fmt.Errorf("failed acquire connection from db pool: %v", err)
	}

	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT id, user_id, match_user_id, match_cat_id, user_cat_id, message, status, created_at FROM matches WHERE (match_cat_id = $1 OR user_cat_id = $2) AND (status != 'removed' OR status != 'approved')`, catId, catId)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var match models.Match

		err := rows.Scan(&match.Id, &match.UserId, &match.MatchUserId, &match.MatchCatId, &match.UserCatId, &match.Message, &match.Status, &match.CreatedAt)
		if err != nil {
			return []models.Match{}, err
		}

		result = append(result, match)
	}

	return result, nil
}

func (m *Match) Delete(ctx context.Context, userId, accId string) error {
	conn, err := m.dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("faield acquire connection from dbpool: %v", err)
	}

	defer conn.Release()

	var match models.Match

	err = conn.QueryRow(ctx, `select user_id from matches where id = $1`, accId).Scan(
		&match.UserId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRow
		}

		return err
	}

	if strconv.Itoa(match.UserId) != userId {
		return ErrUnauthorized
	}

	_, err = conn.Exec(ctx, "delete from matches where id = $1", accId)

	return err
}

func (m *Match) UpdateStatus(ctx context.Context, e models.Match, status string) error {
	conn, err := m.dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("faield acquire connection from dbpool: %v", err)
	}

	defer conn.Release()

	var match models.Match

	err = conn.QueryRow(ctx, `select user_id from matches where id = $1`, e.Id).Scan(
		&match.UserId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRow
		}

		return err
	}

	if match.UserId != e.UserId {
		return ErrUnauthorized
	}

	_, err = conn.Exec(ctx, `update matches set status = $1 where id = $2`, status, e.Id)

	return err
}
