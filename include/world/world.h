/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:39:31
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-28 10:02:27
 * @FilePath: \CityLive\include\world\world.h
 */
#pragma once
#ifndef WORLD_H
#define WORLD_H

#include <iostream>
#include <string>
#include <vector>
typedef std::vector<class Thing *> vthing;
class building;
class Bank;
class Controller;
class Thing;
class View;
class World {
  public:
    World();
    ~World();
    static World *getInstance();
    void start();
    void changewhere(int where);
    void DepositingMoney(int money);
    void WithdrawMoney(int money);
    int GetLen(building a, building b);

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
    std::vector<std::string> BuildingNames;
    std::vector<building> buildings;
    std::vector<vthing> ToDoThings;
};
class building {
  public:
    building() {}
    building(int id, int x, int y, std::string name) : id{id}, x{x}, y{y}, name{name} {}
    int x;
    int y;
    int id;
    std::string name;
};
#endif // WORLD_H