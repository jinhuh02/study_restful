package service

import (
	"context"
	"errors"
	"strings"

	"todo-api/internal/model"
	"todo-api/internal/repository"
)

// ErrValidation은 입력값이 규칙에 맞지 않을 때 반환합니다.
var ErrValidation = errors.New("validation failed")

// TodoService는 비즈니스 로직을 담당하는 계층입니다.
// "제목은 비어 있으면 안 된다" 같은 규칙을 여기서 판단하고,
// 실제 저장은 Repository에게 위임합니다.
type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) List(ctx context.Context) ([]model.Todo, error) {
	return s.repo.FindAll(ctx)
}

func (s *TodoService) Get(ctx context.Context, id int64) (*model.Todo, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TodoService) Create(ctx context.Context, title string) (*model.Todo, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrValidation
	}
	return s.repo.Create(ctx, title)
}

func (s *TodoService) Update(ctx context.Context, id int64, title string, completed bool) (*model.Todo, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrValidation
	}
	return s.repo.Update(ctx, id, title, completed)
}

func (s *TodoService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
