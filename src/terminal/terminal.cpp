#include "terminal/terminal.h"

#include <cctype>
#include <cstdio>
#include <iostream>
#include <vector>
#include <sstream>

#ifndef _WIN32
#include <unistd.h>
#include <sys/ioctl.h>
#endif

Terminal::KeyEvent::KeyEvent(Type type, char ch)
    : type(type), ch(ch)
{
}

Terminal::Terminal()
#ifdef _WIN32
    : rawModeEnabled(false), consoleInputHandle(INVALID_HANDLE_VALUE), originalConsoleMode(0)
#else
    : rawModeEnabled(false)
#endif
{
}

Terminal::~Terminal()
{
    disableRawMode();
}

void Terminal::clearOutputLog()
{
    outputLog.clear();
}

void Terminal::writeRaw(const std::string &text)
{
    outputLog.push_back(text);
    std::cout << text;
    std::cout.flush();
}

void Terminal::writeLine(const std::string &line)
{
    writeRaw(line + "\n");
}

const std::vector<std::string> &Terminal::getOutputLog() const
{
    return outputLog;
}

void Terminal::enableRawMode()
{
    if (rawModeEnabled)
    {
        return;
    }
#ifdef _WIN32
    consoleInputHandle = GetStdHandle(STD_INPUT_HANDLE);
    if (consoleInputHandle == INVALID_HANDLE_VALUE)
    {
        throw std::runtime_error("GetStdHandle failed");
    }

    if (!GetConsoleMode(consoleInputHandle, &originalConsoleMode))
    {
        throw std::runtime_error("GetConsoleMode failed");
    }

    DWORD rawMode = originalConsoleMode;
    rawMode &= ~(ENABLE_LINE_INPUT | ENABLE_ECHO_INPUT);
    rawMode |= ENABLE_EXTENDED_FLAGS;
    rawMode &= ~ENABLE_QUICK_EDIT_MODE;

    if (!SetConsoleMode(consoleInputHandle, rawMode))
    {
        throw std::runtime_error("SetConsoleMode failed");
    }

    FlushConsoleInputBuffer(consoleInputHandle);
    rawModeEnabled = true;
#else
    if (!isatty(STDIN_FILENO))
    {
        throw std::runtime_error("STDIN is not a terminal");
    }
    if (tcgetattr(STDIN_FILENO, &originalTermios) == -1)
    {
        throw std::runtime_error("tcgetattr failed");
    }

    struct termios raw = originalTermios;
    raw.c_lflag &= ~(ICANON | ECHO);
    raw.c_cc[VMIN] = 1;
    raw.c_cc[VTIME] = 0;

    if (tcsetattr(STDIN_FILENO, TCSANOW, &raw) == -1)
    {
        throw std::runtime_error("tcsetattr failed");
    }
    rawModeEnabled = true;
#endif
}

void Terminal::disableRawMode()
{
    if (!rawModeEnabled)
    {
        return;
    }
#ifdef _WIN32
    if (consoleInputHandle != INVALID_HANDLE_VALUE)
    {
        SetConsoleMode(consoleInputHandle, originalConsoleMode);
    }
    consoleInputHandle = INVALID_HANDLE_VALUE;
    rawModeEnabled = false;
#else
    tcsetattr(STDIN_FILENO, TCSANOW, &originalTermios);
    rawModeEnabled = false;
#endif
}

