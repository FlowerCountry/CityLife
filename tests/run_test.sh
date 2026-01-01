#!/bin/bash
# CityLife API 自动化测试脚本

set -e

# 项目根目录
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_DIR"

# 配置
PORT=18081
BASE_URL="http://localhost:$PORT/api/v2"
PID_FILE="/tmp/citylife_test.pid"

# 测试结果统计
PASSED=0
FAILED=0

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# 清理函数
cleanup() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        kill "$PID" 2>/dev/null || true
        rm -f "$PID_FILE"
    fi
}

# 捕获退出信号
trap cleanup EXIT

# 编译项目
echo "编译项目..."
go build -o bin/citylife ./cmd/citylife

# 清理旧进程
cleanup

# 启动服务器
echo "启动API服务器 (端口 $PORT)..."
./bin/citylife -port "$PORT" &
echo $! > "$PID_FILE"
sleep 2

# 检查服务器是否启动
if ! curl -s "$BASE_URL/sessions" -X POST > /dev/null 2>&1; then
    echo "服务器启动失败"
    exit 1
fi

# 辅助函数：发送请求并返回响应
api_call() {
    local METHOD="$1"
    local API_PATH="$2"
    local DATA="$3"

    if [ -n "$DATA" ]; then
        curl -s -X "$METHOD" "$BASE_URL$API_PATH" -H "Content-Type: application/json" -d "$DATA"
    else
        curl -s -X "$METHOD" "$BASE_URL$API_PATH"
    fi
}

# 辅助函数：检查JSON字段
check_json() {
    local JSON="$1"
    local FIELD="$2"
    local EXPECTED="$3"

    local VALUE=$(echo "$JSON" | jq -r "$FIELD" 2>/dev/null)
    if [ "$VALUE" = "$EXPECTED" ]; then
        return 0
    else
        return 1
    fi
}

# 辅助函数：检查JSON字段包含字符串
check_json_contains() {
    local JSON="$1"
    local FIELD="$2"
    local PATTERN="$3"

    local VALUE=$(echo "$JSON" | jq -r "$FIELD" 2>/dev/null)
    if echo "$VALUE" | grep -q "$PATTERN"; then
        return 0
    else
        return 1
    fi
}

# 创建Session并返回ID
create_session() {
    api_call POST "/sessions" | jq -r '.data.session_id'
}

# 测试框架
run_test() {
    local TEST_NAME="$1"
    local TEST_FUNC="$2"

    echo -n "运行测试: $TEST_NAME ... "

    if $TEST_FUNC; then
        echo -e "${GREEN}[PASS]${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}[FAIL]${NC}"
        FAILED=$((FAILED + 1))
    fi
}

echo ""
echo "========================================"
echo "     CityLife API 自动化测试"
echo "========================================"
echo ""

# ========== 测试用例 ==========

# 测试：创建Session
test_create_session() {
    local RESP=$(api_call POST "/sessions")
    check_json "$RESP" ".success" "true" || return 1
    check_json "$RESP" ".data.session_id" "null" && return 1
    return 0
}

# 测试：获取初始状态
test_initial_state() {
    local SESSION=$(create_session)
    local RESP=$(api_call GET "/sessions/$SESSION/state")

    check_json "$RESP" ".success" "true" || return 1
    check_json "$RESP" ".data.location.id" "0" || return 1
    check_json "$RESP" ".data.location.name" "市中心" || return 1
    check_json "$RESP" ".data.money.wallet_total" "100" || return 1
    check_json "$RESP" ".data.is_alive" "true" || return 1
    return 0
}

# 测试：导航到超市
test_navigation_supermarket() {
    local SESSION=$(create_session)
    local RESP=$(api_call POST "/sessions/$SESSION/navigate/1")

    check_json "$RESP" ".success" "true" || return 1
    check_json_contains "$RESP" ".data.message" "超市" || return 1

    # 验证位置已变更
    RESP=$(api_call GET "/sessions/$SESSION/state")
    check_json "$RESP" ".data.location.id" "1" || return 1
    return 0
}

# 测试：导航到银行
test_navigation_bank() {
    local SESSION=$(create_session)
    local RESP=$(api_call POST "/sessions/$SESSION/navigate/2")

    check_json "$RESP" ".success" "true" || return 1
    check_json_contains "$RESP" ".data.message" "银行" || return 1
    return 0
}

# 测试：导航到医院
test_navigation_hospital() {
    local SESSION=$(create_session)
    local RESP=$(api_call POST "/sessions/$SESSION/navigate/4")

    check_json "$RESP" ".success" "true" || return 1
    check_json_contains "$RESP" ".data.message" "医院" || return 1
    return 0
}

# 测试：获取商品列表
test_get_commodities() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/1" > /dev/null

    local RESP=$(api_call GET "/sessions/$SESSION/shop/commodities")
    check_json "$RESP" ".success" "true" || return 1

    # 检查包含面包
    if ! echo "$RESP" | jq -e '.data[] | select(.name == "面包")' > /dev/null 2>&1; then
        return 1
    fi
    return 0
}

# 测试：购买商品
test_buy_commodity() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/1" > /dev/null

    local RESP=$(api_call POST "/sessions/$SESSION/shop/buy/面包")
    check_json "$RESP" ".success" "true" || return 1

    # 验证余额减少
    RESP=$(api_call GET "/sessions/$SESSION/wallet")
    local TOTAL=$(echo "$RESP" | jq -r '.data.total')
    if [ "$TOTAL" -ge 100 ]; then
        return 1
    fi
    return 0
}

# 测试：银行存款
test_bank_deposit() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/2" > /dev/null

    local RESP=$(api_call POST "/sessions/$SESSION/bank/deposit" '{"amount":100}')
    check_json "$RESP" ".success" "true" || return 1

    # 验证银行余额增加
    RESP=$(api_call GET "/sessions/$SESSION/bank")
    check_json "$RESP" ".data.balance" "100" || return 1
    return 0
}

