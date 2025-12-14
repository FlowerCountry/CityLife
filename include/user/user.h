#pragma once

#include <map>
#include <string>
#include <vector>

class Health;
class User {
  public:
    /**
     * @description: 用户类构造函数
     * @return null
     */
    User();

    /**
     * @description: 用户类析构函数
     * @return null
     */
    ~User();

    /**
     * @description: 返回用户类单例对象
     * @return {*} 用户类单例对象
     */
    static User *GetInstance();

    /**
     * @description: 返回用户的营养信息
     * @return {*} 用户的营养信息
     */
    std::map<std::string, int> *GetUserHealth() { return &UserHealth; }

    /**
     * @description: 获取指定营养值
     * @param name 营养名称
     * @return 营养值（可为负数）
     */
    int GetNutrition(const std::string &name) const;

    /**
     * @description: 设置指定营养值（支持负数）
     * @param name 营养名称
     * @param value 营养值
     */
    void SetNutrition(const std::string &name, int value);

    /**
     * @description: 消耗营养（减少指定数值）
     * @param name 营养名称
     * @param amount 消耗量
     */
    void ConsumeNutrition(const std::string &name, int amount);

    /**
     * @description: 增加营养（考虑上限100）
     * @param name 营养名称
     * @param amount 增加量
     */
    void AddNutrition(const std::string &name, int amount);

    /**
     * @description: 获取所有营养名称
     * @return 营养名称列表
     */
    std::vector<std::string> GetAllNutritionNames() const;

    /**
     * @description: 检查玩家是否存活
     * @return true=存活, false=死亡
     */
    bool IsAlive() const;

    /**
     * @description: 获取死亡风险百分比
     * @return 死亡风险 (0.0-1.0)
     */
    float GetDeathRisk() const;

    /**
     * @description: 获取营养状态等级
     * @param name 营养名称
     * @return 0=安全, 1=警告, 2=危险, 3=透支, 4=濒死
     */
    int GetNutritionLevel(const std::string &name) const;

  private:
    static User *instance;
    std::map<std::string, int> UserHealth;

    /**
     * @description: 获取营养的负数下限
     * @param name 营养名称
     * @return 负数下限值
     */
    int GetMinValue(const std::string &name) const;
};