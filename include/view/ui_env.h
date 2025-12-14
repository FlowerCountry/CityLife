// 轻量级交互环境探测工具：集中管理是否处于交互式 TTY
#pragma once

#ifdef _WIN32
#include <cstdio>
#else
#include <unistd.h>
#endif

namespace UiEnv {

// 返回当前标准输入是否为交互式终端
inline bool IsInteractive()
{
#ifdef _WIN32
    // 简化处理：Windows 测试环境下统一视作交互式
    // 如需精确判断，可改用 _isatty(_fileno(stdin))
    return true;
#else
    return ::isatty(STDIN_FILENO) != 0;
#endif
}

} // namespace UiEnv
