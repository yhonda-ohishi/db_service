# Tasks: Database-repo for ryohi_sub_cal

**Input**: Design documents from `/specs/master/`
**Prerequisites**: plan.md (✓), research.md (✓), data-model.md (✓), contracts/ (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → ✓ Tech stack: Go 1.21, gRPC v1.58.0, GORM v1.25.5
   → ✓ Structure: Single project (src/, tests/)
2. Load optional design documents:
   → ✓ data-model.md: 3 entities (DTakoUriageKeihi, ETCMeisai, DTakoFerryRows)
   → ✓ contracts/ryohi.proto: 3 services, 15 endpoints total
   → ✓ research.md: GORM with MySQL, godotenv config
3. Generate tasks by category:
   → Setup: Go mod init, dependencies, proto compilation
   → Tests: 15 contract tests, 6 integration tests
   → Core: 3 models, 3 repositories, 3 services
   → Integration: DB connection, config, middleware
   → Polish: unit tests, performance, documentation
4. Apply task rules:
   → Different files marked [P] for parallel
   → Same file sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001-T040)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → ✓ All contracts have tests
   → ✓ All entities have models
   → ✓ All endpoints implemented
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- Paths shown below for single project structure

## Phase 3.1: Setup
- [ ] T001 Create project structure (src/models, src/repository, src/service, src/proto, src/config, tests)
- [ ] T002 Initialize Go module and install dependencies (gRPC, GORM, godotenv, MySQL driver)
- [ ] T003 [P] Configure Protocol Buffers compilation in Makefile
- [ ] T004 [P] Setup .env file from .env.example with database credentials
- [ ] T005 Compile Protocol Buffers: `protoc --go_out=. --go-grpc_out=. src/proto/ryohi.proto`

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests - DTakoUriageKeihiService
- [ ] T006 [P] Contract test Create DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_create.go
- [ ] T007 [P] Contract test Get DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_get.go
- [ ] T008 [P] Contract test Update DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_update.go
- [ ] T009 [P] Contract test Delete DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_delete.go
- [ ] T010 [P] Contract test List DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_list.go

### Contract Tests - ETCMeisaiService
- [ ] T011 [P] Contract test Create ETCMeisai in tests/contract/test_etc_meisai_create.go
- [ ] T012 [P] Contract test Get ETCMeisai in tests/contract/test_etc_meisai_get.go
- [ ] T013 [P] Contract test Update ETCMeisai in tests/contract/test_etc_meisai_update.go
- [ ] T014 [P] Contract test Delete ETCMeisai in tests/contract/test_etc_meisai_delete.go
- [ ] T015 [P] Contract test List ETCMeisai in tests/contract/test_etc_meisai_list.go

### Contract Tests - DTakoFerryRowsService
- [ ] T016 [P] Contract test Create DTakoFerryRows in tests/contract/test_dtako_ferry_rows_create.go
- [ ] T017 [P] Contract test Get DTakoFerryRows in tests/contract/test_dtako_ferry_rows_get.go
- [ ] T018 [P] Contract test Update DTakoFerryRows in tests/contract/test_dtako_ferry_rows_update.go
- [ ] T019 [P] Contract test Delete DTakoFerryRows in tests/contract/test_dtako_ferry_rows_delete.go
- [ ] T020 [P] Contract test List DTakoFerryRows in tests/contract/test_dtako_ferry_rows_list.go

### Integration Tests
- [ ] T021 [P] Integration test composite key operations in tests/integration/test_composite_key.go
- [ ] T022 [P] Integration test transaction rollback in tests/integration/test_transaction.go
- [ ] T023 [P] Integration test bulk insert operations in tests/integration/test_bulk_insert.go
- [ ] T024 [P] Integration test date range queries in tests/integration/test_date_range.go
- [ ] T025 [P] Integration test connection pooling in tests/integration/test_connection_pool.go
- [ ] T026 [P] Integration test environment config loading in tests/integration/test_config.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Configuration & Database
- [ ] T027 Config loader with environment variables in src/config/config.go
- [ ] T028 Database connection manager with GORM in src/config/database.go

