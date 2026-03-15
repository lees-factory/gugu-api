# Auth Architecture

회원가입/로그인 영역은 다음 개념으로 분리했다.

- `email account`: 이메일과 패스워드를 가진 로컬 인증 계정
- `email verification`: 메일 인증 토큰 발급과 소모
- `oauth identity`: 외부 provider subject를 내부 유저에 연결하는 식별자
- `tokens`: 로그인 완료 후 발급되는 JWT access token / refresh token 묶음

레이어 흐름은 `core/api -> core/domain`까지만 직접 의존한다. `storage`, `clients`, `support`는 `cmd/api`의 wiring에서 주입된다.

## Core API

`internal/core/api`는 `chi` 라우터와 API 버전별 패키지 구성을 가진 runnable entry package다.

## Core Domain

`internal/core/domain/auth`의 `Service`는 아래 비즈니스 흐름을 가진다.

- 이메일 가입
- 이메일 로그인
- 이메일 인증
- OAuth 로그인

repository port는 개념 패키지별로 나뉜다.

- `internal/core/domain/user`
- `internal/core/domain/verification`
- `internal/core/domain/auth`

도구레이어는 별도 `implement` 패키지를 두지 않고 각 개념 패키지 내부에 둔다.

- `internal/core/domain/user/UserFinder.go`
- `internal/core/domain/user/UserWriter.go`
- `internal/core/domain/verification/EmailVerificationFinder.go`
- `internal/core/domain/verification/EmailVerificationWriter.go`
- `internal/core/domain/auth/OAuthIdentityFinder.go`
- `internal/core/domain/auth/OAuthIdentityWriter.go`
- `internal/core/domain/auth/SessionAppender.go`

의존 방향은 `Service -> 도구 -> Repository -> db-core 구현체`다.

OAuth는 provider 문자열 기반으로 저장하므로 Google 이후 Apple, Kakao 등을 추가해도 `oauth_identities` 테이블과 로그인 프로세서를 그대로 재사용할 수 있다.

## Storage

런타임은 지금 개념별 메모리 저장소로 부트스트랩했다.

- `internal/storage/memory/user`
- `internal/storage/memory/verification`
- `internal/storage/memory/auth`

실제 영속화는 `sqlc.yaml`, `internal/storage/dbcore/sqlc/schema/auth.sql`, `internal/storage/dbcore/sqlc/query/auth.sql`을 기준으로 `sqlc generate` 후 `internal/storage/dbcore`의 개념별 `*SqlcRepository`에 붙이면 된다.

## Clients

`internal/clients`는 외부 API와 크롤러 서버 연동을 위한 영역이다. 현재는 `aliexpress`, `crawler` client contract만 추가해 두었다.
