/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:39:31
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 21:59:12
 * @FilePath: \CityLife\include\world\world.h
 */
#pragma once

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
    static World *getInstance();
    /**
     * @description: 游戏开始
     * @return null
     */
    void start();
    /**
     * @description: 更改位置
     * @param {int} where 目标位置
     * @return null
     */
    void changewhere(int where);
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
     * @description: 计算两地距离
     * @param {Building} a 起始地
     * @param {Building} b 目标地
     * @return {*} 距离
     */
    int GetLen(Building a, Building b);

  private:
    static World *instance;
    int year, month, day, hour, minute, second;
    int money;
    int where;
    float LifeQuality = 1.0;
    Bank *bank;
    std::vector<std::string> BuildingNames;
    std::vector<Building *> Buildings;
    std::vector<vobject> ToDoThings;
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