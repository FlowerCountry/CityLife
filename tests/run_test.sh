#!/bin/bash
# 自动化测试脚本 - 增强版
# 功能：运行多个测试场景并验证输出

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_DIR"

PASS_COUNT=0
FAIL_COUNT=0

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 辅助函数：运行单个测试
# 用法: run_test "测试名称" "输入文件" "断言1" "断言2" ...
run_test() {
    local TEST_NAME="$1"
    local INPUT_FILE="$2"
    local OUTPUT_FILE="tests/${TEST_NAME// /_}_actual.txt"
    shift 2
    local ASSERTIONS=("$@")

    echo ""
    echo -e "${YELLOW}运行测试: $TEST_NAME${NC}"

    # 检查输入文件是否存在
    if [ ! -f "$INPUT_FILE" ]; then
        echo -e "  ${RED}[SKIP]${NC} 输入文件不存在: $INPUT_FILE"
        return
    fi

    # 运行测试
    ./bin/main --test "$INPUT_FILE" "$OUTPUT_FILE" 2>&1 || true

    # 检查输出文件是否生成
    if [ ! -f "$OUTPUT_FILE" ]; then
        echo -e "  ${RED}[FAIL]${NC} 输出文件未生成"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        return
    fi

    local TEST_PASS=1
    for assertion in "${ASSERTIONS[@]}"; do
        if grep -q "$assertion" "$OUTPUT_FILE"; then
            echo -e "  ${GREEN}[OK]${NC} 找到: $assertion"
        else
            echo -e "  ${RED}[FAIL]${NC} 未找到: $assertion"
            TEST_PASS=0
        fi
    done

    if [ "$TEST_PASS" -eq 1 ]; then
        echo -e "  => ${GREEN}测试通过${NC}"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "  => ${RED}测试失败${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi

    # 清理临时文件
    rm -f "$OUTPUT_FILE"
}

echo "=== CityLife 自动化测试 ==="
echo "项目目录: $PROJECT_DIR"

# 检查可执行文件
if [ ! -f "./bin/main" ]; then
    echo -e "${RED}错误: 可执行文件 ./bin/main 不存在${NC}"
    echo "请先运行: cmake -S . -B build && cmake --build build"
    exit 1
fi

# 清理旧的存档文件（防止测试污染）
rm -rf ~/.citylife/saves/save_*.sav 2>/dev/null || true

# ==================== 测试场景 ====================

# 测试1: 基础功能测试
run_test "基础功能" "tests/automated_input.txt" \
    "购买面包成功" \
    "购买牛奶成功"

# 测试2: 边界测试
run_test "边界测试" "tests/boundary_input.txt" \
    "输入无效"

# 测试3: 存档系统测试
run_test "存档系统" "tests/save_load_input.txt" \
    "===== 保存游戏 =====" \
    "游戏已保存到槽位"

# 测试4: 商品覆盖测试
run_test "商品覆盖" "tests/commodity_input.txt" \
    "购买苹果成功" \
    "你感到身体更健康了" \
    "购买番茄成功" \
    "购买矿泉水成功"

# 测试5: 支付系统测试
run_test "支付系统" "tests/payment_input.txt" \
    "购买面包成功" \
    "购买牛奶成功" \
    "钱包里的钱不够这件商品"

# 测试6: 健康属性测试
run_test "健康属性" "tests/health_input.txt" \
    "你感到很饱足" \
    "你感到精力充沛" \
    "你感到身体更健康了"

# 测试7: 银行操作测试
run_test "银行操作" "tests/bank_operations_input.txt" \
    "请输入你要存的钱" \
    "请输入你要取的钱" \
    "你的银行余额"

# 测试8: 存档恢复测试
run_test "存档恢复" "tests/save_restore_input.txt" \
    "游戏已保存到槽位" \
    "存档已加载"

# 测试9: 更多商品测试
run_test "更多商品" "tests/commodity_extra_input.txt" \
    "购买蛋糕成功" \
    "购买鸡蛋成功" \
    "购买意大利面成功"

# 测试10: 导航系统测试
run_test "导航系统" "tests/navigation_input.txt" \
    "你当前位于: 超市" \
    "你当前位于: 银行" \
    "你当前位于: 市中心"

# 测试11: 更多商品覆盖
run_test "商品覆盖2" "tests/commodity_more_input.txt" \
    "购买酸奶成功" \
    "购买巧克力成功" \
    "购买薯片成功"

# 测试12: 银行边界测试
run_test "银行边界" "tests/bank_edge_input.txt" \
    "你没有这么多的钱"

# 测试13: 剩余商品测试
run_test "商品覆盖3" "tests/commodity_final_input.txt" \
    "购买咖啡成功" \
    "购买沙拉成功" \
    "购买米饭成功" \
    "购买燕麦片成功"

# 测试14: 公告系统测试
run_test "公告系统" "tests/announcement_input.txt" \
    "银行" \
    "电信大楼"

# 测试15: 高价商品测试
run_test "高价商品" "tests/commodity_expensive_input.txt" \
    "购买鸡胸肉成功" \
    "你感到精力充沛" \
    "购买巧克力成功" \
    "一股幸福感涌上心头"

# 测试16: 剩余商品测试
run_test "剩余商品" "tests/commodity_remaining_input.txt" \
    "购买橙汁成功" \
    "购买三明治成功" \
    "购买冰淇淋成功"

# ==================== 单元测试 ====================

echo ""
echo -e "${YELLOW}运行单元测试${NC}"

if [ -f "./bin/unit_tests" ]; then
    if ./bin/unit_tests --gtest_brief=1 2>&1; then
        echo -e "  => ${GREEN}单元测试通过${NC}"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "  => ${RED}单元测试失败${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
else
    echo -e "  ${YELLOW}[SKIP]${NC} 单元测试可执行文件不存在"
    echo "  提示: 使用 cmake -DBUILD_TESTS=ON 构建以启用单元测试"
fi

# ==================== 测试汇总 ====================

echo ""
echo "==========================================="
echo "               测试汇总"
echo "==========================================="
echo -e "通过: ${GREEN}$PASS_COUNT${NC}"
echo -e "失败: ${RED}$FAIL_COUNT${NC}"
echo ""

if [ "$FAIL_COUNT" -eq 0 ]; then
    echo -e "${GREEN}所有测试通过!${NC}"
    exit 0
else
    echo -e "${RED}存在 $FAIL_COUNT 个失败的测试${NC}"
    exit 1
fi
