#include "test/test_module.h"

#include <cstdio>
#include <stdexcept>
#include <string>

namespace TestModule {
void ConfigureFileIO(const std::string &input_path, const std::string &output_path)
{
    if (!input_path.empty())
    {
        if (freopen(input_path.c_str(), "r", stdin) == nullptr)
        {
            throw std::runtime_error("无法重定向标准输入到文件: " + input_path);
        }
    }

    if (!output_path.empty())
    {
        if (freopen(output_path.c_str(), "w", stdout) == nullptr)
        {
            throw std::runtime_error("无法重定向标准输出到文件: " + output_path);
        }
    }
}

void AppendLog(const std::string &label, const std::string &message)
{
    fprintf(stdout, "<%s> %s\n", label.c_str(), message.c_str());
    fflush(stdout);
}

} // namespace TestModule