Terminal::KeyEvent Terminal::readKey()
{
    if (!rawModeEnabled)
    {
        enableRawMode();
    }
#ifdef _WIN32
    INPUT_RECORD record;
    DWORD eventsRead = 0;

    while (true)
    {
        if (!ReadConsoleInput(consoleInputHandle, &record, 1, &eventsRead))
        {
            return KeyEvent(KeyEvent::Type::EndOfInput);
        }
        if (eventsRead == 0)
        {
            continue;
        }
        if (record.EventType != KEY_EVENT)
        {
            continue;
        }

        const KEY_EVENT_RECORD &key = record.Event.KeyEvent;
        if (!key.bKeyDown)
        {
            continue;
        }

        switch (key.wVirtualKeyCode)
        {
        case VK_UP:
            return KeyEvent(KeyEvent::Type::ArrowUp);
        case VK_DOWN:
            return KeyEvent(KeyEvent::Type::ArrowDown);
        case VK_LEFT:
            return KeyEvent(KeyEvent::Type::ArrowLeft);
        case VK_RIGHT:
            return KeyEvent(KeyEvent::Type::ArrowRight);
        case VK_RETURN:
            return KeyEvent(KeyEvent::Type::Enter, '\n');
        default:
            break;
        }

        const WCHAR unicode = key.uChar.UnicodeChar;
        if (unicode == 3)
        {
            return KeyEvent(KeyEvent::Type::EndOfInput);
        }

        if (unicode >= 32 && unicode <= 126)
        {
            return KeyEvent(KeyEvent::Type::Character, static_cast<char>(unicode));
        }

        if (unicode == 27)
        {
            return KeyEvent(KeyEvent::Type::Unknown);
        }
    }
#else
    int ch = readSingleUnsafe();
    if (ch == -1)
    {
        return KeyEvent(KeyEvent::Type::EndOfInput);
    }
    if (ch == 27)
    {
        return decodeEscape();
    }
    if (ch == '\r' || ch == '\n')
    {
        return KeyEvent(KeyEvent::Type::Enter, '\n');
    }
    return normalizePrintable(ch);
#endif
}

int Terminal::runMenu(const std::vector<std::string> &options,
                      const std::function<void(int, const std::string &)> &render,
                      std::string &typedBuffer)
{
    // 默认垂直导航，与历史行为保持一致。
    return runMenuEx(options, render, typedBuffer, NavAxis::Vertical, true);
}

int Terminal::runMenuEx(const std::vector<std::string> &options,
                        const std::function<void(int, const std::string &)> &render,
                        std::string &typedBuffer,
                        NavAxis axis,
                        bool wrapAround,
                        int startIndex)
{
    if (options.empty())
    {
        throw std::runtime_error("Terminal::runMenu requires non-empty options");
    }
    if (!render)
    {
        throw std::runtime_error("Terminal::runMenu requires a render callback");
    }

    // 以 RAII 管理原始模式，避免多处分支恢复。
    RawSessionGuard guard(*this);

    typedBuffer.clear();
    int index = 0;
    const auto clamp = [&](int v) {
        const int n = static_cast<int>(options.size());
        if (wrapAround)
        {
            return (v % n + n) % n;
        }
        if (v < 0) return 0;
        if (v >= n) return n - 1;
        return v;
    };

    index = clamp(startIndex);
    render(index, typedBuffer);

    while (true)
    {
        KeyEvent event = readKey();
        bool requestRepaint = false;

        switch (event.type)
        {
        case KeyEvent::Type::ArrowUp:
            if (axis == NavAxis::Vertical || axis == NavAxis::Both)
            {
                index = clamp(index - 1);
                requestRepaint = true;
            }
            break;
        case KeyEvent::Type::ArrowDown:
            if (axis == NavAxis::Vertical || axis == NavAxis::Both)
            {
                index = clamp(index + 1);
                requestRepaint = true;
            }
            break;
        case KeyEvent::Type::ArrowLeft:
            if (axis == NavAxis::Horizontal || axis == NavAxis::Both)
            {
                index = clamp(index - 1);
                requestRepaint = true;
            }
            break;
        case KeyEvent::Type::ArrowRight:
            if (axis == NavAxis::Horizontal || axis == NavAxis::Both)
            {
                index = clamp(index + 1);
                requestRepaint = true;
            }
            break;
        case KeyEvent::Type::Character:
            typedBuffer.push_back(event.ch);
            requestRepaint = true;
            break;
        case KeyEvent::Type::Enter:
            return index;
        case KeyEvent::Type::EndOfInput:
            throw std::runtime_error("INPUT_EOF");
        case KeyEvent::Type::Unknown:
            break;
        }

        if (requestRepaint)
        {
            render(index, typedBuffer);
        }
    }
}

