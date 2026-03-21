# Front Local Run Guide

이 저장소에는 프론트엔드 앱 코드가 포함되어 있지 않다.
대신 프론트가 로컬에서 이 API를 붙여서 개발할 수 있도록 백엔드 실행 방법과 프록시 기준을 정리한다.

## 1. 백엔드 실행

### 환경 변수 준비

```bash
cp .env.example .env
# .env 파일을 열어서 값 채우기
```

`.env.example`을 복사해서 값을 채운다. 주요 환경변수:

| 변수 | 필수 | 기본값 | 설명 |
|------|------|--------|------|
| `DATABASE_URL` | 선택 | (빈값=in-memory) | Supabase PostgreSQL 접속 URL |
| `ALIEXPRESS_APP_KEY` | 필수 | | AliExpress Affiliate API 키 |
| `ALIEXPRESS_APP_SECRET` | 필수 | | AliExpress Affiliate API 시크릿 |
| `CRAWLER_BASE_URL` | 선택 | `http://localhost:8000` | 크롤러 서버 주소 |
| `HTTP_ADDRESS` | 선택 | `:8080` | API 서버 포트 |
| `CORS_ALLOWED_ORIGINS` | 선택 | `http://localhost:3000,http://localhost:5173` | 프론트 허용 origin |
| `JWT_SECRET` | 선택 | `change-me` | JWT 서명 키 |

### Supabase 초기 스키마 적용

Supabase SQL Editor 에서 [init.sql](/Users/LJJ/Desktop/project/go/gugu/gugu-api/docs/init.sql) 을 실행한다.

### 로컬 실행

이 프로젝트는 시작 시 루트의 `.env` 를 자동으로 읽는다.
`.env` 파일이 없으면 조용히 넘어가고, 파일이 있는데 파싱 오류가 있을 때만 경고 로그를 남긴다.

```bash
go run ./cmd/api
```

수동으로 shell 환경 변수에 올리는 방식도 그대로 사용할 수 있다.

```bash
set -a
source .env
set +a

go run ./cmd/api
```

정상 실행되면 API 는 `http://localhost:8080` 에서 뜬다.
별도 빌드 과정 없이 `go run`이 자동으로 컴파일 후 실행한다.

바이너리를 만들고 싶은 경우:

```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

### 크롤러 서버 연동

크롤러 서버(gugu-crawler)는 AliExpress Affiliate API가 실패할 때 fallback으로 사용된다.

```bash
# .env에 크롤러 서버 주소 설정 (기본값: http://localhost:8000)
CRAWLER_BASE_URL=http://localhost:8000
```

- 크롤러 서버가 꺼져 있어도 API 서버는 정상 기동된다.
- Affiliate API가 성공하면 크롤러를 호출하지 않는다.
- 크롤러 서버를 로컬에서 띄우려면 gugu-crawler 저장소의 안내를 따른다.

확인:

```bash
curl http://localhost:8080/health
```

예상 응답:

```json
{"result":"SUCCESS","data":{"status":"ok"}}
```

## 2. 프론트에서 붙이는 방법

이 백엔드는 기본적으로 `http://localhost:3000`, `http://localhost:5173` 에 대한 CORS 를 허용한다.
필요하면 `CORS_ALLOWED_ORIGINS` 환경 변수에 쉼표로 구분한 origin 목록을 넣어서 변경할 수 있다.

로컬 개발에서는 CORS 를 직접 써도 되지만, 프론트 dev server 의 proxy 기능으로 `/api` 요청을 `http://localhost:8080` 으로 넘기면 배포 환경과 더 비슷하게 맞출 수 있다.

권장 기준:

- 프론트 주소: `http://localhost:3000` 또는 `http://localhost:5173`
- 백엔드 주소: `http://localhost:8080`
- 프론트 호출 경로: `/api/...`
- 프록시 대상: `http://localhost:8080/...`

## 3. Vite 예시

`vite.config.ts`

```ts
import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
})
```

프론트 코드 예시:

```ts
const response = await fetch('/api/health')
const data = await response.json()
```

## 4. Next.js 예시

`next.config.ts`

```ts
import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/:path*',
      },
    ]
  },
}

export default nextConfig
```

프론트 코드 예시:

```ts
const response = await fetch('/api/health')
const data = await response.json()
```

## 5. 주요 API

명세는 [openapi.yml](/Users/LJJ/Desktop/project/go/gugu/gugu-api/openapi.yml) 에 있다.

