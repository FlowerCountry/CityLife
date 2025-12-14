#include "disease/disease_manager.h"
#include "locale/locale_manager.h"
#include "user/user.h"
#include "world/world.h"

#include <algorithm>
#include <cstdlib>
#include <ctime>
#include <sstream>

DiseaseManager *DiseaseManager::instance = nullptr;

DiseaseManager *DiseaseManager::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new DiseaseManager();
    }
    return instance;
}

DiseaseManager::DiseaseManager()
{
    std::srand(static_cast<unsigned int>(std::time(nullptr)));
    InitDiseases();
}

void DiseaseManager::InitDiseases()
{
    LocaleManager *locale = LocaleManager::GetInstance();

    // 感冒 - 维生素C不足
    Disease cold;
    cold.id = "cold";
    cold.name = locale->Get("disease.name", "cold");
    cold.description = locale->Get("disease.description", "cold");
    cold.severity = DiseaseSeverity::Mild;
    cold.triggerConditions = {{"维生素C", 20}};
    cold.triggerDuration = 10;
    cold.triggerChance = 0.15f;
    cold.decayMultiplier = 1.2f;
    cold.actionCostMultiplier = 1.2f;
    cold.cureConditions = {{"维生素C", 60}};
    cold.cureDuration = 5;
    cold.treatmentCost = 50;
    diseaseDefinitions["cold"] = cold;

    // 贫血 - 铁不足
    Disease anemia;
    anemia.id = "anemia";
    anemia.name = locale->Get("disease.name", "anemia");
    anemia.description = locale->Get("disease.description", "anemia");
    anemia.severity = DiseaseSeverity::Moderate;
    anemia.triggerConditions = {{"铁", 15}};
    anemia.triggerDuration = 15;
    anemia.triggerChance = 0.2f;
    anemia.decayMultiplier = 1.5f;
    anemia.actionCostMultiplier = 1.5f;
    anemia.cureConditions = {{"铁", 70}, {"蛋白质", 60}};
    anemia.cureDuration = 10;
    anemia.treatmentCost = 200;
    diseaseDefinitions["anemia"] = anemia;

    // 坏血病 - 维生素C严重缺乏
    Disease scurvy;
    scurvy.id = "scurvy";
    scurvy.name = locale->Get("disease.name", "scurvy");
    scurvy.description = locale->Get("disease.description", "scurvy");
    scurvy.severity = DiseaseSeverity::Severe;
    scurvy.triggerConditions = {{"维生素C", -10}};
    scurvy.triggerDuration = 5;
    scurvy.triggerChance = 0.4f;
    scurvy.decayMultiplier = 2.0f;
    scurvy.actionCostMultiplier = 2.0f;
    scurvy.cureConditions = {{"维生素C", 80}};
    scurvy.cureDuration = 15;
    scurvy.treatmentCost = 500;
    diseaseDefinitions["scurvy"] = scurvy;

    // 食物中毒 - 特殊触发
    Disease foodPoisoning;
    foodPoisoning.id = "food_poisoning";
    foodPoisoning.name = locale->Get("disease.name", "food_poisoning");
    foodPoisoning.description = locale->Get("disease.description", "food_poisoning");
    foodPoisoning.severity = DiseaseSeverity::Moderate;
    foodPoisoning.triggerConditions = {};  // 特殊触发，不由营养阈值触发
    foodPoisoning.triggerDuration = 1;
    foodPoisoning.triggerChance = 1.0f;
    foodPoisoning.decayMultiplier = 2.0f;
    foodPoisoning.actionCostMultiplier = 1.8f;
    foodPoisoning.cureConditions = {};  // 自然康复
    foodPoisoning.cureDuration = 8;
    foodPoisoning.treatmentCost = 100;
    diseaseDefinitions["food_poisoning"] = foodPoisoning;

    // 营养不良 - 蛋白质不足
    Disease malnutrition;
    malnutrition.id = "malnutrition";
    malnutrition.name = locale->Get("disease.name", "malnutrition");
    malnutrition.description = locale->Get("disease.description", "malnutrition");
    malnutrition.severity = DiseaseSeverity::Moderate;
    malnutrition.triggerConditions = {{"蛋白质", 20}};
    malnutrition.triggerDuration = 12;
    malnutrition.triggerChance = 0.2f;
    malnutrition.decayMultiplier = 1.4f;
    malnutrition.actionCostMultiplier = 1.4f;
    malnutrition.cureConditions = {{"蛋白质", 70}, {"碳水化合物", 50}};
    malnutrition.cureDuration = 8;
    malnutrition.treatmentCost = 150;
    diseaseDefinitions["malnutrition"] = malnutrition;

    // 抑郁症 - 幸福感长期低
    Disease depression;
    depression.id = "depression";
    depression.name = locale->Get("disease.name", "depression");
    depression.description = locale->Get("disease.description", "depression");
    depression.severity = DiseaseSeverity::Moderate;
    depression.triggerConditions = {{"幸福感", 10}};
    depression.triggerDuration = 20;
    depression.triggerChance = 0.15f;
    depression.decayMultiplier = 1.3f;
    depression.actionCostMultiplier = 1.3f;
    depression.cureConditions = {{"幸福感", 60}};
    depression.cureDuration = 15;
    depression.treatmentCost = 300;
    diseaseDefinitions["depression"] = depression;
}

