package repository

import (
	"context"      // 요청 취소/타임아웃 정보를 담아 DB 호출에 넘기는 용도
	"database/sql" // Go 표준 DB 라이브러리
	"errors"       // 에러를 만들고 비교할 때 사용

	"player-api/internal/model"
)

var ErrNotFound = errors.New("player not found")

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Create(ctx context.Context, name string, age int64) (*model.Player, error) {
	// 1) INSERT 실행. ?는 값이 들어갈 자리(플레이스홀더). SQL 인젝션을 막아줍니다.
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO players (name, age) VALUES (?, ?)`, name, age)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.ReadByID(ctx, id)
}

func (r *PlayerRepository) ReadAll(ctx context.Context) ([]model.Player, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, age FROM players ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]model.Player, 0)

	for rows.Next() {
		var p model.Player
		if err := rows.Scan(&p.ID, &p.Name, &p.Age); err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return players, rows.Err()
}

func (r *PlayerRepository) ReadByID(ctx context.Context, id int64) (*model.Player, error) {
	var p model.Player

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, age FROM players WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Age)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PlayerRepository) Update(ctx context.Context, id int64, name string, age int64) (*model.Player, error) {
	result, err := r.db.ExecContext(ctx,
		`UPDATE players SET name = ?, age = ? WHERE id = ?`, name, age, id)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrNotFound
	}

	return r.ReadByID(ctx, id)
}

func (r *PlayerRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM players WHERE id = ?`, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
