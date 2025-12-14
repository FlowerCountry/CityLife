#pragma once

#include <functional>
#include <stdexcept>
#include <string>
#include <vector>

#ifdef _WIN32
#ifndef NOMINMAX
#define NOMINMAX
#endif
#include <windows.h>
#ifndef ENABLE_QUICK_EDIT_MODE
#define ENABLE_QUICK_EDIT_MODE 0x0040
#endif
#ifndef ENABLE_EXTENDED_FLAGS
#define ENABLE_EXTENDED_FLAGS 0x0080
#endif
#else
#include <termios.h>
#endif

// Terminal 提供统一的终端输入管理，支持方向键、回车以及常规字符输入，兼容 Linux/macOS。
class Terminal
{
  public:
    struct WindowSize
    {
        int columns;
        int rows;
    };

    struct GridLayout
    {
        int columns;
        int rows;
        int columnWidth;
        int itemCount;
    };

    struct MenuGridOption
    {
        std::vector<std::vector<std::string>> rows;
        bool wrapRows = true;
        bool wrapCols = true;
    };

    struct MenuGridResult
    {
        int row = -1;
        int col = -1;
        int flatIndex = -1;
    };

    // 导航轴定义：控制 runMenuEx 如何响应方向键。
    enum class NavAxis
    {
        Vertical,   // 仅响应上下键
        Horizontal, // 仅响应左右键
        Both        // 上下左右均可改变选中项
    };
    // KeyEvent 表示一次键盘输入事件。
    struct KeyEvent
    {
        // Type 枚举描述输入类别。
        enum class Type
        {
            ArrowUp,
            ArrowDown,
            ArrowLeft,
            ArrowRight,
            Enter,
            Character,
            EndOfInput,
            Unknown
        };

        // 构造函数，type 为事件类型，ch 用于 Character 类型保存原始字符。
        KeyEvent(Type type, char ch = '\0');

        Type type;
        char ch;
    };

    Terminal();
    ~Terminal();

    Terminal(const Terminal &) = delete;
    Terminal &operator=(const Terminal &) = delete;

    // enableRawMode 开启原始输入模式，关闭回显与行缓冲。
    // Throws: std::runtime_error 当系统调用失败。
    void enableRawMode();

    // disableRawMode 恢复终端为普通模式，可重复调用且幂等。
    void disableRawMode();

    // readKey 读取并解析下一次键盘输入。
    // Returns: KeyEvent，包含方向键、回车、可打印字符或 EOF 信息。
    KeyEvent readKey();

    WindowSize getWindowSize() const;

    GridLayout computeGridLayout(const std::vector<std::string> &items, int padding = 4, int maxColumns = 0) const;

    MenuGridResult runMenuGrid(const MenuGridOption &option,
                               const std::function<void(int, int, const std::vector<std::vector<std::string>> &, const std::vector<int> &)> &render = nullptr);

    // clearOutputLog 清空当前的输出日志。
    // 调用后：内部缓存的输出段落被移除，但不会清空屏幕。
    void clearOutputLog();

    // writeRaw 追加一段原始输出内容并立即写入终端。
    // text: 需要写入的原始文本，可包含 ANSI 控制序列与换行符。
    // Effects: text 会被追加到输出缓存，同时同步写到标准输出并刷新。
    void writeRaw(const std::string &text);

    // writeLine 追加一行文本并自动补充换行符。
    // line: 单行业务文本，不需要包含换行符。
    void writeLine(const std::string &line);

    // getOutputLog 访问当前累积的输出内容。
    // Returns: const 引用，记录 writeRaw/writeLine 调用顺序的文本段。
    const std::vector<std::string> &getOutputLog() const;

    // runMenu 启动菜单循环，根据上下键更新选中项，并在回车时退出。
    // options: 菜单项集合，不能为空。
    // render: 渲染回调，参数为当前高亮索引与累积的字符缓冲。
    // typedBuffer: 收集用户输入的普通字符，外部可用于过滤或搜索。
    // Returns: 用户最终确定的选项索引（0-based）。
    // Throws: std::runtime_error 当输入流结束或 options 为空。
    int runMenu(const std::vector<std::string> &options,
                const std::function<void(int, const std::string &)> &render,
                std::string &typedBuffer);

    // runMenuEx 是 runMenu 的可配置版本，可按需启用左右选择或双轴选择。
    // options: 菜单项集合（用于计数/返回索引）。
    // render: 渲染回调，参数为当前高亮索引与累积的字符缓冲。
    // typedBuffer: 收集普通字符输入（可用于模糊搜索等场景）。
    // axis: 指定导航轴（垂直/水平/双轴）。
    // wrapAround: 是否允许索引越界后环回。
    // Returns: 用户按回车时的选项索引。
    // Throws: std::runtime_error 当输入流结束或 options 为空。
    int runMenuEx(const std::vector<std::string> &options,
                  const std::function<void(int, const std::string &)> &render,
                  std::string &typedBuffer,
                  NavAxis axis,
                  bool wrapAround = true,
                  int startIndex = 0);

    // runInteractiveTest 提供一个内置的手动测试流程，验证方向键与字符输入。
    // 该流程会渲染示例菜单并打印最终选择及累积输入，不接受参数也不返回值。
    static void runInteractiveTest();

  private:
    bool rawModeEnabled;
#ifdef _WIN32
    HANDLE consoleInputHandle;
    DWORD originalConsoleMode;
#else
    termios originalTermios;
#endif

#ifndef _WIN32
    // readSingleUnsafe 从标准输入读取单个字节（未解析）。
    // Returns: -1 当遇到 EOF，否则返回 0-255。
    int readSingleUnsafe();

    // decodeEscape 解析 ESC 开头的序列，识别方向键等特殊按键。
    // Returns: 匹配到的 KeyEvent，否则返回 Unknown。
    KeyEvent decodeEscape();
#endif

    // normalizePrintable 将普通字符转换为 KeyEvent::Character 或特殊类型。
    // ch: 原始字符码。
    // Returns: 规范化后的 KeyEvent。
    KeyEvent normalizePrintable(int ch);

    // RawSessionGuard 以 RAII 方式管理原始输入模式的生存期。
    // 构造时：若当前未启用 raw 模式，则启用之并记录需要恢复；
    // 析构时：仅在构造阶段启用过的情况下才恢复为普通模式，避免干扰外部会话。
    class RawSessionGuard
    {
      public:
        // 构造函数
        // owner: 关联的 Terminal 实例。
        explicit RawSessionGuard(Terminal &owner);
        // 析构函数：按需恢复终端模式。
        ~RawSessionGuard();

      private:
        Terminal &owner;
        bool needRestore;
    };

    std::vector<std::string> outputLog;
};
