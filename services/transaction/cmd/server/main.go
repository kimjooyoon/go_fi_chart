package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aske/go_fi_chart/services/transaction/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	serverCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 환경 변수 로드
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 라우터 설정
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// 핸들러 설정
	handler := api.NewHandler()
	handler.RegisterRoutes(r)

	// 서버 종료 시그널 처리
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	// 서버 설정
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// 서버 시작
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("서버 시작 실패: %v", err)
		}
	}()

	// 서버 종료 대기
	<-serverCtx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel() // shutdown context의 cancel 함수를 호출합니다.

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("서버 종료 중 오류 발생: %v", err)
	}
}
