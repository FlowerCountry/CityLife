# Test Module Notes

- 运行 `./tests/run_test.sh` 可一键完成基础自动化回归：脚本会调用 `./bin/main --test tests/automated_input.txt tests/actual_output.txt` 并在比较结果前过滤随机公告行，适合日常验证。
- 边界测试输入 `tests/boundary_input.txt` 聚焦极端交互：包含超范围商品编号、异常存取款金额以及余额不足购买等场景，可通过 `./bin/main --test tests/boundary_input.txt tests/boundary_output.txt` 复现并审阅 `tests/boundary_output.txt`。
- 市中心公告的排序在每次运行时都会随机化，写测试断言时需忽略公告行或进行过滤。
- 若未指定测试模式，程序将使用 ncurses 进行正常交互。
- `TestModule::AppendLog` 会直接把日志写入测试输出文件，方便核对断言。
- 初始化时若输入文件不存在会自动生成并写入 `0`，以确保 `View::Scan` 至少有一个默认值。
