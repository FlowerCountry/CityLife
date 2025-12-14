// tests/unit/test_cashier.cpp
// Cashier 类的单元测试

#include <gtest/gtest.h>
#include "payment/cashier.h"

// 测试 fixture，提供公共的面额配置
class CashierTest : public ::testing::Test
{
  protected:
    // 标准面额：100, 50, 20, 10, 5, 1
    const std::array<int, 6> denominations = {100, 50, 20, 10, 5, 1};
    Cashier cashier{denominations};
};

// 测试零金额支付
TEST_F(CashierTest, ZeroAmountPayment)
{
    std::array<int, 6> wallet = {1, 0, 0, 0, 0, 0}; // 100 元
    auto result = cashier.pay(0, wallet);

    EXPECT_TRUE(result.success);
    EXPECT_EQ(result.paid, 0);
    // 不应使用任何纸币
    for (int count : result.used)
    {
        EXPECT_EQ(count, 0);
    }
}

// 测试负数金额支付
TEST_F(CashierTest, NegativeAmountPayment)
{
    std::array<int, 6> wallet = {1, 0, 0, 0, 0, 0};
    auto result = cashier.pay(-10, wallet);

    EXPECT_TRUE(result.success);
    EXPECT_EQ(result.paid, 0);
}

// 测试精确支付（无需找零）
TEST_F(CashierTest, ExactPayment)
{
    std::array<int, 6> wallet = {1, 0, 0, 0, 0, 0}; // 100 元
    auto result = cashier.pay(100, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 100);
    EXPECT_EQ(result.used[0], 1);  // 使用 1 张 100 元

    // 无找零
    for (int count : result.change)
    {
        EXPECT_EQ(count, 0);
    }
}

// 测试需要找零的支付
TEST_F(CashierTest, PaymentWithChange)
{
    std::array<int, 6> wallet = {1, 0, 0, 0, 0, 0}; // 100 元
    auto result = cashier.pay(15, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 100);
    // 找零 85 元 = 50 + 20 + 10 + 5
    EXPECT_EQ(result.change[1], 1); // 50 元
    EXPECT_EQ(result.change[2], 1); // 20 元
    EXPECT_EQ(result.change[3], 1); // 10 元
    EXPECT_EQ(result.change[4], 1); // 5 元
}

// 测试余额不足
TEST_F(CashierTest, InsufficientFunds)
{
    std::array<int, 6> wallet = {0, 0, 0, 1, 0, 0}; // 只有 10 元
    auto result = cashier.pay(15, wallet);

    EXPECT_FALSE(result.success);
}

// 测试使用多种面额
TEST_F(CashierTest, MultipleDenominations)
{
    std::array<int, 6> wallet = {0, 1, 1, 0, 0, 0}; // 50 + 20 = 70 元
    auto result = cashier.pay(70, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 70);
    EXPECT_EQ(result.used[1], 1); // 50 元
    EXPECT_EQ(result.used[2], 1); // 20 元
}

// 测试小面额组合支付
TEST_F(CashierTest, SmallDenominationsPayment)
{
    std::array<int, 6> wallet = {0, 0, 0, 2, 1, 5}; // 10*2 + 5 + 1*5 = 30 元
    auto result = cashier.pay(26, wallet);

    ASSERT_TRUE(result.success);
    // 实付应 >= 26
    EXPECT_GE(result.paid, 26);
}

// 测试空钱包
TEST_F(CashierTest, EmptyWallet)
{
    std::array<int, 6> wallet = {0, 0, 0, 0, 0, 0};
    auto result = cashier.pay(10, wallet);

    EXPECT_FALSE(result.success);
}

// 测试大额支付
TEST_F(CashierTest, LargePayment)
{
    std::array<int, 6> wallet = {5, 0, 0, 0, 0, 0}; // 500 元
    auto result = cashier.pay(350, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 400);  // 使用 4 张 100 元
    EXPECT_EQ(result.used[0], 4);
    // 找零 50 元
    EXPECT_EQ(result.change[1], 1);
}

// 测试找零计算精度
TEST_F(CashierTest, ChangeCalculationAccuracy)
{
    std::array<int, 6> wallet = {1, 0, 0, 0, 0, 0}; // 100 元
    auto result = cashier.pay(37, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 100);

    // 找零 63 元 = 50 + 10 + 1*3
    int changeTotal = 0;
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        changeTotal += result.change[i] * denominations[i];
    }
    EXPECT_EQ(changeTotal, 63);
}

// 测试刚好凑足
TEST_F(CashierTest, ExactFit)
{
    std::array<int, 6> wallet = {0, 0, 1, 1, 1, 0}; // 20 + 10 + 5 = 35 元
    auto result = cashier.pay(35, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 35);
    // 无找零
    for (int count : result.change)
    {
        EXPECT_EQ(count, 0);
    }
}

// 测试需要凑钱的场景（贪心不够时需要多出一张）
TEST_F(CashierTest, NeedOverpayForChange)
{
    std::array<int, 6> wallet = {0, 1, 0, 0, 0, 0}; // 只有 50 元
    auto result = cashier.pay(30, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 50);
    // 找零 20 元
    EXPECT_EQ(result.change[2], 1);
}

// 测试边界情况：刚好1元
TEST_F(CashierTest, OneYuanPayment)
{
    std::array<int, 6> wallet = {0, 0, 0, 0, 0, 5}; // 5张1元
    auto result = cashier.pay(1, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 1);
    EXPECT_EQ(result.used[5], 1);
}

// 测试需要多张同面额
TEST_F(CashierTest, MultipleSameDenomination)
{
    std::array<int, 6> wallet = {0, 0, 5, 0, 0, 0}; // 5张20元
    auto result = cashier.pay(60, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 60);
    EXPECT_EQ(result.used[2], 3);
}

// 测试复杂找零场景
TEST_F(CashierTest, ComplexChangeScenario)
{
    std::array<int, 6> wallet = {2, 0, 0, 0, 0, 0}; // 200元
    auto result = cashier.pay(137, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 200);
    // 找零 63 元
    int changeTotal = 0;
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        changeTotal += result.change[i] * denominations[i];
    }
    EXPECT_EQ(changeTotal, 63);
}

// 测试只有小面额时的支付
TEST_F(CashierTest, OnlySmallDenominations)
{
    std::array<int, 6> wallet = {0, 0, 0, 0, 10, 20}; // 50+20=70元
    auto result = cashier.pay(67, wallet);

    ASSERT_TRUE(result.success);
    // 验证支付成功
    int paidTotal = 0;
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        paidTotal += result.used[i] * denominations[i];
    }
    EXPECT_GE(paidTotal, 67);
}

// 测试支付金额等于钱包总额
TEST_F(CashierTest, PayExactWalletTotal)
{
    std::array<int, 6> wallet = {1, 1, 1, 1, 1, 1}; // 186元
    auto result = cashier.pay(186, wallet);

    ASSERT_TRUE(result.success);
    EXPECT_EQ(result.paid, 186);
    // 无找零
    for (int c : result.change)
    {
        EXPECT_EQ(c, 0);
    }
}

// 测试支付金额超过钱包总额1元
TEST_F(CashierTest, PayOnMoreThanWallet)
{
    std::array<int, 6> wallet = {1, 0, 0, 0, 0, 0}; // 100元
    auto result = cashier.pay(101, wallet);

    EXPECT_FALSE(result.success);
}
