# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Freelance platform (程序员接单平台) — a marketplace where clients post projects and developers bid on them. Supports escrow payments, milestone-based releases, real-time chat, reviews, and admin management.

## Build & Run Commands

### Server (Go)
```bash
cd server
go mod tidy
go run cmd/server/main.go          # dev server on :8080
go build ./cmd/server              # build binary
```
Requires PostgreSQL with database `freelancedb`. Run migrations: `psql -U postgres -d freelancedb -f server/migrations/000001_init_schema.up.sql`

### Web (Vue 3)
```bash
cd web
npm install
npm run dev                        # dev server, Vite proxies /api → localhost:8080
npm run build                      # vue-tsc type check + vite build
npm run test:e2e                   # Playwright E2E tests
npm run test:e2e -- --grep "test name"  # run single E2E test
```

### Miniapp (Taro)
```bash
cd miniapp
npm install
npm run dev:weapp                  # watch mode, output to dist/
npm run build:weapp                # production build
```
Open `dist/` in WeChat DevTools.

### Android
Open `android/` in Android Studio. Sync Gradle, then run. `./gradlew assembleDebug` from CLI.

### Docker
```bash
docker-compose up --build          # server (:8080) + web (:3002)
```

## Architecture

### Server — Clean Architecture (Go/Gin)
```
cmd/server/main.go     → entrypoint, wires all dependencies
internal/config/       → Viper config from config.yaml + env overrides
internal/model/        → data models, DTOs, request/response structs
internal/repository/   → pgx5 PostgreSQL data access (21 files, one per domain)
internal/service/      → business logic (18 files), coordinates repositories
internal/handler/      → Gin HTTP handlers (16 files), calls services
internal/middleware/    → JWT auth, CORS, RequireAdmin, RequireUserType
internal/ws/           → gorilla/websocket hub for real-time chat
```

**Request flow**: Handler → Service → Repository → pgx5 pool. Handlers use shared `response.go` helpers (`Success`, `BadRequest`, `Unauthorized`, etc.) returning `{code, message, data, meta}`.

**Config**: `server/config.yaml` with env var overrides (e.g. `DATABASE_HOST`). JWT secret, expiry, storage path, WeChat credentials.

**Auth**: JWT access tokens (15min) + refresh tokens (7d). Middleware sets user context on `gin.Context`.

**WebSocket**: Hub-and-client pattern at `WS /ws`. Auth via JWT in query param. Supports heartbeats, direct messages, broadcast.

### Web — Vue 3 + Element Plus + Pinia
```
src/api/index.ts       → Axios instance, base URL /api/v1, auto token refresh on 401
src/api/modules.ts     → Domain API modules (auth, user, project, bid, contract, etc.)
src/store/             → Pinia stores (user.ts, chat.ts)
src/router/index.ts    → Route guards: guest/auth/role meta fields
src/views/             → Pages by domain (auth/, project/, developer/, chat/, admin/, user/)
src/styles/global.scss → CSS variables, utility classes
```

**Auto-imports**: Element Plus components and Vue APIs via `unplugin-auto-import` / `unplugin-vue-components`. No manual imports needed for Element Plus or Vue composition APIs.

**Dev proxy**: Vite proxies `/api` → `http://localhost:8080` in dev mode.

### Miniapp — Taro 4 + React
```
src/services/api.ts    → Centralized request function with JWT, 401 refresh queue, loading indicators
src/services/auth.ts   → WeChat login (wx.login → code exchange), silent login, logout
src/store/user.ts      → React Context + useUser() hook, persists to Taro storage
src/pages/             → 14+ pages by domain
src/components/        → ProjectCard, DeveloperCard, Empty
```

**Auth flow**: WeChat `wx.login()` → server code exchange → JWT tokens stored in Taro storage. Miniapp also supports email/phone login.

### Android — MVVM + Jetpack Compose
```
data/api/ApiService.kt       → Retrofit interface, 30+ endpoints
data/api/RetrofitClient.kt   → Singleton OkHttp + Retrofit
data/api/TokenInterceptor.kt → Auto-auth header, 401 refresh
data/local/TokenManager.kt   → EncryptedSharedPreferences for JWT
data/repository/             → AuthRepository, ProjectRepository
ui/auth/                     → LoginScreen, RegisterScreen, AuthViewModel
```

**DI**: Hilt. **Navigation**: NavHost with bottom tabs (Home, Projects, Chat, Profile).

## Key Domain Concepts

- **Dual roles**: Client (需求方) and Developer (开发者). Some routes require specific role via `RequireUserType` middleware.
- **Project lifecycle**: Draft → Published → Bidding → In Progress → Delivered → Completed
- **Escrow payments**: Client deposits → platform holds → released per milestone acceptance. 10% platform commission.
- **API prefix**: All REST endpoints at `/api/v1/`, ~95 endpoints total.
- **API response format**: `{code: 0, message: "success", data: ..., meta: {page, page_size, total}}`. Code 0 = success, non-zero = error (40000, 40100, 40300, 40400, 50000).
- **Database**: 25 tables, UUID primary keys. Core: users, developer_profiles, client_profiles, projects, bids, contracts, payments, platform_wallets, chat_conversations, chat_messages, reviews, notifications, disputes.

## Conventions

- Server follows Go standard layout: `internal/` prevents external imports, no business logic in handlers.
- All four clients (web, miniapp, android) share the same REST API contract.
- The miniapp uses BEM naming for CSS classes (e.g., `project-card__header`).
- Web uses Element Plus Chinese locale by default.
