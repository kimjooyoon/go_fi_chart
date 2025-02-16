package event

import "errors"

// 이벤트 시스템 관련 에러
var (
	// ErrEventBusClosed 이벤트 버스가 닫혔을 때 발생하는 에러
	ErrEventBusClosed = errors.New("event bus is closed")
	// ErrInvalidEventType 유효하지 않은 이벤트 타입
	ErrInvalidEventType = errors.New("invalid event type")
	// ErrInvalidPayload 유효하지 않은 페이로드
	ErrInvalidPayload = errors.New("invalid payload")
	// ErrHandlerNotFound 핸들러를 찾을 수 없음
	ErrHandlerNotFound = errors.New("handler not found")
	// ErrDuplicateHandler 이미 등록된 핸들러
	ErrDuplicateHandler = errors.New("duplicate handler")
)
