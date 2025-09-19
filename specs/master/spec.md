# Feature Specification: Database-repo Implementation

## Overview
Implementation of a Database repository service following gRPC-First architecture pattern for managing ryohi_sub_cal database tables.

## User Stories
1. As a developer, I want to access database operations through gRPC API
2. As a system, I need to manage dtako_uriage_keihi and etc_meisai tables
3. As an operator, I need secure database access without hardcoded credentials

## Functional Requirements
1. **gRPC Service Implementation**
   - Protocol Buffers contract definition
   - CRUD operations for dtako_uriage_keihi table
   - CRUD operations for etc_meisai table

2. **Database Tables**
   - `dtako_uriage_keihi`: Expense tracking table with composite primary key
   - `etc_meisai`: ETC toll details with auto-increment ID

3. **Data Access Layer**
   - GORM-based repository implementation
   - Connection pooling and management
   - Transaction support

## Non-Functional Requirements
1. **Security**
   - No hardcoded secrets
   - Environment variable configuration
   - Secure connection strings

2. **Architecture**
   - Single responsibility: Data access and schema management only
   - gRPC-First approach
   - Clean separation of concerns

## Technical Context
Based on https://github.com/yhonda-ohishi/db-handler-server/blob/main/README.md:
- Use GORM for ORM
- golang-migrate for schema management
- Protocol Buffers for contract definition
- bufconn for in-memory gRPC communication

## Database Schema (from C:\go\ryohi_sub_cal\mysqldata)
### dtako_uriage_keihi
- Primary Key: srch_id, datetime, keihi_c
- Contains expense tracking data with relations to dtako_row

### etc_meisai
- Primary Key: id (auto-increment)
- Contains ETC toll transaction details

## Success Criteria
1. Functional gRPC server on port 50051
2. All CRUD operations working for both tables
3. Environment-based configuration
4. No hardcoded credentials
5. Proper error handling and logging