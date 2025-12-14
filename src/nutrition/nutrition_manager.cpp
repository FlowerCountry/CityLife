#include "nutrition/nutrition_manager.h"
#include "locale/locale_manager.h"
#include "user/user.h"

#include <algorithm>
#include <cstdlib>
#include <ctime>
#include <sstream>

NutritionManager *NutritionManager::instance = nullptr;

NutritionManager *NutritionManager::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new NutritionManager();
    }
    return instance;
}

NutritionManager::NutritionManager()
{
    // 初始化随机数种子
    std::srand(static_cast<unsigned int>(std::time(nullptr)));
    InitConfigs();
    InitInteractions();
    InitHintConfigs();
    InitPreviousLevels();
}

void NutritionManager::InitConfigs()
{
    // 核心属性（Primary）- 衰减快，影响大
    // 游戏平衡：从100降到0约需12-20次行动
    configs["饱腹感"] = {"饱腹感", NutritionCategory::Primary, 6.0f, -30, 100};
    configs["饥渴"] = {"饥渴", NutritionCategory::Primary, 5.0f, -20, 100};
    configs["蛋白质"] = {"蛋白质", NutritionCategory::Primary, 2.5f, -15, 100};
    configs["碳水化合物"] = {"碳水化合物", NutritionCategory::Primary, 3.0f, -15, 100};

    // 微量元素（Micro）- 衰减慢，长期影响
    // 游戏平衡：从100降到0约需50-200次行动
    configs["钙"] = {"钙", NutritionCategory::Micro, 0.5f, -10, 100};
    configs["糖分"] = {"糖分", NutritionCategory::Special, 2.0f, -10, 100};
    configs["脂肪"] = {"脂肪", NutritionCategory::Special, 0.5f, 0, 100};  // 脂肪不可为负
    configs["纤维素"] = {"纤维素", NutritionCategory::Micro, 1.0f, -10, 100};
    configs["铁"] = {"铁", NutritionCategory::Micro, 0.5f, -10, 100};
    configs["维生素A"] = {"维生素A", NutritionCategory::Micro, 0.5f, -10, 100};
    configs["维生素B"] = {"维生素B", NutritionCategory::Micro, 0.5f, -10, 100};
    configs["维生素C"] = {"维生素C", NutritionCategory::Micro, 1.0f, -10, 100};
    configs["维生素D"] = {"维生素D", NutritionCategory::Micro, 0.3f, -10, 100};
    configs["维生素E"] = {"维生素E", NutritionCategory::Micro, 0.3f, -10, 100};
    configs["钾"] = {"钾", NutritionCategory::Micro, 0.5f, -10, 100};
    configs["硒"] = {"硒", NutritionCategory::Micro, 0.3f, -10, 100};
    configs["锌"] = {"锌", NutritionCategory::Micro, 0.3f, -10, 100};

    // 特殊属性（Special）
    configs["幸福感"] = {"幸福感", NutritionCategory::Special, 1.5f, -20, 100};
    configs["精神振奋"] = {"精神振奋", NutritionCategory::Special, 3.0f, -10, 100};
    configs["益生菌"] = {"益生菌", NutritionCategory::Micro, 0.5f, -10, 100};
    configs["镁"] = {"镁", NutritionCategory::Micro, 0.3f, -10, 100};

    // 兼容旧的"饥饿"属性（等同于饥渴）
    configs["饥饿"] = {"饥饿", NutritionCategory::Primary, 5.0f, -20, 100};
}

void NutritionManager::InitInteractions()
{
    // 铁不足影响蛋白质吸收
    interactions.push_back({"铁", "蛋白质", 0.3f, 30, true});

    // 维生素D不足影响钙吸收
    interactions.push_back({"维生素D", "钙", 0.4f, 30, true});

    // 维生素B不足影响碳水化合物代谢
    interactions.push_back({"维生素B", "碳水化合物", 0.2f, 30, true});

    // 糖分充足临时提升幸福感（正向影响）
    interactions.push_back({"糖分", "幸福感", 0.5f, 50, false});

    // 脂肪作为能量缓冲（当碳水不足时）
    interactions.push_back({"脂肪", "碳水化合物", 0.3f, 20, true});
}

