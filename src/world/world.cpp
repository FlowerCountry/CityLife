#include "world/world.h"
#include "payment/cashier.h"
#include "bank/bank.h"
#include "controller/controller.h"
#include "disease/disease_manager.h"
#include "health/health.h"
#include "locale/locale_manager.h"
#include "nutrition/nutrition_manager.h"
#include "object/object.h"
#include "view/view.h"
#include "user/user.h"
#include "center/center.h"

#include <algorithm>
#include <clocale>

// 定义支持的纸币面额（从大到小）。
const std::array<int, 6> World::denominations = {100, 50, 20, 10, 5, 1};

World::World()
{
    setlocale(LC_ALL, "");
    year = 2010;
    month = 10;
    day = 10;
    hour = 11;
    minute = 3;
    second = 0;
    where = 0;
    // 初始给予玩家一张 100 元纸币，其他面额为空。
    wallet = {1, 0, 0, 0, 0, 0};
    RefreshWalletTotal();
    bank = new Bank();
    BuildingNames = {
        "市中心",
        "超市",
        "银行",
        "超市内",
        "医院",
    };
    Buildings = {
        new class Building(0, 5, 5, "市中心"),
        new class Building(1, 3, 5, "超市"),
        new class Building(2, 5, 3, "银行"),
        new class Building(4, 7, 5, "医院"),
    };
    ToDoThings = {
        {new class Information("查看公告", "公告")},
        {new class Buy("购买物品")},
        {new class DepositingMoney(), new class WithdrawMoney(), new class CheckBankBalance()},
        {new class Commodity("面包", 15, {Health("饱腹感", 18), Health("饥饿", 18), Health("蛋白质", 9), Health("维生素B", 9)}),
         new class Commodity("牛奶", 25, {Health("饱腹感", 10), Health("饥饿", 16), Health("钙", 24), Health("蛋白质", 14), Health("维生素A", 12)}),
         new class Commodity("蛋糕", 40, {Health("饱腹感", 24), Health("饥饿", 10), Health("糖分", 30), Health("脂肪", 20), Health("维生素E", 14)}),
         new class Commodity("苹果", 10, {Health("饱腹感", 16), Health("维生素C", 20), Health("纤维素", 14), Health("钾", 10)}),
         new class Commodity("牛排", 150, {Health("饱腹感", 20), Health("蛋白质", 25), Health("铁", 15), Health("脂肪", 8), Health("锌", 32)}),
         new class Commodity("橙汁", 18, {Health("饱腹感", 7), Health("维生素C", 22), Health("饥渴", 14), Health("糖分", 7), Health("钙", 10)}),
         new class Commodity("沙拉", 35, {Health("饱腹感", 9), Health("维生素A", 30), Health("纤维素", 20), Health("钾", 14)}),
         new class Commodity("鸡蛋", 12, {Health("饱腹感", 16), Health("蛋白质", 24), Health("硒", 12)}),
         new class Commodity("意大利面", 30, {Health("饱腹感", 24), Health("碳水化合物", 30), Health("维生素B", 12)}),
         new class Commodity("咖啡", 20, {Health("饱腹感", 4), Health("精神振奋", 20), Health("饥渴", 10), Health("钾", 12)}),
         new class Commodity("巧克力", 25, {Health("饱腹感", 6), Health("幸福感", 24), Health("糖分", 30), Health("脂肪", 16), Health("镁", 12)}),
         new class Commodity("番茄", 7, {Health("饱腹感", 8), Health("维生素C", 16), Health("纤维素", 12), Health("钾", 8)}),
         new class Commodity("鸡胸肉", 70, {Health("饱腹感", 24), Health("蛋白质", 36), Health("脂肪", 6), Health("铁", 26)}),
         new class Commodity("矿泉水", 3, {Health("饱腹感", 2), Health("饥渴", 45), Health("钙", 5)}),
         new class Commodity("米饭", 15, {Health("饱腹感", 24), Health("碳水化合物", 36), Health("维生素B", 12)}),
         new class Commodity("燕麦片", 22, {Health("饱腹感", 28), Health("纤维素", 20), Health("钾", 14)}),
         new class Commodity("酸奶", 18, {Health("饱腹感", 12), Health("钙", 24), Health("蛋白质", 16), Health("益生菌", 16)}),
         new class Commodity("三明治", 40, {Health("饱腹感", 32), Health("蛋白质", 20), Health("铁", 18)}),
         new class Commodity("薯片", 12, {Health("饱腹感", 6), Health("脂肪", 24), Health("碳水化合物", 16), Health("维生素C", 12)}),
         new class Commodity("冰淇淋", 35, {Health("饱腹感", 6), Health("幸福感", 30), Health("糖分", 36), Health("脂肪", 20), Health("维生素D", 14)})},
        {new class SeeDoctor()},  // 医院的操作
    };

    for (auto &list : ToDoThings)
    {
        list.push_back(new class CheckCash());
    }

    // 仅在市中心（ToDoThings[0]）添加存档/读档功能
    ToDoThings[0].push_back(new class SaveGame());
    ToDoThings[0].push_back(new class LoadGame());

    for (auto i : Buildings)
    {
        for (auto j : Buildings)
        {
            if (i->id != j->id)
            {
                ToDoThings[i->id].push_back(new class GoWhere("前往" + j->name, i->id, j->id, GetLen(*i, *j)));
            }
        }
    }

    auto categoryRank = [](EventCategory cat) {
        switch (cat)
        {
        case EventCategory::Primary:
            return 0;
        case EventCategory::Insight:
            return 1;
        case EventCategory::Navigation:
            return 2;
        }
        return 3;
    };

    for (auto &list : ToDoThings)
    {
        std::stable_sort(list.begin(), list.end(), [&](Object *lhs, Object *rhs) {
            return categoryRank(lhs->GetCategory()) < categoryRank(rhs->GetCategory());
        });
    }
}
World::~World()
{
    delete bank;
    for (auto it : Buildings) delete it;
    for (auto &list : ToDoThings)
    {
        for (auto item : list) delete item;
    }
}

