# ---- 1단계: 빌드 (Go 이미지에서 컴파일) ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 의존성 먼저 복사 → 캐시 활용 (코드만 바뀌면 이 단계는 재실행 안 됨)
COPY go.mod go.sum ./
RUN go mod download

# 소스 복사 후 정적 바이너리로 빌드
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

# ---- 2단계: 실행 (가벼운 Alpine에 바이너리만 담기) ----
FROM alpine:3.20

WORKDIR /app

# 바이너리만 복사 → 최종 이미지가 매우 작아짐 (수십 MB)
COPY --from=builder /server /app/server

EXPOSE 8080

CMD ["/app/server"]
