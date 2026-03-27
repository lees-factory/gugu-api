package aliexpress

import (
	"context"

	alidto "github.com/ljj/gugu-api/internal/clients/aliexpress/dto"
)

type TokenExchangeInput = alidto.TokenExchangeInput
type RefreshTokenInput = alidto.RefreshTokenInput
type TokenSet = alidto.TokenSet

type ProductLookupInput = alidto.ProductLookupInput
type CategoryResult = alidto.CategoryResult
type AffiliateCategory = alidto.AffiliateCategory
type ProductQueryInput = alidto.ProductQueryInput
type ProductQueryResult = alidto.ProductQueryResult
type ProductDetailInput = alidto.ProductDetailInput
type ProductDetailResult = alidto.ProductDetailResult
type ProductSKUDetailInput = alidto.ProductSKUDetailInput
type ProductSKUDetailResult = alidto.ProductSKUDetailResult
type ProductSnapshot = alidto.ProductSnapshot
type AffiliateProduct = alidto.AffiliateProduct
type PromoCodeInfo = alidto.PromoCodeInfo
type AffiliateSKUItemInfo = alidto.AffiliateSKUItemInfo
type AffiliateSKUInfo = alidto.AffiliateSKUInfo

type AuthClient interface {
	BuildAuthorizationURL() (string, error)
	ExchangeCode(ctx context.Context, input TokenExchangeInput) (*TokenSet, error)
	RefreshAccessToken(ctx context.Context, input RefreshTokenInput) (*TokenSet, error)
}

type ProductClient interface {
	GetAffiliateCategories(ctx context.Context) (*CategoryResult, error)
	GetAffiliateProducts(ctx context.Context, input ProductQueryInput) (*ProductQueryResult, error)
	GetAffiliateHotProducts(ctx context.Context, input ProductQueryInput) (*ProductQueryResult, error)
	GetAffiliateProductDetail(ctx context.Context, input ProductDetailInput) (*ProductDetailResult, error)
	GetAffiliateProductSKUDetail(ctx context.Context, input ProductSKUDetailInput) (*ProductSKUDetailResult, error)
	GetProductSnapshot(ctx context.Context, input ProductLookupInput) (*ProductSnapshot, error)
}

type Client interface {
	AuthClient
	ProductClient
}
