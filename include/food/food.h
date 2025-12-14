#pragma once

#include "health/health.h"
#include <string>
#include <vector>

class Health;

/**
 * @description: 食物类型枚举，决定新鲜度衰减速度
 */
enum class FoodType
{
    Fresh,      // 生鲜食品：25小时完全腐烂（衰减率 4.0/小时）
    Beverage,   // 饮料：50小时（衰减率 2.0/小时）
    Processed,  // 加工食品：100小时（衰减率 1.0/小时）
    Canned      // 罐头食品：1000小时（衰减率 0.1/小时）
};

/**
 * @description: 食物类，包含新鲜度系统
 */
class Food
{
  public:
    /**
     * @description: 构造函数
     * @param name 食物名称
     * @param health 营养成分列表
     * @param price 价格
     * @param type 食物类型
     * @param initialFreshness 初始新鲜度（默认100）
     */
    Food(const std::string &name,
         const std::vector<Health> &health,
         int price,
         FoodType type,
         int initialFreshness = 100);

    ~Food() = default;

    // 新鲜度管理
    /**
     * @description: 获取当前新鲜度
     * @return 新鲜度值 (可为负数表示已腐烂)
     */
    int GetFreshness() const { return freshness; }

    /**
     * @description: 基于游戏时间更新新鲜度
     * @param gameMinutes 游戏内经过的分钟数
     */
    void UpdateFreshness(int gameMinutes);

    /**
     * @description: 检查是否已过期
     * @return true=已过期
     */
    bool IsExpired() const { return freshness <= 0; }

    /**
     * @description: 检查是否快要过期
     * @return true=快过期（新鲜度<=25）
     */
    bool IsSpoiling() const { return freshness <= 25 && freshness > 0; }

    /**
     * @description: 获取新鲜度等级描述
     * @return 描述字符串（新鲜/较新鲜/不太新鲜/快过期/已过期）
     */
    std::string GetFreshnessLabel() const;

    /**
     * @description: 获取考虑新鲜度后的实际营养效果
     * @return 营养列表（值可能被折扣或变为负数）
     */
    std::vector<Health> GetEffectiveHealth() const;

    // Getters
    std::string GetName() const { return name; }
    int GetPrice() const { return price; }
    FoodType GetType() const { return type; }
    const std::vector<Health> &GetBaseHealth() const { return baseHealth; }

    /**
     * @description: 获取显示信息（包含新鲜度）
     * @return 格式化的显示字符串
     */
    std::string GetInfo() const;

    /**
     * @description: 获取营养乘数（基于新鲜度）
     * @return 乘数值（0.0-1.0，过期时为负数）
     */
    float GetNutritionMultiplier() const;

    /**
     * @description: 检查食用后是否会导致食物中毒
     * @return 食物中毒概率 (0.0-1.0)
     */
    float GetFoodPoisoningRisk() const;

  private:
    std::string name;              // 食物名称
    std::vector<Health> baseHealth;// 基础营养成分
    int price;                     // 价格
    FoodType type;                 // 食物类型
    int freshness;                 // 当前新鲜度 (0-100，可为负表示腐烂程度)
    int maxFreshness;              // 最大新鲜度

    /**
     * @description: 获取基于食物类型的衰减速率
     * @return 每小时衰减的新鲜度值
     */
    float GetDecayRate() const;
};

/**
 * @description: 将 FoodType 转换为字符串
 * @param type 食物类型
 * @return 类型名称字符串
 */
std::string FoodTypeToString(FoodType type);
