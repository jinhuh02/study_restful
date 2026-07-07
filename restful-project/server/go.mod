// go mod init github.com/user/myproject 명령어로 자동 생성
module player-api // 이 프로젝트의 고유 이름

go 1.22 // go 버전. 

// require : 외부에서 가져다 쓰는 라이브러리 목록. 자동으로 채워짐
// github.com/go-sql-driver/mysql <- Go에서 MySQL 데이터베이스에 접속하게 해주는 드라이버(driver) 라이브러리
require github.com/go-sql-driver/mysql v1.7.1 