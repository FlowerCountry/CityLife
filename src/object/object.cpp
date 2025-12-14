/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 09:42:34
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-26 11:02:15
 * @FilePath: \CityLife\src\object\object.cpp
 */
#include "object/object.h"
#include "center/center.h"
#include "disease/disease_manager.h"
#include "health/health.h"
#include "nutrition/nutrition_manager.h"
#include "object/wallet_inspector.h"
#include "payment/cashier.h"
#include "save/save_manager.h"
#include "terminal/terminal.h"
#include "user/user.h"
#include "view/view.h"
#include "view/ui_env.h"
#include "world/world.h"

#include <algorithm>
#include <array>
#include <sstream>


void GoWhere::ToDoIt()
{
    World *world = World::GetInstance();
    User *user = User::GetInstance();

    // 移动到目标位置
    world->ChangeWhere(this->to);

    // 根据距离消耗时间（每单位距离 = 1 分钟）
    world->UpdateTime(hungry * 60);

    // 根据距离消耗饱腹感和饥渴
    // 移动消耗：饱腹感 = 距离 * 2，饥渴 = 距离 * 1.5
    user->ConsumeNutrition("饱腹感", hungry * 2);
    user->ConsumeNutrition("饥渴", static_cast<int>(hungry * 1.5f));

    // 移动也会消耗少量蛋白质和碳水化合物
    user->ConsumeNutrition("蛋白质", hungry / 2);
    user->ConsumeNutrition("碳水化合物", hungry / 2);
}

EventCategory GoWhere::GetCategory() const
{
    return EventCategory::Navigation;
}

void Buy::ToDoIt()
{
    World::GetInstance()->ChangeWhere(3);
}

EventCategory Buy::GetCategory() const
{
    return EventCategory::Primary;
}
Information::Information(const std::string &name, const std::string &content) : Object(name), content{content}
{
    things["公告"] = std::bind(&Center::PrintAnnouncement, Center::GetInstance());
}
void Information::ToDoIt()
{
    things[content]();
}

EventCategory Information::GetCategory() const
{
    return EventCategory::Insight;
}
void DepositingMoney::ToDoIt()
{
    // 非交互：退回到旧的数字输入流程，保证自动化测试兼容
    if (!UiEnv::IsInteractive())
    {
        World *world = World::GetInstance();
        while (true)
        {
            View::GetInstance()->Clear();
            View::GetInstance()->Print(0, 0, "请输入你要存的钱:");
            int money = View::GetInstance()->Scan(0, 0);
            if (world->DepositingMoney(money))
            {
                return;
            }
        }
    }

    // 交互式 TTY：使用“拿出100/放回去100/确认存入”的左右布局
    World *world = World::GetInstance();
    Terminal terminal;
    std::string typed;

    int stagedHundreds = 0; // 暂存准备存入的 100 元张数
    std::string notice;
    int focus = 0;

    auto getHundredsInWallet = [&]() { return world->GetWalletBreakdown()[0]; };
    const std::vector<std::string> optionTokens = {"take100", "putback100", "confirm"};

    auto render = [&](int active, const std::string &) {
        const int availableHundreds = getHundredsInWallet();
        const int amount = stagedHundreds * 100;

        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "银行存款 (仅整百)\n\n";
        frame << "你当前持有 100 元纸币: " << availableHundreds << " 张\n";
        frame << "准备存入金额: " << amount << " 元\n\n";

        const std::string left = (active == 0 ? "> 拿出100" : "  拿出100");
        const std::string right = (active == 1 ? "> 放回去100" : "  放回去100");
        frame << left << "      " << right << "\n";
        frame << (active == 2 ? "> 确认存入" : "  确认存入") << "\n";

        if (!notice.empty())
        {
            frame << "\n" << notice << "\n";
        }
        terminal.writeRaw(frame.str());
    };

    while (true)
    {
        const int chosen = terminal.runMenuEx(optionTokens, render, typed, Terminal::NavAxis::Both, true, focus);
        focus = chosen;

        const int availableHundreds = getHundredsInWallet();
        if (chosen == 0)
        {
            if (stagedHundreds < availableHundreds)
            {
                ++stagedHundreds;
                notice.clear();
            }
            else
            {
                notice = "没有更多整百可供存入。";
            }
        }
        else if (chosen == 1)
        {
            if (stagedHundreds > 0)
            {
                --stagedHundreds;
                notice.clear();
            }
            else
            {
                notice = "没有可以放回去的整百。";
            }
        }
        else
        {
            const int money = stagedHundreds * 100;
            if (money == 0)
            {
                return;
            }
            if (world->DepositingMoney(money))
            {
                return;
            }
            notice = "存入失败，请根据提示调整金额。";
        }
    }
}

