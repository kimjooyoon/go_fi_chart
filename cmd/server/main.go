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

	"github.com/aske/go_fi_chart/internal/api"
	"github.com/aske/go_fi_chart/internal/config"
	"github.com/aske/go_fi_chart/internal/domain/asset"
	"github.com/aske/go_fi_chart/internal/domain/gamification"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("FIN-RPG 서버 시작 중...")

	// 설정 로드
	cfg := config.NewDefaultConfig()

	// 저장소 생성
	assetRepo := asset.NewMemoryAssetRepository()
	transactionRepo := asset.NewMemoryTransactionRepository()
	portfolioRepo := asset.NewMemoryPortfolioRepository()
	gamificationRepo := gamification.NewMemoryRepository()

	// API 핸들러 생성
	apiHandler := api.NewHandler(assetRepo, transactionRepo, portfolioRepo, gamificationRepo)

	// 라우터 설정
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// 헬스 체크
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("헬스 체크 응답 작성 실패: %v", err)
		}
	})

	// API 라우터 그룹
	r.Route("/api", func(r chi.Router) {
		apiHandler.RegisterRoutes(r)
	})

	// 서버 설정
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:           r,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		ReadHeaderTimeout: 2 * time.Second,
	}

	// 종료 시그널 처리
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("서버가 시작되었습니다. %s에서 대기 중...", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("서버 시작 실패: %v", err)
		}
	}()

	<-done
	log.Println("서버 종료 중...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("서버 종료 실패: %v", err)
	}

	log.Println("서버가 정상적으로 종료되었습니다.")
}
