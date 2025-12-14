#include "food/food.h"
#include "health/health.h"
#include "locale/locale_manager.h"

#include <algorithm>
#include <sstream>

Food::Food(const std::string &name,
           const std::vector<Health> &health,
           int price,
           FoodType type,
           int initialFreshness)
    : name(name),
      baseHealth(health),
      price(price),
      type(type),
      freshness(initialFreshness),
      maxFreshness(initialFreshness)
{
}

void Food::UpdateFreshness(int gameMinutes)
{
    // 将分钟转换为小时
    float hours = static_cast<float>(gameMinutes) / 60.0f;
    float decay = GetDecayRate() * hours;
    freshness -= static_cast<int>(decay);
    // 不设下限，允许负数表示腐烂程度
}

std::string Food::GetFreshnessLabel() const
{
    LocaleManager *locale = LocaleManager::GetInstance();
    if (freshness > 75)
    {
        return locale->Get("food.freshness", "fresh");
    }
    else if (freshness > 50)
    {
        return locale->Get("food.freshness", "fairly_fresh");
    }
    else if (freshness > 25)
    {
        return locale->Get("food.freshness", "not_fresh");
    }
    else if (freshness > 0)
    {
        return locale->Get("food.freshness", "expiring");
    }
    else
    {
        return locale->Get("food.freshness", "expired");
    }
}

float Food::GetNutritionMultiplier() const
{
    if (freshness > 75)
    {
        return 1.0f;  // 100% 营养
    }
    else if (freshness > 50)
    {
        return 0.8f;  // 80% 营养
    }
    else if (freshness > 25)
    {
        return 0.5f;  // 50% 营养
    }
    else if (freshness > 0)
    {
        return 0.2f;  // 20% 营养
    }
    else
    {
        return -0.5f; // 过期：负效果
    }
}

std::vector<Health> Food::GetEffectiveHealth() const
{
    std::vector<Health> result;
    float multiplier = GetNutritionMultiplier();

    if (multiplier > 0)
    {
        // 正常效果：按比例计算营养
        for (const auto &h : baseHealth)
        {
            int effectiveValue = static_cast<int>(h.GetReserves() * multiplier);
            if (effectiveValue > 0)
            {
                result.push_back(Health(h.GetInfo(), effectiveValue));
            }
        }
    }
    else
    {
        // 过期食物：减少饱腹感和幸福感，其他营养为0
        result.push_back(Health("饱腹感", -10));
        result.push_back(Health("幸福感", -15));
        result.push_back(Health("饥渴", -5));
    }

    return result;
}

float Food::GetFoodPoisoningRisk() const
{
    if (freshness > 25)
    {
        return 0.0f;  // 不会食物中毒
    }
    else if (freshness > 0)
    {
        return 0.1f;  // 10% 概率
    }
    else if (freshness > -25)
    {
        return 0.3f;  // 30% 概率
    }
    else if (freshness > -50)
    {
        return 0.6f;  // 60% 概率
    }
    else
    {
        return 0.9f;  // 90% 概率
    }
}

std::string Food::GetInfo() const
{
    std::ostringstream ss;
    ss << name << " " << price << "元";

    // 添加新鲜度标签
    std::string label = GetFreshnessLabel();
    if (label != "新鲜")
    {
        ss << " [" << label << "]";
    }

    return ss.str();
}

float Food::GetDecayRate() const
{
    switch (type)
    {
    case FoodType::Fresh:
        return 4.0f;  // 25小时完全腐烂
    case FoodType::Beverage:
        return 2.0f;  // 50小时
    case FoodType::Processed:
        return 1.0f;  // 100小时
    case FoodType::Canned:
        return 0.1f;  // 1000小时
    default:
        return 1.0f;
    }
}

std::string FoodTypeToString(FoodType type)
{
    LocaleManager *locale = LocaleManager::GetInstance();
    switch (type)
    {
    case FoodType::Fresh:
        return locale->Get("food.type", "fresh");
    case FoodType::Beverage:
        return locale->Get("food.type", "beverage");
    case FoodType::Processed:
        return locale->Get("food.type", "processed");
    case FoodType::Canned:
        return locale->Get("food.type", "canned");
    default:
        return locale->Get("food.type", "unknown");
    }
}