World *World::instance = nullptr;

World *World::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new World();
    }
    return instance;
}

void World::Start()
{
    NutritionManager *nutrition = NutritionManager::GetInstance();
    DiseaseManager *disease = DiseaseManager::GetInstance();
    View *view = View::GetInstance();
    LocaleManager *locale = LocaleManager::GetInstance();

    while (true)
    {
        // 获取并构建身体信号提示
        std::string header;
        auto hints = nutrition->GetPendingHints();
        if (!hints.empty())
        {
            header += locale->Get("game.status", "body_signal") + "\n";
            for (const auto &hint : hints)
            {
                header += hint + "\n";
            }
            header += "\n";
            nutrition->ClearPendingHints();
        }

        // 构建位置标题
        std::string locationPrefix = locale->Get("game.location", "current");
        // 确保前缀与地点名之间有空格（INI 解析会去掉末尾空白）
        if (!locationPrefix.empty() && locationPrefix.back() != ' ')
        {
            locationPrefix += " ";
        }
        header += locationPrefix + BuildingNames[where];

        // 添加疾病状态
        auto diseases = disease->GetActiveDiseases();
        if (!diseases.empty())
        {
            header += "\n" + locale->Get("game.status", "disease") + disease->GetStatusSummary();
        }

        // 添加营养状态警告
        auto criticals = nutrition->GetCriticals();
        if (!criticals.empty())
        {
            header += "\n" + locale->Get("game.status", "danger") + criticals[0];
            if (criticals.size() > 1)
            {
                header += locale->Get("game.status", "items_suffix") + std::to_string(criticals.size()) + locale->Get("game.status", "items_count");
            }
        }
        else
        {
            auto warnings = nutrition->GetWarnings();
            if (!warnings.empty())
            {
                header += "\n" + locale->Get("game.status", "warning") + warnings[0];
                if (warnings.size() > 1)
                {
                    header += locale->Get("game.status", "items_suffix") + std::to_string(warnings.size()) + locale->Get("game.status", "items_count");
                }
            }
        }

        int choose = Controller::GetInstance()->Choose(header, ToDoThings[where]);
        ToDoThings[where][choose]->ToDoIt();

        // 获取疾病导致的行动消耗倍率
        float diseaseMultiplier = disease->GetActionCostMultiplier();

        // 每次行动后触发营养衰减（考虑疾病倍率）
        nutrition->OnAction(diseaseMultiplier);

        // 处理营养相互影响
        nutrition->ProcessInteractions();

        // 检测营养等级变化，生成下一轮提示
        nutrition->CheckLevelChanges();

        // 检查疾病状态（触发和康复）
        disease->OnAction();

        // 检查随机死亡风险
        if (nutrition->CheckRandomDeath())
        {
            view->Clear();
            view->Print(0, 0, locale->Get("game.death", "random"));
            view->Print(1, 0, locale->Get("game.death", "malnutrition"));
            view->Print(3, 0, locale->Get("game.death", "game_over"));
            view->WaitForEnter();
            break;
        }

        // 检查是否死亡
        if (!nutrition->IsAlive())
        {
            view->Clear();
            view->Print(0, 0, locale->Get("game.death", "starved"));
            view->Print(2, 0, locale->Get("game.death", "game_over"));
            view->WaitForEnter();
            break;
        }
    }
}

