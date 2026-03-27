package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	supportauth "github.com/ljj/gugu-api/internal/support/auth"
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

func (i JWTTokenIssuer) IssueAccessToken(userID string, now time.Time) (supportauth.IssuedAccessToken, error) {
	accessExpiresAt := now.Add(i.accessTokenDuration)

	accessToken, err := i.sign(jwtClaims{
		Iss:  i.issuer,
		Sub:  userID,
		Type: accessTokenType,
		Iat:  now.Unix(),
		Exp:  accessExpiresAt.Unix(),
	})
	if err != nil {
		return supportauth.IssuedAccessToken{}, fmt.Errorf("sign access token: %w", err)
	}

	return supportauth.IssuedAccessToken{
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

func (i JWTTokenIssuer) VerifyAccessToken(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, i.secret)
	if _, err := mac.Write([]byte(signingInput)); err != nil {
		return "", fmt.Errorf("compute signature: %w", err)
	}
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return "", fmt.Errorf("invalid token signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode claims: %w", err)
	}

	var claims jwtClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return "", fmt.Errorf("unmarshal claims: %w", err)
	}

	if claims.Type != accessTokenType {
		return "", fmt.Errorf("invalid token type: %s", claims.Type)
	}

	if time.Now().Unix() > claims.Exp {
		return "", fmt.Errorf("token expired")
	}

	return claims.Sub, nil
}

func encodeJWTPart(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal jwt part: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(payload), nil
}
