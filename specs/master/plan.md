# Implementation Plan: Database-repo for ryohi_sub_cal

**Branch**: `master` | **Date**: 2025-09-19 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/master/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → ✓ Feature spec loaded from C:/go/db_service/specs/master/spec.md
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → ✓ Project Type: Single (gRPC service)
   → ✓ Structure Decision: Option 1 (single project with src/)
3. Fill the Constitution Check section based on the content of the constitution document.
   → ✓ Constitution template loaded (no specific rules defined yet)
4. Evaluate Constitution Check section below
   → ✓ No violations (constitution uses template placeholders)
   → ✓ Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → In progress
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file
7. Re-evaluate Constitution Check section
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

## Summary
Implementation of a gRPC-based Database repository service for managing ryohi_sub_cal database tables (dtako_uriage_keihi, etc_meisai, dtako_ferry_rows) using GORM ORM, with environment-based configuration and no hardcoded secrets.

## Technical Context
**Language/Version**: Go 1.21
**Primary Dependencies**: gRPC v1.58.0, GORM v1.25.5, Protocol Buffers v1.31.0
**Storage**: MySQL/MariaDB (db1 database)
**Testing**: Go testing package
**Target Platform**: Linux/Windows server
**Project Type**: single (gRPC service)
**Performance Goals**: <200ms response time for CRUD operations
**Constraints**: No hardcoded secrets, environment-based configuration
**Scale/Scope**: 3 main tables, ~10 gRPC endpoints

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Since the constitution file contains only template placeholders, no specific gates are defined yet. However, following general best practices:
- ✓ Single responsibility: Data access only
- ✓ Clear contracts: Protocol Buffers definition
- ✓ Environment-based configuration: No hardcoded secrets
- ✓ Standard patterns: Repository pattern with GORM

## Project Structure

### Documentation (this feature)
```
specs/master/
├── plan.md              # This file (/plan command output)
├── spec.md              # Feature specification (created)
├── research.md          # Phase 0 output (to be created)
├── data-model.md        # Phase 1 output (to be created)
├── quickstart.md        # Phase 1 output (to be created)
├── contracts/           # Phase 1 output (to be created)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 1: Single project (SELECTED)
src/
├── models/              # GORM models
├── repository/          # Repository implementations
├── service/             # gRPC service implementations
├── proto/               # Protocol Buffer definitions
└── config/              # Configuration management

tests/
├── contract/            # Contract tests
├── integration/         # Integration tests
└── unit/                # Unit tests
```

**Structure Decision**: Option 1 - Single project structure (gRPC service)

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - ✓ Database schema confirmed from SQL files
   - ✓ Table structures identified (dtako_uriage_keihi, etc_meisai, dtako_ferry_rows)
   - ✓ Technology stack defined (Go, gRPC, GORM)

2. **Generate and dispatch research agents**:
   - Research GORM best practices for composite primary keys
   - Research gRPC service patterns for database operations
   - Research environment configuration patterns in Go

3. **Consolidate findings** in `research.md` using format:
   - Decision: GORM with MySQL driver
   - Rationale: Native support for composite keys, mature ecosystem
   - Alternatives considered: sqlx (lower level), ent (more complex)

**Output**: research.md with all clarifications resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - DTakoUriageKeihi entity with composite key
   - ETCMeisai entity with auto-increment ID
   - DTakoFerryRows entity with auto-increment ID

2. **Generate API contracts** from functional requirements:
   - Protocol Buffers definitions for all entities
   - CRUD service definitions for each table
   - Output to `/contracts/ryohi.proto`

3. **Generate contract tests** from contracts:
   - Test files for each service method
   - Schema validation tests

4. **Extract test scenarios** from user stories:
   - Create and retrieve expense records
   - Handle composite key operations
   - Transaction rollback scenarios

5. **Update agent file incrementally**:
   - CLAUDE.md with project context
   - Technology stack details
   - Recent changes tracking

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, CLAUDE.md

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Generate tasks from contracts and data model
- Each entity → model creation task [P]
- Each service method → implementation task
- Each test scenario → test implementation task

**Ordering Strategy**:
- Config and database connection first
- Models before repositories
- Services after repositories
- Tests throughout (TDD approach)

**Estimated Output**: 20-25 numbered, ordered tasks in tasks.md

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)
**Phase 4**: Implementation (execute tasks.md)
**Phase 5**: Validation (run tests, execute quickstart.md)

## Complexity Tracking
*No violations requiring justification*

## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (none)

---
*Based on Constitution (template) - See `/memory/constitution.md`*