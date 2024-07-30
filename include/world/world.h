/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:39:31
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-30 08:34:53
 * @FilePath: \CityLife\include\world\world.h
 */
#pragma once
#ifndef WORLD_H
#define WORLD_H

#include <string>
#include <vector>
typedef std::vector<class Thing *> vthing;
class Building;
class Commodity;
class Bank;
class Controller;
class Thing;
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
     * 2表示银行
     */
    int year, month, day, hour, minute, second;
    int money;
    int where;
    Bank *bank;
    std::vector<Commodity *> Commodities;
    std::vector<std::string> BuildingNames;
    std::vector<Building> Buildings;
    std::vector<vthing> ToDoThings;
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

class Commodity {
  public:
    Commodity() {}
    Commodity(std::string name, int price, Health *health) : name{name}, price{price}, health{health} {}
    std::string name;
    int price;
    Health *health;
};
#endif // WORLD_H