int Terminal::readSingleUnsafe()
{
#ifndef _WIN32
    unsigned char ch = 0;
    const ssize_t readBytes = ::read(STDIN_FILENO, &ch, 1);
    if (readBytes <= 0)
    {
        return -1;
    }
    return static_cast<int>(ch);
#else
    return -1;
#endif
}

Terminal::KeyEvent Terminal::decodeEscape()
{
#ifndef _WIN32
    int second = readSingleUnsafe();
    if (second == -1)
    {
        return KeyEvent(KeyEvent::Type::EndOfInput);
    }
    if (second == '[' || second == 'O')
    {
        int third = readSingleUnsafe();
        if (third == -1)
        {
            return KeyEvent(KeyEvent::Type::EndOfInput);
        }
        if (third == 'A')
        {
            return KeyEvent(KeyEvent::Type::ArrowUp);
        }
        if (third == 'B')
        {
            return KeyEvent(KeyEvent::Type::ArrowDown);
        }
        if (third == 'C')
        {
            return KeyEvent(KeyEvent::Type::ArrowRight);
        }
        if (third == 'D')
        {
            return KeyEvent(KeyEvent::Type::ArrowLeft);
        }
        return KeyEvent(KeyEvent::Type::Unknown);
    }
    if (std::isprint(static_cast<unsigned char>(second)) != 0)
    {
        return KeyEvent(KeyEvent::Type::Character, static_cast<char>(second));
    }
    return KeyEvent(KeyEvent::Type::Unknown);
#else
    return KeyEvent(KeyEvent::Type::Unknown);
#endif
}

Terminal::KeyEvent Terminal::normalizePrintable(int ch)
{
    if (ch == -1)
    {
        return KeyEvent(KeyEvent::Type::EndOfInput);
    }
    if (ch == '\b' || ch == 127)
    {
        return KeyEvent(KeyEvent::Type::Unknown);
    }
    if (std::isprint(static_cast<unsigned char>(ch)) != 0)
    {
        return KeyEvent(KeyEvent::Type::Character, static_cast<char>(ch));
    }
    return KeyEvent(KeyEvent::Type::Unknown);
}

Terminal::WindowSize Terminal::getWindowSize() const
{
    WindowSize size{80, 24};
#ifdef _WIN32
    CONSOLE_SCREEN_BUFFER_INFO info;
    HANDLE handle = GetStdHandle(STD_OUTPUT_HANDLE);
    if (handle != INVALID_HANDLE_VALUE && GetConsoleScreenBufferInfo(handle, &info))
    {
        size.columns = static_cast<int>(info.srWindow.Right - info.srWindow.Left + 1);
        size.rows = static_cast<int>(info.srWindow.Bottom - info.srWindow.Top + 1);
    }
#else
    struct winsize ws
    {
    };
    if (ioctl(STDOUT_FILENO, TIOCGWINSZ, &ws) == 0)
    {
        if (ws.ws_col > 0)
        {
            size.columns = ws.ws_col;
        }
        if (ws.ws_row > 0)
        {
            size.rows = ws.ws_row;
        }
    }
#endif
    return size;
}

Terminal::GridLayout Terminal::computeGridLayout(const std::vector<std::string> &items, int padding, int maxColumns) const
{
    GridLayout layout{};
    layout.itemCount = static_cast<int>(items.size());
    if (items.empty())
    {
        layout.columns = 0;
        layout.rows = 0;
        layout.columnWidth = 0;
        return layout;
    }

    int maxLen = 0;
    for (const auto &item : items)
    {
        maxLen = std::max(maxLen, static_cast<int>(item.size()));
    }

    const int baseWidth = std::max(1, maxLen + padding);
    WindowSize ws = getWindowSize();
    int cols = ws.columns / baseWidth;
    if (cols <= 0)
    {
        cols = 1;
    }
    if (maxColumns > 0)
    {
        cols = std::min(cols, std::max(1, maxColumns));
    }
    layout.columns = cols;
    layout.columnWidth = baseWidth;
    layout.rows = (layout.itemCount + cols - 1) / cols;
    return layout;
}

