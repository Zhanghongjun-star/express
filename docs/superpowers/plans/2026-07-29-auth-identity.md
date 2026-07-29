# Auth Identity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the identity/authentication business from the technical review: verification code, register, login, refresh, logout, password change, and real-name authentication.

**Architecture:** Follow the existing Kratos layering. `service` converts proto DTOs, `biz` owns business rules and repository interfaces, `data` owns MySQL/Redis access, and `model` documents the table shapes. SQL is written with `database/sql`; Redis stores short-lived verification codes, refresh sessions, access-token blacklist entries, and password-version invalidation hints.

**Tech Stack:** Go, Kratos v3, Protobuf/buf, MySQL via `github.com/go-sql-driver/mysql`, Redis via `github.com/redis/go-redis/v9`, bcrypt via `golang.org/x/crypto/bcrypt`.

---

### Task 1: Auth API Contract

**Files:**
- Create: `api/auth/v1/auth.proto`
- Create: `api/auth/v1/error_reason.proto`
- Modify generated files through `make api`

- [ ] Write proto for review endpoints under `/api/v1/auth` and `/api/v1/users/me/real-name-auth`.
- [ ] Generate API code with `make api`.
- [ ] Keep generated files unedited.

### Task 2: Biz Rules

**Files:**
- Create: `internal/biz/auth.go`
- Test: `internal/biz/auth_test.go`

- [ ] Write failing tests for register duplicate account, password length, login lockout after five failures, password change rejecting old password mismatch, and real-name ID-card validation.
- [ ] Implement simple domain structs and `AuthRepo` interface.
- [ ] Implement `AuthUsecase` methods with readable validation helpers and short comments per module.
- [ ] Run `go test ./internal/biz`.

### Task 3: Data Access

**Files:**
- Modify: `internal/data/data.go`
- Create: `internal/data/auth.go`
- Create: `internal/data/auth_schema.sql`

- [ ] Add MySQL and Redis clients to `Data`.
- [ ] Open clients from `configs/config.yaml` through `internal/conf`.
- [ ] Implement SQL repository methods for the three identity tables.
- [ ] Implement Redis helpers for verification code, session, blacklist, and login failure TTL.
- [ ] Run `go test ./internal/data`.

### Task 4: Service Wiring

**Files:**
- Create: `internal/service/auth.go`
- Test: `internal/service/auth_test.go`
- Modify: `internal/biz/biz.go`
- Modify: `internal/data/data.go`
- Modify: `internal/service/service.go`
- Modify: `internal/server/http.go`
- Modify: `internal/server/grpc.go`
- Modify generated Wire code through `go generate ./cmd/server`

- [ ] Write service tests for register/login/real-name happy paths using a fake repo.
- [ ] Implement DTO to DO conversion only in service.
- [ ] Register Auth service in HTTP and gRPC servers.
- [ ] Regenerate Wire.
- [ ] Run `go test ./internal/service`.

### Task 5: Final Verification

**Files:**
- Whole repository

- [ ] Run `go test ./...`.
- [ ] Run `go build ./...`.
- [ ] Run `go run ./cmd/server -conf ./configs` long enough to confirm config loading and startup, then stop it.
- [ ] Report exactly which commands passed or failed.