void DiseaseManager::OnAction()
{
    CheckTriggerConditions();
    CheckCureConditions();
}

void DiseaseManager::CheckTriggerConditions()
{
    for (const auto &pair : diseaseDefinitions)
    {
        const Disease &disease = pair.second;

        // 跳过特殊触发的疾病（如食物中毒）
        if (disease.triggerConditions.empty())
        {
            continue;
        }

        // 检查是否已经患病
        auto it = activeDiseases.find(disease.id);
        if (it != activeDiseases.end() && it->second.isActive)
        {
            continue;  // 已经激活，跳过
        }

        // 检查触发条件
        if (CheckSingleTrigger(disease))
        {
            // 满足条件，增加触发进度
            if (it == activeDiseases.end())
            {
                ActiveDisease ad;
                ad.id = disease.id;
                ad.triggerProgress = 1;
                ad.cureProgress = 0;
                ad.isActive = false;
                activeDiseases[disease.id] = ad;
            }
            else
            {
                it->second.triggerProgress++;
            }

            // 检查是否达到触发阈值
            auto &ad = activeDiseases[disease.id];
            if (ad.triggerProgress >= disease.triggerDuration)
            {
                // 随机检查是否触发
                float roll = static_cast<float>(std::rand()) / static_cast<float>(RAND_MAX);
                if (roll < disease.triggerChance)
                {
                    ad.isActive = true;
                }
            }
        }
        else
        {
            // 不满足条件，重置触发进度
            if (it != activeDiseases.end() && !it->second.isActive)
            {
                it->second.triggerProgress = 0;
            }
        }
    }
}

void DiseaseManager::CheckCureConditions()
{
    for (auto &pair : activeDiseases)
    {
        ActiveDisease &ad = pair.second;
        if (!ad.isActive)
        {
            continue;
        }

        const Disease *disease = GetDiseaseInfo(ad.id);
        if (disease == nullptr)
        {
            continue;
        }

        // 食物中毒自然康复（不需要条件）
        if (disease->cureConditions.empty())
        {
            ad.cureProgress++;
            if (ad.cureProgress >= disease->cureDuration)
            {
                ad.isActive = false;
                ad.cureProgress = 0;
                ad.triggerProgress = 0;
            }
            continue;
        }

        // 检查康复条件
        if (CheckSingleCure(*disease))
        {
            ad.cureProgress++;
            if (ad.cureProgress >= disease->cureDuration)
            {
                ad.isActive = false;
                ad.cureProgress = 0;
                ad.triggerProgress = 0;
            }
        }
        else
        {
            ad.cureProgress = 0;  // 条件不满足，重置康复进度
        }
    }
}

bool DiseaseManager::CheckSingleTrigger(const Disease &disease) const
{
    User *user = User::GetInstance();
    for (const auto &trigger : disease.triggerConditions)
    {
        int value = user->GetNutrition(trigger.nutritionName);
        if (value >= trigger.threshold)
        {
            return false;  // 任一条件不满足
        }
    }
    return true;
}

bool DiseaseManager::CheckSingleCure(const Disease &disease) const
{
    User *user = User::GetInstance();
    for (const auto &cure : disease.cureConditions)
    {
        int value = user->GetNutrition(cure.nutritionName);
        if (value < cure.threshold)
        {
            return false;  // 任一条件不满足
        }
    }
    return true;
}

void DiseaseManager::TriggerFoodPoisoning(float severity)
{
    float roll = static_cast<float>(std::rand()) / static_cast<float>(RAND_MAX);
    if (roll < severity)
    {
        ActiveDisease ad;
        ad.id = "food_poisoning";
        ad.triggerProgress = 0;
        ad.cureProgress = 0;
        ad.isActive = true;
        activeDiseases["food_poisoning"] = ad;
    }
}

std::vector<std::string> DiseaseManager::GetActiveDiseases() const
{
    std::vector<std::string> result;
    for (const auto &pair : activeDiseases)
    {
        if (pair.second.isActive)
        {
            const Disease *disease = GetDiseaseInfo(pair.first);
            if (disease != nullptr)
            {
                result.push_back(disease->name);
            }
        }
    }
    return result;
}

bool DiseaseManager::HasDisease(const std::string &diseaseId) const
{
    auto it = activeDiseases.find(diseaseId);
    return it != activeDiseases.end() && it->second.isActive;
}

