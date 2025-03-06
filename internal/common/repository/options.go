package repository

// FindOption 조회 옵션 인터페이스
type FindOption interface {
	Apply(options *FindOptions)
}

// FindOptions 조회 옵션 집합
type FindOptions struct {
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  SortOrder
	Filters    []Filter
	Pagination *Pagination
}

// SortOrder 정렬 순서 타입
type SortOrder string

const (
	// SortAscending 오름차순 정렬
	SortAscending SortOrder = "ASC"
	// SortDescending 내림차순 정렬
	SortDescending SortOrder = "DESC"
)

// Filter 필터링 조건
type Filter struct {
	Field    string      // 필드명
	Operator string      // 연산자 (eq, ne, gt, lt, gte, lte, like 등)
	Value    interface{} // 비교 값
}

// Pagination 페이지네이션 정보
type Pagination struct {
	Page     int // 현재 페이지 (1부터 시작)
	PageSize int // 페이지 크기
}

// NewFindOptions 기본 옵션으로 FindOptions 생성
func NewFindOptions() *FindOptions {
	return &FindOptions{
		Limit:     100,
		Offset:    0,
		SortOrder: SortAscending,
	}
}

// WithLimit 조회 결과 개수 제한 옵션
func WithLimit(limit int) FindOption {
	return limitOption{limit: limit}
}

type limitOption struct {
	limit int
}

func (o limitOption) Apply(options *FindOptions) {
	options.Limit = o.limit
}

// WithOffset 조회 시작 위치 옵션
func WithOffset(offset int) FindOption {
	return offsetOption{offset: offset}
}

type offsetOption struct {
	offset int
}

func (o offsetOption) Apply(options *FindOptions) {
	options.Offset = o.offset
}

// WithSort 정렬 옵션
func WithSort(field string, order SortOrder) FindOption {
	return sortOption{field: field, order: order}
}

type sortOption struct {
	field string
	order SortOrder
}

func (o sortOption) Apply(options *FindOptions) {
	options.SortBy = o.field
	options.SortOrder = o.order
}

// WithFilter 필터링 옵션
func WithFilter(field, operator string, value interface{}) FindOption {
	return filterOption{
		filter: Filter{
			Field:    field,
			Operator: operator,
			Value:    value,
		},
	}
}

type filterOption struct {
	filter Filter
}

func (o filterOption) Apply(options *FindOptions) {
	options.Filters = append(options.Filters, o.filter)
}

// WithPagination 페이지네이션 옵션
func WithPagination(page, pageSize int) FindOption {
	return paginationOption{
		pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
		},
	}
}

type paginationOption struct {
	pagination Pagination
}

func (o paginationOption) Apply(options *FindOptions) {
	options.Pagination = &o.pagination
}
