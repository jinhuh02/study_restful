package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // database/sql에 MySQL 드라이버 등록

	"player-api/internal/repository"
	"player-api/internal/resolver"
	"player-api/internal/service"
)

func main() {
	// 1) 환경변수에서 DB 접속 정보 읽기 (docker-compose에서 주입)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		env("DB_USER", "myuser"),
		env("DB_PASSWORD", "userpassword"),
		env("DB_HOST", "mysql"),
		env("DB_PORT", "3306"),
		env("DB_NAME", "testdb"),
	)

	// 2) DB 연결 (MySQL이 아직 준비 안 됐을 수 있으니 재시도)
	db, err := connectWithRetry(dsn, 10, 3*time.Second)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("connected to MySQL")

	// 3) 계층 조립: Repository → Service → resolver
	repo := repository.NewPlayerRepository(db)
	svc := service.NewPlayerService(repo)
	h := resolver.NewPlayerHandler(svc)

	// 4) 라우터 구성
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// 5) 서버 실행
	addr := ":" + env("PORT", "8080")
	log.Printf("server listening on %s", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// connectWithRetry는 DB가 응답할 때까지 일정 횟수 재시도합니다.
// 컨테이너 기동 시 MySQL이 아직 준비되지 않은 경우를 대비한 것입니다.
func connectWithRetry(dsn string, attempts int, delay time.Duration) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn) // 여기서는 실제 연결이 아니라 핸들만 생성됨
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	for i := 1; i <= attempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = db.PingContext(ctx) // 실제 연결 확인
		cancel()
		if err == nil {
			return db, nil
		}
		log.Printf("waiting for database... (%d/%d): %v", i, attempts, err)
		time.Sleep(delay)
	}
	return nil, err
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
