#pragma once

#include <string>

namespace TestModule {

// 配置文件输入输出，便于批量测试或记录日志。
void ConfigureFileIO(const std::string &input_path, const std::string &output_path);

// 追加一行测试日志；在测试模式下写入输出文件。
void AppendLog(const std::string &label, const std::string &message);

} // namespace TestModule
