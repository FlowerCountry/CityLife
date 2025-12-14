#include "user/user.h"

#include <algorithm>

User::User()
{
    // 核心属性（Primary）- 直接影响生存
    UserHealth = {
        {"饱腹感", 100},      // 负数下限 -30
        {"饥渴", 100},        // 负数下限 -20（原名"饥饿"改为"饥渴"表示水分需求）
        {"蛋白质", 100},      // 负数下限 -15
        {"碳水化合物", 100},  // 负数下限 -15

        // 微量元素（Micro）- 长期健康
        {"钙", 100},          // 负数下限 -10
        {"糖分", 100},        // 负数下限 -10
        {"脂肪", 100},        // 负数下限 0（脂肪不可为负）
        {"纤维素", 100},      // 负数下限 -10
        {"铁", 100},          // 负数下限 -10
        {"维生素A", 100},     // 负数下限 -10
        {"维生素B", 100},     // 负数下限 -10
        {"维生素C", 100},     // 负数下限 -10
        {"维生素D", 100},     // 负数下限 -10
        {"维生素E", 100},     // 负数下限 -10
        {"钾", 100},          // 负数下限 -10
        {"硒", 100},          // 负数下限 -10
        {"锌", 100},          // 负数下限 -10

        // 特殊属性（Special）- 心理/状态
        {"幸福感", 100},      // 负数下限 -20
        {"精神振奋", 100},    // 负数下限 -10

        // 食物定义中使用的额外属性
        {"益生菌", 100},      // 负数下限 -10
        {"镁", 100},          // 负数下限 -10
        {"饥饿", 100}         // 兼容旧食物定义（作为"饱腹感"的别名处理）
    };
}

User::~User() {}

User *User::instance = nullptr;

User *User::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new User();
    }
    return instance;
}

int User::GetNutrition(const std::string &name) const
{
    auto it = UserHealth.find(name);
    if (it != UserHealth.end())
    {
        return it->second;
    }
    return 0;  // 未知营养默认返回0
}

void User::SetNutrition(const std::string &name, int value)
{
    // 获取该营养的负数下限
    int minVal = GetMinValue(name);
    // 应用下限限制
    UserHealth[name] = std::max(value, minVal);
}

void User::ConsumeNutrition(const std::string &name, int amount)
{
    int current = GetNutrition(name);
    int minVal = GetMinValue(name);
    UserHealth[name] = std::max(current - amount, minVal);
}

void User::AddNutrition(const std::string &name, int amount)
{
    int current = GetNutrition(name);
    // 上限固定为100
    UserHealth[name] = std::min(current + amount, 100);
}

std::vector<std::string> User::GetAllNutritionNames() const
{
    std::vector<std::string> names;
    names.reserve(UserHealth.size());
    for (const auto &pair : UserHealth)
    {
        names.push_back(pair.first);
    }
    return names;
}

bool User::IsAlive() const
{
    // 检查核心属性是否低于死亡线
    // 饱腹感 < -30 或 饥渴 < -20 会导致死亡
    if (GetNutrition("饱腹感") <= -30)
    {
        return false;
    }
    if (GetNutrition("饥渴") <= -20)
    {
        return false;
    }
    // 蛋白质和碳水化合物低于 -15 也会死亡
    if (GetNutrition("蛋白质") <= -15)
    {
        return false;
    }
    if (GetNutrition("碳水化合物") <= -15)
    {
        return false;
    }
    return true;
}

float User::GetDeathRisk() const
{
    float risk = 0.0f;

    // 检查核心属性的濒死区间（-10 到 -20 之间）
    // 每个属性在濒死区间贡献 5% 死亡风险
    int satiety = GetNutrition("饱腹感");
    if (satiety < -10 && satiety > -30)
    {
        risk += 0.05f;
    }

    int thirst = GetNutrition("饥渴");
    if (thirst < -10 && thirst > -20)
    {
        risk += 0.05f;
    }

    int protein = GetNutrition("蛋白质");
    if (protein < -10 && protein > -15)
    {
        risk += 0.03f;
    }

    int carbs = GetNutrition("碳水化合物");
    if (carbs < -10 && carbs > -15)
    {
        risk += 0.03f;
    }

    return std::min(risk, 1.0f);
}

int User::GetNutritionLevel(const std::string &name) const
{
    int value = GetNutrition(name);

    if (value > 50)
    {
        return 0;  // 安全区 (50-100)
    }
    else if (value > 20)
    {
        return 1;  // 警告区 (20-50)
    }
    else if (value > 0)
    {
        return 2;  // 危险区 (0-20)
    }
    else if (value > -10)
    {
        return 3;  // 透支区 (-10 到 0)
    }
    else
    {
        return 4;  // 濒死区 (< -10)
    }
}

int User::GetMinValue(const std::string &name) const
{
    // 核心属性的负数下限
    if (name == "饱腹感")
    {
        return -30;
    }
    if (name == "饥渴" || name == "饥饿")
    {
        return -20;
    }
    if (name == "蛋白质" || name == "碳水化合物")
    {
        return -15;
    }
    // 特殊属性
    if (name == "幸福感")
    {
        return -20;
    }
    if (name == "脂肪")
    {
        return 0;  // 脂肪不可为负
    }
    // 其他微量元素默认下限 -10
    return -10;
}