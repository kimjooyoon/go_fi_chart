package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"github.com/aske/go_fi_chart/internal/common/repository"
	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
	assetRepo "github.com/aske/go_fi_chart/services/asset/internal/domain/asset"
)

// AssetRepository MongoDB 기반 자산 저장소 구현체입니다.
type AssetRepository struct {
	collection *mongo.Collection
	eventBus   events.EventBus
}

// 자산 문서 구조체
type assetDocument struct {
	ID        string            `bson:"_id"`
	UserID    string            `bson:"user_id"`
	Type      domain.AssetType  `bson:"type"`
	Name      string            `bson:"name"`
	Amount    amountDocument    `bson:"amount"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
	IsDeleted bool              `bson:"is_deleted"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty"`
}

type amountDocument struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}

// NewAssetRepository는 새로운 MongoDB 자산 저장소를 생성합니다.
func NewAssetRepository(db *mongo.Database, eventBus events.EventBus) *AssetRepository {
	return &AssetRepository{
		collection: db.Collection("assets"),
		eventBus:   eventBus,
	}
}

// Init MongoDB 컬렉션에 필요한 인덱스를 초기화합니다.
func (r *AssetRepository) Init(ctx context.Context) error {
	// 사용자 ID에 대한 인덱스 생성
	userIDIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetBackground(true),
	}
	
	// 자산 유형에 대한 인덱스 생성
	assetTypeIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "type", Value: 1}},
		Options: options.Index().SetBackground(true),
	}
	
	// 삭제 상태에 대한 인덱스 생성
	isDeletedIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "is_deleted", Value: 1}},
		Options: options.Index().SetBackground(true),
	}
	
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		userIDIndex,
		assetTypeIndex,
		isDeletedIndex,
	})
	
	return err
}

// toDocument는 Asset 엔티티를 MongoDB 문서로 변환합니다.
func toDocument(asset *domain.Asset) assetDocument {
	doc := assetDocument{
		ID:        asset.ID,
		UserID:    asset.UserID,
		Type:      asset.Type,
		Name:      asset.Name,
		Amount: amountDocument{
			Amount:   asset.Amount.Amount,
			Currency: asset.Amount.Currency,
		},
		CreatedAt: primitive.NewDateTimeFromTime(asset.CreatedAt),
		UpdatedAt: primitive.NewDateTimeFromTime(asset.UpdatedAt),
		IsDeleted: asset.IsDeleted,
	}
	
	if asset.DeletedAt != nil {
		deletedAt := primitive.NewDateTimeFromTime(*asset.DeletedAt)
		doc.DeletedAt = &deletedAt
	}
	
	return doc
}

// fromDocument는 MongoDB 문서를 Asset 엔티티로 변환합니다.
func fromDocument(doc assetDocument) (*domain.Asset, error) {
	amount, err := valueobjects.NewMoney(doc.Amount.Amount, doc.Amount.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid money value: %w", err)
	}
	
	asset := &domain.Asset{
		ID:        doc.ID,
		UserID:    doc.UserID,
		Type:      doc.Type,
		Name:      doc.Name,
		Amount:    amount,
		CreatedAt: doc.CreatedAt.Time(),
		UpdatedAt: doc.UpdatedAt.Time(),
		IsDeleted: doc.IsDeleted,
	}
	
	if doc.DeletedAt != nil {
		deletedAt := doc.DeletedAt.Time()
		asset.DeletedAt = &deletedAt
	}
	
	return asset, nil
}

// FindByID는 ID로 자산을 조회합니다.
func (r *AssetRepository) FindByID(ctx context.Context, id string) (*domain.Asset, error) {
	var doc assetDocument
	
	filter := bson.M{
		"_id":        id,
		"is_deleted": false,
	}
	
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.NewError(
				"FindByID",
				"Asset",
				fmt.Sprintf("asset not found with id: %s", id),
				repository.ErrEntityNotFound,
			)
		}
		return nil, repository.NewError(
			"FindByID",
			"Asset",
			fmt.Sprintf("failed to find asset: %v", err),
			repository.ErrRepositoryError,
		)
	}
	
	return fromDocument(doc)
}

// FindAll은 모든 자산을 조회합니다. 옵션을 통해 필터링, 정렬, 페이지네이션을 적용할 수 있습니다.
func (r *AssetRepository) FindAll(ctx context.Context, opts ...repository.FindOption) ([]*domain.Asset, error) {
	options := repository.NewFindOptions()
	for _, opt := range opts {
		opt.Apply(options)
	}
	
	// 기본 필터: 삭제되지 않은 자산만 조회
	filter := bson.M{"is_deleted": false}
	
	// 사용자 지정 필터 적용
	for field, value := range options.Filters {
		filter[field] = value
	}
	
	// MongoDB 옵션 설정
	findOptions := options.ToMongoOptions()
	
	// 쿼리 실행
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, repository.NewError(
			"FindAll",
			"Asset",
			fmt.Sprintf("failed to find assets: %v", err),
			repository.ErrRepositoryError,
		)
	}
	defer cursor.Close(ctx)
	
	// 결과 처리
	var assets []*domain.Asset
	for cursor.Next(ctx) {
		var doc assetDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, repository.NewError(
				"FindAll",
				"Asset",
				fmt.Sprintf("failed to decode asset: %v", err),
				repository.ErrRepositoryError,
			)
		}
		
		asset, err := fromDocument(doc)
		if err != nil {
			return nil, repository.NewError(
				"FindAll",
				"Asset",
				fmt.Sprintf("failed to convert document to asset: %v", err),
				repository.ErrRepositoryError,
			)
		}
		
		assets = append(assets, asset)
	}
	
	if err := cursor.Err(); err != nil {
		return nil, repository.NewError(
			"FindAll",
			"Asset",
			fmt.Sprintf("cursor error: %v", err),
			repository.ErrRepositoryError,
		)
	}
	
	return assets, nil
}

