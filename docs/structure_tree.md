```
├── Dockerfile
├── build.gradle
├── clients
│   └── client-aligo
│       ├── build.gradle
│       └── src
│           └── main
│               └── java
│                   └── com
│                       └── nexsol
│                           └── tpa
│                               └── client
│                                   └── aligo
│                                       ├── AligoClient.java
│                                       ├── AligoSmsSender.java
│                                       └── config
│                                           ├── AligoClientAutoConfiguration.java
│                                           └── AligoFeignConfig.java
├── core
│   ├── core-api
│   │   ├── build.gradle
│   │   └── src
│   │       ├── main
│   │       │   ├── java
│   │       │   │   └── com
│   │       │   │       └── nexsol
│   │       │   │           └── tpa
│   │       │   │               ├── CoreApiApplication.java
│   │       │   │               └── core
│   │       │   │                   ├── DocsController.java
│   │       │   │                   └── api
│   │       │   │                       ├── config
│   │       │   │                       │   ├── AsyncConfig.java
│   │       │   │                       │   ├── AsyncExceptionHandler.java
│   │       │   │                       │   └── WebMvcConfig.java
│   │       │   │                       ├── controller
│   │       │   │                       │   ├── ApiControllerAdvice.java
│   │       │   │                       │   ├── HealthController.java
│   │       │   │                       │   └── v1
│   │       │   │                       │       ├── AccidentController.java
│   │       │   │                       │       ├── AuthController.java
│   │       │   │                       │       ├── FileController.java
│   │       │   │                       │       ├── InsuranceController.java
│   │       │   │                       │       ├── PaymentController.java
│   │       │   │                       │       ├── UserController.java
│   │       │   │                       │       ├── request
│   │       │   │                       │       │   ├── AccidentAttachmentRequest.java
│   │       │   │                       │       │   ├── AccidentInfoRequest.java
│   │       │   │                       │       │   ├── AccidentReportRequest.java
│   │       │   │                       │       │   ├── AccidentRequest.java
│   │       │   │                       │       │   ├── CheckPlantNameRequest.java
│   │       │   │                       │       │   ├── DepositRequest.java
│   │       │   │                       │       │   ├── DocumentRequest.java
│   │       │   │                       │       │   ├── DocumentSetRequest.java
│   │       │   │                       │       │   ├── EmailSendRequest.java
│   │       │   │                       │       │   ├── EmailVerifyRequest.java
│   │       │   │                       │       │   ├── InsuranceCompleteRequest.java
│   │       │   │                       │       │   ├── InsuranceConditionRequest.java
│   │       │   │                       │       │   ├── InsurancePlantRequest.java
│   │       │   │                       │       │   ├── InsuranceSearchRequest.java
│   │       │   │                       │       │   ├── InsuranceStartRequest.java
│   │       │   │                       │       │   ├── ModifyUserRequest.java
│   │       │   │                       │       │   ├── PaymentCancelRequest.java
│   │       │   │                       │       │   ├── PledgeRequest.java
│   │       │   │                       │       │   ├── SignInRequest.java
│   │       │   │                       │       │   └── SignUpRequest.java
│   │       │   │                       │       └── response
│   │       │   │                       │           ├── AccidentReportDetailResponse.java
│   │       │   │                       │           ├── AccidentReportListResponse.java
│   │       │   │                       │           ├── AccidentReportResponse.java
│   │       │   │                       │           ├── AuthResponse.java
│   │       │   │                       │           ├── FileResponse.java
│   │       │   │                       │           ├── InsuranceListResponse.java
│   │       │   │                       │           ├── InsuranceResponse.java
│   │       │   │                       │           ├── InsuranceStatusResponse.java
│   │       │   │                       │           ├── InsuranceSummaryResponse.java
│   │       │   │                       │           ├── PaymentResponse.java
│   │       │   │                       │           ├── PlantNameCheckResponse.java
│   │       │   │                       │           ├── SignUpResponse.java
│   │       │   │                       │           └── UserResponse.java
│   │       │   │                       ├── scheduler
│   │       │   │                       │   └── DraftCleanupScheduler.java
│   │       │   │                       └── support
│   │       │   │                           ├── error
│   │       │   │                           │   ├── CoreApiErrorCode.java
│   │       │   │                           │   ├── CoreApiErrorMessage.java
│   │       │   │                           │   ├── CoreApiErrorType.java
│   │       │   │                           │   └── CoreApiException.java
│   │       │   │                           └── response
│   │       │   │                               ├── ApiResponse.java
│   │       │   │                               ├── InsuranceDivisionCalculator.java
│   │       │   │                               ├── PageResponse.java
│   │       │   │                               └── ResultType.java
│   │       │   └── resources
│   │       │       └── application.yml
│   │       └── test
│   │           ├── java
│   │           │   └── com
│   │           │       └── nexsol
│   │           │           └── tpa
│   │           │               └── core
│   │           │                   ├── DocsControllerTest.java
│   │           │                   └── api
│   │           │                       └── controller
│   │           │                           ├── HealthControllerTest.java
│   │           │                           └── v1
│   │           │                               ├── AccidentControllerTest.java
│   │           │                               ├── AuthControllerTest.java
│   │           │                               ├── FileControllerTest.java
│   │           │                               ├── InsuranceControllerTest.java
│   │           │                               ├── PaymentControllerTest.java
│   │           │                               └── UserControllerTest.java
│   │           └── resources
│   ├── core-domain
│   │   ├── build.gradle
│   │   └── src
│   │       ├── main
│   │       │   └── java
│   │       │       └── com
│   │       │           └── nexsol
│   │       │               └── tpa
│   │       │                   └── core
│   │       │                       ├── domain
│   │       │                       │   ├── accident
│   │       │                       │   │   ├── Accident.java
│   │       │                       │   │   ├── AccidentAttachment.java
│   │       │                       │   │   ├── AccidentContractValidator.java
│   │       │                       │   │   ├── AccidentInfo.java
│   │       │                       │   │   ├── AccidentInsured.java
│   │       │                       │   │   ├── AccidentNumberGenerator.java
│   │       │                       │   │   ├── AccidentPlant.java
│   │       │                       │   │   ├── AccidentReport.java
│   │       │                       │   │   ├── AccidentReportAppender.java
│   │       │                       │   │   ├── AccidentReportDetail.java
│   │       │                       │   │   ├── AccidentReportReader.java
│   │       │                       │   │   ├── AccidentReportRepository.java
│   │       │                       │   │   ├── AccidentReportService.java
│   │       │                       │   │   └── NewAccidentReport.java
│   │       │                       │   ├── auth
│   │       │                       │   │   ├── AuthService.java
│   │       │                       │   │   ├── AuthToken.java
│   │       │                       │   │   ├── JwtPayload.java
│   │       │                       │   │   ├── RefreshToken.java
│   │       │                       │   │   ├── RefreshTokenRepository.java
│   │       │                       │   │   ├── TokenAppender.java
│   │       │                       │   │   ├── TokenIssuer.java
│   │       │                       │   │   ├── TokenReader.java
│   │       │                       │   │   └── TokenRemover.java
│   │       │                       │   ├── channel
│   │       │                       │   │   ├── Channel.java
│   │       │                       │   │   ├── ChannelReader.java
│   │       │                       │   │   └── ChannelRepository.java
│   │       │                       │   ├── file
│   │       │                       │   │   ├── DocumentFile.java
│   │       │                       │   │   ├── DocumentInfo.java
│   │       │                       │   │   ├── FileInfo.java
│   │       │                       │   │   ├── FileService.java
│   │       │                       │   │   └── FileStorageClient.java
│   │       │                       │   ├── insurance
│   │       │                       │   │   ├── Agreement.java
│   │       │                       │   │   ├── Applicant.java
│   │       │                       │   │   ├── InsuranceApplication.java
│   │       │                       │   │   ├── InsuranceApplicationReader.java
│   │       │                       │   │   ├── InsuranceApplicationRepository.java
│   │       │                       │   │   ├── InsuranceApplicationResult.java
│   │       │                       │   │   ├── InsuranceApplicationValidator.java
│   │       │                       │   │   ├── InsuranceApplicationWriter.java
│   │       │                       │   │   ├── InsuranceAttachment.java
│   │       │                       │   │   ├── InsuranceCommandService.java
│   │       │                       │   │   ├── InsuranceCondition.java
│   │       │                       │   │   ├── InsuranceConditionPolicy.java
│   │       │                       │   │   ├── InsuranceDocument.java
│   │       │                       │   │   ├── InsuranceInspector.java
│   │       │                       │   │   ├── InsuranceListResult.java
│   │       │                       │   │   ├── InsurancePlant.java
│   │       │                       │   │   ├── InsuranceQueryService.java
│   │       │                       │   │   ├── InsuranceStatusResult.java
│   │       │                       │   │   ├── PlantNameCheckResult.java
│   │       │                       │   │   ├── PlantNameOwnerInfo.java
│   │       │                       │   │   └── Pledge.java
│   │       │                       │   ├── insurer
│   │       │                       │   │   ├── Insurer.java
│   │       │                       │   │   ├── InsurerQueryService.java
│   │       │                       │   │   ├── InsurerReader.java
│   │       │                       │   │   └── InsurerRepository.java
│   │       │                       │   ├── notification
│   │       │                       │   │   ├── EmailSender.java
│   │       │                       │   │   └── SmsSender.java
│   │       │                       │   ├── partner
│   │       │                       │   │   ├── Partner.java
│   │       │                       │   │   ├── PartnerReader.java
│   │       │                       │   │   └── PartnerRepository.java
│   │       │                       │   ├── payment
│   │       │                       │   │   ├── Payment.java
│   │       │                       │   │   ├── PaymentAppender.java
│   │       │                       │   │   ├── PaymentCancel.java
│   │       │                       │   │   ├── PaymentCancelAppender.java
│   │       │                       │   │   ├── PaymentCancelReader.java
│   │       │                       │   │   ├── PaymentCancelRepository.java
│   │       │                       │   │   ├── PaymentCancelService.java
│   │       │                       │   │   ├── PaymentCommandService.java
│   │       │                       │   │   ├── PaymentConfirmProcessor.java
│   │       │                       │   │   ├── PaymentCreator.java
│   │       │                       │   │   ├── PaymentDepositProcessor.java
│   │       │                       │   │   ├── PaymentQueryService.java
│   │       │                       │   │   ├── PaymentReader.java
│   │       │                       │   │   ├── PaymentRepository.java
│   │       │                       │   │   ├── PaymentValidator.java
│   │       │                       │   │   ├── PaymentWriter.java
│   │       │                       │   │   └── RefundCalculator.java
│   │       │                       │   ├── policy
│   │       │                       │   │   ├── InsurancePolicy.java
│   │       │                       │   │   ├── InsurancePolicyReader.java
│   │       │                       │   │   └── InsurancePolicyRepository.java
│   │       │                       │   ├── premium
│   │       │                       │   │   ├── AreaKeyResolver.java
│   │       │                       │   │   ├── InsurancePremiumCalculator.java
│   │       │                       │   │   ├── InsuranceRatePolicy.java
│   │       │                       │   │   ├── InsuranceRateReader.java
│   │       │                       │   │   ├── InsuranceRateRepository.java
│   │       │                       │   │   ├── PremiumQuote.java
│   │       │                       │   │   └── RateTable.java
│   │       │                       │   ├── staff
│   │       │                       │   │   ├── AccidentStaff.java
│   │       │                       │   │   ├── AccidentStaffFinder.java
│   │       │                       │   │   ├── AccidentStaffRepository.java
│   │       │                       │   │   └── AccidentStaffService.java
│   │       │                       │   ├── user
│   │       │                       │   │   ├── ModifyUser.java
│   │       │                       │   │   ├── NewUser.java
│   │       │                       │   │   ├── User.java
│   │       │                       │   │   ├── UserAppender.java
│   │       │                       │   │   ├── UserFinder.java
│   │       │                       │   │   ├── UserReader.java
│   │       │                       │   │   ├── UserRepository.java
│   │       │                       │   │   └── UserService.java
│   │       │                       │   └── verification
│   │       │                       │       ├── EmailGenerateCode.java
│   │       │                       │       ├── EmailSendValidator.java
│   │       │                       │       ├── EmailUpdateValidator.java
│   │       │                       │       ├── EmailVerification.java
│   │       │                       │       ├── EmailVerificationAppender.java
│   │       │                       │       ├── EmailVerificationFinder.java
│   │       │                       │       ├── EmailVerificationReader.java
│   │       │                       │       ├── EmailVerificationRepository.java
│   │       │                       │       └── EmailVerifiedService.java
│   │       │                       ├── error
│   │       │                       │   ├── CoreErrorCode.java
│   │       │                       │   ├── CoreErrorKind.java
│   │       │                       │   ├── CoreErrorLevel.java
│   │       │                       │   ├── CoreErrorType.java
│   │       │                       │   └── CoreException.java
│   │       │                       └── support
│   │       │                           ├── PageResult.java
│   │       │                           ├── SliceResult.java
│   │       │                           └── SortPage.java
│   │       ├── test
│   │       │   └── java
│   │       │       └── com
│   │       │           └── nexsol
│   │       │               └── tpa
│   │       │                   └── core
│   │       │                       └── domain
│   │       │                           ├── AuthServiceTest.java
│   │       │                           ├── EmailVerifiedServiceTest.java
│   │       │                           ├── InsurancePremiumCalculatorTest.java
│   │       │                           ├── NewUserTest.java
│   │       │                           └── UserServiceTest.java
│   │       └── testFixtures
│   │           └── java
│   │               └── com
│   │                   └── nexsol
│   │                       └── tpa
│   │                           └── core
│   │                               └── domain
│   │                                   ├── EmailVerificationFixture.java
│   │                                   ├── NewUserBuilder.java
│   │                                   └── UserFixture.java
│   └── core-enum
│       ├── build.gradle
│       └── src
│           └── main
│               └── java
│                   └── com
│                       └── nexsol
│                           └── tpa
│                               └── core
│                                   └── enums
│                                       ├── AccidentStatus.java
│                                       ├── BondSendStatus.java
│                                       ├── EmailVerifiedType.java
│                                       ├── EntityStatus.java
│                                       ├── InsuranceDocumentType.java
│                                       ├── InsuranceStatus.java
│                                       ├── PaymentStatus.java
│                                       └── RateType.java
├── docker-compose.app.yml
├── docker-compose.yml
├── gradle
│   └── wrapper
│       ├── gradle-wrapper.jar
│       └── gradle-wrapper.properties
├── gradle.properties
├── gradlew
├── gradlew.bat
├── lint.gradle
├── scripts
│   └── deploy.sh
├── settings.gradle
├── storage
│   ├── db-core
│   │   ├── build.gradle
│   │   └── src
│   │       ├── main
│   │       │   ├── java
│   │       │   │   └── com
│   │       │   │       └── nexsol
│   │       │   │           └── tpa
│   │       │   │               └── storage
│   │       │   │                   └── db
│   │       │   │                       └── core
│   │       │   │                           ├── config
│   │       │   │                           │   ├── BondSendStatusConverter.java
│   │       │   │                           │   ├── CoreDataSourceConfig.java
│   │       │   │                           │   └── CoreJpaConfig.java
│   │       │   │                           ├── entity
│   │       │   │                           │   ├── AccidentAttachmentEntity.java
│   │       │   │                           │   ├── AccidentHistoryEntity.java
│   │       │   │                           │   ├── AccidentInfoEmbeddable.java
│   │       │   │                           │   ├── AccidentReportEntity.java
│   │       │   │                           │   ├── AccidentStaffEntity.java
│   │       │   │                           │   ├── AgreementInfoEmbeddable.java
│   │       │   │                           │   ├── ApplicantInfoEmbeddable.java
│   │       │   │                           │   ├── BaseEntity.java
│   │       │   │                           │   ├── ChannelEntity.java
│   │       │   │                           │   ├── EmailVerificationEntity.java
│   │       │   │                           │   ├── InsuranceApplicationEntity.java
│   │       │   │                           │   ├── InsuranceAttachmentEntity.java
│   │       │   │                           │   ├── InsuranceConditionEntity.java
│   │       │   │                           │   ├── InsurancePlantEmbeddable.java
│   │       │   │                           │   ├── InsurancePolicyEntity.java
│   │       │   │                           │   ├── InsuranceRateEntity.java
│   │       │   │                           │   ├── InsurerEntity.java
│   │       │   │                           │   ├── PartnerEntity.java
│   │       │   │                           │   ├── PaymentCancelEntity.java
│   │       │   │                           │   ├── PaymentEntity.java
│   │       │   │                           │   ├── PledgeEmbeddable.java
│   │       │   │                           │   ├── QuoteEmbeddable.java
│   │       │   │                           │   ├── RefreshTokenEntity.java
│   │       │   │                           │   └── UserEntity.java
│   │       │   │                           └── repository
│   │       │   │                               ├── AccidentAttachmentJpaRepository.java
│   │       │   │                               ├── AccidentHistoryJpaRepository.java
│   │       │   │                               ├── AccidentReportJpaRepository.java
│   │       │   │                               ├── AccidentReportRepositoryImpl.java
│   │       │   │                               ├── AccidentStaffJpaRepository.java
│   │       │   │                               ├── AccidentStaffRepositoryImpl.java
│   │       │   │                               ├── ChannelJpaRepository.java
│   │       │   │                               ├── ChannelRepositoryImpl.java
│   │       │   │                               ├── EmailVerificationJpaRepository.java
│   │       │   │                               ├── EmailVerificationRepositoryImpl.java
│   │       │   │                               ├── InsuranceApplicationJpaRepository.java
│   │       │   │                               ├── InsuranceApplicationRepositoryImpl.java
│   │       │   │                               ├── InsuranceAttachmentJpaRepository.java
│   │       │   │                               ├── InsuranceConditionJpaRepository.java
│   │       │   │                               ├── InsurancePolicyJpaRepository.java
│   │       │   │                               ├── InsurancePolicyRepositoryImpl.java
│   │       │   │                               ├── InsuranceRateJpaRepository.java
│   │       │   │                               ├── InsuranceRateRepositoryImpl.java
│   │       │   │                               ├── InsurerJpaRepository.java
│   │       │   │                               ├── InsurerRepositoryImpl.java
│   │       │   │                               ├── PartnerJpaRepository.java
│   │       │   │                               ├── PartnerRepositoryImpl.java
│   │       │   │                               ├── PaymentCancelJpaRepository.java
│   │       │   │                               ├── PaymentCancelRepositoryImpl.java
│   │       │   │                               ├── PaymentJpaRepository.java
│   │       │   │                               ├── PaymentRepositoryImpl.java
│   │       │   │                               ├── RefreshTokenJpaRepository.java
│   │       │   │                               ├── RefreshTokenRepositoryImpl.java
│   │       │   │                               ├── UserJpaRepository.java
│   │       │   │                               └── UserRepositoryImpl.java
│   │       │   └── resources
│   │       │       ├── db
│   │       │       └── db-core.yml
│   │       └── test
│   │           └── java
│   │               └── com
│   │                   └── nexsol
│   │                       └── tpa
│   │                           └── storage
│   │                               └── db
│   │                                   └── core
│   └── file-core
│       ├── build.gradle
│       └── src
│           └── main
│               └── java
│                   └── com
│                       └── nexsol
│                           └── tpa
│                               └── storage
│                                   └── file
│                                       └── core
│                                           ├── S3FileStorageClient.java
│                                           └── config
│                                               ├── S3ClientConfig.java
│                                               └── S3Properties.java
├── support
│   ├── mailer
│   │   ├── build.gradle
│   │   └── src
│   │       └── main
│   │           └── java
│   │               └── com
│   │                   └── nexsol
│   │                       └── tpa
│   │                           └── support
│   │                               └── mailer
│   │                                   └── SmtpEmailSender.java
│   └── security
│       ├── build.gradle
│       └── src
│           └── main
│               └── java
│                   └── com
│                       └── nexsol
│                           └── tpa
│                               └── support
│                                   └── security
│                                       ├── JwtAuthenticationFilter.java
│                                       ├── JwtProperties.java
│                                       ├── JwtTokenIssuer.java
│                                       └── config
│                                           └── SecurityConfig.java
└── tests
    └── api-docs
        ├── build.gradle
        └── src
            └── main
                └── java
                    └── com
                        └── nexsol
                            └── tpa
                                └── test
                                    └── api
                                        ├── RestDocsTest.java
                                        └── RestDocsUtils.java


```