# Terminal 交互模块接入指南

该文档说明如何在 CityLife 项目中使用 `Terminal` 类实现统一的终端输入处理。`Terminal` 支持 Linux/macOS（POSIX）和 Windows 平台，可无缝读取方向键、回车键以及普通字符。

## 1. 基本概念

- **原始模式**：在 POSIX 平台关闭行缓冲和回显，在 Windows 平台通过 WinAPI 禁用行输入/回显，确保可以逐键读取。
- **KeyEvent**：`Terminal::KeyEvent` 封装了一次按键，包含类型与可选的字符值。
- **菜单循环**：`runMenu` 根据上下键切换当前索引，回车返回结果，并将普通字符累计进 `typedBuffer`。

## 2. 引用与初始化

```cpp
#include "terminal/terminal.h"

Terminal terminal;
```

构造函数不会立即修改终端状态；只有调用 `enableRawMode()`、`readKey()` 或 `runMenu()` 时才会进入原始模式。

### 手动控制原始模式（可选）

```cpp
terminal.enableRawMode();
// 执行自定义读取逻辑
terminal.disableRawMode();
```

`enableRawMode()`/`disableRawMode()` 成对出现，且是幂等的；`runMenu()` 会在函数内部自动管理原始模式，外层通常无需手动控制。

## 3. 菜单示例

```cpp
std::vector<std::string> options = {"储物箱", "钱包", "返回"};
std::string typedBuffer;

auto render = [&options](int activeIndex, const std::string &buffer) {
    // 根据 activeIndex 重绘界面；buffer 为累计输入的普通字符
};

const int selected = terminal.runMenu(options, render, typedBuffer);
```

`render` 在任意索引或缓冲变化时触发；实现通常清屏并重绘选项。返回值为最终选中的索引（0 开始）。若读到 EOF（如脚本输入结束），函数会抛出 `std::runtime_error("INPUT_EOF")`。

#### 横向/双轴导航（runMenuEx）

当需要用左右键切换（或上下左右都能切换）时，使用 `runMenuEx`：

```cpp
std::vector<std::string> options = {"拿出100", "放回去100", "确认"};
std::string typed;

auto render = [&](int active, const std::string&) {
    // 自行渲染：将前两个选项画在同一行，实现“左右分布”
};

// 仅水平导航：左右键切换，回车返回索引
int i = terminal.runMenuEx(options, render, typed, Terminal::NavAxis::Horizontal);

// 同时支持上下与左右：
// int i = terminal.runMenuEx(options, render, typed, Terminal::NavAxis::Both);
```

`runMenu` 等价于 `runMenuEx(..., NavAxis::Vertical)`，保持向后兼容。

### 处理字符缓冲

`typedBuffer` 收集用户在菜单中输入的普通字符。常见用途：

- 过滤菜单选项（模糊查找）
- 提示用户已输入的搜索关键字
- 组合热键，例如 `q` 退出

如需响应退格，需要在 `runMenu` 外自行解析（默认忽略退格），可通过修改 `typedBuffer` 并强制重绘实现。

## 4. 直接读取 KeyEvent

如果不需要完整菜单，可循环调用 `readKey()`：

```cpp
while (true) {
    Terminal::KeyEvent ev = terminal.readKey();
    if (ev.type == Terminal::KeyEvent::Type::EndOfInput) {
        break;
    }
    // 根据 ev.type 或 ev.ch 执行行为
}
```

`readKey()` 在首次调用时自动进入原始模式，析构或显式调用 `disableRawMode()` 后会恢复终端。

## 5. 跨平台注意事项

| 平台        | 处理方式                                                          |
|-------------|-------------------------------------------------------------------|
| Windows     | 使用 `ReadConsoleInput` 监听 `KEY_EVENT`，屏蔽 `Quick Edit` 并恢复原模式 |
| Linux/macOS | 通过 `termios` 关闭 `ICANON` 与 `ECHO`，使用 `read` 读取单字节            |

- Windows 下的 Unicode 输入仅在 ASCII 范围内回传给调用方；若需全字符集支持，需在 `Terminal` 中扩展映射逻辑。
- `CTRL+C` 被视作 `EndOfInput`，调用方可按需捕获或重启菜单。

## 6. 内置测试入口

编译后执行：

```sh
./bin/main --terminal-test
```

测试脚本会渲染一个示例菜单，并在结束时打印最终选项与累计字符。用于快速验证终端配置、方向键和字母输入是否正常。

## 7. 常见问题

1. **界面闪烁严重**：渲染回调中建议使用最小必要的输出，可配合 ncurses 或手动控制清屏。
2. **Windows 控制台无响应**：确认程序以控制台模式运行，且未被远程桌面等环境劫持输入焦点。
3. **脚本输入导致异常**：以文件/管道驱动时可能读到 EOF，捕获 `std::runtime_error("INPUT_EOF")` 做降级处理。

如需扩展新键位或组合键，直接在 `readKey()` 的平台分支中添加解析即可。