EventCategory DepositingMoney::GetCategory() const
{
    return EventCategory::Primary;
}
void WithdrawMoney::ToDoIt()
{
    // 银行取款流程：输入必须是 100 的整数倍，World 会处理失败提示。
    while (true)
    {
        View::GetInstance()->Clear();
        // 为保持测试输出兼容，提示精简为旧格式
        View::GetInstance()->Print(0, 0, "请输入你要取的钱:");
        int money = View::GetInstance()->Scan(0, 28);
        if (World::GetInstance()->WithdrawMoney(money))
        {
            return;
        }
    }
}

EventCategory WithdrawMoney::GetCategory() const
{
    return EventCategory::Primary;
}

void CheckCash::ToDoIt()
{
    WalletInspector inspector;
    const int walletTotal = World::GetInstance()->GetWallet();
    const auto &walletBreakdown = World::GetInstance()->GetWalletBreakdown();
    const bool interactive = UiEnv::IsInteractive();

    if (!interactive)
    {
        View *view = View::GetInstance();
        view->Clear();
        if (walletTotal == 0)
        {
            view->Print(0, 0, inspector.PickZeroWalletMessage());
            view->WaitForEnter();
            return;
        }

        const WalletInspector::Summary summary = inspector.BuildSummary(walletBreakdown);
        const std::string summaryLead = inspector.PickSummaryLead();
        const std::string moreMessage = summary.hasMore ? inspector.PickHasMoreMessage() : std::string();
        const std::array<std::string, 2> menuOptions = {"仔细点点", "收好钱包"};

        Terminal terminal;
        int choice = 0;

        auto renderLegacy = [&](int focus) {
            view->Clear();
            int line = 0;
            view->Print(line++, 0, summaryLead);
            for (const auto &entry : summary.lines)
            {
                view->Print(line++, 0, entry);
            }
            if (!moreMessage.empty())
            {
                view->Print(line++, 0, moreMessage);
            }
            view->Print(line++, 0, "");
            for (std::size_t i = 0; i < menuOptions.size(); ++i)
            {
                const bool isActive = static_cast<int>(i) == focus;
                const std::string prefix = isActive ? "> " : "  ";
                view->Print(line++, 0, prefix + menuOptions[i]);
            }
        };

        renderLegacy(choice);
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            switch (event.type)
            {
            case Terminal::KeyEvent::Type::ArrowUp:
                choice = (choice - 1 + static_cast<int>(menuOptions.size())) % static_cast<int>(menuOptions.size());
                renderLegacy(choice);
                break;
            case Terminal::KeyEvent::Type::ArrowDown:
                choice = (choice + 1) % static_cast<int>(menuOptions.size());
                renderLegacy(choice);
                break;
            case Terminal::KeyEvent::Type::Enter:
                terminal.disableRawMode();
                if (choice == 0)
                {
                    const WalletInspector::Detail detail = inspector.BuildDetail(walletBreakdown);
                    view->Clear();
                    view->Print(0, 0, inspector.PickDetailLead());
                    int detailLine = 1;
                    for (const auto &entry : detail.lines)
                    {
                        view->Print(detailLine++, 0, entry);
                    }
                    view->Print(detailLine, 0, inspector.PickDetailEnd());
                    view->WaitForEnter();
                    return;
                }

                view->Clear();
                view->Print(0, 0, inspector.PickCloseMessage());
                view->WaitForEnter();
                return;
            case Terminal::KeyEvent::Type::EndOfInput:
                terminal.disableRawMode();
                throw std::runtime_error("INPUT_EOF");
            default:
                break;
            }
        }
    }

    Terminal terminal;

    auto waitForEnter = [&terminal]() {
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
            {
                return;
            }
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
            {
                throw std::runtime_error("INPUT_EOF");
            }
        }
    };

    if (walletTotal == 0)
    {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << inspector.PickZeroWalletMessage() << "\n\n按回车返回";
        terminal.writeRaw(frame.str());
        waitForEnter();
        return;
    }

    const WalletInspector::Summary summary = inspector.BuildSummary(walletBreakdown);
    const std::string summaryLead = inspector.PickSummaryLead();
    const std::string moreMessage = summary.hasMore ? inspector.PickHasMoreMessage() : std::string();
    const std::array<std::string, 2> menuOptions = {"仔细点点", "收好钱包"};

    int choice = 0;

    auto renderFrame = [&](int focus) {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << summaryLead << '\n';
        for (const auto &entry : summary.lines)
        {
            frame << entry << '\n';
        }
        if (!moreMessage.empty())
        {
            frame << moreMessage << '\n';
        }
        frame << '\n';
        for (std::size_t i = 0; i < menuOptions.size(); ++i)
        {
            const bool isActive = static_cast<int>(i) == focus;
            frame << (isActive ? "> " : "  ") << menuOptions[i] << '\n';
        }
        terminal.writeRaw(frame.str());
    };

    renderFrame(choice);
    while (true)
    {
        Terminal::KeyEvent event = terminal.readKey();
        switch (event.type)
        {
        case Terminal::KeyEvent::Type::ArrowUp:
            choice = (choice - 1 + static_cast<int>(menuOptions.size())) % static_cast<int>(menuOptions.size());
            renderFrame(choice);
            break;
        case Terminal::KeyEvent::Type::ArrowDown:
            choice = (choice + 1) % static_cast<int>(menuOptions.size());
            renderFrame(choice);
            break;
        case Terminal::KeyEvent::Type::Enter:
            if (choice == 0)
            {
                const WalletInspector::Detail detail = inspector.BuildDetail(walletBreakdown);
                std::ostringstream detailFrame;
                detailFrame << "\033[H\033[2J";
                detailFrame << inspector.PickDetailLead() << '\n';
                for (const auto &entry : detail.lines)
                {
                    detailFrame << entry << '\n';
                }
                detailFrame << inspector.PickDetailEnd() << "\n\n按回车返回";
                terminal.writeRaw(detailFrame.str());
                waitForEnter();
                renderFrame(choice);
            }
            else
            {
                std::ostringstream closeFrame;
                closeFrame << "\033[H\033[2J";
                closeFrame << inspector.PickCloseMessage() << "\n\n按回车返回";
                terminal.writeRaw(closeFrame.str());
                waitForEnter();
                return;
            }
            break;
        case Terminal::KeyEvent::Type::EndOfInput:
            throw std::runtime_error("INPUT_EOF");
        default:
            break;
        }
    }
}


