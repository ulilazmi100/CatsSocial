package functions

import (
	"CatsSocial/db/models"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Cat struct {
	dbPool *pgxpool.Pool
}

func NewCatFn(dbPool *pgxpool.Pool) *Cat {
	return &Cat{
		dbPool: dbPool,
	}
}

// FilterGetCats struct {
// 	Id                 string `json:"id"`
// 	Limit              int    `json:"limit"`
// 	Offset             int    `json:"offset"`
// 	Race               string `json:"Race"`
// 	Sex                string `json:"sex"`
// 	HasMatched         bool   `json:"hasMatched"`
// 	AgeInMonthOperator string `json:"ageInMonthOperator"`
// 	AgeInMonthValue    int    `json:"ageInMonthValue"`
// 	Owned              bool   `json:"owned"`
// 	Search             string `json:"search"`
// }

func (p *Cat) constructWhereQuery(ctx context.Context, filter models.FilterGetCats, userID int) string {
	whereSQL := []string{}
	if filter.Owned {
		whereSQL = append(whereSQL, " user_id = "+fmt.Sprintf("%d", userID))
	}

	if filter.Id != "" {
		whereSQL = append(whereSQL, " id = '"+filter.Id+"'")
	}

	if filter.Race != "" {
		whereSQL = append(whereSQL, " race = '"+filter.Race+"'")
	}

	if filter.Sex != "" {
		whereSQL = append(whereSQL, " sex = '"+filter.Sex+"'")
	}

	if filter.HasMatched {
		whereSQL = append(whereSQL, " hasMatched = '"+"1"+"'")
	}
	// else {
	// 	whereSQL = append(whereSQL, " hasMatched = '"+"0"+"'")
	// }

	if filter.AgeInMonthOperator != "" {
		if filter.AgeInMonthOperator == "=>" {
			whereSQL = append(whereSQL, " age_in_month >= "+fmt.Sprintf("%d", filter.AgeInMonthValue))
		} else if filter.AgeInMonthOperator == "=<" {
			whereSQL = append(whereSQL, " age_in_month <= "+fmt.Sprintf("%d", filter.AgeInMonthValue))
		} else {
			whereSQL = append(whereSQL, " age_in_month = "+fmt.Sprintf("%d", filter.AgeInMonthValue))
		}
	}

	if filter.Search != "" {
		whereSQL = append(whereSQL, " name ILIKE '%"+filter.Search+"%'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}

func (p *Cat) FindAll(ctx context.Context, filter models.FilterGetCats, userID int) ([]models.Cat, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `SELECT id, user_id, name, race, sex, age_in_month, description, image_urls, has_matched, created_at FROM cats`

	sql += p.constructWhereQuery(ctx, filter, userID)

	sql += " ORDER BY " + "created_at" + " " + "DESC"

	if filter.Limit > 0 {
		sql += " LIMIT " + fmt.Sprintf("%d", filter.Limit)
	}

	if filter.Offset >= 0 {
		sql += " OFFSET " + fmt.Sprintf("%d", filter.Offset)
	}

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed get cats: %v", err)
	}

	defer rows.Close()

	cats := []models.Cat{}

	for rows.Next() {
		cat := models.Cat{}
		err := rows.Scan(&cat.Id, &cat.UserId, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.Description, &cat.ImageUrls, &cat.HasMatched, &cat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed scan cats: %v", err)
		}
		cats = append(cats, cat)
	}

	return cats, nil
}

func (p *Cat) Count(ctx context.Context, filter models.FilterGetCats, userID int) (int, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `SELECT COUNT(id) FROM cats`

	sql += p.constructWhereQuery(ctx, filter, userID)

	var count int
	err = conn.QueryRow(ctx, sql).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed get cats count: %v", err)
	}

	return count, nil
}

func (p *Cat) Add(ctx context.Context, cat models.Cat) (models.Cat, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return models.Cat{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `
		INSERT INTO cats (user_id, name, race, sex, age_in_month, description, image_urls) 
		values ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at
	`

	var result models.Cat

	err = conn.QueryRow(ctx, sql, cat.UserId,
		cat.Name,
		cat.Race,
		cat.Sex,
		cat.AgeInMonth,
		cat.Description,
		cat.ImageUrls).Scan(&result.Id, &result.CreatedAt)

	if err != nil {
		return models.Cat{}, fmt.Errorf("failed insert cat: %v", err)
	}

	cat.Id = result.Id
	cat.CreatedAt = result.CreatedAt

	return cat, nil
}

func (p *Cat) Update(ctx context.Context, cat models.Cat) error {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `
		update cats set name = $1, race = $2, sex = $3, age_in_month = $4, description = $5, image_urls = $6, has_matched = $7, updated_at = now()
		where id = $8 and user_id = $9
	`

	_, err = conn.Exec(ctx, sql,
		cat.Name,
		cat.Race,
		cat.Sex,
		cat.AgeInMonth,
		cat.Description,
		cat.ImageUrls,
		cat.HasMatched,
		cat.Id,
		cat.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRow
		}
		return fmt.Errorf("failed update cat: %v", err)
	}

	return nil
}

func (p *Cat) FindByID(ctx context.Context, catID int) (models.Cat, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return models.Cat{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	var cat models.Cat

	err = conn.QueryRow(ctx, `SELECT id, user_id, name, race, sex, age_in_month, description, image_urls, has_matched, created_at FROM cats WHERE id = $1`, catID).Scan(
		&cat.Id, &cat.UserId, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.Description, &cat.ImageUrls, &cat.HasMatched, &cat.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Cat{}, ErrNoRow
		}
		return models.Cat{}, fmt.Errorf("failed get cat: %v", err)
	}

	return cat, nil
}

func (p *Cat) FindByIDUser(ctx context.Context, catID int, userID int) (models.Cat, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return models.Cat{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	var cat models.Cat

	err = conn.QueryRow(ctx, `SELECT id, user_id, name, race, sex, age_in_month, description, image_urls, has_matched, created_at FROM cats WHERE id = $1 AND user_id = $2`, catID, userID).Scan(
		&cat.Id, &cat.UserId, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.Description, &cat.ImageUrls, &cat.HasMatched, &cat.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Cat{}, ErrNoRow
		}
		return models.Cat{}, fmt.Errorf("failed get cat: %v", err)
	}

	return cat, nil
}

func (p *Cat) DeleteByID(ctx context.Context, catID int) error {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `delete from cats where id = $1`
	_, err = conn.Exec(ctx, sql, catID)
	if err != nil {
		return fmt.Errorf("failed delete cat: %v", err)
	}

	return nil
}
