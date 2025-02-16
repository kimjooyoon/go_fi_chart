package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestServerStartupAndShutdown(t *testing.T) {
	// 테스트용 포트 설정
	os.Setenv("PORT", "8081")
	defer os.Unsetenv("PORT")

	// 서버 시작을 위한 채널
	serverReady := make(chan struct{})
	serverClosed := make(chan struct{})

	// 서버 실행
	go func() {
		// 라우터 설정
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))

		// 헬스 체크 엔드포인트 추가
		r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// 서버 설정
		server := &http.Server{
			Addr:              ":8081",
			Handler:           r,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// 서버가 준비되었음을 알림
		close(serverReady)

		// 서버 종료 시그널 처리
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			server.Shutdown(shutdownCtx)
			close(serverClosed)
		}()

		// 서버 시작
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("서버 시작 실패: %v", err)
		}
	}()

	// 서버가 준비될 때까지 대기
	select {
	case <-serverReady:
		// 서버가 준비됨
	case <-time.After(5 * time.Second):
		t.Fatal("서버 시작 타임아웃")
	}

	// 서버가 응답하는지 확인
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get("http://localhost:8081/health")
	if err != nil {
		t.Fatalf("서버 응답 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("예상치 못한 상태 코드: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	// 서버에 종료 시그널 전송
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("프로세스 찾기 실패: %v", err)
	}
	process.Signal(syscall.SIGTERM)

	// 서버가 정상적으로 종료되었는지 확인
	select {
	case <-serverClosed:
		// 서버가 정상적으로 종료됨
	case <-time.After(5 * time.Second):
		t.Fatal("서버 종료 타임아웃")
	}
}
