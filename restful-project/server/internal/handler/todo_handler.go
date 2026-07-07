package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"todo-api/internal/repository"
	"todo-api/internal/service"
)

// TodoHandler는 요청/응답(HTTP, JSON)만 다루는 최상단 계층입니다.
// REST에서는 보통 Handler(또는 Controller)라고 부르며,
// 요청을 해석해서 Service를 호출하고, 결과를 JSON으로 내려줍니다.
type TodoHandler struct {
	svc *service.TodoService
}

func NewTodoHandler(svc *service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// RegisterRoutes는 Go 1.22+ ServeMux의 "메서드 + 경로" 패턴으로 라우트를 등록합니다.
func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /todos", h.list)
	mux.HandleFunc("POST /todos", h.create)
	mux.HandleFunc("GET /todos/{id}", h.get)
	mux.HandleFunc("PUT /todos/{id}", h.update)
	mux.HandleFunc("DELETE /todos/{id}", h.delete)
}

// ---- 요청 본문(JSON) 구조 ----

type createRequest struct {
	Title string `json:"title"`
}

type updateRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// ---- 핸들러들 ----

func (h *TodoHandler) list(w http.ResponseWriter, r *http.Request) {
	todos, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list todos")
		return
	}
	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	todo, err := h.svc.Get(r.Context(), id)
	if handleServiceError(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	todo, err := h.svc.Create(r.Context(), req.Title)
	if handleServiceError(w, err) {
		return
	}
	writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	todo, err := h.svc.Update(r.Context(), id, req.Title, req.Completed)
	if handleServiceError(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) delete(w http.ResponseWriter, r *http.Request) {
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

// ---- 공통 헬퍼 ----

// parseID는 URL 경로의 {id}를 int64로 변환합니다.
func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

// handleServiceError는 계층에서 올라온 에러를 적절한 HTTP 상태로 변환합니다.
// 처리했으면 true를 반환합니다(=호출부는 그대로 return하면 됨).
func handleServiceError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "todo not found")
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusBadRequest, "title must not be empty")
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
