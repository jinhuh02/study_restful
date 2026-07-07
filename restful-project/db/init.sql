-- init.sql은 MySQL 볼륨이 처음 만들어질때 한번만 실행됨
-- 이전에 docker compose up을 한 적이 있는데, 테이블 수정사항을 업데이트 하고 싶으면
-- docker compose down -v로 볼륨을 지우고 다시 docker compose up -d --build 하여 init.sql을 실행시켜야 함
-- (docker-compose에서 /docker-entrypoint-initdb.d/ 에 마운트하고 있음)

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS players (
    id   INT AUTO_INCREMENT PRIMARY KEY,     -- 자동 증가 기본키
    name VARCHAR(255) NOT NULL,              -- 이름 (필수)
    age  INT          NOT NULL               -- 나이 (필수)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 확인용 샘플 데이터 몇 개
INSERT INTO players (name, age) VALUES
    ('허진', 300),
    ('홍길동', 25),
    ('김철수', 40);
