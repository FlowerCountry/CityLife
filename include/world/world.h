#pragma once

#include <array>
#include <string>
#include <vector>
typedef std::vector<class Object *> vobject;
class Building;
class Commodity;
class Bank;
class Controller;
class Object;
class View;
class Health;
class World {
  public:
    /**
     * @description: 世界类构造函数
     * @return null
     */
    World();
    /**
     * @description: 世界类析构函数
     * @return null
     */
    ~World();
    /**
     * @description: 返回世界类单例对象
     * @return {*} 世界类单例对象
     */
    static World *GetInstance();
    /**
     * @description: 游戏开始
     * @return null
     */
    void Start();
    /**
     * @description: 更改位置
     * @param {int} where 目标位置
     * @return null
     */
    void ChangeWhere(int where);
    /**
     * @description: 花钱
     * @param {int} money 花的钱数
     * @return {*} 是否成功
     */
    bool SpendMoney(int money);
    /**
     * @description: 存钱
     * @param {int} money 存的钱数
     * @return {*} 是否成功
     */
    bool DepositingMoney(int money);
    /**
     * @description: 取钱
     * @param {int} money 取得钱数
     * @return {*} 是否成功
     */
    bool WithdrawMoney(int money);
    /**
     * @description: 获取玩家随身现金
     * @return {*} 随身现金
     */
    int GetWallet() const;
    const std::array<int, 6> &GetWalletBreakdown() const;
    bool RemoveBills(const std::array<int, 6> &counts);
    void AddBills(const std::array<int, 6> &counts);
    /**
     * @description: 获取玩家银行余额
     * @return {*} 银行余额
     */
    int GetBankBalance() const;
    /**
     * @description: 计算两地距离
     * @param {Building} a 起始地
     * @param {Building} b 目标地
     * @return {*} 距离
     */
    int GetLen(Building a, Building b);
    /**
     * @description: 更新游戏时间
     * @param {int} seconds 增加的秒数
     * @return null
     */
    void UpdateTime(int seconds);

    // ===== 存档系统 getter =====
    /**
     * @description: 获取当前年份
     * @return {*} 年份
     */
    int GetYear() const;
    /**
     * @description: 获取当前月份
     * @return {*} 月份
     */
    int GetMonth() const;
    /**
     * @description: 获取当前日期
     * @return {*} 日期
     */
    int GetDay() const;
    /**
     * @description: 获取当前小时
     * @return {*} 小时
     */
    int GetHour() const;
    /**
     * @description: 获取当前分钟
     * @return {*} 分钟
     */
    int GetMinute() const;
    /**
     * @description: 获取当前秒数
     * @return {*} 秒
     */
    int GetSecond() const;
    /**
     * @description: 获取当前位置
     * @return {*} 位置ID
     */
    int GetWhere() const;
    /**
     * @description: 获取生活质量系数
     * @return {*} 生活质量
     */
    float GetLifeQuality() const;

    // ===== 存档系统 setter =====
    /**
     * @description: 设置时间（用于存档加载）
     * @param {int} year 年份
     * @param {int} month 月份
     * @param {int} day 日期
     * @param {int} hour 小时
     * @param {int} minute 分钟
     * @param {int} second 秒
     * @return null
     */
    void SetTime(int year, int month, int day, int hour, int minute, int second);
    /**
     * @description: 设置位置（用于存档加载）
     * @param {int} where 位置ID
     * @return null
     */
    void SetWhere(int where);
    /**
     * @description: 设置生活质量系数（用于存档加载）
     * @param {float} quality 生活质量
     * @return null
     */
    void SetLifeQuality(float quality);
    /**
     * @description: 设置钱包面额分布（用于存档加载）
     * @param {array} wallet 面额分布数组
     * @return null
     */
    void SetWallet(const std::array<int, 6> &wallet);
    /**
     * @description: 设置银行存款（用于存档加载）
     * @param {int} amount 存款金额
     * @return null
     */
    void SetBankDeposit(int amount);

  private:
    static World *instance;
    int year, month, day, hour, minute, second;
    int money;
    std::array<int, 6> wallet;
    static const std::array<int, 6> denominations;
    int where;
    float LifeQuality = 1.0;
    Bank *bank;
    std::vector<std::string> BuildingNames;
    std::vector<Building *> Buildings;
    std::vector<vobject> ToDoThings;

    bool RemoveAmount(int amount);
    static int SumWallet(const std::array<int, 6> &counts);
    void RefreshWalletTotal();
};
class Building {
  public:
    Building() {}
    Building(int id, int x, int y, std::string name) : id{id}, x{x}, y{y}, name{name} {}
    int x;
    int y;
    int id;
    std::string name;
};