void World::ChangeWhere(int where)
{
    this->where = where;
}

bool World::SpendMoney(int money)
{
    LocaleManager *locale = LocaleManager::GetInstance();
    if (this->money < money)
    {
        View::GetInstance()->Clear();
        View::GetInstance()->Print(0, 0, locale->Get("bank.error", "insufficient"));
        View::GetInstance()->WaitForEnter();
        return false;
    }

    Cashier cashier(denominations);
    const Cashier::Result result = cashier.pay(money, wallet);
    if (!result.success)
    {
        View::GetInstance()->Clear();
        View::GetInstance()->Print(0, 0, locale->Get("bank.error", "cannot_compose"));
        View::GetInstance()->WaitForEnter();
        return false;
    }

    for (std::size_t i = 0; i < wallet.size(); ++i)
    {
        wallet[i] -= result.used[i];
    }

    const bool hasChange = std::any_of(result.change.begin(), result.change.end(), [](int count) {
        return count > 0;
    });

    if (hasChange)
    {
        AddBills(result.change);
    }
    else
    {
        RefreshWalletTotal();
    }

    return true;
}

bool World::DepositingMoney(int money)
{
    LocaleManager *locale = LocaleManager::GetInstance();
    if (this->money < money)
    {
        View::GetInstance()->Clear();
        View::GetInstance()->Print(0, 0, locale->Get("bank.error", "insufficient"));
        View::GetInstance()->WaitForEnter();
        return false;
    }

    if (money % 100 != 0)
    {
        View::GetInstance()->Clear();
        View::GetInstance()->Print(0, 0, locale->Get("bank.error", "deposit_hundreds"));
        View::GetInstance()->WaitForEnter();
        return false;
    }

    // 银行只接收整百纸币，不允许由小面额拼凑。
    const int neededHundreds = money / 100;
    if (wallet[0] < neededHundreds)
    {
        View::GetInstance()->Clear();
        View::GetInstance()->Print(0, 0, locale->Get("bank.error", "no_hundreds"));
        View::GetInstance()->WaitForEnter();
        return false;
    }

    std::array<int, 6> depositBills{};
    depositBills[0] = neededHundreds;
    if (!RemoveBills(depositBills))
    {
        // 正常情况下不会触发，属于数据不一致的兜底保护。
        View::GetInstance()->Clear();
        View::GetInstance()->Print(0, 0, locale->Get("bank.error", "remove_failed"));
        View::GetInstance()->WaitForEnter();
        return false;
    }

    bank->DepositingMoney(money);
    return true;
}

