// tests/unit/test_wallet_inspector.cpp
// WalletInspector 类的单元测试

#include <gtest/gtest.h>
#include "locale/locale_manager.h"
#include "object/wallet_inspector.h"

class WalletInspectorTest : public ::testing::Test
{
  protected:
    WalletInspector inspector;

    static void SetUpTestSuite()
    {
        // 初始化 LocaleManager，使用默认语言文件
        LocaleManager::GetInstance()->Initialize("zh_CN", ".");
    }
};

// 测试空钱包的摘要
TEST_F(WalletInspectorTest, EmptyWalletSummary)
{
    std::array<int, 6> emptyWallet = {0, 0, 0, 0, 0, 0};
    auto summary = inspector.BuildSummary(emptyWallet);

    EXPECT_TRUE(summary.lines.empty());
    EXPECT_FALSE(summary.hasMore);
}

// 测试空钱包的详单
TEST_F(WalletInspectorTest, EmptyWalletDetail)
{
    std::array<int, 6> emptyWallet = {0, 0, 0, 0, 0, 0};
    auto detail = inspector.BuildDetail(emptyWallet);

    ASSERT_EQ(detail.lines.size(), 1u);
    EXPECT_NE(detail.lines[0].find("一张纸币都没有"), std::string::npos);
}

// 测试单一面额摘要
TEST_F(WalletInspectorTest, SingleDenominationSummary)
{
    std::array<int, 6> wallet = {1, 0, 0, 0, 0, 0}; // 1张100元
    auto summary = inspector.BuildSummary(wallet);

    ASSERT_EQ(summary.lines.size(), 1u);
    EXPECT_FALSE(summary.hasMore);
    // 验证包含面额信息
    EXPECT_NE(summary.lines[0].find("100"), std::string::npos);
    EXPECT_NE(summary.lines[0].find("1 张"), std::string::npos);
}

// 测试多种面额摘要（不超过5种）
TEST_F(WalletInspectorTest, MultipleDenominationsSummary)
{
    std::array<int, 6> wallet = {1, 1, 1, 1, 0, 0}; // 4种面额
    auto summary = inspector.BuildSummary(wallet);

    EXPECT_EQ(summary.lines.size(), 4u);
    EXPECT_FALSE(summary.hasMore);
}

// 测试超过5种面额时有更多标记
TEST_F(WalletInspectorTest, HasMoreWhenExceedsLimit)
{
    std::array<int, 6> wallet = {1, 1, 1, 1, 1, 1}; // 6种面额
    auto summary = inspector.BuildSummary(wallet);

    EXPECT_EQ(summary.lines.size(), 5u);
    EXPECT_TRUE(summary.hasMore);
}

// 测试详单包含所有面额
TEST_F(WalletInspectorTest, DetailContainsAllDenominations)
{
    std::array<int, 6> wallet = {1, 2, 3, 4, 5, 6};
    auto detail = inspector.BuildDetail(wallet);

    EXPECT_EQ(detail.lines.size(), 6u);
    // 验证包含各种面额
    bool has100 = false, has50 = false, has20 = false;
    bool has10 = false, has5 = false, has1 = false;
    for (const auto &line : detail.lines)
    {
        if (line.find("100") != std::string::npos) has100 = true;
        if (line.find("50") != std::string::npos) has50 = true;
        if (line.find("20") != std::string::npos) has20 = true;
        if (line.find("10") != std::string::npos) has10 = true;
        if (line.find("5 元") != std::string::npos) has5 = true;
        if (line.find("1 元") != std::string::npos) has1 = true;
    }
    EXPECT_TRUE(has100);
    EXPECT_TRUE(has50);
    EXPECT_TRUE(has20);
    EXPECT_TRUE(has10);
    EXPECT_TRUE(has5);
    EXPECT_TRUE(has1);
}

// 测试随机消息不为空
TEST_F(WalletInspectorTest, RandomMessagesNotEmpty)
{
    EXPECT_FALSE(inspector.PickZeroWalletMessage().empty());
    EXPECT_FALSE(inspector.PickSummaryLead().empty());
    EXPECT_FALSE(inspector.PickHasMoreMessage().empty());
    EXPECT_FALSE(inspector.PickCloseMessage().empty());
    EXPECT_FALSE(inspector.PickDetailLead().empty());
    EXPECT_FALSE(inspector.PickDetailEnd().empty());
}

// 测试大数量纸币格式化
TEST_F(WalletInspectorTest, LargeCountFormatting)
{
    std::array<int, 6> wallet = {999, 0, 0, 0, 0, 0};
    auto summary = inspector.BuildSummary(wallet);

    ASSERT_EQ(summary.lines.size(), 1u);
    EXPECT_NE(summary.lines[0].find("999"), std::string::npos);
}

// 测试只有小面额的钱包
TEST_F(WalletInspectorTest, OnlySmallDenominations)
{
    std::array<int, 6> wallet = {0, 0, 0, 0, 3, 7}; // 只有5元和1元
    auto summary = inspector.BuildSummary(wallet);

    EXPECT_EQ(summary.lines.size(), 2u);
    EXPECT_FALSE(summary.hasMore);
}

// 测试跳过零数量的面额
TEST_F(WalletInspectorTest, SkipsZeroCountDenominations)
{
    std::array<int, 6> wallet = {1, 0, 1, 0, 1, 0}; // 隔一个有一种
    auto summary = inspector.BuildSummary(wallet);

    EXPECT_EQ(summary.lines.size(), 3u);
    // 确保不包含0的描述
    for (const auto &line : summary.lines)
    {
        EXPECT_EQ(line.find("0 张"), std::string::npos);
    }
}