# 测试：银行取款
test_bank_withdraw() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/2" > /dev/null
    api_call POST "/sessions/$SESSION/bank/deposit" '{"amount":100}' > /dev/null

    local RESP=$(api_call POST "/sessions/$SESSION/bank/withdraw" '{"amount":100}')
    check_json "$RESP" ".success" "true" || return 1

    # 验证银行余额为0
    RESP=$(api_call GET "/sessions/$SESSION/bank")
    check_json "$RESP" ".data.balance" "0" || return 1
    return 0
}

# 测试：看病
test_see_doctor() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/4" > /dev/null

    local RESP=$(api_call POST "/sessions/$SESSION/hospital/doctor")
    check_json "$RESP" ".success" "true" || return 1
    return 0
}

# 测试：获取体检项目
test_get_checkups() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/4" > /dev/null

    local RESP=$(api_call GET "/sessions/$SESSION/hospital/checkups")
    check_json "$RESP" ".success" "true" || return 1

    # 检查包含套餐和单项
    if ! echo "$RESP" | jq -e '.data.packages' > /dev/null 2>&1; then
        return 1
    fi
    if ! echo "$RESP" | jq -e '.data.singles' > /dev/null 2>&1; then
        return 1
    fi
    return 0
}

# 测试：执行体检
test_do_checkup() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/4" > /dev/null

    local RESP=$(api_call POST "/sessions/$SESSION/hospital/checkup/p_core")
    check_json "$RESP" ".success" "true" || return 1
    return 0
}

# 测试：获取药品列表
test_get_medicines() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/4" > /dev/null

    local RESP=$(api_call GET "/sessions/$SESSION/hospital/medicines")
    check_json "$RESP" ".success" "true" || return 1

    # 检查包含OTC药品
    if ! echo "$RESP" | jq -e '.data.otc | length > 0' > /dev/null 2>&1; then
        return 1
    fi
    return 0
}

# 测试：购买药品
test_buy_medicine() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/4" > /dev/null

    local RESP=$(api_call POST "/sessions/$SESSION/hospital/medicine/vitamin_c")
    check_json "$RESP" ".success" "true" || return 1
    return 0
}

# 测试：存档
test_save_game() {
    local SESSION=$(create_session)
    api_call POST "/sessions/$SESSION/navigate/1" > /dev/null

    local RESP=$(api_call POST "/sessions/$SESSION/saves/1")
    check_json "$RESP" ".success" "true" || return 1
    return 0
}

# 测试：读档
test_load_game() {
    local SESSION=$(create_session)
    # 先导航到超市并保存
    api_call POST "/sessions/$SESSION/navigate/1" > /dev/null
    api_call POST "/sessions/$SESSION/saves/1" > /dev/null

    # 导航到银行
    api_call POST "/sessions/$SESSION/navigate/2" > /dev/null

    # 加载存档
    local RESP=$(api_call POST "/sessions/$SESSION/saves/1/load")
    check_json "$RESP" ".success" "true" || return 1

    # 验证位置恢复到超市
    RESP=$(api_call GET "/sessions/$SESSION/state")
    check_json "$RESP" ".data.location.id" "1" || return 1
    return 0
}

# 测试：获取健康状态
test_get_health() {
    local SESSION=$(create_session)

    local RESP=$(api_call GET "/sessions/$SESSION/health")
    check_json "$RESP" ".success" "true" || return 1
    check_json "$RESP" ".data.is_alive" "true" || return 1

    # 检查核心营养存在
    if ! echo "$RESP" | jq -e '.data.core' > /dev/null 2>&1; then
        return 1
    fi
    return 0
}

# 测试：获取可用行动
test_get_actions() {
    local SESSION=$(create_session)

    local RESP=$(api_call GET "/sessions/$SESSION/actions")
    check_json "$RESP" ".success" "true" || return 1

    # 检查包含导航行动
    if ! echo "$RESP" | jq -e '.data[] | select(.category == "navigation")' > /dev/null 2>&1; then
        return 1
    fi
    return 0
}

# 测试：获取地图
test_get_map() {
    local SESSION=$(create_session)

    local RESP=$(api_call GET "/sessions/$SESSION/map")
    check_json "$RESP" ".success" "true" || return 1
    check_json "$RESP" ".data.current_location" "0" || return 1

    # 检查位置列表
    if ! echo "$RESP" | jq -e '.data.locations | length > 0' > /dev/null 2>&1; then
        return 1
    fi
    return 0
}

# 运行所有测试
run_test "创建Session" test_create_session
run_test "初始状态" test_initial_state
run_test "导航到超市" test_navigation_supermarket
run_test "导航到银行" test_navigation_bank
run_test "导航到医院" test_navigation_hospital
run_test "获取商品列表" test_get_commodities
run_test "购买商品" test_buy_commodity
run_test "银行存款" test_bank_deposit
run_test "银行取款" test_bank_withdraw
run_test "看病" test_see_doctor
run_test "获取体检项目" test_get_checkups
run_test "执行体检" test_do_checkup
run_test "获取药品列表" test_get_medicines
run_test "购买药品" test_buy_medicine
run_test "存档" test_save_game
run_test "读档" test_load_game
run_test "获取健康状态" test_get_health
run_test "获取可用行动" test_get_actions
run_test "获取地图" test_get_map

echo ""
echo "========================================"
echo "     测试结果"
echo "========================================"
echo -e "  ${GREEN}通过: $PASSED${NC}"
echo -e "  ${RED}失败: $FAILED${NC}"
echo "========================================"
echo ""

# 返回适当的退出码
if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
