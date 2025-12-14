# Repository Guidelines

## Project Structure & Module Organization
CityLife is built as a modular C++11 project. Gameplay sources live in `src/`, grouped by feature folders such as `center/`, `controller/`, `food/`, and `world/`; matching headers sit under `include/` with the same directory layout. Build artifacts are emitted to `bin/` and `build/`, while automation scripts and Python acceptance tests (`test_game.py`, `test_game_detailed.py`, `test_game_visual.py`) remain in the repository root.

## Build, Test, and Development Commands
Run `cmake -S . -B build` to configure the project, and `cmake --build build` to compile and place the executable at `bin/main`. Use `./bin/main` for manual play-throughs. Execute `python3 test_game.py` for the default automated journey, `python3 test_game_detailed.py` when you need verbose logging, and `./run_visual_test.sh` to capture a full ncurses session in `game_test_output.txt`.

## Coding Style & Naming Conventions
Match the existing four-space indentation and Allman brace placement (`brace on its own line`). Classes and singletons use `PascalCase` (e.g., `World`, `Controller`), member functions use `camelCase`, and free functions or utilities follow the same pattern. Files and directories stay lowercase with underscores only when readability requires it. Prefer `std::` facilities over raw pointers unless ownership or ncurses integration demands otherwise, and mirror every new `.cpp` with a public header in `include/` when symbols must be shared.
You should put the c++ lib to the .h file, and put this project lib to .cpp file
When you create a new function, you should write the information of it, about the args, return, and more

## Testing Guidelines
Our Python harness relies on `pexpect` to drive the ncurses UI; keep timing-sensitive flows resilient by avoiding unnecessary output changes inside scripted paths. When adding tests, mirror the naming scheme (`test_game_<focus>.py`) and ensure each script exits cleanly via `child.sendcontrol('c')` after assertions. For manual verification, include the resulting snippet from `game_test_output.txt` or terminal captures in the PR description whenever behaviour changes.

## Commit & Pull Request Guidelines
Recent history mixes English and Chinese summaries (e.g., `更改变量传入类型 增加食物类`); keep using a single concise sentence that explains the change. Reference related issues with `#ID` when applicable. Pull requests should summarise the scenario, list the commands run (build + relevant tests), and attach screenshots or log excerpts for gameplay-affecting tweaks. Highlight compatibility considerations such as ncurses version expectations or changes to scripted key sequences.

## Dependencies & Environment
Linux builds link against `ncurses`; Windows contributors must configure PDCurses 3.9 and match the include paths defined in `CMakeLists.txt`. Ensure Python 3 with `pexpect` is available before running automated tests. Document any new external requirement in `readme.md` and update the scripts if the executable path changes from the default `bin/main`.
