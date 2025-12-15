# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build executable (outputs to bin/citylife)
make build

# Run the game
make run

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

# Test mode with file I/O
./bin/citylife --test tests/input/basic_test.txt tests/output/output.txt
```

**Test Scenarios (run_test.sh):**
- Basic functionality, navigation, shopping, bank operations
- Payment system, save/load, health, hospital, medicine

## Architecture Overview

### Core Loop: `main → tea.Program → ui.Model → action.Execute`

**Execution Flow:**
1. **Entry** (`cmd/citylife/main.go`): Initializes path manager, parses args, creates game state
2. **UI Loop** (`internal/ui/model.go`): Bubble Tea event loop handles keyboard input and renders UI
3. **Actions** (`internal/action/`): Each action implements the Action interface with `Execute()` and `Info()`
4. **State** (`internal/game/state.go`): Aggregates World, User, Nutrition, and Disease managers

### Data Flow & Ownership

**Game State (Composition Pattern):**
- `game.State`: Aggregates all game systems
  - `World`: Time, money, wallet, bank, location management
  - `User`: 20+ nutrition attributes
  - `Nutrition`: Decay rules and interactions
  - `Disease`: 6 disease types with triggers and cures

**Key Packages:**
- `internal/world/`: Time, money, buildings, wallet (6 denominations)
- `internal/user/`: Nutrition management
- `internal/action/`: Action interface and implementations (GoWhere, Commodity, Banking, Medical, Save/Load)
- `internal/ui/`: Bubble Tea UI model
- `internal/food/`: 20+ food items with nutrition effects
- `internal/medicine/`: OTC and prescription medicines
- `internal/checkup/`: Medical checkup packages
- `internal/disease/`: Disease triggers, symptoms, treatments
- `internal/save/`: JSON-based save/load system
- `internal/path/`: Cross-platform path management (~/.citylife/)

### Critical Design Decisions

**1. Bubble Tea UI Framework**
- TUI rendering handled by charmbracelet/bubbletea
- Model-Update-View pattern
- Keyboard navigation: 1-9 for menu items, arrows for scrolling

**2. Action Interface Pattern**
- All actions implement `Action` interface with `Execute()`, `Info()`, `Category()`
- Categories: Primary (shopping, banking), Insight (view status), Navigation (movement)
- Actions sorted by category in menus

**3. Wallet System (6 Denominations)**
- 1, 5, 10, 20, 50, 100 yuan notes
- Greedy payment algorithm with change calculation
- `World.GetWalletTotal()`, `World.SpendMoney()`, `World.AddMoney()`

**4. Save System (JSON)**
- 3 save slots stored in `~/.citylife/saves/`
- Serializes game state including World, User, Disease status

**5. Test Mode**
- `--test [input] [output]` redirects I/O for automation
- Input format: `enter`, `up`, `down`, `1-9`, `q` commands

## File Structure

```
cmd/
└── citylife/main.go          # Entry point

internal/
├── action/                   # Action implementations
│   ├── action.go             # Action interface
│   ├── actions.go            # Basic actions (GoWhere, CheckCash, etc.)
│   ├── banking.go            # Bank operations
│   ├── shopping.go           # Commodity purchases
│   ├── medical.go            # Hospital actions
│   └── save.go               # Save/Load actions
├── game/state.go             # Game state aggregator
├── world/                    # World management
│   ├── world.go              # Time, money, location
│   └── building.go           # Building definitions
├── user/user.go              # User nutrition
├── nutrition/manager.go      # Nutrition decay and interactions
├── disease/manager.go        # Disease system
├── food/                     # Food definitions
├── medicine/medicine.go      # Medicine system
├── checkup/checkup.go        # Medical checkups
├── doctor/doctor.go          # Doctor diagnosis
├── save/manager.go           # Save/load system
├── path/manager.go           # Path utilities
├── ui/model.go               # Bubble Tea UI
├── center/center.go          # City center announcements
├── mapview/mapview.go        # ASCII map rendering
├── locale/locale.go          # i18n support
└── testutil/                 # Test utilities

tests/
├── run_test.sh               # Integration test runner
└── input/                    # Test input files
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
    Info() string                    // Display text
    Execute(state interface{}) tea.Cmd // Execute action
    Category() EventCategory         // For menu sorting
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

- **Go 1.21+**
- **charmbracelet/bubbletea** - TUI framework
- **charmbracelet/lipgloss** - Styling

## Common Pitfalls

1. **Nutrition Decay:** Happens on each action via `Nutrition.OnAction()`
2. **Disease System:** Triggers based on nutrition thresholds
3. **Menu Numbers:** 1-9 directly select items; 10+ requires arrow navigation
4. **Save Path:** Uses `~/.citylife/saves/` cross-platform
5. **Test Input:** Numbers 1-9 don't need Enter; other inputs do

## Commit Guidelines

- Commit messages can use Chinese or English
- Keep messages concise (single sentence explaining the change)
- Reference related issues with `#ID` when applicable