Terminal::MenuGridResult Terminal::runMenuGrid(const MenuGridOption &option,
                                               const std::function<void(int, int, const std::vector<std::vector<std::string>> &, const std::vector<int> &)> &render)
{
    if (option.rows.empty())
    {
        throw std::runtime_error("Terminal::runMenuGrid requires non-empty rows");
    }

    std::size_t maxCols = 0;
    for (const auto &row : option.rows)
    {
        maxCols = std::max(maxCols, row.size());
    }
    if (maxCols == 0)
    {
        throw std::runtime_error("Terminal::runMenuGrid requires at least one option");
    }

    std::vector<int> columnWidths(maxCols, 0);
    for (const auto &row : option.rows)
    {
        for (std::size_t col = 0; col < row.size(); ++col)
        {
            const int width = static_cast<int>(row[col].size()) + 2; // prefix width
            columnWidths[col] = std::max(columnWidths[col], width);
        }
    }
    for (std::size_t col = 0; col < columnWidths.size(); ++col)
    {
        if (col + 1 < columnWidths.size())
        {
            columnWidths[col] += 2; // gap between columns
        }
        if (columnWidths[col] <= 0)
        {
            columnWidths[col] = 2;
        }
    }

    auto isRowValid = [&](int rowIndex) {
        return rowIndex >= 0 && rowIndex < static_cast<int>(option.rows.size()) && !option.rows[static_cast<std::size_t>(rowIndex)].empty();
    };

    MenuGridResult result;
    int focusRow = 0;
    while (focusRow < static_cast<int>(option.rows.size()) && option.rows[static_cast<std::size_t>(focusRow)].empty())
    {
        ++focusRow;
    }
    if (focusRow >= static_cast<int>(option.rows.size()))
    {
        throw std::runtime_error("Terminal::runMenuGrid requires at least one non-empty row");
    }
    int focusCol = 0;

    auto clampCol = [&](int rowIndex, int desiredCol) {
        const auto &row = option.rows[static_cast<std::size_t>(rowIndex)];
        const int size = static_cast<int>(row.size());
        if (size == 0)
        {
            return -1;
        }
        if (option.wrapCols)
        {
            int wrapped = desiredCol % size;
            if (wrapped < 0)
            {
                wrapped += size;
            }
            return wrapped;
        }
        return std::max(0, std::min(desiredCol, size - 1));
    };

    auto renderFrame = [&]() {
        if (render)
        {
            render(focusRow, focusCol, option.rows, columnWidths);
            return;
        }

        std::ostringstream frame;
        frame << "\033[H\033[2J";
        for (std::size_t row = 0; row < option.rows.size(); ++row)
        {
            for (std::size_t col = 0; col < maxCols; ++col)
            {
                std::string cell;
                if (col < option.rows[row].size())
                {
                    const bool active = static_cast<int>(row) == focusRow && static_cast<int>(col) == focusCol;
                    cell = (active ? "> " : "  ") + option.rows[row][col];
                }
                if (static_cast<int>(cell.size()) < columnWidths[col])
                {
                    cell += std::string(columnWidths[col] - static_cast<int>(cell.size()), ' ');
                }
                frame << cell;
            }
            frame << '\n';
        }
        writeRaw(frame.str());
    };

    RawSessionGuard guard(*this);
    renderFrame();

    auto moveRow = [&](int delta) {
        int preferredCol = focusCol;
        int attempts = static_cast<int>(option.rows.size());
        int newRow = focusRow;
        while (attempts-- > 0)
        {
            newRow += delta;
            if (delta > 0 && newRow >= static_cast<int>(option.rows.size()))
            {
                if (option.wrapRows)
                {
                    newRow = 0;
                }
                else
                {
                    newRow = static_cast<int>(option.rows.size()) - 1;
                    break;
                }
            }
            else if (delta < 0 && newRow < 0)
            {
                if (option.wrapRows)
                {
                    newRow = static_cast<int>(option.rows.size()) - 1;
                }
                else
                {
                    newRow = 0;
                    break;
                }
            }

            if (isRowValid(newRow))
            {
                int adjusted = clampCol(newRow, preferredCol);
                if (adjusted >= 0)
                {
                    focusRow = newRow;
                    focusCol = adjusted;
                    return;
                }
            }

            if (!option.wrapRows && (newRow == 0 || newRow == static_cast<int>(option.rows.size()) - 1))
            {
                break;
            }
        }
    };

    auto moveCol = [&](int delta) {
        const auto &row = option.rows[static_cast<std::size_t>(focusRow)];
        if (row.empty())
        {
            return;
        }
        int newCol = focusCol + delta;
        if (option.wrapCols)
        {
            newCol %= static_cast<int>(row.size());
            if (newCol < 0)
            {
                newCol += static_cast<int>(row.size());
            }
        }
        else
        {
            newCol = std::max(0, std::min(newCol, static_cast<int>(row.size()) - 1));
        }
        focusCol = newCol;
    };

    while (true)
    {
        KeyEvent event = readKey();
        switch (event.type)
        {
        case KeyEvent::Type::ArrowUp:
            moveRow(-1);
            renderFrame();
            break;
        case KeyEvent::Type::ArrowDown:
            moveRow(1);
            renderFrame();
            break;
        case KeyEvent::Type::ArrowLeft:
            moveCol(-1);
            renderFrame();
            break;
        case KeyEvent::Type::ArrowRight:
            moveCol(1);
            renderFrame();
            break;
        case KeyEvent::Type::Enter:
        {
            result.row = focusRow;
            result.col = focusCol;
            int flat = 0;
            for (int r = 0; r < focusRow; ++r)
            {
                flat += static_cast<int>(option.rows[static_cast<std::size_t>(r)].size());
            }
            result.flatIndex = flat + focusCol;
            return result;
        }
        case KeyEvent::Type::EndOfInput:
            throw std::runtime_error("INPUT_EOF");
        default:
            break;
        }
    }
}