void NutritionManager::OnAction(float actionCost)
{
    User *user = User::GetInstance();

    // 遍历所有配置的营养进行衰减
    for (const auto &pair : configs)
    {
        const std::string &name = pair.first;
        const NutritionConfig &config = pair.second;

        // 计算衰减倍率（负数区间时加速衰减）
        float multiplier = GetDecayMultiplier(name);

        // 计算实际衰减量
        float decayAmount = config.baseDecay * actionCost * multiplier;

        // 应用衰减
        int current = user->GetNutrition(name);
        int newValue = current - static_cast<int>(decayAmount);

        // 限制在最小值以上
        user->SetNutrition(name, std::max(newValue, config.minValue));
    }
}

void NutritionManager::ProcessInteractions()
{
    User *user = User::GetInstance();

    for (const auto &interaction : interactions)
    {
        int sourceValue = user->GetNutrition(interaction.source);

        if (interaction.isPositive)
        {
            // 源营养不足时，影响目标营养的衰减
            if (sourceValue < interaction.threshold)
            {
                // 计算惩罚：源营养越低，目标营养衰减越快
                float penalty = static_cast<float>(interaction.threshold - sourceValue) / 100.0f * interaction.ratio;
                int targetValue = user->GetNutrition(interaction.target);
                int extraDecay = static_cast<int>(targetValue * penalty * 0.1f);
                if (extraDecay > 0)
                {
                    user->ConsumeNutrition(interaction.target, extraDecay);
                }
            }
        }
        else
        {
            // 源营养充足时，正向影响目标
            if (sourceValue > interaction.threshold)
            {
                // 小幅提升目标营养
                int boost = static_cast<int>((sourceValue - interaction.threshold) * interaction.ratio * 0.05f);
                if (boost > 0)
                {
                    user->AddNutrition(interaction.target, boost);
                }
            }
        }
    }
}

const NutritionConfig *NutritionManager::GetConfig(const std::string &name) const
{
    auto it = configs.find(name);
    if (it != configs.end())
    {
        return &it->second;
    }
    return nullptr;
}

float NutritionManager::GetDecayMultiplier(const std::string &name) const
{
    User *user = User::GetInstance();
    int value = user->GetNutrition(name);

    // 根据营养值区间返回衰减倍率
    if (value > 20)
    {
        return 1.0f;  // 正常衰减
    }
    else if (value > 0)
    {
        return 1.2f;  // 危险区：衰减加速 20%
    }
    else if (value > -10)
    {
        return 1.5f;  // 透支区：衰减加速 50%
    }
    else
    {
        return 2.0f;  // 濒死区：衰减加速 100%
    }
}

std::vector<std::string> NutritionManager::GetWarnings() const
{
    std::vector<std::string> warnings;
    User *user = User::GetInstance();
    LocaleManager *locale = LocaleManager::GetInstance();

    for (const auto &pair : configs)
    {
        const std::string &name = pair.first;
        int level = user->GetNutritionLevel(name);

        if (level == 1)  // 警告区
        {
            warnings.push_back(name + locale->Get("nutrition.level", "low"));
        }
    }
    return warnings;
}

std::vector<std::string> NutritionManager::GetCriticals() const
{
    std::vector<std::string> criticals;
    User *user = User::GetInstance();
    LocaleManager *locale = LocaleManager::GetInstance();

    for (const auto &pair : configs)
    {
        const std::string &name = pair.first;
        int level = user->GetNutritionLevel(name);

        if (level >= 2)  // 危险区及以上
        {
            std::string severity;
            switch (level)
            {
            case 2:
                severity = locale->Get("nutrition.level", "danger");
                break;
            case 3:
                severity = locale->Get("nutrition.level", "overdraft");
                break;
            case 4:
                severity = locale->Get("nutrition.level", "critical");
                break;
            default:
                severity = locale->Get("nutrition.level", "abnormal");
            }
            criticals.push_back(name + severity);
        }
    }
    return criticals;
}

std::string NutritionManager::GetStatusSummary() const
{
    User *user = User::GetInstance();
    std::ostringstream ss;

    // 只显示核心属性的状态
    ss << "饱腹:" << user->GetNutrition("饱腹感");
    ss << " 饥渴:" << user->GetNutrition("饥渴");
    ss << " 蛋白:" << user->GetNutrition("蛋白质");
    ss << " 碳水:" << user->GetNutrition("碳水化合物");

    return ss.str();
}

bool NutritionManager::IsAlive() const
{
    return User::GetInstance()->IsAlive();
}

