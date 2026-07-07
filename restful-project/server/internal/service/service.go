package service

import (
	"context" // 요청 취소/타임아웃 정보를 아래 계층까지 그대로 전달하는 용도
	"errors"  // 에러를 만들 때 사용
	"strings" // 문자열 앞뒤 공백 제거(TrimSpace)에 사용

	"player-api/internal/model"
	"player-api/internal/repository"
)

// ErrValidation은 입력값이 규칙에 맞지 않을 때 반환하는 에러입니다.
// 예: 이름이 비어 있거나 나이가 음수일 때.
// 상위 계층(Resolver)이 이 에러를 보고 "400 Bad Request"를 내려줄 수 있습니다.
var ErrValidation = errors.New("validation failed")

// PlayerService는 "비즈니스 규칙"을 담당하는 중간 계층입니다.
// 여기서 입력값이 올바른지 판단하고, 실제 DB 작업은 Repository에게 시킵니다.
type PlayerService struct {
	repo *repository.PlayerRepository // 아래 계층인 Repository를 들고 있음
}

// NewPlayerService는 repo를 받아 PlayerService를 만들어 반환합니다.
func NewPlayerService(repo *repository.PlayerRepository) *PlayerService {
	return &PlayerService{repo: repo}
}

// List: 전체 목록 조회. 검증할 게 없으니 그대로 Repository에 넘깁니다.
func (s *PlayerService) List(ctx context.Context) ([]model.Player, error) {
	return s.repo.ReadAll(ctx)
}

// Get: id로 한 명 조회. 이것도 그대로 위임합니다.
func (s *PlayerService) Get(ctx context.Context, id int64) (*model.Player, error) {
	return s.repo.ReadByID(ctx, id)
}

// Create: 새 Player 생성. 저장하기 전에 입력값을 검사합니다.
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

// Update: 기존 Player 수정. Create와 같은 검증을 한 뒤 위임합니다.
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

// Delete: id로 삭제. 검증할 게 없어 그대로 위임합니다.
func (s *PlayerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
