/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:32:47
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-08-30 21:14:52
 * @FilePath: \CityLife\include\controller\controller.h
 */
#pragma once
#ifndef CONTROLLER_H
#define CONTROLLER_H

#include <string>
#include <vector>
class Object;
class View;
class World;
class Controller {
  public:
    Controller();
    ~Controller();
    static Controller *getInstance();
    int choose(std::string str, std::vector<class Object *> options);

  private:
    static Controller *instance;
};

#endif // CONTROLLER_H