float NutritionManager::GetDeathRisk() const
{
    return User::GetInstance()->GetDeathRisk();
}

bool NutritionManager::CheckRandomDeath()
{
    float risk = GetDeathRisk();
    if (risk <= 0.0f)
    {
        return false;
    }

    // 生成 0.0 到 1.0 之间的随机数
    float roll = static_cast<float>(std::rand()) / static_cast<float>(RAND_MAX);
    return roll < risk;
}

void NutritionManager::InitHintConfigs()
{
    LocaleManager *locale = LocaleManager::GetInstance();

    // 饱腹感提示
    hintMessages[{"饱腹感", 1}] = locale->Get("nutrition.warning.hunger", "level1");
    hintMessages[{"饱腹感", 2}] = locale->Get("nutrition.warning.hunger", "level2");
    hintMessages[{"饱腹感", 3}] = locale->Get("nutrition.warning.hunger", "level3");
    hintMessages[{"饱腹感", 4}] = locale->Get("nutrition.warning.hunger", "level4");

    // 饥渴提示
    hintMessages[{"饥渴", 1}] = locale->Get("nutrition.warning.thirst", "level1");
    hintMessages[{"饥渴", 2}] = locale->Get("nutrition.warning.thirst", "level2");
    hintMessages[{"饥渴", 3}] = locale->Get("nutrition.warning.thirst", "level3");
    hintMessages[{"饥渴", 4}] = locale->Get("nutrition.warning.thirst", "level4");

    // 蛋白质提示
    hintMessages[{"蛋白质", 1}] = locale->Get("nutrition.warning.protein", "level1");
    hintMessages[{"蛋白质", 2}] = locale->Get("nutrition.warning.protein", "level2");
    hintMessages[{"蛋白质", 3}] = locale->Get("nutrition.warning.protein", "level3");
    hintMessages[{"蛋白质", 4}] = locale->Get("nutrition.warning.protein", "level4");

    // 碳水化合物提示
    hintMessages[{"碳水化合物", 1}] = locale->Get("nutrition.warning.carbs", "level1");
    hintMessages[{"碳水化合物", 2}] = locale->Get("nutrition.warning.carbs", "level2");
    hintMessages[{"碳水化合物", 3}] = locale->Get("nutrition.warning.carbs", "level3");
    hintMessages[{"碳水化合物", 4}] = locale->Get("nutrition.warning.carbs", "level4");

    // 幸福感提示
    hintMessages[{"幸福感", 1}] = locale->Get("nutrition.warning.happiness", "level1");
    hintMessages[{"幸福感", 2}] = locale->Get("nutrition.warning.happiness", "level2");
    hintMessages[{"幸福感", 3}] = locale->Get("nutrition.warning.happiness", "level3");
    hintMessages[{"幸福感", 4}] = locale->Get("nutrition.warning.happiness", "level4");

    // 精神振奋提示
    hintMessages[{"精神振奋", 1}] = locale->Get("nutrition.warning.energy", "level1");
    hintMessages[{"精神振奋", 2}] = locale->Get("nutrition.warning.energy", "level2");
    hintMessages[{"精神振奋", 3}] = locale->Get("nutrition.warning.energy", "level3");
    hintMessages[{"精神振奋", 4}] = locale->Get("nutrition.warning.energy", "level4");

    // 维生素C提示
    hintMessages[{"维生素C", 1}] = locale->Get("nutrition.warning.vitaminc", "level1");
    hintMessages[{"维生素C", 2}] = locale->Get("nutrition.warning.vitaminc", "level2");

    // 铁提示
    hintMessages[{"铁", 1}] = locale->Get("nutrition.warning.iron", "level1");
    hintMessages[{"铁", 2}] = locale->Get("nutrition.warning.iron", "level2");

    // 钙提示
    hintMessages[{"钙", 1}] = locale->Get("nutrition.warning.calcium", "level1");
    hintMessages[{"钙", 2}] = locale->Get("nutrition.warning.calcium", "level2");
}

void NutritionManager::InitPreviousLevels()
{
    User *user = User::GetInstance();

    // 初始化所有营养的等级记录
    for (const auto &pair : configs)
    {
        const std::string &name = pair.first;
        previousLevels[name] = user->GetNutritionLevel(name);
    }
}

