package model

import "time"

// Todo는 할 일 항목 하나를 나타내는 도메인 모델입니다.
// 이 구조체가 모든 계층(Handler / Service / Repository)에서 공통으로 오갑니다.
type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}