EventCategory CheckCash::GetCategory() const
{
    return EventCategory::Insight;
}

void CheckBankBalance::ToDoIt()
{
    if (!UiEnv::IsInteractive())
    {
        View::GetInstance()->Clear();
        View::GetInstance()->Print(0, 0, "你的银行余额: " + std::to_string(World::GetInstance()->GetBankBalance()));
        View::GetInstance()->WaitForEnter();
        return;
    }

    Terminal terminal;
    const auto waitForEnter = [&terminal]() {
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
            {
                return;
            }
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
            {
                throw std::runtime_error("INPUT_EOF");
            }
        }
    };

    std::ostringstream frame;
    frame << "\033[H\033[2J";
    frame << "你的银行余额: " << World::GetInstance()->GetBankBalance() << "\n\n按回车返回";
    terminal.writeRaw(frame.str());
    waitForEnter();
}

EventCategory CheckBankBalance::GetCategory() const
{
    return EventCategory::Insight;
}

Commodity::Commodity(const std::string &name, const int &price, const std::vector<Health> &health) : Object{name}, price{price}, health{health} {}

std::vector<Health> Commodity::GetHealth() { return health; }

void Commodity::ToDoIt()
{
    World *world = World::GetInstance();
    View *view = View::GetInstance();
    static const std::array<int, 6> denom = {100, 50, 20, 10, 5, 1};
    const bool interactive = UiEnv::IsInteractive();

    if (world->GetWallet() < price)
    {
        if (!interactive)
        {
            view->Clear();
            view->Print(0, 0, "钱包里的钱不够这件商品");
            view->WaitForEnter();
        }
        else
        {
            Terminal terminal;
            const auto waitForEnter = [&terminal]() {
                while (true)
                {
                    Terminal::KeyEvent event = terminal.readKey();
                    if (event.type == Terminal::KeyEvent::Type::Enter)
                    {
                        return;
                    }
                    if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                    {
                        throw std::runtime_error("INPUT_EOF");
                    }
                }
            };

            std::ostringstream frame;
            frame << "\033[H\033[2J";
            frame << "钱包里的钱不够这件商品" << "\n\n按回车返回";
            terminal.writeRaw(frame.str());
            waitForEnter();
        }
        world->ChangeWhere(1);
        return;
    }

    // 非交互（测试）环境：调用扣款流程+打印结果，绕过交互式点钞
    if (!interactive)
    {
        if (!world->SpendMoney(price))
        {
            return; // 失败提示已由 SpendMoney 输出
        }
        for (auto i : health)
        {
            (*User::GetInstance()->GetUserHealth())[i.GetInfo()] = std::min((*User::GetInstance()->GetUserHealth())[i.GetInfo()] + i.GetReserves(), 100);
        }
        view->Clear();
        int st = 0;
        view->Print(st, 0, "购买" + name + "成功");
        NutritionManager *nutrition = NutritionManager::GetInstance();
        for (auto i : health)
        {
            int currentValue = (*User::GetInstance()->GetUserHealth())[i.GetInfo()];
            view->Print(++st, 0, nutrition->GetFeelingDescription(i.GetInfo(), currentValue));
        }
        view->WaitForEnter();
        world->ChangeWhere(1);
        return;
    }

    Terminal terminal;
    const auto waitForEnter = [&terminal]() {
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
            {
                return;
            }
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
            {
                throw std::runtime_error("INPUT_EOF");
            }
        }
    };

    auto computeSum = [&](const std::array<int, 6> &counts) {
        int total = 0;
        for (std::size_t i = 0; i < denom.size(); ++i)
        {
            total += counts[i] * denom[i];
        }
        return total;
    };

    std::array<int, 6> selection{};
    std::array<int, 6> visible{};
    int lastPaid = 0;
    int lastChange = 0;

    const auto &initialWallet = world->GetWalletBreakdown();
    for (std::size_t i = 0; i < denom.size(); ++i)
    {
        const int total = initialWallet[i];
        visible[i] = std::min(total, total >= 3 ? 3 : total);
    }

    enum class EntryType
    {
        Denom,
        Commit,
        Reset,
        Flip,
        Cancel
    };

    struct MenuEntry
    {
        EntryType type;
        std::size_t index;
    };

    std::string notice;
    bool finished = false;

    while (!finished)
    {
        const auto &wallet = world->GetWalletBreakdown();
        for (std::size_t i = 0; i < denom.size(); ++i)
        {
            visible[i] = std::min(visible[i], wallet[i]);
            if (selection[i] > wallet[i])
            {
                selection[i] = wallet[i];
            }
        }

        const int current = computeSum(selection);

        std::vector<std::vector<MenuEntry>> gridEntries(denom.size());
        std::vector<std::vector<std::string>> gridLabels(denom.size());

        auto formatDenomination = [&](std::size_t bill) {
            std::string line = std::to_string(denom[bill]) + " 元";
            const int total = wallet[bill];
            const int discovered = visible[bill];
            const int picked = selection[bill];
            const int availableNow = std::max(0, discovered - picked);

            if (picked > 0)
            {
                line += " *已选";
            }
            if (total == 0)
            {
                line += " [没有]";
            }
            else if (availableNow <= 0)
            {
                if (discovered < total)
                {
                    line += " [还没翻出来]";
                }
                else
                {
                    line += " [暂时摸不到]";
                }
            }
            else
            {
                const int remainingTotal = total - picked;
                if (remainingTotal < 3)
                {
                    line += " (剩余 " + std::to_string(remainingTotal) + " 张)";
                }
                else if (availableNow < 3)
                {
                    line += " (可拿 " + std::to_string(availableNow) + " 张)";
                }
            }
            return line;
        };

        for (std::size_t i = 0; i < denom.size(); ++i)
        {
            gridEntries[i].push_back({EntryType::Denom, i});
            gridLabels[i].push_back(formatDenomination(i));
        }

        const std::array<MenuEntry, 4> actionEntries = {
            MenuEntry{EntryType::Commit, 0},
            MenuEntry{EntryType::Reset, 0},
            MenuEntry{EntryType::Flip, 0},
            MenuEntry{EntryType::Cancel, 0}
        };
        const std::array<const char *, 4> actionLabels = {
            "完成支付",
            "重新选择",
            "翻翻钱包",
            "放弃购买"
        };

        for (std::size_t j = 0; j < actionEntries.size(); ++j)
        {
            if (j >= gridEntries.size())
            {
                gridEntries.resize(j + 1);
                gridLabels.resize(j + 1);
            }
            gridEntries[j].push_back(actionEntries[j]);
            gridLabels[j].push_back(actionLabels[j]);
        }

        Terminal::MenuGridOption option;
        option.rows = gridLabels;
        option.wrapRows = true;
        option.wrapCols = true;

        auto renderGrid = [&](int activeRow, int activeCol, const std::vector<std::vector<std::string>> &rows, const std::vector<int> &columnWidths) {
            std::ostringstream frame;
            frame << "\033[H\033[2J";
            frame << "购买 " << name << " 需要 " << price << " 元\n";
            frame << "当前已选金额: " << current;
            if (current < price)
            {
                frame << " (还差 " << (price - current) << " 元)";
            }
            frame << "\n\n";

            const std::size_t columnCount = columnWidths.size();
            for (std::size_t row = 0; row < rows.size(); ++row)
            {
                for (std::size_t col = 0; col < columnCount; ++col)
                {
                    std::string cell;
                    if (col < rows[row].size())
                    {
                        const bool active = static_cast<int>(row) == activeRow && static_cast<int>(col) == activeCol;
                        cell = (active ? "> " : "  ") + rows[row][col];
                    }
                    if (static_cast<int>(cell.size()) < columnWidths[static_cast<std::size_t>(col)])
                    {
                        cell += std::string(columnWidths[static_cast<std::size_t>(col)] - static_cast<int>(cell.size()), ' ');
                    }
                    frame << cell;
                }
                frame << '\n';
            }

            if (!notice.empty())
            {
                frame << "\n" << notice << '\n';
            }
            terminal.writeRaw(frame.str());
        };

        const auto gridResult = terminal.runMenuGrid(option, renderGrid);
        const MenuEntry &entry = gridEntries[static_cast<std::size_t>(gridResult.row)][static_cast<std::size_t>(gridResult.col)];

        switch (entry.type)
        {
        case EntryType::Denom:
        {
            const std::size_t bill = entry.index;
            const int discovered = visible[bill];
            const int picked = selection[bill];
            const int availableNow = std::max(0, discovered - picked);
            if (availableNow <= 0)
            {
                notice = "这一面额暂时掏不出来。";
                break;
            }
            ++selection[bill];
            notice.clear();
            break;
        }
        case EntryType::Commit:
        {
            if (current <= 0)
            {
                notice = "还没选任何钞票。";
                break;
            }

            Cashier cashier(denom);
            const Cashier::Result payResult = cashier.pay(price, selection);
            if (!payResult.success)
            {
                notice = "这组钞票凑不出目标金额。";
                break;
            }

            if (!world->RemoveBills(payResult.used))
            {
                notice = "钱包里找不到这些面额组合。";
                break;
            }

            if (std::any_of(payResult.change.begin(), payResult.change.end(), [](int c) { return c > 0; }))
            {
                world->AddBills(payResult.change);
            }

            lastPaid = payResult.paid;
            lastChange = payResult.paid - price;
            finished = true;
            notice.clear();
            break;
        }
        case EntryType::Reset:
            selection.fill(0);
            notice = "已清空选择。";
            break;
        case EntryType::Flip:
        {
            bool revealed = false;
            for (std::size_t i = 0; i < denom.size(); ++i)
            {
                if (visible[i] < wallet[i])
                {
                    ++visible[i];
                    revealed = true;
                }
            }
            notice = revealed ? "翻翻钱包，又摸出几张皱巴巴的钞票。" : "翻遍钱包也没多的了。";
            break;
        }
        case EntryType::Cancel:
        {
            std::ostringstream frame;
            frame << "\033[H\033[2J";
            frame << "你放弃了购买" << name << "。" << "\n\n按回车返回";
            terminal.writeRaw(frame.str());
            waitForEnter();
            world->ChangeWhere(1);
            return;
        }
        }
    }

    for (auto i : health)
    {
        (*User::GetInstance()->GetUserHealth())[i.GetInfo()] = std::min((*User::GetInstance()->GetUserHealth())[i.GetInfo()] + i.GetReserves(), 100);
    }
    std::ostringstream summaryFrame;
    summaryFrame << "\033[H\033[2J";
    summaryFrame << "购买" << name << "成功" << '\n';
    if (lastPaid > 0)
    {
        summaryFrame << "实付 " << lastPaid << " 元" << '\n';
    }
    if (lastChange > 0)
    {
        summaryFrame << "找零 " << lastChange << " 元" << '\n';
    }
    NutritionManager *nutrition = NutritionManager::GetInstance();
    for (auto i : health)
    {
        int currentValue = (*User::GetInstance()->GetUserHealth())[i.GetInfo()];
        summaryFrame << nutrition->GetFeelingDescription(i.GetInfo(), currentValue) << '\n';
    }
    summaryFrame << "\n按回车返回";
    terminal.writeRaw(summaryFrame.str());
    waitForEnter();
    world->ChangeWhere(1);
}

