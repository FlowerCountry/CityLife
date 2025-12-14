#include "object/wallet_inspector.h"
#include "locale/locale_manager.h"

bool WalletInspector::seeded = false;

// 统一的可用纸币面额，按金额从大到小排列。
const std::array<int, 6> WalletInspector::denominations = {100, 50, 20, 10, 5, 1};
// 每个情境下提供多条台词，实际使用时随机挑一条。
const std::array<std::string, 3> WalletInspector::zeroWalletMessages = {
    "钱包空空的，一打开只剩下风吹过。",
    "口袋拍了拍，只剩点空气作伴。",
    "翻遍钱包，结果只有比心还干净的底。"};
const std::array<std::string, 3> WalletInspector::summaryLeads = {
    "随手一翻，看到这些钞票：",
    "掀开钱包盖，里面排着：",
    "哗啦啦翻看钱包，里面静静躺着："};
const std::array<std::string, 2> WalletInspector::hasMoreMessages = {
    "剩下的零票得仔细点点才放心。",
    "还有一些散票，想看的话不妨数一遍。"};
const std::array<std::string, 3> WalletInspector::optionCloseMessages = {
    "好吧，先把这些钞票塞回去。",
    "收好钱包，别让风再偷走什么。",
    "合上钱包，留着以后慢慢花。"};
const std::array<std::string, 2> WalletInspector::detailLeadMessages = {
    "仔细点点，所有面额如下：",
    "认真清点一遍，钱包里其实是这样："};
const std::array<std::string, 2> WalletInspector::detailEndMessages = {
    "点好啦，按回车把它们收回去。",
    "数得清清楚楚，回车收好钱包。"};
// 面额对应的形容词，帮助输出更生动的文案。
const std::array<std::string, 6> WalletInspector::denominationDescriptions = {
    "鲜红的",
    "紫色的",
    "翠绿的",
    "蓝灰的",
    "褐色的",
    "浅绿色的"};

WalletInspector::Summary WalletInspector::BuildSummary(const std::array<int, 6> &counts) const
{
    Summary summary;
    summary.hasMore = false;
    // 统一口径：按非零面额种类判断是否空。
    if (CountNonZero(counts) == 0)
    {
        return summary;
    }

    // 直接根据当前钱包面额列出摘要。
    constexpr int kSummaryLimit = 5;
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        const int count = counts[i];
        if (count == 0)
        {
            continue;
        }
        summary.lines.emplace_back(FormatEntryAt(i, count));
        if (static_cast<int>(summary.lines.size()) == kSummaryLimit)
        {
            break;
        }
    }
    summary.hasMore = CountNonZero(counts) > static_cast<int>(summary.lines.size());
    return summary;
}

WalletInspector::Detail WalletInspector::BuildDetail(const std::array<int, 6> &counts) const
{
    Detail detail;
    bool any = false;
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        const int count = counts[i];
        if (count == 0)
        {
            continue;
        }
        any = true;
        detail.lines.emplace_back(FormatEntryAt(i, count));
    }
    if (!any)
    {
        LocaleManager *locale = LocaleManager::GetInstance();
        detail.lines.emplace_back(locale->Get("wallet.misc", "no_bills"));
    }
    return detail;
}

std::string WalletInspector::PickZeroWalletMessage() const
{
    LocaleManager *locale = LocaleManager::GetInstance();
    auto values = locale->GetSectionValues("wallet.empty");
    if (values.empty())
    {
        return PickRandom(zeroWalletMessages);
    }
    EnsureSeed();
    return values[static_cast<std::size_t>(std::rand() % values.size())];
}

std::string WalletInspector::PickSummaryLead() const
{
    LocaleManager *locale = LocaleManager::GetInstance();
    auto values = locale->GetSectionValues("wallet.summary_lead");
    if (values.empty())
    {
        return PickRandom(summaryLeads);
    }
    EnsureSeed();
    return values[static_cast<std::size_t>(std::rand() % values.size())];
}

std::string WalletInspector::PickHasMoreMessage() const
{
    LocaleManager *locale = LocaleManager::GetInstance();
    auto values = locale->GetSectionValues("wallet.has_more");
    if (values.empty())
    {
        return PickRandom(hasMoreMessages);
    }
    EnsureSeed();
    return values[static_cast<std::size_t>(std::rand() % values.size())];
}

std::string WalletInspector::PickCloseMessage() const
{
    LocaleManager *locale = LocaleManager::GetInstance();
    auto values = locale->GetSectionValues("wallet.close");
    if (values.empty())
    {
        return PickRandom(optionCloseMessages);
    }
    EnsureSeed();
    return values[static_cast<std::size_t>(std::rand() % values.size())];
}

std::string WalletInspector::PickDetailLead() const
{
    LocaleManager *locale = LocaleManager::GetInstance();
    auto values = locale->GetSectionValues("wallet.detail_lead");
    if (values.empty())
    {
        return PickRandom(detailLeadMessages);
    }
    EnsureSeed();
    return values[static_cast<std::size_t>(std::rand() % values.size())];
}

std::string WalletInspector::PickDetailEnd() const
{
    LocaleManager *locale = LocaleManager::GetInstance();
    auto values = locale->GetSectionValues("wallet.detail_end");
    if (values.empty())
    {
        return PickRandom(detailEndMessages);
    }
    EnsureSeed();
    return values[static_cast<std::size_t>(std::rand() % values.size())];
}

int WalletInspector::CountNonZero(const Breakdown &counts)
{
    // 统计仍然存在的面额种类数量。
    int total = 0;
    for (int count : counts)
    {
        if (count > 0)
        {
            ++total;
        }
    }
    return total;
}

std::string WalletInspector::FormatEntryAt(std::size_t index, int count)
{
    LocaleManager *locale = LocaleManager::GetInstance();
    const int denomination = denominations[index];

    // 从 LocaleManager 获取面额描述
    std::string desc;
    switch (denomination)
    {
    case 100:
        desc = locale->Get("wallet.denomination", "d100");
        break;
    case 50:
        desc = locale->Get("wallet.denomination", "d50");
        break;
    case 20:
        desc = locale->Get("wallet.denomination", "d20");
        break;
    case 10:
        desc = locale->Get("wallet.denomination", "d10");
        break;
    case 5:
        desc = locale->Get("wallet.denomination", "d5");
        break;
    case 1:
        desc = locale->Get("wallet.denomination", "d1");
        break;
    default:
        desc = denominationDescriptions[index];
    }

    return std::to_string(count) + " 张" + desc + std::to_string(denomination) + " 元";
}

void WalletInspector::EnsureSeed()
{
    if (!seeded)
    {
        std::srand(static_cast<unsigned int>(std::time(nullptr)));
        seeded = true;
    }
}