void NutritionManager::CheckLevelChanges()
{
    User *user = User::GetInstance();

    for (const auto &pair : configs)
    {
        const std::string &name = pair.first;
        int currentLevel = user->GetNutritionLevel(name);
        int previousLevel = previousLevels[name];

        // 等级数值增大 = 状态恶化（0=安全, 1=警告, 2=危险, 3=透支, 4=濒死）
        if (currentLevel > previousLevel)
        {
            // 查找对应等级的提示
            auto it = hintMessages.find({name, currentLevel});
            if (it != hintMessages.end())
            {
                pendingHints.push_back(it->second);
            }
        }

        // 更新记录的等级
        previousLevels[name] = currentLevel;
    }

    // 限制最多显示3条提示，避免屏幕过于拥挤
    if (pendingHints.size() > 3)
    {
        pendingHints.resize(3);
    }
}

std::vector<std::string> NutritionManager::GetPendingHints() const
{
    return pendingHints;
}

void NutritionManager::ClearPendingHints()
{
    pendingHints.clear();
}

std::string NutritionManager::GetFeelingDescription(const std::string &name, int value) const
{
    LocaleManager *locale = LocaleManager::GetInstance();

    // 饱腹感
    if (name == "饱腹感" || name == "饥饿")
    {
        if (value >= 90) return locale->Get("nutrition.feedback.hunger", "high");
        if (value >= 70) return locale->Get("nutrition.feedback.hunger", "medium");
        if (value >= 50) return locale->Get("nutrition.feedback.hunger", "low");
        return locale->Get("nutrition.feedback.hunger", "minimal");
    }

    // 饥渴
    if (name == "饥渴")
    {
        if (value >= 90) return locale->Get("nutrition.feedback.thirst", "high");
        if (value >= 70) return locale->Get("nutrition.feedback.thirst", "medium");
        if (value >= 50) return locale->Get("nutrition.feedback.thirst", "low");
        return locale->Get("nutrition.feedback.thirst", "minimal");
    }

    // 蛋白质
    if (name == "蛋白质")
    {
        if (value >= 80) return locale->Get("nutrition.feedback.protein", "high");
        if (value >= 60) return locale->Get("nutrition.feedback.protein", "medium");
        return locale->Get("nutrition.feedback.protein", "low");
    }

    // 碳水化合物
    if (name == "碳水化合物")
    {
        if (value >= 80) return locale->Get("nutrition.feedback.carbs", "high");
        if (value >= 60) return locale->Get("nutrition.feedback.carbs", "medium");
        return locale->Get("nutrition.feedback.carbs", "low");
    }

    // 幸福感
    if (name == "幸福感")
    {
        if (value >= 80) return locale->Get("nutrition.feedback.happiness", "high");
        if (value >= 60) return locale->Get("nutrition.feedback.happiness", "medium");
        return locale->Get("nutrition.feedback.happiness", "low");
    }

    // 精神振奋
    if (name == "精神振奋")
    {
        if (value >= 80) return locale->Get("nutrition.feedback.energy", "high");
        if (value >= 60) return locale->Get("nutrition.feedback.energy", "medium");
        return locale->Get("nutrition.feedback.energy", "low");
    }

    // 维生素类
    if (name.find("维生素") != std::string::npos)
    {
        if (value >= 80) return locale->Get("nutrition.feedback.vitamin", "high");
        return locale->Get("nutrition.feedback.vitamin", "low");
    }

    // 矿物质类（钙、铁、锌等）
    if (name == "钙" || name == "铁" || name == "锌" || name == "钾" ||
        name == "硒" || name == "镁")
    {
        if (value >= 80) return locale->Get("nutrition.feedback.mineral", "high");
        return locale->Get("nutrition.feedback.mineral", "low");
    }

    // 其他（糖分、脂肪、纤维素、益生菌）
    if (name == "糖分")
    {
        if (value >= 70) return locale->Get("nutrition.feedback.sugar", "high");
        return locale->Get("nutrition.feedback.sugar", "low");
    }

    if (name == "纤维素")
    {
        if (value >= 70) return locale->Get("nutrition.feedback.fiber", "high");
        return locale->Get("nutrition.feedback.fiber", "low");
    }

    if (name == "益生菌")
    {
        if (value >= 70) return locale->Get("nutrition.feedback.probiotic", "high");
        return locale->Get("nutrition.feedback.probiotic", "low");
    }

    // 默认
    return locale->Get("nutrition.feedback.default", "msg");
}
