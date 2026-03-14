package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type JWTTokenIssuer struct {
	secret               []byte
	issuer               string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtClaims struct {
	Iss  string `json:"iss"`
	Sub  string `json:"sub"`
	Type string `json:"type"`
	Iat  int64  `json:"iat"`
	Exp  int64  `json:"exp"`
}

func NewJWTTokenIssuer(secret string, issuer string) JWTTokenIssuer {
	return JWTTokenIssuer{
		secret:               []byte(secret),
		issuer:               issuer,
		accessTokenDuration:  15 * time.Minute,
		refreshTokenDuration: 14 * 24 * time.Hour,
	}
}

func (i JWTTokenIssuer) IssueAccessToken(userID string, now time.Time) (domainauth.IssuedAccessToken, error) {
	accessExpiresAt := now.Add(i.accessTokenDuration)

	accessToken, err := i.sign(jwtClaims{
		Iss:  i.issuer,
		Sub:  userID,
		Type: accessTokenType,
		Iat:  now.Unix(),
		Exp:  accessExpiresAt.Unix(),
	})
	if err != nil {
		return domainauth.IssuedAccessToken{}, fmt.Errorf("sign access token: %w", err)
	}

	return domainauth.IssuedAccessToken{
		Token:     accessToken,
		ExpiresAt: accessExpiresAt,
	}, nil
}

func (i JWTTokenIssuer) sign(claims jwtClaims) (string, error) {
	if len(i.secret) == 0 {
		return "", fmt.Errorf("jwt secret is empty")
	}

	headerPart, err := encodeJWTPart(jwtHeader{Alg: "HS256", Typ: "JWT"})
	if err != nil {
		return "", err
	}
	claimsPart, err := encodeJWTPart(claims)
	if err != nil {
		return "", err
	}

	signingInput := strings.Join([]string{headerPart, claimsPart}, ".")
	mac := hmac.New(sha256.New, i.secret)
	if _, err := mac.Write([]byte(signingInput)); err != nil {
		return "", fmt.Errorf("write signing input: %w", err)
	}

	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signingInput + "." + signature, nil
}

func encodeJWTPart(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal jwt part: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(payload), nil
}
