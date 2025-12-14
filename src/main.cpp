#include "locale/locale_manager.h"
#include "path/path_manager.h"
#include "terminal/terminal.h"
#include "test/test_module.h"
#include "world/world.h"

#include <exception>
#include <iostream>
#include <string>

namespace {

std::string GetExecutableDirectory(const char *argv0)
{
    std::string path(argv0);
    // 查找最后一个路径分隔符
    std::size_t lastSlash = path.find_last_of("/\\");
    if (lastSlash != std::string::npos)
    {
        return path.substr(0, lastSlash);
    }
    // 没有路径分隔符，使用当前目录
    return ".";
}

void ConfigureTestIO(int argc, char **argv)
{
    if (argc < 2)
    {
        return;
    }

    const std::string flag(argv[1]);
    if (flag != "--test")
    {
        return;
    }

    const std::string input_path = argc > 2 ? std::string(argv[2]) : std::string();
    const std::string output_path = argc > 3 ? std::string(argv[3]) : std::string();

    TestModule::ConfigureFileIO(input_path, output_path);
}

bool HandleTerminalTest(int argc, char **argv)
{
    if (argc < 2)
    {
        return false;
    }

    if (std::string(argv[1]) != "--terminal-test")
    {
        return false;
    }

    Terminal::runInteractiveTest();
    return true;
}
} // namespace

int main(int argc, char **argv)
{
    try
    {
        if (HandleTerminalTest(argc, argv))
        {
            return 0;
        }

        // 初始化路径管理器
        std::string exeDir = GetExecutableDirectory(argv[0]);
        PathManager::GetInstance()->Initialize(exeDir);
        if (!PathManager::GetInstance()->EnsureDirectories())
        {
            std::cerr << "警告: 无法创建数据目录" << std::endl;
        }

        // 迁移旧数据到新位置
        PathManager::GetInstance()->MigrateOldData();

        // 初始化语言资源管理器
        if (!LocaleManager::GetInstance()->Initialize("zh_CN", exeDir))
        {
            std::cerr << "警告: 无法加载语言文件，将使用 fallback 值" << std::endl;
        }

        ConfigureTestIO(argc, argv);
        World::GetInstance()->Start();
    } catch (const std::runtime_error &ex)
    {
        if (std::string(ex.what()) == "INPUT_EOF")
        {
            return 0;
        }
        std::cerr << "运行失败: " << ex.what() << std::endl;
        return 1;
    } catch (const std::exception &ex)
    {
        std::cerr << "运行失败: " << ex.what() << std::endl;
        return 1;
    } catch (...)
    {
        std::cerr << "运行失败: 未知异常" << std::endl;
        return 1;
    }
    return 0;
}
