---
name: project_architecture_and_coding_guidelines
description: 프로젝트의 도메인 모델링, 레이어드 아키텍처 참조 규칙, 멀티 모듈 구성 및 데이터 액세스 제약 사항을 정의하는 핵심 가이드라인입니다.
---

# 프로젝트 아키텍처 및 코딩 가이드라인 (AI Agent Instructions)

본 문서는 이 프로젝트에서 코드를 생성, 수정, 리팩토링할 때 AI가 반드시 준수해야 하는 핵심 설계 원칙과 제약 사항입니다. 코드를 제안하기 전에 항상 이 규칙을 검증하십시오.

## 1. Domain Modeling: 개념(Concept)과 격벽(Wall)

### 1.1. 개념 (Concept) 정의 규칙
- **작은 개념 지향**: 하나의 개념(예: 대출)에 너무 많은 행위가 몰려 있다면, 반드시 개념을 분리하십시오. (예: 대출 -> 대출, 연체 등으로 분리)
- **상태의 개념화**: 특정 상태(예: 상환 실패)로 인해 새로운 비즈니스 흐름(재시도, 추가 이자)이 파생된다면, 이는 단순한 '상태'가 아니라 새로운 '개념'(예: 연체)으로 독립시켜 격벽을 세워야 합니다.
- **행위/상태 ≠ 개념**: 상태나 행위 자체를 개념으로 착각하지 마십시오. 개념은 상태를 표기하는 도메인 객체입니다.
- **외부 연동의 격리**: 외부 연동사나 외부 DB와 관련된 요소는 우리의 핵심 개념이 아닙니다. 별도의 시스템으로 격리하십시오.
    - 외부 api는 개념의 영역이 아니다.
    - 외부 api 와 연동되는 내부 db를 통해서 관리한다.
    - 외부 연동사(api) 는 참조 불가의 벽이 있다.

### 1.2. 격벽 (Wall) 규칙
- **책임 분리**: 한 개념의 변경이 다른 개념의 코드 수정을 유발해서는 안 됩니다.
- 코드를 작성할 때, 현재 수정하는 코드가 "어떤 개념의 격벽 안에 있는지" 명확히 인지하고, 선을 넘는 참조를 만들지 마십시오.

---

## 2. Layered Architecture: 레이어 통제 규칙

프로젝트는 다음 4개의 레이어로 구성되며, 아래의 참조 제약 규칙을 절대적으로(STRICTLY) 준수해야 합니다.

### 2.1. 레이어 정의
- **Presentation Layer**: 외부 요청/응답 처리, 외부 변화에 민감한 영역.
- **Business Layer**: 비즈니스 로직을 표현하고 흐름을 중계(Coordinator)하는 영역.
- **Implement Layer**: 상세 구현 로직을 가진 도구 클래스(Collaborators)들의 영역. (가장 재사용성이 높음)
    - **네이밍 컨벤션**: *Reader, *Finder, *Writer, *Appender, *Generator 등 책임을 명확히 할 것. 모호할 때만 *Process, *Manager 사용.
- **Data Access Layer**: 기술 의존성을 격리하고 순수 인터페이스를 제공하는 영역.

### 2.2. 🚫 레이어 참조 절대 규칙 (Crucial)
- **단방향 참조 (MUST)**: Presentation -> Business -> Implement -> Data Access 순서로만 참조해야 합니다.
- **역류 금지 (NEVER)**: 하위 레이어는 상위 레이어를 절대 참조할 수 없습니다. (예: Implement에서 Business 참조 불가)
- **건너뛰기 금지 (NEVER)**: 상위 레이어가 중간을 건너뛰고 하위 레이어를 직접 참조할 수 없습니다. (예: Business가 Data Access를 직접 호출 불가. 반드시 Implement를 거칠 것)
- **동일 레이어 참조 금지 (NEVER)**: 같은 레이어 클래스 간의 참조는 원칙적으로 금지됩니다.
    - **[예외 허용]**: 단, Implement Layer 내에서 도구 클래스 간의 조합/재사용을 위한 상호 참조는 허용합니다.

---

## 3. Business Logic Implementation: 중계자와 협력자

- **Service 클래스는 중계자(Coordinator)입니다.**
- Service (Business Layer) 내부에 직접 데이터 조회, 검증, 생성 등의 상세 로직을 구현하지 마십시오.
- 코드를 읽었을 때 상세 구현을 몰라도 "비즈니스의 흐름"이 한눈에 파악되도록 작성하십시오.
- **상세 로직은 협력 도구(Collaborator)에 위임하십시오.**
- 실제 작업은 단일 책임을 가진 Implement Layer의 도구 클래스들에게 위임하십시오.
- **Business Layer 오염 금지**: Business Layer는 Data Access의 구현 기술(JPA 등)을 전혀 몰라야 합니다.

---

## 4. Multi-Module Configuration: 모듈화 규칙

Gradle 멀티 모듈 환경에서 결합도를 낮추기 위한 규칙입니다.

### 4.1. 의존성 선언 규칙 (Gradle)
- 상위 모듈로의 의존성 전파를 막기 위해 `api` 키워드 사용을 지양하고 반드시 `implementation` 키워드를 사용하십시오.
- 예: `db-core` 모듈이 JPA를 사용하더라도, 이를 참조하는 `core-api` 모듈은 JPA의 존재(어노테이션, 엔티티 매니저 등)를 알 수 없어야 합니다.

### 4.2. 모듈 vs 레이어 vs 개념의 독립성 (WARNING)
- **모듈 ≠ 레이어**: 대칭되지 않습니다.
- **모듈/레이어 ≠ 개념/격벽**: 물리적 분리(모듈/레이어)와 논리적 분리(개념/격벽)를 혼동하지 마십시오.

---

## 5. Directory Structure Mapping

새로운 클래스나 패키지를 생성할 때 다음 모듈 구조의 역할을 엄격히 따르십시오.

- **clients/* (예: exchange-rate)**: 외부 API 연동 레이어 바깥 영역.
- **core/core-api**: Presentation Layer (Controller, Request/Response DTO).
- **core/core-domain**: Business Layer (Service) 및 Implement Layer (Reader, Writer 등 도구). 상위 레이어가 필요하면 여기에 추가.
- **core/core-enum**: 공통 열거형 (Enums).
- **storage/db-core**: Data Access Layer (JPA Entity, Repository 구현체 등). 데이터베이스 접근 기술 격리.
- **support/* (logging, monitoring)**: 인프라/지원 도구.
- **web/security**: 보안 관련 설정(spring-security 를 사용할경우). spring-security는 횡단관심사 영역으로 레이어 개념에 귀속되지 않음.

---

## 6. Data Access Layer
- JPA 사용 시 연관관계 사용을 하지 않는다.
- Join을 쓰게 된다면 Join 자체가 격벽을 침범한다고 가정하고 기준을 적절히 유지하기 위해 Join table도 응집 기준으로 쓴다.
- Join으로 인해 God Class가 탄생하지 않아야 한다.

---

**[AI 최종 확인 사항]**
코드를 제안하기 전에 1) 비즈니스 흐름이 쉽게 읽히는지, 2) 레이어 참조 규칙(건너뛰기/역류)을 위반하지 않았는지, 3) 도메인 개념이 비대해지지 않았는지 스스로 검증하십시오.