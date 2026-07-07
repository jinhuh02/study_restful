FROM golang:1.22-alpine

WORKDIR /app

# 프로젝트 파일 전부 복사
COPY . .

# 의존성 다운로드
RUN go mod download

# 빌드해서 server라는 실행 파일 생성
RUN go build -o server .
# .(현재 폴더의 소스코드) 를
# -o(output) server라는 이름으로 빌드해라

EXPOSE 8080

CMD ["/app/server"]
