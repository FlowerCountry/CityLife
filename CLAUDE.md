# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build executable (outputs to bin/citylife)
make build

# Run the API server (default port 8080)
make run

# Run with custom port
./bin/citylife --port 3000

# Run unit tests
make test

# Format code
make fmt

# Run go vet
make vet

# Cross-platform builds (Linux/macOS/Windows)
make build-all
```

## Test Commands

```bash
# Run all unit tests
make test

# Run integration tests
./tests/run_test.sh
```

**Test Scenarios (run_test.sh):**
- Basic functionality, navigation, shopping, bank operations
- Payment system, save/load, health, hospital, medicine

## Architecture Overview

### Core Loop: `main → gin.Engine → handler → action.Execute`

**Execution Flow:**
1. **Entry** (`cmd/citylife/main.go`): Initializes path manager, creates session manager, registers Gin routes
2. **Router** (`internal/api/router.go`): Registers all API endpoints with middleware
3. **Handlers** (`internal/api/handler/`): Process HTTP requests, interact with game state
4. **Actions** (`internal/action/`): Each action implements the Action interface with `Execute()` returning `*Result`
5. **State** (`internal/game/state.go`): Aggregates World, User, Nutrition, and Disease managers

### Data Flow & Ownership

**Game State (Composition Pattern):**
- `game.State`: Aggregates all game systems
  - `World`: Time, money, wallet, bank, location management
  - `User`: 20+ nutrition attributes
  - `Nutrition`: Decay rules and interactions
  - `Disease`: 6 disease types with triggers and cures

**Key Packages:**
- `internal/api/`: REST API layer (Gin handlers, middleware, response)
- `internal/session/`: Session management (30-minute timeout)
- `internal/action/`: Action interface and implementations
- `internal/executor/`: Action execution orchestration
- `internal/world/`: Time, money, buildings, wallet (6 denominations)
- `internal/user/`: Nutrition management
- `internal/payment/`: Wallet and payment processing
- `internal/food/`: 20+ food items with nutrition effects
- `internal/medicine/`: OTC and prescription medicines
- `internal/checkup/`: Medical checkup packages
- `internal/disease/`: Disease triggers, symptoms, treatments
- `internal/save/`: JSON-based save/load system
- `internal/path/`: Cross-platform path management (~/.citylife/)

### Critical Design Decisions

**1. Gin REST API Framework**
- RESTful API design with JSON responses
- Session-based state management (UUID session IDs)
- CORS middleware for cross-origin requests
- Route group: `/api/v2/sessions/:session_id/*`

**2. Action Interface Pattern**
- All actions implement `Action` interface with `ID()`, `Info()`, `Execute()`, `Category()`
- Categories: Primary (shopping, banking), Insight (view status), Navigation (movement)
- `Execute()` returns `*Result` with Message, Success, TimeElapsed

**3. Wallet System (6 Denominations)**
- 1, 5, 10, 20, 50, 100 yuan notes
- Greedy payment algorithm with change calculation
- `World.GetWalletTotal()`, `World.SpendMoney()`, `World.AddMoney()`

**4. Save System (JSON)**
- 3 save slots stored in `~/.citylife/saves/`
- Serializes game state including World, User, Disease status

**5. Session Management**
- UUID-based session IDs
- 30-minute idle timeout
- Thread-safe session storage

## API Endpoints

```
POST   /api/v2/sessions                         Create game session
GET    /api/v2/sessions/:id                     Get session info
DELETE /api/v2/sessions/:id                     Delete session
GET    /api/v2/sessions/:id/state               Get full game state
GET    /api/v2/sessions/:id/status              Get player status
GET    /api/v2/sessions/:id/actions             Get available actions
POST   /api/v2/sessions/:id/actions/:action_id  Execute action
POST   /api/v2/sessions/:id/navigate/:loc       Navigate to location
POST   /api/v2/sessions/:id/bank/*              Bank operations
GET    /api/v2/sessions/:id/shop/*              Shopping (commodities)
POST   /api/v2/sessions/:id/shop/*              Shopping (buy)
GET    /api/v2/sessions/:id/hospital/*          Medical (lists)
POST   /api/v2/sessions/:id/hospital/*          Medical (actions)
GET    /api/v2/sessions/:id/saves               Save slots
POST   /api/v2/sessions/:id/saves/*             Save/Load
GET    /api/v2/sessions/:id/housing/*           Housing (status/offers)
POST   /api/v2/sessions/:id/housing/*           Housing actions
GET    /api/v2/sessions/:id/jobs                Job list
POST   /api/v2/sessions/:id/jobs/:id            Do job
GET    /api/v2/sessions/:id/restaurant/*        Restaurant menu
POST   /api/v2/sessions/:id/restaurant/:id      Restaurant action
GET    /api/v2/sessions/:id/park/*              Park activities
POST   /api/v2/sessions/:id/park/:id            Park action
GET    /api/v2/sessions/:id/hotel/*             Hotel services
POST   /api/v2/sessions/:id/hotel/:id           Hotel action
GET    /api/v2/sessions/:id/wallet              Wallet detail
GET    /api/v2/sessions/:id/bank                Bank balance
GET    /api/v2/sessions/:id/health              Health detail
GET    /api/v2/sessions/:id/diseases            Disease status
GET    /api/v2/sessions/:id/map                 Map locations
```

## File Structure

```
cmd/
└── citylife/main.go          # Entry point (Gin server)

internal/
├── api/                      # REST API layer
│   ├── router.go             # Route registration
│   ├── handler/              # Request handlers
│   │   └── v2/                # v2 API handlers
│   ├── middleware/           # CORS, session validation
│   └── response/             # Unified JSON responses
├── session/manager.go        # Session lifecycle
├── action/                   # Action implementations
│   ├── action.go             # Action interface
│   └── actions.go            # Concrete actions
├── executor/                 # Action execution
├── game/state.go             # Game state aggregator
├── world/                    # World management
│   ├── world.go              # Time, money, location
│   └── building.go           # Building definitions
├── user/user.go              # User nutrition
├── nutrition/manager.go      # Nutrition decay
├── disease/manager.go        # Disease system
├── payment/                  # Payment processing
│   ├── cashier.go            # Transaction handling
│   └── wallet_inspector.go   # Wallet queries
├── food/                     # Food definitions
├── medicine/medicine.go      # Medicine system
├── checkup/checkup.go        # Medical checkups
├── doctor/doctor.go          # Doctor diagnosis
├── save/manager.go           # Save/load system
├── path/manager.go           # Path utilities
├── center/center.go          # City center
├── mapview/mapview.go        # ASCII map
├── locale/locale.go          # i18n support
└── testutil/                 # Test utilities

tests/
├── run_test.sh               # Integration test runner
└── integration/              # Integration tests
```

## Coding Conventions

- **Indentation:** Tabs (Go standard)
- **Naming:** Go conventions (CamelCase for exported, camelCase for internal)
- **Packages:** Single-purpose packages in `internal/`
- **Interfaces:** Small, focused interfaces (e.g., `Action`)
- **Testing:** Table-driven tests with `*_test.go` files
- **Comments:** Chinese comments acceptable, `@description`, `@param`, `@return` annotations

## Key Interfaces

### Action Interface
```go
type Action interface {
    ID() string                      // Unique identifier
    Info() string                    // Display text
    Execute(state *game.State) *Result // Execute action
    Category() EventCategory         // For sorting
}

type Result struct {
    Message     string
    Success     bool
    TimeElapsed int
}
```

### World Money API
```go
func (w *World) GetWalletTotal() int
func (w *World) SpendMoney(amount int) bool
func (w *World) AddMoney(amount int)
func (w *World) DepositToBank(amount int) bool
func (w *World) WithdrawFromBank(amount int) bool
```

### User Nutrition API
```go
func (u *User) GetNutrition(name string) int
func (u *User) AddNutrition(name string, amount int)
func (u *User) ConsumeNutrition(name string, amount int)
func (u *User) IsAlive() bool
```

## Dependencies

- **Go 1.23+**
- **github.com/gin-gonic/gin** - REST API framework
- **github.com/google/uuid** - Session ID generation

## Common Pitfalls

1. **Nutrition Decay:** Happens on each action via `Nutrition.OnAction()`
2. **Disease System:** Triggers based on nutrition thresholds
3. **Session Timeout:** Sessions expire after 30 minutes of inactivity
4. **Save Path:** Uses `~/.citylife/saves/` cross-platform
5. **CORS:** Enabled by default for all origins

## Commit Guidelines

- Commit messages can use Chinese or English
- Keep messages concise (single sentence explaining the change)
- Reference related issues with `#ID` when applicable
