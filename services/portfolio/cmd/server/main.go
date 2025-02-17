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

	"log/slog"

	"github.com/aske/go_fi_chart/services/portfolio/internal/api"
	"github.com/aske/go_fi_chart/services/portfolio/internal/domain"
	"github.com/gorilla/mux"
)

func main() {
	serverCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 환경 변수 로드
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 라우터 설정
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(recoveryMiddleware)
	r.Use(timeoutMiddleware)

	// 핸들러 설정
	portfolioRepo := domain.NewMemoryPortfolioRepository()
	handler := api.NewHandler(portfolioRepo, logger)
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
	logger.Info("starting server", "port", fmt.Sprintf(":%s", port))
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// 서버 종료 대기
	<-serverCtx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("서버 종료 중 오류 발생: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("패닉 발생: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func timeoutMiddleware(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 60*time.Second, "시간 초과")
}