// ===== RAII: RawSessionGuard =====
Terminal::RawSessionGuard::RawSessionGuard(Terminal &owner)
    : owner(owner), needRestore(false)
{
    if (!owner.rawModeEnabled)
    {
        owner.enableRawMode();
        needRestore = true;
    }
}

Terminal::RawSessionGuard::~RawSessionGuard()
{
    if (needRestore)
    {
        owner.disableRawMode();
    }
}

void Terminal::runInteractiveTest()
{
    Terminal terminal;
    std::vector<std::string> options = {
        "示例选项一",
        "示例选项二",
        "示例选项三"
    };
    std::string typedBuffer;

    auto render = [&terminal, &options](int active, const std::string &buffer) {
        terminal.writeRaw("\033[H\033[2J");
        terminal.writeLine("Terminal 自检: 上下方向键切换，回车确认，输入字符记录到缓冲。");
        terminal.writeLine(std::string("当前缓冲: ") + (buffer.empty() ? "<空>" : buffer));
        terminal.writeLine("");
        for (std::size_t i = 0; i < options.size(); ++i)
        {
            if (static_cast<int>(i) == active)
            {
                terminal.writeLine(std::string("> ") + options[i]);
            }
            else
            {
                terminal.writeLine(std::string("  ") + options[i]);
            }
        }
    };

    try
    {
        const int selected = terminal.runMenu(options, render, typedBuffer);
        terminal.writeLine("");
        terminal.writeLine(std::string("选择完成: ") + options[static_cast<std::size_t>(selected)]);
        if (!typedBuffer.empty())
        {
            terminal.writeLine(std::string("键入字符: ") + typedBuffer);
        }
        else
        {
            terminal.writeLine("未输入额外字符。");
        }
    }
    catch (const std::exception &ex)
    {
        terminal.disableRawMode();
        std::cerr << "终端测试中断: " << ex.what() << std::endl;
    }
}
