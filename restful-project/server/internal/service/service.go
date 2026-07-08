package service

import (
	"context" // 요청 취소/타임아웃 정보를 아래 계층까지 그대로 전달하는 용도
	"errors"  // 에러를 만들 때 사용
	"strings" // 문자열 앞뒤 공백 제거(TrimSpace)에 사용

	"player-api/internal/model"
	"player-api/internal/repository"
)

var ErrValidation = errors.New("validation failed")

type PlayerService struct {
	repo *repository.PlayerRepository // 아래 계층인 Repository를 들고 있음
}

func NewPlayerService(repo *repository.PlayerRepository) *PlayerService {
	return &PlayerService{repo: repo}
}

// R1 검증할 게 없으니 그대로 Repository에 넘깁니다.
func (s *PlayerService) List(ctx context.Context) ([]model.Player, error) {
	return s.repo.ReadAll(ctx)
}

// R2 이것도 그대로
func (s *PlayerService) Get(ctx context.Context, id int64) (*model.Player, error) {
	return s.repo.ReadByID(ctx, id)
}

// C 저장하기 전에 입력값을 검사합니다.
func (s *PlayerService) Create(ctx context.Context, name string, age int64) (*model.Player, error) {
	name = strings.TrimSpace(name) // 이름 앞뒤 공백 제거 ("  " 같은 입력 방지)
	if name == "" {                // 이름이 비었으면 규칙 위반
		return nil, ErrValidation
	}
	if age < 0 { // 나이는 음수가 될 수 없음
		return nil, ErrValidation
	}
	return s.repo.Create(ctx, name, age) // 통과하면 저장 위임
}

// U 입력값 검사
func (s *PlayerService) Update(ctx context.Context, id int64, name string, age int64) (*model.Player, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrValidation
	}
	if age < 0 {
		return nil, ErrValidation
	}
	return s.repo.Update(ctx, id, name, age)
}

// D 그대로
func (s *PlayerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
