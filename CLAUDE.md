# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Configure project
cmake -S . -B build

# Build executable (outputs to bin/main)
cmake --build build

# Run manually
./bin/main
```

## Test Commands

```bash
# Basic automated test (recommended for daily regression)
./tests/run_test.sh

# Manual test mode with file I/O
./bin/main --test tests/automated_input.txt tests/automated_output.txt

# Boundary test (extreme inputs, edge cases)
./bin/main --test tests/boundary_input.txt tests/boundary_output.txt
```

**Testing Notes:**
- Center announcements are randomized on each run; test assertions must filter or ignore announcement lines
- `TestModule::ConfigureFileIO` redirects stdin/stdout when `--test` flag is present
- Without `--test`, the program runs with ncurses interactive mode
- Input EOF triggers a `runtime_error("INPUT_EOF")` which exits gracefully (exit code 0)

## Architecture Overview

### Core Loop: `main → World::Start → Controller::Choose → Object::ToDoIt`

**Execution Flow:**
1. **Entry** (`src/main.cpp`): Checks for `--test [input] [output]` CLI args, configures TestModule if present, then calls `World::GetInstance()->Start()`
2. **World Loop** (`src/world/world.cpp`): Infinite loop that calls controller to get player choice at current location (`where`), then executes corresponding `Object::ToDoIt()`
3. **Controller** (`src/controller/controller.cpp`): Renders menu options, validates input, returns selected index
4. **Object Dispatch** (`src/object/object.cpp`): Each action is a subclass of `Object` with its own `ToDoIt()` implementation

### Data Flow & Ownership

**Singletons (Global State):**
- `World`: Owns time, money, bank, location (`where`), all buildings, and action lists (`ToDoThings`)
- `Controller`: Renders menus and validates choices
- `View`: Single point for all I/O (`Print`/`Scan`/`Clear`/`WaitForEnter`)
- `User`: Stores health metrics (hunger, protein, vitamins, etc.), all initialized to 100

**Buildings & Actions:**
- `World::Buildings`: Vector of `Building*` representing locations (city center, market, bank, market interior)
- `World::ToDoThings`: 2D vector mapping location ID to available `Object*` actions
- Actions include: `GoWhere` (travel), `Buy` (enter market), `Information` (view announcements), `DepositingMoney`/`WithdrawMoney`, `Commodity` (purchase item)

**Key Pattern: Dynamic Action Generation**
- `ToDoThings` combines static actions (e.g., view announcements, bank operations) with dynamically generated "Go to X" actions based on location combinations
- Each location gets a set of actions; `Controller::Choose` renders them as numbered options

### Critical Design Decisions

**1. All I/O Centralized in `View` Singleton**
- Rationale: Enables easy UI replacement (console → ncurses → GUI) without touching game logic
- Implementation: `View::Print` uses `std::cout`, `View::Scan` loops until valid integer, `View::WaitForEnter` blocks on `std::getchar()`

**2. Object-Based Action Dispatch**
- Rationale: Each action is self-contained; no giant switch statements
- Pattern: `Object` base class with `virtual void ToDoIt()` and `virtual std::string GetInfo()`
- Derived classes: `GoWhere`, `Buy`, `Information`, `DepositingMoney`, `WithdrawMoney`, `Commodity`

**3. Money Operations Return `bool` for Transaction Validation**
- `World::SpendMoney(int)`: Checks `money >= amount`, returns false on failure
- `World::DepositingMoney(int)`: Deducts from `money`, adds to `bank`, returns false if insufficient funds
- `World::WithdrawMoney(int)`: Checks bank balance, transfers to `money`, returns false if insufficient
- On failure, caller must handle UI feedback (typically `View::Print` + `View::WaitForEnter`)

**4. TestModule Redirects stdio for Automation**
- `ConfigureFileIO(input_path, output_path)` replaces `stdin`/`stdout` with files
- If input file missing, creates it with a single `0` to provide default value for `View::Scan`
- `AppendLog` writes to current stdout for test markers

## File Structure

```
src/
├── main.cpp              # Entry point, TestModule setup
├── world/world.cpp       # Game loop, time/money/location management
├── controller/controller.cpp  # Menu rendering and input validation
├── view/view.cpp         # I/O primitives (Print/Scan/Clear/WaitForEnter)
├── object/object.cpp     # Action implementations (GoWhere, Buy, etc.)
├── center/center.cpp     # Announcement generation and display
├── user/user.cpp         # Player health metrics (hunger, protein, etc.)
└── test/test_module.cpp  # File I/O redirection for automated tests

include/
└── [mirrors src/ structure with .h headers]

tests/
├── run_test.sh           # Automated test runner (filters random announcements)
├── automated_input.txt   # Standard test scenario inputs
├── automated_output.txt  # Expected outputs
├── boundary_input.txt    # Edge case inputs (invalid choices, insufficient funds)
└── boundary_output.txt   # Edge case expected outputs
```

## Coding Conventions (Inferred from Existing Code)

- **Indentation:** 4 spaces (not tabs)
- **Braces:** Allman style (opening brace on own line)
- **Classes/Singletons:** `PascalCase` (e.g., `World`, `Controller`)
- **Member functions:** `camelCase` (e.g., `GetInstance`, `ToDoIt`)
- **C++ Standard:** C++11 (set in CMakeLists.txt)
- **Naming:** Descriptive Chinese comments in headers (`@description`, `@param`, `@return`)
- **Singleton Pattern:** Static `GetInstance()` method, private constructor, static `instance` member

## Key Interfaces

### Object Hierarchy
```cpp
class Object {
    virtual std::string GetInfo();  // Returns display string (e.g., "Apple 5$")
    virtual void ToDoIt();          // Executes action (e.g., updates world state)
};

// Example: GoWhere changes location, updates time, decreases hunger
// Example: Commodity checks payment, updates User health, waits for Enter
```

### World Money API
```cpp
bool SpendMoney(int money);      // Returns false if insufficient funds
bool DepositingMoney(int money); // Transfers money → bank
bool WithdrawMoney(int money);   // Transfers bank → money, checks bank balance
```

### View I/O
```cpp
void Print(const std::string &str);  // Output text
int Scan();                          // Read valid integer (loops on error)
void Clear();                        // Simulate clear screen with newlines
void WaitForEnter();                 // Block until Enter key pressed
```

## Dependencies

- **Linux:** ncurses (linked via CMake)
- **Windows:** PDCurses 3.9 (configure include paths in CMakeLists.txt)
- **Python Tests:** Python 3 with `pexpect` (for driving ncurses UI in tests)

## Common Pitfalls

1. **Announcement Randomization:** Center announcements use `std::random_shuffle`, so output differs between runs. Tests must filter these lines.
2. **Input Validation:** `View::Scan` loops infinitely until valid integer received. In test mode, ensure input files have enough valid numbers.
3. **EOF Handling:** `View::Scan` throws `runtime_error("INPUT_EOF")` when input exhausted. Main catches this and exits with code 0.
4. **Money vs Bank:** Player has two balances: `World::money` (wallet) and `Bank::deposit` (savings). Operations must specify which to use.
5. **Health Clamping:** `User` health values capped at 100 via `std::min(current + delta, 100)` when consuming `Commodity`.
