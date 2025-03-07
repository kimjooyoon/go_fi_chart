package errors

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// ErrorResponse는 API 응답에 포함될 에러 정보 구조체입니다.
type ErrorResponse struct {
	Code    string `json:"code"`    // 에러 코드
	Message string `json:"message"` // 에러 메시지
}

// RespondWithError는 에러를 HTTP 응답으로 변환하여 응답합니다.
func RespondWithError(w http.ResponseWriter, err error) {
	var statusCode int
	var errorCode string
	var message string

	// 도메인 에러 타입 확인
	var domainErr DomainError
	if errors.As(err, &domainErr) {
		statusCode = domainErr.StatusCode()
		errorCode = domainErr.Code()
		message = domainErr.Error()
	} else {
		// 기본 에러 처리
		statusCode = http.StatusInternalServerError
		errorCode = "INTERNAL_SERVER_ERROR"
		message = "서버 내부 오류가 발생했습니다"
	}

	RespondWithErrorCode(w, statusCode, errorCode, message)
}

// RespondWithErrorCode는 상태 코드, 에러 코드, 메시지를 이용하여 HTTP 에러 응답을 생성합니다.
func RespondWithErrorCode(w http.ResponseWriter, statusCode int, code, message string) {
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("에러 응답 인코딩 실패: %v", err)
	}
}
