package domain

import (
	"context"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/google/uuid"
)

// PortfolioAsset 포트폴리오의 자산 구성을 나타냅니다.
type PortfolioAsset struct {
	AssetID string
	Weight  valueobjects.Percentage
}

// Portfolio 포트폴리오를 나타냅니다.
type Portfolio struct {
	ID        string
	UserID    string
	Name      string
	Assets    []PortfolioAsset
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewPortfolio 새로운 포트폴리오를 생성합니다.
func NewPortfolio(userID string, name string) *Portfolio {
	now := time.Now()
	return &Portfolio{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		Assets:    make([]PortfolioAsset, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddAsset 포트폴리오에 자산을 추가합니다.
func (p *Portfolio) AddAsset(assetID string, weight valueobjects.Percentage) error {
	// 기존 자산의 총 가중치 계산
	var totalWeight float64
	for _, asset := range p.Assets {
		totalWeight += asset.Weight.Value
	}

	// 새 자산 추가 시 총 가중치가 100%를 초과하는지 확인
	if totalWeight+weight.Value > 100 {
		return ErrInvalidWeight
	}

	p.Assets = append(p.Assets, PortfolioAsset{
		AssetID: assetID,
		Weight:  weight,
	})
	p.UpdatedAt = time.Now()
	return nil
}

// UpdateAssetWeight 자산의 가중치를 업데이트합니다.
func (p *Portfolio) UpdateAssetWeight(assetID string, weight valueobjects.Percentage) error {
	// 기존 자산의 총 가중치 계산 (업데이트 대상 제외)
	var totalWeight float64
	found := false
	for _, asset := range p.Assets {
		if asset.AssetID != assetID {
			totalWeight += asset.Weight.Value
		} else {
			found = true
		}
	}

	if !found {
		return ErrAssetNotFound
	}

	// 새 가중치 적용 시 총 가중치가 100%를 초과하는지 확인
	if totalWeight+weight.Value > 100 {
		return ErrInvalidWeight
	}

	// 가중치 업데이트
	for i, asset := range p.Assets {
		if asset.AssetID == assetID {
			p.Assets[i].Weight = weight
			break
		}
	}

	p.UpdatedAt = time.Now()
	return nil
}

// RemoveAsset 포트폴리오에서 자산을 제거합니다.
func (p *Portfolio) RemoveAsset(assetID string) error {
	found := false
	newAssets := make([]PortfolioAsset, 0, len(p.Assets)-1)
	for _, asset := range p.Assets {
		if asset.AssetID != assetID {
			newAssets = append(newAssets, asset)
		} else {
			found = true
		}
	}

	if !found {
		return ErrAssetNotFound
	}

	p.Assets = newAssets
	p.UpdatedAt = time.Now()
	return nil
}

// Validate 포트폴리오의 유효성을 검증합니다.
func (p *Portfolio) Validate() error {
	var totalWeight float64
	for _, asset := range p.Assets {
		totalWeight += asset.Weight.Value
	}

	if totalWeight > 100 {
		return ErrInvalidWeight
	}

	return nil
}

// Error 정의
var (
	ErrInvalidWeight = NewDomainError("invalid_weight", "자산 가중치의 총합이 100%를 초과할 수 없습니다")
	ErrAssetNotFound = NewDomainError("asset_not_found", "자산을 찾을 수 없습니다")
)

// Error 도메인 에러를 나타냅니다.
type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// NewDomainError 새로운 도메인 에러를 생성합니다.
func NewDomainError(code string, message string) error {
	return Error{
		Code:    code,
		Message: message,
	}
}

// PortfolioRepository 포트폴리오 저장소 인터페이스입니다.
type PortfolioRepository interface {
	Save(ctx context.Context, portfolio *Portfolio) error
	FindByID(ctx context.Context, id string) (*Portfolio, error)
	Update(ctx context.Context, portfolio *Portfolio) error
	Delete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string) ([]*Portfolio, error)
}