주요 엔드포인트:

| Method | Path | 설명 |
|--------|------|------|
| `GET` | `/health` | 헬스 체크 |
| `POST` | `/v1/auth/register/email` | 이메일 회원가입 |
| `POST` | `/v1/auth/verify-email` | 이메일 인증 |
| `POST` | `/v1/auth/login/email` | 이메일 로그인 |
| `POST` | `/v1/auth/oauth/login` | OAuth 로그인 |
| `POST` | `/v1/auth/refresh` | 토큰 갱신 |
| `POST` | `/v1/auth/logout` | 로그아웃 |
| `GET` | `/v1/tracked-items?user_id=` | 추적 상품 목록 |
| `POST` | `/v1/tracked-items` | 추적 상품 추가 |
| `GET` | `/v1/tracked-items/{id}?user_id=` | 추적 상품 상세 |
| `DELETE` | `/v1/tracked-items/{id}?user_id=` | 추적 상품 삭제 |
| `PATCH` | `/v1/tracked-items/{id}/sku` | SKU 선택 |
| `GET` | `/v1/products/{id}?user_id=` | 상품 상세 (SKU 포함) |
| `GET` | `/v1/products/{id}/skus` | 상품 SKU 목록 |
| `POST` | `/v1/integrations/aliexpress/authorize-url` | AliExpress 인가 URL 생성 |
| `POST` | `/v1/integrations/aliexpress/exchange-code` | AliExpress 토큰 교환 |
| `POST` | `/v1/integrations/aliexpress/refresh-token` | AliExpress 토큰 갱신 |
| `GET` | `/v1/integrations/aliexpress/connection-status` | AliExpress 연결 상태 |
| `GET` | `/v1/integrations/aliexpress/product-detail?product_id=` | AliExpress 상품 상세 (개발자용) |
| `GET` | `/v1/integrations/aliexpress/product-sku-detail?product_id=` | AliExpress SKU 상세 (개발자용) |
| `POST` | `/v1/batch/update-prices` | 전체 상품 가격 일괄 업데이트 |

## 6. 요청 예시

회원가입:

```bash
curl -X POST http://localhost:8080/v1/auth/register/email \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "user@example.com",
    "password": "secret123!",
    "display_name": "Gugu User"
  }'
```

이메일 로그인:

```bash
curl -X POST http://localhost:8080/v1/auth/login/email \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "user@example.com",
    "password": "secret123!"
  }'
```

로그인 응답은 서버 세션 대신 JWT 토큰 묶음을 반환한다.

이메일 인증 코드는 Gmail SMTP로 발송한다.
개발 중 발송 없이 로그만 보려면 `MAIL_PROVIDER='log'` 로 바꾸면 된다.

## 7. AliExpress 토큰 설정

AliExpress Affiliate API는 `access_token`이 필수다. 최초 1회 OAuth 인증이 필요하다.

```bash
# 1. 인가 URL 생성
curl -X POST http://localhost:8080/v1/integrations/aliexpress/authorize-url

# 2. 응답의 authorization_url을 브라우저에서 열고 로그인 → callback URL에서 code 확인

# 3. code를 토큰으로 교환
curl -X POST http://localhost:8080/v1/integrations/aliexpress/exchange-code \
  -H 'Content-Type: application/json' \
  -d '{"code":"받은_code"}'
```

access_token은 24시간 만료. refresh_token(6개월)으로 갱신 가능:

```bash
curl -X POST http://localhost:8080/v1/integrations/aliexpress/refresh-token
```

토큰 상태 확인:

```bash
curl http://localhost:8080/v1/integrations/aliexpress/connection-status
```

## 8. 가격 배치 업데이트

등록된 모든 상품의 가격을 외부 API에서 일괄 조회하여 업데이트한다.

수동 실행:

```bash
curl -X POST http://localhost:8080/v1/batch/update-prices
```

AWS Lightsail 배포 시 crontab으로 하루 1회 자동 실행:

```bash
# crontab -e
0 3 * * * curl -sS -X POST http://localhost:8080/v1/batch/update-prices
```

## 9. 현재 제한 사항

- 프론트 앱 자체는 이 저장소에 없다.
- `DATABASE_URL` 이 없으면 저장소는 메모리 기반으로 동작한다.
- `DATABASE_URL` 이 있으면 Supabase(Postgres)를 실제 사용자/이메일 인증/OAuth 식별자 저장소로 사용한다.
