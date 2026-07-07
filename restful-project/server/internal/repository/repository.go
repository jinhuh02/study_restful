package repository

import (
	"context"      // 요청 취소/타임아웃 정보를 담아 DB 호출에 넘기는 용도
	"database/sql" // Go 표준 DB 라이브러리
	"errors"       // 에러를 만들고 비교할 때 사용

	"player-api/internal/model"
)

// ErrNotFound는 해당 id의 Player가 DB에 없을 때 돌려주는 에러입니다.
// 상위 계층(Service/Resolver)이 이 에러를 보고 "404 Not Found"를 내려줄 수 있습니다.
var ErrNotFound = errors.New("player not found")

// PlayerRepository는 DB 접근(SQL 실행)하여 "저장/조회/수정/삭제"만 합니다.
type PlayerRepository struct {
	db *sql.DB // 실제 DB 연결 핸들. 바깥(main.go)에서 만들어 넣어줍니다.
}

// NewPlayerRepository는 db 핸들을 받아 PlayerRepository를 만들어 반환합니다.
// Go에는 생성자 문법이 따로 없어서, 이렇게 New○○ 함수를 관례로 씁니다.
func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Create(ctx context.Context, name string, age int64) (*model.Player, error) {
	// 1) INSERT 실행. ?는 값이 들어갈 자리(플레이스홀더). SQL 인젝션을 막아줍니다.
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO players (name, age) VALUES (?, ?)`, name, age)
	if err != nil {
		return nil, err // DB 오류가 나면 그대로 위로 전달
	}

	// 2) 방금 INSERT된 행의 자동 증가 id(AUTO_INCREMENT)를 가져옵니다.
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 3) 그 id로 다시 조회해서, 완전한 Player 객체를 돌려줍니다.
	return r.ReadByID(ctx, id)
}

// ReadAll: 모든 Player를 id 내림차순(최신순)으로 조회합니다.
func (r *PlayerRepository) ReadAll(ctx context.Context) ([]model.Player, error) {
	// 여러 행을 반환하는 조회는 QueryContext를 씁니다.
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, age FROM players ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 함수가 끝날 때 결과 커서를 꼭 닫아줍니다(자원 누수 방지).

	// nil이 아닌 "빈 슬라이스"로 시작 → 결과가 없어도 JSON에서 null이 아니라 [] 로 나옵니다.
	players := make([]model.Player, 0)

	// rows.Next()로 한 줄씩 앞으로 이동하며 읽습니다. 더 없으면 false.
	for rows.Next() {
		var p model.Player
		// Scan: SELECT한 컬럼 순서(id, name, age) 그대로 변수 주소에 담습니다.
		if err := rows.Scan(&p.ID, &p.Name, &p.Age); err != nil {
			return nil, err
		}
		players = append(players, p) // 슬라이스 뒤에 추가
	}

	// 반복 도중 발생한 오류가 있었는지 마지막에 확인합니다.
	return players, rows.Err()
}

// ReadByID: id 하나로 Player를 조회합니다. 없으면 ErrNotFound를 반환합니다.
func (r *PlayerRepository) ReadByID(ctx context.Context, id int64) (*model.Player, error) {
	var p model.Player

	// 딱 한 줄만 기대할 때는 QueryRowContext를 씁니다.
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, age FROM players WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Age)

	// 조회 결과가 0줄이면 sql.ErrNoRows가 옵니다 → 우리 에러로 바꿔서 반환.
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil // 주소(&p)를 반환해서 *model.Player 타입으로 돌려줍니다.
}

// Update: id에 해당하는 Player의 name, age를 수정합니다. 대상이 없으면 ErrNotFound.
func (r *PlayerRepository) Update(ctx context.Context, id int64, name string, age int64) (*model.Player, error) {
	result, err := r.db.ExecContext(ctx,
		`UPDATE players SET name = ?, age = ? WHERE id = ?`, name, age, id)
	if err != nil {
		return nil, err
	}

	// RowsAffected: 실제로 수정된 행 수. 0이면 그 id가 없었다는 뜻.
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrNotFound
	}

	// 수정 후 최신 상태를 다시 조회해서 반환합니다.
	return r.ReadByID(ctx, id)
}

// Delete: id로 Player를 삭제합니다. 대상이 없으면 ErrNotFound.
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
	return nil // 성공 시 에러 없음(nil)을 반환
}