EventCategory Commodity::GetCategory() const
{
    return EventCategory::Primary;
}

void SaveGame::ToDoIt()
{
    View *view = View::GetInstance();
    SaveManager *sm = SaveManager::GetInstance();
    const bool interactive = UiEnv::IsInteractive();

    if (!interactive)
    {
        // 非交互模式：简单的数字输入
        view->Clear();
        view->Print(0, 0, "===== 保存游戏 =====");
        auto descriptions = sm->GetAllSaveDescriptions();
        for (std::size_t i = 0; i < descriptions.size(); ++i)
        {
            view->Print(static_cast<int>(i + 1), 0,
                        std::to_string(i + 1) + ". " + descriptions[i]);
        }
        view->Print(static_cast<int>(descriptions.size() + 2), 0, "0. 返回");
        view->Print(static_cast<int>(descriptions.size() + 3), 0, "请选择槽位: ");

        int slot = view->Scan(0, 0);
        if (slot < 1 || slot > 3)
        {
            return;
        }

        if (sm->SaveGame(slot))
        {
            view->Clear();
            view->Print(0, 0, "游戏已保存到槽位 " + std::to_string(slot));
            view->WaitForEnter();
        }
        else
        {
            view->Clear();
            view->Print(0, 0, "保存失败，请检查磁盘空间");
            view->WaitForEnter();
        }
        return;
    }

    // 交互模式：使用方向键选择
    Terminal terminal;
    auto descriptions = sm->GetAllSaveDescriptions();
    std::vector<std::string> options;
    for (std::size_t i = 0; i < descriptions.size(); ++i)
    {
        options.push_back("槽位" + std::to_string(i + 1) + ": " + descriptions[i]);
    }
    options.push_back("返回");

    auto render = [&](int active, const std::string &) {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "===== 保存游戏 =====\n\n";
        for (std::size_t i = 0; i < options.size(); ++i)
        {
            frame << (static_cast<int>(i) == active ? "> " : "  ") << options[i] << "\n";
        }
        terminal.writeRaw(frame.str());
    };

    std::string typed;
    int choice = terminal.runMenu(options, render, typed);

    if (choice < 0 || choice >= static_cast<int>(descriptions.size()))
    {
        return; // 返回
    }

    int slot = choice + 1;
    if (sm->SaveGame(slot))
    {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "游戏已保存到槽位 " << slot << "\n\n按回车返回";
        terminal.writeRaw(frame.str());
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
                break;
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                throw std::runtime_error("INPUT_EOF");
        }
    }
    else
    {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "保存失败\n\n按回车返回";
        terminal.writeRaw(frame.str());
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
                break;
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                throw std::runtime_error("INPUT_EOF");
        }
    }
}