### Models
- [ ] T029 [P] DTakoUriageKeihi GORM model with composite key in src/models/dtako_uriage_keihi.go
- [ ] T030 [P] ETCMeisai GORM model with auto-increment in src/models/etc_meisai.go
- [ ] T031 [P] DTakoFerryRows GORM model in src/models/dtako_ferry_rows.go

### Repositories
- [ ] T032 [P] DTakoUriageKeihi repository implementation in src/repository/dtako_uriage_keihi_repo.go
- [ ] T033 [P] ETCMeisai repository implementation in src/repository/etc_meisai_repo.go
- [ ] T034 [P] DTakoFerryRows repository implementation in src/repository/dtako_ferry_rows_repo.go

### gRPC Services
- [ ] T035 DTakoUriageKeihi gRPC service implementation in src/service/dtako_uriage_keihi_service.go
- [ ] T036 ETCMeisai gRPC service implementation in src/service/etc_meisai_service.go
- [ ] T037 DTakoFerryRows gRPC service implementation in src/service/dtako_ferry_rows_service.go

## Phase 3.4: Integration
- [ ] T038 gRPC server setup with all services in cmd/server/main.go
- [ ] T039 Error handling and structured logging middleware in src/middleware/logging.go
- [ ] T040 Database migration runner in cmd/migrate/main.go

## Phase 3.5: Polish
- [ ] T041 [P] Unit tests for model validations in tests/unit/test_model_validation.go
- [ ] T042 [P] Performance test for bulk operations (<200ms) in tests/performance/test_bulk_perf.go
- [ ] T043 [P] Update README.md with setup and usage instructions
- [ ] T044 Run quickstart.md validation scenarios
- [ ] T045 Add health check endpoint for monitoring

## Dependencies
- Setup (T001-T005) must complete first
- Tests (T006-T026) before implementation (T027-T037)
- Config (T027-T028) before models and repositories
- Models (T029-T031) before repositories (T032-T034)
- Repositories before services (T035-T037)
- All services before server setup (T038)
- Implementation before polish (T041-T045)

## Parallel Example
```bash
# Launch contract tests together (T006-T020):
Task: "Contract test Create DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_create.go"
Task: "Contract test Get DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_get.go"
Task: "Contract test Update DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_update.go"
Task: "Contract test Delete DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_delete.go"
Task: "Contract test List DTakoUriageKeihi in tests/contract/test_dtako_uriage_keihi_list.go"
# ... (continue with T011-T020)

# Launch model creation together (T029-T031):
Task: "DTakoUriageKeihi GORM model with composite key in src/models/dtako_uriage_keihi.go"
Task: "ETCMeisai GORM model with auto-increment in src/models/etc_meisai.go"
Task: "DTakoFerryRows GORM model in src/models/dtako_ferry_rows.go"

# Launch repository implementation together (T032-T034):
Task: "DTakoUriageKeihi repository implementation in src/repository/dtako_uriage_keihi_repo.go"
Task: "ETCMeisai repository implementation in src/repository/etc_meisai_repo.go"
Task: "DTakoFerryRows repository implementation in src/repository/dtako_ferry_rows_repo.go"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing (TDD requirement)
- Commit after each task completion
- Environment variables must never be hardcoded
- Use composite key handling for DTakoUriageKeihi
- Maintain <200ms response time requirement

## Task Generation Rules Applied
1. **From Contracts**:
   - ✓ ryohi.proto → 15 contract test tasks (5 per service)
   - ✓ Each endpoint → implementation in service files

2. **From Data Model**:
   - ✓ DTakoUriageKeihi → model with composite key support
   - ✓ ETCMeisai → model with auto-increment
   - ✓ DTakoFerryRows → model with auto-increment

3. **From User Stories**:
   - ✓ CRUD operations → integration tests
   - ✓ Transaction handling → rollback test
   - ✓ Performance requirement → bulk operation test

4. **Ordering**:
   - ✓ Setup → Tests → Models → Services → Integration → Polish
   - ✓ Dependencies properly sequenced

## Validation Checklist
- ✓ All contracts have corresponding tests (T006-T020)
- ✓ All entities have model tasks (T029-T031)
- ✓ All tests come before implementation (Phase 3.2 before 3.3)
- ✓ Parallel tasks truly independent (different files)
- ✓ Each task specifies exact file path
- ✓ No task modifies same file as another [P] task