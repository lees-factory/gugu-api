# Auth API

초기 구현 범위는 이메일 가입/로그인/인증과 확장 가능한 OAuth 로그인입니다.

## Endpoints

### POST `/v1/auth/register/email`

```json
{
  "email": "user@example.com",
  "password": "strong-password",
  "display_name": "LJJ"
}
```

응답:

```json
{
  "user": {
    "id": "generated-user-id",
    "email": "user@example.com",
    "display_name": "LJJ",
    "email_verified": false,
    "created_at": "2026-03-09T00:00:00Z"
  },
  "verification_token": "token",
  "verification_dispatched": true
}
```

### POST `/v1/auth/verify-email`

```json
{
  "token": "verification-token"
}
```

### POST `/v1/auth/login/email`

```json
{
  "email": "user@example.com",
  "password": "strong-password"
}
```

### POST `/v1/auth/oauth/login`

현재는 provider-neutral contract만 제공합니다. 실제 Google OAuth2 callback 핸들링은 다음 단계에서 이 endpoint 앞단에 붙이면 됩니다.

```json
{
  "provider": "google",
  "subject": "google-user-sub",
  "email": "user@example.com",
  "display_name": "Google User"
}
```