EventCategory SaveGame::GetCategory() const
{
    return EventCategory::Insight;
}

void LoadGame::ToDoIt()
{
    View *view = View::GetInstance();
    SaveManager *sm = SaveManager::GetInstance();
    const bool interactive = UiEnv::IsInteractive();

    if (!interactive)
    {
        // 非交互模式
        view->Clear();
        view->Print(0, 0, "===== 读取存档 =====");
        auto descriptions = sm->GetAllSaveDescriptions();
        bool hasAny = false;
        for (std::size_t i = 0; i < descriptions.size(); ++i)
        {
            view->Print(static_cast<int>(i + 1), 0,
                        std::to_string(i + 1) + ". " + descriptions[i]);
            if (sm->HasSave(static_cast<int>(i + 1)))
            {
                hasAny = true;
            }
        }
        view->Print(static_cast<int>(descriptions.size() + 2), 0, "0. 返回");

        if (!hasAny)
        {
            view->Print(static_cast<int>(descriptions.size() + 3), 0, "暂无存档");
            view->WaitForEnter();
            return;
        }

        view->Print(static_cast<int>(descriptions.size() + 3), 0, "请选择槽位: ");
        int slot = view->Scan(0, 0);

        if (slot < 1 || slot > 3)
        {
            return;
        }

        if (!sm->HasSave(slot))
        {
            view->Clear();
            view->Print(0, 0, "该槽位没有存档");
            view->WaitForEnter();
            return;
        }

        if (sm->LoadGame(slot))
        {
            view->Clear();
            view->Print(0, 0, "存档已加载");
            view->WaitForEnter();
        }
        else
        {
            view->Clear();
            view->Print(0, 0, "加载失败，存档可能已损坏");
            view->WaitForEnter();
        }
        return;
    }

    // 交互模式
    Terminal terminal;
    auto descriptions = sm->GetAllSaveDescriptions();
    std::vector<std::string> options;
    for (std::size_t i = 0; i < descriptions.size(); ++i)
    {
        options.push_back("槽位" + std::to_string(i + 1) + ": " + descriptions[i]);
    }
    options.push_back("返回");

    auto render = [&](int active, const std::string &) {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "===== 读取存档 =====\n\n";
        for (std::size_t i = 0; i < options.size(); ++i)
        {
            frame << (static_cast<int>(i) == active ? "> " : "  ") << options[i] << "\n";
        }
        terminal.writeRaw(frame.str());
    };

    std::string typed;
    int choice = terminal.runMenu(options, render, typed);

    if (choice < 0 || choice >= static_cast<int>(descriptions.size()))
    {
        return; // 返回
    }

    int slot = choice + 1;
    if (!sm->HasSave(slot))
    {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "该槽位没有存档\n\n按回车返回";
        terminal.writeRaw(frame.str());
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
                break;
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                throw std::runtime_error("INPUT_EOF");
        }
        return;
    }

    if (sm->LoadGame(slot))
    {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "存档已加载\n\n按回车返回";
        terminal.writeRaw(frame.str());
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
                break;
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                throw std::runtime_error("INPUT_EOF");
        }
    }
    else
    {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "加载失败，存档可能已损坏\n\n按回车返回";
        terminal.writeRaw(frame.str());
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
                break;
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                throw std::runtime_error("INPUT_EOF");
        }
    }
}

