package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("FIN-RPG 서버 시작 중...")

	// TODO: 의존성 주입 설정
	// TODO: 라우터 설정
	// TODO: 미들웨어 설정
	// TODO: 데이터베이스 연결

	srv := &http.Server{
		Addr:              ":8080",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("서버가 시작되었습니다. :8080 포트에서 대기 중...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
