package model

// 도메인 모델입니다.
// 모든 계층(resolver / Service / Repository)에서 공통으로 쓰이는 구조체
type Player struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

//`json:"name"` <- Name을 JSON으로 변환 시 name(소문자)로 변환해줌.
//Go는 대문자 규칙이지만 JSON(웹API)쪽은 소문자여야하는 규칙이 있음.
