-- 이 파일은 MySQL 컨테이너가 "처음" 생성될 때 자동 실행됩니다.
-- (docker-compose에서 /docker-entrypoint-initdb.d/ 에 마운트)
-- 이미 volume에 데이터가 있으면 실행되지 않습니다.

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS todos (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    completed  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 확인용 샘플 데이터 몇 개
INSERT INTO todos (title, completed) VALUES
    ('Go REST API 만들기', TRUE),
    ('Docker로 실행하기', FALSE),
    ('CRUD 테스트하기', FALSE);
