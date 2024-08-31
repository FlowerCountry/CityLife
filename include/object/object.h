/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 08:56:57
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-08-30 19:50:34
 * @FilePath: \CityLife\include\object\object.h
 */
#pragma once
#include "health/health.h"
#ifndef THING_H
#define THING_H

#include <functional>
#include <map>
#include <string>
#include <vector>

class World;
class Center;
class Health;
typedef std::function<void(World *)> VoidFunction;
class Object {
  public:
    Object(std::string name) : name{name} {}
    virtual ~Object() {}
    virtual std::string GetInfo() { return name; }
    virtual void ToDoIt(World *world) {}
    std::string name;
};
class GoWhere : public Object {
  public:
    GoWhere(std::string name, int from, int to, int hungry) : Object(name), from{from}, to{to}, hungry{hungry} {}
    void ToDoIt(World *world) override;

  private:
    int from;
    int to;
    int hungry;
};
class Buy : public Object {
  public:
    Buy(std::string name) : Object(name) {}
    void ToDoIt(World *world) override;

  private:
};
class Information : public Object {
  public:
    Information(std::string name, std::string content);
    void ToDoIt(World *world) override;

  private:
    std::string content;
    std::map<std::string, VoidFunction> things;
};
class DepositingMoney : public Object {
  public:
    DepositingMoney(std::string name) : Object(name) {}
    void ToDoIt(World *world) override;

  private:
};
class WithdrawMoney : public Object {
  public:
    WithdrawMoney(std::string name) : Object(name) {}
    void ToDoIt(World *world) override;

  private:
};

class Commodity : public Object {
  public:
    Commodity(std::string name, int price, std::vector<Health *> health) : Object{name}, price{price}, health{health} {}
    void ToDoIt(World *world) override;
    std::string GetInfo() override { return name + " " + std::to_string(price) + "$"; }
    int price;
    std::vector<Health *> health;
};
#endif