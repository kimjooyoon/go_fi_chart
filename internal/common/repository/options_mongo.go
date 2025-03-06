package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ToMongoOptions는 FindOptions를 MongoDB의 FindOptions로 변환합니다.
func (opts *FindOptions) ToMongoOptions() *options.FindOptions {
	mongoOpts := options.Find()

	// 페이지네이션 적용
	if opts.Limit > 0 {
		mongoOpts.SetLimit(int64(opts.Limit))
		mongoOpts.SetSkip(int64(opts.Offset))
	}

	// 정렬 적용
	if len(opts.Sort) > 0 {
		sort := bson.D{}
		for field, order := range opts.Sort {
			// MongoDB 정렬 방식: 1은 오름차순, -1은 내림차순
			var sortOrder int
			if order == SortAscending {
				sortOrder = 1
			} else {
				sortOrder = -1
			}
			sort = append(sort, bson.E{Key: field, Value: sortOrder})
		}
		if len(sort) > 0 {
			mongoOpts.SetSort(sort)
		}
	}

	return mongoOpts
}

// WithMongoFilter는 MongoDB에 특화된 필터를 추가하는 옵션을 생성합니다.
func WithMongoFilter(filter bson.M) FindOption {
	return mongoFilterOption{filter: filter}
}

type mongoFilterOption struct {
	filter bson.M
}

func (o mongoFilterOption) Apply(options *FindOptions) {
	if options.ExtraFilters == nil {
		options.ExtraFilters = make(map[string]interface{})
	}
	options.ExtraFilters["mongodb_filter"] = o.filter
}
