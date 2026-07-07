package resolver

import (
	"encoding/json" // JSON <-> Go 구조체 변환
	"errors"        // 에러 종류 비교(errors.Is)에 사용
	"net/http"      // HTTP 서버 관련 타입들
	"strconv"       // URL 경로의 문자열 id를 숫자(int64)로 변환

	"player-api/internal/repository"
	"player-api/internal/service"
)

// PlayerHandler는 가장 바깥 계층입니다.
// HTTP 요청을 해석하고, Service를 호출하고, 결과를 JSON으로 내려주는 일만 합니다.
type PlayerHandler struct {
	svc *service.PlayerService // 아래 계층인 Service를 들고 있음
}

// NewPlayerHandler는 svc를 받아 PlayerHandler를 만들어 반환합니다.
func NewPlayerHandler(svc *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

// RegisterRoutes는 "어떤 HTTP 메서드 + 어떤 경로"가 오면 어떤 함수를 실행할지 등록합니다.
// (Go 1.22+ 부터 "GET /players" 처럼 메서드까지 한 번에 지정할 수 있습니다.)
func (h *PlayerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /players", h.list)           // 전체 목록
	mux.HandleFunc("POST /players", h.create)        // 새로 만들기
	mux.HandleFunc("GET /players/{id}", h.get)       // 하나 조회
	mux.HandleFunc("PUT /players/{id}", h.update)    // 수정
	mux.HandleFunc("DELETE /players/{id}", h.delete) // 삭제
}

// ---- 요청 본문(JSON) 모양 ----
// 클라이언트가 보내는 JSON을 담을 구조체입니다.

type createRequest struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

type updateRequest struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

// ---- 실제 요청을 처리하는 핸들러들 ----

// list: GET /players
func (h *PlayerHandler) list(w http.ResponseWriter, r *http.Request) {
	players, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list players")
		return
	}
	writeJSON(w, http.StatusOK, players) // 200 OK + 목록
}

// get: GET /players/{id}
func (h *PlayerHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return // parseID 안에서 이미 에러 응답을 보냈음
	}
	player, err := h.svc.Get(r.Context(), id)
	if handleServiceError(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, player)
}

// create: POST /players
func (h *PlayerHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	// 요청 본문(JSON)을 req 구조체로 풀어냅니다.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	player, err := h.svc.Create(r.Context(), req.Name, req.Age)
	if handleServiceError(w, err) {
		return
	}
	writeJSON(w, http.StatusCreated, player) // 201 Created
}

// update: PUT /players/{id}
func (h *PlayerHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	player, err := h.svc.Update(r.Context(), id, req.Name, req.Age)
	if handleServiceError(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, player)
}

// delete: DELETE /players/{id}
func (h *PlayerHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	err := h.svc.Delete(r.Context(), id)
	if handleServiceError(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204: 본문 없이 성공
}

// ---- 공통 헬퍼 함수들 ----

// parseID는 URL 경로의 {id}(문자열)를 숫자(int64)로 바꿉니다.
// 실패하면 400 응답을 보내고 false를 돌려줍니다.
func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := r.PathValue("id") // 경로에서 {id} 부분 꺼내기
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

// handleServiceError는 아래 계층에서 올라온 에러를 알맞은 HTTP 상태 코드로 바꿉니다.
// 에러를 처리했으면 true를 반환합니다(호출한 쪽은 그대로 return하면 됨).
func handleServiceError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false // 에러 없음 → 아무것도 안 함
	}
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "player not found") // 404
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusBadRequest, "invalid name or age") // 400
	default:
		writeError(w, http.StatusInternalServerError, "internal error") // 500
	}
	return true
}

// writeJSON은 데이터를 JSON으로 바꿔 응답으로 내려줍니다.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError는 {"error": "메시지"} 형태의 JSON 에러 응답을 내려줍니다.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
