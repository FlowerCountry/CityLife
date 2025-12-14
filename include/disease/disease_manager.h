#pragma once

#include <map>
#include <string>
#include <vector>

/**
 * @description: 疾病严重程度枚举
 */
enum class DiseaseSeverity
{
    Mild,      // 轻微（可能自愈）
    Moderate,  // 中等（需要治疗）
    Severe     // 严重（需要住院）
};

/**
 * @description: 疾病触发条件
 */
struct DiseaseTrigger
{
    std::string nutritionName;  // 营养名称
    int threshold;              // 触发阈值（低于此值）
};

/**
 * @description: 疾病数据结构
 */
struct Disease
{
    std::string id;                    // 疾病ID
    std::string name;                  // 疾病名称
    std::string description;           // 描述
    DiseaseSeverity severity;          // 严重程度

    // 触发条件
    std::vector<DiseaseTrigger> triggerConditions;  // 营养触发条件
    int triggerDuration;               // 需要持续的行动次数
    float triggerChance;               // 触发概率 (0.0-1.0)

    // 效果
    float decayMultiplier;             // 营养衰减倍率
    float actionCostMultiplier;        // 行动消耗倍率

    // 治愈
    std::vector<DiseaseTrigger> cureConditions;  // 治愈所需营养条件
    int cureDuration;                  // 需要满足条件的行动次数
    int treatmentCost;                 // 医院治疗费用
};

/**
 * @description: 活动中的疾病状态
 */
struct ActiveDisease
{
    std::string id;           // 疾病ID
    int triggerProgress;      // 触发进度（累计满足条件的行动次数）
    int cureProgress;         // 康复进度
    bool isActive;            // 是否已激活（症状显现）
};

/**
 * @description: 疾病管理器（单例）
 */
class DiseaseManager
{
  public:
    /**
     * @description: 获取单例实例
     */
    static DiseaseManager *GetInstance();

    /**
     * @description: 每次行动后检查疾病
     */
    void OnAction();

    /**
     * @description: 触发食物中毒
     * @param severity 严重程度 (0.0-1.0)
     */
    void TriggerFoodPoisoning(float severity);

    /**
     * @description: 获取当前活动疾病列表
     */
    std::vector<std::string> GetActiveDiseases() const;

    /**
     * @description: 检查是否患有特定疾病
     */
    bool HasDisease(const std::string &diseaseId) const;

    /**
     * @description: 获取当前行动消耗倍率
     */
    float GetActionCostMultiplier() const;

    /**
     * @description: 获取指定营养的衰减倍率
     */
    float GetNutritionDecayMultiplier(const std::string &nutrition) const;

    /**
     * @description: 尝试自然康复
     * @return true=康复成功
     */
    bool TryNaturalCure(const std::string &diseaseId);

    /**
     * @description: 在医院治疗
     * @return true=治疗成功
     */
    bool TreatAtHospital(const std::string &diseaseId);

    /**
     * @description: 获取疾病信息
     */
    const Disease *GetDiseaseInfo(const std::string &diseaseId) const;

    /**
     * @description: 获取所有疾病状态摘要
     */
    std::string GetStatusSummary() const;

    /**
     * @description: 序列化疾病状态
     */
    std::string Serialize() const;

    /**
     * @description: 反序列化疾病状态
     */
    void Deserialize(const std::string &data);

  private:
    DiseaseManager();
    static DiseaseManager *instance;

    // 疾病定义表
    std::map<std::string, Disease> diseaseDefinitions;

    // 当前疾病状态
    std::map<std::string, ActiveDisease> activeDiseases;

    /**
     * @description: 初始化疾病定义
     */
    void InitDiseases();

    /**
     * @description: 检查触发条件
     */
    void CheckTriggerConditions();

    /**
     * @description: 检查康复条件
     */
    void CheckCureConditions();

    /**
     * @description: 检查单个疾病的触发条件
     */
    bool CheckSingleTrigger(const Disease &disease) const;

    /**
     * @description: 检查单个疾病的康复条件
     */
    bool CheckSingleCure(const Disease &disease) const;
};

/**
 * @description: 将疾病严重程度转换为字符串
 */
std::string SeverityToString(DiseaseSeverity severity);
