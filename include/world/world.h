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
    World();
    ~World();
    static World *getInstance();
    void start();
    void changewhere(int where);
    bool SpendMoney(int money);
    bool DepositingMoney(int money);
    bool WithdrawMoney(int money);
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