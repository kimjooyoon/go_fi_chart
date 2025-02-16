package domain

import "time"

// Event 도메인 이벤트 인터페이스입니다.
type Event interface {
	// EventType 이벤트의 타입을 반환합니다.
	EventType() string
	// Timestamp 이벤트가 발생한 시간을 반환합니다.
	Timestamp() time.Time
	// Payload 이벤트의 페이로드를 반환합니다.
	Payload() interface{}
	// Source 이벤트의 발생 소스를 반환합니다.
	Source() string
	// Metadata 이벤트의 메타데이터를 반환합니다.
	Metadata() map[string]string
}

// BaseEvent 기본 이벤트 구현체입니다.
type BaseEvent struct {
	Type string
	Time time.Time
	Data interface{}
	Src  string
	Meta map[string]string
}

func (e *BaseEvent) EventType() string {
	return e.Type
}

func (e *BaseEvent) Timestamp() time.Time {
	return e.Time
}

func (e *BaseEvent) Payload() interface{} {
	return e.Data
}

func (e *BaseEvent) Source() string {
	return e.Src
}

func (e *BaseEvent) Metadata() map[string]string {
	return e.Meta
}
