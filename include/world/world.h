/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:39:31
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-08-30 19:30:26
 * @FilePath: \CityLife\include\world\world.h
 */
#pragma once
#ifndef WORLD_H
#define WORLD_H

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
    World();
    ~World();
    static World *getInstance();
    void start();
    void changewhere(int where);
    void DepositingMoney(int money);
    void WithdrawMoney(int money);
    int GetLen(Building a, Building b);
    void buy();

  private:
    static World *instance;
    /**
     * 0表示市中央
     * 1表示超市
     * 2表示银行f
     */
    int year, month, day, hour, minute, second;
    int money;
    int where;
    float LifeQuality = 1.0;
    Bank *bank;
    std::vector<Object *> Commodities;
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

#endif // WORLD_H