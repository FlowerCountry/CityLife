#pragma once

#include <array>
#include <functional>
#include <map>
#include <string>
#include <vector>

class World;
class Center;
class Health;
typedef std::function<void()> VoidFunction;
enum class EventCategory
{
    Primary = 0,
    Insight,
    Navigation
};
class Object {
  public:
    /**
     * @description: 事件基类初始化
     * @param {string} name 事件名称
     * @return null
     */
    Object(const std::string &name) : name{name} {}
    /**
     * @description: 事件基类虚构函数
     * @return null
     */
    virtual ~Object() {}
    /**
     * @description: 获取事件信息
     * @return null
     */
    virtual std::string GetInfo() { return name; }
    /**
     * @description: 去做当前事件
     * @return null
     */
    virtual void ToDoIt() {}
    /**
     * @description: 返回事件分类
     * @return {*} 事件分类
     */
    virtual EventCategory GetCategory() const { return EventCategory::Insight; }
    std::string name;
};
class GoWhere : public Object {
  public:
    /**
     * @description: 去别地事件构造函数
     * @param {string} name 事件名称
     * @param {int} from 来源地
     * @param {int} to 目标地
     * @param {int} hungry 消耗的饥饿值
     * @return null
     */
    GoWhere(const std::string &name, const int &from, const int &to, const int &hungry) : Object(name), from{from}, to{to}, hungry{hungry} {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;

  private:
    int from;
    int to;
    int hungry;
};
class Buy : public Object {
  public:
    /**
     * @description: 购买物品事件构造函数
     * @param {string} name 所购买的物品名称
     * @return null
     */
    Buy(const std::string &name) : Object(name) {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;

  private:
};
class Information : public Object {
  public:
    /**
     * @description: 公告事件构造函数
     * @param {string} name 公告名称
     * @param {string} content 公告内容
     * @return null
     */
    Information(const std::string &name, const std::string &content);
    void ToDoIt() override;
    EventCategory GetCategory() const override;

  private:
    std::string content;
    std::map<std::string, VoidFunction> things;
};
class DepositingMoney : public Object {
  public:
    /**
     * @description: 存钱事件构造函数
     * @return null
     */
    DepositingMoney() : Object("存钱") {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;

  private:
};
class WithdrawMoney : public Object {
  public:
    /**
     * @description: 取钱事件构造函数
     * @return null
     */
    WithdrawMoney() : Object("取钱") {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;

  private:
};
class CheckCash : public Object {
  public:
    /**
     * @description: 查看随身现金事件构造函数
     * @return null
     */
    CheckCash() : Object("查看身上现金") {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;

  private:
};
class CheckBankBalance : public Object {
  public:
    /**
     * @description: 查看银行余额事件构造函数
     * @return null
     */
    CheckBankBalance() : Object("查看账户余额") {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;

  private:
};
class Commodity : public Object {
  public:
    /**
     * @description: 商品构造函数
     * @param {string} name 商品名称
     * @param {int} price 商品价格
     * @param {vector<Health>} health 商品所含营养
     * @return null
     */
    Commodity(const std::string &name, const int &price, const std::vector<Health> &health);
    void ToDoIt() override;
    std::string GetInfo() override { return name + " " + std::to_string(price) + "$"; }
    EventCategory GetCategory() const override;

    /**
     * @description: 返回商品所含营养
     * @return {*} 商品所含营养
     */
    std::vector<Health> GetHealth();

  private:
    int price;
    std::vector<Health> health;
};

class SaveGame : public Object {
  public:
    /**
     * @description: 保存游戏事件构造函数
     * @return null
     */
    SaveGame() : Object("保存游戏") {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;
};

class LoadGame : public Object {
  public:
    /**
     * @description: 读取存档事件构造函数
     * @return null
     */
    LoadGame() : Object("读取存档") {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;
};

class SeeDoctor : public Object {
  public:
    /**
     * @description: 看病事件构造函数
     * @return null
     */
    SeeDoctor() : Object("看病") {}
    void ToDoIt() override;
    EventCategory GetCategory() const override;
};
