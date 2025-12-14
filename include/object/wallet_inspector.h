#pragma once

#include <array>
#include <cstdlib>
#include <ctime>
#include <string>
#include <vector>

// WalletInspector 负责处理“查看身上现金”事件中的面额统计与文案随机。
// 把金额拆分、生成摘要/详单文本以及不同场景的提示语集中管理，避免事件主体膨胀。
class WalletInspector {
  public:
    // 概览数据：最多五条面额摘要以及是否还有未展示的面额。
    struct Summary {
        std::vector<std::string> lines;
        bool hasMore;
    };

    // 详细数据：列出所有非零面额的统计。
    struct Detail {
        std::vector<std::string> lines;
    };

    // 根据当前钱包面额生成摘要（最多展示五种面额）。
    // counts: 各面额纸币的数量，顺序与 denominations 保持一致。
    Summary BuildSummary(const std::array<int, 6> &counts) const;
    // 生成完整面额明细。
    // counts: 各面额纸币的数量，顺序与 denominations 保持一致。
    Detail BuildDetail(const std::array<int, 6> &counts) const;

    // 各类随机台词选择器。
    // 返回 “钱包空空” 场景的提示语。
    std::string PickZeroWalletMessage() const;
    // 返回摘要部分的开场句。
    std::string PickSummaryLead() const;
    // 返回摘要列满 5 条后提示“还有更多”的句子。
    std::string PickHasMoreMessage() const;
    // 返回玩家选择“收好钱包”时的收尾句。
    std::string PickCloseMessage() const;
    // 返回进入“仔细点点”时的开场句。
    std::string PickDetailLead() const;
    // 返回“仔细点点”结束时的收尾句。
    std::string PickDetailEnd() const;

  private:
    using Breakdown = std::array<int, 6>;

    // 统计非零面额的数量，用于判断是否还有更多面额可展示。
    static int CountNonZero(const Breakdown &counts);
    // 格式化单项描述（按索引）："N 张{描述}{面额} 元"。
    // index: 面额在 denominations 数组中的索引（0..5）。
    // count: 该面额纸币的数量，必须 >= 0。
    // Returns: 已格式化的文本，例如 "2 张鲜艳的100 元"。
    static std::string FormatEntryAt(std::size_t index, int count);

    template <typename Container>
    static const std::string &PickRandom(const Container &container);

    // 保证伪随机种子只播一次。
    static void EnsureSeed();
    static bool seeded;

    // 固定面额与对应文案库。
    // 面额数组，按从大到小排列。
    static const std::array<int, 6> denominations;
    // “钱包空空” 场景台词候选。
    static const std::array<std::string, 3> zeroWalletMessages;
    // 摘要开场台词候选。
    static const std::array<std::string, 3> summaryLeads;
    // 通知还有更多面额时的台词候选。
    static const std::array<std::string, 2> hasMoreMessages;
    // 玩家选择“收好钱包”时的台词候选。
    static const std::array<std::string, 3> optionCloseMessages;
    // “仔细点点”模式的开场台词候选。
    static const std::array<std::string, 2> detailLeadMessages;
    // “仔细点点”模式的结束台词候选。
    static const std::array<std::string, 2> detailEndMessages;
    // 每种面额的修饰词（与 denominations 顺序对应）。
    static const std::array<std::string, 6> denominationDescriptions;
};

template <typename Container>
const std::string &WalletInspector::PickRandom(const Container &container)
{
    EnsureSeed();
    const std::size_t index = static_cast<std::size_t>(std::rand() % container.size());
    return container[index];
}
