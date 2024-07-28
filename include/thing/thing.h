/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 08:56:57
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-28 10:46:19
 * @FilePath: \CityLive\include\thing\thing.h
 */
#pragma once
#ifndef THING_H
#define THING_H

#include <functional>
#include <map>
#include <string>

class World;
class Center;
typedef std::function<void(World *)> VoidFunction;
class Thing {
  public:
    Thing(std::string name) : name{name} {}
    virtual ~Thing() {}
    virtual std::string GetInfo() { return name; }
    virtual void ToDoThing(World *world) {}
    std::string name;
};
class GoWhere : public Thing {
  public:
    GoWhere(std::string name, int from, int to, int hungry) : Thing(name), from{from}, to{to}, hungry{hungry} {}
    void ToDoThing(World *world) override;

  private:
    int from;
    int to;
    int hungry;
};
class Buy : public Thing {
  public:
    Buy(std::string name) : Thing(name) {}
    void ToDoThing(World *world) override;

  private:
};
class Information : public Thing {
  public:
    Information(std::string name, std::string content);
    void ToDoThing(World *world) override;

  private:
    std::string content;
    std::map<std::string, VoidFunction> things;
};
class DepositingMoney : public Thing {
  public:
    DepositingMoney(std::string name) : Thing(name) {}
    void ToDoThing(World *world) override;

  private:
};
class WithdrawMoney : public Thing {
  public:
    WithdrawMoney(std::string name) : Thing(name) {}
    void ToDoThing(World *world) override;

  private:
};
#endif