EventCategory LoadGame::GetCategory() const
{
    return EventCategory::Insight;
}

void SeeDoctor::ToDoIt()
{
    View *view = View::GetInstance();
    World *world = World::GetInstance();
    DiseaseManager *disease = DiseaseManager::GetInstance();
    const bool interactive = UiEnv::IsInteractive();

    // 获取当前疾病列表
    auto diseases = disease->GetActiveDiseases();

    if (diseases.empty())
    {
        if (!interactive)
        {
            view->Clear();
            view->Print(0, 0, "医生说：你很健康，不需要治疗。");
            view->WaitForEnter();
        }
        else
        {
            Terminal terminal;
            std::ostringstream frame;
            frame << "\033[H\033[2J";
            frame << "医生说：你很健康，不需要治疗。\n\n按回车返回";
            terminal.writeRaw(frame.str());
            while (true)
            {
                Terminal::KeyEvent event = terminal.readKey();
                if (event.type == Terminal::KeyEvent::Type::Enter)
                    break;
                if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                    throw std::runtime_error("INPUT_EOF");
            }
        }
        return;
    }

    // 构建疾病选项列表（包含治疗费用）
    std::vector<std::string> options;
    std::vector<std::string> diseaseIds;

    // 根据疾病名称查找 ID
    auto findDiseaseId = [&disease](const std::string &name) -> std::string {
        // 遍历常见疾病ID
        const std::vector<std::string> ids = {"cold", "anemia", "scurvy", "food_poisoning", "malnutrition", "depression"};
        for (const auto &id : ids)
        {
            const Disease *d = disease->GetDiseaseInfo(id);
            if (d && d->name == name)
            {
                return id;
            }
        }
        return "";
    };

    for (const auto &diseaseName : diseases)
    {
        std::string id = findDiseaseId(diseaseName);
        if (!id.empty())
        {
            const Disease *info = disease->GetDiseaseInfo(id);
            if (info)
            {
                options.push_back(diseaseName + " - 治疗费: ¥" + std::to_string(info->treatmentCost));
                diseaseIds.push_back(id);
            }
        }
    }
    options.push_back("暂不治疗");

    if (!interactive)
    {
        view->Clear();
        int line = 0;
        view->Print(line++, 0, "===== 医院 =====");
        view->Print(line++, 0, "你的账户余额: ¥" + std::to_string(world->GetWallet()));
        line++;
        view->Print(line++, 0, "你有以下疾病需要治疗：");
        for (std::size_t i = 0; i < options.size(); ++i)
        {
            view->Print(line++, 0, std::to_string(i) + ". " + options[i]);
        }
        view->Print(line++, 0, "请选择要治疗的疾病：");

        int choice = view->Scan(0, 0);
        if (choice < 0 || choice >= static_cast<int>(diseaseIds.size()))
        {
            return; // 暂不治疗
        }

        const std::string &selectedId = diseaseIds[static_cast<std::size_t>(choice)];
        if (disease->TreatAtHospital(selectedId))
        {
            view->Clear();
            view->Print(0, 0, "治疗成功！你的" + diseases[static_cast<std::size_t>(choice)] + "已经痊愈。");
            view->WaitForEnter();
        }
        else
        {
            view->Clear();
            view->Print(0, 0, "治疗失败，可能是资金不足。");
            view->WaitForEnter();
        }
        return;
    }

    // 交互模式
    Terminal terminal;
    std::string typed;

    auto render = [&](int active, const std::string &) {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "===== 医院 =====\n";
        frame << "你的账户余额: ¥" << world->GetWallet() << "\n\n";
        frame << "你有以下疾病需要治疗：\n\n";
        for (std::size_t i = 0; i < options.size(); ++i)
        {
            frame << (static_cast<int>(i) == active ? "> " : "  ") << options[i] << "\n";
        }
        terminal.writeRaw(frame.str());
    };

    int choice = terminal.runMenu(options, render, typed);

    if (choice < 0 || choice >= static_cast<int>(diseaseIds.size()))
    {
        return; // 暂不治疗
    }

    const std::string &selectedId = diseaseIds[static_cast<std::size_t>(choice)];
    if (disease->TreatAtHospital(selectedId))
    {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "治疗成功！你的" << diseases[static_cast<std::size_t>(choice)] << "已经痊愈。\n\n按回车返回";
        terminal.writeRaw(frame.str());
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
                break;
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                throw std::runtime_error("INPUT_EOF");
        }
    }
    else
    {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << "治疗失败，可能是资金不足。\n\n按回车返回";
        terminal.writeRaw(frame.str());
        while (true)
        {
            Terminal::KeyEvent event = terminal.readKey();
            if (event.type == Terminal::KeyEvent::Type::Enter)
                break;
            if (event.type == Terminal::KeyEvent::Type::EndOfInput)
                throw std::runtime_error("INPUT_EOF");
        }
    }
}

EventCategory SeeDoctor::GetCategory() const
{
    return EventCategory::Primary;
}