// Save는 자산을 저장합니다.
func (r *AssetRepository) Save(ctx context.Context, asset *domain.Asset) error {
	// 이미 존재하는지 확인
	existingAsset, err := r.FindByID(ctx, asset.ID)
	if err == nil && existingAsset != nil {
		return repository.NewError(
			"Save",
			"Asset",
			fmt.Sprintf("asset already exists with id: %s", asset.ID),
			repository.ErrDuplicateEntity,
		)
	} else if err != nil && !errors.Is(err, repository.ErrEntityNotFound) {
		return err
	}
	
	// 문서로 변환
	doc := toDocument(asset)
	
	// 저장
	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		return repository.NewError(
			"Save",
			"Asset",
			fmt.Sprintf("failed to save asset: %v", err),
			repository.ErrRepositoryError,
		)
	}
	
	// 이벤트 발행
	events := asset.Events()
	for _, event := range events {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return repository.NewError(
				"Save",
				"Asset",
				fmt.Sprintf("failed to publish event: %v", err),
				repository.ErrRepositoryError,
			)
		}
	}
	asset.ClearEvents()
	
	return nil
}

// Update는 자산을 업데이트합니다.
func (r *AssetRepository) Update(ctx context.Context, asset *domain.Asset) error {
	// 존재하는지 확인
	_, err := r.FindByID(ctx, asset.ID)
	if err != nil {
		return err
	}
	
	// 문서로 변환
	doc := toDocument(asset)
	
	// 업데이트
	filter := bson.M{"_id": asset.ID}
	update := bson.M{"$set": doc}
	
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return repository.NewError(
			"Update",
			"Asset",
			fmt.Sprintf("failed to update asset: %v", err),
			repository.ErrRepositoryError,
		)
	}
	
	// 이벤트 발행
	events := asset.Events()
	for _, event := range events {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return repository.NewError(
				"Update",
				"Asset",
				fmt.Sprintf("failed to publish event: %v", err),
				repository.ErrRepositoryError,
			)
		}
	}
	asset.ClearEvents()
	
	return nil
}

// Delete는 자산을 삭제합니다.
func (r *AssetRepository) Delete(ctx context.Context, id string) error {
	// 존재하는지 확인
	asset, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 논리적 삭제 처리
	asset.MarkAsDeleted()
	
	// 문서로 변환
	doc := toDocument(asset)
	
	// 업데이트
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": doc.DeletedAt,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return repository.NewError(
			"Delete",
			"Asset",
			fmt.Sprintf("failed to delete asset: %v", err),
			repository.ErrRepositoryError,
		)
	}
	
	// 이벤트 발행
	events := asset.Events()
	for _, event := range events {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return repository.NewError(
				"Delete",
				"Asset",
				fmt.Sprintf("failed to publish event: %v", err),
				repository.ErrRepositoryError,
			)
		}
	}
	
	return nil
}

// Count는 조건에 맞는 자산의 총 개수를 반환합니다.
func (r *AssetRepository) Count(ctx context.Context, opts ...repository.FindOption) (int64, error) {
	options := repository.NewFindOptions()
	for _, opt := range opts {
		opt.Apply(options)
	}
	
	// 기본 필터: 삭제되지 않은 자산만 계산
	filter := bson.M{"is_deleted": false}
	
	// 사용자 지정 필터 적용
	for field, value := range options.Filters {
		filter[field] = value
	}
	
	// 개수 조회
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, repository.NewError(
			"Count",
			"Asset",
			fmt.Sprintf("failed to count assets: %v", err),
			repository.ErrRepositoryError,
		)
	}
	
	return count, nil
}

// FindByUserID는 사용자 ID로 자산 목록을 조회합니다.
func (r *AssetRepository) FindByUserID(ctx context.Context, userID string, opts ...repository.FindOption) ([]*domain.Asset, error) {
	// 필터에 사용자 ID 추가
	filterOpt := repository.WithFilter("user_id", userID)
	
	// 기존 옵션에 사용자 ID 필터 추가
	opts = append(opts, filterOpt)
	
	// FindAll 메서드 호출
	return r.FindAll(ctx, opts...)
}

// FindByType은 자산 유형으로 자산 목록을 조회합니다.
func (r *AssetRepository) FindByType(ctx context.Context, assetType domain.AssetType, opts ...repository.FindOption) ([]*domain.Asset, error) {
	// 필터에 자산 유형 추가
	filterOpt := repository.WithFilter("type", assetType)
	
	// 기존 옵션에 자산 유형 필터 추가
	opts = append(opts, filterOpt)
	
	// FindAll 메서드 호출
	return r.FindAll(ctx, opts...)
}

// CountByUserID는 사용자 ID에 해당하는 자산 개수를 반환합니다.
func (r *AssetRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	// 필터에 사용자 ID 추가
	filterOpt := repository.WithFilter("user_id", userID)
	
	// Count 메서드 호출
	return r.Count(ctx, filterOpt)
}

// CountByType은 자산 유형에 해당하는 자산 개수를 반환합니다.
func (r *AssetRepository) CountByType(ctx context.Context, assetType domain.AssetType) (int64, error) {
	// 필터에 자산 유형 추가
	filterOpt := repository.WithFilter("type", assetType)
	
	// Count 메서드 호출
	return r.Count(ctx, filterOpt)
}