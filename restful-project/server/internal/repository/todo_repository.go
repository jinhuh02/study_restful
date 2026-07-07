package repository

import (
	"context"
	"database/sql"
	"errors"

	"todo-api/internal/model"
)

// ErrNotFound는 해당 id의 Todo가 없을 때 반환합니다.
// 상위 계층(Service/Handler)이 이 에러를 보고 404를 내려줄 수 있게 합니다.
var ErrNotFound = errors.New("todo not found")

// TodoRepository는 DB 접근만 담당하는 계층입니다.
// 여기서는 SQL을 실행하는 것 외의 판단(비즈니스 규칙)은 하지 않습니다.
type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// FindAll: 모든 Todo를 최신순으로 조회합니다.
func (r *TodoRepository) FindAll(ctx context.Context) ([]model.Todo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, completed, created_at FROM todos ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// nil이 아닌 빈 슬라이스로 시작 → JSON에서 null 대신 [] 로 나오게 함
	todos := make([]model.Todo, 0)
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

// FindByID: id로 Todo 하나를 조회합니다. 없으면 ErrNotFound.
func (r *TodoRepository) FindByID(ctx context.Context, id int64) (*model.Todo, error) {
	var t model.Todo
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, completed, created_at FROM todos WHERE id = ?`, id).
		Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Create: 새 Todo를 삽입하고, 방금 만든 레코드를 반환합니다.
func (r *TodoRepository) Create(ctx context.Context, title string) (*model.Todo, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO todos (title, completed) VALUES (?, ?)`, title, false)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	// 삽입 직후 created_at 같은 DB 기본값까지 채워서 돌려주려고 다시 조회
	return r.FindByID(ctx, id)
}

// Update: title / completed를 수정합니다. 대상이 없으면 ErrNotFound.
func (r *TodoRepository) Update(ctx context.Context, id int64, title string, completed bool) (*model.Todo, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE todos SET title = ?, completed = ? WHERE id = ?`, title, completed, id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrNotFound
	}
	return r.FindByID(ctx, id)
}

// Delete: id로 삭제합니다. 대상이 없으면 ErrNotFound.
func (r *TodoRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
