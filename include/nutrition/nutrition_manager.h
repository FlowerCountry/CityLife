#pragma once

#include <map>
#include <string>
#include <vector>
#include <utility>

/**
 * @description: 营养类别枚举
 */
enum class NutritionCategory
{
    Primary,  // 核心属性：饱腹感、饥渴、蛋白质、碳水化合物
    Micro,    // 微量元素：维生素、矿物质等
    Special   // 特殊属性：幸福感、精神振奋等
};

/**
 * @description: 营养配置结构
 */
struct NutritionConfig
{
    std::string name;          // 营养名称
    NutritionCategory category;// 类别
    float baseDecay;           // 每次行动的基础衰减值
    int minValue;              // 最小值（负数下限）
    int maxValue;              // 最大值（默认100）
};

/**
 * @description: 营养相互影响规则
 */
struct NutritionInteraction
{
    std::string source;   // 源营养
    std::string target;   // 目标营养
    float ratio;          // 影响比例
    int threshold;        // 触发阈值
    bool isPositive;      // true=正向影响, false=负向影响
};

/**
 * @description: 营养管理器（单例）
 * 负责处理营养衰减、相互影响和状态检查
 */
class NutritionManager
{
  public:
    /**
     * @description: 获取单例实例
     * @return 单例指针
     */
    static NutritionManager *GetInstance();

    /**
     * @description: 每次行动触发的营养衰减
     * @param actionCost 行动消耗系数（默认1.0）
     */
    void OnAction(float actionCost = 1.0f);

    /**
     * @description: 处理营养相互影响
     */
    void ProcessInteractions();

    /**
     * @description: 获取营养配置
     * @param name 营养名称
     * @return 营养配置
     */
    const NutritionConfig *GetConfig(const std::string &name) const;

    /**
     * @description: 获取所有警告信息
     * @return 警告信息列表
     */
    std::vector<std::string> GetWarnings() const;

    /**
     * @description: 获取所有危急信息
     * @return 危急信息列表
     */
    std::vector<std::string> GetCriticals() const;

    /**
     * @description: 计算当前衰减倍率（基于负数区间）
     * @param name 营养名称
     * @return 衰减倍率
     */
    float GetDecayMultiplier(const std::string &name) const;

    /**
     * @description: 获取营养状态摘要字符串
     * @return 状态摘要
     */
    std::string GetStatusSummary() const;

    /**
     * @description: 检查玩家是否存活（委托给User）
     * @return true=存活
     */
    bool IsAlive() const;

    /**
     * @description: 获取死亡风险
     * @return 风险值 (0.0-1.0)
     */
    float GetDeathRisk() const;

    /**
     * @description: 随机检查是否触发死亡
     * @return true=触发死亡
     */
    bool CheckRandomDeath();

    /**
     * @description: 检测营养等级变化，生成提示
     */
    void CheckLevelChanges();

    /**
     * @description: 获取待显示的提示列表
     * @return 提示字符串列表
     */
    std::vector<std::string> GetPendingHints() const;

    /**
     * @description: 清空待显示提示队列
     */
    void ClearPendingHints();

    /**
     * @description: 根据营养值获取沉浸式感受描述（用于购买反馈）
     * @param name 营养名称
     * @param value 当前营养值
     * @return 感受描述字符串
     */
    std::string GetFeelingDescription(const std::string &name, int value) const;

  private:
    NutritionManager();
    static NutritionManager *instance;

    // 营养配置表
    std::map<std::string, NutritionConfig> configs;

    // 营养相互影响规则
    std::vector<NutritionInteraction> interactions;

    // 上一次各营养的等级（用于检测首次进入阈值）
    std::map<std::string, int> previousLevels;

    // 待显示的提示队列
    std::vector<std::string> pendingHints;

    // 提示配置表：<营养名称, 等级> -> 提示文案
    std::map<std::pair<std::string, int>, std::string> hintMessages;

    /**
     * @description: 初始化营养配置
     */
    void InitConfigs();

    /**
     * @description: 初始化相互影响规则
     */
    void InitInteractions();

    /**
     * @description: 初始化提示文案配置
     */
    void InitHintConfigs();

    /**
     * @description: 初始化等级记录（全部为安全等级0）
     */
    void InitPreviousLevels();

    /**
     * @description: 应用单个营养的衰减
     * @param name 营养名称
     * @param multiplier 额外倍率
     */
    void ApplyDecay(const std::string &name, float multiplier);
};
