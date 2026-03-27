package product

import (
	"context"
	"log"
	"sync"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

// CompositeProvider는 상품 기본 정보와 SKU를 각각 다른 소스에서 병렬 수집하여 조합한다.
//   - metadataProvider: 상품 기본 정보 (title, price, image, product_url) + 부분 SKU
//   - skuProvider: 전체 SKU 수집 (전량 보장)
//
// 조합 규칙:
//   - 양쪽 성공 → 상품 정보는 metadataProvider, SKU는 skuProvider 사용
//   - metadata 성공 + sku 실패 → metadata의 부분 SKU로 저장 (부분 수집)
//   - metadata 실패 + sku 성공 → skuProvider 데이터로 상품 + SKU 저장
//   - 양쪽 실패 → 에러
type CompositeProvider struct {
	metadataProvider domainproduct.ProductProvider
	skuProvider      domainproduct.ProductProvider
}

func NewCompositeProvider(metadataProvider domainproduct.ProductProvider, skuProvider domainproduct.ProductProvider) *CompositeProvider {
	return &CompositeProvider{
		metadataProvider: metadataProvider,
		skuProvider:      skuProvider,
	}
}

func (p *CompositeProvider) Provide(ctx context.Context, market enum.Market, externalProductID string, originalURL string) (*domainproduct.NewProduct, error) {
	var (
		metadataResult *domainproduct.NewProduct
		metadataErr    error
		skuResult      *domainproduct.NewProduct
		skuErr         error
		wg             sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		metadataResult, metadataErr = p.metadataProvider.Provide(ctx, market, externalProductID, originalURL)
	}()

	go func() {
		defer wg.Done()
		skuResult, skuErr = p.skuProvider.Provide(ctx, market, externalProductID, originalURL)
	}()

	wg.Wait()

	// 양쪽 성공: 상품 정보는 metadata, SKU는 skuProvider
	if metadataErr == nil && metadataResult != nil && skuErr == nil && skuResult != nil {
		metadataResult.SKUs = skuResult.SKUs
		return metadataResult, nil
	}

	// metadata 성공 + sku 실패: metadata의 부분 SKU 사용
	if metadataErr == nil && metadataResult != nil {
		if skuErr != nil {
			log.Printf("sku provider failed, using metadata skus: %v", skuErr)
		}
		return metadataResult, nil
	}

	// metadata 실패 + sku 성공: skuProvider 데이터 사용
	if skuErr == nil && skuResult != nil {
		if metadataErr != nil {
			log.Printf("metadata provider failed, using sku provider data: %v", metadataErr)
		}
		return skuResult, nil
	}

	// 양쪽 실패
	log.Printf("both providers failed: metadata=%v, sku=%v", metadataErr, skuErr)
	return nil, coreerror.New(coreerror.ProductNotFound)
}