bool World::WithdrawMoney(int money)
{
    LocaleManager *locale = LocaleManager::GetInstance();
    if (money % 100 != 0)
    {
        View::GetInstance()->Clear();
        View::GetInstance()->Print(0, 0, locale->Get("bank.error", "withdraw_hundreds"));
        View::GetInstance()->WaitForEnter();
        return false;
    }

    if (bank->GetMoney() >= money)
    {
        bank->WithdrawMoney(money);

        std::array<int, 6> gained{};
        int remaining = money;
        for (std::size_t i = 0; i < denominations.size(); ++i)
        {
            const int denom = denominations[i];
            int take = remaining / denom;
            gained[i] = take;
            remaining -= take * denom;
        }
        // 取款默认按最大面额分配。
        AddBills(gained);
        return true;
    }

    View::GetInstance()->Clear();
    View::GetInstance()->Print(0, 0, locale->Get("bank.error", "bank_insufficient"));
    View::GetInstance()->WaitForEnter();
    return false;
}

int World::GetWallet() const
{
    return money;
}

const std::array<int, 6> &World::GetWalletBreakdown() const
{
    // 返回当前各面额纸币的张数分布。
    return wallet;
}

bool World::RemoveBills(const std::array<int, 6> &counts)
{
    // 检查是否持有足够张数的指定面额。
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        if (counts[i] > wallet[i])
        {
            return false;
        }
    }

    // 扣减对应张数后刷新总金额。
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        wallet[i] -= counts[i];
    }
    RefreshWalletTotal();
    return true;
}

void World::AddBills(const std::array<int, 6> &counts)
{
    // 将指定面额张数加入钱包后刷新总金额。
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        wallet[i] += counts[i];
    }
    RefreshWalletTotal();
}

bool World::RemoveAmount(int amount)
{
    std::array<int, 6> temp = wallet;

    // 贪心策略：优先使用大面额钞票尝试扣款。
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        const int denom = denominations[i];
        int maxUse = std::min(amount / denom, temp[i]);
        if (maxUse > 0)
        {
            amount -= maxUse * denom;
            temp[i] -= maxUse;
        }
    }

    if (amount != 0)
    {
        return false;
    }

    wallet = temp;
    RefreshWalletTotal();
    return true;
}

int World::SumWallet(const std::array<int, 6> &counts)
{
    // 根据面额数组与张数求和。
    int total = 0;
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        total += counts[i] * denominations[i];
    }
    return total;
}

void World::RefreshWalletTotal()
{
    // 重新计算钱袋总额，确保 GetWallet() 返回值正确。
    money = SumWallet(wallet);
}

int World::GetBankBalance() const
{
    return bank->GetMoney();
}

int World::GetLen(Building a, Building b)
{
    return abs(a.x - b.x) + abs(a.y - b.y);
}

void World::UpdateTime(int seconds)
{
    second += seconds;
    while (second >= 60)
    {
        second -= 60;
        minute++;
    }
    while (minute >= 60)
    {
        minute -= 60;
        hour++;
    }
    while (hour >= 24)
    {
        hour -= 24;
        day++;
    }
    // 简化：每月固定30天
    while (day > 30)
    {
        day -= 30;
        month++;
    }
    while (month > 12)
    {
        month -= 12;
        year++;
    }
}

// ===== 存档系统 getter 实现 =====
int World::GetYear() const { return year; }
int World::GetMonth() const { return month; }
int World::GetDay() const { return day; }
int World::GetHour() const { return hour; }
int World::GetMinute() const { return minute; }
int World::GetSecond() const { return second; }
int World::GetWhere() const { return where; }
float World::GetLifeQuality() const { return LifeQuality; }

// ===== 存档系统 setter 实现 =====
void World::SetTime(int y, int m, int d, int h, int min, int sec)
{
    year = y;
    month = m;
    day = d;
    hour = h;
    minute = min;
    second = sec;
}

void World::SetWhere(int w)
{
    where = w;
}

void World::SetLifeQuality(float quality)
{
    LifeQuality = quality;
}

void World::SetWallet(const std::array<int, 6> &w)
{
    wallet = w;
    RefreshWalletTotal();
}

void World::SetBankDeposit(int amount)
{
    bank->SetMoney(amount);
}