float DiseaseManager::GetActionCostMultiplier() const
{
    float multiplier = 1.0f;
    for (const auto &pair : activeDiseases)
    {
        if (pair.second.isActive)
        {
            const Disease *disease = GetDiseaseInfo(pair.first);
            if (disease != nullptr)
            {
                multiplier *= disease->actionCostMultiplier;
            }
        }
    }
    return multiplier;
}

float DiseaseManager::GetNutritionDecayMultiplier(const std::string &nutrition) const
{
    (void)nutrition;  // 当前所有疾病对所有营养使用相同倍率
    float multiplier = 1.0f;
    for (const auto &pair : activeDiseases)
    {
        if (pair.second.isActive)
        {
            const Disease *disease = GetDiseaseInfo(pair.first);
            if (disease != nullptr)
            {
                multiplier *= disease->decayMultiplier;
            }
        }
    }
    return multiplier;
}

bool DiseaseManager::TryNaturalCure(const std::string &diseaseId)
{
    auto it = activeDiseases.find(diseaseId);
    if (it == activeDiseases.end() || !it->second.isActive)
    {
        return false;
    }

    const Disease *disease = GetDiseaseInfo(diseaseId);
    if (disease == nullptr)
    {
        return false;
    }

    // 检查康复条件
    if (CheckSingleCure(*disease))
    {
        it->second.cureProgress++;
        if (it->second.cureProgress >= disease->cureDuration)
        {
            it->second.isActive = false;
            it->second.cureProgress = 0;
            it->second.triggerProgress = 0;
            return true;
        }
    }
    return false;
}

bool DiseaseManager::TreatAtHospital(const std::string &diseaseId)
{
    auto it = activeDiseases.find(diseaseId);
    if (it == activeDiseases.end() || !it->second.isActive)
    {
        return false;
    }

    const Disease *disease = GetDiseaseInfo(diseaseId);
    if (disease == nullptr)
    {
        return false;
    }

    // 检查是否有足够的钱
    World *world = World::GetInstance();
    if (world->GetWallet() < disease->treatmentCost)
    {
        return false;
    }

    // 扣款
    if (!world->SpendMoney(disease->treatmentCost))
    {
        return false;
    }

    // 立即治愈
    it->second.isActive = false;
    it->second.cureProgress = 0;
    it->second.triggerProgress = 0;
    return true;
}

const Disease *DiseaseManager::GetDiseaseInfo(const std::string &diseaseId) const
{
    auto it = diseaseDefinitions.find(diseaseId);
    if (it != diseaseDefinitions.end())
    {
        return &it->second;
    }
    return nullptr;
}

std::string DiseaseManager::GetStatusSummary() const
{
    auto diseases = GetActiveDiseases();
    if (diseases.empty())
    {
        LocaleManager *locale = LocaleManager::GetInstance();
        return locale->Get("disease.status", "healthy");
    }

    std::ostringstream ss;
    for (std::size_t i = 0; i < diseases.size(); ++i)
    {
        if (i > 0)
        {
            ss << ", ";
        }
        ss << diseases[i];
    }
    return ss.str();
}

std::string DiseaseManager::Serialize() const
{
    std::ostringstream ss;
    for (const auto &pair : activeDiseases)
    {
        const ActiveDisease &ad = pair.second;
        ss << ad.id << ":" << ad.triggerProgress << ":"
           << ad.cureProgress << ":" << (ad.isActive ? 1 : 0) << ";";
    }
    return ss.str();
}

void DiseaseManager::Deserialize(const std::string &data)
{
    activeDiseases.clear();
    if (data.empty())
    {
        return;
    }

    std::istringstream iss(data);
    std::string token;
    while (std::getline(iss, token, ';'))
    {
        if (token.empty())
        {
            continue;
        }

        std::istringstream tokenStream(token);
        std::string id;
        int triggerProgress, cureProgress, isActive;

        std::getline(tokenStream, id, ':');
        tokenStream >> triggerProgress;
        tokenStream.ignore(1);  // skip ':'
        tokenStream >> cureProgress;
        tokenStream.ignore(1);
        tokenStream >> isActive;

        ActiveDisease ad;
        ad.id = id;
        ad.triggerProgress = triggerProgress;
        ad.cureProgress = cureProgress;
        ad.isActive = (isActive != 0);
        activeDiseases[id] = ad;
    }
}

std::string SeverityToString(DiseaseSeverity severity)
{
    LocaleManager *locale = LocaleManager::GetInstance();
    switch (severity)
    {
    case DiseaseSeverity::Mild:
        return locale->Get("disease.severity", "mild");
    case DiseaseSeverity::Moderate:
        return locale->Get("disease.severity", "moderate");
    case DiseaseSeverity::Severe:
        return locale->Get("disease.severity", "severe");
    default:
        return locale->Get("disease.severity", "unknown");
    }
}
